package actions

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	cmgrmodels "github.com/cyberark/ark-sdk-golang/pkg/services/cmgr/models"
)

// TerraformActionNetworkResource is a struct that defines the CMGR action for the Ark service for Terraform.
var TerraformActionNetworkResource = &actions.ArkServiceTerraformResourceActionDefinition{
	ArkServiceBaseTerraformActionDefinition: actions.ArkServiceBaseTerraformActionDefinition{
		ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
			ActionName:        "cmgr-network",
			ActionDescription: "Connector manager network resource, manages network associated to pools.",
			ActionVersion:     1,
			Schemas:           ActionToSchemaMap,
		},
		StateSchema: &cmgrmodels.ArkCmgrNetwork{},
	},
	SupportedOperations: []actions.ArkServiceActionOperation{
		actions.CreateOperation,
		actions.ReadOperation,
		actions.UpdateOperation,
		actions.DeleteOperation,
		actions.StateOperation,
	},
	ActionsMappings: map[actions.ArkServiceActionOperation]string{
		actions.CreateOperation: "add-network",
		actions.ReadOperation:   "network",
		actions.UpdateOperation: "update-network",
		actions.DeleteOperation: "delete-network",
	},
}

// TerraformActionPoolResource is a struct that defines the CMGR pool action for the Ark service for Terraform.
var TerraformActionPoolResource = &actions.ArkServiceTerraformResourceActionDefinition{
	ArkServiceBaseTerraformActionDefinition: actions.ArkServiceBaseTerraformActionDefinition{
		ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
			ActionName:        "cmgr-pool",
			ActionDescription: "Connector manager pool resource, manages pool of SIA / system connectors.",
			ActionVersion:     1,
			Schemas:           ActionToSchemaMap,
		},
		ExtraRequiredAttributes: []string{
			"assigned_network_ids",
		},
		StateSchema: &cmgrmodels.ArkCmgrPool{},
	},
	SupportedOperations: []actions.ArkServiceActionOperation{
		actions.CreateOperation,
		actions.ReadOperation,
		actions.UpdateOperation,
		actions.DeleteOperation,
		actions.StateOperation,
	},
	ActionsMappings: map[actions.ArkServiceActionOperation]string{
		actions.CreateOperation: "add-pool",
		actions.ReadOperation:   "pool",
		actions.UpdateOperation: "update-pool",
		actions.DeleteOperation: "delete-pool",
	},
}

// TerraformActionPoolIdentifierResource is a struct that defines the CMGR pool identifier action for the Ark service for Terraform.
var TerraformActionPoolIdentifierResource = &actions.ArkServiceTerraformResourceActionDefinition{
	ArkServiceBaseTerraformActionDefinition: actions.ArkServiceBaseTerraformActionDefinition{
		ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
			ActionName:        "cmgr-pool-identifier",
			ActionDescription: "Connector manager pool identifier resource, is associated to a pool for identifying the pool in a simplified manner, and not only via the network name",
			ActionVersion:     1,
			Schemas:           ActionToSchemaMap,
		},
		ExtraRequiredAttributes: []string{
			"pool_id",
			"type",
			"value",
		},
		StateSchema: &cmgrmodels.ArkCmgrPoolIdentifier{},
	},
	SupportedOperations: []actions.ArkServiceActionOperation{
		actions.CreateOperation,
		actions.ReadOperation,
		actions.UpdateOperation,
		actions.DeleteOperation,
		actions.StateOperation,
	},
	ActionsMappings: map[actions.ArkServiceActionOperation]string{
		actions.CreateOperation: "add-pool-identifier",
		actions.ReadOperation:   "pool-identifier",
		actions.UpdateOperation: "update-pool-identifier",
		actions.DeleteOperation: "delete-pool-identifier",
	},
}

// TerraformActionNetworkDataSource is a struct that defines the CMGR network action for the Ark service for Terraform.
var TerraformActionNetworkDataSource = &actions.ArkServiceTerraformDataSourceActionDefinition{
	ArkServiceBaseTerraformActionDefinition: actions.ArkServiceBaseTerraformActionDefinition{
		ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
			ActionName:        "cmgr-network",
			ActionDescription: "Cmgr Network data source, reads network information and metadata, based on the id of the network.",
			ActionVersion:     1,
			Schemas:           ActionToSchemaMap,
		},
		ExtraRequiredAttributes: []string{
			"network_id",
		},
		StateSchema: &cmgrmodels.ArkCmgrNetwork{},
	},
	DataSourceAction: "network",
}

// TerraformActionPoolDataSource is a struct that defines the CMGR pool action for the Ark service for Terraform.
var TerraformActionPoolDataSource = &actions.ArkServiceTerraformDataSourceActionDefinition{
	ArkServiceBaseTerraformActionDefinition: actions.ArkServiceBaseTerraformActionDefinition{
		ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
			ActionName:        "cmgr-pool",
			ActionDescription: "Cmgr Pool data source, reads pool information and metadata, based on the id of the pool.",
			ActionVersion:     1,
			Schemas:           ActionToSchemaMap,
		},
		ExtraRequiredAttributes: []string{
			"pool_id",
		},
		StateSchema: &cmgrmodels.ArkCmgrPool{},
	},
	DataSourceAction: "pool",
}

// TerraformActionPoolIdentifierDataSource is a struct that defines the CMGR pool identifier action for the Ark service for Terraform.
var TerraformActionPoolIdentifierDataSource = &actions.ArkServiceTerraformDataSourceActionDefinition{
	ArkServiceBaseTerraformActionDefinition: actions.ArkServiceBaseTerraformActionDefinition{
		ArkServiceBaseActionDefinition: actions.ArkServiceBaseActionDefinition{
			ActionName:        "cmgr-pool-identifier",
			ActionDescription: "Cmgr Pool data source, reads pool information and metadata, based on the id of the pool.",
			ActionVersion:     1,
			Schemas:           ActionToSchemaMap,
		},
		ExtraRequiredAttributes: []string{
			"pool_id",
			"identifier_id",
		},
		StateSchema: &cmgrmodels.ArkCmgrPoolIdentifier{},
	},
	DataSourceAction: "pool-identifier",
}
