package models

// ArkSMSessionsStats represents statistics about sessions over the last 30 days.
type ArkSMSessionsStats struct {
	SessionsCount                   int                        `json:"sessions_count" mapstructure:"sessions_count" desc:"Sessions count in the last 30 days"`
	SessionsCountPerApplicationCode map[string]int             `json:"sessions_count_per_application_code" mapstructure:"sessions_count_per_application_code" desc:"Sessions count per application code"`
	SessionsCountPerPlatform        map[string]int             `json:"sessions_count_per_platform" mapstructure:"sessions_count_per_platform" desc:"Sessions count per platform"`
	SessionsCountPerStatus          map[ArkSMSessionStatus]int `json:"sessions_count_per_status" mapstructure:"sessions_count_per_status" desc:"Sessions count per status"`
	SessionsCountPerProtocol        map[string]int             `json:"sessions_count_per_protocol" mapstructure:"sessions_count_per_protocol" desc:"Sessions count per protocol"`
	SessionsFailureCount            int                        `json:"sessions_failure_count" mapstructure:"sessions_failure_count" desc:"Sessions count with failures"`
}
