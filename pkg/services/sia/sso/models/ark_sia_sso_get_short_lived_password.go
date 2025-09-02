package models

// ArkSIASSOGetShortLivedPassword is a struct that represents the request for getting a short-lived password from the Ark SIA SSO service.
type ArkSIASSOGetShortLivedPassword struct {
	AllowCaching bool   `json:"allow_caching" mapstructure:"allow_caching" flag:"allow-caching" desc:"Allow short lived token caching" default:"false"`
	Service      string `json:"service" mapstructure:"service" flag:"service" desc:"Which service to get the token info for" choices:"DPA-DB,DPA-RDP" default:"DPA-DB"`
}
