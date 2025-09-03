package models

import (
	uapcommonmodels "github.com/cyberark/ark-sdk-golang/pkg/services/uap/common/models"
)

// ArkUAPSCAConditions represents SCA-specific conditions.
// It is currently identical to ArkUAPConditions but defined separately for clarity
// and future extensibility without refactoring the ArkUAPSCACloudConsoleAccessPolicy model.
type ArkUAPSCAConditions struct {
	uapcommonmodels.ArkUAPConditions `mapstructure:",squash"`
}
