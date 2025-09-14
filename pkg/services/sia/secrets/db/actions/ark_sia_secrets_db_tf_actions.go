package actions

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	secretsdbmodels "github.com/cyberark/ark-sdk-golang/pkg/services/sia/secrets/db/models"
)

// TerraformActionSecretsDBResource is a struct that defines the SIA secrets db resource action for the Ark service for Terraform.
var TerraformActionSecretsDBResource = &actions.ArkServiceTerraformResourceActionDefinition{
	ArkServiceBaseTerraformActionDefinition: actions.ArkServiceBaseTerraformActionDefinition{
		ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
			ActionName:        "sia-secrets-db",
			ActionDescription: "SIA Secrets DB resource, manages DB secrets information and metadata, based on the type of secret.",
			ActionVersion:     1,
			Schemas:           ActionToSchemaMap,
		},
		SensitiveAttributes: []string{
			"password",
			"iam_secret_access_key",
			"atlas_private_key",
		},
		StateSchema: &secretsdbmodels.ArkSIADBSecretMetadata{},
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
		actions.UpdateOperation: "update-secret",
		actions.DeleteOperation: "delete-secret",
	},
}

// TerraformActionSecretsDBDataSource is a struct that defines the sia secrets db data source action for the Ark service for Terraform.
var TerraformActionSecretsDBDataSource = &actions.ArkServiceTerraformDataSourceActionDefinition{
	ArkServiceBaseTerraformActionDefinition: actions.ArkServiceBaseTerraformActionDefinition{
		ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
			ActionName:        "sia-secrets-db",
			ActionDescription: "SIA Secrets DB data source, reads DB secrets information and metadata, based on the id of the secret.",
			ActionVersion:     1,
			Schemas:           ActionToSchemaMap,
		},
		StateSchema: &secretsdbmodels.ArkSIADBSecretMetadata{},
	},
	DataSourceAction: "secret",
}
