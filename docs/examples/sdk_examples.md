---
title: SDK Examples
description: SDK Examples
---

# SDK Examples
Using the SDK is similar to using the CLI

## Short lived password example

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
}
```

In the above example we authenticate to our ISP tenant and create a short lived password

## Target set example

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
	secret, err := siaAPI.VMSecrets().AddSecret(
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
	targetSet, err := siaAPI.TargetSets().AddTargetSet(
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

In the above example we authenticate to our ISP tenant and create a target set with a VM secret

## CMGR example

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
	pool, err := cmgrService.AddPool(&cmgrmodels.ArkCmgrAddPool{Name: "tlvpool", AssignedNetworkIDs: []string{network.ID}})
	if err != nil {
		panic(err)
	}
	identifier, err := cmgrService.AddPoolIdentifier(&cmgrmodels.ArkCmgrAddPoolSingleIdentifier{PoolID: pool.ID, Type: cmgrmodels.GeneralFQDN, Value: "mymachine.tlv.com"})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Added pool: %s\n", pool.ID)
	fmt.Printf("Added network: %s\n", network.ID)
	fmt.Printf("Added identifier: %s\n", identifier.ID)
}
```

In the above example we authenticate to our ISP tenant and create a network, pool and identifier

## List PCloud Accounts

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

## List Identities

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

In the above example we authenticate to our ISP tenant and list all of the accounts
