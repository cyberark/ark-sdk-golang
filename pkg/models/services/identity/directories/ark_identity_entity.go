package directories

import "github.com/cyberark/ark-sdk-golang/pkg/models/common/identity"

// Possible entity types
const (
	Role  = "ROLE"
	User  = "USER"
	Group = "GROUP"
)

// ArkIdentityEntity is an interface that defines the methods for an identity entity.
type ArkIdentityEntity interface {
	GetEntityType() string
}

// ArkIdentityBaseEntity represents the schema for an identity entity.
type ArkIdentityBaseEntity struct {
	ArkIdentityEntity        `json:"-" mapstructure:"-"`
	ID                       string `json:"id" mapstructure:"id" flag:"id" desc:"ID of the entity" required:"true"`
	Name                     string `json:"name" mapstructure:"name" flag:"name" desc:"Name of the entity" required:"true"`
	EntityType               string `json:"entity_type" mapstructure:"entity_type" flag:"entity-type" desc:"Type of the entity" required:"true" choices:"USER,ROLE,GROUP"`
	DirectoryServiceType     string `json:"directory_service_type" mapstructure:"directory_service_type" flag:"directory-service-type" desc:"Directory type of the entity" required:"true" choices:"AdProxy,CDS,FDS"`
	DisplayName              string `json:"display_name,omitempty" mapstructure:"display_name" flag:"display-name" desc:"Display name of the entity"`
	ServiceInstanceLocalized string `json:"service_instance_localized" mapstructure:"service_instance_localized" flag:"service-instance-localized" desc:"Display directory service name" required:"true"`
}

// GetEntityType returns the entity type of the ArkIdentityUserEntity.
func (a *ArkIdentityBaseEntity) GetEntityType() string {
	return a.EntityType
}

// ArkIdentityUserEntity represents the schema for a user entity.
type ArkIdentityUserEntity struct {
	ArkIdentityBaseEntity
	Email       string `json:"email,omitempty" mapstructure:"email" flag:"email" desc:"Email of the user"`
	Description string `json:"description,omitempty" mapstructure:"description" flag:"description" desc:"Description of the user"`
}

// ArkIdentityGroupEntity represents the schema for a group entity.
type ArkIdentityGroupEntity struct {
	ArkIdentityBaseEntity
}

// GetEntityType returns the entity type of the ArkIdentityGroupEntity.
func (a *ArkIdentityGroupEntity) GetEntityType() string {
	return a.EntityType
}

// ArkIdentityRoleEntity represents the schema for a role entity.
type ArkIdentityRoleEntity struct {
	ArkIdentityBaseEntity
	AdminRights []identity.RoleAdminRight `json:"admin_rights,omitempty" mapstructure:"admin_rights" flag:"admin-rights" desc:"Admin rights of the role"`
	IsHidden    bool                      `json:"is_hidden" mapstructure:"is_hidden" flag:"is-hidden" desc:"Whether this role is hidden or not" required:"true"`
	Description string                    `json:"description,omitempty" mapstructure:"description" flag:"description" desc:"Description of the role"`
}

// GetEntityType returns the entity type of the ArkIdentityRoleEntity.
func (a *ArkIdentityRoleEntity) GetEntityType() string {
	return a.EntityType
}
