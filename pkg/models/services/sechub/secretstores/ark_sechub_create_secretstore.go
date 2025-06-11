package secretstores

type ArkSecHubCreateSecretStoreData struct {
	// AWS ASM Specific Fields
	AccountAlias string `json:"account_alias,omitempty" mapstructure:"account_alias,omitempty" flag:"aws-account-alias" desc:"The alias of your AWS account"`
	AccountID    string `json:"account_id,omitempty" mapstructure:"account_id,omitempty" flag:"aws-account-id" desc:"The 12-digit account ID of the AWS account that has the AWS Secrets Manager where you store secrets"`
	RegionID     string `json:"region_id,omitempty" mapstructure:"region_id,omitempty" flag:"aws-region-id" desc:"The region ID where the AWS account is managed"`
	// Common Fields
	// Used by AWS and HashiCorp Vault
	RoleName string `json:"role_name,omitempty" mapstructure:"role_name,omitempty" flag:"rolename" desc:"The role used to allow Secrets Hub to manage secrets in your AWS Secrets Manager or Hashi Vault"`
}

type ArkSecHubCreateSecretStore struct {
	Target      string                         `json:"target" mapstructure:"target" flag:"target" validate:"required" desc:"The target for the secrets, either where the secrets are scanned or where the secrets are syncing to. Valid values: AWS_ASM, AZURE_AKV, GCP_GSM, HASHI_HCV, PAM_PCLOUD, PAM_SELF_HOSTED"`
	Source      string                         `json:"source,omitempty" mapstructure:"source,omitempty" flag:"source" desc:"The source of the secrets, meaning where the secrets are syncing from. Valid values: PAM_PCLOUD, PAM_SELF_HOSTED"`
	Description string                         `json:"description,omitempty" mapstructure:"description,omitempty" flag:"description" desc:"A description of the secret store."`
	Name        string                         `json:"name" mapstructure:"name" desc:"The secret store name." flag:"name" validate:"required"`
	State       string                         `json:"state,omitempty" mapstructure:"state,omitempty" flag:"state" desc:"The secret store state. Valid Values: ENABLED, DISABLED. Default Value: ENABLED"`
	Data        ArkSecHubCreateSecretStoreData `json:"data" mapstructure:"data" desc:"The data of the secret store depends on the secret store type."`
}
