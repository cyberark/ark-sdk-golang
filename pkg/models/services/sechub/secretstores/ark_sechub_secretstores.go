package secretstores

// ArkSecHubSecretStoreConnectionConfig defines the connection configuration for a secret store in the Ark Secrets Hub.
type ArkSecHubSecretStoreConnectionConfig struct {
	ConnectionType string `json:"connection_type,omitempty" mapstructure:"connection_type,omitempty" desc:"The connection type (CONNECTOR,PUBLIC)"`
	// Required if you choose 'CONNECTOR' as the connection type.
	// If you choose 'PUBLIC', these fields are not required.
	ConnectorID     string `json:"connector_id,omitempty" mapstructure:"connector_id,omitempty" desc:"The connector unique identifier used to connect Secrets Hub and the Cloud Vendor. Example: ManagementAgent_90c63827-7315-4284-8559-ac8d24f2666d"`
	ConnectorPoolID string `json:"connector_pool_id,omitempty" mapstructure:"connector_pool_id,omitempty" desc:"The connector pool unique identifier used to connect PAM Self-Hosted and Secrets Hub.. Example: c389961d-a0cd-46ab-9f69-877f756a59c1"`
}

// ArkSecHubSecretStoreData defines the data structure for a secret store in the Ark Secrets Hub.
type ArkSecHubSecretStoreData struct {
	// AWS ASM Specific Fields
	AccountAlias string `json:"account_alias,omitempty" mapstructure:"account_alias,omitempty" flag:"aws-account-alias"`
	AccountID    string `json:"account_id,omitempty" mapstructure:"account_id,omitempty"`
	RegionID     string `json:"region_id,omitempty" mapstructure:"region_id,omitempty"`
	// GCP GSM Specific Fields
	GcpProjectName            string `json:"gcp_project_name,omitempty" mapstructure:"gcp_project_name,omitempty" desc:"The name of the GCP project where the GCP Secret Manager is stored"`
	GcpProjectNumber          string `json:"gcp_project_number,omitempty" mapstructure:"gcp_project_number,omitempty" desc:"The number of the GCP project where the GCP Secret Manager is stored"`
	GcpWorkloadIdentityPoolID string `json:"gcp_workload_identity_pool_id,omitempty" mapstructure:"gcp_workload_identity_pool_id,omitempty" desc:"The GCP workload identity pool ID created for Secrets Hub to access the GCP Secret Manager"`
	GcpPoolProviderID         string `json:"gcp_pool_provider_id,omitempty" mapstructure:"gcp_pool_provider_id,omitempty" desc:"The GCP pool provider ID created for Secrets Hub to access the GCP Secret Manager"`
	ServiceAccountEmail       string `json:"service_account_email,omitempty" mapstructure:"service_account_email,omitempty" desc:"The service account email created for Secrets Hub to access the GCP Secret Manager"`
	// Hashi Vault Specific Fields
	HashiVaultURL    string `json:"hashi_vault_url,omitempty" mapstructure:"hashi_vault_url,omitempty" desc:"The URL of the HashiCorp Vault"`
	EnginePath       string `json:"engine_path,omitempty" mapstructure:"engine_path,omitempty" desc:"The path of the engine in HashiCorp Vault"`
	EngineType       string `json:"engine_type,omitempty" mapstructure:"engine_type,omitempty" desc:"The type of the engine in HashiCorp Vault. Valid values: KV, PKI, SSH"`
	EngineAPIVersion string `json:"engine_api_version,omitempty" mapstructure:"engine_api_version,omitempty" desc:"The API version of the engine in HashiCorp Vault. Valid values: 1, 2"`
	// Privilege Cloud and Self-Hosted Specific Fields
	URL             string `json:"url,omitempty" mapstructure:"url,omitempty"`
	UserName        string `json:"user_name,omitempty" mapstructure:"user_name,omitempty" desc:"The user used for Secrets Hub to get secrets from PAM source. Should be 'SecretsHub'. This user should be created by REST API in PAM"`
	ConnectorID     string `json:"connector_id,omitempty" mapstructure:"connector_id,omitempty" desc:"The connector unique identifier used to connect Secrets Hub and the Cloud Vendor. Example: ManagementAgent_90c63827-7315-4284-8559-ac8d24f2666d"`
	ConnectorPoolID string `json:"connector_pool_id,omitempty" mapstructure:"connector_pool_id,omitempty" desc:"The connector pool unique identifier used to connect PAM Self-Hosted and Secrets Hub.. Example: c389961d-a0cd-46ab-9f69-877f756a59c1"`
	// Azure AKV Specific Fields
	AppClientDirectoryID string `json:"app_client_directory_id,omitempty" mapstructure:"app_client_directory_id,omitempty"`
	AzureVaultURL        string `json:"azure_vault_url,omitempty" mapstructure:"azure_vault_url,omitempty"`
	AppClientID          string `json:"app_client_id,omitempty" mapstructure:"app_client_id,omitempty"`
	SubscriptionID       string `json:"subscription_id,omitempty" mapstructure:"subscription_id,omitempty"`
	SubscriptionName     string `json:"subscription_name,omitempty" mapstructure:"subscription_name,omitempty"`
	ResourceGroupName    string `json:"resource_group_name,omitempty" mapstructure:"resource_group_name,omitempty"`
	// Common Fields
	// Used by AWS and HashiCorp Vault
	RoleName string `json:"role_name,omitempty" mapstructure:"role_name,omitempty"`
	// Used by Azure, GCP, and HashiCorp Vault
	ConnectionConfig *ArkSecHubSecretStoreConnectionConfig `json:"connection_config,omitempty" mapstructure:"connection_config,omitempty"`
}

// ArkSecHubSecretStoreScan defines the scan status of a secret store in the Ark Secrets Hub.
type ArkSecHubSecretStoreScan struct {
	ID         string `json:"id" mapstructure:"id" validate:"required" desc:"The unique identifier of the scan"`
	Status     string `json:"status" mapstructure:"status" validate:"required" desc:"The status of the scan (IN_PROGRESS,SUCCESS,FAILED)"`
	Message    string `json:"message,omitempty" mapstructure:"message,omitempty" desc:"More information on the scan status."`
	FinishedAt string `json:"finished_at,omitempty" mapstructure:"finished_at,omitempty" desc:"The date and time the scan ended. Example: 2023-07-06T15:45:00.103000"`
}

// ArkSecHubSecretStore defines the structure for a secret store in the Ark Secrets Hub.
type ArkSecHubSecretStore struct {
	ID                 string                   `json:"id" mapstructure:"id" desc:"The unique identifier of the secret store" validate:"required"`
	Type               string                   `json:"type" mapstructure:"type" desc:"The type of secret store (PAM_PCLOUD,PAM_SELF_HOSTED,AWS_ASM,AZURE_AKV,GCP_GSM,HASHI_HCV)" validate:"required"`
	Behaviors          []string                 `json:"behaviors" mapstructure:"behaviors" desc:"Whether the secret store is used as a source or a target. There can be only one source secret store per tenant. Valid values: SECRETS_SOURCE, SECRETS_TARGET"`
	CreatedAt          string                   `json:"created_at" mapstructure:"created_at" desc:"The secret store creation date." validate:"required"`
	CreatedBy          string                   `json:"created_by" mapstructure:"created_by" desc:"The user who created the secret store." validate:"required"`
	Data               ArkSecHubSecretStoreData `json:"data" mapstructure:"data" desc:"Data related to the secret store as defined in the cloud platform." validate:"required"`
	Description        string                   `json:"description,omitempty" mapstructure:"description,omitempty" desc:"A description of the secret store."`
	Name               string                   `json:"name" mapstructure:"name" desc:"The secret store name." validate:"required"`
	UpdatedAt          string                   `json:"updated_at" mapstructure:"updated_at" desc:"The last date the secret store was updated" validate:"required"`
	UpdatedBy          string                   `json:"updated_by" mapstructure:"updated_by" desc:"The last user to update the secret store." validate:"required"`
	CreationDetails    string                   `json:"creation_details,omitempty" mapstructure:"creation_details,omitempty" desc:"Allowed Values: Secrets Hub, Connect Cloud Environment"`
	OrganizationID     string                   `json:"organization_id,omitempty" mapstructure:"organization_id,omitempty"`
	Scan               ArkSecHubSecretStoreScan `json:"scan" mapstructure:"scan"`
	TotalPoliciesCount int                      `json:"total_policies_count,omitempty" mapstructure:"total_policies_count,omitempty" desc:"The total amount of policies in the secret store"`
	TotalSecretsCount  int                      `json:"total_secrets_count,omitempty" mapstructure:"total_secrets_count,omitempty" desc:"The total amount of secrets in the secret store"`
}
