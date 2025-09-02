package main

import (
	"fmt"

	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	authmodels "github.com/cyberark/ark-sdk-golang/pkg/models/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/services/sia"
	dbsecretsmodels "github.com/cyberark/ark-sdk-golang/pkg/services/sia/secrets/db/models"
	workspacesdbmodels "github.com/cyberark/ark-sdk-golang/pkg/services/sia/workspaces/db/models"

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

	// Add a DB secret
	siaAPI, err := sia.NewArkSIAAPI(ispAuth.(*auth.ArkISPAuth))
	if err != nil {
		panic(err)
	}
	secret, err := siaAPI.SecretsDB().AddSecret(
		&dbsecretsmodels.ArkSIADBAddSecret{
			SecretType: "username_password",
			Username:   "CoolUser",
			Password:   "CoolPassword",
		},
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("Secret ID:", secret.SecretID)

	// Add the database with the created secret
	database, err := siaAPI.WorkspacesDB().AddDatabase(
		&workspacesdbmodels.ArkSIADBAddDatabase{
			Name:              "MyDatabase",
			ProviderEngine:    workspacesdbmodels.EngineTypeAuroraMysql,
			ReadWriteEndpoint: "myrds.com",
			SecretID:          secret.SecretID,
		},
	)

	if err != nil {
		panic(err)
	}
	fmt.Printf("Database: %v\n", database)
}
