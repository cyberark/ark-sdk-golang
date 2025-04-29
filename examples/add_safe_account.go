package main

import (
	"fmt"
	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	authmodels "github.com/cyberark/ark-sdk-golang/pkg/models/auth"
	accountsmodels "github.com/cyberark/ark-sdk-golang/pkg/models/services/pcloud/accounts"
	safesmodels "github.com/cyberark/ark-sdk-golang/pkg/models/services/pcloud/safes"
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

	// Add a new safe and account
	pcloudAPI, err := pcloud.NewArkPCloudAPI(ispAuth.(*auth.ArkISPAuth))
	if err != nil {
		panic(err)
	}
	safe, err := pcloudAPI.Safes().AddSafe(&safesmodels.ArkPCloudAddSafe{
		SafeName: "mysafe",
	})
	if err != nil {
		panic(err)
	}
	account, err := pcloudAPI.Accounts().AddAccount(&accountsmodels.ArkPCloudAddAccount{
		SafeName:   safe.SafeName,
		Secret:     "mysecret",
		UserName:   "myuser",
		Address:    "myaddr.com",
		PlatformID: "UnixSSH",
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Safe added: %s\n", safe.SafeName)
	fmt.Printf("Account added: %s\n", account.ID)
}
