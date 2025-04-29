package connections

// ArkConnectionCommand represents a command to be executed on a remote server.
type ArkConnectionCommand struct {
	Command          string                 `json:"command" mapstructure:"command"`                       // The command to actually run
	ExpectedRC       int                    `json:"expected_rc" mapstructure:"expected_rc"`               // Expected return code
	ExtraCommandData map[string]interface{} `json:"extra_command_data" mapstructure:"extra_command_data"` // Extra data for the command
}
