package dbsecrets

import (
	"context"
	"errors"
	"fmt"
	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/common"
	"github.com/cyberark/ark-sdk-golang/pkg/common/isp"
	dbsecretsmodels "github.com/cyberark/ark-sdk-golang/pkg/models/services/sia/secrets/db"
	dbworkspacemodels "github.com/cyberark/ark-sdk-golang/pkg/models/services/sia/workspaces/db"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	"github.com/mitchellh/mapstructure"
	"io"
	"net/http"
	"regexp"
)

const (
	secretsURL       = "/api/adb/secretsmgmt/secrets"
	secretURL        = "/api/adb/secretsmgmt/secrets/%s"
	enableSecretURL  = "/api/adb/secretsmgmt/secrets/%s/enable"
	disableSecretURL = "/api/adb/secretsmgmt/secrets/%s/disable"
)

// SIASecretsDBServiceConfig is the configuration for the SIA DB secrets service.
var SIASecretsDBServiceConfig = services.ArkServiceConfig{
	ServiceName:                "sia-secrets-db",
	RequiredAuthenticatorNames: []string{"isp"},
	OptionalAuthenticatorNames: []string{},
}

// ArkSIASecretsDBService is the service for managing db secrets.
type ArkSIASecretsDBService struct {
	services.ArkService
	*services.ArkBaseService
	ispAuth *auth.ArkISPAuth
	client  *isp.ArkISPServiceClient
}

// NewArkSIASecretsDBService creates a new instance of ArkSIASecretsDBService.
func NewArkSIASecretsDBService(authenticators ...auth.ArkAuth) (*ArkSIASecretsDBService, error) {
	secretsDBService := &ArkSIASecretsDBService{}
	var secretsDBServiceInterface services.ArkService = secretsDBService
	baseService, err := services.NewArkBaseService(secretsDBServiceInterface, authenticators...)
	if err != nil {
		return nil, err
	}
	ispBaseAuth, err := baseService.Authenticator("isp")
	if err != nil {
		return nil, err
	}
	ispAuth := ispBaseAuth.(*auth.ArkISPAuth)
	client, err := isp.FromISPAuth(ispAuth, "dpa", ".", "", secretsDBService.refreshSIAAuth)
	if err != nil {
		return nil, err
	}
	secretsDBService.client = client
	secretsDBService.ispAuth = ispAuth
	secretsDBService.ArkBaseService = baseService
	return secretsDBService, nil
}

func (s *ArkSIASecretsDBService) refreshSIAAuth(client *common.ArkClient) error {
	err := isp.RefreshClient(client, s.ispAuth)
	if err != nil {
		return err
	}
	return nil
}

func (s *ArkSIASecretsDBService) listSecretsWithFilters(secretType string, tags map[string]string) (*dbsecretsmodels.ArkSIADBSecretMetadataList, error) {
	params := make(map[string]string)
	if secretType != "" {
		params["secret_type"] = secretType
	}
	if tags != nil {
		for key, value := range tags {
			params[key] = value
		}
	}
	response, err := s.client.Get(context.Background(), secretsURL, params)
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
		return nil, fmt.Errorf("failed to list secrets - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	secretsJSON, err := common.DeserializeJSONSnake(response.Body)
	if err != nil {
		return nil, err
	}
	var secretsList dbsecretsmodels.ArkSIADBSecretMetadataList
	err = mapstructure.Decode(secretsJSON, &secretsList)

	if err != nil {
		s.Logger.Error("Failed to parse list secrets response [%v]", err)
		return nil, fmt.Errorf("failed to parse list secrets response: [%v]", err)
	}
	return &secretsList, nil
}

// AddSecret adds a new secret to the Ark SIA DB.
func (s *ArkSIASecretsDBService) AddSecret(addSecret *dbsecretsmodels.ArkSIADBAddSecret) (*dbsecretsmodels.ArkSIADBSecretMetadata, error) {
	if addSecret.StoreType == "" {
		storeType, ok := dbsecretsmodels.SecretTypeToStoreDict[addSecret.SecretType]
		if !ok {
			return nil, errors.New("invalid secret type")
		}
		addSecret.StoreType = storeType
	}
	addSecretJSON := map[string]interface{}{
		"secret_store": map[string]interface{}{
			"store_type": addSecret.StoreType,
		},
		"secret_name": addSecret.SecretName,
		"secret_type": addSecret.SecretType,
	}
	if addSecret.Description != "" {
		addSecretJSON["description"] = addSecret.Description
	}
	if addSecret.Purpose != "" {
		addSecretJSON["purpose"] = addSecret.Purpose
	}
	if addSecret.Tags != nil {
		addSecretJSON["tags"] = make([]dbworkspacemodels.ArkSIADBTag, len(addSecret.Tags))
		for key, value := range addSecret.Tags {
			addSecretJSON["tags"] = append(addSecretJSON["tags"].([]dbworkspacemodels.ArkSIADBTag), dbworkspacemodels.ArkSIADBTag{
				Key:   key,
				Value: value,
			})
		}
	}
	switch addSecret.SecretType {
	case dbsecretsmodels.UsernamePassword:
		if addSecret.Username == "" || addSecret.Password == "" {
			return nil, errors.New("username and password must be supplied for username_password type")
		}
		addSecretJSON["secret_data"] = map[string]interface{}{
			"username": addSecret.Username,
			"password": addSecret.Password,
		}
	case dbsecretsmodels.CyberArkPAM:
		if addSecret.PAMSafe == "" || addSecret.PAMAccountName == "" {
			return nil, errors.New("pam safe and pam account name must be supplied for pam type")
		}
		addSecretJSON["secret_link"] = map[string]interface{}{
			"safe":         addSecret.PAMSafe,
			"account_name": addSecret.PAMAccountName,
		}
	case dbsecretsmodels.IAMUser:
		if addSecret.IAMAccessKeyID == "" || addSecret.IAMSecretAccessKey == "" || addSecret.IAMAccount == "" || addSecret.IAMUsername == "" {
			return nil, errors.New("all IAM parameters must be supplied for iam_user type")
		}
		addSecretJSON["secret_data"] = map[string]interface{}{
			"account":           addSecret.IAMAccount,
			"username":          addSecret.IAMUsername,
			"access_key_id":     addSecret.IAMAccessKeyID,
			"secret_access_key": addSecret.IAMSecretAccessKey,
		}
	case dbsecretsmodels.AtlasAccessKeys:
		if addSecret.AtlasPublicKey == "" || addSecret.AtlasPrivateKey == "" {
			return nil, errors.New("public key and private key must be supplied for atlas type")
		}
		addSecretJSON["secret_data"] = map[string]interface{}{
			"public_key":  addSecret.AtlasPublicKey,
			"private_key": addSecret.AtlasPrivateKey,
		}
	default:
		return nil, fmt.Errorf("unsupported secret type: %s", addSecret.SecretType)
	}
	response, err := s.client.Post(context.Background(), secretsURL, addSecretJSON)
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
		return nil, fmt.Errorf("failed to add secret - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	secretJSON, err := common.DeserializeJSONSnake(response.Body)
	if err != nil {
		return nil, err
	}
	var secret dbsecretsmodels.ArkSIADBSecretMetadata
	err = mapstructure.Decode(secretJSON, &secret)
	if err != nil {
		return nil, err
	}
	return &secret, nil
}

// UpdateSecret updates an existing secret in the Ark SIA DB.
func (s *ArkSIASecretsDBService) UpdateSecret(updateSecret *dbsecretsmodels.ArkSIADBUpdateSecret) (*dbsecretsmodels.ArkSIADBSecretMetadata, error) {
	if updateSecret.SecretName != "" && updateSecret.SecretID == "" {
		secrets, err := s.ListSecretsBy(&dbsecretsmodels.ArkSIADBSecretsFilter{SecretName: updateSecret.SecretName})
		if err != nil || len(secrets.Secrets) == 0 {
			return nil, fmt.Errorf("failed to find secret by name: %v", err)
		}
		updateSecret.SecretID = secrets.Secrets[0].SecretID
	}
	s.Logger.Info(fmt.Sprintf("Updating existing db secret with id [%s]", updateSecret.SecretID))
	updateSecretMap := make(map[string]interface{})
	if updateSecret.NewSecretName != "" {
		updateSecretMap["secret_name"] = updateSecret.NewSecretName
	}
	if updateSecret.Description != "" {
		updateSecretMap["description"] = updateSecret.Description
	}
	if updateSecret.Purpose != "" {
		updateSecretMap["purpose"] = updateSecret.Purpose
	}
	if updateSecret.Tags != nil {
		updateSecretMap["tags"] = make([]dbworkspacemodels.ArkSIADBTag, len(updateSecret.Tags))
		for key, value := range updateSecret.Tags {
			updateSecretMap["tags"] = append(updateSecretMap["tags"].([]dbworkspacemodels.ArkSIADBTag), dbworkspacemodels.ArkSIADBTag{
				Key:   key,
				Value: value,
			})
		}
	}
	if updateSecret.PAMAccountName != "" || updateSecret.PAMSafe != "" {
		if updateSecret.PAMAccountName == "" || updateSecret.PAMSafe == "" {
			return nil, errors.New("both pam safe and pam account name must be supplied for pam secret")
		}
		updateSecretMap["secret_link"] = map[string]interface{}{
			"safe":         updateSecret.PAMSafe,
			"account_name": updateSecret.PAMAccountName,
		}
	}
	if updateSecret.Username != "" || updateSecret.Password != "" {
		if updateSecret.Username == "" || updateSecret.Password == "" {
			return nil, errors.New("both username and password must be supplied for username_password secret")
		}
		updateSecretMap["secret_data"] = map[string]interface{}{
			"username": updateSecret.Username,
			"password": updateSecret.Password,
		}
	}

	if updateSecret.IAMAccessKeyID != "" || updateSecret.IAMSecretAccessKey != "" || updateSecret.IAMAccount != "" || updateSecret.IAMUsername != "" {
		if updateSecret.IAMAccessKeyID == "" || updateSecret.IAMSecretAccessKey == "" || updateSecret.IAMAccount == "" || updateSecret.IAMUsername == "" {
			return nil, errors.New("all IAM parameters must be supplied for iam_user secret")
		}
		updateSecretMap["secret_data"] = map[string]interface{}{
			"account":           updateSecret.IAMAccount,
			"username":          updateSecret.IAMUsername,
			"access_key_id":     updateSecret.IAMAccessKeyID,
			"secret_access_key": updateSecret.IAMSecretAccessKey,
		}
	}

	if updateSecret.AtlasPublicKey != "" || updateSecret.AtlasPrivateKey != "" {
		if updateSecret.AtlasPublicKey == "" || updateSecret.AtlasPrivateKey == "" {
			return nil, errors.New("both public key and private key must be supplied for atlas secret")
		}
		updateSecretMap["secret_data"] = map[string]interface{}{
			"public_key":  updateSecret.AtlasPublicKey,
			"private_key": updateSecret.AtlasPrivateKey,
		}
	}
	response, err := s.client.Patch(context.Background(), fmt.Sprintf(secretURL, updateSecret.SecretID), updateSecretMap)
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
		return nil, fmt.Errorf("failed to update secret - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	secretJSON, err := common.DeserializeJSONSnake(response.Body)
	if err != nil {
		return nil, err
	}
	var secret dbsecretsmodels.ArkSIADBSecretMetadata
	err = mapstructure.Decode(secretJSON, &secret)
	if err != nil {
		return nil, err
	}

	return &secret, nil
}

// DeleteSecret deletes a secret from the Ark SIA DB.
func (s *ArkSIASecretsDBService) DeleteSecret(deleteSecret *dbsecretsmodels.ArkSIADBDeleteSecret) error {
	if deleteSecret.SecretName != "" && deleteSecret.SecretID == "" {
		secrets, err := s.ListSecretsBy(&dbsecretsmodels.ArkSIADBSecretsFilter{SecretName: deleteSecret.SecretName})
		if err != nil || len(secrets.Secrets) == 0 {
			return fmt.Errorf("failed to find secret by name: %v", err)
		}
		deleteSecret.SecretID = secrets.Secrets[0].SecretID
	}
	s.Logger.Info(fmt.Sprintf("Deleting db secret by id [%s]", deleteSecret.SecretID))
	response, err := s.client.Delete(context.Background(), fmt.Sprintf(secretURL, deleteSecret.SecretID), nil)
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
		return fmt.Errorf("failed to delete db secret [%s] - [%d]", common.SerializeResponseToJSON(response.Body), response.StatusCode)
	}
	return nil
}

// ListSecrets lists all secrets in the Ark SIA DB.
func (s *ArkSIASecretsDBService) ListSecrets() (*dbsecretsmodels.ArkSIADBSecretMetadataList, error) {
	return s.listSecretsWithFilters("", nil)
}

// ListSecretsBy lists secrets in the Ark SIA DB by the given filter.
func (s *ArkSIASecretsDBService) ListSecretsBy(filter *dbsecretsmodels.ArkSIADBSecretsFilter) (*dbsecretsmodels.ArkSIADBSecretMetadataList, error) {
	secrets, err := s.listSecretsWithFilters(filter.SecretType, filter.Tags)
	if err != nil {
		return nil, err
	}
	if filter.StoreType != "" {
		var filteredSecrets []dbsecretsmodels.ArkSIADBSecretMetadata
		for _, secret := range secrets.Secrets {
			if secret.SecretStore.StoreType == filter.StoreType {
				filteredSecrets = append(filteredSecrets, secret)
			}
		}
		secrets.Secrets = filteredSecrets
	}
	if filter.SecretName != "" {
		var filteredSecrets []dbsecretsmodels.ArkSIADBSecretMetadata
		for _, secret := range secrets.Secrets {
			if secret.SecretName != "" {
				matched, _ := regexp.MatchString(filter.SecretName, secret.SecretName)
				if matched {
					filteredSecrets = append(filteredSecrets, secret)
				}
			}
		}
		secrets.Secrets = filteredSecrets
	}
	if filter.IsActive {
		var filteredSecrets []dbsecretsmodels.ArkSIADBSecretMetadata
		for _, secret := range secrets.Secrets {
			if secret.IsActive == filter.IsActive {
				filteredSecrets = append(filteredSecrets, secret)
			}
		}
		secrets.Secrets = filteredSecrets
	}
	secrets.TotalCount = len(secrets.Secrets)
	return secrets, nil
}

// EnableSecret enables a secret in the Ark SIA DB.
func (s *ArkSIASecretsDBService) EnableSecret(enableSecret *dbsecretsmodels.ArkSIADBEnableSecret) error {
	if enableSecret.SecretName != "" && enableSecret.SecretID == "" {
		secrets, err := s.ListSecretsBy(&dbsecretsmodels.ArkSIADBSecretsFilter{SecretName: enableSecret.SecretName})
		if err != nil || len(secrets.Secrets) == 0 {
			return fmt.Errorf("failed to find secret by name: %v", err)
		}
		enableSecret.SecretID = secrets.Secrets[0].SecretID
	}
	s.Logger.Info(fmt.Sprintf("Enabling db secret by id [%s]", enableSecret.SecretID))
	response, err := s.client.Post(context.Background(), fmt.Sprintf(enableSecretURL, enableSecret.SecretID), nil)
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
		return fmt.Errorf("failed to enable db secret [%s] - [%d]", common.SerializeResponseToJSON(response.Body), response.StatusCode)
	}
	return nil
}

// DisableSecret disables a secret in the Ark SIA DB.
func (s *ArkSIASecretsDBService) DisableSecret(enableSecret *dbsecretsmodels.ArkSIADBDisableSecret) error {
	if enableSecret.SecretName != "" && enableSecret.SecretID == "" {
		secrets, err := s.ListSecretsBy(&dbsecretsmodels.ArkSIADBSecretsFilter{SecretName: enableSecret.SecretName})
		if err != nil || len(secrets.Secrets) == 0 {
			return fmt.Errorf("failed to find secret by name: %v", err)
		}
		enableSecret.SecretID = secrets.Secrets[0].SecretID
	}
	s.Logger.Info(fmt.Sprintf("Disabling db secret by id [%s]", enableSecret.SecretID))
	response, err := s.client.Post(context.Background(), fmt.Sprintf(disableSecretURL, enableSecret.SecretID), nil)
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
		return fmt.Errorf("failed to disable db secret [%s] - [%d]", common.SerializeResponseToJSON(response.Body), response.StatusCode)
	}
	return nil
}

// Secret retrieves a secret from the Ark SIA DB by its ID.
func (s *ArkSIASecretsDBService) Secret(getSecret *dbsecretsmodels.ArkSIADBGetSecret) (*dbsecretsmodels.ArkSIADBSecretMetadata, error) {
	if getSecret.SecretName != "" && getSecret.SecretID == "" {
		secrets, err := s.ListSecretsBy(&dbsecretsmodels.ArkSIADBSecretsFilter{SecretName: getSecret.SecretName})
		if err != nil || len(secrets.Secrets) == 0 {
			return nil, fmt.Errorf("failed to find secret by name: %v", err)
		}
		getSecret.SecretID = secrets.Secrets[0].SecretID
	}
	s.Logger.Info("Retrieving db secret by id [%s]", getSecret.SecretID)
	response, err := s.client.Get(context.Background(), fmt.Sprintf(secretURL, getSecret.SecretID), nil)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.GlobalLogger.Warning("Error closing response body")
		}
	}(response.Body)

	// Check response status
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to retrieve db secret [%s] - [%d]", common.SerializeResponseToJSON(response.Body), response.StatusCode)
	}
	secretJSON, err := common.DeserializeJSONSnake(response.Body)
	if err != nil {
		s.Logger.Error("Failed to parse db secret response [%v]", err)
		return nil, fmt.Errorf("failed to parse db secret response: %w", err)
	}
	var secret dbsecretsmodels.ArkSIADBSecretMetadata
	err = mapstructure.Decode(secretJSON, &secret)
	if err != nil {
		return nil, err
	}
	return &secret, nil
}

// SecretsStats retrieves the statistics of secrets in the Ark SIA DB.
func (s *ArkSIASecretsDBService) SecretsStats() (*dbsecretsmodels.ArkSIADBSecretsStats, error) {
	s.Logger.Info("Calculating secrets statistics")
	secretsList, err := s.ListSecrets()
	if err != nil {
		return nil, err
	}
	secretsStats := &dbsecretsmodels.ArkSIADBSecretsStats{
		SecretsCountBySecretType: make(map[string]int),
		SecretsCountByStoreType:  make(map[string]int),
	}
	secretsStats.SecretsCount = len(secretsList.Secrets)
	for _, secret := range secretsList.Secrets {
		if secret.IsActive {
			secretsStats.ActiveSecretsCount++
		} else {
			secretsStats.InactiveSecretsCount++
		}
		if secret.SecretType != "" {
			if _, ok := secretsStats.SecretsCountBySecretType[secret.SecretType]; !ok {
				secretsStats.SecretsCountBySecretType[secret.SecretType] = 0
			}
			secretsStats.SecretsCountBySecretType[secret.SecretType]++
		}
		if secret.SecretStore.StoreType != "" {
			if _, ok := secretsStats.SecretsCountByStoreType[secret.SecretStore.StoreType]; !ok {
				secretsStats.SecretsCountByStoreType[secret.SecretStore.StoreType] = 0
			}
			secretsStats.SecretsCountByStoreType[secret.SecretStore.StoreType]++
		}
	}
	return secretsStats, nil
}

// ServiceConfig returns the service configuration for the ArkSIASecretsVMService.
func (s *ArkSIASecretsDBService) ServiceConfig() services.ArkServiceConfig {
	return SIASecretsDBServiceConfig
}
