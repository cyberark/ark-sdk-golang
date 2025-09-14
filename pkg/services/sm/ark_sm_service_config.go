package sm

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	smactions "github.com/cyberark/ark-sdk-golang/pkg/services/sm/actions"
)

// ServiceConfig is the configuration for the Session Monitoring service.
var ServiceConfig = services.ArkServiceConfig{
	ServiceName:                "sm",
	RequiredAuthenticatorNames: []string{"isp"},
	OptionalAuthenticatorNames: []string{},
	ActionsConfigurations: map[actions.ArkServiceActionType][]actions.ArkServiceActionDefinition{
		actions.ArkServiceActionTypeCLI: {
			smactions.CLIAction,
		},
	},
}

// ServiceGenerator is the function that creates a new instance of the Session Monitoring service.
var ServiceGenerator = NewArkSMService

// Module init, registers the service configuration.
func init() {
	err := services.Register(ServiceConfig, true)
	if err != nil {
		panic(err)
	}
}
