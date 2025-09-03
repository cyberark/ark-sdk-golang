package common

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestNewBasicKeyring(t *testing.T) {
	tests := []struct {
		name         string
		setupFunc    func() (string, func()) // Returns temp dir and cleanup func
		envVar       string
		expectedNil  bool
		validateFunc func(t *testing.T, keyring *BasicKeyring, tempDir string)
	}{
		{
			name: "success_creates_keyring_with_default_folder",
			setupFunc: func() (string, func()) {
				tempDir, _ := os.MkdirTemp("", "ark_test_home_")
				originalHome := os.Getenv("HOME")
				os.Setenv("HOME", tempDir)
				os.Unsetenv(ArkBasicKeyringFolderEnvVar)
				return tempDir, func() {
					os.Setenv("HOME", originalHome)
					os.RemoveAll(tempDir)
				}
			},
			expectedNil: false,
			validateFunc: func(t *testing.T, keyring *BasicKeyring, tempDir string) {
				expectedPath := filepath.Join(tempDir, DefaultBasicKeyringFolder)
				if keyring.basicFolderPath != expectedPath {
					t.Errorf("Expected basicFolderPath '%s', got '%s'", expectedPath, keyring.basicFolderPath)
				}
				if keyring.keyringFilePath != filepath.Join(expectedPath, "keyring") {
					t.Errorf("Expected keyringFilePath to end with 'keyring', got '%s'", keyring.keyringFilePath)
				}
				if keyring.macFilePath != filepath.Join(expectedPath, "mac") {
					t.Errorf("Expected macFilePath to end with 'mac', got '%s'", keyring.macFilePath)
				}
				// Verify folder was created
				if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
					t.Error("Expected keyring folder to be created")
				}
			},
		},
		{
			name: "success_creates_keyring_with_env_var_folder",
			setupFunc: func() (string, func()) {
				tempDir, _ := os.MkdirTemp("", "ark_test_custom_")
				customPath := filepath.Join(tempDir, "custom_keyring")
				os.Setenv(ArkBasicKeyringFolderEnvVar, customPath)
				return customPath, func() {
					os.Unsetenv(ArkBasicKeyringFolderEnvVar)
					os.RemoveAll(tempDir)
				}
			},
			expectedNil: false,
			validateFunc: func(t *testing.T, keyring *BasicKeyring, tempDir string) {
				if keyring.basicFolderPath != tempDir {
					t.Errorf("Expected basicFolderPath '%s', got '%s'", tempDir, keyring.basicFolderPath)
				}
				// Verify folder was created
				if _, err := os.Stat(tempDir); os.IsNotExist(err) {
					t.Error("Expected custom keyring folder to be created")
				}
			},
		},
		{
			name: "success_handles_existing_folder",
			setupFunc: func() (string, func()) {
				tempDir, _ := os.MkdirTemp("", "ark_test_existing_")
				existingPath := filepath.Join(tempDir, "existing_keyring")
				os.MkdirAll(existingPath, 0755)
				os.Setenv(ArkBasicKeyringFolderEnvVar, existingPath)
				return existingPath, func() {
					os.Unsetenv(ArkBasicKeyringFolderEnvVar)
					os.RemoveAll(tempDir)
				}
			},
			expectedNil: false,
			validateFunc: func(t *testing.T, keyring *BasicKeyring, tempDir string) {
				if keyring.basicFolderPath != tempDir {
					t.Errorf("Expected basicFolderPath '%s', got '%s'", tempDir, keyring.basicFolderPath)
				}
			},
		},
		{
			name: "error_returns_nil_on_folder_creation_failure",
			setupFunc: func() (string, func()) {
				// Create a path that will fail to create (invalid characters)
				invalidPath := "/proc/invalid/path/that/cannot/be/created"
				os.Setenv(ArkBasicKeyringFolderEnvVar, invalidPath)
				return invalidPath, func() {
					os.Unsetenv(ArkBasicKeyringFolderEnvVar)
				}
			},
			expectedNil: true,
			validateFunc: func(t *testing.T, keyring *BasicKeyring, tempDir string) {
				// No validation needed for nil case
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test environment
			tempDir, cleanup := tt.setupFunc()
			defer cleanup()

			// Execute function
			keyring := NewBasicKeyring()

			// Validate result
			if tt.expectedNil {
				if keyring != nil {
					t.Errorf("Expected nil keyring, got %+v", keyring)
				}
				return
			}

			if keyring == nil {
				t.Error("Expected non-nil keyring")
				return
			}

			// Run custom validation
			if tt.validateFunc != nil {
				tt.validateFunc(t, keyring, tempDir)
			}
		})
	}
}

func TestBasicKeyring_SetPassword(t *testing.T) {
	tests := []struct {
		name          string
		setupFunc     func() (*BasicKeyring, func())
		serviceName   string
		username      string
		password      string
		expectedError bool
		validateFunc  func(t *testing.T, keyring *BasicKeyring)
	}{
		{
			name: "success_sets_password_new_keyring",
			setupFunc: func() (*BasicKeyring, func()) {
				tempDir, _ := os.MkdirTemp("", "ark_test_set_")
				os.Setenv(ArkBasicKeyringFolderEnvVar, tempDir)
				keyring := NewBasicKeyring()
				return keyring, func() {
					os.Unsetenv(ArkBasicKeyringFolderEnvVar)
					os.RemoveAll(tempDir)
				}
			},
			serviceName:   "github",
			username:      "testuser",
			password:      "testpassword",
			expectedError: false,
			validateFunc: func(t *testing.T, keyring *BasicKeyring) {
				// Verify keyring file exists
				if _, err := os.Stat(keyring.keyringFilePath); os.IsNotExist(err) {
					t.Error("Expected keyring file to be created")
				}
				// Verify MAC file exists
				if _, err := os.Stat(keyring.macFilePath); os.IsNotExist(err) {
					t.Error("Expected MAC file to be created")
				}
			},
		},
		{
			name: "success_sets_password_existing_keyring",
			setupFunc: func() (*BasicKeyring, func()) {
				tempDir, _ := os.MkdirTemp("", "ark_test_set_existing_")
				os.Setenv(ArkBasicKeyringFolderEnvVar, tempDir)
				keyring := NewBasicKeyring()
				// Set an initial password
				keyring.SetPassword("service1", "user1", "pass1")
				return keyring, func() {
					os.Unsetenv(ArkBasicKeyringFolderEnvVar)
					os.RemoveAll(tempDir)
				}
			},
			serviceName:   "service2",
			username:      "user2",
			password:      "pass2",
			expectedError: false,
			validateFunc: func(t *testing.T, keyring *BasicKeyring) {
				// Verify both passwords can be retrieved
				pass1, err := keyring.GetPassword("service1", "user1")
				if err != nil {
					t.Errorf("Error retrieving first password: %v", err)
				}
				if pass1 != "pass1" {
					t.Errorf("Expected first password 'pass1', got '%s'", pass1)
				}
				pass2, err := keyring.GetPassword("service2", "user2")
				if err != nil {
					t.Errorf("Error retrieving second password: %v", err)
				}
				if pass2 != "pass2" {
					t.Errorf("Expected second password 'pass2', got '%s'", pass2)
				}
			},
		},
		{
			name: "success_overwrites_existing_password",
			setupFunc: func() (*BasicKeyring, func()) {
				tempDir, _ := os.MkdirTemp("", "ark_test_overwrite_")
				os.Setenv(ArkBasicKeyringFolderEnvVar, tempDir)
				keyring := NewBasicKeyring()
				// Set an initial password
				keyring.SetPassword("github", "testuser", "oldpassword")
				return keyring, func() {
					os.Unsetenv(ArkBasicKeyringFolderEnvVar)
					os.RemoveAll(tempDir)
				}
			},
			serviceName:   "github",
			username:      "testuser",
			password:      "newpassword",
			expectedError: false,
			validateFunc: func(t *testing.T, keyring *BasicKeyring) {
				// Verify new password is retrieved
				password, err := keyring.GetPassword("github", "testuser")
				if err != nil {
					t.Errorf("Error retrieving password: %v", err)
				}
				if password != "newpassword" {
					t.Errorf("Expected password 'newpassword', got '%s'", password)
				}
			},
		},
		{
			name: "edge_case_empty_service_name",
			setupFunc: func() (*BasicKeyring, func()) {
				tempDir, _ := os.MkdirTemp("", "ark_test_empty_service_")
				os.Setenv(ArkBasicKeyringFolderEnvVar, tempDir)
				keyring := NewBasicKeyring()
				return keyring, func() {
					os.Unsetenv(ArkBasicKeyringFolderEnvVar)
					os.RemoveAll(tempDir)
				}
			},
			serviceName:   "",
			username:      "testuser",
			password:      "testpassword",
			expectedError: false,
			validateFunc: func(t *testing.T, keyring *BasicKeyring) {
				// Verify password can be retrieved with empty service name
				password, err := keyring.GetPassword("", "testuser")
				if err != nil {
					t.Errorf("Error retrieving password: %v", err)
				}
				if password != "testpassword" {
					t.Errorf("Expected password 'testpassword', got '%s'", password)
				}
			},
		},
		{
			name: "edge_case_empty_username",
			setupFunc: func() (*BasicKeyring, func()) {
				tempDir, _ := os.MkdirTemp("", "ark_test_empty_user_")
				os.Setenv(ArkBasicKeyringFolderEnvVar, tempDir)
				keyring := NewBasicKeyring()
				return keyring, func() {
					os.Unsetenv(ArkBasicKeyringFolderEnvVar)
					os.RemoveAll(tempDir)
				}
			},
			serviceName:   "github",
			username:      "",
			password:      "testpassword",
			expectedError: false,
			validateFunc: func(t *testing.T, keyring *BasicKeyring) {
				// Verify password can be retrieved with empty username
				password, err := keyring.GetPassword("github", "")
				if err != nil {
					t.Errorf("Error retrieving password: %v", err)
				}
				if password != "testpassword" {
					t.Errorf("Expected password 'testpassword', got '%s'", password)
				}
			},
		},
		{
			name: "edge_case_empty_password",
			setupFunc: func() (*BasicKeyring, func()) {
				tempDir, _ := os.MkdirTemp("", "ark_test_empty_pass_")
				os.Setenv(ArkBasicKeyringFolderEnvVar, tempDir)
				keyring := NewBasicKeyring()
				return keyring, func() {
					os.Unsetenv(ArkBasicKeyringFolderEnvVar)
					os.RemoveAll(tempDir)
				}
			},
			serviceName:   "github",
			username:      "testuser",
			password:      "",
			expectedError: false,
			validateFunc: func(t *testing.T, keyring *BasicKeyring) {
				// Verify empty password can be retrieved
				password, err := keyring.GetPassword("github", "testuser")
				if err != nil {
					t.Errorf("Error retrieving password: %v", err)
				}
				if password != "" {
					t.Errorf("Expected empty password, got '%s'", password)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test environment
			keyring, cleanup := tt.setupFunc()
			defer cleanup()

			if keyring == nil {
				t.Fatal("Failed to create keyring for test")
			}

			// Execute function
			err := keyring.SetPassword(tt.serviceName, tt.username, tt.password)

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

			// Run custom validation
			if tt.validateFunc != nil {
				tt.validateFunc(t, keyring)
			}
		})
	}
}

func TestBasicKeyring_GetPassword(t *testing.T) {
	tests := []struct {
		name             string
		setupFunc        func() (*BasicKeyring, func())
		serviceName      string
		username         string
		expectedPassword string
		expectedError    bool
	}{
		{
			name: "success_gets_existing_password",
			setupFunc: func() (*BasicKeyring, func()) {
				tempDir, _ := os.MkdirTemp("", "ark_test_get_")
				os.Setenv(ArkBasicKeyringFolderEnvVar, tempDir)
				keyring := NewBasicKeyring()
				keyring.SetPassword("github", "testuser", "testpassword")
				return keyring, func() {
					os.Unsetenv(ArkBasicKeyringFolderEnvVar)
					os.RemoveAll(tempDir)
				}
			},
			serviceName:      "github",
			username:         "testuser",
			expectedPassword: "testpassword",
			expectedError:    false,
		},
		{
			name: "success_returns_empty_for_nonexistent_keyring",
			setupFunc: func() (*BasicKeyring, func()) {
				tempDir, _ := os.MkdirTemp("", "ark_test_no_keyring_")
				os.Setenv(ArkBasicKeyringFolderEnvVar, tempDir)
				keyring := NewBasicKeyring()
				return keyring, func() {
					os.Unsetenv(ArkBasicKeyringFolderEnvVar)
					os.RemoveAll(tempDir)
				}
			},
			serviceName:      "github",
			username:         "testuser",
			expectedPassword: "",
			expectedError:    false,
		},
		{
			name: "success_returns_empty_for_nonexistent_service",
			setupFunc: func() (*BasicKeyring, func()) {
				tempDir, _ := os.MkdirTemp("", "ark_test_no_service_")
				os.Setenv(ArkBasicKeyringFolderEnvVar, tempDir)
				keyring := NewBasicKeyring()
				keyring.SetPassword("github", "testuser", "testpassword")
				return keyring, func() {
					os.Unsetenv(ArkBasicKeyringFolderEnvVar)
					os.RemoveAll(tempDir)
				}
			},
			serviceName:      "gitlab",
			username:         "testuser",
			expectedPassword: "",
			expectedError:    false,
		},
		{
			name: "success_returns_empty_for_nonexistent_username",
			setupFunc: func() (*BasicKeyring, func()) {
				tempDir, _ := os.MkdirTemp("", "ark_test_no_user_")
				os.Setenv(ArkBasicKeyringFolderEnvVar, tempDir)
				keyring := NewBasicKeyring()
				keyring.SetPassword("github", "testuser", "testpassword")
				return keyring, func() {
					os.Unsetenv(ArkBasicKeyringFolderEnvVar)
					os.RemoveAll(tempDir)
				}
			},
			serviceName:      "github",
			username:         "otheruser",
			expectedPassword: "",
			expectedError:    false,
		},
		{
			name: "success_gets_multiple_passwords",
			setupFunc: func() (*BasicKeyring, func()) {
				tempDir, _ := os.MkdirTemp("", "ark_test_multiple_")
				os.Setenv(ArkBasicKeyringFolderEnvVar, tempDir)
				keyring := NewBasicKeyring()
				keyring.SetPassword("github", "user1", "pass1")
				keyring.SetPassword("github", "user2", "pass2")
				keyring.SetPassword("gitlab", "user1", "pass3")
				return keyring, func() {
					os.Unsetenv(ArkBasicKeyringFolderEnvVar)
					os.RemoveAll(tempDir)
				}
			},
			serviceName:      "gitlab",
			username:         "user1",
			expectedPassword: "pass3",
			expectedError:    false,
		},
		{
			name: "edge_case_empty_service_name",
			setupFunc: func() (*BasicKeyring, func()) {
				tempDir, _ := os.MkdirTemp("", "ark_test_empty_service_get_")
				os.Setenv(ArkBasicKeyringFolderEnvVar, tempDir)
				keyring := NewBasicKeyring()
				keyring.SetPassword("", "testuser", "testpassword")
				return keyring, func() {
					os.Unsetenv(ArkBasicKeyringFolderEnvVar)
					os.RemoveAll(tempDir)
				}
			},
			serviceName:      "",
			username:         "testuser",
			expectedPassword: "testpassword",
			expectedError:    false,
		},
		{
			name: "edge_case_empty_username",
			setupFunc: func() (*BasicKeyring, func()) {
				tempDir, _ := os.MkdirTemp("", "ark_test_empty_user_get_")
				os.Setenv(ArkBasicKeyringFolderEnvVar, tempDir)
				keyring := NewBasicKeyring()
				keyring.SetPassword("github", "", "testpassword")
				return keyring, func() {
					os.Unsetenv(ArkBasicKeyringFolderEnvVar)
					os.RemoveAll(tempDir)
				}
			},
			serviceName:      "github",
			username:         "",
			expectedPassword: "testpassword",
			expectedError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test environment
			keyring, cleanup := tt.setupFunc()
			defer cleanup()

			if keyring == nil {
				t.Fatal("Failed to create keyring for test")
			}

			// Execute function
			password, err := keyring.GetPassword(tt.serviceName, tt.username)

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

			// Validate result
			if password != tt.expectedPassword {
				t.Errorf("Expected password '%s', got '%s'", tt.expectedPassword, password)
			}
		})
	}
}

func TestBasicKeyring_DeletePassword(t *testing.T) {
	tests := []struct {
		name          string
		setupFunc     func() (*BasicKeyring, func())
		serviceName   string
		username      string
		expectedError bool
		validateFunc  func(t *testing.T, keyring *BasicKeyring)
	}{
		{
			name: "success_deletes_existing_password",
			setupFunc: func() (*BasicKeyring, func()) {
				tempDir, _ := os.MkdirTemp("", "ark_test_delete_")
				os.Setenv(ArkBasicKeyringFolderEnvVar, tempDir)
				keyring := NewBasicKeyring()
				keyring.SetPassword("github", "testuser", "testpassword")
				return keyring, func() {
					os.Unsetenv(ArkBasicKeyringFolderEnvVar)
					os.RemoveAll(tempDir)
				}
			},
			serviceName:   "github",
			username:      "testuser",
			expectedError: false,
			validateFunc: func(t *testing.T, keyring *BasicKeyring) {
				// Verify password no longer exists
				password, err := keyring.GetPassword("github", "testuser")
				if err != nil {
					t.Errorf("Error checking deleted password: %v", err)
				}
				if password != "" {
					t.Errorf("Expected empty password after deletion, got '%s'", password)
				}
			},
		},
		{
			name: "success_idempotent_nonexistent_keyring",
			setupFunc: func() (*BasicKeyring, func()) {
				tempDir, _ := os.MkdirTemp("", "ark_test_delete_no_keyring_")
				os.Setenv(ArkBasicKeyringFolderEnvVar, tempDir)
				keyring := NewBasicKeyring()
				return keyring, func() {
					os.Unsetenv(ArkBasicKeyringFolderEnvVar)
					os.RemoveAll(tempDir)
				}
			},
			serviceName:   "github",
			username:      "testuser",
			expectedError: false,
			validateFunc: func(t *testing.T, keyring *BasicKeyring) {
				// No validation needed - should be idempotent
			},
		},
		{
			name: "success_idempotent_nonexistent_service",
			setupFunc: func() (*BasicKeyring, func()) {
				tempDir, _ := os.MkdirTemp("", "ark_test_delete_no_service_")
				os.Setenv(ArkBasicKeyringFolderEnvVar, tempDir)
				keyring := NewBasicKeyring()
				keyring.SetPassword("github", "testuser", "testpassword")
				return keyring, func() {
					os.Unsetenv(ArkBasicKeyringFolderEnvVar)
					os.RemoveAll(tempDir)
				}
			},
			serviceName:   "gitlab",
			username:      "testuser",
			expectedError: false,
			validateFunc: func(t *testing.T, keyring *BasicKeyring) {
				// Verify original password still exists
				password, err := keyring.GetPassword("github", "testuser")
				if err != nil {
					t.Errorf("Error checking original password: %v", err)
				}
				if password != "testpassword" {
					t.Errorf("Expected original password 'testpassword', got '%s'", password)
				}
			},
		},
		{
			name: "success_idempotent_nonexistent_username",
			setupFunc: func() (*BasicKeyring, func()) {
				tempDir, _ := os.MkdirTemp("", "ark_test_delete_no_user_")
				os.Setenv(ArkBasicKeyringFolderEnvVar, tempDir)
				keyring := NewBasicKeyring()
				keyring.SetPassword("github", "testuser", "testpassword")
				return keyring, func() {
					os.Unsetenv(ArkBasicKeyringFolderEnvVar)
					os.RemoveAll(tempDir)
				}
			},
			serviceName:   "github",
			username:      "otheruser",
			expectedError: false,
			validateFunc: func(t *testing.T, keyring *BasicKeyring) {
				// Verify original password still exists
				password, err := keyring.GetPassword("github", "testuser")
				if err != nil {
					t.Errorf("Error checking original password: %v", err)
				}
				if password != "testpassword" {
					t.Errorf("Expected original password 'testpassword', got '%s'", password)
				}
			},
		},
		{
			name: "success_deletes_one_of_multiple_passwords",
			setupFunc: func() (*BasicKeyring, func()) {
				tempDir, _ := os.MkdirTemp("", "ark_test_delete_multiple_")
				os.Setenv(ArkBasicKeyringFolderEnvVar, tempDir)
				keyring := NewBasicKeyring()
				keyring.SetPassword("github", "user1", "pass1")
				keyring.SetPassword("github", "user2", "pass2")
				keyring.SetPassword("gitlab", "user1", "pass3")
				return keyring, func() {
					os.Unsetenv(ArkBasicKeyringFolderEnvVar)
					os.RemoveAll(tempDir)
				}
			},
			serviceName:   "github",
			username:      "user1",
			expectedError: false,
			validateFunc: func(t *testing.T, keyring *BasicKeyring) {
				// Verify deleted password is gone
				password, err := keyring.GetPassword("github", "user1")
				if err != nil {
					t.Errorf("Error checking deleted password: %v", err)
				}
				if password != "" {
					t.Errorf("Expected empty password after deletion, got '%s'", password)
				}

				// Verify other passwords still exist
				password2, err := keyring.GetPassword("github", "user2")
				if err != nil {
					t.Errorf("Error checking remaining password: %v", err)
				}
				if password2 != "pass2" {
					t.Errorf("Expected password 'pass2', got '%s'", password2)
				}

				password3, err := keyring.GetPassword("gitlab", "user1")
				if err != nil {
					t.Errorf("Error checking remaining password: %v", err)
				}
				if password3 != "pass3" {
					t.Errorf("Expected password 'pass3', got '%s'", password3)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test environment
			keyring, cleanup := tt.setupFunc()
			defer cleanup()

			if keyring == nil {
				t.Fatal("Failed to create keyring for test")
			}

			// Execute function
			err := keyring.DeletePassword(tt.serviceName, tt.username)

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

			// Run custom validation
			if tt.validateFunc != nil {
				tt.validateFunc(t, keyring)
			}
		})
	}
}

func TestBasicKeyring_Integration(t *testing.T) {
	tests := []struct {
		name         string
		actions      []func(*BasicKeyring) error
		validateFunc func(t *testing.T, keyring *BasicKeyring)
	}{
		{
			name: "integration_complete_lifecycle",
			actions: []func(*BasicKeyring) error{
				func(k *BasicKeyring) error { return k.SetPassword("service1", "user1", "pass1") },
				func(k *BasicKeyring) error { return k.SetPassword("service1", "user2", "pass2") },
				func(k *BasicKeyring) error { return k.SetPassword("service2", "user1", "pass3") },
				func(k *BasicKeyring) error { return k.DeletePassword("service1", "user1") },
			},
			validateFunc: func(t *testing.T, keyring *BasicKeyring) {
				// Check deleted password
				pass1, err := keyring.GetPassword("service1", "user1")
				if err != nil {
					t.Errorf("Error getting deleted password: %v", err)
				}
				if pass1 != "" {
					t.Errorf("Expected empty password for deleted entry, got '%s'", pass1)
				}

				// Check remaining passwords
				pass2, err := keyring.GetPassword("service1", "user2")
				if err != nil {
					t.Errorf("Error getting password: %v", err)
				}
				if pass2 != "pass2" {
					t.Errorf("Expected 'pass2', got '%s'", pass2)
				}

				pass3, err := keyring.GetPassword("service2", "user1")
				if err != nil {
					t.Errorf("Error getting password: %v", err)
				}
				if pass3 != "pass3" {
					t.Errorf("Expected 'pass3', got '%s'", pass3)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test environment
			tempDir, _ := os.MkdirTemp("", "ark_test_integration_")
			os.Setenv(ArkBasicKeyringFolderEnvVar, tempDir)
			defer func() {
				os.Unsetenv(ArkBasicKeyringFolderEnvVar)
				os.RemoveAll(tempDir)
			}()

			keyring := NewBasicKeyring()
			if keyring == nil {
				t.Fatal("Failed to create keyring for test")
			}

			// Execute actions
			for i, action := range tt.actions {
				if err := action(keyring); err != nil {
					t.Errorf("Action %d failed: %v", i, err)
				}
			}

			// Run validation
			if tt.validateFunc != nil {
				tt.validateFunc(t, keyring)
			}
		})
	}
}

func TestConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant interface{}
		expected interface{}
	}{
		{
			name:     "nonce_size_correct_value",
			constant: nonceSize,
			expected: 16,
		},
		{
			name:     "tag_size_correct_value",
			constant: tagSize,
			expected: 16,
		},
		{
			name:     "block_size_correct_value",
			constant: blockSize,
			expected: 32,
		},
		{
			name:     "default_folder_correct_value",
			constant: DefaultBasicKeyringFolder,
			expected: ".ark_cache/keyring",
		},
		{
			name:     "env_var_correct_value",
			constant: ArkBasicKeyringFolderEnvVar,
			expected: "ARK_KEYRING_FOLDER",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if !reflect.DeepEqual(tt.constant, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, tt.constant)
			}
		})
	}
}
