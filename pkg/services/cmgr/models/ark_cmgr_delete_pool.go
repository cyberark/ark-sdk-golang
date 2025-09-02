package models

// ArkCmgrDeletePool is a struct representing the filter for deleting a specific pool in the Ark CMGR service.
type ArkCmgrDeletePool struct {
	PoolID string `json:"pool_id" mapstructure:"pool_id" flag:"pool-id" desc:"ID of the pool to delete"`
}
