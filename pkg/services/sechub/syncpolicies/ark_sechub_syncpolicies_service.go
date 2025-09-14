package syncpolicies

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/common"
	"github.com/cyberark/ark-sdk-golang/pkg/common/isp"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	syncpoliciesmodels "github.com/cyberark/ark-sdk-golang/pkg/services/sechub/syncpolicies/models"
	"github.com/mitchellh/mapstructure"
)

const (
	sechubURL      = "/api/policies"
	policyURL      = "/api/policies/%s"
	policyStateURL = "/api/policies/%s/state"
)

// ArkSecHubSyncPoliciesPage is a page of ArkSecHubPolicy items.
type ArkSecHubSyncPoliciesPage = common.ArkPage[syncpoliciesmodels.ArkSecHubPolicy]

// ArkSecHubSyncPoliciesService is the service for retrieve Secrets Hub Sync Policies
type ArkSecHubSyncPoliciesService struct {
	services.ArkService
	*services.ArkBaseService
	ispAuth *auth.ArkISPAuth
	client  *isp.ArkISPServiceClient
}

// NewArkSecHubSyncPoliciesService creates a new instance of ArkSecHubSyncPoliciesService.
func NewArkSecHubSyncPoliciesService(authenticators ...auth.ArkAuth) (*ArkSecHubSyncPoliciesService, error) {
	syncPoliciesService := &ArkSecHubSyncPoliciesService{}
	var syncPoliciesServiceInterface services.ArkService = syncPoliciesService
	baseService, err := services.NewArkBaseService(syncPoliciesServiceInterface, authenticators...)
	if err != nil {
		return nil, err
	}
	ispBaseAuth, err := baseService.Authenticator("isp")
	if err != nil {
		return nil, err
	}
	ispAuth := ispBaseAuth.(*auth.ArkISPAuth)
	client, err := isp.FromISPAuth(ispAuth, "secretshub", ".", "", syncPoliciesService.refreshSecHubAuth)
	if err != nil {
		return nil, err
	}
	syncPoliciesService.client = client
	syncPoliciesService.ispAuth = ispAuth
	syncPoliciesService.ArkBaseService = baseService
	return syncPoliciesService, nil
}

func (s *ArkSecHubSyncPoliciesService) refreshSecHubAuth(client *common.ArkClient) error {
	err := isp.RefreshClient(client, s.ispAuth)
	if err != nil {
		return err
	}
	return nil
}

func (s *ArkSecHubSyncPoliciesService) getSyncPoliciesWithFilters(
	projection string,
	filter string,
) (<-chan *ArkSecHubSyncPoliciesPage, error) {
	query := map[string]string{}
	if projection != "" {
		query["projection"] = projection
	}
	if filter != "" {
		query["filter"] = filter
	}
	results := make(chan *ArkSecHubSyncPoliciesPage)
	go func() {
		defer close(results)
		for {
			response, err := s.client.Get(context.Background(), sechubURL, query)
			if err != nil {
				s.Logger.Error("Failed to list Sync Policies: %v", err)
				return
			}
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					common.GlobalLogger.Warning("Error closing response body")
				}
			}(response.Body)
			if response.StatusCode != http.StatusOK {
				s.Logger.Error("Failed to list Sync Policies - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
				return
			}
			result, err := common.DeserializeJSONSnake(response.Body)
			if err != nil {
				s.Logger.Error("Failed to decode response: %v", err)
				return
			}
			resultMap := result.(map[string]interface{})
			var syncPoliciesJSON []interface{}
			if syncPolicy, ok := resultMap["policies"]; ok {
				syncPoliciesJSON = syncPolicy.([]interface{})
			} else {
				s.Logger.Error("Failed to list Sync Policies, unexpected result")
				return
			}
			for i, syncPolicy := range syncPoliciesJSON {
				if syncPoliciesMap, ok := syncPolicy.(map[string]interface{}); ok {
					if syncPolicyID, ok := syncPoliciesMap["id"]; ok {
						syncPoliciesJSON[i].(map[string]interface{})["id"] = syncPolicyID
					}
				}
			}
			var syncPolicies []*syncpoliciesmodels.ArkSecHubPolicy
			if err := mapstructure.Decode(syncPoliciesJSON, &syncPolicies); err != nil {
				s.Logger.Error("Failed to validate Sync Policies: %v", err)
				return
			}
			results <- &ArkSecHubSyncPoliciesPage{Items: syncPolicies}
			if nextLink, ok := resultMap["nextLink"].(string); ok {
				nextQuery, _ := url.Parse(nextLink)
				queryValues := nextQuery.Query()
				query = make(map[string]string)
				for key, values := range queryValues {
					if len(values) > 0 {
						query[key] = values[0]
					}
				}
			} else {
				break
			}
		}
	}()
	return results, nil
}

// ListSyncPolicies returns a channel of ArkSecHubSyncPoliciesPage containing all Sync Policies.
func (s *ArkSecHubSyncPoliciesService) ListSyncPolicies(syncPolicies *syncpoliciesmodels.ArkSecHubGetSyncPolicies) (<-chan *ArkSecHubSyncPoliciesPage, error) {
	var projection string
	if syncPolicies.Projection != "" {
		projection = syncPolicies.Projection
	}
	return s.getSyncPoliciesWithFilters(
		projection,
		"",
	)
}

// ListSyncPoliciesBy returns a channel of ArkSecHubSyncPoliciesPage containing secrets filtered by the given filters.
func (s *ArkSecHubSyncPoliciesService) ListSyncPoliciesBy(syncPoliciesFilters *syncpoliciesmodels.ArkSecHubSyncPoliciesFilters) (<-chan *ArkSecHubSyncPoliciesPage, error) {
	var projection string
	if syncPoliciesFilters.Projection != "" {
		projection = syncPoliciesFilters.Projection
	}
	return s.getSyncPoliciesWithFilters(
		projection,
		syncPoliciesFilters.Filters,
	)
}

// SyncPolicy returns an individual sync policy
// https://api-docs.cyberark.com/docs/secretshub-api/f5jjh0rv9ivfs-get-sync-policy
func (s *ArkSecHubSyncPoliciesService) SyncPolicy(
	getSyncPolicy *syncpoliciesmodels.ArkSecHubGetSyncPolicy) (*syncpoliciesmodels.ArkSecHubPolicy, error) {
	s.Logger.Info("Retrieving sync policy [%s]", getSyncPolicy.PolicyID)
	query := map[string]string{}
	if getSyncPolicy.Projection != "" {
		query["projection"] = getSyncPolicy.Projection
	} else {
		query["projection"] = "REGULAR"
	}
	response, err := s.client.Get(context.Background(), fmt.Sprintf(policyURL, getSyncPolicy.PolicyID), query)
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
		return nil, fmt.Errorf("failed to retrieve sync policy - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	syncPolicyJSON, err := common.DeserializeJSONSnake(response.Body)
	if err != nil {
		return nil, err
	}
	syncPolicyJSONMap := syncPolicyJSON.(map[string]interface{})
	if syncPolicyID, ok := syncPolicyJSONMap["id"]; ok {
		syncPolicyJSONMap["id"] = syncPolicyID
	}
	var syncPolicy syncpoliciesmodels.ArkSecHubPolicy
	err = mapstructure.Decode(syncPolicyJSONMap, &syncPolicy)
	if err != nil {
		return nil, err
	}
	return &syncPolicy, nil
}

// CreateSyncPolicy creates a new sync policy
// https://api-docs.cyberark.com/docs/secretshub-api/3kf2d2n01bm5x-create-sync-policy
func (s *ArkSecHubSyncPoliciesService) CreateSyncPolicy(syncPolicy *syncpoliciesmodels.ArkSechubCreateSyncPolicy) (*syncpoliciesmodels.ArkSecHubPolicy, error) {
	s.Logger.Info("Creating sync policy [%s]", syncPolicy.Name)
	createSyncPolicyJSON, err := common.SerializeJSONCamel(syncPolicy)
	if err != nil {
		return nil, err
	}
	if syncPolicy.Description != "" {
		delete(createSyncPolicyJSON, "description")
		createSyncPolicyJSON["description"] = syncPolicy.Description
	}
	// Documentation states that default is allowed however in testing this seems to cause failures.
	if syncPolicy.Transformation.Predefined == "default" {
		delete(createSyncPolicyJSON, "transformation")
	}
	response, err := s.client.Post(context.Background(), sechubURL, createSyncPolicyJSON)
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
		return nil, fmt.Errorf("failed to create sync policy - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	syncPolicyJSON, err := common.DeserializeJSONSnake(response.Body)
	if err != nil {
		return nil, err
	}
	var syncPolicyResponse syncpoliciesmodels.ArkSecHubPolicy
	err = mapstructure.Decode(syncPolicyJSON, &syncPolicyResponse)
	if err != nil {
		return nil, err
	}
	return &syncPolicyResponse, nil
}

// SetSyncPolicyState sets the state of a sync policy.
// https://api-docs.cyberark.com/docs/secretshub-api/by05aodbep6xy-set-sync-policy-state
func (s *ArkSecHubSyncPoliciesService) SetSyncPolicyState(
	setSyncPolicyState *syncpoliciesmodels.ArkSecHubSetSyncPolicyState) error {
	s.Logger.Info("Setting sync policy state [%s] to [%s]", setSyncPolicyState.PolicyID, setSyncPolicyState.Action)
	bodyMap := map[string]string{
		"action": setSyncPolicyState.Action,
	}
	response, err := s.client.Put(context.Background(), fmt.Sprintf(policyStateURL, setSyncPolicyState.PolicyID), bodyMap)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.GlobalLogger.Warning("Error closing response body")
		}
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to set sync policy state - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	return nil
}

// DeleteSyncPolicy deletes a specified secret store based on ID
// https://api-docs.cyberark.com/docs/secretshub-api/lgbolpf4ka7oa-delete-sync-policy
func (s *ArkSecHubSyncPoliciesService) DeleteSyncPolicy(syncPolicy *syncpoliciesmodels.ArkSecHubDeleteSyncPolicy) error {
	s.Logger.Info("Deleting secret store")
	response, err := s.client.Delete(context.Background(), fmt.Sprintf(policyURL, syncPolicy.PolicyID), nil)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.GlobalLogger.Warning("Error closing response body")
		}
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete sync policy - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	return nil
}

// SyncPoliciesStats retrieves statistics about sync policies.
func (s *ArkSecHubSyncPoliciesService) SyncPoliciesStats() (*syncpoliciesmodels.ArkSecHubSyncPoliciesStats, error) {
	s.Logger.Info("Retrieving sync policy stats")
	var projection = syncpoliciesmodels.ArkSecHubGetSyncPolicies{
		Projection: "REGULAR",
	}
	syncPoliciesChan, err := s.ListSyncPolicies(&projection)
	if err != nil {
		return nil, err
	}
	syncPolicies := make([]*syncpoliciesmodels.ArkSecHubPolicy, 0)
	for page := range syncPoliciesChan {
		syncPolicies = append(syncPolicies, page.Items...)
	}
	var syncPoliciesStats syncpoliciesmodels.ArkSecHubSyncPoliciesStats
	syncPoliciesStats.SyncPoliciesCount = len(syncPolicies)
	syncPoliciesStats.SyncPoliciesCountByCreator = make(map[string]int)
	for _, syncPolicy := range syncPolicies {
		if _, ok := syncPoliciesStats.SyncPoliciesCountByCreator[syncPolicy.CreatedBy]; !ok {
			syncPoliciesStats.SyncPoliciesCountByCreator[syncPolicy.CreatedBy] = 0
		}
		syncPoliciesStats.SyncPoliciesCountByCreator[syncPolicy.CreatedBy]++
	}
	return &syncPoliciesStats, nil
}

// ServiceConfig returns the service configuration for the ArkSecHubSyncPoliciesService.
func (s *ArkSecHubSyncPoliciesService) ServiceConfig() services.ArkServiceConfig {
	return ServiceConfig
}
