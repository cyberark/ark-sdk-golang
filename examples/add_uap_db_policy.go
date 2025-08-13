package main

import (
	"fmt"
	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	authmodels "github.com/cyberark/ark-sdk-golang/pkg/models/auth"
	commonmodels "github.com/cyberark/ark-sdk-golang/pkg/models/common"
	"github.com/cyberark/ark-sdk-golang/pkg/models/services/sia/workspaces/db"
	commonuapmodels "github.com/cyberark/ark-sdk-golang/pkg/models/services/uap/common"
	uapsia "github.com/cyberark/ark-sdk-golang/pkg/models/services/uap/sia/common"
	uapdbmodels "github.com/cyberark/ark-sdk-golang/pkg/models/services/uap/sia/db"
	"github.com/cyberark/ark-sdk-golang/pkg/services/uap"
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

	uapAPI, err := uap.NewArkUAPAPI(ispAuth.(*auth.ArkISPAuth))
	if err != nil {
		panic(err)
	}
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
							InstanceName:         "example-db-instance",
							InstanceType:         db.FamilyTypeMSSQL,
							InstanceID:           "1",
							AuthenticationMethod: uapdbmodels.AuthMethodLDAPAuth,
							LDAPAuthProfile: &uapdbmodels.ArkUAPSIADBLDAPAuthProfile{
								AssignGroups: []string{"mygroup"},
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
