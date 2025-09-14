package db

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	siadbactions "github.com/cyberark/ark-sdk-golang/pkg/services/sia/db/actions"
)

// ServiceConfig is the configuration for the ArkSIADBService.
var ServiceConfig = services.ArkServiceConfig{
	ServiceName:                "sia-db",
	RequiredAuthenticatorNames: []string{"isp"},
	OptionalAuthenticatorNames: []string{},
	ActionsConfigurations: map[actions.ArkServiceActionType][]actions.ArkServiceActionDefinition{
		actions.ArkServiceActionTypeCLI: {
			siadbactions.CLIAction,
		},
	},
}

// ServiceGenerator is the function that creates a new instance of the SIA DB service.
var ServiceGenerator = NewArkSIADBService

// Module init, registers the service configuration.
func init() {
	err := services.Register(ServiceConfig, false)
	if err != nil {
		panic(err)
	}
}
