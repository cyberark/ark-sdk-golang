package common

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHeaderTransport_RoundTrip(t *testing.T) {
	tests := []struct {
		name             string
		headers          map[string]string
		requestHeaders   map[string]string
		expectedHeaders  map[string]string
		mockTransportErr error
		expectedError    bool
		validateFunc     func(t *testing.T, req *http.Request, resp *http.Response)
	}{
		{
			name: "success_adds_single_header",
			headers: map[string]string{
				"User-Agent": "MyApp/1.0",
			},
			requestHeaders: map[string]string{},
			expectedHeaders: map[string]string{
				"User-Agent": "MyApp/1.0",
			},
			expectedError: false,
		},
		{
			name: "success_overwrites_existing_headers",
			headers: map[string]string{
				"User-Agent": "MyApp/1.0",
				"X-Custom":   "new-value",
			},
			requestHeaders: map[string]string{
				"User-Agent": "OldApp/0.1",
				"X-Custom":   "old-value",
				"X-Keep":     "keep-value",
			},
			expectedHeaders: map[string]string{
				"User-Agent": "MyApp/1.0",
				"X-Custom":   "new-value",
				"X-Keep":     "keep-value",
			},
			expectedError: false,
		},
		{
			name:            "success_handles_empty_headers_map",
			headers:         map[string]string{},
			requestHeaders:  map[string]string{"Existing": "value"},
			expectedHeaders: map[string]string{"Existing": "value"},
			expectedError:   false,
		},
		{
			name:            "success_handles_nil_headers_map",
			headers:         nil,
			requestHeaders:  map[string]string{"Existing": "value"},
			expectedHeaders: map[string]string{"Existing": "value"},
			expectedError:   false,
		},
		{
			name: "edge_case_empty_header_values",
			headers: map[string]string{
				"Empty-Header": "",
				"X-Test":       "value",
			},
			requestHeaders: map[string]string{},
			expectedHeaders: map[string]string{
				"Empty-Header": "",
				"X-Test":       "value",
			},
			expectedError: false,
		},
		{
			name: "edge_case_special_characters_in_headers",
			headers: map[string]string{
				"X-Special": "value with spaces and symbols!@#$%",
				"X-Unicode": "测试中文",
			},
			requestHeaders: map[string]string{},
			expectedHeaders: map[string]string{
				"X-Special": "value with spaces and symbols!@#$%",
				"X-Unicode": "测试中文",
			},
			expectedError: false,
		},
		{
			name: "error_propagates_transport_error",
			headers: map[string]string{
				"User-Agent": "MyApp/1.0",
			},
			requestHeaders:   map[string]string{},
			mockTransportErr: &mockTransportError{message: "transport error"},
			expectedError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mock transport
			mockTransport := &mockRoundTripper{
				err: tt.mockTransportErr,
			}

			// Create HeaderTransport
			transport := &HeaderTransport{
				Transport: mockTransport,
				Headers:   tt.headers,
			}

			// Create request with initial headers
			req := httptest.NewRequest("GET", "http://example.com", nil)
			for key, value := range tt.requestHeaders {
				req.Header.Set(key, value)
			}

			// Execute RoundTrip
			resp, err := transport.RoundTrip(req)

			// Validate error expectation
			if tt.expectedError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
				return
			}

			// Validate response
			if resp == nil {
				t.Error("Expected response, got nil")
				return
			}

			// Validate headers were set correctly
			for expectedKey, expectedValue := range tt.expectedHeaders {
				actualValue := req.Header.Get(expectedKey)
				if actualValue != expectedValue {
					t.Errorf("Expected header %s='%s', got '%s'", expectedKey, expectedValue, actualValue)
				}
			}

			// Verify no unexpected headers were added (only check if we have expected headers)
			if len(tt.expectedHeaders) > 0 {
				for headerKey := range req.Header {
					if _, expected := tt.expectedHeaders[headerKey]; !expected {
						// This header wasn't in our expected list, but might be from initial request
						if _, wasInitial := tt.requestHeaders[headerKey]; !wasInitial {
							t.Errorf("Unexpected header found: %s", headerKey)
						}
					}
				}
			}

			// Run custom validation if provided
			if tt.validateFunc != nil {
				tt.validateFunc(t, req, resp)
			}
		})
	}
}

func TestBasicAuthTransport_RoundTrip(t *testing.T) {
	tests := []struct {
		name               string
		username           string
		password           string
		requestHeaders     map[string]string
		expectedAuthHeader string
		mockTransportErr   error
		expectedError      bool
		validateFunc       func(t *testing.T, req *http.Request, resp *http.Response)
	}{
		{
			name:               "success_sets_basic_auth_normal_credentials",
			username:           "testuser",
			password:           "testpass",
			requestHeaders:     map[string]string{},
			expectedAuthHeader: "Basic dGVzdHVzZXI6dGVzdHBhc3M=", // base64("testuser:testpass")
			expectedError:      false,
		},
		{
			name:               "success_sets_basic_auth_with_special_characters",
			username:           "user@domain.com",
			password:           "p@ssw0rd!",
			requestHeaders:     map[string]string{},
			expectedAuthHeader: "Basic dXNlckBkb21haW4uY29tOnBAc3N3MHJkIQ==", // base64("user@domain.com:p@ssw0rd!")
			expectedError:      false,
		},
		{
			name:               "success_overwrites_existing_auth_header",
			username:           "newuser",
			password:           "newpass",
			requestHeaders:     map[string]string{"Authorization": "Bearer oldtoken"},
			expectedAuthHeader: "Basic bmV3dXNlcjpuZXdwYXNz", // base64("newuser:newpass")
			expectedError:      false,
		},
		{
			name:               "edge_case_empty_username",
			username:           "",
			password:           "password",
			requestHeaders:     map[string]string{},
			expectedAuthHeader: "Basic OnBhc3N3b3Jk", // base64(":password")
			expectedError:      false,
		},
		{
			name:               "edge_case_empty_password",
			username:           "username",
			password:           "",
			requestHeaders:     map[string]string{},
			expectedAuthHeader: "Basic dXNlcm5hbWU6", // base64("username:")
			expectedError:      false,
		},
		{
			name:               "edge_case_both_empty",
			username:           "",
			password:           "",
			requestHeaders:     map[string]string{},
			expectedAuthHeader: "Basic Og==", // base64(":")
			expectedError:      false,
		},
		{
			name:               "success_preserves_other_headers",
			username:           "user",
			password:           "pass",
			requestHeaders:     map[string]string{"X-Custom": "value", "User-Agent": "MyApp"},
			expectedAuthHeader: "Basic dXNlcjpwYXNz", // base64("user:pass")
			expectedError:      false,
			validateFunc: func(t *testing.T, req *http.Request, resp *http.Response) {
				if req.Header.Get("X-Custom") != "value" {
					t.Errorf("Expected X-Custom header to be preserved")
				}
				if req.Header.Get("User-Agent") != "MyApp" {
					t.Errorf("Expected User-Agent header to be preserved")
				}
			},
		},
		{
			name:             "error_propagates_transport_error",
			username:         "user",
			password:         "pass",
			requestHeaders:   map[string]string{},
			mockTransportErr: &mockTransportError{message: "transport failed"},
			expectedError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mock transport
			mockTransport := &mockRoundTripper{
				err: tt.mockTransportErr,
			}

			// Create BasicAuthTransport
			transport := &BasicAuthTransport{
				Transport: mockTransport,
				Username:  tt.username,
				Password:  tt.password,
			}

			// Create request with initial headers
			req := httptest.NewRequest("GET", "http://example.com", nil)
			for key, value := range tt.requestHeaders {
				req.Header.Set(key, value)
			}

			// Execute RoundTrip
			resp, err := transport.RoundTrip(req)

			// Validate error expectation
			if tt.expectedError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
				return
			}

			// Validate response
			if resp == nil {
				t.Error("Expected response, got nil")
				return
			}

			// Validate Authorization header was set correctly
			authHeader := req.Header.Get("Authorization")
			if authHeader != tt.expectedAuthHeader {
				t.Errorf("Expected Authorization header '%s', got '%s'", tt.expectedAuthHeader, authHeader)
			}

			// Run custom validation if provided
			if tt.validateFunc != nil {
				tt.validateFunc(t, req, resp)
			}
		})
	}
}

func TestHeaderTransport_Integration(t *testing.T) {
	tests := []struct {
		name         string
		headers      map[string]string
		validateFunc func(t *testing.T, req *http.Request)
	}{
		{
			name: "integration_real_http_server",
			headers: map[string]string{
				"User-Agent":   "TestApp/1.0",
				"X-Test-Flag":  "true",
				"Content-Type": "application/json",
			},
			validateFunc: func(t *testing.T, req *http.Request) {
				if req.Header.Get("User-Agent") != "TestApp/1.0" {
					t.Errorf("Expected User-Agent header to be set")
				}
				if req.Header.Get("X-Test-Flag") != "true" {
					t.Errorf("Expected X-Test-Flag header to be set")
				}
				if req.Header.Get("Content-Type") != "application/json" {
					t.Errorf("Expected Content-Type header to be set")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server to capture request
			var capturedRequest *http.Request
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				capturedRequest = r
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			// Create client with HeaderTransport
			transport := &HeaderTransport{
				Transport: http.DefaultTransport,
				Headers:   tt.headers,
			}
			client := &http.Client{Transport: transport}

			// Make request
			resp, err := client.Get(server.URL)
			if err != nil {
				t.Errorf("Request failed: %v", err)
				return
			}
			defer resp.Body.Close()

			// Validate request was captured
			if capturedRequest == nil {
				t.Error("Expected request to be captured")
				return
			}

			// Run custom validation
			if tt.validateFunc != nil {
				tt.validateFunc(t, capturedRequest)
			}
		})
	}
}

func TestBasicAuthTransport_Integration(t *testing.T) {
	tests := []struct {
		name         string
		username     string
		password     string
		validateFunc func(t *testing.T, req *http.Request)
	}{
		{
			name:     "integration_real_http_server",
			username: "testuser",
			password: "testpass",
			validateFunc: func(t *testing.T, req *http.Request) {
				authHeader := req.Header.Get("Authorization")
				expectedAuth := "Basic dGVzdHVzZXI6dGVzdHBhc3M=" // base64("testuser:testpass")
				if authHeader != expectedAuth {
					t.Errorf("Expected Authorization header '%s', got '%s'", expectedAuth, authHeader)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server to capture request
			var capturedRequest *http.Request
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				capturedRequest = r
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			// Create client with BasicAuthTransport
			transport := &BasicAuthTransport{
				Transport: http.DefaultTransport,
				Username:  tt.username,
				Password:  tt.password,
			}
			client := &http.Client{Transport: transport}

			// Make request
			resp, err := client.Get(server.URL)
			if err != nil {
				t.Errorf("Request failed: %v", err)
				return
			}
			defer resp.Body.Close()

			// Validate request was captured
			if capturedRequest == nil {
				t.Error("Expected request to be captured")
				return
			}

			// Run custom validation
			if tt.validateFunc != nil {
				tt.validateFunc(t, capturedRequest)
			}
		})
	}
}

func TestTransportStructures(t *testing.T) {
	tests := []struct {
		name           string
		createStruct   func() interface{}
		validateFields func(t *testing.T, obj interface{})
	}{
		{
			name: "header_transport_struct_fields",
			createStruct: func() interface{} {
				return &HeaderTransport{
					Transport: http.DefaultTransport,
					Headers:   map[string]string{"test": "value"},
				}
			},
			validateFields: func(t *testing.T, obj interface{}) {
				transport := obj.(*HeaderTransport)
				if transport.Transport == nil {
					t.Error("Expected Transport field to be set")
				}
				if transport.Headers == nil {
					t.Error("Expected Headers field to be set")
				}
				if len(transport.Headers) != 1 {
					t.Errorf("Expected 1 header, got %d", len(transport.Headers))
				}
			},
		},
		{
			name: "basic_auth_transport_struct_fields",
			createStruct: func() interface{} {
				return &BasicAuthTransport{
					Transport: http.DefaultTransport,
					Username:  "user",
					Password:  "pass",
				}
			},
			validateFields: func(t *testing.T, obj interface{}) {
				transport := obj.(*BasicAuthTransport)
				if transport.Transport == nil {
					t.Error("Expected Transport field to be set")
				}
				if transport.Username != "user" {
					t.Errorf("Expected Username 'user', got '%s'", transport.Username)
				}
				if transport.Password != "pass" {
					t.Errorf("Expected Password 'pass', got '%s'", transport.Password)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			obj := tt.createStruct()
			if tt.validateFields != nil {
				tt.validateFields(t, obj)
			}
		})
	}
}

// Helper functions and mocks for tests

// mockRoundTripper is a mock implementation of http.RoundTripper for testing
type mockRoundTripper struct {
	response *http.Response
	err      error
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.err != nil {
		return nil, m.err
	}

	if m.response != nil {
		return m.response, nil
	}

	// Default successful response
	return &http.Response{
		StatusCode: 200,
		Header:     http.Header{},
		Body:       http.NoBody,
		Request:    req,
	}, nil
}

// mockTransportError is a mock error for testing
type mockTransportError struct {
	message string
}

func (e *mockTransportError) Error() string {
	return e.message
}

// TestRoundTripperInterface verifies that our transports implement http.RoundTripper
func TestRoundTripperInterface(t *testing.T) {
	tests := []struct {
		name      string
		transport http.RoundTripper
	}{
		{
			name: "header_transport_implements_interface",
			transport: &HeaderTransport{
				Transport: http.DefaultTransport,
				Headers:   map[string]string{},
			},
		},
		{
			name: "basic_auth_transport_implements_interface",
			transport: &BasicAuthTransport{
				Transport: http.DefaultTransport,
				Username:  "user",
				Password:  "pass",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Verify it implements http.RoundTripper interface
			var _ http.RoundTripper = tt.transport

			// Test that RoundTrip method exists and can be called
			req := httptest.NewRequest("GET", "http://example.com", nil)
			_, err := tt.transport.RoundTrip(req)
			if err != nil {
				t.Errorf("RoundTrip method failed: %v", err)
			}
		})
	}
}
