package actions

import (
	uapcommonmodels "github.com/cyberark/ark-sdk-golang/pkg/services/uap/common/models"
	uapsiadbmodels "github.com/cyberark/ark-sdk-golang/pkg/services/uap/sia/db/models"
)

// ActionToSchemaMap defines the mapping of actions to schemas for the UAP SIA DB service.
var ActionToSchemaMap = map[string]interface{}{
	"add-policy":       &uapsiadbmodels.ArkUAPSIADBAccessPolicy{},
	"delete-policy":    &uapcommonmodels.ArkUAPDeletePolicyRequest{},
	"update-policy":    &uapsiadbmodels.ArkUAPSIADBAccessPolicy{},
	"policy":           &uapcommonmodels.ArkUAPGetPolicyRequest{},
	"list-policies":    nil,
	"list-policies-by": &uapsiadbmodels.ArkUAPSIADBFilters{},
	"policies-stats":   nil,
	"policy-status":    &uapcommonmodels.ArkUAPGetPolicyStatus{},
}
