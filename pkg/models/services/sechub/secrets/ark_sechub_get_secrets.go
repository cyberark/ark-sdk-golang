package secrets

type ArkSecHubSecretVendorDataAnnotations struct {
}

type ArkSecHubSecretVendorDataTags struct {
	CyberArkPam      string `json:"cyberark_pam,omitempty" mapstructure:"cyberark_pam,omitempty" desc:""`
	CyberArkSafe     string ``
	CyberArkAccount  string ``
	CyberArkSecretID string ``
	SourceByCyberArk string ``
}

// ArkSecHubSecretVendorData represents the vendor-specific data for a secret.
type ArkSecHubSecretVendorData struct {
	CreatedAt             string                               `json:"created_at" mapstructure:"created_at"`
	UpdatedAt             string                               `json:"updated_at" mapstructure:"updated_at"`
	Enabled               bool                                 `json:"enabled" mapstructure:"enabled" desc:""`
	ProjectName           string                               `json:"project_name,omitempty" mapstructure:"project_name,omitempty" desc:""`
	ProjectNumber         string                               ``
	SecretEnabledVersions int                                  ``
	SecretType            string                               `json:"secret_type,omitempty" mapstructure:"secret_type,omitempty" desc:"E.g GLOBAL"`
	Tags                  ArkSecHubSecretVendorDataTags        `json:"tags" mapstructure:"tags"`
	Annotations           ArkSecHubSecretVendorDataAnnotations ``
	ReplicationMethod     string                               //`json:"secret_type,omitempty" mapstructure:"secret_type,omitempty" desc:"E.g GLOBAL"`
	LastRetrievedAt       string                               `json:"last_retrieved_at" mapstructure:"last_retrieved_at"`
	KmsKeyID              string                               `json:"kms_key_id,omitempty" mapstructure:"kms_key_id,omitempty"`
	AwsAccountID          string                               `json:"aws_account_id" mapstructure:"aws_account_id"`
	Region                string                               `json:"region" mapstructure:"region"`
}

// ArkSecHubSecret represents a single secret in the response.
type ArkSecHubSecret struct {
	VendorType       string                    `json:"vendor_type" mapstructure:"vendor_type" desc:"The vendor type of the store where the secret was found, valid values: AWS, AZURE, GCP" validate:"required"`
	VendorSubType    string                    `json:"vendor_sub_type" mapstructure:"vendor_sub_type" desc:"The subtype of the secret store where the secret was discovered, valid values: ASM, AKV, GSM" validate:"required"`
	ID               string                    `json:"id" mapstructure:"id" desc:"The unique identifier of the secret in Secrets Hub (internal). " validate:"required"`
	OriginID         string                    `json:"origin_id" mapstructure:"origin_id" desc:"The unique identifier of the secret as defined in the secret store." validate:"required"`
	Name             string                    `json:"name,omitempty" mapstructure:"name,omitempty" desc:"The name of the secret as defined in the secret store."`
	StoreID          string                    `json:"store_id" mapstructure:"store_id" desc:"The unique identifier of the secret store"`
	DiscoveredAt     string                    `json:"discovered_at" mapstructure:"discovered_at" desc:"The date and time that the secret was discovered by the Secrets Hub scan."`
	VendorData       ArkSecHubSecretVendorData `json:"vendor_data" mapstructure:"vendor_data" desc:"Data related to the secret as defined in the cloud platform."`
	LastScannedAt    string                    `json:"last_scanned_at" mapstructure:"last_scanned_at" desc:"The last date and time the secret was scanned by Secrets Hub, example: 2023-07-06T15:43:48.103000+00:00"`
	StoreName        string                    `json:"store_name,omitempty" mapstructure:"store_name,omitempty"`
	Onboarded        bool                      `json:"onboarded,omitempty" mapstructure:"onboarded,omitempty" desc:"Indicates whether the secret is onboarded to PAM."`
	SyncedByCyberArk bool                      `json:"synced_by_cyberark" mapstructure:"synced_by_cyberark" desc:""`
}

// ArkSecHubGetSecrets represents the response when requesting secrets from Ark Secrets Hub.
type ArkSecHubGetSecrets struct {
	Count      int               `json:"count,omitempty" mapstructure:"count,omitempty" desc:"Number of secrets in the result"`
	TotalCount int               `json:"total_count,omitempty" mapstructure:"total_count,omitempty" desc:"Number of secrets in all result sets"`
	Secrets    []ArkSecHubSecret `json:"secrets" mapstructure:"secrets" desc:"Secrets returned" validate:"required"`
}
