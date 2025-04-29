package safes

// ArkPCloudSafesStats represents statistics about safes.
type ArkPCloudSafesStats struct {
	SafesCount           int            `json:"safes_count" mapstructure:"safes_count" desc:"Overall safes count"`
	SafesCountByLocation map[string]int `json:"safes_count_by_location" mapstructure:"safes_count_by_location" desc:"Safes count by locations"`
	SafesCountByCreator  map[string]int `json:"safes_count_by_creator" mapstructure:"safes_count_by_creator" desc:"Safes count by creator"`
}
