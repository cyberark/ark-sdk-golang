package secretsdata

// ArkSIADBUserPasswordSecretData represents the user password secret data in the Ark SIA DB.
type ArkSIADBUserPasswordSecretData struct {
	ArkSIADBSecretData
	Username string                 `json:"username,omitempty" mapstructure:"username" desc:"Name or id of the user"`
	Password string                 `json:"password,omitempty" mapstructure:"password" desc:"Password of the user"`
	Metadata map[string]interface{} `json:"metadata,omitempty" mapstructure:"metadata" desc:"Extra secret details"`
}

// GetDataSecretType returns the secret type of the secret data.
func (s *ArkSIADBUserPasswordSecretData) GetDataSecretType() string {
	return "username_password"
}

// ArkSIADBExposedUserPasswordSecretData represents the exposed user password secret data in the Ark SIA DB.
type ArkSIADBExposedUserPasswordSecretData struct {
	ArkSIADBSecretData
	Username string `json:"username,omitempty" mapstructure:"username" desc:"Name or id of the user"`
}

// GetDataSecretType returns the secret type of the secret data.
func (s *ArkSIADBExposedUserPasswordSecretData) GetDataSecretType() string {
	return "username_password"
}
