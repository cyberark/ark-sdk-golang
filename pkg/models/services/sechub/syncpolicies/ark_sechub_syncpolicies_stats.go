package syncpolicies

// ArkSecHubSyncPoliciesStats represents the statistics of sync policies in the Ark SecHub system.
// It includes the total count of secret stores, a breakdown by creator, and a breakdown by type.
type ArkSecHubSyncPoliciesStats struct {
	SyncPoliciesCount          int            `json:"sync_policies_count" mapstructure:"sync_policies_count" desc:"Overall sync policies count"`
	SyncPoliciesCountByCreator map[string]int `json:"sync_policies__count_by_creator" mapstructure:"sync_policies__count_by_creator" desc:"Sync Policies count by creator"`
}
