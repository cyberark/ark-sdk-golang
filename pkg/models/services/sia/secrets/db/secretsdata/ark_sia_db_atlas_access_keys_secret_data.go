package secretsdata

// ArkSIADBAtlasAccessKeysSecretData represents the Mongo Atlas access keys secret data in the Ark SIA DB.
type ArkSIADBAtlasAccessKeysSecretData struct {
	ArkSIADBSecretData
	PublicKey  string                 `json:"public_key" mapstructure:"public_key" desc:"Public part of mongo atlas access keys"`
	PrivateKey string                 `json:"private_key" mapstructure:"private_key" desc:"Private part of mongo atlas access keys"`
	Metadata   map[string]interface{} `json:"metadata,omitempty" mapstructure:"metadata" desc:"Extra secret details"`
}

// GetDataSecretType returns the secret type of the secret data.
func (s *ArkSIADBAtlasAccessKeysSecretData) GetDataSecretType() string {
	return "atlas_access_keys"
}

// ArkSIADBExposedAtlasAccessKeysSecretData represents the exposed Mongo Atlas access keys secret data in the Ark SIA DB.
type ArkSIADBExposedAtlasAccessKeysSecretData struct {
	ArkSIADBSecretData
	PublicKey string `json:"public_key" mapstructure:"public_key" desc:"Public part of mongo atlas access keys"`
}

// GetDataSecretType returns the secret type of the secret data.
func (s *ArkSIADBExposedAtlasAccessKeysSecretData) GetDataSecretType() string {
	return "atlas_access_keys"
}
