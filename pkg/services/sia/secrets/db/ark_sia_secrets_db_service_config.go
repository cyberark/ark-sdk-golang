package dbsecrets

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	siasecretsdbactions "github.com/cyberark/ark-sdk-golang/pkg/services/sia/secrets/db/actions"
)

// ServiceConfig is the configuration for the SIA DB secrets service.
var ServiceConfig = services.ArkServiceConfig{
	ServiceName:                "sia-secrets-db",
	RequiredAuthenticatorNames: []string{"isp"},
	OptionalAuthenticatorNames: []string{},
	ActionsConfigurations: map[actions.ArkServiceActionType][]actions.ArkServiceActionDefinition{
		actions.ArkServiceActionTypeCLI: {
			siasecretsdbactions.CLIAction,
		},
		actions.ArkServiceActionTypeTerraformResource: {
			siasecretsdbactions.TerraformActionSecretsDBResource,
		},
		actions.ArkServiceActionTypeTerraformDataSource: {
			siasecretsdbactions.TerraformActionSecretsDBDataSource,
		},
	},
}

// ServiceGenerator is the function that creates a new instance of the SIA DB secrets service.
var ServiceGenerator = NewArkSIASecretsDBService

// Module init, registers the service configuration.
func init() {
	err := services.Register(ServiceConfig, false)
	if err != nil {
		panic(err)
	}
}
