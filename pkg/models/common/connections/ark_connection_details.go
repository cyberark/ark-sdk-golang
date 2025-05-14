package connections

// ArkConnectionType represents the type of connection.
type ArkConnectionType string

// ArkConnectionType values.
const (
	SSH   ArkConnectionType = "SSH"
	WinRM ArkConnectionType = "WinRM"
)

// ArkConnectionDetails represents the details of a connection.
type ArkConnectionDetails struct {
	Address           string                    `json:"address" mapstructure:"address"`
	Port              int                       `json:"port" mapstructure:"port"`
	ConnectionType    ArkConnectionType         `json:"connection_type" mapstructure:"connection_type"`
	Credentials       *ArkConnectionCredentials `json:"credentials" mapstructure:"credentials"`
	ConnectionData    interface{}               `json:"connection_data" mapstructure:"connection_data"`
	ConnectionRetries int                       `json:"connection_retries" mapstructure:"connection_retries"`
	RetryTickPeriod   int                       `json:"retry_tick_period" mapstructure:"retry_tick_period"`
}
