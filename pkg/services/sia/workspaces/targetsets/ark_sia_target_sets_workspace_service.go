package targetsets

import (
	"context"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/common"
	"github.com/cyberark/ark-sdk-golang/pkg/common/isp"
	targetsetsmodels "github.com/cyberark/ark-sdk-golang/pkg/models/services/sia/workspaces/targetsets"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	"io"
	"net/http"
	"regexp"
)

const (
	targetSetsURL     = "/api/targetsets"
	bulkTargetSetsURL = "/api/targetsets/bulk"
	targetSetURL      = "/api/targetsets/%s"
)

// SIATargetSetsWorkspaceServiceConfig is the configuration for the SIA target sets workspace service.
var SIATargetSetsWorkspaceServiceConfig = services.ArkServiceConfig{
	ServiceName:                "sia-workspaces-target-sets",
	RequiredAuthenticatorNames: []string{"isp"},
	OptionalAuthenticatorNames: []string{},
}

// ArkSIATargetSetsWorkspaceService is the service for managing target sets in a workspace.
type ArkSIATargetSetsWorkspaceService struct {
	services.ArkService
	*services.ArkBaseService
	ispAuth *auth.ArkISPAuth
	client  *isp.ArkISPServiceClient
}

// NewArkSIATargetSetsWorkspaceService creates a new instance of ArkSIATargetSetsWorkspaceService.
func NewArkSIATargetSetsWorkspaceService(authenticators ...auth.ArkAuth) (*ArkSIATargetSetsWorkspaceService, error) {
	targetSetsService := &ArkSIATargetSetsWorkspaceService{}
	var targetSetsServiceInterface services.ArkService = targetSetsService
	baseService, err := services.NewArkBaseService(targetSetsServiceInterface, authenticators...)
	if err != nil {
		return nil, err
	}
	ispBaseAuth, err := baseService.Authenticator("isp")
	if err != nil {
		return nil, err
	}
	ispAuth := ispBaseAuth.(*auth.ArkISPAuth)
	client, err := isp.FromISPAuth(ispAuth, "dpa", ".", "", targetSetsService.refreshSIAAuth)
	if err != nil {
		return nil, err
	}
	targetSetsService.client = client
	targetSetsService.ispAuth = ispAuth
	targetSetsService.ArkBaseService = baseService
	return targetSetsService, nil
}

func (s *ArkSIATargetSetsWorkspaceService) refreshSIAAuth(client *common.ArkClient) error {
	err := isp.RefreshClient(client, s.ispAuth)
	if err != nil {
		return err
	}
	return nil
}

// AddTargetSet adds a new target set with related strong account.
func (s *ArkSIATargetSetsWorkspaceService) AddTargetSet(addTargetSet *targetsetsmodels.ArkSIAAddTargetSet) (*targetsetsmodels.ArkSIATargetSet, error) {
	s.Logger.Info("Adding target set [%s]", addTargetSet.Name)
	var addTargetSetJSON map[string]interface{}
	err := mapstructure.Decode(addTargetSet, &addTargetSetJSON)
	if err != nil {
		return nil, err
	}
	response, err := s.client.Post(context.Background(), targetSetsURL, addTargetSetJSON)
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
		return nil, fmt.Errorf("failed to add target set - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	targetSetJSON, err := common.DeserializeJSONSnake(response.Body)
	if err != nil {
		return nil, err
	}
	var targetSet targetsetsmodels.ArkSIATargetSet
	err = mapstructure.Decode(targetSetJSON["target_set"], &targetSet)
	if err != nil {
		return nil, err
	}
	return &targetSet, nil
}

// BulkAddTargetSets adds multiple target sets with related strong account.
func (s *ArkSIATargetSetsWorkspaceService) BulkAddTargetSets(bulkAddTargetSets *targetsetsmodels.ArkSIABulkAddTargetSets) (*targetsetsmodels.ArkSIABulkTargetSetResponse, error) {
	s.Logger.Info("Bulk adding target set [%v]", bulkAddTargetSets)
	var bulkAddTargetSetsJSON map[string]interface{}
	err := mapstructure.Decode(bulkAddTargetSets, &bulkAddTargetSetsJSON)
	if err != nil {
		return nil, err
	}
	response, err := s.client.Post(context.Background(), bulkTargetSetsURL, bulkAddTargetSetsJSON)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.GlobalLogger.Warning("Error closing response body")
		}
	}(response.Body)
	if response.StatusCode != http.StatusMultiStatus {
		return nil, fmt.Errorf("failed to bulk add target set - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	bulkTargetSetRespJSON, err := common.DeserializeJSONSnake(response.Body)
	if err != nil {
		return nil, err
	}
	var bulkTargetSetsResp targetsetsmodels.ArkSIABulkTargetSetResponse
	err = mapstructure.Decode(bulkTargetSetRespJSON, &bulkTargetSetsResp)
	if err != nil {
		return nil, err
	}
	return &bulkTargetSetsResp, nil
}

// DeleteTargetSet deletes a target set.
func (s *ArkSIATargetSetsWorkspaceService) DeleteTargetSet(deleteTargetSet *targetsetsmodels.ArkSIADeleteTargetSet) error {
	s.Logger.Info("Deleting target set [%s]", deleteTargetSet.Name)
	response, err := s.client.Delete(context.Background(), fmt.Sprintf(targetSetURL, deleteTargetSet.Name), nil)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete target set - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	return nil
}

// BulkDeleteTargetSets deletes multiple target sets.
func (s *ArkSIATargetSetsWorkspaceService) BulkDeleteTargetSets(bulkDeleteTargetSets *targetsetsmodels.ArkSIABulkDeleteTargetSets) (*targetsetsmodels.ArkSIABulkTargetSetResponse, error) {
	s.Logger.Info("Bulk deleting target set [%v]", bulkDeleteTargetSets)
	var bulkDeleteTargetSetsJSON map[string]interface{}
	err := mapstructure.Decode(bulkDeleteTargetSets, &bulkDeleteTargetSetsJSON)
	if err != nil {
		return nil, err
	}
	response, err := s.client.Delete(context.Background(), bulkTargetSetsURL, bulkDeleteTargetSetsJSON)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.GlobalLogger.Warning("Error closing response body")
		}
	}(response.Body)
	if response.StatusCode != http.StatusMultiStatus {
		return nil, fmt.Errorf("failed to bulk delete target set - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	bulkTargetSetRespJSON, err := common.DeserializeJSONSnake(response.Body)
	if err != nil {
		return nil, err
	}
	var bulkTargetSetsResp targetsetsmodels.ArkSIABulkTargetSetResponse
	err = mapstructure.Decode(bulkTargetSetRespJSON, &bulkTargetSetsResp)
	if err != nil {
		return nil, err
	}
	return &bulkTargetSetsResp, nil
}

// UpdateTargetSet updates a target set.
func (s *ArkSIATargetSetsWorkspaceService) UpdateTargetSet(updateTargetSet *targetsetsmodels.ArkSIAUpdateTargetSet) (*targetsetsmodels.ArkSIATargetSet, error) {
	s.Logger.Info("Updating target set [%s]", updateTargetSet.Name)
	var updateTargetSetJSON map[string]interface{}
	err := mapstructure.Decode(updateTargetSet, &updateTargetSetJSON)
	if err != nil {
		return nil, err
	}
	delete(updateTargetSetJSON, "name")
	delete(updateTargetSetJSON, "new_name")
	if updateTargetSet.NewName != "" {
		updateTargetSetJSON["name"] = updateTargetSet.NewName
	}
	response, err := s.client.Post(context.Background(), fmt.Sprintf(targetSetURL, updateTargetSet.Name), updateTargetSetJSON)
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
		return nil, fmt.Errorf("failed to update target set - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	targetSetJSON, err := common.DeserializeJSONSnake(response.Body)
	if err != nil {
		return nil, err
	}
	var targetSet targetsetsmodels.ArkSIATargetSet
	err = mapstructure.Decode(targetSetJSON["target_set"], &targetSet)
	if err != nil {
		return nil, err
	}
	return &targetSet, nil
}

// ListTargetSets lists all target sets.
func (s *ArkSIATargetSetsWorkspaceService) ListTargetSets() ([]*targetsetsmodels.ArkSIATargetSet, error) {
	s.Logger.Info("Listing all target sets")
	response, err := s.client.Get(context.Background(), targetSetsURL, nil)
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
		return nil, fmt.Errorf("failed to list target sets - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	targetSetsResponseJSON, err := common.DeserializeJSONSnake(response.Body)
	if err != nil {
		return nil, err
	}
	var targetSets []*targetsetsmodels.ArkSIATargetSet
	err = mapstructure.Decode(targetSetsResponseJSON["target_sets"], &targetSets)
	if err != nil {
		return nil, err
	}
	return targetSets, nil
}

// ListTargetSetsBy filters target sets by the provided filter.
func (s *ArkSIATargetSetsWorkspaceService) ListTargetSetsBy(targetSetsFilter *targetsetsmodels.ArkSIATargetSetsFilter) ([]*targetsetsmodels.ArkSIATargetSet, error) {
	s.Logger.Info("Listing target sets by filter [%v]", targetSetsFilter)
	targetSets, err := s.ListTargetSets()
	if err != nil {
		return nil, err
	}
	if targetSetsFilter.Name != "" {
		var filteredTargetSets []*targetsetsmodels.ArkSIATargetSet
		for _, targetSet := range targetSets {
			if match, err := regexp.MatchString(targetSetsFilter.Name, targetSet.Name); err == nil && match {
				filteredTargetSets = append(filteredTargetSets, targetSet)
			}
		}
		targetSets = filteredTargetSets
	}
	if targetSetsFilter.SecretType != "" {
		var filteredTargetSets []*targetsetsmodels.ArkSIATargetSet
		for _, targetSet := range targetSets {
			if match, err := regexp.MatchString(targetSetsFilter.SecretType, targetSet.SecretType); err == nil && match {
				filteredTargetSets = append(filteredTargetSets, targetSet)
			}
		}
		targetSets = filteredTargetSets
	}
	return targetSets, nil
}

// TargetSet retrieves a target set by name.
func (s *ArkSIATargetSetsWorkspaceService) TargetSet(getTargetSet *targetsetsmodels.ArkSIAGetTargetSet) (*targetsetsmodels.ArkSIATargetSet, error) {
	s.Logger.Info("Getting target set [%s]", getTargetSet.Name)
	response, err := s.client.Get(context.Background(), fmt.Sprintf(targetSetURL, getTargetSet.Name), nil)
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
		return nil, fmt.Errorf("failed to get target set - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	targetSetJSON, err := common.DeserializeJSONSnake(response.Body)
	if err != nil {
		return nil, err
	}
	var targetSet targetsetsmodels.ArkSIATargetSet
	err = mapstructure.Decode(targetSetJSON["target_set"], &targetSet)
	if err != nil {
		return nil, err
	}
	return &targetSet, nil
}

// TargetSetsStats retrieves statistics about target sets.
func (s *ArkSIATargetSetsWorkspaceService) TargetSetsStats() (*targetsetsmodels.ArkSIATargetSetsStats, error) {
	targetSets, err := s.ListTargetSets()
	if err != nil {
		return nil, err
	}
	var targetSetsStats targetsetsmodels.ArkSIATargetSetsStats
	targetSetsStats.TargetSetsCount = len(targetSets)
	targetSetsStats.TargetSetsCountPerSecretType = make(map[string]int)
	for _, targetSet := range targetSets {
		if _, ok := targetSetsStats.TargetSetsCountPerSecretType[targetSet.SecretType]; !ok {
			targetSetsStats.TargetSetsCountPerSecretType[targetSet.SecretType] = 0
		}
		targetSetsStats.TargetSetsCountPerSecretType[targetSet.SecretType]++
	}
	return &targetSetsStats, nil
}

// ServiceConfig returns the service configuration for the ArkSIATargetSetsWorkspaceService.
func (s *ArkSIATargetSetsWorkspaceService) ServiceConfig() services.ArkServiceConfig {
	return SIATargetSetsWorkspaceServiceConfig
}
