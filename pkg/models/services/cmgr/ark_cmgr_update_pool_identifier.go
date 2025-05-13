package cmgr

// ArkCmgrUpdatePoolIdentifier is a struct representing the filter for updating identifiers in a pool in the Ark CMGR service.
type ArkCmgrUpdatePoolIdentifier struct {
	Type         string `json:"type" mapstructure:"type" flag:"type" desc:"Type of the identifier to update (GENERAL_FQDN,GENERAL_HOSTNAME,AWS_ACCOUNT_ID,AWS_VPC,AWS_SUBNET,AZURE_SUBSCRIPTION,AZURE_VNET,AZURE_SUBNET,GCP_PROJECT,GCP_NETWORK,GCP_SUBNET)" choices:"GENERAL_FQDN,GENERAL_HOSTNAME,AWS_ACCOUNT_ID,AWS_VPC,AWS_SUBNET,AZURE_SUBSCRIPTION,AZURE_VNET,AZURE_SUBNET,GCP_PROJECT,GCP_NETWORK,GCP_SUBNET"`
	Value        string `json:"value" mapstructure:"value" flag:"value" desc:"Value of the identifier"`
	IdentifierID string `json:"identifier_id" mapstructure:"identifier_id" flag:"identifier-id" desc:"ID of the identifier to update from the pool"`
	PoolID       string `json:"pool_id" mapstructure:"pool_id" flag:"pool-id" desc:"ID of the pool to update the identifier to"`
}
