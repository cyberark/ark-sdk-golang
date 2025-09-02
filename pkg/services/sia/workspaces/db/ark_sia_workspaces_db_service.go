package db

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"slices"

	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/common"
	"github.com/cyberark/ark-sdk-golang/pkg/common/isp"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	workspacesdbmodels "github.com/cyberark/ark-sdk-golang/pkg/services/sia/workspaces/db/models"
	"github.com/mitchellh/mapstructure"
)

const (
	resourcesURL = "/api/adb/resources"
	resourceURL  = "/api/adb/resources/%d"
)

// SIADBWorkspaceServiceConfig is the configuration for the SIA db workspace service.
var SIADBWorkspaceServiceConfig = services.ArkServiceConfig{
	ServiceName:                "sia-workspaces-db",
	RequiredAuthenticatorNames: []string{"isp"},
	OptionalAuthenticatorNames: []string{},
}

// ArkSIAWorkspacesDBService is the service for managing databases in a workspace.
type ArkSIAWorkspacesDBService struct {
	services.ArkService
	*services.ArkBaseService
	ispAuth *auth.ArkISPAuth
	client  *isp.ArkISPServiceClient
}

// NewArkSIAWorkspacesDBService creates a new instance of ArkSIAWorkspacesDBService.
func NewArkSIAWorkspacesDBService(authenticators ...auth.ArkAuth) (*ArkSIAWorkspacesDBService, error) {
	dbService := &ArkSIAWorkspacesDBService{}
	var dbServiceInterface services.ArkService = dbService
	baseService, err := services.NewArkBaseService(dbServiceInterface, authenticators...)
	if err != nil {
		return nil, err
	}
	ispBaseAuth, err := baseService.Authenticator("isp")
	if err != nil {
		return nil, err
	}
	ispAuth := ispBaseAuth.(*auth.ArkISPAuth)
	client, err := isp.FromISPAuth(ispAuth, "dpa", ".", "", dbService.refreshSIAAuth)
	if err != nil {
		return nil, err
	}
	dbService.client = client
	dbService.ispAuth = ispAuth
	dbService.ArkBaseService = baseService
	return dbService, nil
}

func (s *ArkSIAWorkspacesDBService) parseDatabaseTagsIntoMap(databaseJSONMap map[string]interface{}) {
	if tags, ok := databaseJSONMap["tags"].([]interface{}); ok {
		parsedTags := make(map[string]string)
		for _, tag := range tags {
			if tagMap, ok := tag.(map[string]interface{}); ok {
				key, keyOk := tagMap["key"].(string)
				value, valueOk := tagMap["value"].(string)
				if keyOk && valueOk {
					parsedTags[key] = value
				}
			}
		}
		databaseJSONMap["tags"] = parsedTags
	}
}

func (s *ArkSIAWorkspacesDBService) refreshSIAAuth(client *common.ArkClient) error {
	err := isp.RefreshClient(client, s.ispAuth)
	if err != nil {
		return err
	}
	return nil
}

func (s *ArkSIAWorkspacesDBService) listDatabasesWithFilters(providerFamily string, tags []workspacesdbmodels.ArkSIADBTag) (*workspacesdbmodels.ArkSIADBDatabaseInfoList, error) {
	params := make(map[string]string)
	if providerFamily != "" {
		params["provider-family"] = providerFamily
	}
	if tags != nil {
		for _, tag := range tags {
			params[fmt.Sprintf("key.%s", tag.Key)] = tag.Value
		}
	}
	response, err := s.client.Get(context.Background(), resourcesURL, params)
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
		return nil, fmt.Errorf("failed to list databases with filters - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}

	databasesJSON, err := common.DeserializeJSONSnake(response.Body)
	if err != nil {
		return nil, err
	}
	var databases workspacesdbmodels.ArkSIADBDatabaseInfoList
	err = mapstructure.Decode(databasesJSON, &databases)
	return &databases, nil
}

// AddDatabase adds a new database to the SIA workspace.
func (s *ArkSIAWorkspacesDBService) AddDatabase(addDatabase *workspacesdbmodels.ArkSIADBAddDatabase) (*workspacesdbmodels.ArkSIADBDatabase, error) {
	s.Logger.Info("Adding database [%s]", addDatabase.Name)
	// Validate ProviderEngine
	if !slices.Contains(workspacesdbmodels.DatabaseEngineTypes, addDatabase.ProviderEngine) {
		return nil, fmt.Errorf("invalid provider engine: %s", addDatabase.ProviderEngine)
	}
	// Set default port if not provided
	if addDatabase.Port == 0 {
		family, ok := workspacesdbmodels.DatabasesEnginesToFamily[addDatabase.ProviderEngine]
		if !ok {
			return nil, fmt.Errorf("unknown provider engine: %s", addDatabase.ProviderEngine)
		}
		addDatabase.Port = workspacesdbmodels.DatabaseFamiliesDefaultPorts[family]
	}
	if addDatabase.Services == nil {
		addDatabase.Services = []string{}
	}
	var addDatabaseJSON map[string]interface{}
	err := mapstructure.Decode(addDatabase, &addDatabaseJSON)
	if err != nil {
		return nil, err
	}
	if addDatabase.Tags != nil {
		addDatabaseJSON["tags"] = make([]workspacesdbmodels.ArkSIADBTag, len(addDatabase.Tags))
		idx := 0
		for key, value := range addDatabase.Tags {
			if key == "" {
				continue
			}
			addDatabaseJSON["tags"].([]workspacesdbmodels.ArkSIADBTag)[idx] = workspacesdbmodels.ArkSIADBTag{
				Key:   key,
				Value: value,
			}
			idx++
		}
	}
	response, err := s.client.Post(context.Background(), resourcesURL, addDatabaseJSON)
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
		return nil, fmt.Errorf("failed to database - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	databaseJSON, err := common.DeserializeJSONSnake(response.Body)
	if err != nil {
		return nil, err
	}
	databaseJSONMap := databaseJSON.(map[string]interface{})
	databaseID, ok := databaseJSONMap["target_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("missing target_id in response")
	}
	getDatabase := &workspacesdbmodels.ArkSIADBGetDatabase{ID: int(databaseID)}
	return s.Database(getDatabase)
}

// DeleteDatabase deletes a database.
func (s *ArkSIAWorkspacesDBService) DeleteDatabase(deleteDatabase *workspacesdbmodels.ArkSIADBDeleteDatabase) error {
	if deleteDatabase.Name != "" && deleteDatabase.ID == 0 {
		databases, err := s.ListDatabasesBy(&workspacesdbmodels.ArkSIADBDatabasesFilter{Name: deleteDatabase.Name})
		if err != nil {
			return fmt.Errorf("failed to fetch database ID by name: %w", err)
		}
		if len(databases.Items) == 0 || len(databases.Items) != 1 {
			return fmt.Errorf("no database found with name: %s", deleteDatabase.Name)
		}
		deleteDatabase.ID = databases.Items[0].ID
	}
	s.Logger.Info("Deleting database [%d]", deleteDatabase.ID)
	response, err := s.client.Delete(context.Background(), fmt.Sprintf(resourceURL, deleteDatabase.ID), nil)
	if err != nil {
		return fmt.Errorf("failed to delete database: %w", err)
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			s.Logger.Warning("Error closing response body")
		}
	}(response.Body)

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete database - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}

	return nil
}

// UpdateDatabase updates a database.
func (s *ArkSIAWorkspacesDBService) UpdateDatabase(updateDatabase *workspacesdbmodels.ArkSIADBUpdateDatabase) (*workspacesdbmodels.ArkSIADBDatabase, error) {
	if updateDatabase.Name != "" && updateDatabase.ID == 0 {
		databases, err := s.ListDatabasesBy(&workspacesdbmodels.ArkSIADBDatabasesFilter{Name: updateDatabase.Name})
		if err != nil {
			return nil, fmt.Errorf("failed to fetch database ID by name: %w", err)
		}
		if len(databases.Items) == 0 || len(databases.Items) != 1 {
			return nil, fmt.Errorf("failed to update database - name [%s] not found", updateDatabase.Name)
		}
		updateDatabase.ID = databases.Items[0].ID
	}
	// Validate ProviderEngine
	if updateDatabase.ProviderEngine != "" && !slices.Contains(workspacesdbmodels.DatabaseEngineTypes, updateDatabase.ProviderEngine) {
		return nil, fmt.Errorf("invalid provider engine: %s", updateDatabase.ProviderEngine)
	}
	existingDatabase, err := s.Database(&workspacesdbmodels.ArkSIADBGetDatabase{ID: updateDatabase.ID, Name: updateDatabase.Name})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve existing database: %w", err)
	}

	// Merge the existing database details with the update request
	mergedDatabase := make(map[string]interface{})
	existingDatabaseMap := make(map[string]interface{})
	updateDatabaseMap := make(map[string]interface{})

	// As update database it PUT, we need first to fetch the existing database,
	// and merge the update params with the existing database, so that all fields will be sent.
	if err := mapstructure.Decode(existingDatabase, &existingDatabaseMap); err != nil {
		return nil, fmt.Errorf("failed to decode existing database: %w", err)
	}
	if err := mapstructure.Decode(updateDatabase, &updateDatabaseMap); err != nil {
		return nil, fmt.Errorf("failed to decode update database payload: %w", err)
	}
	// Merge the maps
	for key, value := range existingDatabaseMap {
		mergedDatabase[key] = value
	}
	for key, value := range updateDatabaseMap {
		mergedDatabase[key] = value
	}

	// Remove unnecessary fields and handle renaming
	delete(mergedDatabase, "name")
	delete(mergedDatabase, "new_name")
	if updateDatabase.NewName != "" {
		mergedDatabase["name"] = updateDatabase.NewName
	} else if updateDatabase.Name != "" {
		mergedDatabase["name"] = updateDatabase.Name
	} else {
		mergedDatabase["name"] = existingDatabase.Name
	}

	// Handling configured auth method
	delete(mergedDatabase, "configured_auth_method")
	if updateDatabase.ConfiguredAuthMethodType == "" {
		mergedDatabase["configured_auth_method_type"] = existingDatabase.ConfiguredAuthMethod.DatabaseAuthMethod.AuthMethod.AuthMethodType
	}

	// Handling provider engine
	delete(mergedDatabase, "provider_details")
	if updateDatabase.ProviderEngine == "" {
		mergedDatabase["provider_engine"] = existingDatabase.ProviderDetails.Engine
	}

	if updateDatabase.Tags != nil {
		mergedDatabase["tags"] = make([]workspacesdbmodels.ArkSIADBTag, len(updateDatabase.Tags))
		idx := 0
		for key, value := range updateDatabase.Tags {
			if key == "" {
				continue
			}
			mergedDatabase["tags"].([]workspacesdbmodels.ArkSIADBTag)[idx] = workspacesdbmodels.ArkSIADBTag{
				Key:   key,
				Value: value,
			}
			idx++
		}
	}

	s.Logger.Info("Updating database [%d]", updateDatabase.ID)
	response, err := s.client.Put(context.Background(), fmt.Sprintf(resourceURL, updateDatabase.ID), mergedDatabase)
	if err != nil {
		return nil, fmt.Errorf("failed to update database: %w", err)
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			s.Logger.Warning("Error closing response body")
		}
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to update database - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	return s.Database(&workspacesdbmodels.ArkSIADBGetDatabase{ID: updateDatabase.ID})
}

// Database retrieves a database by id or name.
func (s *ArkSIAWorkspacesDBService) Database(getDatabase *workspacesdbmodels.ArkSIADBGetDatabase) (*workspacesdbmodels.ArkSIADBDatabase, error) {
	// If Name is provided but ID is not, fetch the ID by filtering databases
	if getDatabase.Name != "" && getDatabase.ID == 0 {
		filter := &workspacesdbmodels.ArkSIADBDatabasesFilter{Name: getDatabase.Name}
		databases, err := s.ListDatabasesBy(filter)
		if err != nil {
			return nil, fmt.Errorf("failed to list databases: %w", err)
		}
		if len(databases.Items) == 0 || len(databases.Items) != 1 {
			return nil, fmt.Errorf("failed to get database - name [%s] not found", getDatabase.Name)
		}
		getDatabase.ID = databases.Items[0].ID
	}
	s.Logger.Info("Getting database [%d]", getDatabase.ID)
	response, err := s.client.Get(context.Background(), fmt.Sprintf(resourceURL, getDatabase.ID), nil)
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
		return nil, fmt.Errorf("failed to get database - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}

	databaseJSON, err := common.DeserializeJSONSnake(response.Body)
	if err != nil {
		return nil, err
	}
	databaseJSONMap := databaseJSON.(map[string]interface{})
	s.parseDatabaseTagsIntoMap(databaseJSONMap)
	var database workspacesdbmodels.ArkSIADBDatabase
	err = mapstructure.Decode(databaseJSONMap, &database)
	return &database, nil
}

// ListDatabases lists all databases.
func (s *ArkSIAWorkspacesDBService) ListDatabases() (*workspacesdbmodels.ArkSIADBDatabaseInfoList, error) {
	s.Logger.Info("Listing all databases")
	return s.listDatabasesWithFilters("", nil)
}

// ListDatabasesBy filters databases by the given filters.
func (s *ArkSIAWorkspacesDBService) ListDatabasesBy(databasesFilter *workspacesdbmodels.ArkSIADBDatabasesFilter) (*workspacesdbmodels.ArkSIADBDatabaseInfoList, error) {
	if databasesFilter.ProviderEngine != "" && !slices.Contains(workspacesdbmodels.DatabaseEngineTypes, databasesFilter.ProviderEngine) {
		return nil, fmt.Errorf("invalid provider engine: %s", databasesFilter.ProviderEngine)
	}
	s.Logger.Info("Listing databases by filters [%+v]", databasesFilter)
	databases, err := s.listDatabasesWithFilters(databasesFilter.ProviderFamily, databasesFilter.Tags)
	if err != nil {
		return nil, fmt.Errorf("failed to list databases with filters: %w", err)
	}
	var filteredItems []workspacesdbmodels.ArkSIADBDatabaseInfo
	for _, database := range databases.Items {
		if databasesFilter.Name != "" {
			matched, err := regexp.MatchString(databasesFilter.Name, database.Name)
			if err != nil || !matched {
				continue
			}
		}
		if databasesFilter.ProviderEngine != "" && database.ProviderInfo.Engine != databasesFilter.ProviderEngine {
			continue
		}
		if databasesFilter.ProviderFamily != "" && database.ProviderInfo.Family != databasesFilter.ProviderFamily {
			continue
		}
		if databasesFilter.ProviderWorkspace != "" && database.ProviderInfo.Workspace != databasesFilter.ProviderWorkspace {
			continue
		}
		if len(databasesFilter.AuthMethods) > 0 {
			matchesAuthMethod := false
			for _, authMethod := range databasesFilter.AuthMethods {
				if database.ConfiguredAuthMethodType == authMethod {
					matchesAuthMethod = true
					break
				}
			}
			if !matchesAuthMethod {
				continue
			}
		}
		if databasesFilter.DBWarningsFilter != "" {
			if (databasesFilter.DBWarningsFilter == workspacesdbmodels.AnyError || databasesFilter.DBWarningsFilter == workspacesdbmodels.NoCertificates) && database.Certificate == "" {
				continue
			}
			if (databasesFilter.DBWarningsFilter == workspacesdbmodels.AnyError || databasesFilter.DBWarningsFilter == workspacesdbmodels.NoSecrets) && database.SecretID == "" {
				continue
			}
		}
		// Add to filtered items if all conditions are met
		filteredItems = append(filteredItems, database)
	}
	databases.Items = filteredItems
	databases.TotalCount = len(filteredItems)
	return databases, nil
}

// DatabasesStats calculates statistics about databases.
func (s *ArkSIAWorkspacesDBService) DatabasesStats() (*workspacesdbmodels.ArkSIADBDatabasesStats, error) {
	s.Logger.Info("Calculating databases stats")
	databases, err := s.ListDatabases()
	if err != nil {
		return nil, fmt.Errorf("failed to list databases: %w", err)
	}
	// Initialize the stats object
	databasesStats := &workspacesdbmodels.ArkSIADBDatabasesStats{
		DatabasesCount:             len(databases.Items),
		DatabasesCountByEngine:     make(map[string]int),
		DatabasesCountByWorkspace:  make(map[string]int),
		DatabasesCountByFamily:     make(map[string]int),
		DatabasesCountByAuthMethod: make(map[string]int),
		DatabasesCountByWarning:    make(map[string]int),
	}
	// Calculate databases per engine
	for _, database := range databases.Items {
		engine := database.ProviderInfo.Engine
		databasesStats.DatabasesCountByEngine[engine]++
	}
	// Calculate databases per workspace
	for _, database := range databases.Items {
		workspace := database.ProviderInfo.Workspace
		databasesStats.DatabasesCountByWorkspace[workspace]++
	}
	// Calculate databases per family
	for _, database := range databases.Items {
		family := database.ProviderInfo.Family
		databasesStats.DatabasesCountByFamily[family]++
	}
	// Calculate databases per auth method
	for _, database := range databases.Items {
		authMethod := database.ConfiguredAuthMethodType
		databasesStats.DatabasesCountByAuthMethod[authMethod]++
	}
	// Calculate databases per warning
	for _, database := range databases.Items {
		if database.Certificate == "" {
			databasesStats.DatabasesCountByWarning[workspacesdbmodels.NoCertificates]++
		}
		if database.SecretID == "" {
			databasesStats.DatabasesCountByWarning[workspacesdbmodels.NoSecrets]++
		}
	}
	return databasesStats, nil
}

// ListEngineTypes returns all possible database engine types.
func (s *ArkSIAWorkspacesDBService) ListEngineTypes() []string {
	return workspacesdbmodels.DatabaseEngineTypes
}

// ListFamilyTypes returns all possible database family types.
func (s *ArkSIAWorkspacesDBService) ListFamilyTypes() []string {
	return workspacesdbmodels.DatabaseFamilyTypes
}

// ServiceConfig returns the service configuration for the ArkSIATargetSetsWorkspaceService.
func (s *ArkSIAWorkspacesDBService) ServiceConfig() services.ArkServiceConfig {
	return SIADBWorkspaceServiceConfig
}
