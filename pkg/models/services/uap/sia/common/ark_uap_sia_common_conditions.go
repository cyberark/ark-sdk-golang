package sia

import "github.com/cyberark/ark-sdk-golang/pkg/models/services/uap/common"

// ArkUAPSIACommonConditions represents common conditions for SIA policies.
type ArkUAPSIACommonConditions struct {
	common.ArkUAPConditions `mapstructure:",squash"`
	IdleTime                int `json:"idle_time,omitempty" mapstructure:"idle_time,omitempty" flag:"idle-time" desc:"The maximum idle time before the session ends, in minutes." validate:"gt=0,lte=120" default:"10"`
}
