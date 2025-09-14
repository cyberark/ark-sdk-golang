package users

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	identityusersactions "github.com/cyberark/ark-sdk-golang/pkg/services/identity/users/actions"
)

// ServiceConfig is the configuration for the identity users service.
var ServiceConfig = services.ArkServiceConfig{
	ServiceName:                "identity-users",
	RequiredAuthenticatorNames: []string{"isp"},
	OptionalAuthenticatorNames: []string{},
	ActionsConfigurations: map[actions.ArkServiceActionType][]actions.ArkServiceActionDefinition{
		actions.ArkServiceActionTypeCLI: {
			identityusersactions.CLIAction,
		},
	},
}

// ServiceGenerator is the function that generates a new instance of the ArkIdentityUsersService.
var ServiceGenerator = NewArkIdentityUsersService

// Module init, registers the service configuration.
func init() {
	err := services.Register(ServiceConfig, false)
	if err != nil {
		panic(err)
	}
}
