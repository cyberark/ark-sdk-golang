package cmgr

// ArkCmgrUpdatePoolIdentifier is a struct representing the filter for updating identifiers in a pool in the Ark CMGR service.
type ArkCmgrUpdatePoolIdentifier struct {
	Type   string `json:"type" mapstructure:"type" flag:"type" desc:"Type of identifier to update"`
	Value  string `json:"value" mapstructure:"value" flag:"value" desc:"Value of the identifier"`
	ID     string `json:"id" mapstructure:"id" flag:"id" desc:"ID of the identifier to update from the pool"`
	PoolID string `json:"pool_id" mapstructure:"pool_id" flag:"pool-id" desc:"ID of the pool to update the identifier to"`
}
