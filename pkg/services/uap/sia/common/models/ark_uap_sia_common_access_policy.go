package models

import (
	uapcommonmodels "github.com/cyberark/ark-sdk-golang/pkg/services/uap/common/models"
)

// ArkUAPSIACommonAccessPolicy represents a common access policy for SIA.
type ArkUAPSIACommonAccessPolicy struct {
	uapcommonmodels.ArkUAPCommonAccessPolicy `mapstructure:",squash"`
	Conditions                               ArkUAPSIACommonConditions `json:"conditions" mapstructure:"conditions" flag:"conditions" desc:"The time, session, and idle time conditions of the policy"`
}
