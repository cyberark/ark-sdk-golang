package models

// ArkSIASSOAcquireTokenResponse is a struct that represents the response from the Ark SIA SSO service for acquiring a token.
type ArkSIASSOAcquireTokenResponse struct {
	Token    map[string]interface{} `json:"token" validate:"required" mapstructure:"token"`
	Metadata map[string]interface{} `json:"metadata" validate:"required" mapstructure:"metadata"`
}
