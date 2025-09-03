package main

import (
	"fmt"

	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	authmodels "github.com/cyberark/ark-sdk-golang/pkg/models/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/services/sia"
	accessmodels "github.com/cyberark/ark-sdk-golang/pkg/services/sia/access/models"

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

	// Install a connector on the pool above
	siaAPI, err := sia.NewArkSIAAPI(ispAuth.(*auth.ArkISPAuth))
	if err != nil {
		panic(err)
	}
	testReachabilityResponse, err := siaAPI.Access().TestConnectorReachability(
		&accessmodels.ArkSIATestConnectorReachability{
			ConnectorID:           "CMSConnector",
			TargetHostname:        "google.com",
			TargetPort:            443,
			CheckBackendEndpoints: true,
		},
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Reachability response: %v\n", testReachabilityResponse)
}
