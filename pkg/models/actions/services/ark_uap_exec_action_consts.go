package services

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	"github.com/cyberark/ark-sdk-golang/pkg/models/services/uap/common"
	"github.com/cyberark/ark-sdk-golang/pkg/models/services/uap/sca"
	"github.com/cyberark/ark-sdk-golang/pkg/models/services/uap/sia/db"
	"github.com/cyberark/ark-sdk-golang/pkg/models/services/uap/sia/vm"
)

// UAPSCAActionToSchemaMap defines the mapping of actions to schemas for the UAP SCA service.
var UAPSCAActionToSchemaMap = map[string]interface{}{
	"add-policy":       &sca.ArkUAPSCACloudConsoleAccessPolicy{},
	"delete-policy":    &common.ArkUAPDeletePolicyRequest{},
	"update-policy":    &sca.ArkUAPSCACloudConsoleAccessPolicy{},
	"policy":           &common.ArkUAPGetPolicyRequest{},
	"list-policies":    nil,
	"list-policies-by": &sca.ArkUAPSCAFilters{},
	"policies-stats":   nil,
	"policy-status":    &common.ArkUAPGetPolicyStatus{},
}

// UAPSCAActions defines the actions for the UAP SCA service.
var UAPSCAActions = actions.ArkServiceActionDefinition{
	ActionName: "sca",
	Schemas:    UAPSCAActionToSchemaMap,
}

// UAPSIADBActionToSchemaMap defines the mapping of actions to schemas for the UAP SIA DB service.
var UAPSIADBActionToSchemaMap = map[string]interface{}{
	"add-policy":       &db.ArkUAPSIADBAccessPolicy{},
	"delete-policy":    &common.ArkUAPDeletePolicyRequest{},
	"update-policy":    &db.ArkUAPSIADBAccessPolicy{},
	"policy":           &common.ArkUAPGetPolicyRequest{},
	"list-policies":    nil,
	"list-policies-by": &db.ArkUAPSIADBFilters{},
	"policies-stats":   nil,
	"policy-status":    &common.ArkUAPGetPolicyStatus{},
}

// UAPSIADBAction defines the actions for the UAP SIA DB service.
var UAPSIADBAction = actions.ArkServiceActionDefinition{
	ActionName: "db",
	Schemas:    UAPSIADBActionToSchemaMap,
}

// UAPSIAVMActionToSchemaMap defines the mapping of actions to schemas for the UAP SIA VM service.
var UAPSIAVMActionToSchemaMap = map[string]interface{}{
	"add-policy":       &vm.ArkUAPSIAVMAccessPolicy{},
	"delete-policy":    &common.ArkUAPDeletePolicyRequest{},
	"update-policy":    &vm.ArkUAPSIAVMAccessPolicy{},
	"policy":           &common.ArkUAPGetPolicyRequest{},
	"list-policies":    nil,
	"list-policies-by": &vm.ArkUAPSIAVMFilters{},
	"policies-stats":   nil,
	"policy-status":    &common.ArkUAPGetPolicyStatus{},
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
	"list-policies-by": &common.ArkUAPFilters{},
	"policy-status":    &common.ArkUAPGetPolicyStatus{},
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
