package services

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	uapcommonmodels "github.com/cyberark/ark-sdk-golang/pkg/services/uap/common/models"
	scamodels "github.com/cyberark/ark-sdk-golang/pkg/services/uap/sca/models"
	dbmodels "github.com/cyberark/ark-sdk-golang/pkg/services/uap/sia/db/models"
	vmmodels "github.com/cyberark/ark-sdk-golang/pkg/services/uap/sia/vm/models"
)

// UAPSCAActionToSchemaMap defines the mapping of actions to schemas for the UAP SCA service.
var UAPSCAActionToSchemaMap = map[string]interface{}{
	"add-policy":       &scamodels.ArkUAPSCACloudConsoleAccessPolicy{},
	"delete-policy":    &uapcommonmodels.ArkUAPDeletePolicyRequest{},
	"update-policy":    &scamodels.ArkUAPSCACloudConsoleAccessPolicy{},
	"policy":           &uapcommonmodels.ArkUAPGetPolicyRequest{},
	"list-policies":    nil,
	"list-policies-by": &scamodels.ArkUAPSCAFilters{},
	"policies-stats":   nil,
	"policy-status":    &uapcommonmodels.ArkUAPGetPolicyStatus{},
}

// UAPSCAActions defines the actions for the UAP SCA service.
var UAPSCAActions = actions.ArkServiceActionDefinition{
	ActionName: "sca",
	Schemas:    UAPSCAActionToSchemaMap,
}

// UAPSIADBActionToSchemaMap defines the mapping of actions to schemas for the UAP SIA DB service.
var UAPSIADBActionToSchemaMap = map[string]interface{}{
	"add-policy":       &dbmodels.ArkUAPSIADBAccessPolicy{},
	"delete-policy":    &uapcommonmodels.ArkUAPDeletePolicyRequest{},
	"update-policy":    &dbmodels.ArkUAPSIADBAccessPolicy{},
	"policy":           &uapcommonmodels.ArkUAPGetPolicyRequest{},
	"list-policies":    nil,
	"list-policies-by": &dbmodels.ArkUAPSIADBFilters{},
	"policies-stats":   nil,
	"policy-status":    &uapcommonmodels.ArkUAPGetPolicyStatus{},
}

// UAPSIADBAction defines the actions for the UAP SIA DB service.
var UAPSIADBAction = actions.ArkServiceActionDefinition{
	ActionName: "db",
	Schemas:    UAPSIADBActionToSchemaMap,
}

// UAPSIAVMActionToSchemaMap defines the mapping of actions to schemas for the UAP SIA VM service.
var UAPSIAVMActionToSchemaMap = map[string]interface{}{
	"add-policy":       &vmmodels.ArkUAPSIAVMAccessPolicy{},
	"delete-policy":    &uapcommonmodels.ArkUAPDeletePolicyRequest{},
	"update-policy":    &vmmodels.ArkUAPSIAVMAccessPolicy{},
	"policy":           &uapcommonmodels.ArkUAPGetPolicyRequest{},
	"list-policies":    nil,
	"list-policies-by": &vmmodels.ArkUAPSIAVMFilters{},
	"policies-stats":   nil,
	"policy-status":    &uapcommonmodels.ArkUAPGetPolicyStatus{},
}

// UAPSIAVMAction defines the actions for the UAP SIA VM service.
var UAPSIAVMAction = actions.ArkServiceActionDefinition{
	ActionName: "vm",
	Schemas:    UAPSIAVMActionToSchemaMap,
}

// UAPActionToSchemaMap defines the mapping of actions to schemas for the UAP service.
var UAPActionToSchemaMap = map[string]interface{}{
	"policies-stats":   nil,
	"list-policies":    nil,
	"list-policies-by": &uapcommonmodels.ArkUAPFilters{},
	"policy-status":    &uapcommonmodels.ArkUAPGetPolicyStatus{},
}

// UAPActions defines the actions for the UAP service, including subactions for SCA and SIA DB.
var UAPActions = &actions.ArkServiceActionDefinition{
	ActionName: "uap",
	Schemas:    UAPActionToSchemaMap,
	Subactions: []*actions.ArkServiceActionDefinition{
		&UAPSCAActions,
		&UAPSIADBAction,
		&UAPSIAVMAction,
	},
}
