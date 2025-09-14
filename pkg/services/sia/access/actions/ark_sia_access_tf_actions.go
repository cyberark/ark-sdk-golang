package actions

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	accessmodels "github.com/cyberark/ark-sdk-golang/pkg/services/sia/access/models"
)

// TerraformActionAccessConnectorResource is a struct that defines the SIA access resource action for the Ark service for Terraform.
var TerraformActionAccessConnectorResource = &actions.ArkServiceTerraformResourceActionDefinition{
	ArkServiceBaseTerraformActionDefinition: actions.ArkServiceBaseTerraformActionDefinition{
		ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
			ActionName:        "sia-access-connector",
			ActionDescription: "SIA access connector resource, manages SIA access connector installation and removal on SIA and target machines.",
			ActionVersion:     1,
			Schemas:           ActionToSchemaMap,
		},
		ExtraRequiredAttributes: []string{
			"connector_os",
			"connector_type",
			"target_machine",
			"username",
		},
		SensitiveAttributes: []string{
			"password",
			"private_key_contents",
		},
		StateSchema: &accessmodels.ArkSIAAccessConnectorID{},
	},
	RawStateInference: true,
	SupportedOperations: []actions.ArkServiceActionOperation{
		actions.CreateOperation,
		actions.DeleteOperation,
		actions.StateOperation,
	},
	ActionsMappings: map[actions.ArkServiceActionOperation]string{
		actions.CreateOperation: "install-connector",
		actions.DeleteOperation: "uninstall-connector",
	},
}
