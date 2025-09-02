package models

// ArkCmgrUpdatePool is a struct representing the request to update a pool in the Ark CMGR service.
type ArkCmgrUpdatePool struct {
	PoolID             string   `json:"pool_id" mapstructure:"pool_id" flag:"pool-id" desc:"ID of the pool to update"`
	Name               string   `json:"name,omitempty" mapstructure:"name,omitempty" flag:"name" desc:"Name of the pool to update"`
	Description        string   `json:"description,omitempty" mapstructure:"description,omitempty" flag:"description" desc:"Pool description to update"`
	AssignedNetworkIDs []string `json:"assigned_network_ids,omitempty" mapstructure:"assigned_network_ids,omitempty" flag:"assigned-network-ids" desc:"Assigned networks to the pool to update"`
}
