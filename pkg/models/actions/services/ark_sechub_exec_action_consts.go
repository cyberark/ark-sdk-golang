package services

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	sechubconfiguration "github.com/cyberark/ark-sdk-golang/pkg/services/sechub/configuration/models"
	filtersmodels "github.com/cyberark/ark-sdk-golang/pkg/services/sechub/filters/models"
	sechubscans "github.com/cyberark/ark-sdk-golang/pkg/services/sechub/scans/models"
	sechubsecrets "github.com/cyberark/ark-sdk-golang/pkg/services/sechub/secrets/models"
	storesmodels "github.com/cyberark/ark-sdk-golang/pkg/services/sechub/secretstores/models"
	policiesmodels "github.com/cyberark/ark-sdk-golang/pkg/services/sechub/syncpolicies/models"
)

// ServiceInfoActionToSchemaMap is a map that defines the mapping between Sec Hub Service Info action names and their corresponding schema types.
var ServiceInfoActionToSchemaMap = map[string]interface{}{
	"service-info": nil,
}

// ServiceInfoAction is a struct that defines the Secrets Hub Service Info action for the Ark service.
var ServiceInfoAction = actions.ArkServiceActionDefinition{
	ActionName: "service-info",
	Schemas:    ServiceInfoActionToSchemaMap,
}

// ScansActionToSchemaMap is a map that defines the mapping between Sec Hub scans action names and their corresponding schema types.
var ScansActionToSchemaMap = map[string]interface{}{
	"scans":        nil,
	"scans-stats":  nil,
	"trigger-scan": &sechubscans.ArkSecHubTriggerScans{},
}

// ScansAction is a struct that defines the scans action for the
var ScansAction = actions.ArkServiceActionDefinition{
	ActionName: "scans",
	Schemas:    ScansActionToSchemaMap,
}

// ConfigurationActionToSchemaMap is a map that defines the mapping between Sec Hub configuration action names and their corresponding schema types.
var ConfigurationActionToSchemaMap = map[string]interface{}{
	"configuration":     nil,
	"set-configuration": &sechubconfiguration.ArkSecHubSetConfiguration{},
}

// ConfigurationAction is a struct that defines the configuration action for the Ark service.
var ConfigurationAction = actions.ArkServiceActionDefinition{
	ActionName: "configuration",
	Schemas:    ConfigurationActionToSchemaMap,
}

// FiltersActionToSchemaMap is a map that defines the mapping between Sec Hub filters action names and their corresponding schema types.
var FiltersActionToSchemaMap = map[string]interface{}{
	"filter":        &filtersmodels.ArkSecHubGetFilter{},
	"list-filters":  &filtersmodels.ArkSecHubGetFilters{},
	"add-filter":    &filtersmodels.ArkSecHubAddFilter{},
	"delete-filter": &filtersmodels.ArkSecHubDeleteFilter{},
}

// FiltersAction is a struct that defines the filters action for the Ark service.
var FiltersAction = actions.ArkServiceActionDefinition{
	ActionName: "filters",
	Schemas:    FiltersActionToSchemaMap,
}

// SecretsStoresActionToSchemaMap is a map that defines the mapping between Sec Hub secrets stores action names and their corresponding schema types.
var SecretsStoresActionToSchemaMap = map[string]interface{}{
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

// SecretsStoresAction is a struct that defines the secrets stores action for the Ark service.
var SecretsStoresAction = actions.ArkServiceActionDefinition{
	ActionName: "secret-stores",
	Schemas:    SecretsStoresActionToSchemaMap,
}

// SecretsSHActionToSchemaMap is a map that defines the mapping between Sec Hub secrets action names and their corresponding schema types.
var SecretsSHActionToSchemaMap = map[string]interface{}{
	"secrets":         nil,
	"list-secrets-by": &sechubsecrets.ArkSecHubSecretsFilter{},
	"secrets-stats":   nil,
}

// SecretsSHAction is a struct that defines the secrets action for the Ark service.
var SecretsSHAction = actions.ArkServiceActionDefinition{
	ActionName: "secrets",
	Schemas:    SecretsSHActionToSchemaMap,
}

// SyncPoliciesActionToSchemaMap is a map that defines the mapping between Sec Hub sync policies action names and their corresponding schema types.
var SyncPoliciesActionToSchemaMap = map[string]interface{}{
	"create-sync-policy":    &policiesmodels.ArkSechubCreateSyncPolicy{},
	"delete-sync-policy":    &policiesmodels.ArkSecHubDeleteSyncPolicy{},
	"sync-policy":           &policiesmodels.ArkSecHubGetSyncPolicy{},
	"list-sync-policies":    &policiesmodels.ArkSecHubGetSyncPolicies{},
	"list-sync-policies-by": &policiesmodels.ArkSecHubSyncPoliciesFilters{},
	"set-sync-policy-state": &policiesmodels.ArkSecHubSetSyncPolicyState{},
	"sync-policies-stats":   nil,
}

// SyncPoliciesAction is a struct that defines the sync policies action for the Ark service.
var SyncPoliciesAction = actions.ArkServiceActionDefinition{
	ActionName: "sync-policies",
	Schemas:    SyncPoliciesActionToSchemaMap,
}

// SecHubActions is a struct that defines the SecHub action for the Ark service.
var SecHubActions = &actions.ArkServiceActionDefinition{
	ActionName:        "sechub",
	ActionAliases:     []string{"secretshub", "sh"},
	ActionDescription: "Secrets Hub is a CyberArk SaaS solution that simplifies managing and consuming secrets in the Cloud Service Providers’ native secret managers.",
	Subactions: []*actions.ArkServiceActionDefinition{
		&ServiceInfoAction,
		&ConfigurationAction,
		&ScansAction,
		&FiltersAction,
		&SecretsStoresAction,
		&SecretsSHAction,
		&SyncPoliciesAction,
	},
}
