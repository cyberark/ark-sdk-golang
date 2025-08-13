package sca

import "github.com/cyberark/ark-sdk-golang/pkg/models/services/uap/common"

// ArkUAPSCAFilters represents filters specific to the SCA (Security Cloud Access) policies
// within the UAP (Unified Access Policies) service.
//
// You can set the following fields:
//
//   - TargetCategory: []common.ArkCategoryType
//     A list of target categories to filter the policies by.
//
//   - PolicyType: []common.ArkUAPPolicyType
//     A list of policy types to filter the policies by.
//
//   - PolicyTags: []string
//     A list of policy tags to filter the policies by.
//
//   - Identities: []string
//     A list of identities to filter the policies by.
//
//   - Status: []common.ArkUAPStatusType
//     A list of policy statuses to filter the policies by.
//
//   - TextSearch: *string
//     A text value to apply as a search filter across policies.
//
//   - ShowEditablePolicies: *bool
//     Whether to show only policies that are editable by the current user.
type ArkUAPSCAFilters struct {
	common.ArkUAPFilters `mapstructure:",squash"`
}
