package roles

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	identityrolesactions "github.com/cyberark/ark-sdk-golang/pkg/services/identity/roles/actions"
)

// ServiceConfig is the configuration for the identity users service.
var ServiceConfig = services.ArkServiceConfig{
	ServiceName:                "identity-roles",
	RequiredAuthenticatorNames: []string{"isp"},
	OptionalAuthenticatorNames: []string{},
	ActionsConfigurations: map[actions.ArkServiceActionType][]actions.ArkServiceActionDefinition{
		actions.ArkServiceActionTypeCLI: {
			identityrolesactions.CLIAction,
		},
	},
}

// ServiceGenerator is the function that generates a new instance of the ArkIdentityRolesService.
var ServiceGenerator = NewArkIdentityRolesService

// Module init, registers the service configuration.
func init() {
	err := services.Register(ServiceConfig, false)
	if err != nil {
		panic(err)
	}
}
