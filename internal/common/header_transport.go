package common

import "net/http"

// HeaderTransport is a custom HTTP transport that adds headers to requests.
type HeaderTransport struct {
	Transport http.RoundTripper
	Headers   map[string]string
}

// RoundTrip NewHeaderTransport creates a new HeaderTransport with the specified headers.
func (t *HeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for key, value := range t.Headers {
		req.Header.Set(key, value)
	}
	return t.Transport.RoundTrip(req)
}

// BasicAuthTransport is a custom HTTP transport that adds Basic Authentication to requests.
type BasicAuthTransport struct {
	Transport http.RoundTripper
	Username  string
	Password  string
}

// RoundTrip NewBasicAuthTransport creates a new BasicAuthTransport with the specified username and password.
func (t *BasicAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(t.Username, t.Password)
	return t.Transport.RoundTrip(req)
}
