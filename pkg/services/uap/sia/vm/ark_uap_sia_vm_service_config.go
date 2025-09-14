package vm

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	uapsiavmactions "github.com/cyberark/ark-sdk-golang/pkg/services/uap/sia/vm/actions"
)

// ServiceConfig defines the service configuration for ArkUAPSIAVMService.
var ServiceConfig = services.ArkServiceConfig{
	ServiceName:                "uap-vm",
	RequiredAuthenticatorNames: []string{"isp"},
	OptionalAuthenticatorNames: []string{},
	ActionsConfigurations: map[actions.ArkServiceActionType][]actions.ArkServiceActionDefinition{
		actions.ArkServiceActionTypeCLI: {
			uapsiavmactions.CLIAction,
		},
		actions.ArkServiceActionTypeTerraformResource: {
			uapsiavmactions.TerraformActionVMResource,
		},
		actions.ArkServiceActionTypeTerraformDataSource: {
			uapsiavmactions.TerraformActionVMDataSource,
		},
	},
}

// ServiceGenerator is the function that generates a new instance of the UAP SIA VM service.
var ServiceGenerator = NewArkUAPSIAVMService

// Module init, registers the service configuration.
func init() {
	err := services.Register(ServiceConfig, false)
	if err != nil {
		panic(err)
	}
}
