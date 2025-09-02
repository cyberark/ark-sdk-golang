package models

// ArkSIAInstallConnector represents the details required to install a connector.
type ArkSIAInstallConnector struct {
	ConnectorType      string `json:"connector_type" mapstructure:"connector_type" flag:"connector-type" desc:"The type of the platform for the connector to be installed in (ON-PREMISE,AWS,AZURE,GCP)" default:"ON-PREMISE" choices:"ON-PREMISE,AWS,AZURE,GCP"`
	ConnectorOS        string `json:"connector_os" mapstructure:"connector_os" flag:"connector-os" desc:"The type of the operating system for the connector to be installed on (linux,windows)" default:"linux" choices:"linux,windows"`
	ConnectorPoolID    string `json:"connector_pool_id" mapstructure:"connector_pool_id" flag:"connector-pool-id" desc:"The connector pool which the connector will be part of, if not given, the connector will be assigned to the default one" validate:"required"`
	TargetMachine      string `json:"target_machine" mapstructure:"target_machine" desc:"Target machine on which to install the connector on"`
	Username           string `json:"username" mapstructure:"username" desc:"Username to connect with to the target machine"`
	Password           string `json:"password,omitempty" mapstructure:"password" desc:"Password to connect with to the target machine"`
	PrivateKeyPath     string `json:"private_key_path,omitempty" mapstructure:"private_key_path" desc:"Private key file path to use for connecting to the target machine via ssh"`
	PrivateKeyContents string `json:"private_key_contents,omitempty" mapstructure:"private_key_contents" desc:"Private key contents to use for connecting to the target machine via ssh"`
	RetryCount         int    `json:"retry_count" mapstructure:"retry_count" flag:"retry-count" desc:"Number of times to retry to connect if it fails" default:"10"`
	RetryDelay         int    `json:"retry_delay" mapstructure:"retry_delay" flag:"retry-delay" desc:"Delay in seconds between retries" default:"5"`
}
