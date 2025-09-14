package db

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	uapsiadbactions "github.com/cyberark/ark-sdk-golang/pkg/services/uap/sia/db/actions"
)

// ServiceConfig defines the service configuration for ArkUAPSIADBService.
var ServiceConfig = services.ArkServiceConfig{
	ServiceName:                "uap-db",
	RequiredAuthenticatorNames: []string{"isp"},
	OptionalAuthenticatorNames: []string{},
	ActionsConfigurations: map[actions.ArkServiceActionType][]actions.ArkServiceActionDefinition{
		actions.ArkServiceActionTypeCLI: {
			uapsiadbactions.CLIAction,
		},
		actions.ArkServiceActionTypeTerraformResource: {
			uapsiadbactions.TerraformActionDBResource,
		},
		actions.ArkServiceActionTypeTerraformDataSource: {
			uapsiadbactions.TerraformActionDBDataSource,
		},
	},
}

// ServiceGenerator is the function that generates a new instance of the ArkUAPSIADBService.
var ServiceGenerator = NewArkUAPSIADBService

// Module init, registers the service configuration.
func init() {
	err := services.Register(ServiceConfig, false)
	if err != nil {
		panic(err)
	}
}
