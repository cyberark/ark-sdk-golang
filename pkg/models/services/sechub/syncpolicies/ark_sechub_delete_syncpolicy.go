package syncpolicies

// ArkSecHubDeleteSyncPolicy contains the policy id for the policy to delete
type ArkSecHubDeleteSyncPolicy struct {
	PolicyID string `json:"policy_id" mapstructure:"policy_id" desc:"Unique identifier of the referenced policy" flag:"policy-id" validate:"required"`
}
