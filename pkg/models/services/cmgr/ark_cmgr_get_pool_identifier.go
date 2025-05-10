package cmgr

// ArkCmgrGetPoolIdentifier is a struct representing the filter for getting a specific component in a pool in the Ark CMGR service.
type ArkCmgrGetPoolIdentifier struct {
	ID     string `json:"id" mapstructure:"id" flag:"id" desc:"ID of the identifier to get from the pool"`
	PoolID string `json:"pool_id" mapstructure:"pool_id" flag:"pool-id" desc:"ID of the pool to get"`
}
