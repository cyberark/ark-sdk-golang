package main

import (
	"fmt"
	"os"

	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	authmodels "github.com/cyberark/ark-sdk-golang/pkg/models/auth"
	commonmodels "github.com/cyberark/ark-sdk-golang/pkg/models/common"
	"github.com/cyberark/ark-sdk-golang/pkg/services/uap"
	commonuapmodels "github.com/cyberark/ark-sdk-golang/pkg/services/uap/common/models"
	uapscamodels "github.com/cyberark/ark-sdk-golang/pkg/services/uap/sca/models"
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
	policy, err := uapAPI.Sca().AddPolicy(
		&uapscamodels.ArkUAPSCACloudConsoleAccessPolicy{
			ArkUAPCommonAccessPolicy: commonuapmodels.ArkUAPCommonAccessPolicy{
				Metadata: commonuapmodels.ArkUAPMetadata{
					Name:        "Example SCA Access Policy",
					Description: "This is an example of a SCA access policy.",
					Status: commonuapmodels.ArkUAPPolicyStatus{
						Status: commonuapmodels.StatusTypeValidating,
					},
					PolicyEntitlement: commonuapmodels.ArkUAPPolicyEntitlement{
						TargetCategory: commonmodels.CategoryTypeCloudConsole,
						LocationType:   commonmodels.WorkspaceTypeAWS,
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
			Conditions: uapscamodels.ArkUAPSCAConditions{
				ArkUAPConditions: commonuapmodels.ArkUAPConditions{
					AccessWindow: commonuapmodels.ArkUAPTimeCondition{
						DaysOfTheWeek: []int{1, 2, 3, 4, 5},
						FromHour:      "09:00:00",
						ToHour:        "17:00:00",
					},
					MaxSessionDuration: 4,
				},
			},
			Targets: uapscamodels.ArkUAPSCACloudConsoleTarget{
				AwsAccountTargets: []uapscamodels.ArkUAPSCAAWSAccountTarget{
					{
						uapscamodels.ArkUAPSCATarget{
							RoleID:        "arn:aws:iam::123456789012:role/ExampleRole",
							RoleName:      "ExampleRole",
							WorkspaceID:   "123456789012",
							WorkspaceName: "ExampleWorkspace",
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
