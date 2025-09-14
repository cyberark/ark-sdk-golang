package actions

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	uapsiadbmodels "github.com/cyberark/ark-sdk-golang/pkg/services/uap/sia/db/models"
)

// TerraformActionDBResource is a struct that defines the UAP SIA DB resource action for the Ark service for Terraform.
var TerraformActionDBResource = &actions.ArkServiceTerraformResourceActionDefinition{
	ArkServiceBaseTerraformActionDefinition: actions.ArkServiceBaseTerraformActionDefinition{
		ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
			ActionName:        "uap-db",
			ActionDescription: "UAP SIA DB Policy resource.",
			ActionVersion:     1,
			Schemas:           ActionToSchemaMap,
		},
		StateSchema: &uapsiadbmodels.ArkUAPSIADBAccessPolicy{},
		ComputedAsSetAttributes: []string{
			"days_of_the_week",
		},
	},
	ReadSchemaPath:   "metadata",
	DeleteSchemaPath: "metadata",
	SupportedOperations: []actions.ArkServiceActionOperation{
		actions.CreateOperation,
		actions.ReadOperation,
		actions.UpdateOperation,
		actions.DeleteOperation,
		actions.StateOperation,
	},
	ActionsMappings: map[actions.ArkServiceActionOperation]string{
		actions.CreateOperation: "add-policy",
		actions.ReadOperation:   "policy",
		actions.UpdateOperation: "update-policy",
		actions.DeleteOperation: "delete-policy",
	},
}

// TerraformActionDBDataSource is a struct that defines the UAP SIA DB data source action for the Ark service for Terraform.
var TerraformActionDBDataSource = &actions.ArkServiceTerraformDataSourceActionDefinition{
	ArkServiceBaseTerraformActionDefinition: actions.ArkServiceBaseTerraformActionDefinition{
		ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
			ActionName:        "uap-db",
			ActionDescription: "UAP SIA DB Policy Data Source.",
			ActionVersion:     1,
			Schemas:           ActionToSchemaMap,
		},
		StateSchema: &uapsiadbmodels.ArkUAPSIADBAccessPolicy{},
	},
	DataSourceAction: "policy",
}
