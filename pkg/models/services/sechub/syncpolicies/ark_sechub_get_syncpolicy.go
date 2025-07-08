package syncpolicies

// ArkSecHubGetSyncPolicy contains the policy id for the policy to retrieve
type ArkSecHubGetSyncPolicy struct {
	PolicyID   string `json:"policy_id" mapstructure:"policy_id" desc:"Unique identifier of the referenced policy" flag:"policy-id" validate:"required"`
	Projection string `json:"projection" mapstructure:"projection" desc:"Data representation method (EXTEND, REGULAR)" default:"REGULAR" flag:"projection" choices:"EXTEND,REGULAR"`
}
