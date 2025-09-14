package actions

import directoriesmodels "github.com/cyberark/ark-sdk-golang/pkg/services/identity/directories/models"

// ActionToSchemaMap is a map that defines the mapping between Directories action names and their corresponding schema types.
var ActionToSchemaMap = map[string]interface{}{
	"list-directories":          &directoriesmodels.ArkIdentityListDirectories{},
	"list-directories-entities": &directoriesmodels.ArkIdentityListDirectoriesEntities{},
	"tenant-default-suffix":     nil,
}
