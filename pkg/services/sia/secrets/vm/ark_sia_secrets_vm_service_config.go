package vmsecrets

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	siasecretsvmactions "github.com/cyberark/ark-sdk-golang/pkg/services/sia/secrets/vm/actions"
)

// ServiceConfig is the configuration for the SIA VM secrets service.
var ServiceConfig = services.ArkServiceConfig{
	ServiceName:                "sia-secrets-vm",
	RequiredAuthenticatorNames: []string{"isp"},
	OptionalAuthenticatorNames: []string{},
	ActionsConfigurations: map[actions.ArkServiceActionType][]actions.ArkServiceActionDefinition{
		actions.ArkServiceActionTypeCLI: {
			siasecretsvmactions.CLIAction,
		},
		actions.ArkServiceActionTypeTerraformResource: {
			siasecretsvmactions.TerraformActionSecretsVMResource,
		},
		actions.ArkServiceActionTypeTerraformDataSource: {
			siasecretsvmactions.TerraformActionSecretsVMDataSource,
		},
	},
}

// ServiceGenerator is the function that creates a new instance of the SIA VM secrets service.
var ServiceGenerator = NewArkSIASecretsVMService

// Module init, registers the service configuration.
func init() {
	err := services.Register(ServiceConfig, false)
	if err != nil {
		panic(err)
	}
}
