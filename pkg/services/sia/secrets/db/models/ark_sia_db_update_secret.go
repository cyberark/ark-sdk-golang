package models

// ArkSIADBUpdateSecret is the struct for updating a secret in the Ark SIA DB.
type ArkSIADBUpdateSecret struct {
	SecretID      string            `json:"secret_id,omitempty" mapstructure:"secret_id" flag:"secret-id" desc:"Secret id to update"`
	SecretName    string            `json:"secret_name,omitempty" mapstructure:"secret_name" flag:"secret-name" desc:"Name of the secret to update"`
	NewSecretName string            `json:"new_secret_name,omitempty" mapstructure:"new_secret_name" flag:"new-secret-name" desc:"New secret name to update to"`
	Description   string            `json:"description,omitempty" mapstructure:"description" flag:"description" desc:"Description about the secret to update"`
	Purpose       string            `json:"purpose,omitempty" mapstructure:"purpose" flag:"purpose" desc:"Purpose of the secret to update"`
	Tags          map[string]string `json:"tags,omitempty" mapstructure:"tags" flag:"tags" desc:"Tags of the secret to change to"`

	// Username Password Secret Type
	Username string `json:"username,omitempty" mapstructure:"username" flag:"username" desc:"Name or id of the user for username_password type"`
	Password string `json:"password,omitempty" mapstructure:"password" flag:"password" desc:"Password of the user for username_password type"`

	// PAM Account Secret Type
	PAMSafe        string `json:"pam_safe,omitempty" mapstructure:"pam_safe" flag:"pam-safe" desc:"Safe of the account for pam_account type"`
	PAMAccountName string `json:"pam_account_name,omitempty" mapstructure:"pam_account_name" flag:"pam-account-name" desc:"Account name for pam_account type"`

	// IAM Secret Type
	IAMAccount         string `json:"iam_account,omitempty" mapstructure:"iam_account" flag:"iam-account" desc:"Account number of the iam user"`
	IAMUsername        string `json:"iam_username,omitempty" mapstructure:"iam_username" flag:"iam-username" desc:"Username portion in the ARN of the iam user"`
	IAMAccessKeyID     string `json:"iam_access_key_id,omitempty" mapstructure:"iam_access_key_id" flag:"iam-access-key-id" desc:"Access key id of the user"`
	IAMSecretAccessKey string `json:"iam_secret_access_key,omitempty" mapstructure:"iam_secret_access_key" flag:"iam-secret-access-key" desc:"Secret access key of the user"`

	// Atlas Secret Type
	AtlasPublicKey  string `json:"atlas_public_key,omitempty" mapstructure:"atlas_public_key" flag:"atlas-public-key" desc:"Public part of mongo atlas access keys"`
	AtlasPrivateKey string `json:"atlas_private_key,omitempty" mapstructure:"atlas_private_key" flag:"atlas-private-key" desc:"Private part of mongo atlas access keys"`
}
