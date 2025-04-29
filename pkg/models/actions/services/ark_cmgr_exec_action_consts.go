package services

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	"github.com/cyberark/ark-sdk-golang/pkg/models/services/cmgr"
)

// CmgrActionToSchemaMap is a map that defines the mapping between CMGR action names and their corresponding schema types.
var CmgrActionToSchemaMap = map[string]interface{}{
	"add-network":              &cmgr.ArkCmgrAddNetwork{},
	"update-network":           &cmgr.ArkCmgrUpdateNetwork{},
	"delete-network":           &cmgr.ArkCmgrDeleteNetwork{},
	"list-networks":            nil,
	"list-networks-by":         &cmgr.ArkCmgrNetworksFilter{},
	"network":                  &cmgr.ArkCmgrGetNetwork{},
	"networks-stats":           nil,
	"add-pool":                 &cmgr.ArkCmgrAddPool{},
	"update-pool":              &cmgr.ArkCmgrUpdatePool{},
	"delete-pool":              &cmgr.ArkCmgrDeletePool{},
	"list-pools":               nil,
	"list-pools-by":            &cmgr.ArkCmgrPoolsFilter{},
	"pool":                     &cmgr.ArkCmgrGetPool{},
	"pools-stats":              nil,
	"add-pool-identifier":      &cmgr.ArkCmgrAddPoolSingleIdentifier{},
	"add-pool-identifiers":     &cmgr.ArkCmgrAddPoolBulkIdentifier{},
	"delete-pool-identifier":   &cmgr.ArkCmgrDeletePoolSingleIdentifier{},
	"delete-pool-identifiers":  &cmgr.ArkCmgrDeletePoolBulkIdentifier{},
	"list-pool-identifiers":    &cmgr.ArkCmgrListPoolIdentifiers{},
	"list-pool-identifiers-by": &cmgr.ArkCmgrPoolIdentifiersFilter{},
	"list-pools-components":    nil,
	"list-pools-components-by": &cmgr.ArkCmgrPoolComponentsFilter{},
	"pool-component":           &cmgr.ArkCmgrGetPoolComponent{},
}

// CmgrActions is a struct that defines the CMGR action for the Ark service.
var CmgrActions = &actions.ArkServiceActionDefinition{
	ActionName: "cmgr",
	Schemas:    CmgrActionToSchemaMap,
}
