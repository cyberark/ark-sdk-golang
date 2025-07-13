package ssh

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/cyberark/ark-sdk-golang/pkg/common"
	"github.com/cyberark/ark-sdk-golang/pkg/common/connections"
	connectionsmodels "github.com/cyberark/ark-sdk-golang/pkg/models/common/connections"
	"golang.org/x/crypto/ssh"
	"os"
	"time"
)

const (
	// SSHPort is the default port for SSH connections.
	SSHPort = 22
)

const (
	connectionTimeout = 10 * time.Second
)

// ArkSSHConnection is a struct that implements the ArkConnection interface for SSH connections.
type ArkSSHConnection struct {
	connections.ArkConnection
	isConnected bool
	isSuspended bool
	sshClient   *ssh.Client
	logger      *common.ArkLogger
}

// NewArkSSHConnection creates a new instance of ArkSSHConnection.
func NewArkSSHConnection() *ArkSSHConnection {
	return &ArkSSHConnection{
		isConnected: false,
		isSuspended: false,
		logger:      common.GetLogger("ArkSSHConnection", common.Unknown),
	}
}

// Connect establishes an SSH connection using the provided connection details.
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
	c.logger.Debug(fmt.Sprintf("Connected to SSH server [%s] on port [%d]", connectionDetails.Address, connectionDetails.Port))
	c.sshClient = client
	c.isConnected = true
	c.isSuspended = false
	return nil
}

// Disconnect closes the SSH connection.
func (c *ArkSSHConnection) Disconnect() error {
	if !c.isConnected {
		return nil
	}
	err := c.sshClient.Close()
	if err != nil {
		c.logger.Warning(fmt.Sprintf("Failed to close SSH client: %s", err))
	}
	c.sshClient = nil
	c.isConnected = false
	c.isSuspended = false
	return nil
}

// SuspendConnection suspends the SSH connection.
func (c *ArkSSHConnection) SuspendConnection() error {
	c.isSuspended = true
	return nil
}

// RestoreConnection restores the SSH connection.
func (c *ArkSSHConnection) RestoreConnection() error {
	c.isSuspended = false
	return nil
}

// IsSuspended checks if the SSH connection is suspended.
func (c *ArkSSHConnection) IsSuspended() bool {
	return c.isSuspended
}

// IsConnected checks if the SSH connection is established.
func (c *ArkSSHConnection) IsConnected() bool {
	return c.isConnected
}

// RunCommand executes a command on the connected system.
func (c *ArkSSHConnection) RunCommand(command *connectionsmodels.ArkConnectionCommand) (*connectionsmodels.ArkConnectionResult, error) {
	if !c.isConnected || c.isSuspended {
		return nil, fmt.Errorf("cannot run command while not being connected")
	}
	c.logger.Debug(fmt.Sprintf("Running command [%s]", command.Command))
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

	c.logger.Debug(fmt.Sprintf("Command rc: [%d]", rc))
	c.logger.Debug(fmt.Sprintf("Command stdout: [%s]", stdout))
	c.logger.Debug(fmt.Sprintf("Command stderr: [%s]", stderr))

	return &connectionsmodels.ArkConnectionResult{
		Stdout: stdout,
		Stderr: stderr,
		RC:     rc,
	}, nil
}
