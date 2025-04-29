package sso

// ArkSIASSOGetSSHKey is a struct that represents the request for getting SSH key from the Ark SIA SSO service.
type ArkSIASSOGetSSHKey struct {
	Folder string `json:"folder" mapstructure:"folder" flag:"folder" desc:"Output folder to write the ssh key to" default:"~/.ssh"`
}
