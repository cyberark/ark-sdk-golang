package auth

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/common"
)

// ArkTokenType is a string type that represents the type of token used in the Ark SIA.
type ArkTokenType string

// Toke types supported by Ark
const (
	JWT      ArkTokenType = "JSON Web Token"
	Cookies  ArkTokenType = "Cookies"
	Token    ArkTokenType = "Token"
	Password ArkTokenType = "Password"
	Custom   ArkTokenType = "Custom"
	Internal ArkTokenType = "Internal"
)

// ArkSecret is a struct that represents a secret in the Ark SIA.
type ArkSecret struct {
	Secret string `json:"secret"`
}

// ArkToken is a struct that represents a token in the Ark SIA.
type ArkToken struct {
	Token        string                 `json:"token" mapstructure:"token" validate:"required"`
	TokenType    ArkTokenType           `json:"token_type" mapstructure:"token_type" validate:"required"`
	Username     string                 `json:"username" mapstructure:"username"`
	Endpoint     string                 `json:"endpoint" mapstructure:"endpoint"`
	AuthMethod   ArkAuthMethod          `json:"auth_method" mapstructure:"auth_method"`
	ExpiresIn    common.ArkRFC3339Time  `json:"expires_in" mapstructure:"expires_in"`
	RefreshToken string                 `json:"refresh_token" mapstructure:"refresh_token"`
	Metadata     map[string]interface{} `json:"metadata" mapstructure:"metadata"`
}
