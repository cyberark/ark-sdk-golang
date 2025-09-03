package models

// ArkCmgrPoolsStats is a struct representing the statistics of pools in the Ark CMGR service.
type ArkCmgrPoolsStats struct {
	PoolsCount              int                       `json:"pools_count" mapstructure:"pools_count" flag:"pools-count" desc:"Overall count of pools"`
	NetworksCountPerPool    map[string]int            `json:"networks_count_per_pool" mapstructure:"networks_count_per_pool" flag:"networks-count-per-pool" desc:"Count of networks for each pool"`
	IdentifiersCountPerPool map[string]int            `json:"identifiers_count_per_pool" mapstructure:"identifiers_count_per_pool" flag:"identifiers-count-per-pool" desc:"Count of identifiers for each pool"`
	ComponentsCountPerPool  map[string]map[string]int `json:"components_count_per_pool" mapstructure:"components_count_per_pool" flag:"components-count-per-pool" desc:"Count of components for each pool"`
}
