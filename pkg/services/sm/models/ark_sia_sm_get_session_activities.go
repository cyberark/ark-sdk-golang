package models

// ArkSIASMGetSessionActivities represents the request to get a session activities by ID.
type ArkSIASMGetSessionActivities struct {
	SessionID string `json:"session_id" mapstructure:"session_id" flag:"session-id" desc:"Session identifier to get the activities for'" validate:"required"`
}
