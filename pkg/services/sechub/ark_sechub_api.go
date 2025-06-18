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
	serviceInfoService   *serviceinfo.ArkSecHubServiceInfoService
	secretsService       *secrets.ArkSecHubSecretsService
	secretStoresService  *secretstores.ArkSecHubSecretStoresService
	syncPoliciesService  *syncpolicies.ArkSecHubSyncPoliciesService
}

// NewArkSecHubAPI creates a new instance of ArkSecHubAPI with the provided ArkISPAuth.
func NewArkSecHubAPI(ispAuth *auth.ArkISPAuth) (*ArkSecHubAPI, error) {
	var baseIspAuth auth.ArkAuth = ispAuth
	configurationService, err := configuration.NewArkSecHubConfigurationService(baseIspAuth)
	if err != nil {
		return nil, err
	}
	filtersService, err := filters.NewArkSecHubFiltersService(baseIspAuth)
	if err != nil {
		return nil, err
	}
	scansService, err := scans.NewArkSecHubScansService(baseIspAuth)
	if err != nil {
		return nil, err
	}
	secretsService, err := secrets.NewArkSecHubSecretsService(baseIspAuth)
	if err != nil {
		return nil, err
	}
	secretStoresService, err := secretstores.NewArkSecHubSecretStoresService(baseIspAuth)
	if err != nil {
		return nil, err
	}
	serviceInfoService, err := serviceinfo.NewArkSecHubServiceInfoService(baseIspAuth)
	if err != nil {
		return nil, err
	}
	syncPoliciesService, err := syncpolicies.NewArkSecHubSyncPoliciesService(baseIspAuth)
	if err != nil {
		return nil, err
	}
	return &ArkSecHubAPI{
		serviceInfoService:   serviceInfoService,
		configurationService: configurationService,
		filtersService:       filtersService,
		scansService:         scansService,
		secretStoresService:  secretStoresService,
		secretsService:       secretsService,
		syncPoliciesService:  syncPoliciesService,
	}, nil
}

// Configuration returns the configuration service of the ArkSecHubAPI instance.
func (api *ArkSecHubAPI) Configuration() *configuration.ArkSecHubConfigurationService {
	return api.configurationService
}

// Filters returns the filters service of the ArkSecHubAPI instance.
func (api *ArkSecHubAPI) Filters() *filters.ArkSecHubFiltersService {
	return api.filtersService
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
	return api.serviceInfoService
}

// SyncPolicies returns the sync policies service of the ArkSecHubAPI instance.
func (api *ArkSecHubAPI) SyncPolicies() *syncpolicies.ArkSecHubSyncPoliciesService {
	return api.syncPoliciesService
}
