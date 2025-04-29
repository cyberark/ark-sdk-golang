package access

// ArkSIAUninstallConnector represents the details required to install a connector.
type ArkSIAUninstallConnector struct {
	ConnectorOS        string `json:"connector_os" mapstructure:"connector_os" flag:"connector-os" desc:"The type of the operating system for the connector to uninstall (linux,windows)" default:"linux" choices:"linux,windows"`
	ConnectorID        string `json:"connector_id" mapstructure:"connector_id" flag:"connector-id" desc:"The connector ID to be uninstalled" validate:"required"`
	TargetMachine      string `json:"target_machine" mapstructure:"target_machine" desc:"Target machine on which to uninstall the connector on"`
	Username           string `json:"username" mapstructure:"username" desc:"Username to connect with to the target machine"`
	Password           string `json:"password,omitempty" mapstructure:"password" desc:"Password to connect with to the target machine"`
	PrivateKeyPath     string `json:"private_key_path,omitempty" mapstructure:"private_key_path" desc:"Private key file path to use for connecting to the target machine via ssh"`
	PrivateKeyContents string `json:"private_key_contents,omitempty" mapstructure:"private_key_contents" desc:"Private key contents to use for connecting to the target machine via ssh"`
}
