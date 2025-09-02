package common

import "net/http"

// HeaderTransport is a custom HTTP transport that adds headers to requests.
//
// HeaderTransport wraps an existing http.RoundTripper and automatically adds
// specified headers to every HTTP request. This is useful for adding common
// headers like User-Agent, Content-Type, or custom API headers to all requests
// made through the transport.
//
// The transport preserves the behavior of the underlying transport while
// ensuring that all specified headers are set on outgoing requests. If a
// header already exists in the request, it will be overwritten with the
// value from the Headers map.
//
// Example:
//
//	transport := &HeaderTransport{
//	    Transport: http.DefaultTransport,
//	    Headers: map[string]string{
//	        "User-Agent": "MyApp/1.0",
//	        "X-API-Key": "secret-key",
//	    },
//	}
type HeaderTransport struct {
	// Transport is the underlying http.RoundTripper to wrap
	Transport http.RoundTripper
	// Headers is a map of header names to values that will be added to every request
	Headers map[string]string
}

// RoundTrip implements the http.RoundTripper interface and adds headers to the request.
//
// RoundTrip iterates through all headers in the Headers map and sets them on the
// provided request using req.Header.Set(). This overwrites any existing headers
// with the same name. After setting all headers, it delegates to the underlying
// Transport's RoundTrip method to perform the actual HTTP request.
//
// Parameters:
//   - req: The HTTP request to modify and execute
//
// Returns the HTTP response from the underlying transport and any error encountered.
//
// Example:
//
//	transport := &HeaderTransport{
//	    Transport: http.DefaultTransport,
//	    Headers: map[string]string{"User-Agent": "MyApp/1.0"},
//	}
//	resp, err := transport.RoundTrip(req)
func (t *HeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for key, value := range t.Headers {
		req.Header.Set(key, value)
	}
	return t.Transport.RoundTrip(req)
}

// BasicAuthTransport is a custom HTTP transport that adds Basic Authentication to requests.
//
// BasicAuthTransport wraps an existing http.RoundTripper and automatically adds
// HTTP Basic Authentication headers to every HTTP request using the provided
// username and password credentials. This is useful for APIs that require
// Basic Authentication for all requests.
//
// The transport preserves the behavior of the underlying transport while
// ensuring that the Authorization header is set with properly encoded
// Basic Authentication credentials on outgoing requests.
//
// Example:
//
//	transport := &BasicAuthTransport{
//	    Transport: http.DefaultTransport,
//	    Username: "myuser",
//	    Password: "mypassword",
//	}
type BasicAuthTransport struct {
	// Transport is the underlying http.RoundTripper to wrap
	Transport http.RoundTripper
	// Username is the username for Basic Authentication
	Username string
	// Password is the password for Basic Authentication
	Password string
}

// RoundTrip implements the http.RoundTripper interface and adds Basic Authentication to the request.
//
// RoundTrip sets the Authorization header on the provided request using HTTP Basic
// Authentication with the configured username and password. The credentials are
// automatically encoded using base64 as required by the HTTP Basic Authentication
// standard. After setting the authentication header, it delegates to the underlying
// Transport's RoundTrip method to perform the actual HTTP request.
//
// Parameters:
//   - req: The HTTP request to modify and execute
//
// Returns the HTTP response from the underlying transport and any error encountered.
//
// Example:
//
//	transport := &BasicAuthTransport{
//	    Transport: http.DefaultTransport,
//	    Username: "myuser",
//	    Password: "mypassword",
//	}
//	resp, err := transport.RoundTrip(req)
func (t *BasicAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(t.Username, t.Password)
	return t.Transport.RoundTrip(req)
}
