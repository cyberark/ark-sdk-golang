package main

import (
	"fmt"
	"os"

	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	authmodels "github.com/cyberark/ark-sdk-golang/pkg/models/auth"
	commonmodels "github.com/cyberark/ark-sdk-golang/pkg/models/common"
	"github.com/cyberark/ark-sdk-golang/pkg/services/sia"
	dbsecretsmodels "github.com/cyberark/ark-sdk-golang/pkg/services/sia/secrets/db/models"
	dbmodels "github.com/cyberark/ark-sdk-golang/pkg/services/sia/workspaces/db/models"
	"github.com/cyberark/ark-sdk-golang/pkg/services/uap"
	commonuapmodels "github.com/cyberark/ark-sdk-golang/pkg/services/uap/common/models"
	uapsia "github.com/cyberark/ark-sdk-golang/pkg/services/uap/sia/common/models"
	uapdbmodels "github.com/cyberark/ark-sdk-golang/pkg/services/uap/sia/db/models"
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

	uapAPI, err := uap.NewArkUAPAPI(ispAuth.(*auth.ArkISPAuth))
	if err != nil {
		panic(err)
	}
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
		&dbmodels.ArkSIADBAddDatabase{
			Name:              "MyDatabase",
			ProviderEngine:    dbmodels.EngineTypeAuroraMysql,
			ReadWriteEndpoint: "myrds.com",
			SecretID:          secret.SecretID,
		},
	)

	if err != nil {
		panic(err)
	}
	fmt.Printf("Database: %v\n", database)

	policy, err := uapAPI.Db().AddPolicy(
		&uapdbmodels.ArkUAPSIADBAccessPolicy{
			ArkUAPSIACommonAccessPolicy: uapsia.ArkUAPSIACommonAccessPolicy{
				ArkUAPCommonAccessPolicy: commonuapmodels.ArkUAPCommonAccessPolicy{
					Metadata: commonuapmodels.ArkUAPMetadata{
						Name:        "Example DB Access Policy",
						Description: "This is an example of a DB access policy for SIA.",
						Status: commonuapmodels.ArkUAPPolicyStatus{
							Status: commonuapmodels.StatusTypeActive,
						},
						PolicyEntitlement: commonuapmodels.ArkUAPPolicyEntitlement{
							TargetCategory: commonmodels.CategoryTypeDB,
							LocationType:   commonmodels.WorkspaceTypeFQDNIP,
							PolicyType:     commonuapmodels.PolicyTypeRecurring,
						},
						PolicyTags: []string{},
					},
					Principals: []commonuapmodels.ArkUAPPrincipal{
						{
							Type:                commonuapmodels.PrincipalTypeUser,
							ID:                  "user-id",
							Name:                "user@cyberark.cloud.12345",
							SourceDirectoryName: "CyberArk",
							SourceDirectoryID:   "12345",
						},
					},
				},
				Conditions: uapsia.ArkUAPSIACommonConditions{
					ArkUAPConditions: commonuapmodels.ArkUAPConditions{
						AccessWindow: commonuapmodels.ArkUAPTimeCondition{
							DaysOfTheWeek: []int{1, 2, 3, 4, 5},
							FromHour:      "09:00",
							ToHour:        "17:00",
						},
						MaxSessionDuration: 4,
					},
					IdleTime: 10,
				},
			},
			Targets: map[string]uapdbmodels.ArkUAPSIADBTargets{
				commonmodels.WorkspaceTypeFQDNIP: {
					Instances: []uapdbmodels.ArkUAPSIADBInstanceTarget{
						{
							InstanceName:         database.Name,
							InstanceType:         database.ProviderDetails.Family,
							InstanceID:           string(rune(database.ID)),
							AuthenticationMethod: uapdbmodels.AuthMethodDBAuth,
							DBAuthProfile: &uapdbmodels.ArkUAPSIADBDBAuthProfile{
								Roles: []string{"db_reader", "db_writer"},
							},
						},
					},
				},
			},
		},
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Policy created successfully: %s\n", policy.Metadata.PolicyID)
}
