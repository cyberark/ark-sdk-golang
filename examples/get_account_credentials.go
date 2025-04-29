package main

import (
	"fmt"
	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	authmodels "github.com/cyberark/ark-sdk-golang/pkg/models/auth"
	accountsmodels "github.com/cyberark/ark-sdk-golang/pkg/models/services/pcloud/accounts"
	"github.com/cyberark/ark-sdk-golang/pkg/services/pcloud"
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

	// Retrieve a specific account credentials
	pcloudAPI, err := pcloud.NewArkPCloudAPI(ispAuth.(*auth.ArkISPAuth))
	if err != nil {
		panic(err)
	}
	creds, err := pcloudAPI.Accounts().AccountCredentials(&accountsmodels.ArkPCloudGetAccountCredentials{
		AccountID: "11_1",
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Account credentials: %s\n", creds.Password)
}
