package sca

import "github.com/cyberark/ark-sdk-golang/pkg/models/services/uap/common"

// ArkUAPSCAConditions represents SCA-specific conditions.
// It is currently identical to ArkUAPConditions but defined separately for clarity
// and future extensibility without refactoring the ArkUAPSCACloudConsoleAccessPolicy model.
type ArkUAPSCAConditions struct {
	common.ArkUAPConditions `mapstructure:",squash"`
}
