package services

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	sechubconfiguration "github.com/cyberark/ark-sdk-golang/pkg/models/services/sechub/configuration"
	sechubfilters "github.com/cyberark/ark-sdk-golang/pkg/models/services/sechub/filters"
	sechubscans "github.com/cyberark/ark-sdk-golang/pkg/models/services/sechub/scans"
	sechubsecrets "github.com/cyberark/ark-sdk-golang/pkg/models/services/sechub/secrets"
	sechubsecretstores "github.com/cyberark/ark-sdk-golang/pkg/models/services/sechub/secretstores"
)

// ServiceInfoActionToSchemaMap is a map that defines the mapping between Sec Hub Service Info action names and their corresponding schema types.
var ServiceInfoActionToSchemaMap = map[string]interface{}{
	"get-service-info": nil,
}

// ServiceInfoAction is a struct that defines the Secrets Hub Service Info action for the Ark service.
var ServiceInfoAction = actions.ArkServiceActionDefinition{
	ActionName: "service-info",
	Schemas:    ServiceInfoActionToSchemaMap,
}

// ScansActionToSchemaMap is a map that defines the mapping between Sec Hub scans action names and their corresponding schema types.
var ScansActionToSchemaMap = map[string]interface{}{
	"get-scans":    nil,
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
	//"get-configuration": &sechubconfiguration.ArkSecHubGetConfiguration{},
	"get-configuration": nil,
	"set-configuration": &sechubconfiguration.ArkSecHubSetConfiguration{},
}

// ConfigurationAction is a struct that defines the configuration action for the Ark service.
var ConfigurationAction = actions.ArkServiceActionDefinition{
	ActionName: "configuration",
	Schemas:    ConfigurationActionToSchemaMap,
}

// FiltersActionToSchemaMap is a map that defines the mapping between Sec Hub filters action names and their corresponding schema types.
var FiltersActionToSchemaMap = map[string]interface{}{
	"get-filter":    &sechubfilters.ArkSecHubGetFilter{},
	"get-filters":   &sechubfilters.ArkSecHubGetFilters{},
	"add-filter":    &sechubfilters.ArkSecHubAddFilter{},
	"delete-filter": &sechubfilters.ArkSecHubDeleteFilter{},
	"filter-stats":  nil,
}

// FiltersAction is a struct that defines the filters action for the Ark service.
var FiltersAction = actions.ArkServiceActionDefinition{
	ActionName: "filters",
	Schemas:    FiltersActionToSchemaMap,
}

// SecretsStoresActionToSchemaMap is a map that defines the mapping between Sec Hub secrets stores action names and their corresponding schema types.
var SecretsStoresActionToSchemaMap = map[string]interface{}{
	"get-secret-store":             &sechubsecretstores.ArkSecHubGetSecretStore{},
	"get-secret-stores":            nil,
	"get-secret-stores-by":         &sechubsecretstores.ArkSecHubSecretStoresFilters{},
	"get-secret-store-conn-status": &sechubsecretstores.ArkSecHubGetSecretStoreConnectionStatus{},
	"set-secret-store-state":       &sechubsecretstores.ArkSecHubSetSecretStoreState{},
	"set-secret-stores-state":      &sechubsecretstores.ArkSecHubSetSecretStoresState{},
	"secret-stores-stats":          nil,
	"delete-secret-store":          &sechubsecretstores.ArkSecHubDeleteSecretStore{},
	"create-secret-store":          &sechubsecretstores.ArkSecHubCreateSecretStore{},
	"update-secret-store":          &sechubsecretstores.ArkSecHubUpdateSecretStore{},
}

// SecretsStoresAction is a struct that defines the secrets stores action for the Ark service.
var SecretsStoresAction = actions.ArkServiceActionDefinition{
	ActionName: "secret-stores",
	Schemas:    SecretsStoresActionToSchemaMap,
}

// SecretsSHActionToSchemaMap is a map that defines the mapping between Sec Hub secrets action names and their corresponding schema types.
var SecretsSHActionToSchemaMap = map[string]interface{}{
	"get-secrets":    nil,
	"get-secrets-by": &sechubsecrets.ArkSecHubSecretsFilter{},
	"secrets-stats":  nil,
}

// SecretsSHAction is a struct that defines the secrets action for the Ark service.
var SecretsSHAction = actions.ArkServiceActionDefinition{
	ActionName: "secrets",
	Schemas:    SecretsSHActionToSchemaMap,
}

// SecHubActions is a struct that defines the SecHub action for the Ark service.
var SecHubActions = &actions.ArkServiceActionDefinition{
	ActionName: "sechub",
	Subactions: []*actions.ArkServiceActionDefinition{
		&ServiceInfoAction,
		&ConfigurationAction,
		&ScansAction,
		&FiltersAction,
		&SecretsStoresAction,
		&SecretsSHAction,
	},
}
