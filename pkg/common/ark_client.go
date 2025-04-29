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
)

const (
	refreshRetryCount = 3
)

// ArkClient is a struct that represents a client for the Ark service.
type ArkClient struct {
	BaseURL                   string
	token                     string
	tokenType                 string
	authHeaderName            string
	client                    *http.Client
	headers                   map[string]string
	cookies                   map[string]string
	refreshConnectionCallback func(*ArkClient) error
	logger                    *ArkLogger
}

// NewSimpleArkClient creates a new instance of ArkClient with the specified base URL.
func NewSimpleArkClient(baseURL string) *ArkClient {
	return NewArkClient(baseURL, "", "", "", nil)
}

// NewArkClient creates a new instance of ArkClient with the specified parameters.
func NewArkClient(baseURL string, token string, tokenType string, authHeaderName string, refreshCallback func(*ArkClient) error) *ArkClient {
	if baseURL != "" && !strings.HasPrefix(baseURL, "https://") {
		baseURL = "https://" + baseURL
	}
	client := &ArkClient{
		BaseURL:                   baseURL,
		authHeaderName:            authHeaderName,
		client:                    &http.Client{},
		headers:                   make(map[string]string),
		cookies:                   make(map[string]string),
		refreshConnectionCallback: refreshCallback,
		logger:                    GetLogger("ArkClient", Unknown),
	}
	client.UpdateToken(token, tokenType)
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
	ac.cookies[key] = value
}

// SetCookies sets multiple cookies for the ArkClient.
func (ac *ArkClient) SetCookies(cookies map[string]string) {
	ac.cookies = cookies
}

// UpdateCookies updates the cookies for the ArkClient with the provided cookies.
func (ac *ArkClient) UpdateCookies(cookies map[string]string) {
	for key, value := range cookies {
		ac.cookies[key] = value
	}
}

// GetCookies returns the cookies for the ArkClient.
func (ac *ArkClient) GetCookies() map[string]string {
	return ac.cookies
}

// MarshalCookies marshals the cookies into a JSON byte array.
func (ac *ArkClient) MarshalCookies() ([]byte, error) {
	cookiesBytes, err := json.Marshal(ac.cookies)
	if err != nil {
		return nil, err
	}
	return cookiesBytes, nil
}

// UnmarshalCookies unmarshals the JSON byte array into the cookies map.
func (ac *ArkClient) UnmarshalCookies(cookies []byte) error {
	return json.Unmarshal(cookies, &ac.cookies)
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
	for key, value := range ac.cookies {
		req.AddCookie(&http.Cookie{Name: key, Value: value})
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
	if ac.client.Jar != nil {
		u, _ := url.Parse(ac.BaseURL)
		cookies := ac.client.Jar.Cookies(u)
		for _, cookie := range cookies {
			ac.cookies[cookie.Name] = cookie.Value
		}
	}
	return resp, nil
}

// Get performs a GET request to the specified route with the provided parameters.
func (ac *ArkClient) Get(ctx context.Context, route string, params map[string]string) (*http.Response, error) {
	return ac.doRequest(ctx, http.MethodGet, route, nil, params, refreshRetryCount)
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
