package syncpolicies

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	sechubsyncpoliciesactions "github.com/cyberark/ark-sdk-golang/pkg/services/sechub/syncpolicies/actions"
)

// ServiceConfig is the configuration for the Secrets Hub Sync Policies service.
var ServiceConfig = services.ArkServiceConfig{
	ServiceName:                "sechub-syncpolicies",
	RequiredAuthenticatorNames: []string{"isp"},
	OptionalAuthenticatorNames: []string{},
	ActionsConfigurations: map[actions.ArkServiceActionType][]actions.ArkServiceActionDefinition{
		actions.ArkServiceActionTypeCLI: {
			sechubsyncpoliciesactions.CLIAction,
		},
	},
}

// ServiceGenerator is the function that creates a new instance of the SecHub Sync Policies service.
var ServiceGenerator = NewArkSecHubSyncPoliciesService

// Module init, registers the service configuration.
func init() {
	err := services.Register(ServiceConfig, false)
	if err != nil {
		panic(err)
	}
}
