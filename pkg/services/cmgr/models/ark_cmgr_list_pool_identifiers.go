package models

// ArkCmgrListPoolIdentifiers is a struct representing the filter for listing pool identifiers in the Ark CMGR service.
type ArkCmgrListPoolIdentifiers struct {
	PoolID string `json:"pool_id" mapstructure:"pool_id" flag:"pool-id" desc:"Pool id to get the identifiers for"`
}
