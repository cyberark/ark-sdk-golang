package filters

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/common"
	"github.com/cyberark/ark-sdk-golang/pkg/common/isp"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	filtersmodels "github.com/cyberark/ark-sdk-golang/pkg/services/sechub/filters/models"
	"github.com/mitchellh/mapstructure"
)

const (
	sechubURL = "/api/secret-stores/%s/filters"
	filterURL = "/api/secret-stores/%s/filters/%s"
)

// ArkSecHubFiltersPage is a page of ArkSecHubFilter items.
type ArkSecHubFiltersPage = common.ArkPage[filtersmodels.ArkSecHubFilter]

// SecHubFiltersServiceConfig is the configuration for the Secrets Hub filters service.
var SecHubFiltersServiceConfig = services.ArkServiceConfig{
	ServiceName:                "sechub-filters",
	RequiredAuthenticatorNames: []string{"isp"},
	OptionalAuthenticatorNames: []string{},
}

// ArkSecHubFiltersService is the service for interacting with Secrets Hub filters
type ArkSecHubFiltersService struct {
	services.ArkService
	*services.ArkBaseService
	ispAuth *auth.ArkISPAuth
	client  *isp.ArkISPServiceClient
}

// NewArkSecHubFiltersService creates a new instance of ArkSecHubFiltersService.
func NewArkSecHubFiltersService(authenticators ...auth.ArkAuth) (*ArkSecHubFiltersService, error) {
	filtersService := &ArkSecHubFiltersService{}
	var filtersServiceInterface services.ArkService = filtersService
	baseService, err := services.NewArkBaseService(filtersServiceInterface, authenticators...)
	if err != nil {
		return nil, err
	}
	ispBaseAuth, err := baseService.Authenticator("isp")
	if err != nil {
		return nil, err
	}
	ispAuth := ispBaseAuth.(*auth.ArkISPAuth)
	client, err := isp.FromISPAuth(ispAuth, "secretshub", ".", "", filtersService.refreshSecHubAuth)
	if err != nil {
		return nil, err
	}
	filtersService.client = client
	filtersService.ispAuth = ispAuth
	filtersService.ArkBaseService = baseService
	return filtersService, nil
}

func (s *ArkSecHubFiltersService) refreshSecHubAuth(client *common.ArkClient) error {
	err := isp.RefreshClient(client, s.ispAuth)
	if err != nil {
		return err
	}
	return nil
}

// Filter retrieves the filters info from the Secrets Hub service.
// https://api-docs.cyberark.com/docs/secretshub-api/rqykgubx980ul-get-secrets-filter
func (s *ArkSecHubFiltersService) Filter(getFilters *filtersmodels.ArkSecHubGetFilter) (*filtersmodels.ArkSecHubFilter, error) {
	if getFilters.StoreID == "" {
		s.Logger.Info("Setting Secret Store ID to default")
		getFilters.StoreID = "default"
	}
	if getFilters.FilterID == "" {
		s.Logger.Info("Setting Secret Store Filter ID to default")
		getFilters.FilterID = "default"
	}
	s.Logger.Info("Getting filter")
	response, err := s.client.Get(context.Background(), fmt.Sprintf(filterURL, getFilters.StoreID, getFilters.FilterID), nil)
	if err != nil {
		s.Logger.Error("Failed to list filters: %v", err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.GlobalLogger.Warning("Error closing response body")
		}
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		s.Logger.Error("Failed to list Secret Store Filters - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
		return nil, err
	}
	filterJSON, err := common.DeserializeJSONSnake(response.Body)
	if err != nil {
		s.Logger.Error("Failed to decode response: %v", err)
		return nil, err
	}
	var filter filtersmodels.ArkSecHubFilter
	err = mapstructure.Decode(filterJSON, &filter)
	if err != nil {
		return nil, err
	}
	return &filter, nil
}

// ListFilters retrieves the filters info from the Secrets Hub service.
// https://api-docs.cyberark.com/docs/secretshub-api/punr36gz4tuqe-get-all-secrets-filters
func (s *ArkSecHubFiltersService) ListFilters(getFilters *filtersmodels.ArkSecHubGetFilters) (<-chan *ArkSecHubFiltersPage, error) {
	if getFilters.StoreID == "" {
		s.Logger.Info("Setting Secret Store ID to default")
		getFilters.StoreID = "default"
	}
	s.Logger.Info("Getting filters")

	results := make(chan *ArkSecHubFiltersPage)
	go func() {
		defer close(results)
		response, err := s.client.Get(context.Background(), fmt.Sprintf(sechubURL, getFilters.StoreID), nil)
		if err != nil {
			s.Logger.Error("Failed to list filters: %v", err)
			return
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				common.GlobalLogger.Warning("Error closing response body")
			}
		}(response.Body)
		if response.StatusCode != http.StatusOK {
			s.Logger.Error("Failed to list Secret Store Filters - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
			return
		}
		result, err := common.DeserializeJSONSnake(response.Body)
		if err != nil {
			s.Logger.Error("Failed to decode response: %v", err)
			return
		}
		resultMap := result.(map[string]interface{})
		var filtersJSON []interface{}
		if filters, ok := resultMap["filters"]; ok {
			filtersJSON = filters.([]interface{})
		} else {
			s.Logger.Error("Failed to list Secret Store filters, unexpected result")
			return
		}
		for i, filtersMember := range filtersJSON {
			if filtersMemberMap, ok := filtersMember.(map[string]interface{}); ok {
				if ID, ok := filtersMemberMap["id"]; ok {
					filtersJSON[i].(map[string]interface{})["id"] = ID
				}
			}
		}
		var filters []*filtersmodels.ArkSecHubFilter
		if err := mapstructure.Decode(filtersJSON, &filters); err != nil {
			s.Logger.Error("Failed to validate Secret Store filters: %v", err)
			return
		}

		results <- &ArkSecHubFiltersPage{Items: filters}
	}()
	return results, nil
}

// AddFilter adds a new filter for a specific secret store id
// https://api-docs.cyberark.com/docs/secretshub-api/ifgbuo8tmt1en-create-secrets-filter
func (s *ArkSecHubFiltersService) AddFilter(filter *filtersmodels.ArkSecHubAddFilter) (*filtersmodels.ArkSecHubFilter, error) {
	s.Logger.Info("Adding filter for secret store [%s]", filter.StoreID)
	bodyMap := map[string]interface{}{
		"type": filter.Type,
		"data": map[string]string{
			"safeName": filter.Data.SafeName,
		},
	}
	response, err := s.client.Post(context.Background(), fmt.Sprintf(sechubURL, filter.StoreID), bodyMap)
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
		return nil, fmt.Errorf("failed to create filter - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	filterJSON, err := common.DeserializeJSONSnake(response.Body)
	if err != nil {
		return nil, err
	}
	var filterResponse filtersmodels.ArkSecHubFilter
	err = mapstructure.Decode(filterJSON, &filterResponse)
	if err != nil {
		return nil, err
	}
	return &filterResponse, nil
}

// DeleteFilter deletes a specified filter based on secret store id and filter id
// https://api-docs.cyberark.com/docs/secretshub-api/h8q9q5xtkxqgz-delete-secrets-filter
func (s *ArkSecHubFiltersService) DeleteFilter(filter *filtersmodels.ArkSecHubDeleteFilter) error {
	s.Logger.Info("Deleting secret store [%s] filter [%s]", filter.StoreID, filter.FilterID)
	response, err := s.client.Delete(context.Background(), fmt.Sprintf(filterURL, filter.StoreID, filter.FilterID), nil)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.GlobalLogger.Warning("Error closing response body")
		}
	}(response.Body)
	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete filter - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	return nil
}

// ServiceConfig returns the service configuration for the ArkSecHubFiltersService.
func (s *ArkSecHubFiltersService) ServiceConfig() services.ArkServiceConfig {
	return SecHubFiltersServiceConfig
}
