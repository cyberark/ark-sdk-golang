package auth

import (
	"encoding/json"
	"fmt"
)

// ArkAuthProfile represents the authentication profile for Ark SIA.
type ArkAuthProfile struct {
	Username           string                `json:"username" mapstructure:"username" flag:"username" desc:"Username"`
	AuthMethod         ArkAuthMethod         `json:"auth_method" mapstructure:"auth_method" flag:"-"`
	AuthMethodSettings ArkAuthMethodSettings `json:"auth_method_settings" mapstructure:"auth_method_settings" flag:"-"`
}

// UnmarshalJSON unmarshals the JSON data into the ArkAuthProfile struct.
func (a *ArkAuthProfile) UnmarshalJSON(data []byte) error {
	type Alias ArkAuthProfile
	aux := &struct {
		AuthMethodSettings json.RawMessage `json:"auth_method_settings"`
		*Alias
	}{
		Alias: (*Alias)(a),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var settings ArkAuthMethodSettings
	switch a.AuthMethod {
	case Identity:
		settings = &IdentityArkAuthMethodSettings{}
	case IdentityServiceUser:
		settings = &IdentityServiceUserArkAuthMethodSettings{}
	case Direct:
		settings = &DirectArkAuthMethodSettings{}
	case Default:
		settings = &DefaultArkAuthMethodSettings{}
	default:
		return fmt.Errorf("unknown auth method: %s", a.AuthMethod)
	}

	if err := json.Unmarshal(aux.AuthMethodSettings, settings); err != nil {
		return err
	}

	a.AuthMethodSettings = settings
	return nil
}
