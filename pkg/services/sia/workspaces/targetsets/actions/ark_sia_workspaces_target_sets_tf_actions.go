package actions

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	targetsetsmodels "github.com/cyberark/ark-sdk-golang/pkg/services/sia/workspaces/targetsets/models"
)

// TerraformActionWorkspacesTargetSetsResource is a struct that defines the SIA workspaces target sets resource action for the Ark service for Terraform.
var TerraformActionWorkspacesTargetSetsResource = &actions.ArkServiceTerraformResourceActionDefinition{
	ArkServiceBaseTerraformActionDefinition: actions.ArkServiceBaseTerraformActionDefinition{
		ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
			ActionName:        "sia-workspaces-target-set",
			ActionDescription: "SIA Workspaces Target Set resource, manages target set information about one or more targets and how they are represented, along with association to relevant secret.",
			ActionVersion:     1,
			Schemas:           ActionToSchemaMap,
		},
		StateSchema: &targetsetsmodels.ArkSIATargetSet{},
	},
	SupportedOperations: []actions.ArkServiceActionOperation{
		actions.CreateOperation,
		actions.ReadOperation,
		actions.UpdateOperation,
		actions.DeleteOperation,
		actions.StateOperation,
	},
	ActionsMappings: map[actions.ArkServiceActionOperation]string{
		actions.CreateOperation: "add-target-set",
		actions.ReadOperation:   "target-set",
		actions.UpdateOperation: "update-target-set",
		actions.DeleteOperation: "delete-target-set",
	},
}

// TerraformActionWorkspacesTargetSetsDataSource is a struct that defines the sia workspaces target sets data source action for the Ark service for Terraform.
var TerraformActionWorkspacesTargetSetsDataSource = &actions.ArkServiceTerraformDataSourceActionDefinition{
	ArkServiceBaseTerraformActionDefinition: actions.ArkServiceBaseTerraformActionDefinition{
		ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
			ActionName:        "sia-workspaces-target-set",
			ActionDescription: "SIA Workspaces Target Set data source, reads target set information and metadata, based on the id of the target set.",
			ActionVersion:     1,
			Schemas:           ActionToSchemaMap,
		},
		StateSchema: &targetsetsmodels.ArkSIATargetSet{},
		ExtraRequiredAttributes: []string{
			"id",
		},
	},
	DataSourceAction: "target-set",
}
