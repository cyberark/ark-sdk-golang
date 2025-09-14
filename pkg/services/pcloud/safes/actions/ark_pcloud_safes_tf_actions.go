package actions

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	safesmodels "github.com/cyberark/ark-sdk-golang/pkg/services/pcloud/safes/models"
)

// TerraformActionSafeResource is a struct that defines the pCloud safe resource action for the Ark service for Terraform.
var TerraformActionSafeResource = &actions.ArkServiceTerraformResourceActionDefinition{
	ArkServiceBaseTerraformActionDefinition: actions.ArkServiceBaseTerraformActionDefinition{
		ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
			ActionName:        "pcloud-safe",
			ActionDescription: "pCloud safe resource, manages pCloud safes information and metadata.",
			ActionVersion:     1,
			Schemas:           ActionToSchemaMap,
		},
		StateSchema: &safesmodels.ArkPCloudSafe{},
	},
	SupportedOperations: []actions.ArkServiceActionOperation{
		actions.CreateOperation,
		actions.ReadOperation,
		actions.UpdateOperation,
		actions.DeleteOperation,
		actions.StateOperation,
	},
	ActionsMappings: map[actions.ArkServiceActionOperation]string{
		actions.CreateOperation: "add-safe",
		actions.ReadOperation:   "safe",
		actions.UpdateOperation: "update-safe",
		actions.DeleteOperation: "delete-safe",
	},
}

// TerraformActionSafeMemberResource is a struct that defines the pCloud safe member resource action for the Ark service for Terraform.
var TerraformActionSafeMemberResource = &actions.ArkServiceTerraformResourceActionDefinition{
	ArkServiceBaseTerraformActionDefinition: actions.ArkServiceBaseTerraformActionDefinition{
		ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
			ActionName:        "pcloud-safe-member",
			ActionDescription: "pCloud safe member resource, manages pCloud safe members and their relevant permissions.",
			ActionVersion:     1,
			Schemas:           ActionToSchemaMap,
		},
		StateSchema: &safesmodels.ArkPCloudSafeMember{},
	},
	SupportedOperations: []actions.ArkServiceActionOperation{
		actions.CreateOperation,
		actions.ReadOperation,
		actions.UpdateOperation,
		actions.DeleteOperation,
		actions.StateOperation,
	},
	ActionsMappings: map[actions.ArkServiceActionOperation]string{
		actions.CreateOperation: "add-safe-member",
		actions.ReadOperation:   "safe-member",
		actions.UpdateOperation: "update-safe-member",
		actions.DeleteOperation: "delete-safe-member",
	},
}

// TerraformActionSafeDataSource is a struct that defines the pCloud safe data source action for the Ark service for Terraform.
var TerraformActionSafeDataSource = &actions.ArkServiceTerraformDataSourceActionDefinition{
	ArkServiceBaseTerraformActionDefinition: actions.ArkServiceBaseTerraformActionDefinition{
		ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
			ActionName:        "pcloud-safe",
			ActionDescription: "PCloud Safe data source, reads safe information and metadata, based on the id of the safe.",
			ActionVersion:     1,
			Schemas:           ActionToSchemaMap,
		},
		ExtraRequiredAttributes: []string{
			"safe_id",
		},
		StateSchema: &safesmodels.ArkPCloudSafe{},
	},
	DataSourceAction: "safe",
}

// TerraformActionSafeMemberDataSource is a struct that defines the pCloud safe member data source action for the Ark service for Terraform.
var TerraformActionSafeMemberDataSource = &actions.ArkServiceTerraformDataSourceActionDefinition{
	ArkServiceBaseTerraformActionDefinition: actions.ArkServiceBaseTerraformActionDefinition{
		ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
			ActionName:        "pcloud-safe-member",
			ActionDescription: "PCloud Safe Member data source, reads safe member information and metadata, based on the id of the safe and the member name.",
			ActionVersion:     1,
			Schemas:           ActionToSchemaMap,
		},
		ExtraRequiredAttributes: []string{
			"safe_id",
			"member_name",
		},
		StateSchema: &safesmodels.ArkPCloudSafeMember{},
	},
	DataSourceAction: "safe-member",
}
