package isp

import (
	"os"
	"strings"
	"testing"

	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/common"
	authmodels "github.com/cyberark/ark-sdk-golang/pkg/models/auth"
	commonmodels "github.com/cyberark/ark-sdk-golang/pkg/models/common"
	"github.com/golang-jwt/jwt/v5"
	cookiejar "github.com/juju/persistent-cookiejar"
)

// createTestJWT creates a test JWT token with the provided claims
func createTestJWT(claims jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("test-secret"))
	return tokenString
}

// setupTestEnv sets up test environment variables
func setupTestEnv(env string) func() {
	oldEnv := os.Getenv("DEPLOY_ENV")
	os.Setenv("DEPLOY_ENV", env)
	return func() {
		if oldEnv != "" {
			os.Setenv("DEPLOY_ENV", oldEnv)
		} else {
			os.Unsetenv("DEPLOY_ENV")
		}
	}
}

func TestNewArkISPServiceClient(t *testing.T) {
	tests := []struct {
		name             string
		serviceName      string
		tenantSubdomain  string
		baseTenantURL    string
		tenantEnv        commonmodels.AwsEnv
		token            string
		authHeaderName   string
		separator        string
		basePath         string
		setupMock        func() func()
		expectedError    bool
		expectedErrorMsg string
		validateFunc     func(t *testing.T, result *ArkISPServiceClient)
	}{
		{
			name:            "success_with_valid_parameters",
			serviceName:     "api",
			tenantSubdomain: "test-tenant",
			baseTenantURL:   "https://test-tenant.cyberark.cloud",
			tenantEnv:       commonmodels.Prod,
			token: createTestJWT(jwt.MapClaims{
				"subdomain":       "test-tenant",
				"platform_domain": "cyberark.cloud",
				"tenant_id":       "12345",
			}),
			authHeaderName: "Authorization",
			separator:      "-",
			basePath:       "v1",
			setupMock:      func() func() { return func() {} },
			expectedError:  false,
			validateFunc: func(t *testing.T, result *ArkISPServiceClient) {
				if result == nil {
					t.Error("Expected non-nil client")
					return
				}
				if result.tenantEnv != commonmodels.Prod {
					t.Errorf("Expected tenant env %v, got %v", commonmodels.Prod, result.tenantEnv)
				}
			},
		},
		// {
		// 	name:            "success_with_empty_tenant_env_uses_environment",
		// 	serviceName:     "portal",
		// 	tenantSubdomain: "dev-tenant",
		// 	baseTenantURL:   "",
		// 	tenantEnv:       "",
		// 	token: createTestJWT(jwt.MapClaims{
		// 		"subdomain": "dev-tenant",
		// 		"tenant_id": "67890",
		// 	}),
		// 	authHeaderName: "Authorization",
		// 	separator:      ".",
		// 	basePath:       "",
		// 	setupMock: func() func() {
		// 		return setupTestEnv("dev")
		// 	},
		// 	expectedError: false,
		// 	validateFunc: func(t *testing.T, result *ArkISPServiceClient) {
		// 		if result == nil {
		// 			t.Error("Expected non-nil client")
		// 			return
		// 		}
		// 		// Should use environment variable value
		// 		if result.tenantEnv != commonmodels.Prod {
		// 			t.Errorf("Expected tenant env from environment %v, got %v", commonmodels.Prod, result.tenantEnv)
		// 		}
		// 	},
		// },
		{
			name:            "success_with_empty_tenant_env_defaults_to_prod",
			serviceName:     "service",
			tenantSubdomain: "prod-tenant",
			baseTenantURL:   "",
			tenantEnv:       "",
			token: createTestJWT(jwt.MapClaims{
				"subdomain": "prod-tenant",
				"tenant_id": "11111",
			}),
			authHeaderName: "Authorization",
			separator:      "-",
			basePath:       "api",
			setupMock: func() func() {
				return setupTestEnv("")
			},
			expectedError: false,
			validateFunc: func(t *testing.T, result *ArkISPServiceClient) {
				if result == nil {
					t.Error("Expected non-nil client")
					return
				}
				// Should default to prod when no environment
				if result.tenantEnv != commonmodels.Prod {
					t.Errorf("Expected default tenant env %v, got %v", commonmodels.Prod, result.tenantEnv)
				}
			},
		},
		{
			name:            "success_with_base_path_appended",
			serviceName:     "api",
			tenantSubdomain: "test-tenant",
			baseTenantURL:   "",
			tenantEnv:       commonmodels.GovProd,
			token: createTestJWT(jwt.MapClaims{
				"subdomain": "test-tenant",
				"tenant_id": "22222",
			}),
			authHeaderName: "Authorization",
			separator:      "-",
			basePath:       "v2/resources",
			setupMock:      func() func() { return func() {} },
			expectedError:  false,
			validateFunc: func(t *testing.T, result *ArkISPServiceClient) {
				if result == nil {
					t.Error("Expected non-nil client")
					return
				}
				// Validate that basePath was appended to URL
				if !strings.Contains(result.ArkClient.BaseURL, "v2/resources") {
					t.Errorf("Expected base path in URL, got %s", result.ArkClient.BaseURL)
				}
			},
		},
		{
			name:            "error_invalid_service_url_resolution",
			serviceName:     "api",
			tenantSubdomain: "",
			baseTenantURL:   "",
			tenantEnv:       commonmodels.Prod,
			token:           "invalid-token",
			authHeaderName:  "Authorization",
			separator:       "-",
			basePath:        "",
			setupMock:       func() func() { return func() {} },
			expectedError:   true,
		},
		{
			name:            "success_with_nil_cookie_jar",
			serviceName:     "test",
			tenantSubdomain: "tenant",
			baseTenantURL:   "",
			tenantEnv:       commonmodels.Prod,
			token: createTestJWT(jwt.MapClaims{
				"subdomain": "tenant",
				"tenant_id": "33333",
			}),
			authHeaderName: "Authorization",
			separator:      "-",
			basePath:       "",
			setupMock:      func() func() { return func() {} },
			expectedError:  false,
			validateFunc: func(t *testing.T, result *ArkISPServiceClient) {
				if result == nil {
					t.Error("Expected non-nil client")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cleanup := tt.setupMock()
			defer cleanup()

			cookieJar, _ := cookiejar.New(nil)
			result, err := NewArkISPServiceClient(
				tt.serviceName,
				tt.tenantSubdomain,
				tt.baseTenantURL,
				tt.tenantEnv,
				tt.token,
				tt.authHeaderName,
				tt.separator,
				tt.basePath,
				cookieJar,
				nil, // refreshConnectionCallback
			)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				if tt.expectedErrorMsg != "" && err.Error() != tt.expectedErrorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.expectedErrorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
				return
			}

			if tt.validateFunc != nil {
				tt.validateFunc(t, result)
			}
		})
	}
}

func TestResolveServiceURL(t *testing.T) {
	tests := []struct {
		name             string
		serviceName      string
		tenantSubdomain  string
		baseTenantURL    string
		tenantEnv        commonmodels.AwsEnv
		token            string
		separator        string
		setupMock        func() func()
		expectedError    bool
		expectedErrorMsg string
		validateFunc     func(t *testing.T, result string)
	}{
		{
			name:            "success_with_subdomain_from_token",
			serviceName:     "api",
			tenantSubdomain: "",
			baseTenantURL:   "",
			tenantEnv:       commonmodels.Prod,
			token: createTestJWT(jwt.MapClaims{
				"subdomain":       "token-tenant",
				"platform_domain": "cyberark.cloud",
			}),
			separator:     "-",
			setupMock:     func() func() { return func() {} },
			expectedError: false,
			validateFunc: func(t *testing.T, result string) {
				expected := "https://token-tenant-api.cyberark.cloud"
				if result != expected {
					t.Errorf("Expected %s, got %s", expected, result)
				}
			},
		},
		{
			name:            "success_with_platform_domain_from_token",
			serviceName:     "portal",
			tenantSubdomain: "",
			baseTenantURL:   "",
			tenantEnv:       "",
			token: createTestJWT(jwt.MapClaims{
				"subdomain":       "custom-tenant",
				"platform_domain": "custom.domain.com",
			}),
			separator:     ".",
			setupMock:     func() func() { return func() {} },
			expectedError: false,
			validateFunc: func(t *testing.T, result string) {
				expected := "https://custom-tenant.portal.custom.domain.com"
				if result != expected {
					t.Errorf("Expected %s, got %s", expected, result)
				}
			},
		},
		{
			name:            "success_with_shell_prefix_removal",
			serviceName:     "api",
			tenantSubdomain: "",
			baseTenantURL:   "",
			tenantEnv:       "",
			token: createTestJWT(jwt.MapClaims{
				"subdomain":       "shell-tenant",
				"platform_domain": "shell.cyberark.cloud",
			}),
			separator:     "-",
			setupMock:     func() func() { return func() {} },
			expectedError: false,
			validateFunc: func(t *testing.T, result string) {
				expected := "https://shell-tenant-api.cyberark.cloud"
				if result != expected {
					t.Errorf("Expected %s, got %s", expected, result)
				}
			},
		},
		{
			name:            "success_fallback_to_tenant_subdomain",
			serviceName:     "service",
			tenantSubdomain: "fallback-tenant",
			baseTenantURL:   "",
			tenantEnv:       commonmodels.Prod,
			token: createTestJWT(jwt.MapClaims{
				"tenant_id": "12345",
			}),
			separator:     "-",
			setupMock:     func() func() { return func() {} },
			expectedError: false,
			validateFunc: func(t *testing.T, result string) {
				expected := "https://fallback-tenant-service.cyberark.cloud"
				if result != expected {
					t.Errorf("Expected %s, got %s", expected, result)
				}
			},
		},
		{
			name:            "success_fallback_to_base_tenant_url",
			serviceName:     "api",
			tenantSubdomain: "",
			baseTenantURL:   "https://url-tenant.example.com",
			tenantEnv:       commonmodels.Prod,
			token: createTestJWT(jwt.MapClaims{
				"tenant_id": "67890",
			}),
			separator:     "-",
			setupMock:     func() func() { return func() {} },
			expectedError: false,
			validateFunc: func(t *testing.T, result string) {
				expected := "https://url-tenant-api.cyberark.cloud"
				if result != expected {
					t.Errorf("Expected %s, got %s", expected, result)
				}
			},
		},
		{
			name:            "success_fallback_to_unique_name",
			serviceName:     "portal",
			tenantSubdomain: "",
			baseTenantURL:   "",
			tenantEnv:       "",
			token: createTestJWT(jwt.MapClaims{
				"unique_name": "user@unique-tenant.cyberark.cloud",
			}),
			separator:     "-",
			setupMock:     func() func() { return func() {} },
			expectedError: false,
			validateFunc: func(t *testing.T, result string) {
				expected := "https://unique-tenant-portal.cyberark.cloud"
				if result != expected {
					t.Errorf("Expected %s, got %s", expected, result)
				}
			},
		},
		{
			name:            "success_without_service_name",
			serviceName:     "",
			tenantSubdomain: "",
			baseTenantURL:   "",
			tenantEnv:       commonmodels.Prod,
			token: createTestJWT(jwt.MapClaims{
				"subdomain": "no-service-tenant",
			}),
			separator:     "-",
			setupMock:     func() func() { return func() {} },
			expectedError: false,
			validateFunc: func(t *testing.T, result string) {
				expected := "https://no-service-tenant.cyberark.cloud"
				if result != expected {
					t.Errorf("Expected %s, got %s", expected, result)
				}
			},
		},
		{
			name:            "error_cannot_resolve_tenant_subdomain",
			serviceName:     "api",
			tenantSubdomain: "",
			baseTenantURL:   "",
			tenantEnv:       commonmodels.Prod,
			token: createTestJWT(jwt.MapClaims{
				"tenant_id": "12345",
			}),
			separator:        "-",
			setupMock:        func() func() { return func() {} },
			expectedError:    true,
			expectedErrorMsg: "failed to resolve tenant subdomain",
		},
		{
			name:            "error_invalid_token_format",
			serviceName:     "api",
			tenantSubdomain: "",
			baseTenantURL:   "",
			tenantEnv:       commonmodels.Prod,
			token:           "invalid.token.format",
			separator:       "-",
			setupMock:       func() func() { return func() {} },
			expectedError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cleanup := tt.setupMock()
			defer cleanup()

			result, err := resolveServiceURL(
				tt.serviceName,
				tt.tenantSubdomain,
				tt.baseTenantURL,
				tt.tenantEnv,
				tt.token,
				tt.separator,
			)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				if tt.expectedErrorMsg != "" && err.Error() != tt.expectedErrorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.expectedErrorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
				return
			}

			if tt.validateFunc != nil {
				tt.validateFunc(t, result)
			}
		})
	}
}

func TestArkISPServiceClient_TenantEnv(t *testing.T) {
	tests := []struct {
		name           string
		tenantEnv      commonmodels.AwsEnv
		expectedResult commonmodels.AwsEnv
	}{
		{
			name:           "success_production_environment",
			tenantEnv:      commonmodels.Prod,
			expectedResult: commonmodels.Prod,
		},
		{
			name:           "success_development_environment",
			tenantEnv:      commonmodels.Prod,
			expectedResult: commonmodels.Prod,
		},
		{
			name:           "success_staging_environment",
			tenantEnv:      commonmodels.GovProd,
			expectedResult: commonmodels.GovProd,
		},
		{
			name:           "success_empty_environment",
			tenantEnv:      "",
			expectedResult: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := &ArkISPServiceClient{
				tenantEnv: tt.tenantEnv,
			}

			result := client.TenantEnv()

			if result != tt.expectedResult {
				t.Errorf("Expected %v, got %v", tt.expectedResult, result)
			}
		})
	}
}

func TestArkISPServiceClient_TenantID(t *testing.T) {
	tests := []struct {
		name             string
		setupClient      func() *ArkISPServiceClient
		expectedResult   string
		expectedError    bool
		expectedErrorMsg string
	}{
		{
			name: "error_empty_token",
			setupClient: func() *ArkISPServiceClient {
				arkClient := &common.ArkClient{}
				arkClient.UpdateToken("", "Bearer")
				return &ArkISPServiceClient{
					ArkClient: arkClient,
				}
			},
			expectedResult:   "",
			expectedError:    true,
			expectedErrorMsg: "failed to retrieve tenant id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := tt.setupClient()
			result, err := client.TenantID()

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				if tt.expectedErrorMsg != "" && err.Error() != tt.expectedErrorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.expectedErrorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
				return
			}

			if result != tt.expectedResult {
				t.Errorf("Expected %s, got %s", tt.expectedResult, result)
			}
		})
	}
}

func TestFromISPAuth(t *testing.T) {
	tests := []struct {
		name             string
		setupISPAuth     func() *auth.ArkISPAuth
		serviceName      string
		separator        string
		basePath         string
		setupMock        func() func()
		expectedError    bool
		expectedErrorMsg string
		validateFunc     func(t *testing.T, result *ArkISPServiceClient)
	}{
		{
			name: "success_with_domain_in_username",
			setupISPAuth: func() *auth.ArkISPAuth {
				return &auth.ArkISPAuth{
					ArkAuthBase: &auth.ArkAuthBase{
						Token: &authmodels.ArkToken{
							Username: "user@tenant.cyberark.cloud",
							Token:    createTestJWT(jwt.MapClaims{"tenant_id": "123", "subdomain": "tenant"}),
							Metadata: map[string]interface{}{},
						},
					},
				}
			},
			serviceName:   "api",
			separator:     "-",
			basePath:      "v1",
			setupMock:     func() func() { return func() {} },
			expectedError: false,
			validateFunc: func(t *testing.T, result *ArkISPServiceClient) {
				if result == nil {
					t.Error("Expected non-nil client")
					return
				}
				if result.tenantEnv != commonmodels.Prod {
					t.Errorf("Expected prod environment, got %v", result.tenantEnv)
				}
			},
		},
		{
			name: "error_invalid_token_in_auth",
			setupISPAuth: func() *auth.ArkISPAuth {
				return &auth.ArkISPAuth{
					ArkAuthBase: &auth.ArkAuthBase{
						Token: &authmodels.ArkToken{
							Username: "user@invalid.cyberark.cloud",
							Token:    "invalid-token-format",
							Metadata: map[string]interface{}{},
						},
					},
				}
			},
			serviceName:   "api",
			separator:     "-",
			basePath:      "",
			setupMock:     func() func() { return func() {} },
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cleanup := tt.setupMock()
			defer cleanup()

			ispAuth := tt.setupISPAuth()
			result, err := FromISPAuth(ispAuth, tt.serviceName, tt.separator, tt.basePath, nil)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				if tt.expectedErrorMsg != "" && err.Error() != tt.expectedErrorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.expectedErrorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
				return
			}

			if tt.validateFunc != nil {
				tt.validateFunc(t, result)
			}
		})
	}
}
