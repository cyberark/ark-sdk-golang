package services

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	cmgrmodels "github.com/cyberark/ark-sdk-golang/pkg/services/cmgr/models"
)

// CmgrActionToSchemaMap is a map that defines the mapping between CMGR action names and their corresponding schema types.
var CmgrActionToSchemaMap = map[string]interface{}{
	"add-network":              &cmgrmodels.ArkCmgrAddNetwork{},
	"update-network":           &cmgrmodels.ArkCmgrUpdateNetwork{},
	"delete-network":           &cmgrmodels.ArkCmgrDeleteNetwork{},
	"list-networks":            nil,
	"list-networks-by":         &cmgrmodels.ArkCmgrNetworksFilter{},
	"network":                  &cmgrmodels.ArkCmgrGetNetwork{},
	"networks-stats":           nil,
	"add-pool":                 &cmgrmodels.ArkCmgrAddPool{},
	"update-pool":              &cmgrmodels.ArkCmgrUpdatePool{},
	"delete-pool":              &cmgrmodels.ArkCmgrDeletePool{},
	"list-pools":               nil,
	"list-pools-by":            &cmgrmodels.ArkCmgrPoolsFilter{},
	"pool":                     &cmgrmodels.ArkCmgrGetPool{},
	"pools-stats":              nil,
	"add-pool-identifier":      &cmgrmodels.ArkCmgrAddPoolSingleIdentifier{},
	"add-pool-identifiers":     &cmgrmodels.ArkCmgrAddPoolBulkIdentifier{},
	"update-pool-identifier":   &cmgrmodels.ArkCmgrUpdatePoolIdentifier{},
	"delete-pool-identifier":   &cmgrmodels.ArkCmgrDeletePoolSingleIdentifier{},
	"delete-pool-identifiers":  &cmgrmodels.ArkCmgrDeletePoolBulkIdentifier{},
	"list-pool-identifiers":    &cmgrmodels.ArkCmgrListPoolIdentifiers{},
	"list-pool-identifiers-by": &cmgrmodels.ArkCmgrPoolIdentifiersFilter{},
	"pool-identifier":          &cmgrmodels.ArkCmgrGetPoolIdentifier{},
	"list-pools-components":    nil,
	"list-pools-components-by": &cmgrmodels.ArkCmgrPoolComponentsFilter{},
	"pool-component":           &cmgrmodels.ArkCmgrGetPoolComponent{},
}

// CmgrActions is a struct that defines the CMGR action for the Ark service.
var CmgrActions = &actions.ArkServiceActionDefinition{
	ActionName: "cmgr",
	Schemas:    CmgrActionToSchemaMap,
}
