package actions

import workspacesdbmodels "github.com/cyberark/ark-sdk-golang/pkg/services/sia/workspaces/db/models"

// ActionToSchemaMap is a map that defines the mapping between DB workspace action names and their corresponding schema types.
var ActionToSchemaMap = map[string]interface{}{
	"add-database":      &workspacesdbmodels.ArkSIADBAddDatabase{},
	"delete-database":   &workspacesdbmodels.ArkSIADBDeleteDatabase{},
	"update-database":   &workspacesdbmodels.ArkSIADBUpdateDatabase{},
	"database":          &workspacesdbmodels.ArkSIADBGetDatabase{},
	"list-databases":    nil,
	"list-databases-by": &workspacesdbmodels.ArkSIADBDatabasesFilter{},
	"databases-stats":   nil,
	"list-engine-types": nil,
	"list-family-types": nil,
}
