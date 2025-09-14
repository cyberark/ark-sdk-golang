package sshca

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	siasshcaactions "github.com/cyberark/ark-sdk-golang/pkg/services/sia/sshca/actions"
)

// ServiceConfig is the configuration for the ArkSIASSHCAService.
var ServiceConfig = services.ArkServiceConfig{
	ServiceName:                "sia-ssh-ca",
	RequiredAuthenticatorNames: []string{"isp"},
	OptionalAuthenticatorNames: []string{},
	ActionsConfigurations: map[actions.ArkServiceActionType][]actions.ArkServiceActionDefinition{
		actions.ArkServiceActionTypeCLI: {
			siasshcaactions.CLIAction,
		},
	},
}

// ServiceGenerator is the function that creates a new instance of the SIA SSH CA service.
var ServiceGenerator = NewArkSIASSHCAService

// Module init, registers the service configuration.
func init() {
	err := services.Register(ServiceConfig, false)
	if err != nil {
		panic(err)
	}
}
