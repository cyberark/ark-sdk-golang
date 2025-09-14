package uap

import (
	"context"
	"fmt"

	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/common"
	"github.com/cyberark/ark-sdk-golang/pkg/common/isp"
	uapcommonmodels "github.com/cyberark/ark-sdk-golang/pkg/services/uap/common/models"
	"github.com/mitchellh/mapstructure"

	"io"
	"net/http"
	"reflect"
)

const (
	policiesURL = "/api/policies"
	policyURL   = "/api/policies/%s"
)

// ArkUAPBasePolicyPage is a page of Raw UAP items.
type ArkUAPBasePolicyPage = common.ArkPage[map[string]interface{}]

// ArkUAPBaseService is the base service for managing UAP policies.
type ArkUAPBaseService struct {
	logger  *common.ArkLogger
	ispAuth *auth.ArkISPAuth
	client  *isp.ArkISPServiceClient
}

// NewArkUAPBaseService creates a new instance of ArkUAPBaseService.
func NewArkUAPBaseService(ispAuth *auth.ArkISPAuth) (*ArkUAPBaseService, error) {
	uapService := &ArkUAPBaseService{
		logger:  common.GetLogger("ArkUAPService", common.Unknown),
		ispAuth: ispAuth,
	}
	client, err := isp.FromISPAuth(ispAuth, "uap", ".", "", uapService.refreshUapAuth)
	if err != nil {
		return nil, err
	}
	uapService.client = client
	return uapService, nil
}

func (s *ArkUAPBaseService) refreshUapAuth(client *common.ArkClient) error {
	err := isp.RefreshClient(client, s.ispAuth)
	if err != nil {
		return err
	}
	return nil
}

// BaseAddPolicy adds a new policy.
func (s *ArkUAPBaseService) BaseAddPolicy(addPolicy map[string]interface{}) (*uapcommonmodels.ArkUAPResponse, error) {
	s.logger.Info("Adding new policy")
	response, err := s.client.Post(context.Background(), policiesURL, addPolicy)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.GlobalLogger.Warning("Failed to close response body: %v", err)
		}
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to add policy - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	policyIDJSON, err := common.DeserializeJSONSnake(response.Body)
	if err != nil {
		return nil, err
	}
	var policyResponse uapcommonmodels.ArkUAPResponse
	err = mapstructure.Decode(policyIDJSON, &policyResponse)
	if err != nil {
		return nil, err
	}
	return &policyResponse, nil
}

// BasePolicy retrieves a policy by ID.
func (s *ArkUAPBaseService) BasePolicy(policyID string, schema *reflect.Type) (map[string]interface{}, error) {
	s.logger.Info("Retrieving policy [%s]", policyID)
	response, err := s.client.Get(context.Background(), fmt.Sprintf(policyURL, policyID), nil)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.GlobalLogger.Warning("Failed to close response body: %v", err)
		}
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to retrieve policy - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	policyJSON, err := common.DeserializeJSONSnakeSchema(response.Body, schema)
	if err != nil {
		return nil, err
	}
	return policyJSON.(map[string]interface{}), nil
}

// BaseUpdatePolicy updates an existing policy.
func (s *ArkUAPBaseService) BaseUpdatePolicy(policyID string, updatePolicy map[string]interface{}) error {
	s.logger.Info("Updating policy [%s]", policyID)
	response, err := s.client.Put(context.Background(), fmt.Sprintf(policyURL, policyID), updatePolicy)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.GlobalLogger.Warning("Failed to close response body: %v", err)
		}
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update policy - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	return nil
}

// BaseDeletePolicy deletes a policy by ID.
func (s *ArkUAPBaseService) BaseDeletePolicy(policyID string) error {
	s.logger.Info("Deleting policy [%s]", policyID)
	response, err := s.client.Delete(context.Background(), fmt.Sprintf(policyURL, policyID), nil)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.GlobalLogger.Warning("Failed to close response body: %v", err)
		}
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete policy - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	return nil
}

// BaseListPolicies retrieves all policies with optional filters.
func (s *ArkUAPBaseService) BaseListPolicies(filters *uapcommonmodels.ArkUAPFilters) (<-chan *ArkUAPBasePolicyPage, error) {
	s.logger.Info("Listing policies")
	if filters == nil {
		filters = uapcommonmodels.NewArkUAPFilters()
	}

	pageChannel := make(chan *ArkUAPBasePolicyPage)
	go func() {
		defer close(pageChannel)

		var nextToken string
		var prevToken string
		pageCount := 0

		for {
			if pageCount >= filters.MaxPages {
				break
			}

			pageCount++

			// Build query parameters
			request := uapcommonmodels.ArkUAPGetAccessPoliciesRequest{
				Filters:   filters,
				NextToken: nextToken,
			}
			queryParams := request.BuildGetQueryParams()
			queryParamsJSON, err := common.SerializeJSONCamel(queryParams)
			if err != nil {
				s.logger.Error("Failed to serialize query parameters: %v", err)
				return
			}
			queryParamsJSONParams := make(map[string]string)
			for key, value := range queryParamsJSON {
				queryParamsJSONParams[key] = fmt.Sprintf("%v", value)
			}

			// Make API call
			s.logger.Info("Requesting policies with next_token [%s] [%v]", nextToken, queryParamsJSONParams)
			response, err := s.client.Get(context.Background(), policiesURL, queryParamsJSONParams)
			if err != nil {
				s.logger.Error("Failed to list policies: %v", err)
				return
			}
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					common.GlobalLogger.Warning("Error closing response body")
				}
			}(response.Body)

			// Check response status
			if response.StatusCode != http.StatusOK {
				s.logger.Error("Failed to list policies - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
				return
			}

			// Parse response
			resultJSON, err := common.DeserializeJSONSnake(response.Body)
			if err != nil {
				s.logger.Error("Failed to decode response: %v", err)
				return
			}
			policiesJSONs, ok := resultJSON.(map[string]interface{})["results"].([]interface{})
			if !ok {
				s.logger.Error("Response does not contain 'results' key")
				return
			}
			policiesJSONsOut := make([]*map[string]interface{}, len(policiesJSONs))
			for i, policyJSONInterface := range policiesJSONs {
				// Convert to snake_case
				policyJSON, ok := policyJSONInterface.(map[string]interface{})
				if !ok {
					continue
				}
				policiesJSONsOut[i] = &policyJSON
			}

			// Send page to channel
			pageChannel <- &ArkUAPBasePolicyPage{Items: policiesJSONsOut}

			// Update tokens
			tempNextToken, ok := resultJSON.(map[string]interface{})["next_token"].(string)
			if !ok {
				s.logger.Error("Response does not contain 'next_token' key or it is not a string")
				return
			}
			prevToken, nextToken = nextToken, tempNextToken

			// Break if no next token or pagination loop detected
			if nextToken == "" || nextToken == prevToken {
				if nextToken == prevToken {
					s.logger.Error("Pagination stuck: next_token did not change between requests")
				}
				break
			}
			if len(policiesJSONs) < queryParams.Limit {
				s.logger.Info("No more policies to retrieve, breaking pagination loop")
				break
			}
		}
	}()

	return pageChannel, nil
}

// BasePolicyByName retrieves a policy by its name.
func (s *ArkUAPBaseService) BasePolicyByName(policyName string) (map[string]interface{}, error) {
	s.logger.Info("Retrieving policy by name [%s]", policyName)
	filters := uapcommonmodels.NewArkUAPFilters()
	filters.TextSearch = policyName
	policies, err := s.BaseListPolicies(filters)
	if err != nil {
		return nil, err
	}

	for page := range policies {
		for _, policy := range page.Items {
			metadataJSON, ok := (*policy)["metadata"].(map[string]interface{})
			if !ok {
				continue
			}
			var metadata uapcommonmodels.ArkUAPMetadata
			err = mapstructure.Decode(metadataJSON, &metadata)
			if err != nil {
				continue
			}
			if metadata.Name == policyName {
				return *policy, nil
			}
		}
	}
	return nil, fmt.Errorf("policy with name '%s' not found", policyName)
}

// BasePolicyStatus retrieves the status of a policy by its ID or name.
func (s *ArkUAPBaseService) BasePolicyStatus(policyID string, policyName string, schema *reflect.Type) (string, error) {
	s.logger.Info("Retrieving policy status for [%s] with name [%s]", policyID, policyName)
	var policy map[string]interface{}
	var err error
	if policyID != "" {
		policy, err = s.BasePolicy(policyID, schema)
		if err != nil {
			return "", fmt.Errorf("failed to retrieve policy status for ID '%s' and name '%s': %w", policyID, policyName, err)
		}
	} else if policyName != "" {
		policy, err = s.BasePolicyByName(policyName)
		if err != nil {
			return "", fmt.Errorf("failed to retrieve policy status for ID '%s' and name '%s': %w", policyID, policyName, err)
		}
	} else {
		return "", fmt.Errorf("either policyID or policyName must be provided to retrieve policy status")
	}
	metadataJSON, ok := policy["metadata"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("policy metadata not found for ID '%s' and name '%s'", policyID, policyName)
	}
	var metadata uapcommonmodels.ArkUAPMetadata
	err = mapstructure.Decode(metadataJSON, &metadata)
	if err != nil {
		return "", fmt.Errorf("failed to decode policy metadata for ID '%s' and name '%s': %w", policyID, policyName, err)
	}
	return metadata.Status.Status, nil
}

// BasePoliciesStats retrieves statistics about policies.
func (s *ArkUAPBaseService) BasePoliciesStats(filters *uapcommonmodels.ArkUAPFilters) (*uapcommonmodels.ArkUAPPoliciesStats, error) {
	policiesStats := &uapcommonmodels.ArkUAPPoliciesStats{
		PoliciesCount:            0,
		PoliciesCountPerStatus:   make(map[string]int),
		PoliciesCountPerProvider: make(map[string]int),
	}
	s.logger.Info("Retrieving policies stats")
	policies, err := s.BaseListPolicies(filters)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve policies stats: %w", err)
	}
	for page := range policies {
		for _, policy := range page.Items {
			policiesStats.PoliciesCount++
			metadataJSON, ok := (*policy)["metadata"].(map[string]interface{})
			if !ok {
				continue
			}
			var metadata uapcommonmodels.ArkUAPMetadata
			err = mapstructure.Decode(metadataJSON, &metadata)
			if err != nil {
				continue
			}
			if _, ok = policiesStats.PoliciesCountPerStatus[metadata.Status.Status]; !ok {
				policiesStats.PoliciesCountPerStatus[metadata.Status.Status] = 0
			}
			if _, ok = policiesStats.PoliciesCountPerProvider[metadata.PolicyEntitlement.LocationType]; !ok {
				policiesStats.PoliciesCountPerProvider[metadata.PolicyEntitlement.LocationType] = 0
			}
			policiesStats.PoliciesCountPerStatus[metadata.Status.Status]++
			policiesStats.PoliciesCountPerProvider[metadata.PolicyEntitlement.LocationType]++
		}
	}
	return policiesStats, nil
}
