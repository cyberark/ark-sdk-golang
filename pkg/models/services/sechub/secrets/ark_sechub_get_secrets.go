package secrets

// ArkSecHubSecretVendorDataReplicas defines the structure for the replicas field in the vendor data of a secret.
type ArkSecHubSecretVendorDataReplicas struct {
	Location string `json:"location,omitempty" mapstructure:"location,omitempty" desc:"Locations secrets are replicated to"`
}

// ArkSecHubSecretVendorData represents the vendor-specific data for a secret.
type ArkSecHubSecretVendorData struct {
	// AWS specific vendor data
	AwsAccountID string `json:"aws_account_id,omitempty" mapstructure:"aws_account_id,omitempty" desc:""`
	// Azure specific vendor data
	SubscriptionID    string `json:"subscription_id,omitempty" mapstructure:"subscription_id,omitempty" desc:"The ID of the Azure subscription associated with the Azure Key Vault."`
	SubsciptionName   string `json:"subscription_name,omitempty" mapstructure:"subscription_name,omitempty" desc:"The display name of the Azure subscription associated with the Azure Key Vault."`
	ResourceGroupName string `json:"resource_group_name,omitempty" mapstructure:"resource_group_name,omitempty" desc:"The name of the resource group to which the Azure Key Vault belongs."`
	NotBefore         string `json:"not_before,omitempty" mapstructure:"not_before,omitempty" desc:"The date from which the secret is valid and can be used, as defined in Azure Key Vault."`
	// GCP specific vendor data
	ProjectName           string                            `json:"project_name,omitempty" mapstructure:"project_name,omitempty" desc:"The name of the Google Cloud project where the secret is stored."`
	ProjectNumber         string                            `json:"project_number,omitempty" mapstructure:"project_number,omitempty" desc:"The number of the Google Cloud project where the secret is stored."`
	SecretEnabledVersions int                               `json:"secret_enabled_versions,omitempty" mapstructure:"secret_enabled_versions,omitempty" desc:"The number of versions of the secret that are enabled in external secret store."`
	SecretType            string                            `json:"secret_type,omitempty" mapstructure:"secret_type,omitempty" desc:"(GLOBAL,REGIONAL)" choices:"GLOBAL,REGIONAL"`
	Replicas              ArkSecHubSecretVendorDataReplicas `json:"replicas,omitzero" mapstructure:"replicas,omitzero" desc:"If the secret is being replicated, and to where."`
	NextRotationTime      string                            `json:"next_rotation_time,omitempty" mapstructure:"next_rotation_time,omitempty" desc:"The time at which the secret is scheduled to be rotated by the external secret store."`
	Annotations           map[string]string                 `json:"annotations,omitempty" mapstructure:"" desc:"The annotations that are applied to the secret in external secret store."`
	ReplicationMethod     string                            `json:"replication_method,omitempty" mapstructure:"replication_method,omitempty" desc:"The method used to replicate the secret in external secret store."`
	// Shared or common objects
	Tags            map[string]string `json:"tags,omitzero" mapstructure:"tags,omitzero" desc:"The tags that are applied to the secret in the secret store."`
	CreatedAt       string            `json:"created_at" mapstructure:"created_at" desc:"The date and time the secret was created in the secret store."`
	Enabled         bool              `json:"enabled" mapstructure:"enabled" desc:"Indicates whether the secret is enabled, as defined in the external secret store."`
	ExpiresAt       string            `json:"expires_at,omitempty" mapstructure:"expires_at,omitempty" desc:"The date the secret expires and can no longer be used, as defined in AKV or GSM."`
	KmsKeyID        string            `json:"kms_key_id,omitempty" mapstructure:"kms_key_id,omitempty" desc:"The ID of the encryption key used to encrypt the secret value in the external secret store."`
	UpdatedAt       string            `json:"updated_at" mapstructure:"updated_at" desc:"The date and time the secret was last updated."`
	LastRetrievedAt string            `json:"last_retrieved_at" mapstructure:"last_retrieved_at" desc:"The last date and time the secret was retrieved from the secret store."`
	Region          string            `json:"region,omitempty" mapstructure:"region,omitempty" desc:"Cloud Service Provider Region"`
}

// ArkSecHubSecret represents a single secret in the response.
type ArkSecHubSecret struct {
	VendorType       string                    `json:"vendor_type" mapstructure:"vendor_type" desc:"The vendor type of the store where the secret was found (AWS, AZURE, GCP)" validate:"required"`
	VendorSubType    string                    `json:"vendor_sub_type" mapstructure:"vendor_sub_type" desc:"The subtype of the secret store where the secret was discovered (ASM, AKV, GSM)" validate:"required"`
	ID               string                    `json:"id" mapstructure:"id" desc:"The unique identifier of the secret in Secrets Hub (internal). " validate:"required"`
	OriginID         string                    `json:"origin_id" mapstructure:"origin_id" desc:"The unique identifier of the secret as defined in the secret store." validate:"required"`
	Name             string                    `json:"name,omitempty" mapstructure:"name,omitempty" desc:"The name of the secret as defined in the secret store."`
	StoreID          string                    `json:"store_id" mapstructure:"store_id" desc:"The unique identifier of the secret store"`
	DiscoveredAt     string                    `json:"discovered_at" mapstructure:"discovered_at" desc:"The date and time that the secret was discovered by the Secrets Hub scan."`
	VendorData       ArkSecHubSecretVendorData `json:"vendor_data,omitzero" mapstructure:"vendor_data,omitzero" desc:"Data related to the secret as defined in the cloud platform."`
	LastScannedAt    string                    `json:"last_scanned_at,omitempty" mapstructure:"last_scanned_at,omitempty" desc:"The last date and time the secret was scanned by Secrets Hub, example: 2023-07-06T15:43:48.103000+00:00"`
	StoreName        string                    `json:"store_name,omitempty" mapstructure:"store_name,omitempty" desc:"Name of the secret store"`
	Onboarded        bool                      `json:"onboarded,omitempty" mapstructure:"onboarded,omitempty" desc:"Indicates whether the secret is onboarded to PAM."`
	SyncedByCyberArk bool                      `json:"synced_by_cyberark" mapstructure:"synced_by_cyberark" desc:"Indicates whether the secret has been synced by CyberArk. If not set, the status is unknown."`
}
