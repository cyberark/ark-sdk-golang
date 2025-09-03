package models

// ArkCmgrUpdateNetwork is a struct representing the request to update a network in the Ark CMGR service.
type ArkCmgrUpdateNetwork struct {
	NetworkID string `json:"network_id" mapstructure:"network_id" flag:"network-id" desc:"ID of the network to update"`
	Name      string `json:"name,omitempty" mapstructure:"name,omitempty" flag:"name" desc:"New name of the network to update"`
}
