package actions

import usersmodels "github.com/cyberark/ark-sdk-golang/pkg/services/identity/users/models"

// ActionToSchemaMapIdentityUsers is a map that defines the mapping between Users action names and their corresponding schema types.
var ActionToSchemaMapIdentityUsers = map[string]interface{}{
	"create-user":         &usersmodels.ArkIdentityCreateUser{},
	"update-user":         &usersmodels.ArkIdentityUpdateUser{},
	"delete-user":         &usersmodels.ArkIdentityDeleteUser{},
	"user-by-name":        &usersmodels.ArkIdentityUserByName{},
	"user-id-by-name":     &usersmodels.ArkIdentityUserIDByName{},
	"reset-user-password": &usersmodels.ArkIdentityResetUserPassword{},
}
