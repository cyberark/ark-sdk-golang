package syncpolicies

// ArkSecHubSetSyncPolicyState defines the structure for setting the state of a sync policy in the Ark Secrets Hub.
type ArkSecHubSetSyncPolicyState struct {
	PolicyID string `json:"policy_id" mapstructure:"policy_id" desc:"Unique identifier of the sync policy" validate:"required" flag:"policy-id"`
	Action   string `json:"action" mapstructure:"action" desc:"The requested state for the policy (Allowed values: 'enable', 'disable')" validate:"required" default:"enable" flag:"action" choices:"enable,disable"`
}
