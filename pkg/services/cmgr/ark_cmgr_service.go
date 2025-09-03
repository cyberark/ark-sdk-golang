package cmgr

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/common"
	"github.com/cyberark/ark-sdk-golang/pkg/common/isp"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	cmgrmodels "github.com/cyberark/ark-sdk-golang/pkg/services/cmgr/models"
	"github.com/mitchellh/mapstructure"
)

const (
	networksURL            = "api/pool-service/networks"
	networkURL             = "api/pool-service/networks/%s"
	poolsURL               = "api/pool-service/pools"
	poolURL                = "api/pool-service/pools/%s"
	poolIdentifiersURL     = "api/pool-service/pools/%s/identifiers"
	poolIdentifiersBulkURL = "api/pool-service/pools/%s/identifiers-bulk"
	poolIdentifierURL      = "api/pool-service/pools/%s/identifiers/%s"
	poolsComponentsURL     = "api/pool-service/pools/components"
	poolComponentURL       = "api/pool-service/pools/%s/components/%s"
)

// ArkCmgrNetworkPage is a page of ArkCmgrNetwork items.
type ArkCmgrNetworkPage = common.ArkPage[cmgrmodels.ArkCmgrNetwork]

// ArkCmgrPoolPage is a page of ArkCmgrPool items.
type ArkCmgrPoolPage = common.ArkPage[cmgrmodels.ArkCmgrPool]

// ArkCmgrPoolIdentifierPage is a page of ArkCmgrPoolIdentifier items.
type ArkCmgrPoolIdentifierPage = common.ArkPage[cmgrmodels.ArkCmgrPoolIdentifier]

// ArkCmgrPoolComponentPage is a page of ArkCmgrPoolComponent items.
type ArkCmgrPoolComponentPage = common.ArkPage[cmgrmodels.ArkCmgrPoolComponent]

// CmgrServiceConfig is the configuration for the connector management service.
var CmgrServiceConfig = services.ArkServiceConfig{
	ServiceName:                "cmgr",
	RequiredAuthenticatorNames: []string{"isp"},
	OptionalAuthenticatorNames: []string{},
}

// ArkCmgrService is the service for managing connector management.
type ArkCmgrService struct {
	services.ArkService
	*services.ArkBaseService
	ispAuth *auth.ArkISPAuth
	client  *isp.ArkISPServiceClient
}

// NewArkCmgrService creates a new instance of ArkCmgrService.
func NewArkCmgrService(authenticators ...auth.ArkAuth) (*ArkCmgrService, error) {
	cmgrService := &ArkCmgrService{}
	var cmgrServiceInterface services.ArkService = cmgrService
	baseService, err := services.NewArkBaseService(cmgrServiceInterface, authenticators...)
	if err != nil {
		return nil, err
	}
	ispBaseAuth, err := baseService.Authenticator("isp")
	if err != nil {
		return nil, err
	}
	ispAuth := ispBaseAuth.(*auth.ArkISPAuth)
	client, err := isp.FromISPAuth(ispAuth, "connectormanagement", ".", "", cmgrService.refreshCmgrAuth)
	if err != nil {
		return nil, err
	}
	cmgrService.client = client
	cmgrService.ispAuth = ispAuth
	cmgrService.ArkBaseService = baseService
	return cmgrService, nil
}

func (s *ArkCmgrService) refreshCmgrAuth(client *common.ArkClient) error {
	err := isp.RefreshClient(client, s.ispAuth)
	if err != nil {
		return err
	}
	return nil
}

func listCommonPools[PageItemType any](
	logger *common.ArkLogger,
	client *isp.ArkISPServiceClient,
	name string, route string,
	commonFilter *cmgrmodels.ArkCmgrPoolsCommonFilter,
	idMappings map[string]string) (<-chan *common.ArkPage[PageItemType], error) {
	logger.Info("Listing %s", name)
	pageChannel := make(chan *common.ArkPage[PageItemType])
	go func() {
		defer close(pageChannel)
		filters := map[string]string{
			"projection": "EXTENDED",
		}
		if commonFilter != nil {
			if commonFilter.Filter != "" {
				filters["filter"] = commonFilter.Filter
			}
			if commonFilter.Order != "" {
				filters["order"] = commonFilter.Order
			}
			if commonFilter.PageSize != 0 {
				filters["pageSize"] = fmt.Sprintf("%d", commonFilter.PageSize)
			}
			if commonFilter.Sort != "" {
				filters["sort"] = commonFilter.Sort
			}
			if commonFilter.Projection != "" {
				filters["projection"] = commonFilter.Projection
			}
		}
		var contToken string
		for {
			if contToken != "" {
				filters["continuation_token"] = contToken
			}
			response, err := client.Get(context.Background(), route, filters)
			if err != nil {
				logger.Error("Failed to list %s: %v", name, err)
				return
			}
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					common.GlobalLogger.Warning("Error closing response body")
				}
			}(response.Body)
			if response.StatusCode != http.StatusOK {
				logger.Error("Failed to list %s - [%d] - [%s]", name, response.StatusCode, common.SerializeResponseToJSON(response.Body))
				return
			}
			result, err := common.DeserializeJSONSnake(response.Body)
			if err != nil {
				logger.Error("Failed to decode response for %s: %v", name, err)
				return
			}
			resultMap := result.(map[string]interface{})
			if idMappings != nil && len(idMappings) > 0 {
				for _, resourceItem := range resultMap["resources"].([]interface{}) {
					for key, value := range idMappings {
						if _, ok := resourceItem.(map[string]interface{})[key]; ok {
							resourceItem.(map[string]interface{})[value] = resourceItem.(map[string]interface{})[key]
						}
					}
					if _, ok := resourceItem.(map[string]interface{})["assigned_pools"]; ok {
						for _, pool := range resourceItem.(map[string]interface{})["assigned_pools"].([]interface{}) {
							pool.(map[string]interface{})["pool_id"] = pool.(map[string]interface{})["id"]
						}
					}
				}
			}

			var items []*PageItemType
			err = mapstructure.Decode(resultMap["resources"], &items)
			if err != nil {
				logger.Error("Failed to decode resources for %s: %v", name, err)
				return
			}
			pageChannel <- &common.ArkPage[PageItemType]{Items: items}
			pageInfo, ok := resultMap["page"].(map[string]interface{})
			if !ok || pageInfo["continuation_token"] == nil || pageInfo["continuation_token"] == "" {
				break
			}
			contToken = pageInfo["continuation_token"].(string)
			if totalResources, ok := pageInfo["total_resources_count"].(float64); ok {
				if pageSize, ok := pageInfo["page_size"].(float64); ok && totalResources == pageSize {
					break
				}
			}
		}
	}()
	return pageChannel, nil
}

// AddNetwork adds a new network to the connector management service.
func (s *ArkCmgrService) AddNetwork(addNetwork *cmgrmodels.ArkCmgrAddNetwork) (*cmgrmodels.ArkCmgrNetwork, error) {
	s.Logger.Info("Adding network [%s]", addNetwork.Name)
	var addNetworkJSON map[string]interface{}
	err := mapstructure.Decode(addNetwork, &addNetworkJSON)
	if err != nil {
		return nil, err
	}
	response, err := s.client.Post(context.Background(), networksURL, addNetworkJSON)
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
		return nil, fmt.Errorf("failed to add network - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	networkJSON, err := common.DeserializeJSONSnake(response.Body)
	if err != nil {
		return nil, err
	}
	networkJSONMap := networkJSON.(map[string]interface{})
	networkJSONMap["network_id"] = networkJSONMap["id"]
	var network cmgrmodels.ArkCmgrNetwork
	err = mapstructure.Decode(networkJSONMap, &network)
	if err != nil {
		return nil, err
	}
	return &network, nil
}

// UpdateNetwork updates an existing network in the connector management service.
func (s *ArkCmgrService) UpdateNetwork(updateNetwork *cmgrmodels.ArkCmgrUpdateNetwork) (*cmgrmodels.ArkCmgrNetwork, error) {
	s.Logger.Info("Updating network [%s]", updateNetwork.NetworkID)
	if updateNetwork.Name == "" {
		s.Logger.Info("Nothing to update")
		return s.Network(&cmgrmodels.ArkCmgrGetNetwork{NetworkID: updateNetwork.NetworkID})
	}
	var updateNetworkJSON map[string]interface{}
	err := mapstructure.Decode(updateNetwork, &updateNetworkJSON)
	if err != nil {
		return nil, err
	}
	response, err := s.client.Patch(context.Background(), fmt.Sprintf(networkURL, updateNetwork.NetworkID), updateNetworkJSON)
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
		return nil, fmt.Errorf("failed to update network - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	networkJSON, err := common.DeserializeJSONSnake(response.Body)
	if err != nil {
		return nil, err
	}
	networkJSONMap := networkJSON.(map[string]interface{})
	networkJSONMap["network_id"] = networkJSONMap["id"]
	var network cmgrmodels.ArkCmgrNetwork
	err = mapstructure.Decode(networkJSONMap, &network)
	if err != nil {
		return nil, err
	}
	return &network, nil
}

// DeleteNetwork deletes an existing network from the connector management service.
func (s *ArkCmgrService) DeleteNetwork(deleteNetwork *cmgrmodels.ArkCmgrDeleteNetwork) error {
	s.Logger.Info("Deleting network [%s]", deleteNetwork.NetworkID)
	response, err := s.client.Delete(context.Background(), fmt.Sprintf(networkURL, deleteNetwork.NetworkID), nil)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete network - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	return nil
}

// ListNetworks lists all networks in the connector management service.
func (s *ArkCmgrService) ListNetworks() (<-chan *ArkCmgrNetworkPage, error) {
	s.Logger.Info("Listing all networks")
	return listCommonPools[cmgrmodels.ArkCmgrNetwork](
		s.Logger,
		s.client,
		"networks",
		networksURL,
		nil,
		map[string]string{
			"id": "network_id",
		},
	)
}

// ListNetworksBy lists networks by the specified filter in the connector management service.
func (s *ArkCmgrService) ListNetworksBy(networksFilter *cmgrmodels.ArkCmgrNetworksFilter) (<-chan *ArkCmgrNetworkPage, error) {
	s.Logger.Info("Listing networks by filter [%v]", networksFilter)
	return listCommonPools[cmgrmodels.ArkCmgrNetwork](
		s.Logger,
		s.client,
		"networks",
		networksURL,
		&networksFilter.ArkCmgrPoolsCommonFilter,
		map[string]string{
			"id": "network_id",
		},
	)
}

// Network retrieves a specific network by its ID from the connector management service.
func (s *ArkCmgrService) Network(getNetwork *cmgrmodels.ArkCmgrGetNetwork) (*cmgrmodels.ArkCmgrNetwork, error) {
	s.Logger.Info("Retrieving network [%s]", getNetwork.NetworkID)
	response, err := s.client.Get(context.Background(), fmt.Sprintf(networkURL, getNetwork.NetworkID), nil)
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
		return nil, fmt.Errorf("failed to retrieve network - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	networkJSON, err := common.DeserializeJSONSnake(response.Body)
	if err != nil {
		return nil, err
	}
	networkJSONMap := networkJSON.(map[string]interface{})
	networkJSONMap["network_id"] = networkJSONMap["id"]
	var network cmgrmodels.ArkCmgrNetwork
	err = mapstructure.Decode(networkJSONMap, &network)
	if err != nil {
		return nil, err
	}
	return &network, nil
}

// NetworksStats retrieves statistics about networks in the connector management service.
func (s *ArkCmgrService) NetworksStats() (*cmgrmodels.ArkCmgrNetworksStats, error) {
	s.Logger.Info("Retrieving networks stats")
	networksChan, err := s.ListNetworks()
	if err != nil {
		return nil, err
	}
	networks := make([]*cmgrmodels.ArkCmgrNetwork, 0)
	for page := range networksChan {
		for _, network := range page.Items {
			networks = append(networks, network)
		}
	}
	var networksStats cmgrmodels.ArkCmgrNetworksStats
	networksStats.NetworksCount = len(networks)
	networksStats.PoolsCountPerNetwork = make(map[string]int)
	for _, network := range networks {
		networksStats.PoolsCountPerNetwork[network.Name] = len(network.AssignedPools)
	}
	return &networksStats, nil
}

// AddPool adds a new pool to the connector management service.
func (s *ArkCmgrService) AddPool(addPool *cmgrmodels.ArkCmgrAddPool) (*cmgrmodels.ArkCmgrPool, error) {
	s.Logger.Info("Adding pool [%s]", addPool.Name)
	var addPoolJSON map[string]interface{}
	err := mapstructure.Decode(addPool, &addPoolJSON)
	if err != nil {
		return nil, err
	}
	if addPool.AssignedNetworkIDs == nil || len(addPool.AssignedNetworkIDs) == 0 {
		return nil, fmt.Errorf("no networks assigned to the pool")
	}
	response, err := s.client.Post(context.Background(), poolsURL, addPoolJSON)
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
		return nil, fmt.Errorf("failed to add pool - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	poolJSON, err := common.DeserializeJSONSnake(response.Body)
	if err != nil {
		return nil, err
	}
	poolJSONMap := poolJSON.(map[string]interface{})
	poolJSONMap["pool_id"] = poolJSONMap["id"]
	var pool cmgrmodels.ArkCmgrPool
	err = mapstructure.Decode(poolJSONMap, &pool)
	if err != nil {
		return nil, err
	}
	return &pool, nil
}

// UpdatePool updates an existing pool in the connector management service.
func (s *ArkCmgrService) UpdatePool(updatePool *cmgrmodels.ArkCmgrUpdatePool) (*cmgrmodels.ArkCmgrPool, error) {
	s.Logger.Info("Updating pool [%s]", updatePool.PoolID)
	if updatePool.Name == "" && updatePool.Description == "" && updatePool.AssignedNetworkIDs == nil {
		s.Logger.Info("Nothing to update")
		return s.Pool(&cmgrmodels.ArkCmgrGetPool{PoolID: updatePool.PoolID})
	}
	var updatePoolJSON map[string]interface{}
	err := mapstructure.Decode(updatePool, &updatePoolJSON)
	if err != nil {
		return nil, err
	}
	response, err := s.client.Patch(context.Background(), fmt.Sprintf(poolURL, updatePool.PoolID), updatePoolJSON)
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
		return nil, fmt.Errorf("failed to update pool - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	poolJSON, err := common.DeserializeJSONSnake(response.Body)
	if err != nil {
		return nil, err
	}
	poolJSONMap := poolJSON.(map[string]interface{})
	poolJSONMap["pool_id"] = poolJSONMap["id"]
	var pool cmgrmodels.ArkCmgrPool
	err = mapstructure.Decode(poolJSONMap, &pool)
	if err != nil {
		return nil, err
	}
	return &pool, nil
}

// DeletePool deletes an existing pool from the connector management service.
func (s *ArkCmgrService) DeletePool(deletePool *cmgrmodels.ArkCmgrDeletePool) error {
	s.Logger.Info("Deleting pool [%s]", deletePool.PoolID)
	response, err := s.client.Delete(context.Background(), fmt.Sprintf(poolURL, deletePool.PoolID), nil)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete pool - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	return nil
}

// ListPools lists all pools in the connector management service.
func (s *ArkCmgrService) ListPools() (<-chan *ArkCmgrPoolPage, error) {
	s.Logger.Info("Listing all pools")
	return listCommonPools[cmgrmodels.ArkCmgrPool](
		s.Logger,
		s.client,
		"pools",
		poolsURL,
		nil,
		map[string]string{
			"id": "pool_id",
		},
	)
}

// ListPoolsBy lists pools by the specified filter in the connector management service.
func (s *ArkCmgrService) ListPoolsBy(poolsFilter *cmgrmodels.ArkCmgrPoolsFilter) (<-chan *ArkCmgrPoolPage, error) {
	s.Logger.Info("Listing pools by filter [%v]", poolsFilter)
	return listCommonPools[cmgrmodels.ArkCmgrPool](
		s.Logger,
		s.client,
		"pools",
		poolsURL,
		&poolsFilter.ArkCmgrPoolsCommonFilter,
		map[string]string{
			"id": "pool_id",
		},
	)
}

// Pool retrieves a specific pool by its ID from the connector management service.
func (s *ArkCmgrService) Pool(getPool *cmgrmodels.ArkCmgrGetPool) (*cmgrmodels.ArkCmgrPool, error) {
	s.Logger.Info("Retrieving pool [%s]", getPool.PoolID)
	response, err := s.client.Get(context.Background(), fmt.Sprintf(poolURL, getPool.PoolID), nil)
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
		return nil, fmt.Errorf("failed to retrieve pool - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	poolJSON, err := common.DeserializeJSONSnake(response.Body)
	if err != nil {
		return nil, err
	}
	poolJSONMap := poolJSON.(map[string]interface{})
	poolJSONMap["pool_id"] = poolJSONMap["id"]
	var pool cmgrmodels.ArkCmgrPool
	err = mapstructure.Decode(poolJSONMap, &pool)
	if err != nil {
		return nil, err
	}
	return &pool, nil
}

// PoolsStats retrieves statistics about pools in the connector management service.
func (s *ArkCmgrService) PoolsStats() (*cmgrmodels.ArkCmgrPoolsStats, error) {
	s.Logger.Info("Retrieving pools stats")
	poolsChan, err := s.ListPools()
	if err != nil {
		return nil, err
	}
	pools := make([]*cmgrmodels.ArkCmgrPool, 0)
	for page := range poolsChan {
		for _, pool := range page.Items {
			pools = append(pools, pool)
		}
	}
	var poolsStats cmgrmodels.ArkCmgrPoolsStats
	poolsStats.PoolsCount = len(pools)
	poolsStats.NetworksCountPerPool = make(map[string]int)
	poolsStats.IdentifiersCountPerPool = make(map[string]int)
	poolsStats.ComponentsCountPerPool = make(map[string]map[string]int)
	for _, pool := range pools {
		poolsStats.NetworksCountPerPool[pool.Name] = len(pool.AssignedNetworkIDs)
		poolsStats.IdentifiersCountPerPool[pool.Name] = pool.IdentifiersCount
		poolsStats.ComponentsCountPerPool[pool.Name] = pool.ComponentsCount
	}
	return &poolsStats, nil
}

// AddPoolIdentifier adds a new identifier to a specific pool in the connector management service.
func (s *ArkCmgrService) AddPoolIdentifier(addPoolIdentifier *cmgrmodels.ArkCmgrAddPoolSingleIdentifier) (*cmgrmodels.ArkCmgrPoolIdentifier, error) {
	s.Logger.Info("Adding pool identifier [%v]", addPoolIdentifier)
	var addPoolIdentifierJSON map[string]interface{}
	err := mapstructure.Decode(addPoolIdentifier, &addPoolIdentifierJSON)
	if err != nil {
		return nil, err
	}
	delete(addPoolIdentifierJSON, "pool_id")
	response, err := s.client.Post(context.Background(), fmt.Sprintf(poolIdentifiersURL, addPoolIdentifier.PoolID), addPoolIdentifierJSON)
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
		return nil, fmt.Errorf("failed to add pool identifier - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	poolIdentifierJSON, err := common.DeserializeJSONSnake(response.Body)
	if err != nil {
		return nil, err
	}
	poolIdentifierJSONMap := poolIdentifierJSON.(map[string]interface{})
	poolIdentifierJSONMap["identifier_id"] = poolIdentifierJSONMap["id"]
	var poolIdentifier cmgrmodels.ArkCmgrPoolIdentifier
	err = mapstructure.Decode(poolIdentifierJSONMap, &poolIdentifier)
	if err != nil {
		return nil, err
	}
	return &poolIdentifier, nil
}

// AddPoolIdentifiers adds multiple identifiers to a specific pool in the connector management service.
func (s *ArkCmgrService) AddPoolIdentifiers(addPoolIdentifiers *cmgrmodels.ArkCmgrAddPoolBulkIdentifier) (*cmgrmodels.ArkCmgrPoolIdentifiers, error) {
	s.Logger.Info("Adding pool identifiers [%v]", addPoolIdentifiers)
	requests := make(map[string]interface{})
	for index, identifier := range addPoolIdentifiers.Identifiers {
		identifierMap := make(map[string]interface{})
		err := mapstructure.Decode(identifier, &identifierMap)
		if err != nil {
			return nil, fmt.Errorf("failed to decode identifier: %w", err)
		}
		requests[fmt.Sprintf("%d", index+1)] = identifierMap
	}
	payload := map[string]interface{}{
		"requests": requests,
	}
	response, err := s.client.Post(context.Background(), fmt.Sprintf(poolIdentifiersBulkURL, addPoolIdentifiers.PoolID), payload)
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
		return nil, fmt.Errorf("failed to add pool identifiers - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	bulkResponsesJSON, err := common.DeserializeJSONSnake(response.Body)
	if err != nil {
		return nil, err
	}
	var bulkResponses *cmgrmodels.ArkCmgrBulkResponses
	err = mapstructure.Decode(bulkResponsesJSON, &bulkResponses)
	if err != nil {
		return nil, fmt.Errorf("failed to decode bulk responses: %w", err)
	}
	identifiers := make([]*cmgrmodels.ArkCmgrPoolIdentifier, 0)
	for _, identifierResponse := range bulkResponses.Responses {
		if identifierResponse.StatusCode != http.StatusCreated {
			return nil, fmt.Errorf("failed to add pool identifiers bulk - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
		}
		identifierResponse.Body["identifier_id"] = identifierResponse.Body["id"]
		var identifier cmgrmodels.ArkCmgrPoolIdentifier
		err := mapstructure.Decode(identifierResponse.Body, &identifier)
		if err != nil {
			return nil, fmt.Errorf("failed to decode identifier response body: %w", err)
		}
		identifiers = append(identifiers, &identifier)
	}
	return &cmgrmodels.ArkCmgrPoolIdentifiers{Identifiers: identifiers}, nil
}

// UpdatePoolIdentifier updates an existing identifier in a specific pool in the connector management service.
func (s *ArkCmgrService) UpdatePoolIdentifier(updatePoolIdentifier *cmgrmodels.ArkCmgrUpdatePoolIdentifier) (*cmgrmodels.ArkCmgrPoolIdentifier, error) {
	s.Logger.Info("Updating pool identifier [%s] from pool [%s]", updatePoolIdentifier.IdentifierID, updatePoolIdentifier.PoolID)
	err := s.DeletePoolIdentifier(&cmgrmodels.ArkCmgrDeletePoolSingleIdentifier{
		IdentifierID: updatePoolIdentifier.IdentifierID,
		PoolID:       updatePoolIdentifier.PoolID,
	})
	if err != nil {
		return nil, err
	}
	return s.AddPoolIdentifier(&cmgrmodels.ArkCmgrAddPoolSingleIdentifier{
		Type:   updatePoolIdentifier.Type,
		Value:  updatePoolIdentifier.Value,
		PoolID: updatePoolIdentifier.PoolID,
	})
}

// DeletePoolIdentifier deletes an identifier from a specific pool in the connector management service.
func (s *ArkCmgrService) DeletePoolIdentifier(deletePoolIdentifier *cmgrmodels.ArkCmgrDeletePoolSingleIdentifier) error {
	s.Logger.Info("Deleting pool identifier [%s]", deletePoolIdentifier.IdentifierID)
	response, err := s.client.Delete(context.Background(), fmt.Sprintf(poolIdentifierURL, deletePoolIdentifier.PoolID, deletePoolIdentifier.IdentifierID), nil)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete pool identifier - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	return nil
}

// DeletePoolIdentifiers deletes multiple identifiers from a specific pool in the connector management service.
func (s *ArkCmgrService) DeletePoolIdentifiers(deletePoolIdentifiers *cmgrmodels.ArkCmgrDeletePoolBulkIdentifier) error {
	s.Logger.Info("Deleting pool identifiers [%s]", deletePoolIdentifiers.PoolID)
	requests := make(map[string]interface{})
	for index, identifier := range deletePoolIdentifiers.Identifiers {
		identifierMap := make(map[string]interface{})
		err := mapstructure.Decode(identifier, &identifierMap)
		if err != nil {
			return fmt.Errorf("failed to decode identifier: %w", err)
		}
		requests[fmt.Sprintf("%d", index+1)] = identifierMap
	}
	payload := map[string]interface{}{
		"requests": requests,
	}
	response, err := s.client.Post(context.Background(), fmt.Sprintf(poolIdentifiersBulkURL, deletePoolIdentifiers.PoolID), payload)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.GlobalLogger.Warning("Error closing response body")
		}
	}(response.Body)
	if response.StatusCode != http.StatusMultiStatus {
		return fmt.Errorf("failed to delete pool identifiers - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	bulkResponsesJSON, err := common.DeserializeJSONSnake(response.Body)
	if err != nil {
		return fmt.Errorf("failed to decode bulk responses: %w", err)
	}
	var bulkResponses cmgrmodels.ArkCmgrBulkResponses
	err = mapstructure.Decode(bulkResponsesJSON, &bulkResponses)
	if err != nil {
		return fmt.Errorf("failed to decode bulk responses: %w", err)
	}
	for _, identifierResponse := range bulkResponses.Responses {
		if identifierResponse.StatusCode != http.StatusNoContent {
			return fmt.Errorf("failed to delete pool identifiers - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
		}
	}
	return nil
}

// ListPoolIdentifiers lists all identifiers in a specific pool in the connector management service.
func (s *ArkCmgrService) ListPoolIdentifiers(listPoolIdentifiers *cmgrmodels.ArkCmgrListPoolIdentifiers) (<-chan *ArkCmgrPoolIdentifierPage, error) {
	s.Logger.Info("Listing pool identifiers [%v]", listPoolIdentifiers)
	return listCommonPools[cmgrmodels.ArkCmgrPoolIdentifier](
		s.Logger,
		s.client,
		"pool identifiers",
		fmt.Sprintf(poolIdentifiersURL, listPoolIdentifiers.PoolID),
		nil,
		map[string]string{
			"id": "identifier_id",
		},
	)
}

// ListPoolIdentifiersBy lists identifiers by the specified filter in a specific pool in the connector management service.
func (s *ArkCmgrService) ListPoolIdentifiersBy(identifiersFilters *cmgrmodels.ArkCmgrPoolIdentifiersFilter) (<-chan *ArkCmgrPoolIdentifierPage, error) {
	s.Logger.Info("Listing pool identifiers by filter [%v]", identifiersFilters)
	return listCommonPools[cmgrmodels.ArkCmgrPoolIdentifier](
		s.Logger,
		s.client,
		"pool identifiers",
		fmt.Sprintf(poolIdentifiersURL, identifiersFilters.PoolID),
		&identifiersFilters.ArkCmgrPoolsCommonFilter,
		map[string]string{
			"id": "identifier_id",
		},
	)
}

// PoolIdentifier retrieves a specific identifier by its ID from a specific pool in the connector management service.
func (s *ArkCmgrService) PoolIdentifier(getIdentifier *cmgrmodels.ArkCmgrGetPoolIdentifier) (*cmgrmodels.ArkCmgrPoolIdentifier, error) {
	s.Logger.Info("Retrieving pool identifier [%s] from pool [%s]", getIdentifier.IdentifierID, getIdentifier.PoolID)
	identifiers, err := s.ListPoolIdentifiers(&cmgrmodels.ArkCmgrListPoolIdentifiers{PoolID: getIdentifier.PoolID})
	if err != nil {
		return nil, err
	}
	for page := range identifiers {
		for _, identifier := range page.Items {
			if identifier.IdentifierID == getIdentifier.IdentifierID {
				return identifier, nil
			}
		}
	}
	return nil, fmt.Errorf("failed to retrieve pool identifier - [%s] from pool - [%s]", getIdentifier.IdentifierID, getIdentifier.PoolID)
}

// ListPoolsComponents lists all components in the connector management service.
func (s *ArkCmgrService) ListPoolsComponents() (<-chan *ArkCmgrPoolComponentPage, error) {
	s.Logger.Info("Listing pools components")
	return listCommonPools[cmgrmodels.ArkCmgrPoolComponent](
		s.Logger,
		s.client,
		"pools components",
		fmt.Sprintf(poolsComponentsURL),
		nil,
		map[string]string{
			"id": "component_id",
		},
	)
}

// ListPoolsComponentsBy lists components by the specified filter in the connector management service.
func (s *ArkCmgrService) ListPoolsComponentsBy(componentsFilters *cmgrmodels.ArkCmgrPoolComponentsFilter) (<-chan *ArkCmgrPoolComponentPage, error) {
	s.Logger.Info("Listing pools components by filter [%v]", componentsFilters)
	return listCommonPools[cmgrmodels.ArkCmgrPoolComponent](
		s.Logger,
		s.client,
		"pools components",
		fmt.Sprintf(poolsComponentsURL),
		&componentsFilters.ArkCmgrPoolsCommonFilter,
		map[string]string{
			"id": "component_id",
		},
	)
}

// PoolComponent retrieves a specific component by its ID from the connector management service.
func (s *ArkCmgrService) PoolComponent(getPoolComponent *cmgrmodels.ArkCmgrGetPoolComponent) (*cmgrmodels.ArkCmgrPoolComponent, error) {
	s.Logger.Info("Retrieving pool component [%s]", getPoolComponent.ComponentID)
	response, err := s.client.Get(context.Background(), fmt.Sprintf(poolComponentURL, getPoolComponent.PoolID, getPoolComponent.ComponentID), nil)
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
		return nil, fmt.Errorf("failed to retrieve pool component - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	poolComponentJSON, err := common.DeserializeJSONSnake(response.Body)
	if err != nil {
		return nil, err
	}
	poolComponentJSONMap := poolComponentJSON.(map[string]interface{})
	poolComponentJSONMap["component_id"] = poolComponentJSONMap["id"]
	var poolComponent cmgrmodels.ArkCmgrPoolComponent
	err = mapstructure.Decode(poolComponentJSONMap, &poolComponent)
	if err != nil {
		return nil, err
	}
	return &poolComponent, nil
}

// ServiceConfig returns the service configuration for the ArkCmgrService.
func (s *ArkCmgrService) ServiceConfig() services.ArkServiceConfig {
	return CmgrServiceConfig
}
