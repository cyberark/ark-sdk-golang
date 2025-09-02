package models

// ArkCmgrAddPool is a struct representing the filter for adding a pool in the Ark CMGR service.
type ArkCmgrAddPool struct {
	Name               string   `json:"name" mapstructure:"name" flag:"name" desc:"Name of the pool to add" required:"true"`
	Description        string   `json:"description,omitempty" mapstructure:"description,omitempty" flag:"description" desc:"Pool description"`
	AssignedNetworkIDs []string `json:"assigned_network_ids" mapstructure:"assigned_network_ids" flag:"assigned-network-ids" desc:"Assigned networks to the pool"`
}
