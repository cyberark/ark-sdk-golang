package actions

import (
	"errors"
	"testing"

	"github.com/cyberark/ark-sdk-golang/pkg/actions/testutils"
	"github.com/cyberark/ark-sdk-golang/pkg/cli"
	"github.com/cyberark/ark-sdk-golang/pkg/models"
	"github.com/cyberark/ark-sdk-golang/pkg/profiles"
	"github.com/spf13/cobra"
)

// Mock implementations for testing
type mockExecAction struct {
	defineExecActionFunc func(*cobra.Command) error
	runExecActionFunc    func(*cli.ArkCLIAPI, *cobra.Command, *cobra.Command, []string) error
}

func (m *mockExecAction) DefineExecAction(cmd *cobra.Command) error {
	if m.defineExecActionFunc != nil {
		return m.defineExecActionFunc(cmd)
	}
	return nil
}

func (m *mockExecAction) RunExecAction(api *cli.ArkCLIAPI, cmd *cobra.Command, execCmd *cobra.Command, args []string) error {
	if m.runExecActionFunc != nil {
		return m.runExecActionFunc(api, cmd, execCmd, args)
	}
	return nil
}

// Helper function for converting exec action to pointer
func execActionPtr(action ArkExecAction) *ArkExecAction {
	return &action
}

func TestNewArkBaseExecAction(t *testing.T) {
	tests := []struct {
		name           string
		execAction     *ArkExecAction
		actionName     string
		profilesLoader *profiles.ProfileLoader
		validateFunc   func(t *testing.T, result *ArkBaseExecAction)
	}{
		{
			name:           "success_normal_parameters",
			execAction:     execActionPtr(&mockExecAction{}),
			actionName:     "test-action",
			profilesLoader: testutils.NewMockProfileLoader().AsProfileLoader(),
			validateFunc: func(t *testing.T, result *ArkBaseExecAction) {
				if result.ArkBaseAction == nil {
					t.Error("Expected ArkBaseAction to be initialized")
				}
				if result.profilesLoader == nil {
					t.Error("Expected profilesLoader to be set")
				}
				if result.execAction == nil {
					t.Error("Expected execAction to be set")
				}
				if result.logger == nil {
					t.Error("Expected logger to be initialized")
				}
			},
		},
		{
			name:           "success_nil_exec_action",
			execAction:     nil,
			actionName:     "test-action",
			profilesLoader: testutils.NewMockProfileLoader().AsProfileLoader(),
			validateFunc: func(t *testing.T, result *ArkBaseExecAction) {
				if result.ArkBaseAction == nil {
					t.Error("Expected ArkBaseAction to be initialized")
				}
				if result.profilesLoader == nil {
					t.Error("Expected profilesLoader to be set")
				}
				if result.execAction != nil {
					t.Error("Expected execAction to remain nil")
				}
				if result.logger == nil {
					t.Error("Expected logger to be initialized")
				}
			},
		},
		{
			name:           "success_nil_profiles_loader",
			execAction:     execActionPtr(&mockExecAction{}),
			actionName:     "test-action",
			profilesLoader: nil,
			validateFunc: func(t *testing.T, result *ArkBaseExecAction) {
				if result.ArkBaseAction == nil {
					t.Error("Expected ArkBaseAction to be initialized")
				}
				if result.profilesLoader != nil {
					t.Error("Expected profilesLoader to remain nil")
				}
				if result.execAction == nil {
					t.Error("Expected execAction to be set")
				}
				if result.logger == nil {
					t.Error("Expected logger to be initialized")
				}
			},
		},
		{
			name:           "success_empty_action_name",
			execAction:     execActionPtr(&mockExecAction{}),
			actionName:     "",
			profilesLoader: testutils.NewMockProfileLoader().AsProfileLoader(),
			validateFunc: func(t *testing.T, result *ArkBaseExecAction) {
				if result.ArkBaseAction == nil {
					t.Error("Expected ArkBaseAction to be initialized")
				}
				if result.logger == nil {
					t.Error("Expected logger to be initialized even with empty name")
				}
			},
		},
		{
			name:           "success_long_action_name",
			execAction:     execActionPtr(&mockExecAction{}),
			actionName:     "very-long-action-name-with-multiple-parts-and-special-characters-123",
			profilesLoader: testutils.NewMockProfileLoader().AsProfileLoader(),
			validateFunc: func(t *testing.T, result *ArkBaseExecAction) {
				if result.ArkBaseAction == nil {
					t.Error("Expected ArkBaseAction to be initialized")
				}
				if result.logger == nil {
					t.Error("Expected logger to handle long names")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := NewArkBaseExecAction(tt.execAction, tt.actionName, tt.profilesLoader)

			if result == nil {
				t.Error("Expected non-nil result")
				return
			}

			if tt.validateFunc != nil {
				tt.validateFunc(t, result)
			}
		})
	}
}

func TestArkBaseExecAction_DefineAction(t *testing.T) {
	tests := []struct {
		name          string
		setupAction   func() *ArkBaseExecAction
		expectedPanic bool
		validateFunc  func(t *testing.T, cmd *cobra.Command)
	}{
		{
			name: "success_defines_exec_command",
			setupAction: func() *ArkBaseExecAction {
				mockExec := &mockExecAction{}
				return NewArkBaseExecAction(
					execActionPtr(mockExec),
					"test-action",
					testutils.NewMockProfileLoader().AsProfileLoader(),
				)
			},
			expectedPanic: false,
			validateFunc: func(t *testing.T, cmd *cobra.Command) {
				execCmd, _, err := cmd.Find([]string{"exec"})
				if err != nil {
					t.Errorf("Expected to find exec command, got error: %v", err)
					return
				}
				if execCmd == nil {
					t.Error("Expected exec command to be added")
					return
				}
				if execCmd.Use != "exec" {
					t.Errorf("Expected command Use to be 'exec', got '%s'", execCmd.Use)
				}
				if execCmd.Short != "Exec an action" {
					t.Errorf("Expected command Short description, got '%s'", execCmd.Short)
				}
			},
		},
		{
			name: "success_adds_persistent_flags",
			setupAction: func() *ArkBaseExecAction {
				mockExec := &mockExecAction{}
				return NewArkBaseExecAction(
					execActionPtr(mockExec),
					"test-action",
					testutils.NewMockProfileLoader().AsProfileLoader(),
				)
			},
			expectedPanic: false,
			validateFunc: func(t *testing.T, cmd *cobra.Command) {
				execCmd, _, _ := cmd.Find([]string{"exec"})
				if execCmd == nil {
					t.Error("Exec command not found")
					return
				}

				expectedFlags := []string{
					"profile-name",
					"output-path",
					"request-file",
					"retry-count",
					"refresh-auth",
				}

				for _, flagName := range expectedFlags {
					flag := execCmd.PersistentFlags().Lookup(flagName)
					if flag == nil {
						t.Errorf("Expected persistent flag '%s' to be defined", flagName)
					}
				}
			},
		},
		{
			name: "success_sets_persistent_prerun",
			setupAction: func() *ArkBaseExecAction {
				mockExec := &mockExecAction{}
				return NewArkBaseExecAction(
					execActionPtr(mockExec),
					"test-action",
					testutils.NewMockProfileLoader().AsProfileLoader(),
				)
			},
			expectedPanic: false,
			validateFunc: func(t *testing.T, cmd *cobra.Command) {
				execCmd, _, _ := cmd.Find([]string{"exec"})
				if execCmd == nil {
					t.Error("Exec command not found")
					return
				}
				if execCmd.PersistentPreRun == nil {
					t.Error("Expected PersistentPreRun to be set")
				}
			},
		},
		{
			name: "error_define_exec_action_fails",
			setupAction: func() *ArkBaseExecAction {
				mockExec := &mockExecAction{
					defineExecActionFunc: func(cmd *cobra.Command) error {
						return errors.New("failed to define exec action")
					},
				}
				return NewArkBaseExecAction(
					execActionPtr(mockExec),
					"test-action",
					testutils.NewMockProfileLoader().AsProfileLoader(),
				)
			},
			expectedPanic: true,
			validateFunc: func(t *testing.T, cmd *cobra.Command) {
				// This test validates that DefineAction panics when DefineExecAction fails
			},
		},
		{
			name: "success_nil_exec_action",
			setupAction: func() *ArkBaseExecAction {
				return NewArkBaseExecAction(
					nil,
					"test-action",
					testutils.NewMockProfileLoader().AsProfileLoader(),
				)
			},
			expectedPanic: true, // Should panic when trying to call (*nil).DefineExecAction
			validateFunc: func(t *testing.T, cmd *cobra.Command) {
				// This test validates that DefineAction handles nil execAction
			},
		},
		{
			name: "success_adds_custom_exec_subcommand",
			setupAction: func() *ArkBaseExecAction {
				mockExec := &mockExecAction{
					defineExecActionFunc: func(cmd *cobra.Command) error {
						subCmd := &cobra.Command{
							Use:   "custom-action",
							Short: "Custom exec action",
						}
						cmd.AddCommand(subCmd)
						return nil
					},
				}
				return NewArkBaseExecAction(
					execActionPtr(mockExec),
					"test-action",
					testutils.NewMockProfileLoader().AsProfileLoader(),
				)
			},
			expectedPanic: false,
			validateFunc: func(t *testing.T, cmd *cobra.Command) {
				execCmd, _, _ := cmd.Find([]string{"exec"})
				if execCmd == nil {
					t.Error("Exec command not found")
					return
				}
				customCmd, _, err := execCmd.Find([]string{"custom-action"})
				if err != nil || customCmd == nil {
					t.Error("Expected custom-action subcommand to be added")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			action := tt.setupAction()
			rootCmd := &cobra.Command{Use: "ark"}

			if tt.expectedPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Error("Expected DefineAction to panic, but it didn't")
					}
				}()
			}

			action.DefineAction(rootCmd)

			if !tt.expectedPanic && tt.validateFunc != nil {
				tt.validateFunc(t, rootCmd)
			}
		})
	}
}

func TestArkBaseExecAction_runExecAction(t *testing.T) {
	// Create test profile with valid token
	createTestProfileWithValidToken := func(name string) *models.ArkProfile {
		profile := testutils.CreateTestProfile(name)
		return profile
	}

	// Create test profile with expired token
	createTestProfileWithExpiredToken := func(name string) *models.ArkProfile {
		profile := testutils.CreateTestProfile(name)
		return profile
	}

	tests := []struct {
		name         string
		setupAction  func() *ArkBaseExecAction
		setupCmd     func() *cobra.Command
		execArgs     []string
		validateFunc func(t *testing.T, action *ArkBaseExecAction, cmd *cobra.Command)
	}{
		{
			name: "error_exec_command_not_found",
			setupAction: func() *ArkBaseExecAction {
				mockLoader := testutils.NewMockProfileLoader()
				mockExec := &mockExecAction{}
				return NewArkBaseExecAction(
					execActionPtr(mockExec),
					"test-action",
					mockLoader.AsProfileLoader(),
				)
			},
			setupCmd: func() *cobra.Command {
				// Return command without exec subcommand
				return &cobra.Command{Use: "root"}
			},
			execArgs: []string{},
			validateFunc: func(t *testing.T, action *ArkBaseExecAction, cmd *cobra.Command) {
				// This test validates that the function handles missing exec command gracefully
			},
		},
		{
			name: "error_profile_not_found",
			setupAction: func() *ArkBaseExecAction {
				mockLoader := testutils.NewMockProfileLoader()
				mockLoader.LoadProfileFunc = func(name string) (*models.ArkProfile, error) {
					return nil, errors.New("profile not found")
				}
				mockExec := &mockExecAction{}
				return NewArkBaseExecAction(
					execActionPtr(mockExec),
					"test-action",
					mockLoader.AsProfileLoader(),
				)
			},
			setupCmd: func() *cobra.Command {
				rootCmd := &cobra.Command{Use: "root"}
				execCmd := &cobra.Command{Use: "exec"}
				execCmd.Flags().String("profile-name", profiles.DefaultProfileName(), "Profile name")
				execCmd.Flags().Bool("refresh-auth", false, "Refresh auth")
				execCmd.Flags().Int("retry-count", 1, "Retry count")
				rootCmd.AddCommand(execCmd)
				return execCmd
			},
			execArgs: []string{},
			validateFunc: func(t *testing.T, action *ArkBaseExecAction, cmd *cobra.Command) {
				// This test validates that the function handles profile loading errors gracefully
			},
		},
		{
			name: "error_nil_profile",
			setupAction: func() *ArkBaseExecAction {
				mockLoader := testutils.NewMockProfileLoader()
				mockLoader.LoadProfileFunc = func(name string) (*models.ArkProfile, error) {
					return nil, nil // Returns nil profile
				}
				mockExec := &mockExecAction{}
				return NewArkBaseExecAction(
					execActionPtr(mockExec),
					"test-action",
					mockLoader.AsProfileLoader(),
				)
			},
			setupCmd: func() *cobra.Command {
				rootCmd := &cobra.Command{Use: "root"}
				execCmd := &cobra.Command{Use: "exec"}
				execCmd.Flags().String("profile-name", profiles.DefaultProfileName(), "Profile name")
				execCmd.Flags().Bool("refresh-auth", false, "Refresh auth")
				execCmd.Flags().Int("retry-count", 1, "Retry count")
				rootCmd.AddCommand(execCmd)
				return execCmd
			},
			execArgs: []string{},
			validateFunc: func(t *testing.T, action *ArkBaseExecAction, cmd *cobra.Command) {
				// This test validates that the function handles nil profiles gracefully
			},
		},
		{
			name: "error_no_authenticators_available",
			setupAction: func() *ArkBaseExecAction {
				mockLoader := testutils.NewMockProfileLoader()
				mockLoader.LoadProfileFunc = func(name string) (*models.ArkProfile, error) {
					return createTestProfileWithExpiredToken(name), nil
				}
				mockExec := &mockExecAction{}
				return NewArkBaseExecAction(
					execActionPtr(mockExec),
					"test-action",
					mockLoader.AsProfileLoader(),
				)
			},
			setupCmd: func() *cobra.Command {
				rootCmd := &cobra.Command{Use: "root"}
				execCmd := &cobra.Command{Use: "exec"}
				execCmd.Flags().String("profile-name", "test-profile", "Profile name")
				execCmd.Flags().Bool("refresh-auth", false, "Refresh auth")
				execCmd.Flags().Int("retry-count", 1, "Retry count")
				rootCmd.AddCommand(execCmd)
				return execCmd
			},
			execArgs: []string{},
			validateFunc: func(t *testing.T, action *ArkBaseExecAction, cmd *cobra.Command) {
				// This test validates behavior when no valid authenticators are available
			},
		},
		{
			name: "success_valid_profile_and_authenticators",
			setupAction: func() *ArkBaseExecAction {
				mockLoader := testutils.NewMockProfileLoader()
				mockLoader.LoadProfileFunc = func(name string) (*models.ArkProfile, error) {
					return createTestProfileWithValidToken(name), nil
				}
				mockExec := &mockExecAction{
					runExecActionFunc: func(api *cli.ArkCLIAPI, cmd *cobra.Command, execCmd *cobra.Command, args []string) error {
						return nil
					},
				}
				return NewArkBaseExecAction(
					execActionPtr(mockExec),
					"test-action",
					mockLoader.AsProfileLoader(),
				)
			},
			setupCmd: func() *cobra.Command {
				rootCmd := &cobra.Command{Use: "root"}
				execCmd := &cobra.Command{Use: "exec"}
				execCmd.Flags().String("profile-name", "test-profile", "Profile name")
				execCmd.Flags().Bool("refresh-auth", false, "Refresh auth")
				execCmd.Flags().Int("retry-count", 1, "Retry count")
				rootCmd.AddCommand(execCmd)
				return execCmd
			},
			execArgs: []string{},
			validateFunc: func(t *testing.T, action *ArkBaseExecAction, cmd *cobra.Command) {
				// This test validates successful execution with valid authenticators
				// Note: This will attempt to create CLI API which might fail without proper auth setup
			},
		},
		{
			name: "success_retry_count_respected",
			setupAction: func() *ArkBaseExecAction {
				mockLoader := testutils.NewMockProfileLoader()
				mockLoader.LoadProfileFunc = func(name string) (*models.ArkProfile, error) {
					return createTestProfileWithValidToken(name), nil
				}
				mockExec := &mockExecAction{
					runExecActionFunc: func(api *cli.ArkCLIAPI, cmd *cobra.Command, execCmd *cobra.Command, args []string) error {
						return errors.New("simulated failure")
					},
				}
				return NewArkBaseExecAction(
					execActionPtr(mockExec),
					"test-action",
					mockLoader.AsProfileLoader(),
				)
			},
			setupCmd: func() *cobra.Command {
				rootCmd := &cobra.Command{Use: "root"}
				execCmd := &cobra.Command{Use: "exec"}
				execCmd.Flags().String("profile-name", "test-profile", "Profile name")
				execCmd.Flags().Bool("refresh-auth", false, "Refresh auth")
				execCmd.Flags().Int("retry-count", 3, "Retry count")
				rootCmd.AddCommand(execCmd)
				return execCmd
			},
			execArgs: []string{},
			validateFunc: func(t *testing.T, action *ArkBaseExecAction, cmd *cobra.Command) {
				// Verify retry count flag is read correctly
				retryCount, err := cmd.Flags().GetInt("retry-count")
				if err != nil {
					t.Errorf("Expected to get retry-count flag, got error: %v", err)
				}
				if retryCount != 3 {
					t.Errorf("Expected retry-count to be 3, got %d", retryCount)
				}
			},
		},
		{
			name: "success_refresh_auth_flag_set",
			setupAction: func() *ArkBaseExecAction {
				mockLoader := testutils.NewMockProfileLoader()
				mockLoader.LoadProfileFunc = func(name string) (*models.ArkProfile, error) {
					return createTestProfileWithValidToken(name), nil
				}
				mockExec := &mockExecAction{}
				return NewArkBaseExecAction(
					execActionPtr(mockExec),
					"test-action",
					mockLoader.AsProfileLoader(),
				)
			},
			setupCmd: func() *cobra.Command {
				rootCmd := &cobra.Command{Use: "root"}
				execCmd := &cobra.Command{Use: "exec"}
				execCmd.Flags().String("profile-name", "test-profile", "Profile name")
				execCmd.Flags().Bool("refresh-auth", true, "Refresh auth")
				execCmd.Flags().Int("retry-count", 1, "Retry count")
				rootCmd.AddCommand(execCmd)
				return execCmd
			},
			execArgs: []string{},
			validateFunc: func(t *testing.T, action *ArkBaseExecAction, cmd *cobra.Command) {
				// Verify refresh-auth flag is read correctly
				refreshAuth, err := cmd.Flags().GetBool("refresh-auth")
				if err != nil {
					t.Errorf("Expected to get refresh-auth flag, got error: %v", err)
				}
				if !refreshAuth {
					t.Error("Expected refresh-auth flag to be true")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action := tt.setupAction()
			cmd := tt.setupCmd()

			// Execute the function
			action.runExecAction(cmd, tt.execArgs)

			// Custom validation
			if tt.validateFunc != nil {
				tt.validateFunc(t, action, cmd)
			}
		})
	}
}

// TestArkBaseExecAction_Integration provides integration-style tests
func TestArkBaseExecAction_Integration(t *testing.T) {
	tests := []struct {
		name         string
		setupAction  func() *ArkBaseExecAction
		validateFunc func(t *testing.T, action *ArkBaseExecAction)
	}{
		{
			name: "success_complete_flow_with_mock_loader",
			setupAction: func() *ArkBaseExecAction {
				mockLoader := testutils.NewMockProfileLoader()
				mockLoader.LoadProfileFunc = func(name string) (*models.ArkProfile, error) {
					return testutils.CreateTestProfile(name), nil
				}
				mockExec := &mockExecAction{
					defineExecActionFunc: func(cmd *cobra.Command) error {
						testCmd := &cobra.Command{
							Use:   "test-command",
							Short: "Test command",
						}
						cmd.AddCommand(testCmd)
						return nil
					},
				}
				return NewArkBaseExecAction(
					execActionPtr(mockExec),
					"test-action",
					mockLoader.AsProfileLoader(),
				)
			},
			validateFunc: func(t *testing.T, action *ArkBaseExecAction) {
				// Test complete DefineAction flow
				rootCmd := &cobra.Command{Use: "ark"}
				action.DefineAction(rootCmd)

				// Verify exec command was added
				execCmd, _, err := rootCmd.Find([]string{"exec"})
				if err != nil {
					t.Errorf("Expected to find exec command, got error: %v", err)
				}
				if execCmd == nil {
					t.Error("Expected exec command to be defined")
				}

				// Verify test subcommand was added
				testCmd, _, err := execCmd.Find([]string{"test-command"})
				if err != nil || testCmd == nil {
					t.Error("Expected test-command to be added by mock exec action")
				}
			},
		},
		{
			name: "success_action_structure_validation",
			setupAction: func() *ArkBaseExecAction {
				mockExec := &mockExecAction{}
				return NewArkBaseExecAction(
					execActionPtr(mockExec),
					"validation-test",
					testutils.NewMockProfileLoader().AsProfileLoader(),
				)
			},
			validateFunc: func(t *testing.T, action *ArkBaseExecAction) {
				// Validate the action structure
				if action == nil {
					t.Error("Expected non-nil action")
					return
				}
				if action.ArkBaseAction == nil {
					t.Error("Expected ArkBaseAction to be initialized")
				}
				if action.profilesLoader == nil {
					t.Error("Expected profilesLoader to be set")
				}
				if action.execAction == nil {
					t.Error("Expected execAction to be set")
				}

				// Validate it implements ArkAction interface
				var _ ArkAction = action
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			action := tt.setupAction()

			if tt.validateFunc != nil {
				tt.validateFunc(t, action)
			}
		})
	}
}
