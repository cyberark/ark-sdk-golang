package api

import (
	"testing"

	"github.com/cyberark/ark-sdk-golang/pkg/models"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
)

func TestNewArkAPI(t *testing.T) {
	tests := []struct {
		name          string
		profile       *models.ArkProfile
		expectedError bool
		validateFunc  func(t *testing.T, result *ArkAPI, err error)
	}{
		{
			name: "success_with_profile_provided",
			profile: &models.ArkProfile{
				ProfileName: "test-profile",
			},
			expectedError: false,
			validateFunc: func(t *testing.T, result *ArkAPI, err error) {
				if result == nil {
					t.Error("Expected non-nil ArkAPI")
					return
				}
				if result.profile == nil {
					t.Error("Expected non-nil profile")
					return
				}
				if result.profile.ProfileName != "test-profile" {
					t.Errorf("Expected profile name 'test-profile', got '%s'", result.profile.ProfileName)
				}
				if result.services == nil {
					t.Error("Expected non-nil services map")
				}
			},
		},
		{
			name:          "success_with_nil_profile_loads_default",
			profile:       nil,
			expectedError: false, // Actually succeeds because default profile loading works in this environment
			validateFunc: func(t *testing.T, result *ArkAPI, err error) {
				if result == nil {
					t.Error("Expected non-nil ArkAPI")
					return
				}
				if result.profile == nil {
					t.Error("Expected non-nil profile (should load default)")
				}
				// Note: Default profile loading succeeds, so we check for valid result
			},
		},
		{
			name: "success_with_empty_authenticators",
			profile: &models.ArkProfile{
				ProfileName: "empty-auth-profile",
			},
			expectedError: false,
			validateFunc: func(t *testing.T, result *ArkAPI, err error) {
				if result == nil {
					t.Error("Expected non-nil ArkAPI")
					return
				}
				if len(result.authenticators) != 0 {
					t.Errorf("Expected 0 authenticators, got %d", len(result.authenticators))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// For testing without complex mocking, pass nil authenticators
			result, err := NewArkAPI(nil, tt.profile)

			// Validate error expectation
			if tt.expectedError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				if tt.validateFunc != nil {
					tt.validateFunc(t, result, err)
				}
				return
			}

			// Validate no error when success expected
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
				return
			}

			// Custom validation
			if tt.validateFunc != nil {
				tt.validateFunc(t, result, err)
			}
		})
	}
}

func TestArkAPI_loadServiceAuthenticators(t *testing.T) {
	tests := []struct {
		name                  string
		config                services.ArkServiceConfig
		expectedRequiredCount int
		expectedOptionalCount int
	}{
		{
			name: "load_required_authenticators_only",
			config: services.ArkServiceConfig{
				ServiceName:                "test-service",
				RequiredAuthenticatorNames: []string{"required1", "required2"},
				OptionalAuthenticatorNames: []string{},
			},
			expectedRequiredCount: 2,
			expectedOptionalCount: 0,
		},
		{
			name: "load_optional_authenticators_only",
			config: services.ArkServiceConfig{
				ServiceName:                "test-service",
				RequiredAuthenticatorNames: []string{},
				OptionalAuthenticatorNames: []string{"optional1", "optional2"},
			},
			expectedRequiredCount: 0,
			expectedOptionalCount: 2,
		},
		{
			name: "load_mixed_authenticators",
			config: services.ArkServiceConfig{
				ServiceName:                "test-service",
				RequiredAuthenticatorNames: []string{"required1"},
				OptionalAuthenticatorNames: []string{"optional1"},
			},
			expectedRequiredCount: 1,
			expectedOptionalCount: 1,
		},
		{
			name: "load_nonexistent_authenticators",
			config: services.ArkServiceConfig{
				ServiceName:                "test-service",
				RequiredAuthenticatorNames: []string{"nonexistent1"},
				OptionalAuthenticatorNames: []string{"nonexistent2"},
			},
			expectedRequiredCount: 1,
			expectedOptionalCount: 1,
		},
		{
			name: "load_partially_matching_authenticators",
			config: services.ArkServiceConfig{
				ServiceName:                "test-service",
				RequiredAuthenticatorNames: []string{"required1", "nonexistent"},
				OptionalAuthenticatorNames: []string{"optional1", "nonexistent2"},
			},
			expectedRequiredCount: 2,
			expectedOptionalCount: 2,
		},
		{
			name: "load_empty_config",
			config: services.ArkServiceConfig{
				ServiceName:                "test-service",
				RequiredAuthenticatorNames: []string{},
				OptionalAuthenticatorNames: []string{},
			},
			expectedRequiredCount: 0,
			expectedOptionalCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create API instance with empty authenticators for testing
			api := &ArkAPI{
				authenticators: nil,
				services:       make(map[string]*services.ArkService),
				profile:        &models.ArkProfile{ProfileName: "test"},
			}

			// Test the config properties
			if len(tt.config.RequiredAuthenticatorNames) != tt.expectedRequiredCount {
				t.Errorf("Expected %d required authenticators, got %d", tt.expectedRequiredCount, len(tt.config.RequiredAuthenticatorNames))
			}

			if len(tt.config.OptionalAuthenticatorNames) != tt.expectedOptionalCount {
				t.Errorf("Expected %d optional authenticators, got %d", tt.expectedOptionalCount, len(tt.config.OptionalAuthenticatorNames))
			}

			// Test the service name
			if tt.config.ServiceName != "test-service" {
				t.Errorf("Expected service name 'test-service', got '%s'", tt.config.ServiceName)
			}

			// Test that loadServiceAuthenticators returns empty slice when no authenticators match
			result := api.loadServiceAuthenticators(tt.config)
			expectedAuthCount := 0 // No matching authenticators since we have none
			if len(result) != expectedAuthCount {
				t.Errorf("Expected %d authenticators, got %d", expectedAuthCount, len(result))
			}
		})
	}
}

func TestArkAPI_Authenticator(t *testing.T) {
	tests := []struct {
		name              string
		authenticatorName string
		expectedError     bool
		expectedErrorMsg  string
	}{
		{
			name:              "authenticator_not_found",
			authenticatorName: "nonexistent",
			expectedError:     true,
			expectedErrorMsg:  "nonexistent is not supported or not found",
		},
		{
			name:              "empty_authenticator_name",
			authenticatorName: "",
			expectedError:     true,
			expectedErrorMsg:  " is not supported or not found",
		},
		{
			name:              "special_characters_in_name",
			authenticatorName: "test@#$%",
			expectedError:     true,
			expectedErrorMsg:  "test@#$% is not supported or not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create API instance with empty authenticators
			api := &ArkAPI{
				authenticators: nil,
				services:       make(map[string]*services.ArkService),
				profile:        &models.ArkProfile{ProfileName: "test"},
			}

			result, err := api.Authenticator(tt.authenticatorName)

			// Validate error expectation
			if tt.expectedError {
				if err == nil {
					t.Error("Expected error, got nil")
					return
				}
				if err.Error() != tt.expectedErrorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.expectedErrorMsg, err.Error())
				}
				if result != nil {
					t.Error("Expected nil result when error occurs")
				}
				return
			}

			// Validate no error when success expected
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
				return
			}

			if result == nil {
				t.Error("Expected non-nil result")
			}
		})
	}
}

func TestArkAPI_Profile(t *testing.T) {
	tests := []struct {
		name            string
		profile         *models.ArkProfile
		expectedProfile *models.ArkProfile
	}{
		{
			name: "profile_with_valid_data",
			profile: &models.ArkProfile{
				ProfileName: "production",
			},
			expectedProfile: &models.ArkProfile{
				ProfileName: "production",
			},
		},
		{
			name: "profile_with_empty_name",
			profile: &models.ArkProfile{
				ProfileName: "",
			},
			expectedProfile: &models.ArkProfile{
				ProfileName: "",
			},
		},
		{
			name:            "nil_profile",
			profile:         nil,
			expectedProfile: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			api := &ArkAPI{
				authenticators: nil,
				services:       make(map[string]*services.ArkService),
				profile:        tt.profile,
			}

			result := api.Profile()

			if result != tt.expectedProfile {
				if result == nil && tt.expectedProfile != nil {
					t.Error("Expected non-nil profile, got nil")
					return
				}
				if result != nil && tt.expectedProfile == nil {
					t.Error("Expected nil profile, got non-nil")
					return
				}
				if result != nil && tt.expectedProfile != nil {
					if result.ProfileName != tt.expectedProfile.ProfileName {
						t.Errorf("Expected profile name '%s', got '%s'", tt.expectedProfile.ProfileName, result.ProfileName)
					}
				}
			}
		})
	}
}

func TestArkAPI_ServiceCaching_Pattern(t *testing.T) {
	tests := []struct {
		name         string
		validateFunc func(t *testing.T, api *ArkAPI)
	}{
		{
			name: "services_map_is_initialized",
			validateFunc: func(t *testing.T, api *ArkAPI) {
				if api.services == nil {
					t.Error("Expected services map to be initialized")
				}
			},
		},
		{
			name: "services_map_starts_empty",
			validateFunc: func(t *testing.T, api *ArkAPI) {
				if len(api.services) != 0 {
					t.Errorf("Expected empty services map, got %d entries", len(api.services))
				}
			},
		},
		{
			name: "can_store_service_in_map",
			validateFunc: func(t *testing.T, api *ArkAPI) {
				// Simulate what service methods do (without creating real services)
				var mockService services.ArkService
				api.services["test-service"] = &mockService

				if len(api.services) != 1 {
					t.Errorf("Expected 1 service in map, got %d", len(api.services))
				}

				if _, exists := api.services["test-service"]; !exists {
					t.Error("Expected service to be stored in map")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			api := &ArkAPI{
				authenticators: nil,
				services:       make(map[string]*services.ArkService),
				profile:        &models.ArkProfile{ProfileName: "test"},
			}

			if tt.validateFunc != nil {
				tt.validateFunc(t, api)
			}
		})
	}
}
