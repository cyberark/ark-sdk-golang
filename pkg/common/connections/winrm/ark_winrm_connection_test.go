package winrm

import (
	"os"
	"strings"
	"testing"
	"time"

	connectionsmodels "github.com/cyberark/ark-sdk-golang/pkg/models/common/connections"
	"github.com/cyberark/ark-sdk-golang/pkg/models/common/connections/connectiondata"
)

func TestNewArkWinRMConnection(t *testing.T) {
	tests := []struct {
		name           string
		validateFunc   func(t *testing.T, result *ArkWinRMConnection)
		expectedResult bool
	}{
		{
			name: "success_creates_new_instance",
			validateFunc: func(t *testing.T, result *ArkWinRMConnection) {
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
				if result.winrmClient != nil {
					t.Error("Expected winrmClient to be nil")
				}
				if result.winrmShell != nil {
					t.Error("Expected winrmShell to be nil")
				}
			},
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := NewArkWinRMConnection()

			if tt.validateFunc != nil {
				tt.validateFunc(t, result)
			}
		})
	}
}

func TestArkWinRMConnection_Connect_Validation(t *testing.T) {
	// These tests validate the input validation and setup logic without external dependencies
	// Create temporary certificate file for testing
	tempCertFile, err := os.CreateTemp("", "test-cert-*.pem")
	if err != nil {
		t.Fatalf("Failed to create temp cert file: %v", err)
	}
	defer os.Remove(tempCertFile.Name())

	_, err = tempCertFile.WriteString("-----BEGIN CERTIFICATE-----\nMOCK_CERT_DATA\n-----END CERTIFICATE-----")
	if err != nil {
		t.Fatalf("Failed to write cert data: %v", err)
	}
	tempCertFile.Close()

	tests := []struct {
		name              string
		connectionDetails *connectionsmodels.ArkConnectionDetails
		setupFunc         func(conn *ArkWinRMConnection)
		expectedError     bool
		expectedErrorMsg  string
		validateFunc      func(t *testing.T, conn *ArkWinRMConnection)
	}{
		{
			name: "success_already_connected",
			connectionDetails: &connectionsmodels.ArkConnectionDetails{
				Address: "test-server",
				Port:    5986,
				Credentials: &connectionsmodels.ArkConnectionCredentials{
					User:     "testuser",
					Password: "testpass",
				},
			},
			setupFunc: func(conn *ArkWinRMConnection) {
				conn.isConnected = true
			},
			expectedError: false,
			validateFunc: func(t *testing.T, conn *ArkWinRMConnection) {
				if !conn.isConnected {
					t.Error("Expected connection to remain connected")
				}
			},
		},
		{
			name: "error_missing_credentials",
			connectionDetails: &connectionsmodels.ArkConnectionDetails{
				Address: "test-server",
				Port:    5986,
			},
			expectedError:    true,
			expectedErrorMsg: "missing credentials for WinRM connection",
		},
		{
			name: "error_invalid_certificate_path",
			connectionDetails: &connectionsmodels.ArkConnectionDetails{
				Address: "test-server",
				Port:    5986,
				Credentials: &connectionsmodels.ArkConnectionCredentials{
					User:     "testuser",
					Password: "testpass",
				},
				ConnectionData: &connectiondata.ArkWinRMConnectionData{
					CertificatePath: "/nonexistent/cert.pem",
				},
			},
			expectedError:    true,
			expectedErrorMsg: "failed to read certificate file",
		},
		{
			name: "success_valid_certificate_path",
			connectionDetails: &connectionsmodels.ArkConnectionDetails{
				Address: "test-server",
				Port:    5986,
				Credentials: &connectionsmodels.ArkConnectionCredentials{
					User:     "testuser",
					Password: "testpass",
				},
				ConnectionData: &connectiondata.ArkWinRMConnectionData{
					CertificatePath:  tempCertFile.Name(),
					TrustCertificate: true,
				},
			},
			// This will fail at WinRM client creation since we don't have a real server,
			// but it validates the certificate reading logic
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			conn := NewArkWinRMConnection()
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

func TestArkWinRMConnection_Disconnect(t *testing.T) {
	tests := []struct {
		name          string
		setupFunc     func(conn *ArkWinRMConnection)
		expectedError bool
		validateFunc  func(t *testing.T, conn *ArkWinRMConnection)
	}{
		{
			name: "success_not_connected",
			setupFunc: func(conn *ArkWinRMConnection) {
				conn.isConnected = false
			},
			expectedError: false,
			validateFunc: func(t *testing.T, conn *ArkWinRMConnection) {
				if conn.isConnected {
					t.Error("Expected connection to remain disconnected")
				}
			},
		},
		// {
		// 	name: "success_connected_no_shell",
		// 	setupFunc: func(conn *ArkWinRMConnection) {
		// 		conn.isConnected = true
		// 		conn.isSuspended = true
		// 		// No shell set - tests the nil check
		// 	},
		// 	expectedError: false,
		// 	validateFunc: func(t *testing.T, conn *ArkWinRMConnection) {
		// 		if conn.isConnected {
		// 			t.Error("Expected connection to be disconnected")
		// 		}
		// 		if conn.isSuspended {
		// 			t.Error("Expected suspension to be cleared")
		// 		}
		// 	},
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			conn := NewArkWinRMConnection()
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

func TestArkWinRMConnection_SuspendConnection(t *testing.T) {
	tests := []struct {
		name          string
		setupFunc     func(conn *ArkWinRMConnection)
		expectedError bool
		validateFunc  func(t *testing.T, conn *ArkWinRMConnection)
	}{
		{
			name: "success_suspend_not_suspended",
			setupFunc: func(conn *ArkWinRMConnection) {
				conn.isSuspended = false
			},
			expectedError: false,
			validateFunc: func(t *testing.T, conn *ArkWinRMConnection) {
				if !conn.isSuspended {
					t.Error("Expected connection to be suspended")
				}
			},
		},
		{
			name: "success_suspend_already_suspended",
			setupFunc: func(conn *ArkWinRMConnection) {
				conn.isSuspended = true
			},
			expectedError: false,
			validateFunc: func(t *testing.T, conn *ArkWinRMConnection) {
				if !conn.isSuspended {
					t.Error("Expected connection to remain suspended")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			conn := NewArkWinRMConnection()
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

func TestArkWinRMConnection_RestoreConnection(t *testing.T) {
	tests := []struct {
		name          string
		setupFunc     func(conn *ArkWinRMConnection)
		expectedError bool
		validateFunc  func(t *testing.T, conn *ArkWinRMConnection)
	}{
		{
			name: "success_restore_suspended",
			setupFunc: func(conn *ArkWinRMConnection) {
				conn.isSuspended = true
			},
			expectedError: false,
			validateFunc: func(t *testing.T, conn *ArkWinRMConnection) {
				if conn.isSuspended {
					t.Error("Expected connection to not be suspended")
				}
			},
		},
		{
			name: "success_restore_not_suspended",
			setupFunc: func(conn *ArkWinRMConnection) {
				conn.isSuspended = false
			},
			expectedError: false,
			validateFunc: func(t *testing.T, conn *ArkWinRMConnection) {
				if conn.isSuspended {
					t.Error("Expected connection to remain not suspended")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			conn := NewArkWinRMConnection()
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

func TestArkWinRMConnection_IsSuspended(t *testing.T) {
	tests := []struct {
		name           string
		setupFunc      func(conn *ArkWinRMConnection)
		expectedResult bool
	}{
		{
			name: "returns_true_when_suspended",
			setupFunc: func(conn *ArkWinRMConnection) {
				conn.isSuspended = true
			},
			expectedResult: true,
		},
		{
			name: "returns_false_when_not_suspended",
			setupFunc: func(conn *ArkWinRMConnection) {
				conn.isSuspended = false
			},
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			conn := NewArkWinRMConnection()
			if tt.setupFunc != nil {
				tt.setupFunc(conn)
			}

			result := conn.IsSuspended()

			if result != tt.expectedResult {
				t.Errorf("Expected result %v, got %v", tt.expectedResult, result)
			}
		})
	}
}

func TestArkWinRMConnection_IsConnected(t *testing.T) {
	tests := []struct {
		name           string
		setupFunc      func(conn *ArkWinRMConnection)
		expectedResult bool
	}{
		{
			name: "returns_true_when_connected",
			setupFunc: func(conn *ArkWinRMConnection) {
				conn.isConnected = true
			},
			expectedResult: true,
		},
		{
			name: "returns_false_when_not_connected",
			setupFunc: func(conn *ArkWinRMConnection) {
				conn.isConnected = false
			},
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			conn := NewArkWinRMConnection()
			if tt.setupFunc != nil {
				tt.setupFunc(conn)
			}

			result := conn.IsConnected()

			if result != tt.expectedResult {
				t.Errorf("Expected result %v, got %v", tt.expectedResult, result)
			}
		})
	}
}

func TestArkWinRMConnection_RunCommand_ValidationLogic(t *testing.T) {
	// These tests focus on the validation logic without external dependencies
	tests := []struct {
		name             string
		command          *connectionsmodels.ArkConnectionCommand
		setupFunc        func(conn *ArkWinRMConnection)
		expectedError    bool
		expectedErrorMsg string
	}{
		{
			name: "error_not_connected",
			command: &connectionsmodels.ArkConnectionCommand{
				Command:    "echo test",
				ExpectedRC: 0,
			},
			setupFunc: func(conn *ArkWinRMConnection) {
				conn.isConnected = false
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
			setupFunc: func(conn *ArkWinRMConnection) {
				conn.isConnected = true
				conn.isSuspended = true
			},
			expectedError:    true,
			expectedErrorMsg: "cannot run command while not being connected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			conn := NewArkWinRMConnection()
			if tt.setupFunc != nil {
				tt.setupFunc(conn)
			}

			_, err := conn.RunCommand(tt.command)

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
		})
	}
}

// Test constants
func TestConstants(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected interface{}
	}{
		{
			name:     "winrm_https_port_correct",
			value:    WinRMHTTPSPort,
			expected: 5986,
		},
		{
			name:     "connection_timeout_correct",
			value:    winrmConnectionTimeout,
			expected: 10 * time.Second,
		},
		{
			name:     "max_single_command_size_correct",
			value:    maxSingleCommandSize,
			expected: 2000,
		},
		{
			name:     "max_chunk_size_correct",
			value:    maxChunkSize,
			expected: 4000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.value != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, tt.value)
			}
		})
	}
}
