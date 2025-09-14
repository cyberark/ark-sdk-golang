package db

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	siaworkspacesdbactions "github.com/cyberark/ark-sdk-golang/pkg/services/sia/workspaces/db/actions"
)

// ServiceConfig is the configuration for the SIA db workspace service.
var ServiceConfig = services.ArkServiceConfig{
	ServiceName:                "sia-workspaces-db",
	RequiredAuthenticatorNames: []string{"isp"},
	OptionalAuthenticatorNames: []string{},
	ActionsConfigurations: map[actions.ArkServiceActionType][]actions.ArkServiceActionDefinition{
		actions.ArkServiceActionTypeCLI: {
			siaworkspacesdbactions.CLIAction,
		},
		actions.ArkServiceActionTypeTerraformResource: {
			siaworkspacesdbactions.TerraformActionWorkspacesDBResource,
		},
		actions.ArkServiceActionTypeTerraformDataSource: {
			siaworkspacesdbactions.TerraformActionWorkspacesDBDataSource,
		},
	},
}

// ServiceGenerator is the function that creates a new instance of the SIA Workspaces DB service.
var ServiceGenerator = NewArkSIAWorkspacesDBService

// Module init, registers the service configuration.
func init() {
	err := services.Register(ServiceConfig, false)
	if err != nil {
		panic(err)
	}
}
