package serviceinfo

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/common"
	"github.com/cyberark/ark-sdk-golang/pkg/common/isp"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	serviceinfomodels "github.com/cyberark/ark-sdk-golang/pkg/services/sechub/serviceinfo/models"
	"github.com/mitchellh/mapstructure"
)

const (
	sechubURL = "/api/info"
)

// SecHubServiceInfoServiceConfig is the configuration for the Secrets Hub Service Info service.
var SecHubServiceInfoServiceConfig = services.ArkServiceConfig{
	ServiceName:                "sechub-serviceinfo",
	RequiredAuthenticatorNames: []string{"isp"},
	OptionalAuthenticatorNames: []string{},
}

// ArkSecHubServiceInfoService is the service for retrieve Secrets Hub service Info
type ArkSecHubServiceInfoService struct {
	services.ArkService
	*services.ArkBaseService
	ispAuth *auth.ArkISPAuth
	client  *isp.ArkISPServiceClient
}

// NewArkSecHubServiceInfoService creates a new instance of ArkSecHubServiceInfoService.
func NewArkSecHubServiceInfoService(authenticators ...auth.ArkAuth) (*ArkSecHubServiceInfoService, error) {
	serviceInfoService := &ArkSecHubServiceInfoService{}
	var serviceInfoServiceInterface services.ArkService = serviceInfoService
	baseService, err := services.NewArkBaseService(serviceInfoServiceInterface, authenticators...)
	if err != nil {
		return nil, err
	}
	ispBaseAuth, err := baseService.Authenticator("isp")
	if err != nil {
		return nil, err
	}
	ispAuth := ispBaseAuth.(*auth.ArkISPAuth)
	client, err := isp.FromISPAuth(ispAuth, "secretshub", ".", "", serviceInfoService.refreshSecHubAuth)
	if err != nil {
		return nil, err
	}
	serviceInfoService.client = client
	serviceInfoService.ispAuth = ispAuth
	serviceInfoService.ArkBaseService = baseService
	return serviceInfoService, nil
}

func (s *ArkSecHubServiceInfoService) refreshSecHubAuth(client *common.ArkClient) error {
	err := isp.RefreshClient(client, s.ispAuth)
	if err != nil {
		return err
	}
	return nil
}

// ServiceInfo retrieves the service info from the Secrets Hub service.
// https://api-docs.cyberark.com/docs/secretshub-api/b7c22j9aexv8r-service-info
func (s *ArkSecHubServiceInfoService) ServiceInfo() (*serviceinfomodels.ArkSecHubGetServiceInfo, error) {
	s.Logger.Info("Getting serviceinfo")
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
		return nil, fmt.Errorf("failed to get service info - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	serviceinfoJSON, err := common.DeserializeJSONSnake(response.Body)
	if err != nil {
		return nil, err
	}
	var serviceinfo serviceinfomodels.ArkSecHubGetServiceInfo
	err = mapstructure.Decode(serviceinfoJSON, &serviceinfo)
	if err != nil {
		return nil, err
	}
	return &serviceinfo, nil
}

// ServiceConfig returns the service configuration for the ArkSecHubServiceInfoService.
func (s *ArkSecHubServiceInfoService) ServiceConfig() services.ArkServiceConfig {
	return SecHubServiceInfoServiceConfig
}
