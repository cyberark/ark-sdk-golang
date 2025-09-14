package actions

import uapcommonmodels "github.com/cyberark/ark-sdk-golang/pkg/services/uap/common/models"

// ActionToSchemaMap defines the mapping of actions to schemas for the UAP service.
var ActionToSchemaMap = map[string]interface{}{
	"policies-stats":   nil,
	"list-policies":    nil,
	"list-policies-by": &uapcommonmodels.ArkUAPFilters{},
	"policy-status":    &uapcommonmodels.ArkUAPGetPolicyStatus{},
}
