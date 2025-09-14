package actions

import policiesmodels "github.com/cyberark/ark-sdk-golang/pkg/services/sechub/syncpolicies/models"

// ActionToSchemaMap is a map that defines the mapping between Sec Hub sync policies action names and their corresponding schema types.
var ActionToSchemaMap = map[string]interface{}{
	"create-sync-policy":    &policiesmodels.ArkSechubCreateSyncPolicy{},
	"delete-sync-policy":    &policiesmodels.ArkSecHubDeleteSyncPolicy{},
	"sync-policy":           &policiesmodels.ArkSecHubGetSyncPolicy{},
	"list-sync-policies":    &policiesmodels.ArkSecHubGetSyncPolicies{},
	"list-sync-policies-by": &policiesmodels.ArkSecHubSyncPoliciesFilters{},
	"set-sync-policy-state": &policiesmodels.ArkSecHubSetSyncPolicyState{},
	"sync-policies-stats":   nil,
}
