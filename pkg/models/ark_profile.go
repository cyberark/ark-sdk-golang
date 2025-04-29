package models

import (
	"encoding/json"
	"github.com/cyberark/ark-sdk-golang/pkg/models/auth"
)

// ArkProfile represents a profile for the Ark SDK.
type ArkProfile struct {
	ProfileName        string                          `json:"profile_name" mapstructure:"profile_name" validate:"required" flag:"profile-name" desc:"The name of the profile to use"`
	ProfileDescription string                          `json:"profile_description" mapstructure:"profile_description" validate:"required" flag:"profile-description" desc:"Profile Description"`
	AuthProfiles       map[string]*auth.ArkAuthProfile `json:"auth_profiles" mapstructure:"auth_profile" validate:"required" flag:"-"`
}

// UnmarshalJSON unmarshals the JSON data into the ArkProfile struct.
func (p *ArkProfile) UnmarshalJSON(data []byte) error {
	type Alias ArkProfile
	aux := &struct {
		AuthProfiles map[string]json.RawMessage `json:"auth_profiles"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	p.AuthProfiles = make(map[string]*auth.ArkAuthProfile)
	for key, rawMessage := range aux.AuthProfiles {
		var authProfile auth.ArkAuthProfile
		if err := json.Unmarshal(rawMessage, &authProfile); err != nil {
			return err
		}
		p.AuthProfiles[key] = &authProfile
	}

	return nil
}
