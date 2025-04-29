package cmgr

// ArkCmgrDeleteNetwork is a struct representing the filter for deleting a network in the Ark CMGR service.
type ArkCmgrDeleteNetwork struct {
	NetworkID string `json:"network_id" mapstructure:"network_id" flag:"network-id" desc:"ID of the network to delete"`
}
