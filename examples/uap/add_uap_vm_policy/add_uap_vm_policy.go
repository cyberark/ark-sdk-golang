package main

import (
	"fmt"
	"os"

	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	authmodels "github.com/cyberark/ark-sdk-golang/pkg/models/auth"
	commonmodels "github.com/cyberark/ark-sdk-golang/pkg/models/common"
	"github.com/cyberark/ark-sdk-golang/pkg/services/uap"
	commonuapmodels "github.com/cyberark/ark-sdk-golang/pkg/services/uap/common/models"
	uapsia "github.com/cyberark/ark-sdk-golang/pkg/services/uap/sia/common/models"
	uapvmmodels "github.com/cyberark/ark-sdk-golang/pkg/services/uap/sia/vm/models"
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
	policy, err := uapAPI.VM().AddPolicy(
		&uapvmmodels.ArkUAPSIAVMAccessPolicy{
			ArkUAPSIACommonAccessPolicy: uapsia.ArkUAPSIACommonAccessPolicy{
				ArkUAPCommonAccessPolicy: commonuapmodels.ArkUAPCommonAccessPolicy{
					Metadata: commonuapmodels.ArkUAPMetadata{
						Name:        "Example VM Access Policy",
						Description: "This is an example of a VM access policy for SIA.",
						Status: commonuapmodels.ArkUAPPolicyStatus{
							Status: commonuapmodels.StatusTypeActive,
						},
						PolicyEntitlement: commonuapmodels.ArkUAPPolicyEntitlement{
							TargetCategory: commonmodels.CategoryTypeVM,
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
			Targets: uapvmmodels.ArkUAPSIAVMPlatformTargets{
				FQDNIPResource: &uapvmmodels.ArkUAPSIAVMFQDNIPResource{
					FQDNRules: []uapvmmodels.ArkUAPSIAVMFQDNRule{
						{
							Operator:            uapvmmodels.VMFQDNOperatorExactly,
							ComputernamePattern: "example-vm",
							Domain:              "mydomain.com",
						},
					},
				},
			},
			Behavior: uapvmmodels.ArkUAPSSIAVMBehavior{
				SSHProfile: &uapvmmodels.ArkUAPSSIAVMSSHProfile{
					Username: "root",
				},
				RDPProfile: &uapvmmodels.ArkUAPSSIAVMRDPProfile{
					LocalEphemeralUser: &uapvmmodels.ArkUAPSSIAVMEphemeralUser{
						AssignGroups: []string{"Remote Desktop Users"},
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
