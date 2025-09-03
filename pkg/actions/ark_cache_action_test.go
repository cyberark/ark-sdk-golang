package actions

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	commoninternal "github.com/cyberark/ark-sdk-golang/internal/common"
	"github.com/spf13/cobra"
)

func TestNewArkCacheAction(t *testing.T) {
	tests := []struct {
		name         string
		validateFunc func(t *testing.T, action *ArkCacheAction)
	}{
		{
			name: "success_creates_cache_action_with_embedded_base_action",
			validateFunc: func(t *testing.T, action *ArkCacheAction) {
				if action == nil {
					t.Error("Expected action to be created, got nil")
					return
				}
				if action.ArkBaseAction == nil {
					t.Error("Expected embedded ArkBaseAction to be initialized")
				}
			},
		},
		{
			name: "success_embedded_base_action_has_logger",
			validateFunc: func(t *testing.T, action *ArkCacheAction) {
				if action.ArkBaseAction == nil {
					t.Error("Expected embedded ArkBaseAction to be initialized")
					return
				}
				// Access logger through reflection since it's unexported
				actionValue := reflect.ValueOf(action.ArkBaseAction).Elem()
				loggerField := actionValue.FieldByName("logger")
				if !loggerField.IsValid() {
					t.Error("Expected logger field to exist in embedded ArkBaseAction")
					return
				}
				if loggerField.IsNil() {
					t.Error("Expected logger to be initialized")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			action := NewArkCacheAction()

			if tt.validateFunc != nil {
				tt.validateFunc(t, action)
			}
		})
	}
}

func TestArkCacheAction_DefineAction(t *testing.T) {
	tests := []struct {
		name         string
		validateFunc func(t *testing.T, rootCmd *cobra.Command, action *ArkCacheAction)
	}{
		{
			name: "success_adds_cache_command_to_parent",
			validateFunc: func(t *testing.T, rootCmd *cobra.Command, action *ArkCacheAction) {
				cacheCmd, _, err := rootCmd.Find([]string{"cache"})
				if err != nil {
					t.Errorf("Expected to find cache command, got error: %v", err)
					return
				}
				if cacheCmd == nil {
					t.Error("Expected cache command to be added")
					return
				}
				if cacheCmd.Use != "cache" {
					t.Errorf("Expected cache command Use to be 'cache', got '%s'", cacheCmd.Use)
				}
				if cacheCmd.Short != "Manage cache" {
					t.Errorf("Expected cache command Short to be 'Manage cache', got '%s'", cacheCmd.Short)
				}
			},
		},
		{
			name: "success_adds_clear_subcommand",
			validateFunc: func(t *testing.T, rootCmd *cobra.Command, action *ArkCacheAction) {
				clearCmd, _, err := rootCmd.Find([]string{"cache", "clear"})
				if err != nil {
					t.Errorf("Expected to find cache clear command, got error: %v", err)
					return
				}
				if clearCmd == nil {
					t.Error("Expected cache clear command to be added")
					return
				}
				if clearCmd.Use != "clear" {
					t.Errorf("Expected clear command Use to be 'clear', got '%s'", clearCmd.Use)
				}
				if clearCmd.Short != "Clears all profiles cache" {
					t.Errorf("Expected clear command Short to be 'Clears all profiles cache', got '%s'", clearCmd.Short)
				}
			},
		},
		{
			name: "success_cache_command_has_persistent_prerun",
			validateFunc: func(t *testing.T, rootCmd *cobra.Command, action *ArkCacheAction) {
				cacheCmd, _, err := rootCmd.Find([]string{"cache"})
				if err != nil {
					t.Errorf("Expected to find cache command, got error: %v", err)
					return
				}
				if cacheCmd.PersistentPreRun == nil {
					t.Error("Expected cache command to have PersistentPreRun function")
				}
			},
		},
		{
			name: "success_cache_command_has_persistent_flags",
			validateFunc: func(t *testing.T, rootCmd *cobra.Command, action *ArkCacheAction) {
				cacheCmd, _, err := rootCmd.Find([]string{"cache"})
				if err != nil {
					t.Errorf("Expected to find cache command, got error: %v", err)
					return
				}

				// Check for some common flags that should be added by CommonActionsConfiguration
				expectedFlags := []string{"raw", "silent", "verbose"}
				for _, flagName := range expectedFlags {
					flag := cacheCmd.PersistentFlags().Lookup(flagName)
					if flag == nil {
						t.Errorf("Expected persistent flag '%s' to be present", flagName)
					}
				}
			},
		},
		{
			name: "success_clear_command_has_run_function",
			validateFunc: func(t *testing.T, rootCmd *cobra.Command, action *ArkCacheAction) {
				clearCmd, _, err := rootCmd.Find([]string{"cache", "clear"})
				if err != nil {
					t.Errorf("Expected to find cache clear command, got error: %v", err)
					return
				}
				if clearCmd.Run == nil {
					t.Error("Expected clear command to have Run function")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			action := NewArkCacheAction()
			rootCmd := &cobra.Command{Use: "test"}

			action.DefineAction(rootCmd)

			if tt.validateFunc != nil {
				tt.validateFunc(t, rootCmd, action)
			}
		})
	}
}

func TestArkCacheAction_runClearCacheAction(t *testing.T) {
	tests := []struct {
		name         string
		setupEnv     func() (cleanup func())
		createFiles  func(dir string) error
		validateFunc func(t *testing.T, dir string)
	}{
		{
			name: "success_clears_cache_files_from_default_location",
			setupEnv: func() (cleanup func()) {
				// Create a temporary directory for testing
				tempDir, err := os.MkdirTemp("", "ark-cache-test")
				if err != nil {
					panic(err)
				}

				// Set HOME to our temp directory
				originalHome := os.Getenv("HOME")
				_ = os.Setenv("HOME", tempDir)

				return func() {
					_ = os.Setenv("HOME", originalHome)
					_ = os.RemoveAll(tempDir)
				}
			},
			createFiles: func(dir string) error {
				cacheDir := filepath.Join(dir, commoninternal.DefaultBasicKeyringFolder)
				if err := os.MkdirAll(cacheDir, 0755); err != nil {
					return err
				}

				// Create keyring and mac files
				keyringFile := filepath.Join(cacheDir, "keyring")
				macFile := filepath.Join(cacheDir, "mac")

				if err := os.WriteFile(keyringFile, []byte("test keyring data"), 0600); err != nil {
					return err
				}
				if err := os.WriteFile(macFile, []byte("test mac data"), 0600); err != nil {
					return err
				}

				return nil
			},
			validateFunc: func(t *testing.T, dir string) {
				cacheDir := filepath.Join(dir, commoninternal.DefaultBasicKeyringFolder)
				keyringFile := filepath.Join(cacheDir, "keyring")
				macFile := filepath.Join(cacheDir, "mac")

				// Files should be removed
				if _, err := os.Stat(keyringFile); !os.IsNotExist(err) {
					t.Errorf("Expected keyring file to be removed, but it still exists")
				}
				if _, err := os.Stat(macFile); !os.IsNotExist(err) {
					t.Errorf("Expected mac file to be removed, but it still exists")
				}
			},
		},
		{
			name: "success_clears_cache_files_from_env_override_location",
			setupEnv: func() (cleanup func()) {
				tempDir, err := os.MkdirTemp("", "ark-cache-test-env")
				if err != nil {
					panic(err)
				}

				// Set custom cache directory via environment variable
				originalEnv := os.Getenv(commoninternal.ArkBasicKeyringFolderEnvVar)
				_ = os.Setenv(commoninternal.ArkBasicKeyringFolderEnvVar, tempDir)

				return func() {
					if originalEnv != "" {
						_ = os.Setenv(commoninternal.ArkBasicKeyringFolderEnvVar, originalEnv)
					} else {
						_ = os.Unsetenv(commoninternal.ArkBasicKeyringFolderEnvVar)
					}
					_ = os.RemoveAll(tempDir)
				}
			},
			createFiles: func(dir string) error {
				// dir is the custom cache directory
				if err := os.MkdirAll(dir, 0755); err != nil {
					return err
				}

				keyringFile := filepath.Join(dir, "keyring")
				macFile := filepath.Join(dir, "mac")

				if err := os.WriteFile(keyringFile, []byte("test keyring data"), 0600); err != nil {
					return err
				}
				if err := os.WriteFile(macFile, []byte("test mac data"), 0600); err != nil {
					return err
				}

				return nil
			},
			validateFunc: func(t *testing.T, dir string) {
				keyringFile := filepath.Join(dir, "keyring")
				macFile := filepath.Join(dir, "mac")

				// Files should be removed
				if _, err := os.Stat(keyringFile); !os.IsNotExist(err) {
					t.Errorf("Expected keyring file to be removed, but it still exists")
				}
				if _, err := os.Stat(macFile); !os.IsNotExist(err) {
					t.Errorf("Expected mac file to be removed, but it still exists")
				}
			},
		},
		{
			name: "success_handles_missing_files_gracefully",
			setupEnv: func() (cleanup func()) {
				tempDir, err := os.MkdirTemp("", "ark-cache-test-missing")
				if err != nil {
					panic(err)
				}

				originalHome := os.Getenv("HOME")
				_ = os.Setenv("HOME", tempDir)

				return func() {
					_ = os.Setenv("HOME", originalHome)
					_ = os.RemoveAll(tempDir)
				}
			},
			createFiles: func(dir string) error {
				// Don't create any files - test handling of missing files
				return nil
			},
			validateFunc: func(t *testing.T, dir string) {
				// Function should complete without error even if files don't exist
				// This test just ensures no panic occurs
			},
		},
		{
			name: "success_handles_missing_cache_directory_gracefully",
			setupEnv: func() (cleanup func()) {
				tempDir, err := os.MkdirTemp("", "ark-cache-test-no-dir")
				if err != nil {
					panic(err)
				}

				originalHome := os.Getenv("HOME")
				_ = os.Setenv("HOME", tempDir)

				return func() {
					_ = os.Setenv("HOME", originalHome)
					_ = os.RemoveAll(tempDir)
				}
			},
			createFiles: func(dir string) error {
				// Don't create cache directory
				return nil
			},
			validateFunc: func(t *testing.T, dir string) {
				// Function should complete without error even if directory doesn't exist
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			cleanup := tt.setupEnv()
			defer cleanup()

			var testDir string
			if tt.name == "success_clears_cache_files_from_env_override_location" {
				testDir = os.Getenv(commoninternal.ArkBasicKeyringFolderEnvVar)
			} else {
				testDir = os.Getenv("HOME")
			}

			if tt.createFiles != nil {
				if err := tt.createFiles(testDir); err != nil {
					t.Fatalf("Failed to create test files: %v", err)
				}
			}

			action := NewArkCacheAction()
			cmd := &cobra.Command{}

			// Execute the function - should not panic
			action.runClearCacheAction(cmd, []string{})

			if tt.validateFunc != nil {
				tt.validateFunc(t, testDir)
			}
		})
	}
}

func TestArkCacheAction_StructFields(t *testing.T) {
	tests := []struct {
		name         string
		validateFunc func(t *testing.T, action *ArkCacheAction)
	}{
		{
			name: "success_struct_embeds_arkbaseaction",
			validateFunc: func(t *testing.T, action *ArkCacheAction) {
				actionValue := reflect.ValueOf(action).Elem()
				actionType := actionValue.Type()

				// Check that it embeds ArkBaseAction
				found := false
				for i := 0; i < actionType.NumField(); i++ {
					field := actionType.Field(i)
					if field.Type.String() == "*actions.ArkBaseAction" && field.Anonymous {
						found = true
						break
					}
				}
				if !found {
					t.Error("Expected ArkCacheAction to embed *ArkBaseAction")
				}
			},
		},
		{
			name: "success_implements_arkaction_interface",
			validateFunc: func(t *testing.T, action *ArkCacheAction) {
				// Verify it implements ArkAction interface by checking method exists
				actionValue := reflect.ValueOf(action)
				actionType := actionValue.Type()

				method, exists := actionType.MethodByName("DefineAction")
				if !exists {
					t.Error("Expected DefineAction method to exist")
					return
				}

				// Check method signature: func(cmd *cobra.Command)
				if method.Type.NumIn() != 2 { // receiver + parameter
					t.Errorf("DefineAction should have 1 parameter, got %d", method.Type.NumIn()-1)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			action := NewArkCacheAction()

			if tt.validateFunc != nil {
				tt.validateFunc(t, action)
			}
		})
	}
}

func TestArkCacheAction_IntegrationWithArkAction(t *testing.T) {
	tests := []struct {
		name         string
		validateFunc func(t *testing.T, action *ArkCacheAction)
	}{
		{
			name: "success_can_be_used_as_arkaction_interface",
			validateFunc: func(t *testing.T, action *ArkCacheAction) {
				// This should compile if ArkCacheAction implements ArkAction
				var arkAction ArkAction = action

				// Test that we can call the interface method
				rootCmd := &cobra.Command{Use: "test"}
				arkAction.DefineAction(rootCmd)

				// Verify the command was added
				cacheCmd, _, err := rootCmd.Find([]string{"cache"})
				if err != nil {
					t.Errorf("Expected to find cache command after DefineAction call, got error: %v", err)
				}
				if cacheCmd == nil {
					t.Error("Expected cache command to be added through interface call")
				}
			},
		},
		{
			name: "success_inherits_common_action_methods",
			validateFunc: func(t *testing.T, action *ArkCacheAction) {
				// Verify that methods from ArkBaseAction are accessible
				cmd := &cobra.Command{}

				// This should not panic and should add common flags
				action.CommonActionsConfiguration(cmd)

				// Check that common flags were added
				if flag := cmd.PersistentFlags().Lookup("verbose"); flag == nil {
					t.Error("Expected to inherit CommonActionsConfiguration from ArkBaseAction")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			action := NewArkCacheAction()

			if tt.validateFunc != nil {
				tt.validateFunc(t, action)
			}
		})
	}
}
