package models

// ArkSecHubUpdateSecretStore defines the structure for updating a secret store in the Ark Secrets Hub.
type ArkSecHubUpdateSecretStore struct {
	SecretStoreID string `json:"secret_store_id" mapstructure:"secret_store_id" flag:"secret-store-id" validate:"required" desc:"The unique identifier of the secret store to update"`
	Description   string `json:"description,omitempty" mapstructure:"description,omitempty" flag:"description" desc:"A description of the secret store."`
	Name          string `json:"name,omitempty" mapstructure:"name,omitempty" flag:"name" desc:"The name of the secret store. It should be unique per tenant."`
	// Data contains the specific data for the secret store type.
	Data *ArkSecHubUpdateSecretStoreData `json:"data,omitempty" mapstructure:"data,omitempty" desc:"The data related to the secret store as defined in the cloud platform."`
}

// ArkSecHubUpdateSecretStoreData defines the data structure for updating a secret store in the Ark Secrets Hub.
type ArkSecHubUpdateSecretStoreData struct {
	// AWS ASM Specific Fields
	AccountAlias string `json:"account_alias,omitempty" mapstructure:"account_alias,omitempty" flag:"aws-account-alias" desc:"AWS: The alias of your AWS account"`
	RoleName     string `json:"role_name,omitempty" mapstructure:"role_name,omitempty" flag:"aws-rolename" desc:"AWS: The role used to allow Secrets Hub to manage secrets in your AWS Secrets Manager"`
	// Azure AKV Specific Fields
	AppClientDirectoryID string `json:"app_client_directory_id,omitempty" mapstructure:"app_client_directory_id,omitempty" flag:"azure-app-client-directory-id" desc:"AZURE: The Azure Active Directory ID of the application that has access to the Azure Key Vault"`
	AzureVaultURL        string `json:"azure_vault_url,omitempty" mapstructure:"azure_vault_url,omitempty" flag:"azure-vault-url" desc:"AZURE: The URL of the Azure Key Vault where you store secrets. Example: https://myvault.vault.azure.net/"`
	AppClientID          string `json:"app_client_id,omitempty" mapstructure:"app_client_id,omitempty" flag:"azure-app-client-id" desc:"AZURE: The Azure Active Directory application ID of the application that has access to the Azure Key Vault"`
	SubscriptionID       string `json:"subscription_id,omitempty" mapstructure:"subscription_id,omitempty" flag:"azure-subscription-id" desc:"AZURE: The Azure subscription ID where the Azure Key Vault is stored"`
	SubscriptionName     string `json:"subscription_name,omitempty" mapstructure:"subscription_name,omitempty" flag:"azure-subscription-name" desc:"AZURE: The name of the Azure subscription where the Azure Key Vault is stored"`
	ResourceGroupName    string `json:"resource_group_name,omitempty" mapstructure:"resource_group_name,omitempty" flag:"azure-resource-group-name" desc:"AZURE: The name of the Azure resource group where the Azure Key Vault is stored"`
	// GCP GSM Specific Fields
	GcpProjectName            string `json:"gcp_project_name,omitempty" mapstructure:"gcp_project_name,omitempty" flag:"gcp-project-name" desc:"GCP: The name of the GCP project where the GCP Secret Manager is stored"`
	GcpWorkloadIdentityPoolID string `json:"gcp_workload_identity_pool_id,omitempty" mapstructure:"gcp_workload_identity_pool_id,omitempty" flag:"gcp-workload-identity-pool-id" desc:"GCP: The GCP workload identity pool ID created for Secrets Hub to access the GCP Secret Manager"`
	GcpPoolProviderID         string `json:"gcp_pool_provider_id,omitempty" mapstructure:"gcp_pool_provider_id,omitempty" flag:"gcp-pool-provider-id" desc:"GCP: The GCP pool provider ID created for Secrets Hub to access the GCP Secret Manager"`
	ServiceAccountEmail       string `json:"service_account_email,omitempty" mapstructure:"service_account_email,omitempty" flag:"gcp-service-account-email" desc:"GCP: The service account email created for Secrets Hub to access the GCP Secret Manager"`
	// Self-Hosted Specific Fields
	Password        string `json:"password,omitempty" mapstructure:"password,omitempty" desc:"SELF HOSTED: The password of the user in PAM 'SecretsHub'" flag:"sh-password"`
	ConnectorID     string `json:"connector_id,omitempty" mapstructure:"connector_id,omitempty" desc:"SELF HOSTED: The connector unique identifier used to connect Secrets Hub and the Cloud Vendor." flag:"sh-connector-id"`
	ConnectorPoolID string `json:"connector_pool_id,omitempty" mapstructure:"connector_pool_id,omitempty" desc:"SELF HOSTED: The connector pool unique identifier used to connect PAM Self-Hosted and Secrets Hub." flag:"sh-connector-pool-id"`
	// Used by Azure, GCP
	ConnectionConfig *ArkSecHubUpdateSecretStoreConnectionConfig `json:"connection_config,omitzero" mapstructure:"connection_config,omitzero" desc:"AZURE: The network access configuration set for your target"`
}

// ArkSecHubUpdateSecretStoreConnectionConfig defines the connection configuration for updating a secret store in the Ark Secrets Hub.
type ArkSecHubUpdateSecretStoreConnectionConfig struct {
	ConnectionType string `json:"connection_type,omitempty" mapstructure:"connection_type,omitempty" flag:"connection-type" desc:"AZURE: The type of connector (CONNECTOR,PUBLIC)" default:"CONNECTOR" choices:"CONNECTOR,PUBLIC"`
	// Required if you choose 'CONNECTOR' as the connection type.
	// If you choose 'PUBLIC', these fields are not required.
	ConnectorID     string `json:"connector_id,omitempty" mapstructure:"connector_id,omitempty" flag:"connector-id" desc:"AZURE: The connector unique identifier used to connect Secrets Hub and the Cloud Vendor."`
	ConnectorPoolID string `json:"connector_pool_id,omitempty" mapstructure:"connector_pool_id,omitempty" flag:"connector-pool-id" desc:"AZURE: The connector pool unique identifier used to connect PAM Self-Hosted and Secrets Hub."`
}
