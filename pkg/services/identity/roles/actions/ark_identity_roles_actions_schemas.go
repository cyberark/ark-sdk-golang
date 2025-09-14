package actions

import (
	rolesmodels "github.com/cyberark/ark-sdk-golang/pkg/services/identity/roles/models"
)

// ActionToSchemaMap is a map that defines the mapping between Roles action names and their corresponding schema types.
var ActionToSchemaMap = map[string]interface{}{
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
