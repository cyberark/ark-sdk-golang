package services

import (
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	directoriesmodels "github.com/cyberark/ark-sdk-golang/pkg/services/identity/directories/models"
	rolesmodels "github.com/cyberark/ark-sdk-golang/pkg/services/identity/roles/models"
	usersmodels "github.com/cyberark/ark-sdk-golang/pkg/services/identity/users/models"
)

// DirectoriesActionToSchemaMap is a map that defines the mapping between Directories action names and their corresponding schema types.
var DirectoriesActionToSchemaMap = map[string]interface{}{
	"list-directories":          &directoriesmodels.ArkIdentityListDirectories{},
	"list-directories-entities": &directoriesmodels.ArkIdentityListDirectoriesEntities{},
	"tenant-default-suffix":     nil,
}

// DirectoriesAction is a struct that defines the Directories action for the Ark service.
var DirectoriesAction = actions.ArkServiceActionDefinition{
	ActionName: "directories",
	Schemas:    DirectoriesActionToSchemaMap,
}

// RolesActionToSchemaMap is a map that defines the mapping between Roles action names and their corresponding schema types.
var RolesActionToSchemaMap = map[string]interface{}{
	"add-user-to-role":         &rolesmodels.ArkIdentityAddUserToRole{},
	"add-group-to-role":        &rolesmodels.ArkIdentityAddGroupToRole{},
	"add-role-to-role":         &rolesmodels.ArkIdentityAddRoleToRole{},
	"remove-user-from-role":    &rolesmodels.ArkIdentityRemoveUserFromRole{},
	"remove-group-from-role":   &rolesmodels.ArkIdentityRemoveGroupFromRole{},
	"remove-role-from-role":    &rolesmodels.ArkIdentityRemoveRoleFromRole{},
	"create-role":              &rolesmodels.ArkIdentityCreateRole{},
	"update-role":              &rolesmodels.ArkIdentityUpdateRole{},
	"delete-role":              &rolesmodels.ArkIdentityDeleteRole{},
	"list-role-members":        &rolesmodels.ArkIdentityListRoleMembers{},
	"add-admin-rights-to-role": &rolesmodels.ArkIdentityAddAdminRightsToRole{},
	"role-id-by-name":          &rolesmodels.ArkIdentityRoleIDByName{},
}

// RolesAction is a struct that defines the Roles action for the Ark service.
var RolesAction = actions.ArkServiceActionDefinition{
	ActionName: "roles",
	Schemas:    RolesActionToSchemaMap,
}

// UsersActionToSchemaMap is a map that defines the mapping between Users action names and their corresponding schema types.
var UsersActionToSchemaMap = map[string]interface{}{
	"create-user":         &usersmodels.ArkIdentityCreateUser{},
	"update-user":         &usersmodels.ArkIdentityUpdateUser{},
	"delete-user":         &usersmodels.ArkIdentityDeleteUser{},
	"user-by-name":        &usersmodels.ArkIdentityUserByName{},
	"user-id-by-name":     &usersmodels.ArkIdentityUserIDByName{},
	"reset-user-password": &usersmodels.ArkIdentityResetUserPassword{},
}

// UsersAction is a struct that defines the Users action for the Ark service.
var UsersAction = actions.ArkServiceActionDefinition{
	ActionName: "users",
	Schemas:    UsersActionToSchemaMap,
}

// IdentityActions is a struct that defines the Identity actions for the Ark service.
var IdentityActions = &actions.ArkServiceActionDefinition{
	ActionName: "identity",
	Subactions: []*actions.ArkServiceActionDefinition{
		&DirectoriesAction,
		&RolesAction,
		&UsersAction,
	},
}
