package models

// ArkSIADBBaseExecution defines the base structure for database execution parameters
type ArkSIADBBaseExecution struct {
	TargetAddress  string `json:"target_address" mapstructure:"target_address" flag:"target-address" desc:"Target address to connect to"`
	TargetUsername string `json:"target_username,omitempty" mapstructure:"target_username,omitempty" flag:"target-username" desc:"Target username account to use"`
	NetworkName    string `json:"network_name,omitempty" mapstructure:"network_name,omitempty" flag:"network-name" desc:"Network name to use for the connection, if applicable"`
}
