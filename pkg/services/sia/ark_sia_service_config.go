package sia

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	siaaccessactions "github.com/cyberark/ark-sdk-golang/pkg/services/sia/access/actions"
	siadbactions "github.com/cyberark/ark-sdk-golang/pkg/services/sia/db/actions"
	siak8sactions "github.com/cyberark/ark-sdk-golang/pkg/services/sia/k8s/actions"
	siasecretsactions "github.com/cyberark/ark-sdk-golang/pkg/services/sia/secrets"
	siasshcaactions "github.com/cyberark/ark-sdk-golang/pkg/services/sia/sshca/actions"
	siassoactions "github.com/cyberark/ark-sdk-golang/pkg/services/sia/sso/actions"
	siaworkspacesactions "github.com/cyberark/ark-sdk-golang/pkg/services/sia/workspaces"
)

// CLIAction is a struct that defines the SIA action for the Ark service for the CLI.
var CLIAction = &actions.ArkServiceCLIActionDefinition{
	ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
		ActionName:        "sia",
		ActionDescription: "Secure infrastructure access provides a seamless, agentless SaaS solution for session management, ideal for securing privileged access to targets spread across hybrid and cloud environments. Session management with SIA allows access with Zero Standing Privileges (ZSP) or vaulted credentials",
		ActionVersion:     1,
	},
	ActionAliases: []string{"dpa"},
	Subactions: []*actions.ArkServiceCLIActionDefinition{
		siassoactions.CLIAction,
		siak8sactions.CLIAction,
		siaworkspacesactions.CLIAction,
		siasecretsactions.CLIAction,
		siaaccessactions.CLIAction,
		siasshcaactions.CLIAction,
		siadbactions.CLIAction,
	},
}

// ServiceConfig is the configuration for the sia service.
var ServiceConfig = services.ArkServiceConfig{
	ServiceName:                "sia",
	RequiredAuthenticatorNames: []string{"isp"},
	OptionalAuthenticatorNames: []string{},
	ActionsConfigurations: map[actions.ArkServiceActionType][]actions.ArkServiceActionDefinition{
		actions.ArkServiceActionTypeCLI: {
			CLIAction,
		},
	},
}

// Module init, registers the service configuration.
func init() {
	err := services.Register(ServiceConfig, true)
	if err != nil {
		panic(err)
	}
}
