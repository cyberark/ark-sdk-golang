package secretsdata

// ArkSIADBSecretData represents a secret data in the Ark SIA DB.
type ArkSIADBSecretData interface {
	// GetDataSecretType returns the secret type of the secret data.
	GetDataSecretType() string
}
