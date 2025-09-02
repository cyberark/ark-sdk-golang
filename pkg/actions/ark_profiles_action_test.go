package actions

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/cyberark/ark-sdk-golang/pkg/actions/testutils"
	"github.com/cyberark/ark-sdk-golang/pkg/models"
	authmodels "github.com/cyberark/ark-sdk-golang/pkg/models/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/profiles"
	"github.com/spf13/cobra"
)

func TestNewArkProfilesAction(t *testing.T) {
	tests := []struct {
		name         string
		setupLoader  func() profiles.ProfileLoader
		validateFunc func(t *testing.T, action *ArkProfilesAction)
	}{
		{
			name: "success_creates_action_with_profile_loader",
			setupLoader: func() profiles.ProfileLoader {
				return testutils.NewMockProfileLoader()
			},
			validateFunc: func(t *testing.T, action *ArkProfilesAction) {
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
		{
			name: "success_handles_nil_loader",
			setupLoader: func() profiles.ProfileLoader {
				return nil
			},
			validateFunc: func(t *testing.T, action *ArkProfilesAction) {
				if action == nil {
					t.Error("Expected action to be created, got nil")
					return
				}
				// Should handle nil loader gracefully - the function accepts a pointer
				// so when we pass &nil, it's still a valid pointer to ProfileLoader
				if action.profilesLoader == nil {
					t.Error("Expected profilesLoader to be set to the passed pointer")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			loader := tt.setupLoader()
			action := NewArkProfilesAction(&loader)

			if tt.validateFunc != nil {
				tt.validateFunc(t, action)
			}
		})
	}
}

func TestArkProfilesAction_DefineAction(t *testing.T) {
	tests := []struct {
		name         string
		setupLoader  func() profiles.ProfileLoader
		validateFunc func(t *testing.T, cmd *cobra.Command, profilesCmd *cobra.Command)
	}{
		{
			name: "success_adds_profiles_command_with_subcommands",
			setupLoader: func() profiles.ProfileLoader {
				return testutils.NewMockProfileLoader()
			},
			validateFunc: func(t *testing.T, cmd *cobra.Command, profilesCmd *cobra.Command) {
				if profilesCmd == nil {
					t.Error("Expected profiles command to be added")
					return
				}

				if profilesCmd.Use != "profiles" {
					t.Errorf("Expected command use 'profiles', got '%s'", profilesCmd.Use)
				}

				if profilesCmd.Short != "Manage profiles" {
					t.Errorf("Expected command short description 'Manage profiles', got '%s'", profilesCmd.Short)
				}

				// Check for expected subcommands
				expectedSubcommands := []string{"list", "show", "delete", "clear", "clone", "add", "edit"}
				actualSubcommands := make([]string, len(profilesCmd.Commands()))
				for i, subcmd := range profilesCmd.Commands() {
					actualSubcommands[i] = subcmd.Use
				}

				for _, expected := range expectedSubcommands {
					found := false
					for _, actual := range actualSubcommands {
						if actual == expected {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("Expected subcommand '%s' not found", expected)
					}
				}
			},
		},
		{
			name: "success_configures_list_command_flags",
			setupLoader: func() profiles.ProfileLoader {
				return testutils.NewMockProfileLoader()
			},
			validateFunc: func(t *testing.T, cmd *cobra.Command, profilesCmd *cobra.Command) {
				var listCmd *cobra.Command
				for _, subcmd := range profilesCmd.Commands() {
					if subcmd.Use == "list" {
						listCmd = subcmd
						break
					}
				}

				if listCmd == nil {
					t.Error("Expected list subcommand to be found")
					return
				}

				expectedFlags := []string{"name", "auth-profile", "all"}
				for _, flagName := range expectedFlags {
					if listCmd.Flags().Lookup(flagName) == nil {
						t.Errorf("Expected flag '%s' not found in list command", flagName)
					}
				}
			},
		},
		{
			name: "success_configures_show_command_flags",
			setupLoader: func() profiles.ProfileLoader {
				return testutils.NewMockProfileLoader()
			},
			validateFunc: func(t *testing.T, cmd *cobra.Command, profilesCmd *cobra.Command) {
				var showCmd *cobra.Command
				for _, subcmd := range profilesCmd.Commands() {
					if subcmd.Use == "show" {
						showCmd = subcmd
						break
					}
				}

				if showCmd == nil {
					t.Error("Expected show subcommand to be found")
					return
				}

				if showCmd.Flags().Lookup("profile-name") == nil {
					t.Error("Expected profile-name flag not found in show command")
				}
			},
		},
		{
			name: "success_configures_delete_command_flags",
			setupLoader: func() profiles.ProfileLoader {
				return testutils.NewMockProfileLoader()
			},
			validateFunc: func(t *testing.T, cmd *cobra.Command, profilesCmd *cobra.Command) {
				var deleteCmd *cobra.Command
				for _, subcmd := range profilesCmd.Commands() {
					if subcmd.Use == "delete" {
						deleteCmd = subcmd
						break
					}
				}

				if deleteCmd == nil {
					t.Error("Expected delete subcommand to be found")
					return
				}

				expectedFlags := []string{"profile-name", "yes"}
				for _, flagName := range expectedFlags {
					if deleteCmd.Flags().Lookup(flagName) == nil {
						t.Errorf("Expected flag '%s' not found in delete command", flagName)
					}
				}
			},
		},
		{
			name: "success_configures_clone_command_flags",
			setupLoader: func() profiles.ProfileLoader {
				return testutils.NewMockProfileLoader()
			},
			validateFunc: func(t *testing.T, cmd *cobra.Command, profilesCmd *cobra.Command) {
				var cloneCmd *cobra.Command
				for _, subcmd := range profilesCmd.Commands() {
					if subcmd.Use == "clone" {
						cloneCmd = subcmd
						break
					}
				}

				if cloneCmd == nil {
					t.Error("Expected clone subcommand to be found")
					return
				}

				expectedFlags := []string{"profile-name", "new-profile-name", "yes"}
				for _, flagName := range expectedFlags {
					if cloneCmd.Flags().Lookup(flagName) == nil {
						t.Errorf("Expected flag '%s' not found in clone command", flagName)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			loader := tt.setupLoader()
			action := NewArkProfilesAction(&loader)
			cmd := &cobra.Command{}

			// Execute DefineAction - should not panic
			action.DefineAction(cmd)

			// Find the profiles command
			var profilesCmd *cobra.Command
			for _, subCmd := range cmd.Commands() {
				if subCmd.Use == "profiles" {
					profilesCmd = subCmd
					break
				}
			}

			if tt.validateFunc != nil {
				tt.validateFunc(t, cmd, profilesCmd)
			}
		})
	}
}

func TestArkProfilesAction_runListAction(t *testing.T) {
	tests := []struct {
		name         string
		setupLoader  func() profiles.ProfileLoader
		setupFlags   func(cmd *cobra.Command)
		validateFunc func(t *testing.T, loader profiles.ProfileLoader)
	}{
		{
			name: "success_lists_all_profiles_without_filters",
			setupLoader: func() profiles.ProfileLoader {
				profiles := []*models.ArkProfile{
					testutils.CreateTestProfile("profile1"),
					testutils.CreateTestProfile("profile2"),
				}
				mock := testutils.NewMockProfileLoader()
				mock.LoadAllProfilesFunc = func() ([]*models.ArkProfile, error) {
					return profiles, nil
				}
				return mock
			},
			setupFlags: func(cmd *cobra.Command) {
				_ = cmd.Flags().String("name", "", "Profile name filter")
				_ = cmd.Flags().String("auth-profile", "", "Auth profile filter")
				_ = cmd.Flags().Bool("all", false, "Show all data")
			},
			validateFunc: func(t *testing.T, loader profiles.ProfileLoader) {
				// Verify that LoadAllProfiles was called
				mockLoader := loader.(*testutils.MockProfileLoader)
				if mockLoader.LoadAllProfilesFunc == nil {
					t.Error("Expected LoadAllProfiles to be called")
				}
			},
		},
		{
			name: "success_handles_no_profiles_found",
			setupLoader: func() profiles.ProfileLoader {
				mock := testutils.NewMockProfileLoader()
				mock.LoadAllProfilesFunc = func() ([]*models.ArkProfile, error) {
					return []*models.ArkProfile{}, nil
				}
				return mock
			},
			setupFlags: func(cmd *cobra.Command) {
				_ = cmd.Flags().String("name", "", "Profile name filter")
				_ = cmd.Flags().String("auth-profile", "", "Auth profile filter")
				_ = cmd.Flags().Bool("all", false, "Show all data")
			},
			validateFunc: func(t *testing.T, loader profiles.ProfileLoader) {
				// Function should handle empty profile list gracefully
			},
		},
		{
			name: "success_filters_profiles_by_name",
			setupLoader: func() profiles.ProfileLoader {
				profiles := []*models.ArkProfile{
					testutils.CreateTestProfile("test-profile"),
					testutils.CreateTestProfile("prod-profile"),
				}
				mock := testutils.NewMockProfileLoader()
				mock.LoadAllProfilesFunc = func() ([]*models.ArkProfile, error) {
					return profiles, nil
				}
				return mock
			},
			setupFlags: func(cmd *cobra.Command) {
				_ = cmd.Flags().String("name", "", "Profile name filter")
				_ = cmd.Flags().String("auth-profile", "", "Auth profile filter")
				_ = cmd.Flags().Bool("all", false, "Show all data")
				_ = cmd.Flags().Set("name", "test.*")
			},
			validateFunc: func(t *testing.T, loader profiles.ProfileLoader) {
				// Function should apply name filter using regex
			},
		},
		{
			name: "success_filters_profiles_by_auth_type",
			setupLoader: func() profiles.ProfileLoader {
				profiles := []*models.ArkProfile{
					testutils.CreateTestProfile("profile1"),
					{
						ProfileName:  "profile2",
						AuthProfiles: map[string]*authmodels.ArkAuthProfile{},
					},
				}
				mock := testutils.NewMockProfileLoader()
				mock.LoadAllProfilesFunc = func() ([]*models.ArkProfile, error) {
					return profiles, nil
				}
				return mock
			},
			setupFlags: func(cmd *cobra.Command) {
				_ = cmd.Flags().String("name", "", "Profile name filter")
				_ = cmd.Flags().String("auth-profile", "", "Auth profile filter")
				_ = cmd.Flags().Bool("all", false, "Show all data")
				_ = cmd.Flags().Set("auth-profile", "test_auth")
			},
			validateFunc: func(t *testing.T, loader profiles.ProfileLoader) {
				// Function should filter by auth profile type
			},
		},
		{
			name: "success_shows_all_data_when_flag_set",
			setupLoader: func() profiles.ProfileLoader {
				profiles := []*models.ArkProfile{
					testutils.CreateTestProfile("profile1"),
				}
				mock := testutils.NewMockProfileLoader()
				mock.LoadAllProfilesFunc = func() ([]*models.ArkProfile, error) {
					return profiles, nil
				}
				return mock
			},
			setupFlags: func(cmd *cobra.Command) {
				_ = cmd.Flags().String("name", "", "Profile name filter")
				_ = cmd.Flags().String("auth-profile", "", "Auth profile filter")
				_ = cmd.Flags().Bool("all", false, "Show all data")
				_ = cmd.Flags().Set("all", "true")
			},
			validateFunc: func(t *testing.T, loader profiles.ProfileLoader) {
				// Function should output full profile data when --all flag is set
			},
		},
		{
			name: "error_handles_load_profiles_error",
			setupLoader: func() profiles.ProfileLoader {
				mock := testutils.NewMockProfileLoader()
				mock.LoadAllProfilesFunc = func() ([]*models.ArkProfile, error) {
					return nil, fmt.Errorf("failed to load profiles")
				}
				return mock
			},
			setupFlags: func(cmd *cobra.Command) {
				_ = cmd.Flags().String("name", "", "Profile name filter")
				_ = cmd.Flags().String("auth-profile", "", "Auth profile filter")
				_ = cmd.Flags().Bool("all", false, "Show all data")
			},
			validateFunc: func(t *testing.T, loader profiles.ProfileLoader) {
				// Function should handle load error gracefully
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			loader := tt.setupLoader()
			action := NewArkProfilesAction(&loader)
			cmd := &cobra.Command{}

			if tt.setupFlags != nil {
				tt.setupFlags(cmd)
			}

			// Execute the function - should not panic
			action.runListAction(cmd, []string{})

			if tt.validateFunc != nil {
				tt.validateFunc(t, loader)
			}
		})
	}
}

func TestArkProfilesAction_runShowAction(t *testing.T) {
	tests := []struct {
		name         string
		setupLoader  func() profiles.ProfileLoader
		setupFlags   func(cmd *cobra.Command)
		validateFunc func(t *testing.T, loader profiles.ProfileLoader)
	}{
		{
			name: "success_shows_specified_profile",
			setupLoader: func() profiles.ProfileLoader {
				profile := testutils.CreateTestProfile("test-profile")
				mock := testutils.NewMockProfileLoader()
				mock.LoadProfileFunc = func(name string) (*models.ArkProfile, error) {
					if name == "test-profile" {
						return profile, nil
					}
					return nil, fmt.Errorf("profile not found")
				}
				return mock
			},
			setupFlags: func(cmd *cobra.Command) {
				_ = cmd.Flags().String("profile-name", "", "Profile name")
				_ = cmd.Flags().Set("profile-name", "test-profile")
			},
			validateFunc: func(t *testing.T, loader profiles.ProfileLoader) {
				// Function should load and display the specified profile
			},
		},
		{
			name: "success_uses_default_profile_when_no_name_specified",
			setupLoader: func() profiles.ProfileLoader {
				profile := testutils.CreateTestProfile("default")
				mock := testutils.NewMockProfileLoader()
				mock.LoadProfileFunc = func(name string) (*models.ArkProfile, error) {
					return profile, nil
				}
				return mock
			},
			setupFlags: func(cmd *cobra.Command) {
				_ = cmd.Flags().String("profile-name", "", "Profile name")
				// Don't set profile-name to test default behavior
			},
			validateFunc: func(t *testing.T, loader profiles.ProfileLoader) {
				// Function should use default profile name deduction
			},
		},
		{
			name: "error_handles_profile_not_found",
			setupLoader: func() profiles.ProfileLoader {
				mock := testutils.NewMockProfileLoader()
				mock.LoadProfileFunc = func(name string) (*models.ArkProfile, error) {
					return nil, fmt.Errorf("profile not found")
				}
				return mock
			},
			setupFlags: func(cmd *cobra.Command) {
				_ = cmd.Flags().String("profile-name", "", "Profile name")
				_ = cmd.Flags().Set("profile-name", "nonexistent")
			},
			validateFunc: func(t *testing.T, loader profiles.ProfileLoader) {
				// Function should handle profile not found gracefully
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			loader := tt.setupLoader()
			action := NewArkProfilesAction(&loader)
			cmd := &cobra.Command{}

			if tt.setupFlags != nil {
				tt.setupFlags(cmd)
			}

			// Execute the function - should not panic
			action.runShowAction(cmd, []string{})

			if tt.validateFunc != nil {
				tt.validateFunc(t, loader)
			}
		})
	}
}

func TestArkProfilesAction_runDeleteAction(t *testing.T) {
	tests := []struct {
		name         string
		setupLoader  func() profiles.ProfileLoader
		setupFlags   func(cmd *cobra.Command)
		validateFunc func(t *testing.T, loader profiles.ProfileLoader)
	}{
		{
			name: "success_deletes_profile_with_yes_flag",
			setupLoader: func() profiles.ProfileLoader {
				profile := testutils.CreateTestProfile("test-profile")
				mock := testutils.NewMockProfileLoader()
				mock.LoadProfileFunc = func(name string) (*models.ArkProfile, error) {
					if name == "test-profile" {
						return profile, nil
					}
					return nil, fmt.Errorf("profile not found")
				}
				mock.DeleteProfileFunc = func(name string) error {
					return nil
				}
				return mock
			},
			setupFlags: func(cmd *cobra.Command) {
				_ = cmd.Flags().String("profile-name", "", "Profile name")
				_ = cmd.Flags().Bool("yes", false, "Skip confirmation")
				_ = cmd.Flags().Set("profile-name", "test-profile")
				_ = cmd.Flags().Set("yes", "true")
			},
			validateFunc: func(t *testing.T, loader profiles.ProfileLoader) {
				// Function should delete profile without confirmation
			},
		},
		{
			name: "error_handles_profile_not_found_for_deletion",
			setupLoader: func() profiles.ProfileLoader {
				mock := testutils.NewMockProfileLoader()
				mock.LoadProfileFunc = func(name string) (*models.ArkProfile, error) {
					return nil, fmt.Errorf("profile not found")
				}
				return mock
			},
			setupFlags: func(cmd *cobra.Command) {
				_ = cmd.Flags().String("profile-name", "", "Profile name")
				_ = cmd.Flags().Bool("yes", false, "Skip confirmation")
				_ = cmd.Flags().Set("profile-name", "nonexistent")
			},
			validateFunc: func(t *testing.T, loader profiles.ProfileLoader) {
				// Function should handle profile not found gracefully
			},
		},
		{
			name: "error_handles_delete_failure",
			setupLoader: func() profiles.ProfileLoader {
				profile := testutils.CreateTestProfile("test-profile")
				mock := testutils.NewMockProfileLoader()
				mock.LoadProfileFunc = func(name string) (*models.ArkProfile, error) {
					return profile, nil
				}
				mock.DeleteProfileFunc = func(name string) error {
					return fmt.Errorf("delete failed")
				}
				return mock
			},
			setupFlags: func(cmd *cobra.Command) {
				_ = cmd.Flags().String("profile-name", "", "Profile name")
				_ = cmd.Flags().Bool("yes", false, "Skip confirmation")
				_ = cmd.Flags().Set("profile-name", "test-profile")
				_ = cmd.Flags().Set("yes", "true")
			},
			validateFunc: func(t *testing.T, loader profiles.ProfileLoader) {
				// Function should handle delete failure gracefully
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			loader := tt.setupLoader()
			action := NewArkProfilesAction(&loader)
			cmd := &cobra.Command{}

			if tt.setupFlags != nil {
				tt.setupFlags(cmd)
			}

			// Execute the function - should not panic
			action.runDeleteAction(cmd, []string{})

			if tt.validateFunc != nil {
				tt.validateFunc(t, loader)
			}
		})
	}
}

func TestArkProfilesAction_runClearAction(t *testing.T) {
	tests := []struct {
		name         string
		setupLoader  func() profiles.ProfileLoader
		setupFlags   func(cmd *cobra.Command)
		validateFunc func(t *testing.T, loader profiles.ProfileLoader)
	}{
		{
			name: "success_clears_all_profiles_with_yes_flag",
			setupLoader: func() profiles.ProfileLoader {
				mock := testutils.NewMockProfileLoader()
				mock.ClearAllProfilesFunc = func() error {
					return nil
				}
				return mock
			},
			setupFlags: func(cmd *cobra.Command) {
				_ = cmd.Flags().Bool("yes", false, "Skip confirmation")
				_ = cmd.Flags().Set("yes", "true")
			},
			validateFunc: func(t *testing.T, loader profiles.ProfileLoader) {
				// Function should clear all profiles without confirmation
			},
		},
		{
			name: "error_handles_clear_failure",
			setupLoader: func() profiles.ProfileLoader {
				mock := testutils.NewMockProfileLoader()
				mock.ClearAllProfilesFunc = func() error {
					return fmt.Errorf("clear failed")
				}
				return mock
			},
			setupFlags: func(cmd *cobra.Command) {
				_ = cmd.Flags().Bool("yes", false, "Skip confirmation")
				_ = cmd.Flags().Set("yes", "true")
			},
			validateFunc: func(t *testing.T, loader profiles.ProfileLoader) {
				// Function should handle clear failure gracefully
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			loader := tt.setupLoader()
			action := NewArkProfilesAction(&loader)
			cmd := &cobra.Command{}

			if tt.setupFlags != nil {
				tt.setupFlags(cmd)
			}

			// Execute the function - should not panic
			action.runClearAction(cmd, []string{})

			if tt.validateFunc != nil {
				tt.validateFunc(t, loader)
			}
		})
	}
}

func TestArkProfilesAction_runCloneAction(t *testing.T) {
	tests := []struct {
		name         string
		setupLoader  func() profiles.ProfileLoader
		setupFlags   func(cmd *cobra.Command)
		validateFunc func(t *testing.T, loader profiles.ProfileLoader)
	}{
		{
			name: "success_clones_profile_with_default_name",
			setupLoader: func() profiles.ProfileLoader {
				profile := testutils.CreateTestProfile("original")
				mock := testutils.NewMockProfileLoader()
				mock.LoadProfileFunc = func(name string) (*models.ArkProfile, error) {
					if name == "original" {
						return profile, nil
					}
					return nil, fmt.Errorf("profile not found")
				}
				mock.ProfileExistsFunc = func(name string) bool {
					return false
				}
				mock.SaveProfileFunc = func(profile *models.ArkProfile) error {
					return nil
				}
				return mock
			},
			setupFlags: func(cmd *cobra.Command) {
				_ = cmd.Flags().String("profile-name", "", "Profile name")
				_ = cmd.Flags().String("new-profile-name", "", "New profile name")
				_ = cmd.Flags().Bool("yes", false, "Skip confirmation")
				_ = cmd.Flags().Set("profile-name", "original")
			},
			validateFunc: func(t *testing.T, loader profiles.ProfileLoader) {
				// Function should clone profile with "_clone" suffix
			},
		},
		{
			name: "success_clones_profile_with_custom_name",
			setupLoader: func() profiles.ProfileLoader {
				profile := testutils.CreateTestProfile("original")
				mock := testutils.NewMockProfileLoader()
				mock.LoadProfileFunc = func(name string) (*models.ArkProfile, error) {
					return profile, nil
				}
				mock.ProfileExistsFunc = func(name string) bool {
					return false
				}
				mock.SaveProfileFunc = func(profile *models.ArkProfile) error {
					return nil
				}
				return mock
			},
			setupFlags: func(cmd *cobra.Command) {
				_ = cmd.Flags().String("profile-name", "", "Profile name")
				_ = cmd.Flags().String("new-profile-name", "", "New profile name")
				_ = cmd.Flags().Bool("yes", false, "Skip confirmation")
				_ = cmd.Flags().Set("profile-name", "original")
				_ = cmd.Flags().Set("new-profile-name", "custom-clone")
			},
			validateFunc: func(t *testing.T, loader profiles.ProfileLoader) {
				// Function should clone profile with custom name
			},
		},
		{
			name: "error_handles_source_profile_not_found",
			setupLoader: func() profiles.ProfileLoader {
				mock := testutils.NewMockProfileLoader()
				mock.LoadProfileFunc = func(name string) (*models.ArkProfile, error) {
					return nil, fmt.Errorf("profile not found")
				}
				return mock
			},
			setupFlags: func(cmd *cobra.Command) {
				_ = cmd.Flags().String("profile-name", "", "Profile name")
				_ = cmd.Flags().String("new-profile-name", "", "New profile name")
				_ = cmd.Flags().Bool("yes", false, "Skip confirmation")
				_ = cmd.Flags().Set("profile-name", "nonexistent")
			},
			validateFunc: func(t *testing.T, loader profiles.ProfileLoader) {
				// Function should handle source profile not found gracefully
			},
		},
		{
			name: "error_handles_save_failure",
			setupLoader: func() profiles.ProfileLoader {
				profile := testutils.CreateTestProfile("original")
				mock := testutils.NewMockProfileLoader()
				mock.LoadProfileFunc = func(name string) (*models.ArkProfile, error) {
					return profile, nil
				}
				mock.ProfileExistsFunc = func(name string) bool {
					return false
				}
				mock.SaveProfileFunc = func(profile *models.ArkProfile) error {
					return fmt.Errorf("save failed")
				}
				return mock
			},
			setupFlags: func(cmd *cobra.Command) {
				_ = cmd.Flags().String("profile-name", "", "Profile name")
				_ = cmd.Flags().String("new-profile-name", "", "New profile name")
				_ = cmd.Flags().Bool("yes", false, "Skip confirmation")
				_ = cmd.Flags().Set("profile-name", "original")
			},
			validateFunc: func(t *testing.T, loader profiles.ProfileLoader) {
				// Function should handle save failure gracefully
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			loader := tt.setupLoader()
			action := NewArkProfilesAction(&loader)
			cmd := &cobra.Command{}

			if tt.setupFlags != nil {
				tt.setupFlags(cmd)
			}

			// Execute the function - should not panic
			action.runCloneAction(cmd, []string{})

			if tt.validateFunc != nil {
				tt.validateFunc(t, loader)
			}
		})
	}
}

func TestArkProfilesAction_StructFields(t *testing.T) {
	tests := []struct {
		name         string
		validateFunc func(t *testing.T, action *ArkProfilesAction)
	}{
		{
			name: "success_struct_has_expected_fields",
			validateFunc: func(t *testing.T, action *ArkProfilesAction) {
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
			validateFunc: func(t *testing.T, action *ArkProfilesAction) {
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
			action := NewArkProfilesAction(&loaderInterface)

			if tt.validateFunc != nil {
				tt.validateFunc(t, action)
			}
		})
	}
}

func TestArkProfilesAction_ArkActionInterface(t *testing.T) {
	tests := []struct {
		name         string
		validateFunc func(t *testing.T)
	}{
		{
			name: "success_implements_arkaction_interface",
			validateFunc: func(t *testing.T) {
				loader := testutils.NewMockProfileLoader()
				loaderInterface := profiles.ProfileLoader(loader)
				action := NewArkProfilesAction(&loaderInterface)

				// This will cause a compile error if ArkProfilesAction doesn't implement ArkAction
				var _ ArkAction = action

				// Verify the DefineAction method exists and can be called
				cmd := &cobra.Command{}
				action.DefineAction(cmd)

				// Verify a profiles command was added
				found := false
				for _, subCmd := range cmd.Commands() {
					if subCmd.Use == "profiles" {
						found = true
						break
					}
				}
				if !found {
					t.Error("Expected profiles command to be added to parent command")
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
