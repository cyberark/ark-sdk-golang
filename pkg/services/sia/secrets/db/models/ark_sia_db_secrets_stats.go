package models

// ArkSIADBSecretsStats represents the statistics of secrets in the Ark SIA DB.
type ArkSIADBSecretsStats struct {
	SecretsCount             int            `json:"secrets_count" mapstructure:"secrets_count" desc:"Overall secrets count"`
	ActiveSecretsCount       int            `json:"active_secrets_count" mapstructure:"active_secrets_count" desc:"Overall active secrets count"`
	InactiveSecretsCount     int            `json:"inactive_secrets_count" mapstructure:"inactive_secrets_count" desc:"Overall inactive secrets count"`
	SecretsCountBySecretType map[string]int `json:"secrets_count_by_secret_type" mapstructure:"secrets_count_by_secret_type" desc:"Secrets count by secret type"`
	SecretsCountByStoreType  map[string]int `json:"secrets_count_by_store_type" mapstructure:"secrets_count_by_store_type" desc:"Secrets count by store type"`
}
