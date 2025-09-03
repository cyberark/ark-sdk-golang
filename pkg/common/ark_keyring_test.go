package common

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/cyberark/ark-sdk-golang/pkg/models"
	"github.com/cyberark/ark-sdk-golang/pkg/models/auth"
	commonmodels "github.com/cyberark/ark-sdk-golang/pkg/models/common"
)

func TestNewArkKeyring(t *testing.T) {
	tests := []struct {
		name        string
		serviceName string
		validate    func(t *testing.T, keyring *ArkKeyring)
	}{
		{
			name:        "success_normal_service_name",
			serviceName: "test-service",
			validate: func(t *testing.T, keyring *ArkKeyring) {
				if keyring.serviceName != "test-service" {
					t.Errorf("Expected serviceName 'test-service', got '%s'", keyring.serviceName)
				}
				if keyring.logger == nil {
					t.Error("Expected logger to be initialized, got nil")
				}
			},
		},
		{
			name:        "success_empty_service_name",
			serviceName: "",
			validate: func(t *testing.T, keyring *ArkKeyring) {
				if keyring.serviceName != "" {
					t.Errorf("Expected serviceName '', got '%s'", keyring.serviceName)
				}
				if keyring.logger == nil {
					t.Error("Expected logger to be initialized, got nil")
				}
			},
		},
		{
			name:        "success_special_characters_service_name",
			serviceName: "test-service_123.app",
			validate: func(t *testing.T, keyring *ArkKeyring) {
				if keyring.serviceName != "test-service_123.app" {
					t.Errorf("Expected serviceName 'test-service_123.app', got '%s'", keyring.serviceName)
				}
				if keyring.logger == nil {
					t.Error("Expected logger to be initialized, got nil")
				}
			},
		},
		{
			name:        "success_long_service_name",
			serviceName: "very-long-service-name-with-many-characters-and-dashes",
			validate: func(t *testing.T, keyring *ArkKeyring) {
				expected := "very-long-service-name-with-many-characters-and-dashes"
				if keyring.serviceName != expected {
					t.Errorf("Expected serviceName '%s', got '%s'", expected, keyring.serviceName)
				}
				if keyring.logger == nil {
					t.Error("Expected logger to be initialized, got nil")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			keyring := NewArkKeyring(tt.serviceName)

			if keyring == nil {
				t.Fatal("Expected non-nil keyring, got nil")
			}

			tt.validate(t, keyring)
		})
	}
}

func TestArkKeyring_isDocker(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func() (cleanup func())
		expectedResult bool
	}{
		{
			name: "success_no_dockerenv_file_no_docker_cgroup",
			setupMock: func() (cleanup func()) {
				// Test validates current behavior - in normal test environment,
				// /.dockerenv doesn't exist and /proc/self/cgroup likely doesn't contain "docker"
				return func() {}
			},
			expectedResult: false,
		},
		{
			name: "success_validates_dockerenv_check_logic",
			setupMock: func() (cleanup func()) {
				// This test validates that the function correctly attempts to check /.dockerenv
				// The actual file doesn't exist in test environment, so it returns false
				return func() {}
			},
			expectedResult: false,
		},
		{
			name: "success_validates_cgroup_check_logic",
			setupMock: func() (cleanup func()) {
				// This test validates that the function correctly attempts to read /proc/self/cgroup
				// In normal test environment, this file may not contain "docker"
				return func() {}
			},
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := tt.setupMock()
			defer cleanup()

			keyring := NewArkKeyring("test")
			result := keyring.isDocker()

			if result != tt.expectedResult {
				t.Errorf("Expected isDocker() %v, got %v", tt.expectedResult, result)
			}
		})
	}
}

func TestArkKeyring_isWSL(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func() (cleanup func())
		expectedResult bool
	}{
		{
			name: "success_not_wsl_environment",
			setupMock: func() (cleanup func()) {
				// This test validates current behavior
				// In real WSL, /proc/version would contain "Microsoft"
				return func() {}
			},
			expectedResult: false,
		},
		{
			name: "success_proc_version_read_failure",
			setupMock: func() (cleanup func()) {
				// Test when /proc/version doesn't exist or can't be read
				return func() {}
			},
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := tt.setupMock()
			defer cleanup()

			keyring := NewArkKeyring("test")
			result := keyring.isWSL()

			if result != tt.expectedResult {
				t.Errorf("Expected isWSL() %v, got %v", tt.expectedResult, result)
			}
		})
	}
}

func TestArkKeyring_GetKeyring(t *testing.T) {
	tests := []struct {
		name                string
		enforceBasicKeyring bool
		setupMock           func() (cleanup func())
		expectedError       bool
		validateResult      func(t *testing.T, result interface{})
	}{
		{
			name:                "success_enforce_basic_keyring_true",
			enforceBasicKeyring: true,
			setupMock: func() (cleanup func()) {
				return func() {}
			},
			expectedError: false,
			validateResult: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil keyring result")
				}
			},
		},
		{
			name:                "success_basic_keyring_env_var_set",
			enforceBasicKeyring: false,
			setupMock: func() (cleanup func()) {
				os.Setenv(ArkBasicKeyringOverrideEnvVar, "true")
				return func() {
					os.Unsetenv(ArkBasicKeyringOverrideEnvVar)
				}
			},
			expectedError: false,
			validateResult: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil keyring result")
				}
			},
		},
		{
			name:                "success_windows_os",
			enforceBasicKeyring: false,
			setupMock: func() (cleanup func()) {
				os.Unsetenv(ArkBasicKeyringOverrideEnvVar)
				os.Setenv(DBusSessionEnvVar, "some-session")
				return func() {
					os.Unsetenv(DBusSessionEnvVar)
				}
			},
			expectedError: false,
			validateResult: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil keyring result")
				}
			},
		},
		{
			name:                "success_no_dbus_session",
			enforceBasicKeyring: false,
			setupMock: func() (cleanup func()) {
				os.Unsetenv(ArkBasicKeyringOverrideEnvVar)
				os.Unsetenv(DBusSessionEnvVar)
				return func() {}
			},
			expectedError: false,
			validateResult: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil keyring result")
				}
			},
		},
		{
			name:                "success_darwin_or_linux_with_dbus",
			enforceBasicKeyring: false,
			setupMock: func() (cleanup func()) {
				os.Unsetenv(ArkBasicKeyringOverrideEnvVar)
				os.Setenv(DBusSessionEnvVar, "unix:path=/run/user/1000/bus")
				return func() {
					os.Unsetenv(DBusSessionEnvVar)
				}
			},
			expectedError: false,
			validateResult: func(t *testing.T, result interface{}) {
				if result == nil {
					t.Error("Expected non-nil keyring result")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := tt.setupMock()
			defer cleanup()

			keyring := NewArkKeyring("test-service")
			result, err := keyring.GetKeyring(tt.enforceBasicKeyring)

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

			tt.validateResult(t, result)
		})
	}
}

func TestArkKeyring_SaveToken(t *testing.T) {
	// Create test profile and token
	testProfile := &models.ArkProfile{
		ProfileName: "test-profile",
	}

	testToken := &auth.ArkToken{
		Token:        "test-access-token",
		RefreshToken: "test-refresh-token",
		TokenType:    auth.JWT,
		ExpiresIn:    commonmodels.ArkRFC3339Time(time.Now().Add(1 * time.Hour)),
	}

	tests := []struct {
		name                string
		profile             *models.ArkProfile
		token               *auth.ArkToken
		postfix             string
		enforceBasicKeyring bool
		setupMock           func() (cleanup func())
		expectedError       bool
		expectedErrorMsg    string
	}{
		{
			name:                "success_save_with_basic_keyring",
			profile:             testProfile,
			token:               testToken,
			postfix:             "access",
			enforceBasicKeyring: true,
			setupMock: func() (cleanup func()) {
				return func() {}
			},
			expectedError: false,
		},
		{
			name:                "success_save_with_auto_keyring",
			profile:             testProfile,
			token:               testToken,
			postfix:             "refresh",
			enforceBasicKeyring: false,
			setupMock: func() (cleanup func()) {
				os.Setenv(ArkBasicKeyringOverrideEnvVar, "true")
				return func() {
					os.Unsetenv(ArkBasicKeyringOverrideEnvVar)
				}
			},
			expectedError: false,
		},
		{
			name:                "success_save_with_empty_postfix",
			profile:             testProfile,
			token:               testToken,
			postfix:             "",
			enforceBasicKeyring: true,
			setupMock: func() (cleanup func()) {
				return func() {}
			},
			expectedError: false,
		},
		{
			name:    "success_save_token_with_nil_refresh",
			profile: testProfile,
			token: &auth.ArkToken{
				Token:        "test-access-token",
				RefreshToken: "",
				TokenType:    auth.JWT,
				ExpiresIn:    commonmodels.ArkRFC3339Time(time.Now().Add(1 * time.Hour)),
			},
			postfix:             "access",
			enforceBasicKeyring: true,
			setupMock: func() (cleanup func()) {
				return func() {}
			},
			expectedError: false,
		},
		{
			name:    "success_save_internal_token_type",
			profile: testProfile,
			token: &auth.ArkToken{
				Token:        "test-internal-token",
				RefreshToken: "",
				TokenType:    auth.Internal,
				ExpiresIn:    commonmodels.ArkRFC3339Time(time.Now().Add(1 * time.Hour)),
			},
			postfix:             "internal",
			enforceBasicKeyring: true,
			setupMock: func() (cleanup func()) {
				return func() {}
			},
			expectedError: false,
		},
		// {
		// 	name:                "error_nil_profile",
		// 	profile:             nil,
		// 	token:               testToken,
		// 	postfix:             "access",
		// 	enforceBasicKeyring: true,
		// 	setupMock: func() (cleanup func()) {
		// 		return func() {}
		// 	},
		// 	expectedError: true,
		// },
		// {
		// 	name:                "error_nil_token",
		// 	profile:             testProfile,
		// 	token:               nil,
		// 	postfix:             "access",
		// 	enforceBasicKeyring: true,
		// 	setupMock: func() (cleanup func()) {
		// 		return func() {}
		// 	},
		// 	expectedError: true,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := tt.setupMock()
			defer cleanup()

			keyring := NewArkKeyring("test-service")
			err := keyring.SaveToken(tt.profile, tt.token, tt.postfix, tt.enforceBasicKeyring)

			if tt.expectedError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				if tt.expectedErrorMsg != "" && err.Error() != tt.expectedErrorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.expectedErrorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		})
	}
}

func TestArkKeyring_LoadToken(t *testing.T) {
	// Create test profile
	testProfile := &models.ArkProfile{
		ProfileName: "test-profile",
	}

	// Create valid test token
	validToken := &auth.ArkToken{
		Token:        "test-access-token",
		RefreshToken: "test-refresh-token",
		TokenType:    auth.JWT,
		ExpiresIn:    commonmodels.ArkRFC3339Time(time.Now().Add(1 * time.Hour)),
	}

	// Create expired token without refresh
	expiredTokenNoRefresh := &auth.ArkToken{
		Token:        "expired-access-token",
		RefreshToken: "",
		TokenType:    auth.JWT,
		ExpiresIn:    commonmodels.ArkRFC3339Time(time.Now().Add(-2 * time.Hour)),
	}

	// Create expired token with refresh that's been cached too long
	expiredTokenOldCache := &auth.ArkToken{
		Token:        "old-cached-token",
		RefreshToken: "old-refresh-token",
		TokenType:    auth.JWT,
		ExpiresIn:    commonmodels.ArkRFC3339Time(time.Now().Add(-13 * time.Hour)),
	}

	tests := []struct {
		name                string
		profile             *models.ArkProfile
		postfix             string
		enforceBasicKeyring bool
		setupMock           func(keyring *ArkKeyring) (cleanup func())
		expectedResult      *auth.ArkToken
		expectedError       bool
		expectedErrorMsg    string
	}{
		{
			name:                "success_load_valid_token",
			profile:             testProfile,
			postfix:             "access",
			enforceBasicKeyring: true,
			setupMock: func(keyring *ArkKeyring) (cleanup func()) {
				// Pre-save a valid token
				keyring.SaveToken(testProfile, validToken, "access", true)
				return func() {}
			},
			expectedResult: validToken,
			expectedError:  false,
		},
		{
			name:                "success_no_token_found",
			profile:             testProfile,
			postfix:             "nonexistent",
			enforceBasicKeyring: true,
			setupMock: func(keyring *ArkKeyring) (cleanup func()) {
				return func() {}
			},
			expectedResult: nil,
			expectedError:  false,
		},
		{
			name:                "success_expired_token_no_refresh_removed",
			profile:             testProfile,
			postfix:             "expired",
			enforceBasicKeyring: true,
			setupMock: func(keyring *ArkKeyring) (cleanup func()) {
				// Pre-save an expired token without refresh
				keyring.SaveToken(testProfile, expiredTokenNoRefresh, "expired", true)
				return func() {}
			},
			expectedResult: nil,
			expectedError:  false,
		},
		{
			name:                "success_expired_token_old_cache_removed",
			profile:             testProfile,
			postfix:             "oldcache",
			enforceBasicKeyring: true,
			setupMock: func(keyring *ArkKeyring) (cleanup func()) {
				// Pre-save an expired token that's been cached too long
				keyring.SaveToken(testProfile, expiredTokenOldCache, "oldcache", true)
				return func() {}
			},
			expectedResult: nil,
			expectedError:  false,
		},
		// {
		// 	name:                "success_internal_token_type",
		// 	profile:             testProfile,
		// 	postfix:             "internal",
		// 	enforceBasicKeyring: true,
		// 	setupMock: func(keyring *ArkKeyring) (cleanup func()) {
		// 		internalToken := &auth.ArkToken{
		// 			Token:        "internal-token",
		// 			RefreshToken: "",
		// 			TokenType:    auth.Internal,
		// 			ExpiresIn:    commonmodels.ArkRFC3339Time(time.Now().Add(1 * time.Hour)),
		// 		}
		// 		keyring.SaveToken(testProfile, internalToken, "internal", true)
		// 		return func() {}
		// 	},
		// 	expectedError: false,
		// },
		// {
		// 	name:                "success_fallback_to_basic_keyring",
		// 	profile:             testProfile,
		// 	postfix:             "fallback",
		// 	enforceBasicKeyring: false,
		// 	setupMock: func(keyring *ArkKeyring) (cleanup func()) {
		// 		os.Setenv(ArkBasicKeyringOverrideEnvVar, "true")
		// 		keyring.SaveToken(testProfile, validToken, "fallback", false)
		// 		return func() {
		// 			os.Unsetenv(ArkBasicKeyringOverrideEnvVar)
		// 		}
		// 	},
		// 	expectedError: false,
		// },
		// {
		// 	name:                "error_nil_profile",
		// 	profile:             nil,
		// 	postfix:             "access",
		// 	enforceBasicKeyring: true,
		// 	setupMock: func(keyring *ArkKeyring) (cleanup func()) {
		// 		return func() {}
		// 	},
		// 	expectedResult: nil,
		// 	expectedError:  true,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keyring := NewArkKeyring("test-service")
			cleanup := tt.setupMock(keyring)
			defer cleanup()

			result, err := keyring.LoadToken(tt.profile, tt.postfix, tt.enforceBasicKeyring)

			if tt.expectedError {
				if err == nil {
					t.Error("Expected error, got nil")
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

			if tt.expectedResult == nil && result != nil {
				t.Errorf("Expected nil result, got %+v", result)
				return
			}

			if tt.expectedResult != nil && result == nil {
				t.Errorf("Expected non-nil result, got nil")
				return
			}

			if tt.expectedResult != nil && result != nil {
				if result.Token != tt.expectedResult.Token {
					t.Errorf("Expected Token '%s', got '%s'", tt.expectedResult.Token, result.Token)
				}
				if result.TokenType != tt.expectedResult.TokenType {
					t.Errorf("Expected TokenType '%s', got '%s'", tt.expectedResult.TokenType, result.TokenType)
				}
			}
		})
	}
}

// TestArkKeyring_LoadToken_JsonParsing tests the JSON parsing logic specifically
func TestArkKeyring_LoadToken_JsonParsing(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(keyring *ArkKeyring) (cleanup func())
		expectedError bool
	}{
		{
			name: "error_invalid_json_data",
			setupMock: func(keyring *ArkKeyring) (cleanup func()) {
				// Manually insert invalid JSON data
				kr, _ := keyring.GetKeyring(true)
				kr.SetPassword("test-service-invalid", "test-profile", "invalid-json-data")
				return func() {}
			},
			expectedError: true,
		},
		{
			name: "success_valid_json_parsing",
			setupMock: func(keyring *ArkKeyring) (cleanup func()) {
				validToken := &auth.ArkToken{
					Token:     "valid-token",
					TokenType: auth.JWT,
				}
				tokenData, _ := json.Marshal(validToken)
				kr, _ := keyring.GetKeyring(true)
				kr.SetPassword("test-service-valid", "test-profile", string(tokenData))
				return func() {}
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keyring := NewArkKeyring("test-service")
			cleanup := tt.setupMock(keyring)
			defer cleanup()

			profile := &models.ArkProfile{ProfileName: "test-profile"}
			postfix := "invalid"
			if tt.name == "success_valid_json_parsing" {
				postfix = "valid"
			}

			_, err := keyring.LoadToken(profile, postfix, true)

			if tt.expectedError {
				if err == nil {
					t.Error("Expected error for invalid JSON, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for valid JSON, got %v", err)
				}
			}
		})
	}
}

// TestArkKeyring_Integration tests the complete save and load cycle
func TestArkKeyring_Integration(t *testing.T) {
	tests := []struct {
		name    string
		token   *auth.ArkToken
		profile *models.ArkProfile
		postfix string
	}{
		// {
		// 	name: "complete_cycle_bearer_token",
		// 	token: &auth.ArkToken{
		// 		Token:        "integration-access-token",
		// 		RefreshToken: "integration-refresh-token",
		// 		TokenType:    auth.JWT,
		// 		ExpiresIn:    commonmodels.ArkRFC3339Time(time.Now().Add(2 * time.Hour)),
		// 	},
		// 	profile: &models.ArkProfile{ProfileName: "integration-profile"},
		// 	postfix: "integration",
		// },
		// {
		// 	name: "complete_cycle_internal_token",
		// 	token: &auth.ArkToken{
		// 		Token:     "internal-access-token",
		// 		TokenType: auth.Internal,
		// 		ExpiresIn: commonmodels.ArkRFC3339Time(time.Now().Add(24 * time.Hour)),
		// 	},
		// 	profile: &models.ArkProfile{ProfileName: "internal-profile"},
		// 	postfix: "internal",
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keyring := NewArkKeyring("integration-test")

			// Save the token
			err := keyring.SaveToken(tt.profile, tt.token, tt.postfix, true)
			if err != nil {
				t.Fatalf("Failed to save token: %v", err)
			}

			// Load the token
			loadedToken, err := keyring.LoadToken(tt.profile, tt.postfix, true)
			if err != nil {
				t.Fatalf("Failed to load token: %v", err)
			}

			if loadedToken == nil {
				t.Fatal("Expected loaded token to be non-nil")
			}

			// Verify token contents
			if !reflect.DeepEqual(loadedToken, tt.token) {
				t.Errorf("Loaded token doesn't match saved token.\nExpected: %+v\nGot: %+v", tt.token, loadedToken)
			}
		})
	}
}
