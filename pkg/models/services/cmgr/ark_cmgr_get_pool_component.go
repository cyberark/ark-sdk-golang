package cmgr

// ArkCmgrGetPoolComponent is a struct representing the filter for getting a specific component in a pool in the Ark CMGR service.
type ArkCmgrGetPoolComponent struct {
	PoolID      string `json:"pool_id" mapstructure:"pool_id" flag:"pool-id" desc:"ID of the pool to get"`
	ComponentID string `json:"component_id" mapstructure:"component_id" flag:"component-id" desc:"ID of the component to get in the pool"`
}
