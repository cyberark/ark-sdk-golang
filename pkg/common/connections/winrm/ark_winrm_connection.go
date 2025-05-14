package winrm

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/cyberark/ark-sdk-golang/pkg/common"
	"github.com/cyberark/ark-sdk-golang/pkg/common/connections"
	connectionsmodels "github.com/cyberark/ark-sdk-golang/pkg/models/common/connections"
	"github.com/cyberark/ark-sdk-golang/pkg/models/common/connections/connectiondata"
	"github.com/google/uuid"
	"github.com/masterzen/winrm"
	"golang.org/x/text/encoding/unicode"
	"io"
	"os"
	"time"
)

const (
	// WinRMHTTPSPort is the default port for WinRM HTTPS connections.
	WinRMHTTPSPort = 5986
)

const (
	winrmConnectionTimeout = 10 * time.Second
	maxSingleCommandSize   = 2000
	maxChunkSize           = 4000
)

// ArkWinRMConnection is a struct that implements the ArkConnection interface for WinRM connections.
type ArkWinRMConnection struct {
	connections.ArkConnection
	isConnected bool
	isSuspended bool
	winrmClient *winrm.Client
	winrmShell  *winrm.Shell
	logger      *common.ArkLogger
}

// NewArkWinRMConnection creates a new instance of ArkWinRMConnection.
func NewArkWinRMConnection() *ArkWinRMConnection {
	return &ArkWinRMConnection{
		isConnected: false,
		isSuspended: false,
		logger:      common.GetLogger("ArkWinRMConnection", common.Unknown),
	}
}

// Connect establishes a WinRM connection using the provided connection details.
func (c *ArkWinRMConnection) Connect(connectionDetails *connectionsmodels.ArkConnectionDetails) error {
	if c.isConnected {
		return nil
	}

	targetPort := WinRMHTTPSPort
	if connectionDetails.Port != 0 {
		targetPort = connectionDetails.Port
	}
	if connectionDetails.ConnectionRetries == 0 {
		connectionDetails.ConnectionRetries = 1
	}
	var err error
	var certData []byte
	certPath := ""
	trustCert := false
	if winrmData, ok := connectionDetails.ConnectionData.(*connectiondata.ArkWinRMConnectionData); ok {
		certPath = winrmData.CertificatePath
		trustCert = winrmData.TrustCertificate
	}
	if certPath != "" {
		certData, err = os.ReadFile(certPath)
		if err != nil {
			return fmt.Errorf("failed to read certificate file: %w", err)
		}
	}

	c.logger.Debug(fmt.Sprintf("Connecting to WinRM server [%s] on port [%d]", connectionDetails.Address, targetPort))
	endpoint := winrm.NewEndpoint(
		connectionDetails.Address,
		targetPort,
		true,
		trustCert,
		certData,
		nil,
		nil,
		winrmConnectionTimeout,
	)
	if connectionDetails.Credentials == nil {
		return fmt.Errorf("missing credentials for WinRM connection")
	}

	var client *winrm.Client
	var shell *winrm.Shell
	for i := 0; i < connectionDetails.ConnectionRetries; i++ {
		params := winrm.DefaultParameters
		params.TransportDecorator = func() winrm.Transporter { return &winrm.ClientNTLM{} }
		client, err = winrm.NewClientWithParameters(endpoint, connectionDetails.Credentials.User, connectionDetails.Credentials.Password, params)
		if err != nil {
			if common.IsConnectionRefused(err) && i < connectionDetails.ConnectionRetries-1 {
				c.logger.Info(fmt.Sprintf("Failed to create WinRM client: %s - Retrying...", err))
				time.Sleep(time.Duration(connectionDetails.RetryTickPeriod) * time.Second)
				continue
			}
			return fmt.Errorf("failed to create WinRM client: %w", err)
		}
		shell, err = client.CreateShell()
		if err != nil {
			if common.IsConnectionRefused(err) && i < connectionDetails.ConnectionRetries-1 {
				c.logger.Info(fmt.Sprintf("Failed to create WinRM shell: %s - Retrying...", err))
				time.Sleep(time.Duration(connectionDetails.RetryTickPeriod) * time.Second)
				continue
			}
			return fmt.Errorf("failed to create WinRM shell: %w", err)
		}
		break
	}
	c.logger.Debug(fmt.Sprintf("WinRM client and shell created successfully for [%s]", connectionDetails.Address))
	c.winrmClient = client
	c.winrmShell = shell
	c.isConnected = true
	c.isSuspended = false
	return nil
}

// Disconnect closes the WinRM connection.
func (c *ArkWinRMConnection) Disconnect() error {
	if !c.isConnected {
		return nil
	}
	err := c.winrmShell.Close()
	if err != nil {
		c.logger.Warning(fmt.Sprintf("failed to close WinRM shell: %s", err))
	}
	c.winrmShell = nil
	c.winrmClient = nil
	c.isConnected = false
	c.isSuspended = false
	return nil
}

// SuspendConnection suspends the WinRM connection.
func (c *ArkWinRMConnection) SuspendConnection() error {
	c.isSuspended = true
	return nil
}

// RestoreConnection restores the WinRM connection.
func (c *ArkWinRMConnection) RestoreConnection() error {
	c.isSuspended = false
	return nil
}

// IsSuspended checks if the WinRM connection is suspended.
func (c *ArkWinRMConnection) IsSuspended() bool {
	return c.isSuspended
}

// IsConnected checks if the WinRM connection is established.
func (c *ArkWinRMConnection) IsConnected() bool {
	return c.isConnected
}

// RunCommand executes a command on the remote machine using WinRM.
// It handles command splitting if the command exceeds the maximum size.
// It also manages the creation of a temporary file for large commands.
// The command is executed in PowerShell and the output is returned.
func (c *ArkWinRMConnection) RunCommand(command *connectionsmodels.ArkConnectionCommand) (*connectionsmodels.ArkConnectionResult, error) {
	if !c.isConnected || c.isSuspended {
		return nil, fmt.Errorf("cannot run command while not being connected")
	}

	c.logger.Debug(fmt.Sprintf("Running command [%s]", command.Command))

	if len(command.Command) > maxSingleCommandSize || (command.ExtraCommandData != nil && command.ExtraCommandData["force_command_split"] == true) {
		encoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder()
		utf16EncodedString, err := encoder.String(command.Command)
		if err != nil {
			return nil, fmt.Errorf("failed to encode string: %w", err)
		}
		encodedCommand := []byte(utf16EncodedString)
		maxSize := maxChunkSize
		var chunks [][]byte
		for i := 0; i < len(encodedCommand); i += maxSize {
			end := i + maxSize
			if end > len(encodedCommand) {
				end = len(encodedCommand)
			}
			chunks = append(chunks, encodedCommand[i:end])
		}

		commandUniqueFileName := uuid.New().String()
		commandFile := fmt.Sprintf("C:\\temp\\%s.ps1", commandUniqueFileName)

		// Ensure C:\temp exists
		_, err = c.winrmShell.ExecuteWithContext(context.Background(), "if not exist C:\\temp mkdir C:\\temp")
		if err != nil {
			return nil, fmt.Errorf("failed to create temp directory: %w", err)
		}

		// Write chunks to the file
		for _, chunk := range chunks {
			base64utf16EncodedString := base64.StdEncoding.EncodeToString(chunk)
			if err != nil {
				return nil, fmt.Errorf("failed to encode string: %w", err)
			}
			_, err = c.winrmShell.ExecuteWithContext(context.Background(), fmt.Sprintf(
				`powershell -Command "[System.Text.Encoding]::Unicode.GetString([System.Convert]::FromBase64String('%s')) | Add-Content -Path %s -Encoding Unicode -NoNewline"`,
				base64utf16EncodedString, commandFile))
			if err != nil {
				return nil, fmt.Errorf("failed to write chunk to file: %w", err)
			}
		}

		// Execute the PowerShell script
		commandOutput, err := c.winrmShell.ExecuteWithContext(context.Background(), fmt.Sprintf("powershell -File %s", commandFile))
		if err != nil {
			return nil, fmt.Errorf("failed to execute command: %w", err)
		}

		// Clean up the temporary file
		_, _ = c.winrmShell.ExecuteWithContext(context.Background(), fmt.Sprintf("del /f %s", commandFile))

		stdout := &bytes.Buffer{}
		stderr := &bytes.Buffer{}
		go io.Copy(stdout, commandOutput.Stdout)
		go io.Copy(stderr, commandOutput.Stderr)
		commandOutput.Wait()
		if command.ExpectedRC != commandOutput.ExitCode() {
			return nil, fmt.Errorf("failed to execute command [%s] - [%d] - [%s]", command.Command, commandOutput.ExitCode(), stderr.String())
		}

		c.logger.Debug(fmt.Sprintf("Command rc: [%d]", commandOutput.ExitCode()))
		c.logger.Debug(fmt.Sprintf("Command stdout: [%s]", stdout))
		c.logger.Debug(fmt.Sprintf("Command stderr: [%s]", stderr))

		return &connectionsmodels.ArkConnectionResult{
			Stdout: stdout.String(),
			Stderr: stderr.String(),
			RC:     commandOutput.ExitCode(),
		}, nil
	}

	encoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder()
	utf16EncodedString, err := encoder.String(command.Command)
	if err != nil {
		return nil, fmt.Errorf("failed to encode string: %w", err)
	}
	base64utf16EncodedString := base64.StdEncoding.EncodeToString([]byte(utf16EncodedString))
	encodedCommand := fmt.Sprintf("powershell -encodedcommand \"%s\"", base64utf16EncodedString)
	commandOutput, err := c.winrmShell.ExecuteWithContext(context.Background(), encodedCommand)
	if err != nil {
		return nil, fmt.Errorf("failed to execute command: %w", err)
	}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	go io.Copy(stdout, commandOutput.Stdout)
	go io.Copy(stderr, commandOutput.Stderr)
	commandOutput.Wait()
	if command.ExpectedRC != commandOutput.ExitCode() {
		return nil, fmt.Errorf("failed to execute command [%s] - [%d] - [%s]", command.Command, commandOutput.ExitCode(), stderr.String())
	}

	c.logger.Debug(fmt.Sprintf("Command rc: [%d]", commandOutput.ExitCode()))
	c.logger.Debug(fmt.Sprintf("Command stdout: [%s]", stdout))
	c.logger.Debug(fmt.Sprintf("Command stderr: [%s]", stderr))

	return &connectionsmodels.ArkConnectionResult{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
		RC:     commandOutput.ExitCode(),
	}, nil
}
