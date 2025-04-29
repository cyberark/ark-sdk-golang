package accounts

// ArkPCloudAccountSecretVersion represents the version details of an account secret.
type ArkPCloudAccountSecretVersion struct {
	IsTemporary      bool   `json:"is_temporary" mapstructure:"is_temporary" desc:"Whether the secret is permanent or temporary" flag:"is-temporary"`
	ModificationDate int    `json:"modification_date" mapstructure:"modification_date" desc:"Modification time of the secret" flag:"modification-date"`
	ModifiedBy       string `json:"modified_by" mapstructure:"modified_by" desc:"Username who modified the secret" flag:"modified-by"`
	VersionID        int    `json:"version_id" mapstructure:"version_id" desc:"Version ID of the secret" flag:"version-id"`
}
