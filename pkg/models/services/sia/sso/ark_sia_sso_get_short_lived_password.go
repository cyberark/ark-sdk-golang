package sso

// ArkSIASSOGetShortLivedPassword is a struct that represents the request for getting a short-lived password from the Ark SIA SSO service.
type ArkSIASSOGetShortLivedPassword struct {
	AllowCaching bool `json:"allow_caching" mapstructure:"allow_caching" flag:"allow-caching" desc:"Allow short lived token caching" default:"false"`
}
