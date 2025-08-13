package cmgr

// ArkCmgrPoolIdentifiersFilter is a struct representing the filter for pool identifiers in the Ark CMGR service.
type ArkCmgrPoolIdentifiersFilter struct {
	ArkCmgrListPoolIdentifiers `mapstructure:",squash"`
	ArkCmgrPoolsCommonFilter   `mapstructure:",squash"`
}
