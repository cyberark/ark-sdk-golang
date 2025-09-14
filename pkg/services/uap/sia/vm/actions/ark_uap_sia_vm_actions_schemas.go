package actions

import (
	uapcommonmodels "github.com/cyberark/ark-sdk-golang/pkg/services/uap/common/models"
	uapsiavmmodels "github.com/cyberark/ark-sdk-golang/pkg/services/uap/sia/vm/models"
)

// ActionToSchemaMap defines the mapping of actions to schemas for the UAP SIA VM service.
var ActionToSchemaMap = map[string]interface{}{
	"add-policy":       &uapsiavmmodels.ArkUAPSIAVMAccessPolicy{},
	"delete-policy":    &uapcommonmodels.ArkUAPDeletePolicyRequest{},
	"update-policy":    &uapsiavmmodels.ArkUAPSIAVMAccessPolicy{},
	"policy":           &uapcommonmodels.ArkUAPGetPolicyRequest{},
	"list-policies":    nil,
	"list-policies-by": &uapsiavmmodels.ArkUAPSIAVMFilters{},
	"policies-stats":   nil,
	"policy-status":    &uapcommonmodels.ArkUAPGetPolicyStatus{},
}
