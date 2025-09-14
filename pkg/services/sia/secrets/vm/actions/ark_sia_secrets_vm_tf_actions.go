package actions

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	secretsvmmodels "github.com/cyberark/ark-sdk-golang/pkg/services/sia/secrets/vm/models"
)

// TerraformActionSecretsVMResource is a struct that defines the SIA secrets vm resource action for the Ark service for Terraform.
var TerraformActionSecretsVMResource = &actions.ArkServiceTerraformResourceActionDefinition{
	ArkServiceBaseTerraformActionDefinition: actions.ArkServiceBaseTerraformActionDefinition{
		ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
			ActionName:        "sia-secrets-vm",
			ActionDescription: "SIA Secrets VM resource, manages VM secrets information and metadata, based on the type of secret.",
			ActionVersion:     1,
			Schemas:           ActionToSchemaMap,
		},
		SensitiveAttributes: []string{
			"provisioner_password",
			"secret_data",
		},
		ExtraRequiredAttributes: []string{
			"secret_name",
		},
		StateSchema: &secretsvmmodels.ArkSIAVMSecret{},
	},
	SupportedOperations: []actions.ArkServiceActionOperation{
		actions.CreateOperation,
		actions.ReadOperation,
		actions.UpdateOperation,
		actions.DeleteOperation,
		actions.StateOperation,
	},
	ActionsMappings: map[actions.ArkServiceActionOperation]string{
		actions.CreateOperation: "add-secret",
		actions.ReadOperation:   "secret",
		actions.UpdateOperation: "change-secret",
		actions.DeleteOperation: "delete-secret",
	},
}

// TerraformActionSecretsVMDataSource is a struct that defines the sia secrets vm data source action for the Ark service for Terraform.
var TerraformActionSecretsVMDataSource = &actions.ArkServiceTerraformDataSourceActionDefinition{
	ArkServiceBaseTerraformActionDefinition: actions.ArkServiceBaseTerraformActionDefinition{
		ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
			ActionName:        "sia-secrets-vm",
			ActionDescription: "SIA Secrets VM data source, reads VM secrets information and metadata, based on the id of the secret.",
			ActionVersion:     1,
			Schemas:           ActionToSchemaMap,
		},
		StateSchema: &secretsvmmodels.ArkSIAVMSecret{},
	},
	DataSourceAction: "secret",
}
