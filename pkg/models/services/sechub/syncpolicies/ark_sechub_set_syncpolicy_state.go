package syncpolicies

type ArkSecHubSetSyncPolicyState struct {
	PolicyID string `json:"policy_id" mapstructure:"policy_id" desc:"Unique identifier of the sync policy" validate:"required" flag:"policy-id"`
	Action   string `json:"action" mapstructure:"action" desc:"The requested state for the policy - Allowed values: 'enable, disable'" validate:"required" default:"enable" flag:"action"` // enable, disable
}
