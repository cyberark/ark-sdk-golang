---
title: SDK Examples
description: SDK Examples
---

# SDK Examples
Using the SDK is similar to using the CLI.

## Short-lived password example

In this example we authenticate to our ISP tenant and create a short-lived password:

```go
package main

import (
	"fmt"
	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	authmodels "github.com/cyberark/ark-sdk-golang/pkg/models/auth"
	ssomodels "github.com/cyberark/ark-sdk-golang/pkg/models/services/sia/sso"
	"github.com/cyberark/ark-sdk-golang/pkg/services/sia/sso"
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

	// Create an SSO service from the authenticator above
	ssoService, err := sso.NewArkSIASSOService(ispAuth)
	if err != nil {
		panic(err)
	}

	// Generate a short-lived password
	ssoPassword, err := ssoService.ShortLivedPassword(
		&ssomodels.ArkSIASSOGetShortLivedPassword{},
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", ssoPassword)

	// Generate a short-lived password for RDP
	ssoPassword, err = ssoService.ShortLivedPassword(
		&ssomodels.ArkSIASSOGetShortLivedPassword{
			Service: "DPA-RDP",
		},
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", ssoPassword)
}
```

## Target set example

In this example we authenticate to our ISP tenant and create a target set with a VM secret:

```go
package main

import (
	"fmt"
	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	authmodels "github.com/cyberark/ark-sdk-golang/pkg/models/auth"
	vmsecretsmodels "github.com/cyberark/ark-sdk-golang/pkg/models/services/sia/secrets/vm"
	targetsetsmodels "github.com/cyberark/ark-sdk-golang/pkg/models/services/sia/workspaces/targetsets"
	"github.com/cyberark/ark-sdk-golang/pkg/services/sia"
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

	// Add a VM secret
	siaAPI, err := sia.NewArkSIAAPI(ispAuth.(*auth.ArkISPAuth))
	if err != nil {
		panic(err)
	}
	secret, err := siaAPI.SecretsVM().AddSecret(
		&vmsecretsmodels.ArkSIAVMAddSecret{
			SecretType:          "ProvisionerUser",
			ProvisionerUsername: "CoolUser",
			ProvisionerPassword: "CoolPassword",
		},
	)
	if err != nil {
		panic(err)
	}
	// Add VM target set
	targetSet, err := siaAPI.WorkspacesTargetSets().AddTargetSet(
		&targetsetsmodels.ArkSIAAddTargetSet{
			Name:       "mydomain.com",
			Type:       "Domain",
			SecretID:   secret.SecretID,
			SecretType: secret.SecretType,
		},
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Target set %s created\n", targetSet.Name)
}
```

## CMGR example

In this example we authenticate to our ISP tenant and create a network, pool, and identifier:

```go
package main

import (
	"fmt"
	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	authmodels "github.com/cyberark/ark-sdk-golang/pkg/models/auth"
	cmgrmodels "github.com/cyberark/ark-sdk-golang/pkg/models/services/cmgr"
	"github.com/cyberark/ark-sdk-golang/pkg/services/cmgr"
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

	// Configure a network, pool and identifiers
	cmgrService, err := cmgr.NewArkCmgrService(ispAuth.(*auth.ArkISPAuth))
	if err != nil {
		panic(err)
	}
	network, err := cmgrService.AddNetwork(&cmgrmodels.ArkCmgrAddNetwork{Name: "tlv"})
	if err != nil {
		panic(err)
	}
	pool, err := cmgrService.AddPool(&cmgrmodels.ArkCmgrAddPool{Name: "tlvpool", AssignedNetworkIDs: []string{network.NetworkID}})
	if err != nil {
		panic(err)
	}
	identifier, err := cmgrService.AddPoolIdentifier(&cmgrmodels.ArkCmgrAddPoolSingleIdentifier{PoolID: pool.PoolID, Type: cmgrmodels.GeneralFQDN, Value: "mymachine.tlv.com"})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Added pool: %s\n", pool.PoolID)
	fmt.Printf("Added network: %s\n", network.NetworkID)
	fmt.Printf("Added identifier: %s\n", identifier.IdentifierID)
}
```

## List pCloud Accounts

In this example we authenticate to our ISP tenant and list pCloud accounts:

```go
package main

import (
	"fmt"
	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	authmodels "github.com/cyberark/ark-sdk-golang/pkg/models/auth"
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

	// List all of the accounts
	pcloudAPI, err := pcloud.NewArkPCloudAPI(ispAuth.(*auth.ArkISPAuth))
	if err != nil {
		panic(err)
	}
	accountsChan, err := pcloudAPI.Accounts().ListAccounts()
	if err != nil {
		panic(err)
	}
	for accountsPage := range accountsChan {
		for account := range accountsPage.Items {
			fmt.Printf("Account: %v\n", account)
		}
	}
}
```

## List identities

In this example we authenticate to our ISP tenant and list all of the accounts:

```go
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
```

## Session Monitoring

In this example we authenticate to our ISP tenant and get all the active sessions:
```go
package main

import (
	"fmt"
	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	authmodels "github.com/cyberark/ark-sdk-golang/pkg/models/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/services/sm"
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

	SMAPI, err := sm.NewArkSMService(ispAuth.(*auth.ArkISPAuth))
	if err != nil {
		panic(err)
	}
	filter := &ArkSMSessionsFilter{
		Search: "status IN Active",
	}
	// Get all active sessions
	activeSessions, err := SMAPI.CountSessionsBy(filter)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Total Active Sessions: %d\n", activeSessions)
}
```

## UAP

In this example we authenticate to our ISP tenant and create a UAP DB policy:
```go
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
```

In this example we authenticate to our ISP tenant and create a UAP SIA VM policy:
```go
package main

import (
	"fmt"
	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	authmodels "github.com/cyberark/ark-sdk-golang/pkg/models/auth"
	commonmodels "github.com/cyberark/ark-sdk-golang/pkg/models/common"
	commonuapmodels "github.com/cyberark/ark-sdk-golang/pkg/models/services/uap/common"
	uapsia "github.com/cyberark/ark-sdk-golang/pkg/models/services/uap/sia/common"
	uapvmmodels "github.com/cyberark/ark-sdk-golang/pkg/models/services/uap/sia/vm"
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
	policy, err := uapAPI.Vm().AddPolicy(
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
```

In this example we authenticate to our ISP tenant and create a UAP SCA policy:
```go
package main

import (
	"fmt"
	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	authmodels "github.com/cyberark/ark-sdk-golang/pkg/models/auth"
	commonmodels "github.com/cyberark/ark-sdk-golang/pkg/models/common"
	commonuapmodels "github.com/cyberark/ark-sdk-golang/pkg/models/services/uap/common"
	uapscamodels "github.com/cyberark/ark-sdk-golang/pkg/models/services/uap/sca"
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
```
