package access

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/common"
	"github.com/cyberark/ark-sdk-golang/pkg/common/connections"
	"github.com/cyberark/ark-sdk-golang/pkg/common/connections/ssh"
	"github.com/cyberark/ark-sdk-golang/pkg/common/connections/winrm"
	"github.com/cyberark/ark-sdk-golang/pkg/common/isp"
	commonmodels "github.com/cyberark/ark-sdk-golang/pkg/models/common"
	connectionsmodels "github.com/cyberark/ark-sdk-golang/pkg/models/common/connections"
	"github.com/cyberark/ark-sdk-golang/pkg/models/common/connections/connectiondata"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	accessmodels "github.com/cyberark/ark-sdk-golang/pkg/services/sia/access/models"
	"github.com/mitchellh/mapstructure"
)

const (
	connectorURL                 = "/api/connectors/%s"
	connectorSetupScriptURL      = "/api/connectors/setup-script"
	connectorTestReachabilityURL = "/api/connectors/%s/reachability"

	// Linux / Darwin Commands
	unixStopConnectorServiceCmd   = "sudo systemctl stop cyberark-dpa-connector"
	unixRemoveConnectorServiceCmd = "sudo rm -f /etc/systemd/system/cyberark-dpa-connector.service && sudo rm -f /usr/lib/systemd/system/cyberark-dpa-connector.service && sudo systemctl daemon-reload && sudo systemctl reset-failed"
	unixRemoveConnectorFilesCmd   = "sudo rm -rf /opt/cyberark/connector"
	unixConnectorActiveCmd        = "sudo systemctl is-active --quiet cyberark-dpa-connector"
	unixReadConnectorConfigCmd    = "sudo cat /opt/cyberark/connector/connector.config.json"

	// Windows Commands
	winStopConnectorServiceCmd   = "Stop-Service -Name \"CyberArkDPAConnector\""
	winRemoveConnectorServiceCmd = `$service = Get-WmiObject -Class Win32_Service -Filter "Name='CyberArkDPAConnector'"; $service.delete()`
	winRemoveConnectorFilesCmd   = "Remove-Item -LiteralPath \"C:\\Program Files\\CyberArk\\DPAConnector\" -Force -Recurse"
	winConnectorActiveCmd        = `$result = Get-Service -Name "CyberArkDPAConnector"; if ($result.Status -ne 'Running') { return 1 }`
	winReadConnectorConfigCmd    = "Get-Content -Path \"C:\\Program Files\\CyberArk\\DPAConnector\\connector.config.json\""

	// Retry Constants
	connectorReadyRetryCount = 5
	connectorRetryTick       = 3.0 * time.Second
)

// ConnectorCmdSet maps OS types to their respective command sets.
var connectorCmdSet = map[string]map[string]string{
	commonmodels.OSTypeLinux: {
		"stopConnectorService":   unixStopConnectorServiceCmd,
		"removeConnectorService": unixRemoveConnectorServiceCmd,
		"removeConnectorFiles":   unixRemoveConnectorFilesCmd,
		"connectorActive":        unixConnectorActiveCmd,
		"readConnectorConfig":    unixReadConnectorConfigCmd,
	},
	commonmodels.OSTypeDarwin: {
		"stopConnectorService":   unixStopConnectorServiceCmd,
		"removeConnectorService": unixRemoveConnectorServiceCmd,
		"removeConnectorFiles":   unixRemoveConnectorFilesCmd,
		"connectorActive":        unixConnectorActiveCmd,
		"readConnectorConfig":    unixReadConnectorConfigCmd,
	},
	commonmodels.OSTypeWindows: {
		"stopConnectorService":   winStopConnectorServiceCmd,
		"removeConnectorService": winRemoveConnectorServiceCmd,
		"removeConnectorFiles":   winRemoveConnectorFilesCmd,
		"connectorActive":        winConnectorActiveCmd,
		"readConnectorConfig":    winReadConnectorConfigCmd,
	},
}

// ArkSIAAccessService is a struct that implements the ArkService interface and provides functionality for Connectors of SIA.
type ArkSIAAccessService struct {
	services.ArkService
	*services.ArkBaseService
	ispAuth *auth.ArkISPAuth
	client  *isp.ArkISPServiceClient
}

// NewArkSIAAccessService creates a new instance of ArkSIAAccessService with the provided authenticators.
func NewArkSIAAccessService(authenticators ...auth.ArkAuth) (*ArkSIAAccessService, error) {
	accessService := &ArkSIAAccessService{}
	var accessServiceInterface services.ArkService = accessService
	baseService, err := services.NewArkBaseService(accessServiceInterface, authenticators...)
	if err != nil {
		return nil, err
	}
	ispBaseAuth, err := baseService.Authenticator("isp")
	if err != nil {
		return nil, err
	}
	ispAuth := ispBaseAuth.(*auth.ArkISPAuth)
	client, err := isp.FromISPAuth(ispAuth, "dpa", ".", "", accessService.refreshSIAAuth)
	if err != nil {
		return nil, err
	}
	accessService.client = client
	accessService.ispAuth = ispAuth
	accessService.ArkBaseService = baseService
	return accessService, nil
}

func (s *ArkSIAAccessService) refreshSIAAuth(client *common.ArkClient) error {
	err := isp.RefreshClient(client, s.ispAuth)
	if err != nil {
		return err
	}
	return nil
}

func (s *ArkSIAAccessService) createConnection(
	osType string,
	targetMachine string,
	username string,
	password string,
	privateKeyPath string,
	privateKeyContents string,
	retryCount int,
	retryDelay int,
) (connections.ArkConnection, map[string]string, error) {
	var connection connections.ArkConnection
	var connectionDetails *connectionsmodels.ArkConnectionDetails

	if osType == commonmodels.OSTypeWindows {
		connection = winrm.NewArkWinRMConnection()
		connectionDetails = &connectionsmodels.ArkConnectionDetails{
			Address:        targetMachine,
			Port:           winrm.WinRMHTTPSPort,
			ConnectionType: connectionsmodels.WinRM,
			Credentials: &connectionsmodels.ArkConnectionCredentials{
				User:     username,
				Password: password,
			},
			ConnectionData: &connectiondata.ArkWinRMConnectionData{
				CertificatePath:  "",
				TrustCertificate: true,
			},
		}
	} else {
		connection = ssh.NewArkSSHConnection()
		connectionDetails = &connectionsmodels.ArkConnectionDetails{
			Address:        targetMachine,
			Port:           ssh.SSHPort,
			ConnectionType: connectionsmodels.SSH,
			Credentials: &connectionsmodels.ArkConnectionCredentials{
				User:               username,
				Password:           password,
				PrivateKeyFilepath: privateKeyPath,
				PrivateKeyContents: privateKeyContents,
			},
			ConnectionData:    &connectiondata.ArkSSHConnectionData{},
			ConnectionRetries: retryCount,
			RetryTickPeriod:   retryDelay,
		}
	}

	if err := connection.Connect(connectionDetails); err != nil {
		return nil, nil, fmt.Errorf("failed to connect: %w", err)
	}

	return connection, connectorCmdSet[osType], nil
}

func (s *ArkSIAAccessService) installConnectorOnMachine(
	installScript string,
	osType string,
	targetMachine string,
	username string,
	password string,
	privateKeyPath string,
	privateKeyContents string,
	retryCount int,
	retryDelay int,
) (*accessmodels.ArkSIAAccessConnectorID, error) {
	// Create connection
	connection, cmdSet, err := s.createConnection(
		osType,
		targetMachine,
		username,
		password,
		privateKeyPath,
		privateKeyContents,
		retryCount,
		retryDelay,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection: %w", err)
	}
	defer func(connection connections.ArkConnection) {
		err := connection.Disconnect()
		if err != nil {
			s.Logger.Warning("failed to disconnect: %v", err)
		}
	}(connection)

	// Run commands to stop, remove service, and remove files
	_, _ = connection.RunCommand(&connectionsmodels.ArkConnectionCommand{
		Command: cmdSet["stopConnectorService"],
	})
	_, _ = connection.RunCommand(&connectionsmodels.ArkConnectionCommand{
		Command: cmdSet["removeConnectorService"],
	})
	_, _ = connection.RunCommand(&connectionsmodels.ArkConnectionCommand{
		Command: cmdSet["removeConnectorFiles"],
	})

	// Install the connector
	if osType == commonmodels.OSTypeWindows {
		_, err = connection.RunCommand(&connectionsmodels.ArkConnectionCommand{
			Command:          installScript,
			ExtraCommandData: map[string]interface{}{"force_command_split": true},
		})
	} else {
		_, err = connection.RunCommand(&connectionsmodels.ArkConnectionCommand{
			Command: installScript,
		})
	}
	if err != nil {
		return nil, fmt.Errorf("failed to install connector: %w", err)
	}

	// Retry checking if the connector is active
	currConnReadyRetryCount := connectorReadyRetryCount
	for {
		_, err = connection.RunCommand(&connectionsmodels.ArkConnectionCommand{
			Command: cmdSet["connectorActive"],
		})
		if err == nil {
			break
		}
		if currConnReadyRetryCount > 0 {
			currConnReadyRetryCount--
			time.Sleep(connectorRetryTick)
			continue
		}
		return nil, fmt.Errorf("failed to check if connector is active: %w", err)
	}

	// Read the connector configuration
	result, err := connection.RunCommand(&connectionsmodels.ArkConnectionCommand{
		Command: cmdSet["readConnectorConfig"],
	})
	if err != nil {
		return nil, fmt.Errorf("failed to read connector config: %w", err)
	}

	// Parse the connector configuration and return the ID
	var connectorConfig map[string]interface{}
	if err := json.Unmarshal([]byte(result.Stdout), &connectorConfig); err != nil {
		return nil, fmt.Errorf("failed to parse connector config: %w", err)
	}
	connectorID, ok := connectorConfig["Id"].(string)
	if !ok {
		return nil, fmt.Errorf("connector ID not found in config")
	}
	return &accessmodels.ArkSIAAccessConnectorID{ConnectorID: connectorID}, nil
}

func (s *ArkSIAAccessService) uninstallConnectorOnMachine(
	osType string,
	targetMachine string,
	username string,
	password string,
	privateKeyPath string,
	privateKeyContents string,
	retryCount int,
	retryDelay int,
) error {
	// Create connection
	connection, cmdSet, err := s.createConnection(
		osType,
		targetMachine,
		username,
		password,
		privateKeyPath,
		privateKeyContents,
		retryCount,
		retryDelay,
	)
	if err != nil {
		return fmt.Errorf("failed to create connection: %w", err)
	}
	defer func(connection connections.ArkConnection) {
		err := connection.Disconnect()
		if err != nil {
			s.Logger.Warning("failed to disconnect: %v", err)
		}
	}(connection)

	// Run commands to stop, remove service, and remove files
	_, err = connection.RunCommand(&connectionsmodels.ArkConnectionCommand{
		Command: cmdSet["stopConnectorService"],
	})
	if err != nil {
		return fmt.Errorf("failed to stop connector service: %w", err)
	}

	_, err = connection.RunCommand(&connectionsmodels.ArkConnectionCommand{
		Command: cmdSet["removeConnectorService"],
	})
	if err != nil {
		return fmt.Errorf("failed to remove connector service: %w", err)
	}

	_, err = connection.RunCommand(&connectionsmodels.ArkConnectionCommand{
		Command: cmdSet["removeConnectorFiles"],
	})
	if err != nil {
		return fmt.Errorf("failed to remove connector files: %w", err)
	}

	return nil
}

// TestConnectorReachability tests the reachability of a connector.
func (s *ArkSIAAccessService) TestConnectorReachability(testReachabilityRequest *accessmodels.ArkSIATestConnectorReachability) (*accessmodels.ArkSIAReachabilityTestResponse, error) {
	s.Logger.Info("Starting connector reachability test. ConnectorID: %s", testReachabilityRequest.ConnectorID)
	var testReachabilityRequestJSON = map[string]interface{}{
		"targets": []map[string]interface{}{
			{
				"hostname": testReachabilityRequest.TargetHostname,
				"port":     testReachabilityRequest.TargetPort,
			},
		},
		"checkBackendEndpoints": testReachabilityRequest.CheckBackendEndpoints,
	}
	response, err := s.client.Post(context.Background(), fmt.Sprintf(connectorTestReachabilityURL, testReachabilityRequest.ConnectorID), testReachabilityRequestJSON)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.GlobalLogger.Warning("Error closing response body")
		}
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to test connector reachability - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	reachabilityTestResponseJSON, err := common.DeserializeJSONSnake(response.Body)
	if err != nil {
		return nil, err
	}
	var testResponse accessmodels.ArkSIAReachabilityTestResponse
	err = mapstructure.Decode(reachabilityTestResponseJSON, &testResponse)
	if err != nil {
		return nil, err
	}
	return &testResponse, nil
}

// ConnectorSetupScript creates the setup script for the connector.
func (s *ArkSIAAccessService) ConnectorSetupScript(getConnectorSetupScript *accessmodels.ArkSIAGetConnectorSetupScript) (*accessmodels.ArkSIAConnectorSetupScript, error) {
	s.Logger.Info("Retrieving new connector setup script")
	var getConnectorSetupScriptJSON map[string]interface{}
	err := mapstructure.Decode(getConnectorSetupScript, &getConnectorSetupScriptJSON)
	if err != nil {
		return nil, err
	}
	response, err := s.client.Post(context.Background(), connectorSetupScriptURL, getConnectorSetupScriptJSON)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.GlobalLogger.Warning("Error closing response body")
		}
	}(response.Body)
	if response.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to retrieve connector setup script - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	connectorSetupScriptJSON, err := common.DeserializeJSONSnake(response.Body)
	if err != nil {
		return nil, err
	}
	var setupScript accessmodels.ArkSIAConnectorSetupScript
	err = mapstructure.Decode(connectorSetupScriptJSON, &setupScript)
	if err != nil {
		return nil, err
	}
	return &setupScript, nil
}

// InstallConnector installs the connector on the target machine.
func (s *ArkSIAAccessService) InstallConnector(installConnector *accessmodels.ArkSIAInstallConnector) (*accessmodels.ArkSIAAccessConnectorID, error) {
	s.Logger.Info(
		"Installing connector on machine [%s] of type [%s]",
		installConnector.TargetMachine,
		installConnector.ConnectorOS,
	)
	installationScript, err := s.ConnectorSetupScript(&accessmodels.ArkSIAGetConnectorSetupScript{
		ConnectorOS:     installConnector.ConnectorOS,
		ConnectorPoolID: installConnector.ConnectorPoolID,
		ConnectorType:   installConnector.ConnectorType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve connector setup script: %w", err)
	}
	return s.installConnectorOnMachine(
		installationScript.BashCmd,
		installConnector.ConnectorOS,
		installConnector.TargetMachine,
		installConnector.Username,
		installConnector.Password,
		strings.TrimSuffix(common.ExpandFolder(installConnector.PrivateKeyPath), "/"),
		installConnector.PrivateKeyContents,
		installConnector.RetryCount,
		installConnector.RetryDelay,
	)
}

// UninstallConnector uninstalls the connector from the target machine.
func (s *ArkSIAAccessService) UninstallConnector(uninstallConnector *accessmodels.ArkSIAUninstallConnector) error {
	s.Logger.Info(
		"Uninstalling connector [%s] from machine",
		uninstallConnector.ConnectorID,
	)
	err := s.uninstallConnectorOnMachine(
		uninstallConnector.ConnectorOS,
		uninstallConnector.TargetMachine,
		uninstallConnector.Username,
		uninstallConnector.Password,
		strings.TrimSuffix(common.ExpandFolder(uninstallConnector.PrivateKeyPath), "/"),
		uninstallConnector.PrivateKeyContents,
		uninstallConnector.RetryCount,
		uninstallConnector.RetryDelay,
	)
	if err != nil {
		return err
	}
	return s.DeleteConnector(&accessmodels.ArkSIADeleteConnector{
		ConnectorID: uninstallConnector.ConnectorID,
		RetryCount:  uninstallConnector.RetryCount,
		RetryDelay:  uninstallConnector.RetryDelay,
	})
}

// DeleteConnector deletes the connector from the target machine.
func (s *ArkSIAAccessService) DeleteConnector(deleteConnector *accessmodels.ArkSIADeleteConnector) error {
	s.Logger.Info(
		"Deleting connector [%s] from machine",
		deleteConnector.ConnectorID,
	)
	currentTryCount := 0
	for {
		response, err := s.client.Delete(context.Background(), fmt.Sprintf(connectorURL, deleteConnector.ConnectorID), nil)
		if err != nil {
			return err
		}
		if response.StatusCode != http.StatusOK {
			if currentTryCount < deleteConnector.RetryCount {
				currentTryCount++
				s.Logger.Warning("Failed to delete connector, retrying... [%d/%d]", currentTryCount, deleteConnector.RetryCount)
				time.Sleep(time.Duration(deleteConnector.RetryDelay) * time.Second)
				continue
			}
			return fmt.Errorf("failed to delete connector - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
		}
		break
	}
	return nil
}

// ServiceConfig returns the service configuration for the ArkSIAAccessService.
func (s *ArkSIAAccessService) ServiceConfig() services.ArkServiceConfig {
	return ServiceConfig
}
