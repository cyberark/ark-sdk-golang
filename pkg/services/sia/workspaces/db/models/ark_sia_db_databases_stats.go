package models

// ArkSIADBDatabasesStats represents the statistics of database in a workspace.
type ArkSIADBDatabasesStats struct {
	DatabasesCount                    int            `json:"databases_count" mapstructure:"databases_count" desc:"Databases overall count"`
	DatabasesWithoutSecretCount       int            `json:"databases_without_secret_count" mapstructure:"databases_without_secret_count" desc:"Databases who does not have any secret attached to them"`
	DatabasesWithoutCertificatesCount int            `json:"databases_without_certificates_count" mapstructure:"databases_without_certificates_count" desc:"Databases who does not have any certificates attached to them"`
	DatabasesCountByEngine            map[string]int `json:"databases_count_by_engine" mapstructure:"databases_count_by_engine" desc:"Databases count per engine type"`
	DatabasesCountByFamily            map[string]int `json:"databases_count_by_family" mapstructure:"databases_count_by_family" desc:"Databases count per family type"`
	DatabasesCountByWorkspace         map[string]int `json:"databases_count_by_workspace" mapstructure:"databases_count_by_workspace" desc:"Databases count per workspace type"`
	DatabasesCountByAuthMethod        map[string]int `json:"databases_count_by_auth_method" mapstructure:"databases_count_by_auth_method" desc:"Databases count per auth type"`
	DatabasesCountByWarning           map[string]int `json:"databases_count_by_warning" mapstructure:"databases_count_by_warning" desc:"Databases count per warning type"`
}
