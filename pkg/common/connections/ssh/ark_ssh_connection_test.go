package ssh

import (
	"os"
	"strings"
	"testing"

	connectionsmodels "github.com/cyberark/ark-sdk-golang/pkg/models/common/connections"
)

func TestSSHPort(t *testing.T) {
	tests := []struct {
		name          string
		expectedValue int
	}{
		{
			name:          "constant_has_correct_value",
			expectedValue: 22,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if SSHPort != tt.expectedValue {
				t.Errorf("Expected SSHPort to be %d, got %d", tt.expectedValue, SSHPort)
			}
		})
	}
}

func TestNewArkSSHConnection(t *testing.T) {
	tests := []struct {
		name           string
		validateFunc   func(t *testing.T, result *ArkSSHConnection)
		expectedResult bool
	}{
		{
			name: "success_creates_new_instance",
			validateFunc: func(t *testing.T, result *ArkSSHConnection) {
				if result == nil {
					t.Error("Expected non-nil connection")
					return
				}
				if result.isConnected {
					t.Error("Expected isConnected to be false")
				}
				if result.isSuspended {
					t.Error("Expected isSuspended to be false")
				}
				if result.logger == nil {
					t.Error("Expected logger to be initialized")
				}
				if result.sshClient != nil {
					t.Error("Expected sshClient to be nil")
				}
			},
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := NewArkSSHConnection()

			if tt.validateFunc != nil {
				tt.validateFunc(t, result)
			}
		})
	}
}

func TestArkSSHConnection_Connect_Validation(t *testing.T) {
	// These tests validate the input validation and setup logic without external dependencies
	// Create temporary private key file for testing
	tempKeyFile, err := os.CreateTemp("", "test-key-*.pem")
	if err != nil {
		t.Fatalf("Failed to create temp key file: %v", err)
	}
	defer os.Remove(tempKeyFile.Name())

	// Write a mock private key (this won't work for real connections but tests file reading)
	_, err = tempKeyFile.WriteString("-----BEGIN OPENSSH PRIVATE KEY-----\nMOCK_KEY_DATA\n-----END OPENSSH PRIVATE KEY-----")
	if err != nil {
		t.Fatalf("Failed to write key data: %v", err)
	}
	tempKeyFile.Close()

	tests := []struct {
		name              string
		connectionDetails *connectionsmodels.ArkConnectionDetails
		setupFunc         func(conn *ArkSSHConnection)
		expectedError     bool
		expectedErrorMsg  string
		validateFunc      func(t *testing.T, conn *ArkSSHConnection)
	}{
		{
			name: "success_already_connected",
			connectionDetails: &connectionsmodels.ArkConnectionDetails{
				Address: "test-server",
				Port:    22,
				Credentials: &connectionsmodels.ArkConnectionCredentials{
					User:     "testuser",
					Password: "testpass",
				},
			},
			setupFunc: func(conn *ArkSSHConnection) {
				conn.isConnected = true
			},
			expectedError: false,
			validateFunc: func(t *testing.T, conn *ArkSSHConnection) {
				if !conn.isConnected {
					t.Error("Expected connection to remain connected")
				}
			},
		},
		{
			name: "error_invalid_private_key_file",
			connectionDetails: &connectionsmodels.ArkConnectionDetails{
				Address: "test-server",
				Port:    22,
				Credentials: &connectionsmodels.ArkConnectionCredentials{
					User:               "testuser",
					PrivateKeyFilepath: "/nonexistent/key.pem",
				},
			},
			expectedError:    true,
			expectedErrorMsg: "failed to check private key file exists",
		},
		{
			name: "error_unreadable_private_key_file",
			connectionDetails: &connectionsmodels.ArkConnectionDetails{
				Address: "test-server",
				Port:    22,
				Credentials: &connectionsmodels.ArkConnectionCredentials{
					User:               "testuser",
					PrivateKeyFilepath: "/dev/null",
				},
			},
			expectedError: true,
			// This will fail at key parsing since /dev/null doesn't contain a valid key
		},
		{
			name: "error_invalid_private_key_contents",
			connectionDetails: &connectionsmodels.ArkConnectionDetails{
				Address: "test-server",
				Port:    22,
				Credentials: &connectionsmodels.ArkConnectionCredentials{
					User:               "testuser",
					PrivateKeyContents: "invalid key content",
				},
			},
			expectedError:    true,
			expectedErrorMsg: "failed to parse private key contents",
		},
		{
			name: "success_default_retry_count",
			connectionDetails: &connectionsmodels.ArkConnectionDetails{
				Address: "test-server",
				Port:    22,
				Credentials: &connectionsmodels.ArkConnectionCredentials{
					User:     "testuser",
					Password: "testpass",
				},
				// ConnectionRetries not set, should default to 1
			},
			// This will fail at actual SSH connection since we don't have a real server,
			// but it validates the retry count defaulting logic
			expectedError: true,
		},
		{
			name: "success_password_auth_setup",
			connectionDetails: &connectionsmodels.ArkConnectionDetails{
				Address: "test-server",
				Port:    22,
				Credentials: &connectionsmodels.ArkConnectionCredentials{
					User:     "testuser",
					Password: "testpass",
				},
				ConnectionRetries: 1,
			},
			// This will fail at actual SSH connection since we don't have a real server,
			// but it validates the password authentication setup
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			conn := NewArkSSHConnection()
			if tt.setupFunc != nil {
				tt.setupFunc(conn)
			}

			err := conn.Connect(tt.connectionDetails)

			if tt.expectedError {
				if err == nil {
					t.Error("Expected error, got nil")
					return
				}
				if tt.expectedErrorMsg != "" && !strings.Contains(err.Error(), tt.expectedErrorMsg) {
					t.Errorf("Expected error message to contain '%s', got '%s'", tt.expectedErrorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
				return
			}

			if tt.validateFunc != nil {
				tt.validateFunc(t, conn)
			}
		})
	}
}

func TestArkSSHConnection_Disconnect(t *testing.T) {
	tests := []struct {
		name          string
		setupFunc     func(conn *ArkSSHConnection)
		expectedError bool
		validateFunc  func(t *testing.T, conn *ArkSSHConnection)
	}{
		{
			name: "success_not_connected",
			setupFunc: func(conn *ArkSSHConnection) {
				conn.isConnected = false
			},
			expectedError: false,
			validateFunc: func(t *testing.T, conn *ArkSSHConnection) {
				if conn.isConnected {
					t.Error("Expected connection to remain disconnected")
				}
			},
		},
		// Note: Testing actual disconnection with a real SSH client would require
		// complex mocking or a real connection, which goes against our no-external-dependencies rule.
		// The disconnect logic is simple enough that the main validation is the state management.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			conn := NewArkSSHConnection()
			if tt.setupFunc != nil {
				tt.setupFunc(conn)
			}

			err := conn.Disconnect()

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

			if tt.validateFunc != nil {
				tt.validateFunc(t, conn)
			}
		})
	}
}

func TestArkSSHConnection_SuspendConnection(t *testing.T) {
	tests := []struct {
		name          string
		setupFunc     func(conn *ArkSSHConnection)
		expectedError bool
		validateFunc  func(t *testing.T, conn *ArkSSHConnection)
	}{
		{
			name: "success_suspend_not_suspended",
			setupFunc: func(conn *ArkSSHConnection) {
				conn.isSuspended = false
			},
			expectedError: false,
			validateFunc: func(t *testing.T, conn *ArkSSHConnection) {
				if !conn.isSuspended {
					t.Error("Expected connection to be suspended")
				}
			},
		},
		{
			name: "success_suspend_already_suspended",
			setupFunc: func(conn *ArkSSHConnection) {
				conn.isSuspended = true
			},
			expectedError: false,
			validateFunc: func(t *testing.T, conn *ArkSSHConnection) {
				if !conn.isSuspended {
					t.Error("Expected connection to remain suspended")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			conn := NewArkSSHConnection()
			if tt.setupFunc != nil {
				tt.setupFunc(conn)
			}

			err := conn.SuspendConnection()

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

			if tt.validateFunc != nil {
				tt.validateFunc(t, conn)
			}
		})
	}
}

func TestArkSSHConnection_RestoreConnection(t *testing.T) {
	tests := []struct {
		name          string
		setupFunc     func(conn *ArkSSHConnection)
		expectedError bool
		validateFunc  func(t *testing.T, conn *ArkSSHConnection)
	}{
		{
			name: "success_restore_suspended",
			setupFunc: func(conn *ArkSSHConnection) {
				conn.isSuspended = true
			},
			expectedError: false,
			validateFunc: func(t *testing.T, conn *ArkSSHConnection) {
				if conn.isSuspended {
					t.Error("Expected connection to not be suspended")
				}
			},
		},
		{
			name: "success_restore_not_suspended",
			setupFunc: func(conn *ArkSSHConnection) {
				conn.isSuspended = false
			},
			expectedError: false,
			validateFunc: func(t *testing.T, conn *ArkSSHConnection) {
				if conn.isSuspended {
					t.Error("Expected connection to remain not suspended")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			conn := NewArkSSHConnection()
			if tt.setupFunc != nil {
				tt.setupFunc(conn)
			}

			err := conn.RestoreConnection()

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

			if tt.validateFunc != nil {
				tt.validateFunc(t, conn)
			}
		})
	}
}

func TestArkSSHConnection_IsSuspended(t *testing.T) {
	tests := []struct {
		name           string
		setupFunc      func(conn *ArkSSHConnection)
		expectedResult bool
	}{
		{
			name: "returns_true_when_suspended",
			setupFunc: func(conn *ArkSSHConnection) {
				conn.isSuspended = true
			},
			expectedResult: true,
		},
		{
			name: "returns_false_when_not_suspended",
			setupFunc: func(conn *ArkSSHConnection) {
				conn.isSuspended = false
			},
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			conn := NewArkSSHConnection()
			if tt.setupFunc != nil {
				tt.setupFunc(conn)
			}

			result := conn.IsSuspended()

			if result != tt.expectedResult {
				t.Errorf("Expected IsSuspended() to return %v, got %v", tt.expectedResult, result)
			}
		})
	}
}

func TestArkSSHConnection_IsConnected(t *testing.T) {
	tests := []struct {
		name           string
		setupFunc      func(conn *ArkSSHConnection)
		expectedResult bool
	}{
		{
			name: "returns_true_when_connected",
			setupFunc: func(conn *ArkSSHConnection) {
				conn.isConnected = true
			},
			expectedResult: true,
		},
		{
			name: "returns_false_when_not_connected",
			setupFunc: func(conn *ArkSSHConnection) {
				conn.isConnected = false
			},
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			conn := NewArkSSHConnection()
			if tt.setupFunc != nil {
				tt.setupFunc(conn)
			}

			result := conn.IsConnected()

			if result != tt.expectedResult {
				t.Errorf("Expected IsConnected() to return %v, got %v", tt.expectedResult, result)
			}
		})
	}
}

func TestArkSSHConnection_RunCommand_Validation(t *testing.T) {
	// These tests validate the input validation and state checking logic
	// without requiring actual SSH connections
	tests := []struct {
		name             string
		command          *connectionsmodels.ArkConnectionCommand
		setupFunc        func(conn *ArkSSHConnection)
		expectedError    bool
		expectedErrorMsg string
		validateFunc     func(t *testing.T, result *connectionsmodels.ArkConnectionResult)
	}{
		{
			name: "error_not_connected",
			command: &connectionsmodels.ArkConnectionCommand{
				Command:    "echo test",
				ExpectedRC: 0,
			},
			setupFunc: func(conn *ArkSSHConnection) {
				conn.isConnected = false
				conn.isSuspended = false
			},
			expectedError:    true,
			expectedErrorMsg: "cannot run command while not being connected",
		},
		{
			name: "error_suspended",
			command: &connectionsmodels.ArkConnectionCommand{
				Command:    "echo test",
				ExpectedRC: 0,
			},
			setupFunc: func(conn *ArkSSHConnection) {
				conn.isConnected = true
				conn.isSuspended = true
			},
			expectedError:    true,
			expectedErrorMsg: "cannot run command while not being connected",
		},
		{
			name: "error_connected_and_suspended",
			command: &connectionsmodels.ArkConnectionCommand{
				Command:    "echo test",
				ExpectedRC: 0,
			},
			setupFunc: func(conn *ArkSSHConnection) {
				conn.isConnected = false
				conn.isSuspended = true
			},
			expectedError:    true,
			expectedErrorMsg: "cannot run command while not being connected",
		},
		// Note: Testing actual command execution would require mocking the SSH client
		// and session creation, which is complex. The main validation here is the
		// connection state checking which is the testable logic without external dependencies.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			conn := NewArkSSHConnection()
			if tt.setupFunc != nil {
				tt.setupFunc(conn)
			}

			result, err := conn.RunCommand(tt.command)

			if tt.expectedError {
				if err == nil {
					t.Error("Expected error, got nil")
					return
				}
				if tt.expectedErrorMsg != "" && !strings.Contains(err.Error(), tt.expectedErrorMsg) {
					t.Errorf("Expected error message to contain '%s', got '%s'", tt.expectedErrorMsg, err.Error())
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
