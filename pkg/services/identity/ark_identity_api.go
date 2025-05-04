package identity

import (
	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/services/identity/directories"
	"github.com/cyberark/ark-sdk-golang/pkg/services/identity/roles"
	"github.com/cyberark/ark-sdk-golang/pkg/services/identity/users"
)

// ArkIdentityAPI is a struct that provides access to the Ark Identity API as a wrapped set of services.
type ArkIdentityAPI struct {
	directoriesService *directories.ArkIdentityDirectoriesService
	rolesService       *roles.ArkIdentityRolesService
	usersService       *users.ArkIdentityUsersService
}

// NewArkIdentityAPI creates a new instance of ArkIdentityAPI with the provided ArkISPAuth.
func NewArkIdentityAPI(ispAuth *auth.ArkISPAuth) (*ArkIdentityAPI, error) {
	var baseIspAuth auth.ArkAuth = ispAuth
	directoriesService, err := directories.NewArkIdentityDirectoriesService(baseIspAuth)
	if err != nil {
		return nil, err
	}
	rolesService, err := roles.NewArkIdentityRolesService(baseIspAuth)
	if err != nil {
		return nil, err
	}
	usersService, err := users.NewArkIdentityUsersService(baseIspAuth)
	if err != nil {
		return nil, err
	}
	return &ArkIdentityAPI{
		directoriesService: directoriesService,
		rolesService:       rolesService,
		usersService:       usersService,
	}, nil
}

// Directories returns the Directories service of the ArkIdentityAPI instance.
func (api *ArkIdentityAPI) Directories() *directories.ArkIdentityDirectoriesService {
	return api.directoriesService
}

// Roles returns the Roles service of the ArkIdentityAPI instance.
func (api *ArkIdentityAPI) Roles() *roles.ArkIdentityRolesService {
	return api.rolesService
}

// Users returns the Users service of the ArkIdentityAPI instance.
func (api *ArkIdentityAPI) Users() *users.ArkIdentityUsersService {
	return api.usersService
}
