package models

// ArkSecHubSecretsStats represents the response when getting secrets statistics from SecHub.
type ArkSecHubSecretsStats struct {
	SecretsCount                    int            `json:"secrets_count" mapstructure:"secrets_count" desc:"Overall secrets count"`
	SecretsCountByVendorType        map[string]int `json:"secrets_count_by_vendor_type" mapstructure:"secrets_count_by_vendor_type" desc:"Secrets count by vendor type"`
	SecretsCountByStoreName         map[string]int `json:"secrets_count_by_store_name" mapstructure:"secrets_count_by_store_name" desc:"Secrets count by store name"`
	SecretsCountSyncedByCyberArk    int            `json:"secrets_count_synced_by_cyberark" mapstructure:"secrets_count_synced_by_cyberark" desc:"Secrets count for secrets synced by CyberArk"`
	SecretsCountNotSyncedByCyberArk int            `json:"secrets_count_not_synced_by_cyberark" mapstructure:"secrets_count_not_synced_by_cyberark" desc:"Secrets count for secrets not synced by CyberArk"`
}
