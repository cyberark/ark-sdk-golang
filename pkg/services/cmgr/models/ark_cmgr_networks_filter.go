package models

// ArkCmgrNetworksFilter is a struct representing the filter for networks in the Ark CMGR service.
type ArkCmgrNetworksFilter struct {
	ArkCmgrPoolsCommonFilter `mapstructure:",squash"`
}
