package actions

import (
	"reflect"
	"testing"

	"github.com/cyberark/ark-sdk-golang/pkg/actions/testutils"
	"github.com/cyberark/ark-sdk-golang/pkg/models"
	authmodels "github.com/cyberark/ark-sdk-golang/pkg/models/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/profiles"
	"github.com/spf13/cobra"
)

func TestNewArkConfigureAction(t *testing.T) {
	tests := []struct {
		name         string
		setupLoader  func() profiles.ProfileLoader
		validateFunc func(t *testing.T, action *ArkConfigureAction)
	}{
		{
			name: "success_creates_action_with_profile_loader",
			setupLoader: func() profiles.ProfileLoader {
				return testutils.NewMockProfileLoader()
			},
			validateFunc: func(t *testing.T, action *ArkConfigureAction) {
				if action == nil {
					t.Error("Expected action to be created, got nil")
					return
				}
				if action.ArkBaseAction == nil {
					t.Error("Expected ArkBaseAction to be initialized")
				}
				if action.profilesLoader == nil {
					t.Error("Expected profilesLoader to be set")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			loader := tt.setupLoader()
			action := NewArkConfigureAction(&loader)

			if tt.validateFunc != nil {
				tt.validateFunc(t, action)
			}
		})
	}
}

func TestArkConfigureAction_DefineAction(t *testing.T) {
	tests := []struct {
		name         string
		setupLoader  func() profiles.ProfileLoader
		validateFunc func(t *testing.T, cmd *cobra.Command, confCmd *cobra.Command)
	}{
		{
			name: "success_adds_configure_command_with_flags",
			setupLoader: func() profiles.ProfileLoader {
				return testutils.NewMockProfileLoader()
			},
			validateFunc: func(t *testing.T, cmd *cobra.Command, confCmd *cobra.Command) {
				// Verify configure command was added
				if confCmd == nil {
					t.Error("Expected configure command to be added")
					return
				}

				if confCmd.Use != "configure" {
					t.Errorf("Expected command use 'configure', got '%s'", confCmd.Use)
				}

				if confCmd.Short != "Configure the CLI" {
					t.Errorf("Expected command short description 'Configure the CLI', got '%s'", confCmd.Short)
				}

				if confCmd.Run == nil {
					t.Error("Expected run function to be set")
				}

				if confCmd.PersistentPreRun == nil {
					t.Error("Expected persistent pre-run function to be set")
				}
			},
		},
		{
			name: "success_adds_profile_flags",
			setupLoader: func() profiles.ProfileLoader {
				return testutils.NewMockProfileLoader()
			},
			validateFunc: func(t *testing.T, cmd *cobra.Command, confCmd *cobra.Command) {
				// Check for common flags from ArkBaseAction
				commonFlags := []string{
					"raw", "silent", "allow-output", "verbose",
					"logger-style", "log-level", "disable-cert-verification", "trusted-cert",
				}

				for _, flagName := range commonFlags {
					flag := confCmd.PersistentFlags().Lookup(flagName)
					if flag == nil {
						t.Errorf("Expected common flag '%s' to be present", flagName)
					}
				}
			},
		},
		{
			name: "edge_case_handles_nil_profile_loader",
			setupLoader: func() profiles.ProfileLoader {
				return nil
			},
			validateFunc: func(t *testing.T, cmd *cobra.Command, confCmd *cobra.Command) {
				// Should not panic even with nil profile loader
				if confCmd == nil {
					t.Error("Expected configure command to be added even with nil loader")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			loader := tt.setupLoader()
			action := NewArkConfigureAction(&loader)
			cmd := &cobra.Command{}

			// Execute DefineAction - should not panic
			action.DefineAction(cmd)

			// Find the configure command
			var confCmd *cobra.Command
			for _, subCmd := range cmd.Commands() {
				if subCmd.Use == "configure" {
					confCmd = subCmd
					break
				}
			}

			if tt.validateFunc != nil {
				tt.validateFunc(t, cmd, confCmd)
			}
		})
	}
}

func TestArkConfigureAction_runSilentConfigureAction(t *testing.T) {
	tests := []struct {
		name            string
		setupLoader     func() profiles.ProfileLoader
		setupFlags      func(cmd *cobra.Command)
		expectedProfile *models.ArkProfile
		expectedError   bool
		validateFunc    func(t *testing.T, profile *models.ArkProfile, err error)
	}{
		{
			name: "success_creates_new_profile_with_default_name",
			setupLoader: func() profiles.ProfileLoader {
				mock := testutils.NewMockProfileLoader()
				mock.LoadProfileFunc = func(name string) (*models.ArkProfile, error) {
					return nil, nil // Profile not found
				}
				return mock
			},
			setupFlags: func(cmd *cobra.Command) {
				// Define the profile-name flag that the function expects
				_ = cmd.Flags().String("profile-name", "", "Profile name")
			},
			expectedError: false,
			validateFunc: func(t *testing.T, profile *models.ArkProfile, err error) {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
					return
				}
				if profile == nil {
					t.Error("Expected profile to be created")
					return
				}
				if profile.AuthProfiles == nil {
					t.Error("Expected AuthProfiles to be initialized")
				}
			},
		},
		{
			name: "success_loads_existing_profile",
			setupLoader: func() profiles.ProfileLoader {
				existingProfile := &models.ArkProfile{
					ProfileName:  "test-profile",
					AuthProfiles: map[string]*authmodels.ArkAuthProfile{},
				}
				mock := testutils.NewMockProfileLoader()
				mock.LoadProfileFunc = func(name string) (*models.ArkProfile, error) {
					return existingProfile, nil
				}
				return mock
			},
			setupFlags: func(cmd *cobra.Command) {
				_ = cmd.Flags().String("profile-name", "", "Profile name")
				_ = cmd.Flags().Set("profile-name", "test-profile")
			},
			expectedError: false,
			validateFunc: func(t *testing.T, profile *models.ArkProfile, err error) {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
					return
				}
				if profile == nil {
					t.Error("Expected profile to be loaded")
					return
				}
				if profile.ProfileName != "test-profile" {
					t.Errorf("Expected profile name 'test-profile', got '%s'", profile.ProfileName)
				}
			},
		},
		{
			name: "success_merges_flag_values",
			setupLoader: func() profiles.ProfileLoader {
				mock := testutils.NewMockProfileLoader()
				mock.LoadProfileFunc = func(name string) (*models.ArkProfile, error) {
					return nil, nil // New profile
				}
				return mock
			},
			setupFlags: func(cmd *cobra.Command) {
				_ = cmd.Flags().String("profile-name", "", "Profile name")
				_ = cmd.Flags().String("tenant-url", "", "Tenant URL")
				_ = cmd.Flags().Set("profile-name", "custom-profile")
				_ = cmd.Flags().Set("tenant-url", "https://example.com")
			},
			expectedError: false,
			validateFunc: func(t *testing.T, profile *models.ArkProfile, err error) {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
					return
				}
				if profile.ProfileName != "custom-profile" {
					t.Errorf("Expected profile name 'custom-profile', got '%s'", profile.ProfileName)
				}
			},
		},
		{
			name: "edge_case_handles_flag_parsing_errors",
			setupLoader: func() profiles.ProfileLoader {
				mock := testutils.NewMockProfileLoader()
				mock.LoadProfileFunc = func(name string) (*models.ArkProfile, error) {
					return nil, nil
				}
				return mock
			},
			setupFlags: func(cmd *cobra.Command) {
				// Define the profile-name flag but don't set a value to test default behavior
				_ = cmd.Flags().String("profile-name", "", "Profile name")
			},
			expectedError: false, // Function handles missing flag values gracefully by using defaults
			validateFunc: func(t *testing.T, profile *models.ArkProfile, err error) {
				// Should create a profile even when no flag values are set
				if profile == nil {
					t.Error("Expected profile to be created despite missing flag values")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			loader := tt.setupLoader()
			action := NewArkConfigureAction(&loader)
			cmd := &cobra.Command{}

			if tt.setupFlags != nil {
				tt.setupFlags(cmd)
			}

			profile, err := action.runSilentConfigureAction(cmd, []string{})

			// Validate error expectation
			if tt.expectedError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if tt.validateFunc != nil {
				tt.validateFunc(t, profile, err)
			}
		})
	}
}

func TestArkConfigureAction_StructFields(t *testing.T) {
	tests := []struct {
		name         string
		validateFunc func(t *testing.T, action *ArkConfigureAction)
	}{
		{
			name: "success_struct_has_expected_fields",
			validateFunc: func(t *testing.T, action *ArkConfigureAction) {
				actionValue := reflect.ValueOf(action).Elem()
				actionType := actionValue.Type()

				expectedFields := []string{"ArkBaseAction", "profilesLoader"}
				actualFields := make([]string, actionType.NumField())

				for i := 0; i < actionType.NumField(); i++ {
					actualFields[i] = actionType.Field(i).Name
				}

				for _, expectedField := range expectedFields {
					found := false
					for _, actualField := range actualFields {
						if actualField == expectedField {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("Expected field '%s' not found in struct", expectedField)
					}
				}
			},
		},
		{
			name: "success_profilesloader_field_has_correct_type",
			validateFunc: func(t *testing.T, action *ArkConfigureAction) {
				actionValue := reflect.ValueOf(action).Elem()
				profilesLoaderField := actionValue.FieldByName("profilesLoader")

				if !profilesLoaderField.IsValid() {
					t.Error("profilesLoader field not found")
					return
				}

				expectedType := "*profiles.ProfileLoader"
				actualType := profilesLoaderField.Type().String()
				if actualType != expectedType {
					t.Errorf("Expected profilesLoader field type '%s', got '%s'", expectedType, actualType)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			loader := testutils.NewMockProfileLoader()
			loaderInterface := profiles.ProfileLoader(loader)
			action := NewArkConfigureAction(&loaderInterface)

			if tt.validateFunc != nil {
				tt.validateFunc(t, action)
			}
		})
	}
}

func TestArkConfigureAction_ArkActionInterface(t *testing.T) {
	tests := []struct {
		name         string
		validateFunc func(t *testing.T)
	}{
		{
			name: "success_implements_arkaction_interface",
			validateFunc: func(t *testing.T) {
				loader := testutils.NewMockProfileLoader()
				loaderInterface := profiles.ProfileLoader(loader)
				action := NewArkConfigureAction(&loaderInterface)

				// This will cause a compile error if ArkConfigureAction doesn't implement ArkAction
				var _ ArkAction = action

				// Verify the DefineAction method exists and can be called
				cmd := &cobra.Command{}
				action.DefineAction(cmd)

				// Verify a configure command was added
				found := false
				for _, subCmd := range cmd.Commands() {
					if subCmd.Use == "configure" {
						found = true
						break
					}
				}
				if !found {
					t.Error("Expected configure command to be added to parent command")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.validateFunc != nil {
				tt.validateFunc(t)
			}
		})
	}
}
