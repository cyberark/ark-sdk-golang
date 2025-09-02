package models

// ArkCmgrNetworksStats is a struct representing the statistics of networks in the Ark CMGR service.
type ArkCmgrNetworksStats struct {
	NetworksCount        int            `json:"networks_count" mapstructure:"networks_count" flag:"networks-count" desc:"Overall count of network"`
	PoolsCountPerNetwork map[string]int `json:"pools_count_per_network" mapstructure:"pools_count_per_network" flag:"pools-count-per-network" desc:"Count of pools for each network"`
}
