package db

// ArkSIADBSecretsFilter represents the filters for querying secrets in the Ark SIA DB.
type ArkSIADBSecretsFilter struct {
	SecretName string            `json:"secret_name,omitempty" mapstructure:"secret_name" flag:"secret-name" desc:"Filter by secret name"`
	SecretType string            `json:"secret_type,omitempty" mapstructure:"secret_type" flag:"secret-type" choices:"username_password,iam_user,cyberark_pam,atlas_access_keys" desc:"Filter by type"`
	StoreType  string            `json:"store_type,omitempty" mapstructure:"store_type" flag:"store-type" desc:"Filter by store type"`
	IsActive   bool              `json:"is_active,omitempty" mapstructure:"is_active" flag:"is-active" desc:"Filter by if secret is active"`
	Tags       map[string]string `json:"tags,omitempty" mapstructure:"tags" flag:"tags" desc:"Filter by tags"`
}
