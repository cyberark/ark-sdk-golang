package sca

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	uapscaactions "github.com/cyberark/ark-sdk-golang/pkg/services/uap/sca/actions"
)

// ServiceConfig defines the service configuration for ArkUAPSCAService.
var ServiceConfig = services.ArkServiceConfig{
	ServiceName:                "uap-sca",
	RequiredAuthenticatorNames: []string{"isp"},
	OptionalAuthenticatorNames: []string{},
	ActionsConfigurations: map[actions.ArkServiceActionType][]actions.ArkServiceActionDefinition{
		actions.ArkServiceActionTypeCLI: {
			uapscaactions.CLIAction,
		},
		actions.ArkServiceActionTypeTerraformResource: {
			uapscaactions.TerraformActionSCAResource,
		},
		actions.ArkServiceActionTypeTerraformDataSource: {
			uapscaactions.TerraformActionSCADataSource,
		},
	},
}

// ServiceGenerator is the function that generates a new instance of the ArkUAPSCAService.
var ServiceGenerator = NewArkUAPSCAService

// Module init, registers the service configuration.
func init() {
	err := services.Register(ServiceConfig, false)
	if err != nil {
		panic(err)
	}
}
