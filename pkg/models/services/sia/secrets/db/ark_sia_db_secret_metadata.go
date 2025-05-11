package db

// ArkSIADBSecretMetadata represents the metadata of a secret in the Ark SIA DB.
type ArkSIADBSecretMetadata struct {
	SecretID          string                  `json:"secret_id" mapstructure:"secret_id" desc:"Secret identifier"`
	SecretName        string                  `json:"secret_name" mapstructure:"secret_name" desc:"Name of the secret"`
	Description       string                  `json:"description,omitempty" mapstructure:"description" desc:"Description about the secret"`
	Purpose           string                  `json:"purpose,omitempty" mapstructure:"purpose" desc:"Purpose of the secret"`
	SecretType        string                  `json:"secret_type" mapstructure:"secret_type" desc:"Type of the secret" choices:"username_password,iam_user,cyberark_pam,atlas_access_keys"`
	SecretStore       ArkSIADBStoreDescriptor `json:"secret_store" mapstructure:"secret_store" desc:"Secret store details of the secret"`
	SecretLink        map[string]interface{}  `json:"secret_link,omitempty" mapstructure:"secret_link" desc:"Link details of the secret"`
	SecretExposedData map[string]interface{}  `json:"secret_exposed_data,omitempty" mapstructure:"secret_exposed_data" desc:"Portion of the secret data which can be exposed to the user"`
	Tags              map[string]string       `json:"tags,omitempty" mapstructure:"tags" desc:"Tags of the secret"`
	CreatedBy         string                  `json:"created_by" mapstructure:"created_by" desc:"Who created the secret"`
	CreationTime      string                  `json:"creation_time" mapstructure:"creation_time" desc:"Creation time of the secret"`
	LastUpdatedBy     string                  `json:"last_updated_by" mapstructure:"last_updated_by" desc:"Who last updated the secret"`
	LastUpdateTime    string                  `json:"last_update_time" mapstructure:"last_update_time" desc:"When was the secret last updated"`
	IsActive          bool                    `json:"is_active" mapstructure:"is_active" desc:"Whether the secret is active or not"`
}

// ArkSIADBSecretMetadataList represents a list of secret metadata in the Ark SIA DB.
type ArkSIADBSecretMetadataList struct {
	TotalCount int                      `json:"total_count" mapstructure:"total_count" desc:"Total secrets found"`
	Secrets    []ArkSIADBSecretMetadata `json:"secrets" mapstructure:"secrets" desc:"Actual secrets metadata"`
}
