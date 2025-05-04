---
title: Services
description: Services
---

# Services

SDK services are defined to execute requests on available ISP services (such as SIA). When a service is initialized, a valid authenticator is required to authorize access to the ISP service. To perform service actions, each service exposes a set of classes and methods.

Here's an example that initializes the `ArkCmgrService` service:

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

The above example authenticates to the specified ISP tenant, initializes a CMGR service using the authorized authenticator, and then uses the service to add a network / pool.

## Secure Infrastructure Access service

The Secure Infrastructure Access (sia) service requires the ArkISPAuth authenticator, and exposes these service classes:

- <b>ArkSIAAccessService (access)</b> - SIA access service
- <b>ArkSIAK8SService (kubernetes)</b> - SIA end-user Kubernetes service
- <b>ArkSIASecretsService (secrets)</b> - SIA secrets management
    - <b>ArkSIAVMSecretsService (vm)</b> - SIA VM secrets services
- <b>ArkSIASSOService (sso)</b> - SIA end-user SSO service
- <b>ArkSIADatabasesService (databases)</b> - SIA end-user databases service
- <b>ArkSIAWorkspacesService (workspaces)</b> - SIA workspaces management
    - <b>ArkSIATargetSetsWorkspaceService (db)</b> - SIA Target Sets workspace management


## Identity service
The Identity (identity) service requires ArkISPAuth authenticator, and exposes those service classes:
- <b>ArkIdentityRolesService - Identity roles service
- <b>ArkIdentityUsersService - Identity users service
- <b>ArkIdentityDirectoriesService - Identity directories service


## Privilege Cloud service
The Privilege Cloud (pcloud) service requires ArkISPAuth authenticator, and exposes those service classes:
- <b>ArkPCloudAccountsService</b> - Accounts management service
- <b>ArkPCloudSafesService</b> - Safes management service


## Connector Manager Service
The Connector Manager (cmgr) service requires ArkISPAuth authenticator, and exposes those service classes:
- <b>ArkCmgrService</b> - Connector Manager service
