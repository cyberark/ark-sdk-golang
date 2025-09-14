package actions

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	uapscamodels "github.com/cyberark/ark-sdk-golang/pkg/services/uap/sca/models"
)

// TerraformActionSCAResource is a struct that defines the UAP SCA resource action for the Ark service for Terraform.
var TerraformActionSCAResource = &actions.ArkServiceTerraformResourceActionDefinition{
	ArkServiceBaseTerraformActionDefinition: actions.ArkServiceBaseTerraformActionDefinition{
		ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
			ActionName:        "uap-sca",
			ActionDescription: "UAP SCA Policy resource.",
			ActionVersion:     1,
			Schemas:           ActionToSchemaMap,
		},
		StateSchema: &uapscamodels.ArkUAPSCACloudConsoleAccessPolicy{},
		ComputedAsSetAttributes: []string{
			"days_of_the_week",
		},
	},
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

// TerraformActionSCADataSource is a struct that defines the UAP SCA data source action for the Ark service for Terraform.
var TerraformActionSCADataSource = &actions.ArkServiceTerraformDataSourceActionDefinition{
	ArkServiceBaseTerraformActionDefinition: actions.ArkServiceBaseTerraformActionDefinition{
		ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
			ActionName:        "uap-sca",
			ActionDescription: "UAP SCA Policy Data Source.",
			ActionVersion:     1,
			Schemas:           ActionToSchemaMap,
		},
		StateSchema: &uapscamodels.ArkUAPSCACloudConsoleAccessPolicy{},
	},
	DataSourceAction: "policy",
}
