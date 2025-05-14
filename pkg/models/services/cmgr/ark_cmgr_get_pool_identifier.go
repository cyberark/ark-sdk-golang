package cmgr

// ArkCmgrGetPoolIdentifier is a struct representing the filter for getting a specific component in a pool in the Ark CMGR service.
type ArkCmgrGetPoolIdentifier struct {
	IdentifierID string `json:"identifier_id" mapstructure:"identifier_id" flag:"identifier-id" desc:"ID of the identifier to get from the pool"`
	PoolID       string `json:"pool_id" mapstructure:"pool_id" flag:"pool-id" desc:"ID of the pool to get"`
}
