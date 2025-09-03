// Package common provides shared utilities and types for the ARK SDK.
//
// This package implements a comprehensive HTTP client with features like:
// - Authentication support (token-based, basic auth)
// - Cookie management with persistent storage
// - Automatic token refresh capabilities
// - Request/response logging
// - TLS configuration options
//
// The ArkClient is the primary interface for making HTTP requests to Ark services,
// providing a consistent and feature-rich HTTP client implementation.
package common

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	cookiejar "github.com/juju/persistent-cookiejar"
)

// Number of retry attempts for token refresh operations.
const (
	refreshRetryCount = 3
)

// cookieJSON represents the JSON serializable format of an HTTP cookie.
//
// This structure is used for marshaling and unmarshaling HTTP cookies
// to and from JSON format, enabling persistent cookie storage and
// session management across client instances.
type cookieJSON struct {
	Name        string        `json:"name"`
	Value       string        `json:"value"`
	Quoted      bool          `json:"quoted"`
	Path        string        `json:"path,omitempty"`
	Domain      string        `json:"domain,omitempty"`
	Expires     time.Time     `json:"expires,omitempty"`
	RawExpires  string        `json:"raw_expires,omitempty"`
	MaxAge      int           `json:"max_age,omitempty"`
	Secure      bool          `json:"secure,omitempty"`
	HTTPOnly    bool          `json:"http_only,omitempty"`
	SameSite    http.SameSite `json:"same_site,omitempty"`
	Partitioned bool          `json:"partitioned,omitempty"`
	Raw         string        `json:"raw,omitempty"`
	Unparsed    []string      `json:"unparsed,omitempty"`
}

// ArkClient provides a comprehensive HTTP client for interacting with Ark services.
//
// ArkClient wraps the standard Go HTTP client with additional features specifically
// designed for Ark service interactions. It handles authentication, cookie management,
// request logging, and automatic token refresh capabilities.
//
// Key features:
// - Token-based and basic authentication support
// - Persistent cookie storage with JSON serialization
// - Automatic retry with token refresh on 401 responses
// - Configurable headers for all requests
// - Request/response logging with timing information
// - TLS configuration support
//
// The client maintains state including authentication tokens, custom headers,
// and cookie storage, making it suitable for session-based interactions
// with Ark services.
type ArkClient struct {
	BaseURL                   string
	token                     string
	tokenType                 string
	authHeaderName            string
	client                    *http.Client
	headers                   map[string]string
	cookieJar                 *cookiejar.Jar
	refreshConnectionCallback func(*ArkClient) error
	logger                    *ArkLogger
}

// MarshalCookies serializes a cookie jar into a JSON byte array.
//
// This function converts all cookies from a cookiejar.Jar into a JSON-serializable
// format, enabling persistent storage of cookie state. The resulting byte array
// can be stored to disk or transmitted over networks for session persistence.
//
// Note: This implementation uses the AllCookies() method from persistent-cookiejar
// which provides direct access to all stored cookies.
//
// Parameters:
//   - cookieJar: The cookie jar containing cookies to be marshaled
//
// Returns the JSON byte array representation of all cookies, or an error
// if JSON marshaling fails.
//
// Example:
//
//	cookieData, err := MarshalCookies(client.GetCookieJar())
//	if err != nil {
//	    // handle error
//	}
//	// Save cookieData to file or database
func MarshalCookies(cookieJar *cookiejar.Jar) ([]byte, error) {
	jsonCookies := make([]cookieJSON, len(cookieJar.AllCookies()))
	for i, c := range cookieJar.AllCookies() {
		jsonCookies[i] = cookieJSON{
			Name:        c.Name,
			Value:       c.Value,
			Quoted:      c.Quoted,
			Path:        c.Path,
			Domain:      c.Domain,
			Expires:     c.Expires,
			RawExpires:  c.RawExpires,
			MaxAge:      c.MaxAge,
			Secure:      c.Secure,
			HTTPOnly:    c.HttpOnly,
			SameSite:    c.SameSite,
			Partitioned: c.Partitioned,
			Raw:         c.Raw,
			Unparsed:    c.Unparsed,
		}
	}
	cookiesBytes, err := json.Marshal(jsonCookies)
	if err != nil {
		return nil, err
	}
	return cookiesBytes, nil
}

// UnmarshalCookies deserializes a JSON byte array into a cookie jar.
//
// This function takes a JSON byte array (typically created by MarshalCookies)
// and populates the provided cookie jar with the deserialized cookies.
// The cookies are organized by URL and properly set in the jar for use
// in subsequent HTTP requests.
//
// Parameters:
//   - cookies: JSON byte array containing serialized cookie data
//   - cookieJar: The cookie jar to populate with deserialized cookies
//
// Returns an error if JSON unmarshaling fails or if URL parsing encounters
// invalid cookie data.
//
// Example:
//
//	err := UnmarshalCookies(savedCookieData, client.GetCookieJar())
//	if err != nil {
//	    // handle error
//	}
//	// Cookie jar now contains restored cookies
func UnmarshalCookies(cookies []byte, cookieJar *cookiejar.Jar) error {
	var jsonCookies []cookieJSON
	if err := json.Unmarshal(cookies, &jsonCookies); err != nil {
		return err
	}
	allCookies := make([]*http.Cookie, len(jsonCookies))
	for i, c := range jsonCookies {
		allCookies[i] = &http.Cookie{
			Name:        c.Name,
			Value:       c.Value,
			Quoted:      c.Quoted,
			Path:        c.Path,
			Domain:      c.Domain,
			Expires:     c.Expires,
			RawExpires:  c.RawExpires,
			MaxAge:      c.MaxAge,
			Secure:      c.Secure,
			HttpOnly:    c.HTTPOnly,
			SameSite:    c.SameSite,
			Partitioned: c.Partitioned,
			Raw:         c.Raw,
			Unparsed:    c.Unparsed,
		}
	}
	cookieGroups := make(map[string][]*http.Cookie)
	for _, cookie := range allCookies {
		urlKey := fmt.Sprintf("https://%s%s", cookie.Domain, cookie.Path)
		cookieGroups[urlKey] = append(cookieGroups[urlKey], cookie)
	}
	for urlKey, cookiesGroup := range cookieGroups {
		parsedURL, err := url.Parse(urlKey)
		if err != nil {
			return fmt.Errorf("failed to parse URL %s: %w", urlKey, err)
		}
		cookieJar.SetCookies(parsedURL, cookiesGroup)
	}
	return nil
}

// NewSimpleArkClient creates a basic ArkClient instance with minimal configuration.
//
// This is a convenience constructor for creating an ArkClient with only a base URL.
// It uses default values for all other parameters (no authentication, new cookie jar,
// no refresh callback). This is suitable for simple use cases or as a starting point
// for further configuration.
//
// Parameters:
//   - baseURL: The base URL for the Ark service (HTTPS prefix will be added if missing)
//
// Returns a configured ArkClient instance ready for basic HTTP operations.
//
// Example:
//
//	client := NewSimpleArkClient("api.example.com")
//	response, err := client.Get(ctx, "/users", nil)
func NewSimpleArkClient(baseURL string) *ArkClient {
	return NewArkClient(baseURL, "", "", "", nil, nil)
}

// NewArkClient creates a new ArkClient instance with comprehensive configuration options.
//
// This is the primary constructor for ArkClient, allowing full customization of
// authentication, cookie management, and refresh behavior. The client will automatically
// add HTTPS prefix to the base URL if not present and initialize a new cookie jar
// if none is provided.
//
// Parameters:
//   - baseURL: The base URL for the Ark service
//   - token: Authentication token (empty string for no authentication)
//   - tokenType: Type of token ("Bearer", "Basic", etc.)
//   - authHeaderName: Name of the authorization header (e.g., "Authorization")
//   - cookieJar: Cookie jar for session management (nil for new jar)
//   - refreshCallback: Function to call for token refresh on 401 responses (nil to disable)
//
// Returns a fully configured ArkClient instance.
//
// Example:
//
//	jar, _ := cookiejar.New(nil)
//	client := NewArkClient(
//	    "https://api.example.com",
//	    "abc123",
//	    "Bearer",
//	    "Authorization",
//	    jar,
//	    func(c *ArkClient) error {
//	        // Token refresh logic
//	        return nil
//	    },
//	)
func NewArkClient(baseURL string, token string, tokenType string, authHeaderName string, cookieJar *cookiejar.Jar, refreshCallback func(*ArkClient) error) *ArkClient {
	if baseURL != "" && !strings.HasPrefix(baseURL, "https://") {
		baseURL = "https://" + baseURL
	}
	if cookieJar == nil {
		cookieJar, _ = cookiejar.New(nil)
	}
	client := &ArkClient{
		BaseURL:        baseURL,
		authHeaderName: authHeaderName,
		cookieJar:      cookieJar,
		client: &http.Client{
			Jar: cookieJar,
		},
		headers:                   make(map[string]string),
		refreshConnectionCallback: refreshCallback,
		logger:                    GetLogger("ArkClient", Unknown),
	}
	client.UpdateToken(token, tokenType)
	client.headers["User-Agent"] = UserAgent()
	return client
}

// SetHeader sets a single HTTP header for the ArkClient.
//
// This method adds or updates a single header in the client's header map.
// The header will be included in all subsequent HTTP requests made by this client.
// If a header with the same key already exists, it will be overwritten.
//
// Parameters:
//   - key: The header name (e.g., "Content-Type", "Accept")
//   - value: The header value (e.g., "application/json", "text/plain")
//
// Example:
//
//	client.SetHeader("Content-Type", "application/json")
//	client.SetHeader("Accept", "application/json")
func (ac *ArkClient) SetHeader(key string, value string) {
	ac.headers[key] = value
}

// SetHeaders replaces all existing headers with the provided header map.
//
// This method completely replaces the client's header map with the new headers.
// Any previously set headers will be lost. Use UpdateHeaders() if you want to
// preserve existing headers while adding new ones.
//
// Parameters:
//   - headers: Map of header names to values that will replace all existing headers
//
// Example:
//
//	headers := map[string]string{
//	    "Content-Type": "application/json",
//	    "Accept": "application/json",
//	}
//	client.SetHeaders(headers)
func (ac *ArkClient) SetHeaders(headers map[string]string) {
	ac.headers = headers
}

// UpdateHeaders merges the provided headers into the existing header map.
//
// This method adds new headers or updates existing ones while preserving
// headers that are not specified in the input map. If a header key already
// exists, its value will be overwritten.
//
// Parameters:
//   - headers: Map of header names to values to add or update
//
// Example:
//
//	newHeaders := map[string]string{
//	    "X-Custom-Header": "custom-value",
//	    "Authorization": "Bearer new-token",
//	}
//	client.UpdateHeaders(newHeaders)
func (ac *ArkClient) UpdateHeaders(headers map[string]string) {
	for key, value := range headers {
		ac.headers[key] = value
	}
}

// GetHeaders returns a copy of the current header map.
//
// This method returns the client's current headers. Note that modifying
// the returned map will not affect the client's headers - use SetHeader()
// or UpdateHeaders() to modify headers.
//
// Returns a map containing all current headers.
//
// Example:
//
//	currentHeaders := client.GetHeaders()
//	fmt.Printf("Content-Type: %s\n", currentHeaders["Content-Type"])
func (ac *ArkClient) GetHeaders() map[string]string {
	return ac.headers
}

// SetCookie sets a single cookie in the client's cookie jar.
//
// This method adds a new cookie to the client's cookie jar, which will be
// included in subsequent requests to the appropriate domain. The cookie
// is associated with the client's base URL.
//
// Parameters:
//   - key: The cookie name
//   - value: The cookie value
//
// Example:
//
//	client.SetCookie("session_id", "abc123")
//	client.SetCookie("user_pref", "dark_mode")
func (ac *ArkClient) SetCookie(key string, value string) {
	parsedURL, err := url.Parse(ac.BaseURL)
	if err != nil {
		ac.logger.Error("Fail to parse url %s: %v", ac.BaseURL, err)
		parsedURL = &url.URL{
			Scheme: "https",
			Host:   ac.BaseURL,
		}
	}
	ac.cookieJar.SetCookies(
		parsedURL,
		[]*http.Cookie{
			{
				Name:  key,
				Value: value,
			},
		},
	)
}

// SetCookies replaces all existing cookies with the provided cookie map.
//
// This method removes all existing cookies from the cookie jar and replaces
// them with the new cookies. Use UpdateCookies() if you want to preserve
// existing cookies while adding new ones.
//
// Parameters:
//   - cookies: Map of cookie names to values that will replace all existing cookies
//
// Example:
//
//	cookies := map[string]string{
//	    "session_id": "abc123",
//	    "csrf_token": "xyz789",
//	}
//	client.SetCookies(cookies)
func (ac *ArkClient) SetCookies(cookies map[string]string) {
	ac.cookieJar.RemoveAll()
	for key, value := range cookies {
		ac.SetCookie(key, value)
	}
}

// UpdateCookies adds or updates cookies in the existing cookie jar.
//
// This method adds new cookies or updates existing ones while preserving
// cookies that are not specified in the input map.
//
// Parameters:
//   - cookies: Map of cookie names to values to add or update
//
// Example:
//
//	newCookies := map[string]string{
//	    "new_session": "def456",
//	    "updated_pref": "light_mode",
//	}
//	client.UpdateCookies(newCookies)
func (ac *ArkClient) UpdateCookies(cookies map[string]string) {
	for key, value := range cookies {
		ac.SetCookie(key, value)
	}
}

// GetCookies returns a map of all current cookies.
//
// This method extracts all cookies from the cookie jar and returns them
// as a simple map of names to values. This is useful for inspecting
// current cookie state or for serialization purposes.
//
// Note: This implementation uses the AllCookies() method from persistent-cookiejar
// which provides direct access to all stored cookies.
//
// Returns a map containing all current cookie names and values.
//
// Example:
//
//	cookies := client.GetCookies()
//	sessionID := cookies["session_id"]
func (ac *ArkClient) GetCookies() map[string]string {
	cookies := make(map[string]string)
	for _, cookie := range ac.cookieJar.AllCookies() {
		cookies[cookie.Name] = cookie.Value
	}
	return cookies
}

// GetCookieJar returns the underlying cookie jar instance.
//
// This method provides direct access to the cookiejar.Jar for advanced
// cookie management operations that are not covered by the convenience
// methods. Use this when you need full control over cookie behavior.
//
// Returns the cookie jar instance used by this client.
//
// Example:
//
//	jar := client.GetCookieJar()
//	// Perform advanced cookie operations
//	cookieData, err := MarshalCookies(jar)
func (ac *ArkClient) GetCookieJar() *cookiejar.Jar {
	return ac.cookieJar
}

// doRequest is the internal method that handles the actual HTTP request execution.
//
// This method constructs and executes HTTP requests with comprehensive functionality
// including URL construction, JSON serialization, header application, query parameter
// handling, TLS configuration, request logging, and automatic token refresh on
// authentication failures.
//
// The method performs several key operations:
// - URL construction with proper path escaping
// - JSON marshaling of request body
// - Header and query parameter application
// - TLS configuration based on certificate verification settings
// - Request timing and logging
// - Automatic retry with token refresh on 401 responses
//
// Parameters:
//   - ctx: Context for request cancellation and timeout control
//   - method: HTTP method (GET, POST, PUT, DELETE, etc.)
//   - route: API route/path to append to the base URL
//   - body: Request body to be JSON-serialized (can be nil for methods like GET)
//   - params: Query parameters to include in the request URL (can be nil)
//   - refreshRetryCount: Number of retry attempts remaining for token refresh
//
// Returns the HTTP response or an error if the request fails or retry attempts
// are exhausted.
//
// The method automatically handles:
// - HTTPS URL construction with proper path segment escaping
// - JSON serialization of request bodies
// - Application of all configured headers
// - Query parameter encoding
// - TLS certificate verification based on global settings
// - Request/response timing logging
// - Token refresh retry logic on 401 Unauthorized responses
func (ac *ArkClient) doRequest(ctx context.Context, method string, route string, body interface{}, params map[string]string, refreshRetryCount int) (*http.Response, error) {
	fullURL := ac.BaseURL
	if route != "" {
		segments := strings.Split(route, "/")
		for i, segment := range segments {
			segments[i] = url.PathEscape(segment)
		}
		route = strings.Join(segments, "/")
		if fullURL[len(fullURL)-1] != '/' && route[0] != '/' {
			fullURL += "/"
		}
		fullURL += route
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, fullURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	for key, value := range ac.headers {
		req.Header.Set(key, value)
	}
	if params != nil {
		urlParams := url.Values{}
		for key, value := range params {
			urlParams.Add(key, value)
		}
		req.URL.RawQuery = urlParams.Encode()
	}
	if !IsVerifyingCertificates() {
		ac.client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	} else {
		ac.client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
		}
	}
	ac.logger.Info("Running request to %s", fullURL)
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime)
		ac.logger.Info("Request to %s took %dms", fullURL, duration.Milliseconds())
	}()
	resp, err := ac.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusUnauthorized && ac.refreshConnectionCallback != nil && refreshRetryCount > 0 {
		err = ac.refreshConnectionCallback(ac)
		if err != nil {
			return nil, err
		}
		return ac.doRequest(ctx, method, route, body, params, refreshRetryCount-1)
	}
	return resp, nil
}

// Get performs an HTTP GET request to the specified route.
//
// This method constructs and executes a GET request using the client's base URL,
// headers, and authentication. Query parameters can be provided via the params map.
// The method handles automatic token refresh on 401 responses if a refresh callback
// is configured.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout control
//   - route: API route/path to append to the base URL
//   - params: Query parameters to include in the request (nil for no parameters)
//
// Returns the HTTP response or an error if the request fails.
//
// Example:
//
//	params := map[string]string{"limit": "10", "offset": "0"}
//	response, err := client.Get(ctx, "/users", params)
//	if err != nil {
//	    // handle error
//	}
//	defer response.Body.Close()
func (ac *ArkClient) Get(ctx context.Context, route string, params map[string]string) (*http.Response, error) {
	return ac.doRequest(ctx, http.MethodGet, route, map[string]string{}, params, refreshRetryCount)
}

// Post performs an HTTP POST request to the specified route.
//
// This method constructs and executes a POST request with the provided body
// serialized as JSON. The request includes all configured headers and
// authentication. Automatic token refresh is handled on 401 responses.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout control
//   - route: API route/path to append to the base URL
//   - body: Request body data to be JSON-serialized
//
// Returns the HTTP response or an error if the request fails.
//
// Example:
//
//	userData := map[string]string{"name": "John", "email": "john@example.com"}
//	response, err := client.Post(ctx, "/users", userData)
//	if err != nil {
//	    // handle error
//	}
//	defer response.Body.Close()
func (ac *ArkClient) Post(ctx context.Context, route string, body interface{}) (*http.Response, error) {
	return ac.doRequest(ctx, http.MethodPost, route, body, nil, refreshRetryCount)
}

// Put performs an HTTP PUT request to the specified route.
//
// This method constructs and executes a PUT request with the provided body
// serialized as JSON. PUT requests are typically used for updating or
// replacing existing resources.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout control
//   - route: API route/path to append to the base URL
//   - body: Request body data to be JSON-serialized
//
// Returns the HTTP response or an error if the request fails.
//
// Example:
//
//	updatedUser := map[string]string{"name": "John Doe", "email": "john.doe@example.com"}
//	response, err := client.Put(ctx, "/users/123", updatedUser)
func (ac *ArkClient) Put(ctx context.Context, route string, body interface{}) (*http.Response, error) {
	return ac.doRequest(ctx, http.MethodPut, route, body, nil, refreshRetryCount)
}

// Delete performs an HTTP DELETE request to the specified route.
//
// This method constructs and executes a DELETE request. An optional body
// can be provided for DELETE requests that require additional data.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout control
//   - route: API route/path to append to the base URL
//   - body: Optional request body data to be JSON-serialized (can be nil)
//
// Returns the HTTP response or an error if the request fails.
//
// Example:
//
//	response, err := client.Delete(ctx, "/users/123", nil)
//	// Or with body:
//	deleteOptions := map[string]bool{"force": true}
//	response, err := client.Delete(ctx, "/users/123", deleteOptions)
func (ac *ArkClient) Delete(ctx context.Context, route string, body interface{}) (*http.Response, error) {
	return ac.doRequest(ctx, http.MethodDelete, route, body, nil, refreshRetryCount)
}

// Patch performs an HTTP PATCH request to the specified route.
//
// This method constructs and executes a PATCH request with the provided body
// serialized as JSON. PATCH requests are typically used for partial updates
// of existing resources.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout control
//   - route: API route/path to append to the base URL
//   - body: Request body data to be JSON-serialized
//
// Returns the HTTP response or an error if the request fails.
//
// Example:
//
//	partialUpdate := map[string]string{"email": "newemail@example.com"}
//	response, err := client.Patch(ctx, "/users/123", partialUpdate)
func (ac *ArkClient) Patch(ctx context.Context, route string, body interface{}) (*http.Response, error) {
	return ac.doRequest(ctx, http.MethodPatch, route, body, nil, refreshRetryCount)
}

// Options performs an HTTP OPTIONS request to the specified route.
//
// This method constructs and executes an OPTIONS request, typically used
// to retrieve information about the communication options available for
// the target resource or server.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout control
//   - route: API route/path to append to the base URL
//
// Returns the HTTP response or an error if the request fails.
//
// Example:
//
//	response, err := client.Options(ctx, "/users")
//	// Check response headers for allowed methods, CORS info, etc.
func (ac *ArkClient) Options(ctx context.Context, route string) (*http.Response, error) {
	return ac.doRequest(ctx, http.MethodOptions, route, nil, nil, refreshRetryCount)
}

// UpdateToken updates the authentication token and token type for the client.
//
// This method updates the client's authentication credentials and automatically
// configures the appropriate authorization header. It supports both standard
// token-based authentication and basic authentication. For basic auth, the token
// should be a base64-encoded "username:password" string.
//
// Parameters:
//   - token: The authentication token or base64-encoded credentials
//   - tokenType: The type of token ("Bearer", "Basic", "API-Key", etc.)
//
// The method will automatically set the Authorization header based on the token type:
// - For "Basic" type: Decodes the token and sets "Authorization: Basic <credentials>"
// - For other types: Sets the configured auth header with format "<tokenType> <token>"
//
// Example:
//
//	// Bearer token
//	client.UpdateToken("abc123xyz", "Bearer")
//
//	// Basic auth (token should be base64 encoded "user:pass")
//	client.UpdateToken("dXNlcjpwYXNz", "Basic")
//
//	// API key
//	client.UpdateToken("api-key-value", "API-Key")
func (ac *ArkClient) UpdateToken(token string, tokenType string) {
	ac.token = token
	ac.tokenType = tokenType
	if token != "" {
		if tokenType == "Basic" {
			decoded, err := base64.StdEncoding.DecodeString(token)
			if err != nil {
				return
			}
			userPass := string(decoded)
			ac.headers["Authorization"] = "Basic " + userPass
		} else {
			ac.headers[ac.authHeaderName] = fmt.Sprintf("%s %s", tokenType, token)
		}
	}
}

// GetToken returns the current authentication token.
//
// This method returns the raw token string that was set via UpdateToken().
// For basic authentication, this will be the base64-encoded credentials.
//
// Returns the current authentication token string.
//
// Example:
//
//	currentToken := client.GetToken()
//	if currentToken == "" {
//	    // No authentication token is set
//	}
func (ac *ArkClient) GetToken() string {
	return ac.token
}

// GetTokenType returns the current token type.
//
// This method returns the token type that was set via UpdateToken(),
// such as "Bearer", "Basic", "API-Key", etc.
//
// Returns the current token type string.
//
// Example:
//
//	tokenType := client.GetTokenType()
//	fmt.Printf("Using %s authentication\n", tokenType)
func (ac *ArkClient) GetTokenType() string {
	return ac.tokenType
}
