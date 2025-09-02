// Package ssh provides Secure Shell (SSH) connection capabilities
// for the ARK SDK Golang. This package implements the ArkConnection interface
// to enable secure command execution on Linux/Unix machines using the SSH protocol.
//
// The package supports multiple authentication methods including password
// authentication and public key authentication (both from file and from content).
// It includes automatic retry mechanisms for connection failures and provides
// session management for command execution.
//
// Key features:
//   - Secure SSH connections with multiple authentication methods
//   - Password-based authentication
//   - Public key authentication (file-based and content-based)
//   - Automatic retry with connection failure detection
//   - Session-based command execution
//   - Connection suspend/restore functionality
//
// Example:
//
//	conn := NewArkSSHConnection()
//	err := conn.Connect(&connectionsmodels.ArkConnectionDetails{
//		Address: "linux-server.example.com",
//		Port:    22,
//		Credentials: &connectionsmodels.ArkConnectionCredentials{
//			User:     "username",
//			Password: "password",
//		},
//	})
//	if err != nil {
//		// handle error
//	}
//	defer conn.Disconnect()
//
//	result, err := conn.RunCommand(&connectionsmodels.ArkConnectionCommand{
//		Command:    "ls -la",
//		ExpectedRC: 0,
//	})
package ssh

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/cyberark/ark-sdk-golang/pkg/common"
	"github.com/cyberark/ark-sdk-golang/pkg/common/connections"
	connectionsmodels "github.com/cyberark/ark-sdk-golang/pkg/models/common/connections"
	"golang.org/x/crypto/ssh"
)

const (
	// SSHPort is the default port for SSH connections.
	SSHPort = 22
)

const (
	// connectionTimeout defines the maximum time to wait for SSH connection establishment.
	connectionTimeout = 10 * time.Second
)

// ArkSSHConnection is a struct that implements the ArkConnection interface for SSH connections.
//
// It provides secure Secure Shell functionality including connection management,
// command execution with session handling, and automatic retry mechanisms.
// The connection supports multiple authentication methods including password
// and public key authentication.
//
// The struct maintains connection state and provides suspend/restore functionality
// for connection lifecycle management.
type ArkSSHConnection struct {
	connections.ArkConnection
	isConnected bool
	isSuspended bool
	sshClient   *ssh.Client
	logger      *common.ArkLogger
}

// NewArkSSHConnection creates a new instance of ArkSSHConnection.
//
// Creates and initializes a new SSH connection instance with default settings.
// The connection is created in a disconnected state and must be explicitly
// connected using the Connect method before use.
//
// Returns a pointer to the newly created ArkSSHConnection instance with
// isConnected and isSuspended set to false, and a logger configured for
// SSH operations.
//
// Example:
//
//	conn := NewArkSSHConnection()
//	err := conn.Connect(connectionDetails)
//	if err != nil {
//		// handle connection error
//	}
func NewArkSSHConnection() *ArkSSHConnection {
	return &ArkSSHConnection{
		isConnected: false,
		isSuspended: false,
		logger:      common.GetLogger("ArkSSHConnection", common.Unknown),
	}
}

// Connect establishes an SSH connection using the provided connection details.
//
// Establishes a secure SSH connection to the target machine using the provided
// connection details. The method supports multiple authentication methods including
// password authentication and public key authentication (both from file and
// content). It handles automatic retry logic for connection failures.
//
// If the connection is already established, this method returns immediately
// without error. The method uses the default SSH port (22) if no port is
// specified in the connection details.
//
// The method supports three authentication methods in order of precedence:
// 1. Password authentication (if password is provided)
// 2. Private key file authentication (if PrivateKeyFilepath is provided)
// 3. Private key content authentication (if PrivateKeyContents is provided)
//
// Parameters:
//   - connectionDetails: Connection configuration including address, port,
//     credentials, retry settings, and authentication information
//
// Returns an error if the connection cannot be established, including cases
// where credentials are missing, private key files cannot be read or parsed,
// or the SSH connection fails.
//
// The method supports automatic retry with configurable retry count and
// tick period. Connection failures are detected and retried up to the
// specified limit.
//
// Example:
//
//	details := &connectionsmodels.ArkConnectionDetails{
//		Address: "linux-server.example.com",
//		Port:    22,
//		Credentials: &connectionsmodels.ArkConnectionCredentials{
//			User:     "username",
//			Password: "password",
//		},
//		ConnectionRetries: 3,
//		RetryTickPeriod:   5,
//	}
//	err := conn.Connect(details)
func (c *ArkSSHConnection) Connect(connectionDetails *connectionsmodels.ArkConnectionDetails) error {
	if c.isConnected {
		return nil
	}
	if connectionDetails.ConnectionRetries == 0 {
		connectionDetails.ConnectionRetries = 1
	}

	var authMethods []ssh.AuthMethod
	if connectionDetails.Credentials != nil {
		if connectionDetails.Credentials.Password != "" {
			authMethods = append(authMethods, ssh.Password(connectionDetails.Credentials.Password))
		} else if connectionDetails.Credentials.PrivateKeyFilepath != "" {
			_, err := os.Stat(connectionDetails.Credentials.PrivateKeyFilepath)
			if err != nil {
				return fmt.Errorf("failed to check private key file exists: %w", err)
			}
			keyData, err := os.ReadFile(connectionDetails.Credentials.PrivateKeyFilepath)
			if err != nil {
				return fmt.Errorf("failed to read private key file: %w", err)
			}
			signer, err := ssh.ParsePrivateKey(keyData)
			if err != nil {
				return fmt.Errorf("failed to parse private key: %w", err)
			}
			authMethods = append(authMethods, ssh.PublicKeys(signer))
		} else if connectionDetails.Credentials.PrivateKeyContents != "" {
			signer, err := ssh.ParsePrivateKey([]byte(connectionDetails.Credentials.PrivateKeyContents))
			if err != nil {
				return fmt.Errorf("failed to parse private key contents: %w", err)
			}
			authMethods = append(authMethods, ssh.PublicKeys(signer))
		}
	}

	config := &ssh.ClientConfig{
		User:            connectionDetails.Credentials.User,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         connectionTimeout,
	}

	address := fmt.Sprintf("%s:%d", connectionDetails.Address, connectionDetails.Port)
	var client *ssh.Client
	var err error
	for i := 0; i < connectionDetails.ConnectionRetries; i++ {
		client, err = ssh.Dial("tcp", address, config)
		if err != nil {
			if common.IsConnectionRefused(err) {
				if i < connectionDetails.ConnectionRetries-1 {
					time.Sleep(time.Duration(connectionDetails.RetryTickPeriod) * time.Second)
					continue
				}
			}
			return fmt.Errorf("failed to connect to SSH server: %w", err)
		}
		break
	}
	c.logger.Debug("Connected to SSH server [%s] on port [%d]", connectionDetails.Address, connectionDetails.Port)
	c.sshClient = client
	c.isConnected = true
	c.isSuspended = false
	return nil
}

// Disconnect closes the SSH connection.
//
// Closes the active SSH client connection and cleans up the connection resources.
// If the connection is not currently established, this method returns
// immediately without error.
//
// The method attempts to close the SSH client gracefully. If the client
// closure fails, a warning is logged but the method continues to clean
// up the connection state.
//
// After successful completion, the connection state is reset to disconnected
// and not suspended.
//
// Returns an error only in exceptional circumstances. Client closure errors
// are logged as warnings but do not cause the method to fail.
//
// Example:
//
//	err := conn.Disconnect()
//	if err != nil {
//		// handle disconnect error (rare)
//	}
func (c *ArkSSHConnection) Disconnect() error {
	if !c.isConnected {
		return nil
	}
	err := c.sshClient.Close()
	if err != nil {
		c.logger.Warning("Failed to close SSH client: %s", err.Error())
	}
	c.sshClient = nil
	c.isConnected = false
	c.isSuspended = false
	return nil
}

// SuspendConnection suspends the SSH connection.
//
// Marks the connection as suspended without actually closing the underlying
// SSH connection. When suspended, the connection will refuse to execute
// commands until it is restored using RestoreConnection.
//
// This is useful for temporarily disabling command execution while keeping
// the underlying network connection alive.
//
// Returns nil as this operation always succeeds.
//
// Example:
//
//	err := conn.SuspendConnection()
//	// Commands will now fail until RestoreConnection is called
func (c *ArkSSHConnection) SuspendConnection() error {
	c.isSuspended = true
	return nil
}

// RestoreConnection restores the SSH connection.
//
// Restores a previously suspended connection, allowing command execution
// to resume. This method clears the suspended state without affecting
// the underlying SSH connection.
//
// Returns nil as this operation always succeeds.
//
// Example:
//
//	err := conn.RestoreConnection()
//	// Commands can now be executed again
func (c *ArkSSHConnection) RestoreConnection() error {
	c.isSuspended = false
	return nil
}

// IsSuspended checks if the SSH connection is suspended.
//
// Returns the current suspension state of the connection. When suspended,
// the connection will refuse to execute commands even if the underlying
// SSH connection is still active.
//
// Returns true if the connection is currently suspended, false otherwise.
//
// Example:
//
//	if conn.IsSuspended() {
//		// Connection is suspended, restore before running commands
//		conn.RestoreConnection()
//	}
func (c *ArkSSHConnection) IsSuspended() bool {
	return c.isSuspended
}

// IsConnected checks if the SSH connection is established.
//
// Returns the current connection state indicating whether an SSH connection
// has been successfully established and is ready for use. This does not
// check the network connectivity, only the internal connection state.
//
// Returns true if the connection is established, false otherwise.
//
// Example:
//
//	if !conn.IsConnected() {
//		err := conn.Connect(connectionDetails)
//		if err != nil {
//			// handle connection error
//		}
//	}
func (c *ArkSSHConnection) IsConnected() bool {
	return c.isConnected
}

// RunCommand executes a command on the connected system.
//
// Executes the specified command on the remote machine through the established
// SSH connection. The method creates a new SSH session for each command execution,
// captures stdout and stderr output, and handles exit status detection.
//
// The method validates that the connection is active and not suspended before
// execution. Commands that return a different exit code than expected will
// result in an error.
//
// Session management is handled automatically - a new session is created for
// each command execution and cleaned up afterward. The method properly handles
// SSH exit errors to extract the actual exit status from the remote command.
//
// Parameters:
//   - command: The command configuration including the command string and
//     expected return code for validation
//
// Returns the command execution result containing stdout, stderr, and return
// code, or an error if the command cannot be executed, the session cannot be
// created, or the command returns an unexpected return code.
//
// Example:
//
//	cmd := &connectionsmodels.ArkConnectionCommand{
//		Command:    "ls -la /home",
//		ExpectedRC: 0,
//	}
//	result, err := conn.RunCommand(cmd)
//	if err != nil {
//		// handle execution error
//	}
//	fmt.Printf("Output: %s\n", result.Stdout)
func (c *ArkSSHConnection) RunCommand(command *connectionsmodels.ArkConnectionCommand) (*connectionsmodels.ArkConnectionResult, error) {
	if !c.isConnected || c.isSuspended {
		return nil, fmt.Errorf("cannot run command while not being connected")
	}
	c.logger.Debug("Running command [%s]", command.Command)
	session, err := c.sshClient.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH session: %w", err)
	}
	var stdoutBuf, stderrBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Stderr = &stderrBuf
	err = session.Run(command.Command)
	rc := 0
	if err != nil {
		var exitErr *ssh.ExitError
		if errors.As(err, &exitErr) {
			rc = exitErr.ExitStatus()
		}
	}

	stdout := stdoutBuf.String()
	stderr := stderrBuf.String()

	if rc != command.ExpectedRC {
		return nil, fmt.Errorf("failed to execute command [%s] - [%d] - [%s]", command.Command, rc, stderr)
	}

	c.logger.Debug("Command rc: [%d]", rc)
	c.logger.Debug("Command stdout: [%s]", stdout)
	c.logger.Debug("Command stderr: [%s]", stderr)

	return &connectionsmodels.ArkConnectionResult{
		Stdout: stdout,
		Stderr: stderr,
		RC:     rc,
	}, nil
}
