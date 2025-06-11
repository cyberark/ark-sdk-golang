package secretstores

type ArkSecHubUpdateSecretStore struct {
	SecretStoreID string `json:"secret_store_id" mapstructure:"secret_store_id" flag:"secret-store-id" validate:"required" desc:"The unique identifier of the secret store to update"`
	Description   string `json:"description,omitempty" mapstructure:"description,omitempty" flag:"description" desc:"A description of the secret store."`
	Name          string `json:"name,omitempty" mapstructure:"name,omitempty" flag:"name" desc:"The name of the secret store. It should be unique per tenant."`
	// Data contains the specific data for the secret store type.
	Data *ArkSecHubUpdateSecretStoreData `json:"data,omitempty" mapstructure:"data,omitempty" desc:"The data related to the secret store as defined in the cloud platform."`
}

type ArkSecHubUpdateSecretStoreData struct {
	// AWS ASM Specific Fields
	AccountAlias string `json:"account_alias,omitempty" mapstructure:"account_alias,omitempty" flag:"aws-account-alias"`
	// GCP GSM Specific Fields
	GcpProjectName            string `json:"gcp_project_name,omitempty" mapstructure:"gcp_project_name,omitempty" flag:"gcp-project-name" desc:"The name of the GCP project where the GCP Secret Manager is stored"`
	GcpWorkloadIdentityPoolId string `json:"gcp_workload_identity_pool_id,omitempty" mapstructure:"gcp_workload_identity_pool_id,omitempty" desc:"The GCP workload identity pool ID created for Secrets Hub to access the GCP Secret Manager"`
	GcpPoolProviderId         string `json:"gcp_pool_provider_id,omitempty" mapstructure:"gcp_pool_provider_id,omitempty" flag:"gcp-pool-provider-id" desc:"The GCP pool provider ID created for Secrets Hub to access the GCP Secret Manager"`
	ServiceAccountEmail       string `json:"service_account_email,omitempty" mapstructure:"service_account_email,omitempty" flag:"gcp-service-account-email" desc:"The service account email created for Secrets Hub to access the GCP Secret Manager"`
	// Self-Hosted Specific Fields
	Password        string `json:"password,omitempty" mapstructure:"password,omitempty" desc:"The password of the user in PAM 'SecretsHub'" flag:"password"`
	ConnectorID     string `json:"connector_id,omitempty" mapstructure:"connector_id,omitempty" desc:"The connector unique identifier used to connect Secrets Hub and the Cloud Vendor. Example: ManagementAgent_90c63827-7315-4284-8559-ac8d24f2666d" flag:"sh-connector-id"`
	ConnectorPoolID string `json:"connector_pool_id,omitempty" mapstructure:"connector_pool_id,omitempty" desc:"The connector pool unique identifier used to connect PAM Self-Hosted and Secrets Hub.. Example: c389961d-a0cd-46ab-9f69-877f756a59c1" flag:"sh-connector-pool-id"`
	// Azure AKV Specific Fields
	AppClientDirectoryId string `json:"app_client_directory_id,omitempty" mapstructure:"app_client_directory_id,omitempty" flag:"azure-app-client-directory-id" desc:"The Azure directory/tenant ID where the application (user) for Secrets Hub was created"`
	AppClientId          string `json:"app_client_id,omitempty" mapstructure:"app_client_id,omitempty" flag:"azure-app-client-id" desc:"A unique Application (client) ID assigned to Secrets Hub by Azure AD when the app was registered."`
	AppClientSecret      string `json:"app_client_secret,omitempty" mapstructure:"app_client_secret,omitempty" flag:"azure-app-client-secret" desc:"The user's password that will be used by Secrets Hub to access the Azure Key Vault."`
	SubscriptionId       string `json:"subscription_id,omitempty" mapstructure:"subscription_id,omitempty" flag:"azure-subscription-id" desc:"The Azure subscription ID linked to the Azure Key Vault"`
	SubscriptionName     string `json:"subscription_name,omitempty" mapstructure:"subscription_name,omitempty" flag:"azure-subscription-name" desc:"The Azure subscription name linked to the Azure Key Vault"`
	ResourceGroupName    string `json:"resource_group_name,omitempty" mapstructure:"resource_group_name,omitempty" flag:"azure-resource-group-name" desc:"The Azure resource group name where the Azure Key Vault is stored"`
	// Common Fields
	// Used by AWS and HashiCorp Vault
	RoleName string `json:"role_name,omitempty" mapstructure:"role_name,omitempty" flag:"rolename" desc:"Rolename for AWS and Hashi"`
	// Used by Azure, GCP, and HashiCorp Vault
	ConnectionConfig *ArkSecHubUpdateSecretStoreConnectionConfig `json:"connection_config,omitempty" mapstructure:"connection_config,omitempty" desc:"The network access configuration set for your target"`
}

type ArkSecHubUpdateSecretStoreConnectionConfig struct {
	ConnectionType string `json:"connection_type,omitempty" mapstructure:"connection_type,omitempty" flag:"connection-type" desc:"If your Cloud Vault is not open to public access, choose 'CONNECTOR'. Valid Values: 'CONNECTOR','PUBLIC"`
	// Required if you choose 'CONNECTOR' as the connection type.
	// If you choose 'PUBLIC', these fields are not required.
	ConnectorID     string `json:"connector_id,omitempty" mapstructure:"connector_id,omitempty" flag:"connector-id" desc:"The connector unique identifier used to connect Secrets Hub and the Cloud Vendor. Example: ManagementAgent_90c63827-7315-4284-8559-ac8d24f2666d"`
	ConnectorPoolID string `json:"connector_pool_id,omitempty" mapstructure:"connector_pool_id,omitempty" flag:"connector-pool-id" desc:"The connector pool unique identifier used to connect PAM Self-Hosted and Secrets Hub.. Example: c389961d-a0cd-46ab-9f69-877f756a59c1"`
}
