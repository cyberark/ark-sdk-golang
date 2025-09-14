package actions

import storesmodels "github.com/cyberark/ark-sdk-golang/pkg/services/sechub/secretstores/models"

// ActionToSchemaMap is a map that defines the mapping between Sec Hub secrets stores action names and their corresponding schema types.
var ActionToSchemaMap = map[string]interface{}{
	"secret-store":             &storesmodels.ArkSecHubGetSecretStore{},
	"list-secret-stores":       nil,
	"list-secret-stores-by":    &storesmodels.ArkSecHubSecretStoresFilters{},
	"secret-store-conn-status": &storesmodels.ArkSecHubGetSecretStoreConnectionStatus{},
	"set-secret-store-state":   &storesmodels.ArkSecHubSetSecretStoreState{},
	"set-secret-stores-state":  &storesmodels.ArkSecHubSetSecretStoresState{},
	"secret-stores-stats":      nil,
	"delete-secret-store":      &storesmodels.ArkSecHubDeleteSecretStore{},
	"create-secret-store":      &storesmodels.ArkSecHubCreateSecretStore{},
	"update-secret-store":      &storesmodels.ArkSecHubUpdateSecretStore{},
}
