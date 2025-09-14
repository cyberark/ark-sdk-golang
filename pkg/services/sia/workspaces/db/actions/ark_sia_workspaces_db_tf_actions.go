package actions

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	workspacesdbmodels "github.com/cyberark/ark-sdk-golang/pkg/services/sia/workspaces/db/models"
)

// TerraformActionWorkspacesDBResource is a struct that defines the SIA workspaces db resource action for the Ark service for Terraform.
var TerraformActionWorkspacesDBResource = &actions.ArkServiceTerraformResourceActionDefinition{
	ArkServiceBaseTerraformActionDefinition: actions.ArkServiceBaseTerraformActionDefinition{
		ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
			ActionName:        "sia-workspaces-db",
			ActionDescription: "SIA Workspaces DB resource, manages DB workspaces information and metadata, along with association to relevant secret.",
			ActionVersion:     1,
			Schemas:           ActionToSchemaMap,
		},
		StateSchema: &workspacesdbmodels.ArkSIADBDatabase{},
	},
	SupportedOperations: []actions.ArkServiceActionOperation{
		actions.CreateOperation,
		actions.ReadOperation,
		actions.UpdateOperation,
		actions.DeleteOperation,
		actions.StateOperation,
	},
	ActionsMappings: map[actions.ArkServiceActionOperation]string{
		actions.CreateOperation: "add-database",
		actions.ReadOperation:   "database",
		actions.UpdateOperation: "update-database",
		actions.DeleteOperation: "delete-database",
	},
}

// TerraformActionWorkspacesDBDataSource is a struct that defines the sia workspaces db data source action for the Ark service for Terraform.
var TerraformActionWorkspacesDBDataSource = &actions.ArkServiceTerraformDataSourceActionDefinition{
	ArkServiceBaseTerraformActionDefinition: actions.ArkServiceBaseTerraformActionDefinition{
		ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
			ActionName:        "sia-workspaces-db",
			ActionDescription: "SIA Workspaces DB data source, reads DB information and metadata, based on the id of the database.",
			ActionVersion:     1,
			Schemas:           ActionToSchemaMap,
		},
		StateSchema: &workspacesdbmodels.ArkSIADBDatabase{},
	},
	DataSourceAction: "database",
}
