package cmgr

// ArkCmgrNetworkPool is a struct representing a network pool in the Ark CMGR service.
type ArkCmgrNetworkPool struct {
	ID   string `json:"id" mapstructure:"id" flag:"id" desc:"ID of the pool"`
	Name string `json:"name" mapstructure:"name" flag:"name" desc:"Name of the pool"`
}

// ArkCmgrNetwork is a struct representing a network in the Ark CMGR service.
type ArkCmgrNetwork struct {
	ID            string               `json:"id" mapstructure:"id" flag:"id" desc:"ID of the network"`
	Name          string               `json:"name" mapstructure:"name" flag:"name" desc:"Name of the network"`
	AssignedPools []ArkCmgrNetworkPool `json:"assigned_pools,omitempty" mapstructure:"assigned_pools,omitempty" flag:"assigned-pools" desc:"Assigned pools on this network"`
	CreatedAt     string               `json:"created_at" mapstructure:"created_at" flag:"created-at" desc:"The creation time of the network"`
	UpdatedAt     string               `json:"updated_at" mapstructure:"updated_at" flag:"updated-at" desc:"The last update time of the network"`
}
