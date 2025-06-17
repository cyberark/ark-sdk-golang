package syncpolicies

// ArkSecHubGetSyncPolicy contains the policy id for the policy to retrieve
type ArkSecHubGetSyncPolicy struct {
	PolicyID   string `json:"policy_id" mapstructure:"policy_id" description:"Unique identifier of the referenced policy" validate:"required"`
	Projection string `json:"projection" mapstructure:"projection" description:"Data representation method. Allowed values: 'EXTEND, REGULAR'" default:"REGULAR"`
}
