package models

import (
	"encoding/json"
	"testing"

	"github.com/cyberark/ark-sdk-golang/pkg/models/auth"
)

func TestArkProfile_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name            string
		jsonData        string
		expectedError   bool
		expectedProfile *ArkProfile
		validateFunc    func(t *testing.T, profile *ArkProfile)
	}{
		{
			name: "success_complete_profile_with_identity_auth",
			jsonData: `{
				"profile_name": "test-profile",
				"profile_description": "Test profile description",
				"auth_profiles": {
					"default": {
						"username": "test@example.com",
						"auth_method": "identity",
						"auth_method_settings": {
							"identity_mfa_method": "pf",
							"identity_mfa_interactive": true,
							"identity_application": "test-app",
							"identity_url": "https://identity.example.com",
							"identity_tenant_subdomain": "test-tenant"
						}
					}
				}
			}`,
			expectedError: false,
			expectedProfile: &ArkProfile{
				ProfileName:        "test-profile",
				ProfileDescription: "Test profile description",
				AuthProfiles: map[string]*auth.ArkAuthProfile{
					"default": {
						Username:   "test@example.com",
						AuthMethod: auth.Identity,
						AuthMethodSettings: &auth.IdentityArkAuthMethodSettings{
							IdentityMFAMethod:       "pf",
							IdentityMFAInteractive:  true,
							IdentityURL:             "https://identity.example.com",
							IdentityTenantSubdomain: "test-tenant",
						},
					},
				},
			},
		},
		{
			name: "success_multiple_auth_profiles",
			jsonData: `{
				"profile_name": "multi-profile",
				"profile_description": "Profile with multiple auth methods",
				"auth_profiles": {
					"identity": {
						"username": "user1@example.com",
						"auth_method": "identity",
						"auth_method_settings": {
							"identity_mfa_method": "sms",
							"identity_mfa_interactive": false,
							"identity_application": "app1",
							"identity_url": "https://id1.example.com",
							"identity_tenant_subdomain": "tenant1"
						}
					},
					"direct": {
						"username": "user2@example.com",
						"auth_method": "direct",
						"auth_method_settings": {
							"direct_url": "https://direct.example.com"
						}
					}
				}
			}`,
			expectedError: false,
			validateFunc: func(t *testing.T, profile *ArkProfile) {
				if profile.ProfileName != "multi-profile" {
					t.Errorf("Expected ProfileName 'multi-profile', got '%s'", profile.ProfileName)
				}
				if len(profile.AuthProfiles) != 2 {
					t.Errorf("Expected 2 auth profiles, got %d", len(profile.AuthProfiles))
				}
				if profile.AuthProfiles["identity"].AuthMethod != auth.Identity {
					t.Errorf("Expected identity auth method, got %s", profile.AuthProfiles["identity"].AuthMethod)
				}
				if profile.AuthProfiles["direct"].AuthMethod != auth.Direct {
					t.Errorf("Expected direct auth method, got %s", profile.AuthProfiles["direct"].AuthMethod)
				}
			},
		},
		{
			name: "success_empty_auth_profiles",
			jsonData: `{
				"profile_name": "empty-profile",
				"profile_description": "Profile with no auth profiles",
				"auth_profiles": {}
			}`,
			expectedError: false,
			expectedProfile: &ArkProfile{
				ProfileName:        "empty-profile",
				ProfileDescription: "Profile with no auth profiles",
				AuthProfiles:       map[string]*auth.ArkAuthProfile{},
			},
		},
		{
			name: "success_minimal_profile",
			jsonData: `{
				"profile_name": "minimal",
				"profile_description": "Minimal profile",
				"auth_profiles": {
					"simple": {
						"username": "simple@example.com",
						"auth_method": "default",
						"auth_method_settings": {}
					}
				}
			}`,
			expectedError: false,
			validateFunc: func(t *testing.T, profile *ArkProfile) {
				if profile.ProfileName != "minimal" {
					t.Errorf("Expected ProfileName 'minimal', got '%s'", profile.ProfileName)
				}
				if profile.AuthProfiles["simple"].Username != "simple@example.com" {
					t.Errorf("Expected username 'simple@example.com', got '%s'", profile.AuthProfiles["simple"].Username)
				}
				if profile.AuthProfiles["simple"].AuthMethod != auth.Default {
					t.Errorf("Expected default auth method, got %s", profile.AuthProfiles["simple"].AuthMethod)
				}
			},
		},
		{
			name:          "error_invalid_json",
			jsonData:      `{"profile_name": "test", "invalid_json"`,
			expectedError: true,
		},
		{
			name: "error_missing_required_fields",
			jsonData: `{
				"auth_profiles": {}
			}`,
			expectedError: false, // JSON unmarshaling succeeds, but fields will be empty
			expectedProfile: &ArkProfile{
				ProfileName:        "",
				ProfileDescription: "",
				AuthProfiles:       map[string]*auth.ArkAuthProfile{},
			},
		},
		{
			name: "error_invalid_auth_method",
			jsonData: `{
				"profile_name": "invalid-auth",
				"profile_description": "Profile with invalid auth method",
				"auth_profiles": {
					"invalid": {
						"username": "user@example.com",
						"auth_method": "invalid_method",
						"auth_method_settings": {}
					}
				}
			}`,
			expectedError: true, // This should fail in auth profile unmarshaling
		},
		{
			name: "error_malformed_auth_profile",
			jsonData: `{
				"profile_name": "malformed",
				"profile_description": "Profile with malformed auth profile",
				"auth_profiles": {
					"bad": "not_an_object"
				}
			}`,
			expectedError: true,
		},
		{
			name: "edge_case_null_auth_profiles",
			jsonData: `{
				"profile_name": "null-auth",
				"profile_description": "Profile with null auth profiles",
				"auth_profiles": null
			}`,
			expectedError: false,
			expectedProfile: &ArkProfile{
				ProfileName:        "null-auth",
				ProfileDescription: "Profile with null auth profiles",
				AuthProfiles:       map[string]*auth.ArkAuthProfile{},
			},
		},
		{
			name: "edge_case_special_characters_in_names",
			jsonData: `{
				"profile_name": "special-chars_123!@#",
				"profile_description": "Profile with special characters: üñîçödé & symbols",
				"auth_profiles": {
					"key-with_special.chars": {
						"username": "user+test@example.com",
						"auth_method": "default",
						"auth_method_settings": {}
					}
				}
			}`,
			expectedError: false,
			validateFunc: func(t *testing.T, profile *ArkProfile) {
				if profile.ProfileName != "special-chars_123!@#" {
					t.Errorf("Expected special character profile name, got '%s'", profile.ProfileName)
				}
				if _, exists := profile.AuthProfiles["key-with_special.chars"]; !exists {
					t.Error("Expected auth profile with special characters in key")
				}
			},
		},
		{
			name: "error_invalid_auth_profile_json_structure",
			jsonData: `{
				"profile_name": "invalid-structure",
				"profile_description": "Profile with invalid auth profile JSON structure",
				"auth_profiles": {
					"malformed": {
						"username": "user@example.com",
						"auth_method": "identity",
						"auth_method_settings": {
							"identity_mfa_method": "invalid_json_value",
							"identity_mfa_interactive": "not_a_boolean"
						}
					}
				}
			}`,
			expectedError: true, // This specifically tests line 77: json.Unmarshal(rawMessage, &authProfile)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var profile ArkProfile
			err := json.Unmarshal([]byte(tt.jsonData), &profile)

			// Validate error expectation
			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				return // Don't check result if we expected an error
			}

			// Validate no error when success expected
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
				return
			}

			// Use custom validation if provided
			if tt.validateFunc != nil {
				tt.validateFunc(t, &profile)
				return
			}

			// Validate expected profile if provided
			if tt.expectedProfile != nil {
				if profile.ProfileName != tt.expectedProfile.ProfileName {
					t.Errorf("Expected ProfileName '%s', got '%s'", tt.expectedProfile.ProfileName, profile.ProfileName)
				}
				if profile.ProfileDescription != tt.expectedProfile.ProfileDescription {
					t.Errorf("Expected ProfileDescription '%s', got '%s'", tt.expectedProfile.ProfileDescription, profile.ProfileDescription)
				}

				// Check auth profiles map
				if len(profile.AuthProfiles) != len(tt.expectedProfile.AuthProfiles) {
					t.Errorf("Expected %d auth profiles, got %d", len(tt.expectedProfile.AuthProfiles), len(profile.AuthProfiles))
				}

				for key, expectedAuth := range tt.expectedProfile.AuthProfiles {
					actualAuth, exists := profile.AuthProfiles[key]
					if !exists {
						t.Errorf("Expected auth profile with key '%s' not found", key)
						continue
					}

					if actualAuth.Username != expectedAuth.Username {
						t.Errorf("Expected username '%s', got '%s'", expectedAuth.Username, actualAuth.Username)
					}
					if actualAuth.AuthMethod != expectedAuth.AuthMethod {
						t.Errorf("Expected auth method '%s', got '%s'", expectedAuth.AuthMethod, actualAuth.AuthMethod)
					}

					// Note: AuthMethodSettings comparison is complex due to interface types
					// We rely on the auth package's own unmarshaling tests for detailed validation
					if actualAuth.AuthMethodSettings == nil && expectedAuth.AuthMethodSettings != nil {
						t.Error("Expected auth method settings, got nil")
					}
					if actualAuth.AuthMethodSettings != nil && expectedAuth.AuthMethodSettings == nil {
						t.Error("Expected nil auth method settings, got non-nil")
					}
				}
			}
		})
	}
}

// TestArkProfile_UnmarshalJSON_Integration tests the integration with auth package
func TestArkProfile_UnmarshalJSON_Integration(t *testing.T) {
	tests := []struct {
		name          string
		jsonData      string
		expectedError bool
		validateAuth  func(t *testing.T, authProfile *auth.ArkAuthProfile)
	}{
		{
			name: "integration_identity_service_user_auth",
			jsonData: `{
				"profile_name": "integration-test",
				"profile_description": "Integration test profile",
				"auth_profiles": {
					"service": {
						"username": "service@example.com",
						"auth_method": "identity_service_user",
						"auth_method_settings": {
							"identity_url": "https://identity.example.com",
							"identity_tenant_subdomain": "test-tenant"
						}
					}
				}
			}`,
			expectedError: false,
			validateAuth: func(t *testing.T, authProfile *auth.ArkAuthProfile) {
				if authProfile.AuthMethod != auth.IdentityServiceUser {
					t.Errorf("Expected IdentityServiceUser auth method, got %s", authProfile.AuthMethod)
				}
				settings, ok := authProfile.AuthMethodSettings.(*auth.IdentityServiceUserArkAuthMethodSettings)
				if !ok {
					t.Errorf("Expected IdentityServiceUserArkAuthMethodSettings, got %T", authProfile.AuthMethodSettings)
					return
				}
				if settings.IdentityURL != "https://identity.example.com" {
					t.Errorf("Expected identity URL 'https://identity.example.com', got '%s'", settings.IdentityURL)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var profile ArkProfile
			err := json.Unmarshal([]byte(tt.jsonData), &profile)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
				return
			}

			// Validate auth profile integration
			if tt.validateAuth != nil && len(profile.AuthProfiles) > 0 {
				for _, authProfile := range profile.AuthProfiles {
					tt.validateAuth(t, authProfile)
					break // Test first auth profile
				}
			}
		})
	}
}

// TestArkProfile_RoundTrip tests JSON marshaling and unmarshaling round trip
func TestArkProfile_RoundTrip(t *testing.T) {
	originalProfile := &ArkProfile{
		ProfileName:        "roundtrip-test",
		ProfileDescription: "Test for JSON round trip",
		AuthProfiles: map[string]*auth.ArkAuthProfile{
			"test": {
				Username:           "roundtrip@example.com",
				AuthMethod:         auth.Default,
				AuthMethodSettings: &auth.DefaultArkAuthMethodSettings{},
			},
		},
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(originalProfile)
	if err != nil {
		t.Fatalf("Failed to marshal profile: %v", err)
	}

	// Unmarshal back to struct
	var unmarshaledProfile ArkProfile
	err = json.Unmarshal(jsonData, &unmarshaledProfile)
	if err != nil {
		t.Fatalf("Failed to unmarshal profile: %v", err)
	}

	// Validate round trip
	if unmarshaledProfile.ProfileName != originalProfile.ProfileName {
		t.Errorf("ProfileName mismatch after round trip: expected '%s', got '%s'",
			originalProfile.ProfileName, unmarshaledProfile.ProfileName)
	}
	if unmarshaledProfile.ProfileDescription != originalProfile.ProfileDescription {
		t.Errorf("ProfileDescription mismatch after round trip: expected '%s', got '%s'",
			originalProfile.ProfileDescription, unmarshaledProfile.ProfileDescription)
	}
	if len(unmarshaledProfile.AuthProfiles) != len(originalProfile.AuthProfiles) {
		t.Errorf("AuthProfiles count mismatch after round trip: expected %d, got %d",
			len(originalProfile.AuthProfiles), len(unmarshaledProfile.AuthProfiles))
	}
}
