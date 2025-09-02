package secretlinks

// ArkSIADBSecretLink represents a secret link in the Ark SIA DB.
type ArkSIADBSecretLink interface {
	// GetLinkSecretType returns the secret type of the secret link.
	GetLinkSecretType() string
}
