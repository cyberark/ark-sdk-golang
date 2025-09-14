package sso

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	siassoactions "github.com/cyberark/ark-sdk-golang/pkg/services/sia/sso/actions"
)

// ServiceConfig is the configuration for the SSO service.
var ServiceConfig = services.ArkServiceConfig{
	ServiceName:                "sia-sso",
	RequiredAuthenticatorNames: []string{"isp"},
	OptionalAuthenticatorNames: []string{},
	ActionsConfigurations: map[actions.ArkServiceActionType][]actions.ArkServiceActionDefinition{
		actions.ArkServiceActionTypeCLI: {
			siassoactions.CLIAction,
		},
	},
}

// ServiceGenerator is the function that creates a new instance of the SIA SSO service.
var ServiceGenerator = NewArkSIASSOService

// Module init, registers the service configuration.
func init() {
	err := services.Register(ServiceConfig, false)
	if err != nil {
		panic(err)
	}
}
