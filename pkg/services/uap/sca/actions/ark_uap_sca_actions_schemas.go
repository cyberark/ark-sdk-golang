package actions

import (
	uapcommonmodels "github.com/cyberark/ark-sdk-golang/pkg/services/uap/common/models"
	uapscamodels "github.com/cyberark/ark-sdk-golang/pkg/services/uap/sca/models"
)

// ActionToSchemaMap defines the mapping of actions to schemas for the UAP SCA service.
var ActionToSchemaMap = map[string]interface{}{
	"add-policy":       &uapscamodels.ArkUAPSCACloudConsoleAccessPolicy{},
	"delete-policy":    &uapcommonmodels.ArkUAPDeletePolicyRequest{},
	"update-policy":    &uapscamodels.ArkUAPSCACloudConsoleAccessPolicy{},
	"policy":           &uapcommonmodels.ArkUAPGetPolicyRequest{},
	"list-policies":    nil,
	"list-policies-by": &uapscamodels.ArkUAPSCAFilters{},
	"policies-stats":   nil,
	"policy-status":    &uapcommonmodels.ArkUAPGetPolicyStatus{},
}
