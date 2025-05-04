package main

import (
	"fmt"
	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	authmodels "github.com/cyberark/ark-sdk-golang/pkg/models/auth"
	directoriesmodels "github.com/cyberark/ark-sdk-golang/pkg/models/services/identity/directories"
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

	// List all identities
	identityAPI, err := identity.NewArkIdentityAPI(ispAuth.(*auth.ArkISPAuth))
	if err != nil {
		panic(err)
	}
	identitiesChan, err := identityAPI.Directories().ListDirectoriesEntities(&directoriesmodels.ArkIdentityListDirectoriesEntities{})
	if err != nil {
		panic(err)
	}
	for loadedIdentity := range identitiesChan {
		fmt.Printf("Identity: %v\n", loadedIdentity)
	}
}
