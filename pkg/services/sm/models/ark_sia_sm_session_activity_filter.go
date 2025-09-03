package models

// ArkSMSessionActivitiesFilter represents a filter for session activities based on session ID and command content.
type ArkSMSessionActivitiesFilter struct {
	SessionID       string `json:"session_id" mapstructure:"session_id" flag:"session-id" desc:"Session identifier to get" validate:"required"`
	CommandContains string `json:"command_contains" mapstructure:"command_contains" flag:"command-contains" desc:"String which the command contains"`
}
