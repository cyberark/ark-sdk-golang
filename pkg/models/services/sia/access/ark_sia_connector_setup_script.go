package access

// ArkSIAConnectorSetupScript represents the setup script details for the SIA connector.
type ArkSIAConnectorSetupScript struct {
	ScriptURL string `json:"script_url" mapstructure:"script_url" desc:"URL to manually download the SIA connector installation script. The script contains a secret token that is valid for 15 minutes from the time it is generated."`
	BashCmd   string `json:"bash_cmd" mapstructure:"bash_cmd" desc:"Bash command to automatically download the installation script and run it on the connector host machine. The script contains a secret token that is valid for 15 minutes from the time it is generated."`
}
