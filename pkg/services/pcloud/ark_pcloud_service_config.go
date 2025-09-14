package pcloud

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	pcloudaccountsactions "github.com/cyberark/ark-sdk-golang/pkg/services/pcloud/accounts/actions"
	pcloudsafesactions "github.com/cyberark/ark-sdk-golang/pkg/services/pcloud/safes/actions"
)

// CLIAction is the CLI action definition for the identity service.
var CLIAction = &actions.ArkServiceCLIActionDefinition{
	ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
		ActionName:        "pcloud",
		ActionDescription: "CyberArk Privilege Cloud is a SaaS solution that enables organizations to securely store, rotate and isolate credentials (for both human and non-human users), monitor sessions, and deliver scalable risk reduction to the business.",
		ActionVersion:     1,
	},
	ActionAliases: []string{"privilegecloud", "pc"},
	Subactions: []*actions.ArkServiceCLIActionDefinition{
		pcloudaccountsactions.CLIAction,
		pcloudsafesactions.CLIAction,
	},
}

// ServiceConfig is the configuration for the identity service.
var ServiceConfig = services.ArkServiceConfig{
	ServiceName:                "pcloud",
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
