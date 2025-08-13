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

const (
	refreshRetryCount = 3
)

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

// ArkClient is a struct that represents a client for the Ark service.
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

// MarshalCookies marshals the cookies into a JSON byte array.
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

// UnmarshalCookies unmarshals the JSON byte array into the cookies map.
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

// NewSimpleArkClient creates a new instance of ArkClient with the specified base URL.
func NewSimpleArkClient(baseURL string) *ArkClient {
	return NewArkClient(baseURL, "", "", "", nil, nil)
}

// NewArkClient creates a new instance of ArkClient with the specified parameters.
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

// SetHeader sets a single header for the ArkClient.
func (ac *ArkClient) SetHeader(key string, value string) {
	ac.headers[key] = value
}

// SetHeaders sets multiple headers for the ArkClient.
func (ac *ArkClient) SetHeaders(headers map[string]string) {
	ac.headers = headers
}

// UpdateHeaders updates the headers for the ArkClient with the provided headers.
func (ac *ArkClient) UpdateHeaders(headers map[string]string) {
	for key, value := range headers {
		ac.headers[key] = value
	}
}

// GetHeaders returns the headers for the ArkClient.
func (ac *ArkClient) GetHeaders() map[string]string {
	return ac.headers
}

// SetCookie sets a single cookie for the ArkClient.
func (ac *ArkClient) SetCookie(key string, value string) {
	ac.cookieJar.SetCookies(
		&url.URL{
			Scheme: "https",
			Host:   ac.BaseURL,
		},
		[]*http.Cookie{
			{
				Name:  key,
				Value: value,
			},
		},
	)
}

// SetCookies sets multiple cookies for the ArkClient.
func (ac *ArkClient) SetCookies(cookies map[string]string) {
	ac.cookieJar.RemoveAll()
	for key, value := range cookies {
		ac.SetCookie(key, value)
	}
}

// UpdateCookies updates the cookies for the ArkClient with the provided cookies.
func (ac *ArkClient) UpdateCookies(cookies map[string]string) {
	for key, value := range cookies {
		ac.SetCookie(key, value)
	}
}

// GetCookies returns the cookies for the ArkClient.
func (ac *ArkClient) GetCookies() map[string]string {
	cookies := make(map[string]string)
	for _, cookie := range ac.cookieJar.AllCookies() {
		cookies[cookie.Name] = cookie.Value
	}
	return cookies
}

// GetCookieJar returns the cookie jar for the ArkClient.
func (ac *ArkClient) GetCookieJar() *cookiejar.Jar {
	return ac.cookieJar
}

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

// Get performs a GET request to the specified route with the provided parameters.
func (ac *ArkClient) Get(ctx context.Context, route string, params map[string]string) (*http.Response, error) {
	return ac.doRequest(ctx, http.MethodGet, route, map[string]string{}, params, refreshRetryCount)
}

// Post performs a POST request to the specified route with the provided body.
func (ac *ArkClient) Post(ctx context.Context, route string, body interface{}) (*http.Response, error) {
	return ac.doRequest(ctx, http.MethodPost, route, body, nil, refreshRetryCount)
}

// Put performs a PUT request to the specified route with the provided body.
func (ac *ArkClient) Put(ctx context.Context, route string, body interface{}) (*http.Response, error) {
	return ac.doRequest(ctx, http.MethodPut, route, body, nil, refreshRetryCount)
}

// Delete performs a DELETE request to the specified route.
func (ac *ArkClient) Delete(ctx context.Context, route string, body interface{}) (*http.Response, error) {
	return ac.doRequest(ctx, http.MethodDelete, route, body, nil, refreshRetryCount)
}

// Patch performs a PATCH request to the specified route with the provided body.
func (ac *ArkClient) Patch(ctx context.Context, route string, body interface{}) (*http.Response, error) {
	return ac.doRequest(ctx, http.MethodPatch, route, body, nil, refreshRetryCount)
}

// Options performs an OPTIONS request to the specified route.
func (ac *ArkClient) Options(ctx context.Context, route string) (*http.Response, error) {
	return ac.doRequest(ctx, http.MethodOptions, route, nil, nil, refreshRetryCount)
}

// UpdateToken updates the token and token type for the ArkClient.
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

// GetToken returns the token for the ArkClient.
func (ac *ArkClient) GetToken() string {
	return ac.token
}

// GetTokenType returns the token type for the ArkClient.
func (ac *ArkClient) GetTokenType() string {
	return ac.tokenType
}
