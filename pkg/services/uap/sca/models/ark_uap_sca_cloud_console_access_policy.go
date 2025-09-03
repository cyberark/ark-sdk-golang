package models

import (
	uapcommonmodels "github.com/cyberark/ark-sdk-golang/pkg/services/uap/common/models"
)

// ArkUAPSCACloudConsoleAccessPolicy represents the access policy for the SCA Cloud Console.
type ArkUAPSCACloudConsoleAccessPolicy struct {
	uapcommonmodels.ArkUAPCommonAccessPolicy `mapstructure:",squash"`
	Conditions                               ArkUAPSCAConditions            `json:"conditions" mapstructure:"conditions" flag:"conditions" desc:"The time and session conditions of the policy"`
	Targets                                  ArkUAPSCACloudConsoleTarget    `json:"targets,omitempty" mapstructure:"targets,omitempty" flag:"targets" desc:"The targeted cloud provider and workspace"`
	InvalidResources                         ArkUAPSCACloudInvalidResources `json:"invalid_resources,omitempty" mapstructure:"invalid_resources,omitempty" flag:"invalid-resources" desc:"Resources that are not valid for the policy"`
}
