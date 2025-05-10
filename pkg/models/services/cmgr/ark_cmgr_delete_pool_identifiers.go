package cmgr

// ArkCmgrDeletePoolIdentifier is a struct representing the filter for deleting a specific identifier from a pool in the Ark CMGR service.
type ArkCmgrDeletePoolIdentifier struct {
	IdentifierID string `json:"identifier_id" mapstructure:"identifier_id" flag:"identifier-id" desc:"ID of the identifier to delete"`
}

// ArkCmgrDeletePoolSingleIdentifier is a struct representing the filter for deleting a single identifier from a pool in the Ark CMGR service.
type ArkCmgrDeletePoolSingleIdentifier struct {
	ID     string `json:"id" mapstructure:"id" flag:"id" desc:"ID of the identifier to delete"`
	PoolID string `json:"pool_id" mapstructure:"pool_id" flag:"pool-id" desc:"ID of the pool to delete the identifier from"`
}

// ArkCmgrDeletePoolBulkIdentifier is a struct representing the filter for deleting multiple identifiers from a pool in the Ark CMGR service.
type ArkCmgrDeletePoolBulkIdentifier struct {
	PoolID      string                        `json:"pool_id" mapstructure:"pool_id" flag:"pool-id" desc:"ID of the pool to delete the identifiers from"`
	Identifiers []ArkCmgrDeletePoolIdentifier `json:"identifiers" mapstructure:"identifiers" flag:"identifiers" desc:"Identifiers to delete"`
}
