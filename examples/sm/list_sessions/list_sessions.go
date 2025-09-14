package main

import (
	"fmt"
	"os"

	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	authmodels "github.com/cyberark/ark-sdk-golang/pkg/models/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/services/sm"
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

	// List all sessions
	smService, err := sm.NewArkSMService(ispAuth.(*auth.ArkISPAuth))
	if err != nil {
		panic(err)
	}
	smChan, err := smService.ListSessions()
	if err != nil {
		panic(err)
	}
	for loadedSessions := range smChan {
		fmt.Printf("sessions: %v\n", loadedSessions)
	}
}
