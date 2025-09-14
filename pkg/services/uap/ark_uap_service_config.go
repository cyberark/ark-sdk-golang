package uap

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	uapactions "github.com/cyberark/ark-sdk-golang/pkg/services/uap/actions"
)

// ServiceConfig is the configuration for the uap service.
var ServiceConfig = services.ArkServiceConfig{
	ServiceName:                "uap",
	RequiredAuthenticatorNames: []string{"isp"},
	OptionalAuthenticatorNames: []string{},
	ActionsConfigurations: map[actions.ArkServiceActionType][]actions.ArkServiceActionDefinition{
		actions.ArkServiceActionTypeCLI: {
			uapactions.CLIAction,
		},
	},
}

// ServiceGenerator is the default service generator for the uap service.
var ServiceGenerator = NewArkUAPService

// Module init, registers the service configuration.
func init() {
	err := services.Register(ServiceConfig, true)
	if err != nil {
		panic(err)
	}
}
