package main

import (
	"fmt"
	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	authmodels "github.com/cyberark/ark-sdk-golang/pkg/models/auth"
	rolesmodels "github.com/cyberark/ark-sdk-golang/pkg/models/services/identity/roles"
	usersmodels "github.com/cyberark/ark-sdk-golang/pkg/models/services/identity/users"
	"github.com/cyberark/ark-sdk-golang/pkg/services/identity"
	"os"
)

func main() {
	// Perform authentication using ArkISPAuth to the platform
	// First, create an ISP authentication class
	// Afterwards, perform the authentication
	ispAuth := auth.NewArkISPAuth(false)
	_, err := ispAuth.Authenticate(
		nil,
		&authmodels.ArkAuthProfile{
			Username:           "user@cyberark.cloud.12345",
			AuthMethod:         authmodels.Identity,
			AuthMethodSettings: &authmodels.IdentityArkAuthMethodSettings{},
		},
		&authmodels.ArkSecret{
			Secret: os.Getenv("ARK_SECRET"),
		},
		false,
		false,
	)
	if err != nil {
		panic(err)
	}

	// Add role and user
	identityAPI, err := identity.NewArkIdentityAPI(ispAuth.(*auth.ArkISPAuth))
	if err != nil {
		panic(err)
	}
	role, err := identityAPI.Roles().CreateRole(&rolesmodels.ArkIdentityCreateRole{RoleName: "myrole"})
	if err != nil {
		panic(err)
	}
	user, err := identityAPI.Users().CreateUser(&usersmodels.ArkIdentityCreateUser{Username: "myuser", Roles: []string{role.RoleName}})
	if err != nil {
		panic(err)
	}
	fmt.Printf("User: %v\n", user)
	fmt.Printf("Role: %v\n", role)
}
