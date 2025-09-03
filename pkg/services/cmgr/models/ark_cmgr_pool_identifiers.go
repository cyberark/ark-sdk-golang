package models

// Possible values for the identifier type in ArkCmgrPoolIdentifier
const (
	GeneralFQDN       = "GENERAL_FQDN"
	GeneralHostname   = "GENERAL_HOSTNAME"
	AWSAccountID      = "AWS_ACCOUNT_ID"
	AWSVPC            = "AWS_VPC"
	AWSSubnet         = "AWS_SUBNET"
	AzureSubscription = "AZURE_SUBSCRIPTION"
	AzureVNet         = "AZURE_VNET"
	AzureSubnet       = "AZURE_SUBNET"
	GCPProject        = "GCP_PROJECT"
	GCPNetwork        = "GCP_NETWORK"
	GCPSubnet         = "GCP_SUBNET"
)

// ArkCmgrPoolIdentifier is a struct representing an identifier for a pool in the Ark CMGR service.
type ArkCmgrPoolIdentifier struct {
	IdentifierID string `json:"identifier_id" mapstructure:"identifier_id" flag:"identifier-id" desc:"ID of the identifier"`
	PoolID       string `json:"pool_id" mapstructure:"pool_id" flag:"pool-id" desc:"ID of the pool this identifier is associated to"`
	Type         string `json:"type" mapstructure:"type" flag:"type" desc:"Type of the identifier (GENERAL_FQDN,GENERAL_HOSTNAME,AWS_ACCOUNT_ID,AWS_VPC,AWS_SUBNET,AZURE_SUBSCRIPTION,AZURE_VNET,AZURE_SUBNET,GCP_PROJECT,GCP_NETWORK,GCP_SUBNET)" choices:"GENERAL_FQDN,GENERAL_HOSTNAME,AWS_ACCOUNT_ID,AWS_VPC,AWS_SUBNET,AZURE_SUBSCRIPTION,AZURE_VNET,AZURE_SUBNET,GCP_PROJECT,GCP_NETWORK,GCP_SUBNET"`
	Value        string `json:"value" mapstructure:"value" flag:"value" desc:"Value of the identifier"`
	CreatedAt    string `json:"created_at" mapstructure:"created_at" flag:"created-at" desc:"The creation time of the identifier"`
	UpdatedAt    string `json:"updated_at" mapstructure:"updated_at" flag:"updated-at" desc:"The last update time of the identifier"`
}

// ArkCmgrPoolIdentifiers is a struct representing a list of identifiers for pools in the Ark CMGR service.
type ArkCmgrPoolIdentifiers struct {
	Identifiers []*ArkCmgrPoolIdentifier `json:"identifiers" mapstructure:"identifiers" flag:"identifiers" desc:"Identifiers List"`
}
