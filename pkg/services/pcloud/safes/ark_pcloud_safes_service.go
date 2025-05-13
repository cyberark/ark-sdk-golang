package vmsecrets

import (
	"context"
	"fmt"
	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/common"
	"github.com/cyberark/ark-sdk-golang/pkg/common/isp"
	vmsecretsmodels "github.com/cyberark/ark-sdk-golang/pkg/models/services/sia/secrets/vm"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	"github.com/mitchellh/mapstructure"
	"io"
	"net/http"
	"regexp"
	"slices"
)

const (
	secretsURL = "/api/secrets"
	secretURL  = "/api/secrets/%s"
)

// SIASecretsVMServiceConfig is the configuration for the SIA VM secrets service.
var SIASecretsVMServiceConfig = services.ArkServiceConfig{
	ServiceName:                "sia-secrets-vm",
	RequiredAuthenticatorNames: []string{"isp"},
	OptionalAuthenticatorNames: []string{},
}

// ArkSIASecretsVMService is the service for managing vm secrets.
type ArkSIASecretsVMService struct {
	services.ArkService
	*services.ArkBaseService
	ispAuth *auth.ArkISPAuth
	client  *isp.ArkISPServiceClient
}

// NewArkSIASecretsVMService creates a new instance of ArkSIASecretsVMService.
func NewArkSIASecretsVMService(authenticators ...auth.ArkAuth) (*ArkSIASecretsVMService, error) {
	secretsVMService := &ArkSIASecretsVMService{}
	var secretsVMServiceInterface services.ArkService = secretsVMService
	baseService, err := services.NewArkBaseService(secretsVMServiceInterface, authenticators...)
	if err != nil {
		return nil, err
	}
	ispBaseAuth, err := baseService.Authenticator("isp")
	if err != nil {
		return nil, err
	}
	ispAuth := ispBaseAuth.(*auth.ArkISPAuth)
	client, err := isp.FromISPAuth(ispAuth, "dpa", ".", "", secretsVMService.refreshSIAAuth)
	if err != nil {
		return nil, err
	}
	secretsVMService.client = client
	secretsVMService.ispAuth = ispAuth
	secretsVMService.ArkBaseService = baseService
	return secretsVMService, nil
}

func (s *ArkSIASecretsVMService) refreshSIAAuth(client *common.ArkClient) error {
	err := isp.RefreshClient(client, s.ispAuth)
	if err != nil {
		return err
	}
	return nil
}

// AddSecret adds a new secret to the SIA VM secrets service.
func (s *ArkSIASecretsVMService) AddSecret(addSecret *vmsecretsmodels.ArkSIAVMAddSecret) (*vmsecretsmodels.ArkSIAVMSecret, error) {
	s.Logger.Info("Adding new vm secret")
	addSecretJSON := map[string]interface{}{
		"secret_name": addSecret.SecretName,
		"secret_type": addSecret.SecretType,
		"secret": map[string]interface{}{
			"tenant_encrypted": false,
		},
		"is_active":      !addSecret.IsDisabled,
		"secret_details": addSecret.SecretDetails,
	}
	if addSecret.SecretType == "ProvisionerUser" {
		if addSecret.ProvisionerUsername == "" || addSecret.ProvisionerPassword == "" {
			return nil, fmt.Errorf("provisioner username and password are required for ProvisionerUser secret type")
		}
		addSecretJSON["secret"].(map[string]interface{})["secret_data"] = map[string]interface{}{
			"username": addSecret.ProvisionerUsername,
			"password": addSecret.ProvisionerPassword,
		}
	} else if addSecret.SecretType == "PCloudAccount" {
		if addSecret.PCloudAccountSafe == "" || addSecret.PCloudAccountName == "" {
			return nil, fmt.Errorf("pcloud account safe and name are required for PCloudAccount secret type")
		}
		addSecretJSON["secret"].(map[string]interface{})["secret_data"] = map[string]interface{}{
			"safe":         addSecret.PCloudAccountSafe,
			"account_name": addSecret.PCloudAccountName,
		}
	} else {
		return nil, fmt.Errorf("invalid secret type: %s", addSecret.SecretType)
	}
	if addSecret.SecretDetails != nil {
		addSecretJSON["secret_details"] = addSecret.SecretDetails
	} else {
		addSecretJSON["secret_details"] = map[string]interface{}{}
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
	var secret vmsecretsmodels.ArkSIAVMSecret
	err = mapstructure.Decode(secretJSON, &secret)
	if err != nil {
		return nil, err
	}
	return &secret, nil
}

// ChangeSecret changes an existing secret in the SIA VM secrets service.
func (s *ArkSIASecretsVMService) ChangeSecret(changeSecret *vmsecretsmodels.ArkSIAVMChangeSecret) (*vmsecretsmodels.ArkSIAVMSecret, error) {
	s.Logger.Info("Changing existing vm secret with id [%s]", changeSecret.SecretID)
	changeSecretJSON := map[string]interface{}{
		"is_active": !changeSecret.IsDisabled,
	}
	if changeSecret.ProvisionerUsername != "" && changeSecret.ProvisionerPassword != "" {
		changeSecretJSON["secret"] = map[string]interface{}{
			"secret_data": map[string]interface{}{
				"username": changeSecret.ProvisionerUsername,
				"password": changeSecret.ProvisionerPassword,
			},
		}
	}
	if changeSecret.PCloudAccountSafe != "" && changeSecret.PCloudAccountName != "" {
		changeSecretJSON["secret"] = map[string]interface{}{
			"secret_data": map[string]interface{}{
				"safe":         changeSecret.PCloudAccountSafe,
				"account_name": changeSecret.PCloudAccountName,
			},
		}
	}
	if changeSecret.SecretName != "" {
		changeSecretJSON["secret_name"] = changeSecret.SecretName
	}
	if changeSecret.SecretDetails != nil {
		changeSecretJSON["secret_details"] = changeSecret.SecretDetails
	}
	response, err := s.client.Post(context.Background(), fmt.Sprintf(secretURL, changeSecret.SecretID), changeSecretJSON)
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
		return nil, fmt.Errorf("failed to change secret - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	secretJSON, err := common.DeserializeJSONSnake(response.Body)
	if err != nil {
		return nil, err
	}
	var secret vmsecretsmodels.ArkSIAVMSecret
	err = mapstructure.Decode(secretJSON, &secret)
	if err != nil {
		return nil, err
	}
	return &secret, nil
}

// DeleteSecret deletes a secret from the SIA VM secrets service.
func (s *ArkSIASecretsVMService) DeleteSecret(deleteSecret *vmsecretsmodels.ArkSIAVMDeleteSecret) error {
	s.Logger.Info("Deleting secret [%s]", deleteSecret.SecretID)
	response, err := s.client.Delete(context.Background(), fmt.Sprintf(secretURL, deleteSecret.SecretID), nil)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete secret - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	return nil
}

func (s *ArkSIASecretsVMService) listSecretsWithFilter(secretType string, secretDetails map[string]interface{}) ([]*vmsecretsmodels.ArkSIAVMSecret, error) {
	filterJSON := map[string]string{}
	if secretType != "" {
		filterJSON["secret_type"] = secretType
	}
	if secretDetails != nil {
		for key, value := range secretDetails {
			if value != nil {
				filterJSON[key] = fmt.Sprintf("%v", value)
			}
		}
	}
	response, err := s.client.Get(context.Background(), secretsURL, nil)
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
	secretsResponseJSON, err := common.DeserializeJSONSnake(response.Body)
	if err != nil {
		return nil, err
	}
	var secrets []*vmsecretsmodels.ArkSIAVMSecret
	err = mapstructure.Decode(secretsResponseJSON, &secrets)
	if err != nil {
		return nil, err
	}
	return secrets, nil
}

// ListSecrets lists all secrets in the SIA VM secrets service.
func (s *ArkSIASecretsVMService) ListSecrets() ([]*vmsecretsmodels.ArkSIAVMSecret, error) {
	s.Logger.Info("Listing all secrets")
	return s.listSecretsWithFilter("", nil)
}

// ListSecretsBy lists secrets in the SIA VM secrets service by filter.
func (s *ArkSIASecretsVMService) ListSecretsBy(filter *vmsecretsmodels.ArkSIAVMSecretsFilter) ([]*vmsecretsmodels.ArkSIAVMSecret, error) {
	s.Logger.Info("Listing secrets by filters [%v]", filter)
	secretType := ""
	if filter.SecretTypes != nil && len(filter.SecretTypes) > 0 {
		secretType = filter.SecretTypes[0]
	}
	secrets, err := s.listSecretsWithFilter(secretType, filter.SecretDetails)
	if err != nil {
		return nil, err
	}
	var filteredSecrets []*vmsecretsmodels.ArkSIAVMSecret
	for _, secret := range secrets {
		if secret.IsActive == filter.IsActive {
			filteredSecrets = append(filteredSecrets, secret)
		}
	}
	secrets = filteredSecrets
	if filter.Name != "" {
		filteredSecrets = []*vmsecretsmodels.ArkSIAVMSecret{}
		for _, secret := range secrets {
			if match, err := regexp.MatchString(filter.Name, secret.SecretName); err == nil && match {
				filteredSecrets = append(filteredSecrets, secret)
			}
		}
		secrets = filteredSecrets
	}
	if filter.SecretTypes != nil {
		filteredSecrets = []*vmsecretsmodels.ArkSIAVMSecret{}
		for _, secret := range secrets {
			if slices.Contains(filter.SecretTypes, secret.SecretType) {
				filteredSecrets = append(filteredSecrets, secret)
			}
		}
		secrets = filteredSecrets
	}
	return secrets, nil
}

// Secret retrieves a specific secret from the SIA VM secrets service.
func (s *ArkSIASecretsVMService) Secret(getSecret *vmsecretsmodels.ArkSIAVMGetSecret) (*vmsecretsmodels.ArkSIAVMSecret, error) {
	s.Logger.Info("Getting secret [%s]", getSecret.SecretID)
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
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get secret - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	secretJSON, err := common.DeserializeJSONSnake(response.Body)
	if err != nil {
		return nil, err
	}
	var secret vmsecretsmodels.ArkSIAVMSecret
	err = mapstructure.Decode(secretJSON, &secret)
	if err != nil {
		return nil, err
	}
	return &secret, nil
}

// SecretsStats retrieves statistics about secrets in the SIA VM secrets service.
func (s *ArkSIASecretsVMService) SecretsStats() (*vmsecretsmodels.ArkSIAVMSecretsStats, error) {
	secrets, err := s.ListSecrets()
	if err != nil {
		return nil, err
	}
	var secretsStats vmsecretsmodels.ArkSIAVMSecretsStats
	secretsStats.SecretsCount = len(secrets)
	for _, secret := range secrets {
		if secret.IsActive {
			secretsStats.ActiveSecretsCount++
		} else {
			secretsStats.InactiveSecretsCount++
		}
		if _, ok := secretsStats.SecretsCountByType[secret.SecretType]; !ok {
			secretsStats.SecretsCountByType[secret.SecretType] = 0
		}
		secretsStats.SecretsCountByType[secret.SecretType]++
	}
	return &secretsStats, nil
}

// ServiceConfig returns the service configuration for the ArkSIASecretsVMService.
func (s *ArkSIASecretsVMService) ServiceConfig() services.ArkServiceConfig {
	return SIASecretsVMServiceConfig
}
