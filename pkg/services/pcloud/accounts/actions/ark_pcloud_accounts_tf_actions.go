package actions

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	accountsmodels "github.com/cyberark/ark-sdk-golang/pkg/services/pcloud/accounts/models"
)

// TerraformActionAccountResource is a struct that defines the pCloud account resource action for the Ark service for Terraform.
var TerraformActionAccountResource = &actions.ArkServiceTerraformResourceActionDefinition{
	ArkServiceBaseTerraformActionDefinition: actions.ArkServiceBaseTerraformActionDefinition{
		ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
			ActionName:        "pcloud-account",
			ActionDescription: "pCloud account resource, manages pCloud accounts information / metadata and credentials.",
			ActionVersion:     1,
			Schemas:           ActionToSchemaMap,
		},
		ExtraRequiredAttributes: []string{
			"address",
		},
		StateSchema: &accountsmodels.ArkPCloudAccount{},
	},
	SupportedOperations: []actions.ArkServiceActionOperation{
		actions.CreateOperation,
		actions.ReadOperation,
		actions.UpdateOperation,
		actions.DeleteOperation,
		actions.StateOperation,
	},
	ActionsMappings: map[actions.ArkServiceActionOperation]string{
		actions.CreateOperation: "add-account",
		actions.ReadOperation:   "account",
		actions.UpdateOperation: "update-account",
		actions.DeleteOperation: "delete-account",
	},
}

// TerraformActionAccountDataSource is a struct that defines the pCloud account data source action for the Ark service for Terraform.
var TerraformActionAccountDataSource = &actions.ArkServiceTerraformDataSourceActionDefinition{
	ArkServiceBaseTerraformActionDefinition: actions.ArkServiceBaseTerraformActionDefinition{
		ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
			ActionName:        "pcloud-account",
			ActionDescription: "PCloud Account data source, reads account information and metadata, based on the id of the account.",
			ActionVersion:     1,
			Schemas:           ActionToSchemaMap,
		},
		ExtraRequiredAttributes: []string{
			"account_id",
		},
		StateSchema: &accountsmodels.ArkPCloudAccount{},
	},
	DataSourceAction: "account",
}
