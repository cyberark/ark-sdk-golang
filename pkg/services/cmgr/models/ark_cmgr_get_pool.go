package models

// ArkCmgrGetPool is a struct representing the filter for getting a specific pool in the Ark CMGR service.
type ArkCmgrGetPool struct {
	PoolID string `json:"pool_id" mapstructure:"pool_id" flag:"pool-id" desc:"ID of the pool to get"`
}
