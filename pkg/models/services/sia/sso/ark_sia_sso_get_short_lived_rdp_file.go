package sso

// ArkSIASSOGetShortLivedRDPFile is a struct that represents the request for getting a short-lived RDP file from the Ark SIA SSO service.
type ArkSIASSOGetShortLivedRDPFile struct {
	AllowCaching       bool   `json:"allow_caching" mapstructure:"allow_caching" flag:"allow-caching" desc:"Allow short lived token caching" default:"false"`
	Folder             string `json:"folder" validate:"required" mapstructure:"folder" flag:"folder" desc:"Output folder to write the rdp file to"`
	TargetAddress      string `json:"target_address" validate:"required" mapstructure:"target_address"`
	TargetDomain       string `json:"target_domain" mapstructure:"target_domain" flag:"target-domain" desc:"Target domain to use for the rdp file"`
	TargetUser         string `json:"target_user" mapstructure:"target_user" flag:"target-user" desc:"Target user to use for the rdp file"`
	ElevatedPrivileges bool   `json:"elevated_privileges" mapstructure:"elevated_privileges" flag:"elevated-privileges" desc:"Whether to use elevated privileges or not"`
}
