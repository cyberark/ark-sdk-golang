package cmgr

// ArkCmgrGetNetwork is a struct representing the filter for getting a specific network in the Ark CMGR service.
type ArkCmgrGetNetwork struct {
	NetworkID string `json:"network_id" mapstructure:"network_id" flag:"network-id" desc:"ID of the network to get"`
}
