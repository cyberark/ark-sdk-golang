package models

// ArkSecHubCreateSecretStoreConnectionConfig defines the connection configuration for creating a secret store in the Ark Secrets Hub.
type ArkSecHubCreateSecretStoreConnectionConfig struct {
	ConnectionType string `json:"connection_type,omitempty" mapstructure:"connection_type,omitempty" flag:"connection-type" desc:"COMMON - AKV, GCP: The type of connector (CONNECTOR,PUBLIC)" default:"CONNECTOR" choices:"CONNECTOR,PUBLIC"`
	// Required if you choose 'CONNECTOR' as the connection type.
	// If you choose 'PUBLIC', these fields are not required.
	ConnectorID     string `json:"connector_id,omitempty" mapstructure:"connector_id,omitempty" flag:"connector-id" desc:"AZURE: The connector unique identifier used to connect Secrets Hub and the Cloud Vendor. Example: ManagementAgent_90c63827-7315-4284-8559-ac8d24f2666d"`
	ConnectorPoolID string `json:"connector_pool_id,omitempty" mapstructure:"connector_pool_id,omitempty" flag:"connector-pool-id" desc:"AZURE: The connector pool unique identifier used to connect PAM Self-Hosted and Secrets Hub."`
}

// ArkSecHubCreateSecretStoreData defines the data structure for creating a secret store in the Ark Secrets Hub.
type ArkSecHubCreateSecretStoreData struct {
	// AWS ASM Specific Fields
	AccountAlias string `json:"account_alias,omitempty" mapstructure:"account_alias,omitempty" flag:"aws-account-alias" desc:"AWS: The alias of your AWS account"`
	AccountID    string `json:"account_id,omitempty" mapstructure:"account_id,omitempty" flag:"aws-account-id" desc:"AWS: The 12-digit account ID of the AWS account that has the AWS Secrets Manager where you store secrets"`
	RegionID     string `json:"region_id,omitempty" mapstructure:"region_id,omitempty" flag:"aws-region-id" desc:"AWS: The region ID for the AWS Secrets Manager"`
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
	GcpProjectNumber          string `json:"gcp_project_number,omitempty" mapstructure:"gcp_project_number,omitempty" flag:"gcp-project-number" desc:"GCP: The number of the GCP project where the GCP Secret Manager is stored"`
	GcpWorkloadIdentityPoolID string `json:"gcp_workload_identity_pool_id,omitempty" mapstructure:"gcp_workload_identity_pool_id,omitempty" flag:"gcp-workload-identity-pool-id" desc:"GCP: The GCP workload identity pool ID created for Secrets Hub to access the GCP Secret Manager"`
	GcpPoolProviderID         string `json:"gcp_pool_provider_id,omitempty" mapstructure:"gcp_pool_provider_id,omitempty" flag:"gcp-pool-provider-id" desc:"GCP: The GCP pool provider ID created for Secrets Hub to access the GCP Secret Manager"`
	ServiceAccountEmail       string `json:"service_account_email,omitempty" mapstructure:"service_account_email,omitempty" flag:"gcp-service-account-email" desc:"GCP: The service account email created for Secrets Hub to access the GCP Secret Manager"`
	// Self-Hosted Specific Fields
	Password        string `json:"password,omitempty" mapstructure:"password,omitempty" desc:"SELF HOSTED: The password of the user in PAM 'SecretsHub'" flag:"sh-password"`
	URL             string `json:"url,omitempty" mapstructure:"url,omitempty" flag:"sh-url" desc:"SELF HOSTED: The URL of your PAM Self-Hosted PVWA, or the load balancer for the PVWA"`
	UserName        string `json:"username,omitempty" mapstructure:"username,omitempty" flag:"sh-username" desc:"SELF HOSTED: The user used for Secrets Hub to get secrets from PAM source. Should be 'SecretsHub'. This user should be created by REST API in PAM."`
	ConnectorID     string `json:"connector_id,omitempty" mapstructure:"connector_id,omitempty" desc:"SELF HOSTED: The connector unique identifier used to connect Secrets Hub and the Cloud Vendor." flag:"sh-connector-id"`
	ConnectorPoolID string `json:"connector_pool_id,omitempty" mapstructure:"connector_pool_id,omitempty" desc:"SELF HOSTED: The connector pool unique identifier used to connect PAM Self-Hosted and Secrets Hub." flag:"sh-connector-pool-id"`
	// Used by Azure, GCP
	ConnectionConfig *ArkSecHubCreateSecretStoreConnectionConfig `json:"connection_config,omitzero" mapstructure:"connection_config,omitzero" desc:"COMMON - AZURE, GCP: The network access configuration set for your target"`
}

// ArkSecHubCreateSecretStore defines the structure for creating a secret store in the Ark Secrets Hub.
type ArkSecHubCreateSecretStore struct {
	Type        string                         `json:"type" mapstructure:"type" flag:"type" validate:"required" desc:"The type for the secrets (AWS_ASM, AZURE_AKV,GCP_GSM,HASHI_HCV,PAM_PCLOUD,PAM_SELF_HOSTED)" choices:"AWS_ASM,AZURE_AKV,GCP_GSM,HASHI_HCV,PAM_PCLOUD,PAM_SELF_HOSTED"`
	Description string                         `json:"description,omitempty" mapstructure:"description,omitempty" flag:"description" desc:"A description of the secret store."`
	Name        string                         `json:"name" mapstructure:"name" desc:"The secret store name." flag:"name" validate:"required"`
	State       string                         `json:"state,omitempty" mapstructure:"state,omitempty" flag:"state" desc:"The secret store state (ENABLED,DISABLED)" default:"ENABLED" choices:"ENABLED,DISABLED"`
	Data        ArkSecHubCreateSecretStoreData `json:"data" mapstructure:"data" desc:"The data of the secret store depends on the secret store type."`
}
