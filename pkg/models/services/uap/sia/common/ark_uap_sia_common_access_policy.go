package sia

import "github.com/cyberark/ark-sdk-golang/pkg/models/services/uap/common"

// ArkUAPSIACommonAccessPolicy represents a common access policy for SIA.
type ArkUAPSIACommonAccessPolicy struct {
	common.ArkUAPCommonAccessPolicy `mapstructure:",squash"`
	Conditions                      ArkUAPSIACommonConditions `json:"conditions" mapstructure:"conditions" flag:"conditions" desc:"The time, session, and idle time conditions of the policy"`
}
