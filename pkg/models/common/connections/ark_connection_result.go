package connections

// ArkConnectionResult represents the result of a connection attempt.
type ArkConnectionResult struct {
	Stdout string `json:"stdout" mapstructure:"stdout"`
	Stderr string `json:"stderr" mapstructure:"stderr"`
	RC     int    `json:"rc" mapstructure:"rc"`
}
