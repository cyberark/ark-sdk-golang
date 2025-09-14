package actions

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	uapscaactions "github.com/cyberark/ark-sdk-golang/pkg/services/uap/sca/actions"
	uapsiadbactions "github.com/cyberark/ark-sdk-golang/pkg/services/uap/sia/db/actions"
	uapsiavmactions "github.com/cyberark/ark-sdk-golang/pkg/services/uap/sia/vm/actions"
)

// CLIAction is a struct that defines the uap action for the Ark service CLI.
var CLIAction = &actions.ArkServiceCLIActionDefinition{
	ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
		ActionName:        "uap",
		ActionDescription: "Access policies define when specified users may access particular assets and for how long. You may use access policies for cloud workspaces, virtual machines, and databases.",
		ActionVersion:     1,
		Schemas:           ActionToSchemaMap,
	},
	ActionAliases: []string{"useraccesspolicies"},
	Subactions: []*actions.ArkServiceCLIActionDefinition{
		uapscaactions.CLIAction,
		uapsiadbactions.CLIAction,
		uapsiavmactions.CLIAction,
	},
}
