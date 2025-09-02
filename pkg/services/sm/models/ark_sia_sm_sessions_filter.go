package models

// ArkSMSessionsFilter represents a filter for querying sessions using a search expression.
type ArkSMSessionsFilter struct {
	// Search is a free text query to search sessions by.
	// Examples:
	// - 'duration LE 01:00:00'
	// - 'startTime GE 2023-11-18T06:53:30Z'
	// - 'status IN Failed,Ended AND endReason STARTSWITH Err008'
	// - 'command STARTSWITH ls'
	// - 'protocol IN SSH,RDP'
	Search string `json:"search" mapstructure:"search" flag:"search" desc:"Free text query to search sessions by. For example: 'startTime GE 2023-11-18T06:53:30Z AND status IN Failed,Ended AND endReason STARTSWITH Err008'" validate:"max=4096"`
}
