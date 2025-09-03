package models

import (
	uapcommonmodels "github.com/cyberark/ark-sdk-golang/pkg/services/uap/common/models"
)

// ArkUAPSIACommonConditions represents common conditions for SIA policies.
type ArkUAPSIACommonConditions struct {
	uapcommonmodels.ArkUAPConditions `mapstructure:",squash"`
	IdleTime                         int `json:"idle_time,omitempty" mapstructure:"idle_time,omitempty" flag:"idle-time" desc:"The maximum idle time before the session ends, in minutes." validate:"gt=0,lte=120" default:"10"`
}
