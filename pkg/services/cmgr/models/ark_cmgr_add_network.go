package models

// ArkCmgrAddNetwork is a struct representing the filter for adding a network in the Ark CMGR service.
type ArkCmgrAddNetwork struct {
	Name string `json:"name" mapstructure:"name" flag:"name" desc:"Name of the network to add" required:"true"`
}
