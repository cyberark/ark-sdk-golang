package sechub

import (
	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/services/sechub/configuration"
	"github.com/cyberark/ark-sdk-golang/pkg/services/sechub/filters"
	"github.com/cyberark/ark-sdk-golang/pkg/services/sechub/scans"
	"github.com/cyberark/ark-sdk-golang/pkg/services/sechub/secrets"
	"github.com/cyberark/ark-sdk-golang/pkg/services/sechub/secretstores"
	"github.com/cyberark/ark-sdk-golang/pkg/services/sechub/serviceinfo"
	"github.com/cyberark/ark-sdk-golang/pkg/services/sechub/syncpolicies"
)

// ArkSecHubAPI is a struct that provides access to the Ark SecHub API as a wrapped set of services.
type ArkSecHubAPI struct {
	configurationService *configuration.ArkSecHubConfigurationService
	filtersService       *filters.ArkSecHubFiltersService
	scansService         *scans.ArkSecHubScansService
	serviceinfoService   *serviceinfo.ArkSecHubServiceInfoService
	secretStoresService  *secretstores.ArkSecHubSecretStoresService
	secretsService       *secrets.ArkSecHubSecretsService
	syncPoliciesService  *syncpolicies.ArkSecHubSyncPoliciesService
}

// NewArkSecHubAPI creates a new instance of ArkSecHubAPI with the provided ArkISPAuth.
func NewArkSecHubAPI(ispAuth *auth.ArkISPAuth) (*ArkSecHubAPI, error) {
	var baseIspAuth auth.ArkAuth = ispAuth
	serviceinfoService, err := serviceinfo.NewArkSecHubServiceInfoService(baseIspAuth)
	if err != nil {
		return nil, err
	}
	return &ArkSecHubAPI{
		serviceinfoService: serviceinfoService,
	}, nil
}

// Configuration returns the configuration service of the ArkSecHubAPI instance.
func (api *ArkSecHubAPI) Configuration() *configuration.ArkSecHubConfigurationService {
	return api.configurationService
}

// Scans returns the scans service of the ArkSecHubAPI instance.
func (api *ArkSecHubAPI) Scans() *scans.ArkSecHubScansService {
	return api.scansService
}

// Secrets returns the Secrets service of the ArkSecHubAPI instance.
func (api *ArkSecHubAPI) Secrets() *secrets.ArkSecHubSecretsService {
	return api.secretsService
}

// SecretStores returns the secret stores service of the ArkSecHubAPI instance.
func (api *ArkSecHubAPI) SecretStores() *secretstores.ArkSecHubSecretStoresService {
	return api.secretStoresService
}

// ServiceInfo returns the service info service of the ArkSecHubAPI instance.
func (api *ArkSecHubAPI) ServiceInfo() *serviceinfo.ArkSecHubServiceInfoService {
	return api.serviceinfoService
}

// SyncPolicies returns the sync policies service of the ArkSecHubAPI instance.
func (api *ArkSecHubAPI) SyncPolicies() *syncpolicies.ArkSecHubSyncPoliciesService {
	return api.syncPoliciesService
}

// Filters returns the filters service of the ArkSecHubAPI instance.
func (api *ArkSecHubAPI) Filters() *filters.ArkSecHubFiltersService {
	return api.filtersService
}
