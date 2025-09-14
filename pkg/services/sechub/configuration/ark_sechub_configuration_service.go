package configuration

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/common"
	"github.com/cyberark/ark-sdk-golang/pkg/common/isp"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	configurationmodels "github.com/cyberark/ark-sdk-golang/pkg/services/sechub/configuration/models"
	"github.com/mitchellh/mapstructure"
)

const (
	sechubURL = "/api/configuration"
)

// ArkSecHubConfigurationService is the service for interacting with Secrets Hub configuration
type ArkSecHubConfigurationService struct {
	services.ArkService
	*services.ArkBaseService
	ispAuth *auth.ArkISPAuth
	client  *isp.ArkISPServiceClient
}

// NewArkSecHubConfigurationService creates a new instance of ArkSecHubConfigurationService.
func NewArkSecHubConfigurationService(authenticators ...auth.ArkAuth) (*ArkSecHubConfigurationService, error) {
	configurationService := &ArkSecHubConfigurationService{}
	var configurationServiceInterface services.ArkService = configurationService
	baseService, err := services.NewArkBaseService(configurationServiceInterface, authenticators...)
	if err != nil {
		return nil, err
	}
	ispBaseAuth, err := baseService.Authenticator("isp")
	if err != nil {
		return nil, err
	}
	ispAuth := ispBaseAuth.(*auth.ArkISPAuth)
	client, err := isp.FromISPAuth(ispAuth, "secretshub", ".", "", configurationService.refreshSecHubAuth)
	if err != nil {
		return nil, err
	}
	configurationService.client = client
	configurationService.ispAuth = ispAuth
	configurationService.ArkBaseService = baseService
	return configurationService, nil
}

func (s *ArkSecHubConfigurationService) refreshSecHubAuth(client *common.ArkClient) error {
	err := isp.RefreshClient(client, s.ispAuth)
	if err != nil {
		return err
	}
	return nil
}

// Configuration retrieves the configuration info from the Secrets Hub service.
// https://api-docs.cyberark.com/docs/secretshub-api/r3a0vv9er2enm-view-configuration
func (s *ArkSecHubConfigurationService) Configuration() (*configurationmodels.ArkSecHubGetConfiguration, error) {
	s.Logger.Info("Getting configuration")
	response, err := s.client.Get(context.Background(), sechubURL, nil)
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
		return nil, fmt.Errorf("failed to get configuration - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	configurationJSON, err := common.DeserializeJSONSnake(response.Body)
	if err != nil {
		return nil, err
	}
	var configurationinfo configurationmodels.ArkSecHubGetConfiguration
	err = mapstructure.Decode(configurationJSON, &configurationinfo)
	if err != nil {
		return nil, err
	}
	return &configurationinfo, nil
}

// SetConfiguration updates the configuration info in the Secrets Hub service.
// https://api-docs.cyberark.com/docs/secretshub-api/eko5hfu8sg16o-update-configuration
func (s *ArkSecHubConfigurationService) SetConfiguration(setConfiguration *configurationmodels.ArkSecHubSetConfiguration) (*configurationmodels.ArkSecHubGetConfiguration, error) {
	s.Logger.Info("Updating configuration. Setting secret validity to [%d]", setConfiguration.SyncSettings.SecretValidity)
	setConfigurationJSON, err := common.SerializeJSONCamel(setConfiguration)
	if err != nil {
		return nil, err
	}
	response, err := s.client.Patch(context.Background(), sechubURL, setConfigurationJSON)
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
		return nil, fmt.Errorf("failed to update configuration - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	configurationJSON, err := common.DeserializeJSONSnake(response.Body)
	if err != nil {
		return nil, err
	}
	var configurationinfo configurationmodels.ArkSecHubGetConfiguration
	err = mapstructure.Decode(configurationJSON, &configurationinfo)
	if err != nil {
		return nil, err
	}
	return &configurationinfo, nil
}

// ServiceConfig returns the service configuration for the ArkSecHubConfigurationService.
func (s *ArkSecHubConfigurationService) ServiceConfig() services.ArkServiceConfig {
	return ServiceConfig
}
