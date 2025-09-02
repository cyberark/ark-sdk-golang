package models

// ArkCmgrAddPoolIdentifier is a struct representing the filter for adding identifiers to a pool in the Ark CMGR service.
type ArkCmgrAddPoolIdentifier struct {
	Type  string `json:"type" mapstructure:"type" flag:"type" desc:"Type of the identifier to add (GENERAL_FQDN,GENERAL_HOSTNAME,AWS_ACCOUNT_ID,AWS_VPC,AWS_SUBNET,AZURE_SUBSCRIPTION,AZURE_VNET,AZURE_SUBNET,GCP_PROJECT,GCP_NETWORK,GCP_SUBNET)" choices:"GENERAL_FQDN,GENERAL_HOSTNAME,AWS_ACCOUNT_ID,AWS_VPC,AWS_SUBNET,AZURE_SUBSCRIPTION,AZURE_VNET,AZURE_SUBNET,GCP_PROJECT,GCP_NETWORK,GCP_SUBNET"`
	Value string `json:"value" mapstructure:"value" flag:"value" desc:"Value of the identifier"`
}

// ArkCmgrAddPoolSingleIdentifier is a struct representing the filter for adding a single identifier to a pool in the Ark CMGR service.
type ArkCmgrAddPoolSingleIdentifier struct {
	Type   string `json:"type" mapstructure:"type" flag:"type" desc:"Type of the identifier to add (GENERAL_FQDN,GENERAL_HOSTNAME,AWS_ACCOUNT_ID,AWS_VPC,AWS_SUBNET,AZURE_SUBSCRIPTION,AZURE_VNET,AZURE_SUBNET,GCP_PROJECT,GCP_NETWORK,GCP_SUBNET)" choices:"GENERAL_FQDN,GENERAL_HOSTNAME,AWS_ACCOUNT_ID,AWS_VPC,AWS_SUBNET,AZURE_SUBSCRIPTION,AZURE_VNET,AZURE_SUBNET,GCP_PROJECT,GCP_NETWORK,GCP_SUBNET"`
	Value  string `json:"value" mapstructure:"value" flag:"value" desc:"Value of the identifier"`
	PoolID string `json:"pool_id" mapstructure:"pool_id" flag:"pool-id" desc:"ID of the pool to add the identifier to"`
}

// ArkCmgrAddPoolBulkIdentifier is a struct representing the filter for adding multiple identifiers to a pool in the Ark CMGR service.
type ArkCmgrAddPoolBulkIdentifier struct {
	PoolID      string                     `json:"pool_id" mapstructure:"pool_id" flag:"pool-id" desc:"ID of the pool to add the identifiers to"`
	Identifiers []ArkCmgrAddPoolIdentifier `json:"identifiers" mapstructure:"identifiers" flag:"identifiers" desc:"Identifiers to add"`
}
