package vm

// ArkSIAVMSecretsStats represents the statistics of secrets in a VM.
type ArkSIAVMSecretsStats struct {
	SecretsCount         int            `json:"secrets_count" mapstructure:"secrets_count" flag:"secrets-count" desc:"Overall secrets count"`
	ActiveSecretsCount   int            `json:"active_secrets_count" mapstructure:"active_secrets_count" flag:"active-secrets-count" desc:"Overall active secrets count"`
	InactiveSecretsCount int            `json:"inactive_secrets_count" mapstructure:"inactive_secrets_count" flag:"inactive-secrets-count" desc:"Overall inactive secrets count"`
	SecretsCountByType   map[string]int `json:"secrets_count_by_type" mapstructure:"secrets_count_by_type" flag:"secrets-count-by-type" desc:"Secrets count by type"`
}
