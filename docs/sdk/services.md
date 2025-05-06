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

The above example authenticates to the specified ISP tenant, initializes a CMGR service using the authorized authenticator, and then uses the service to add a network and pool.

## Secure Infrastructure Access service

The Secure Infrastructure Access (sia) service requires the ArkISPAuth authenticator, and exposes these service classes:

- **ArkSIAAccessService** (access) - SIA access service
- **ArkSIAK8SService** (Kubernetes) - SIA end-user Kubernetes service
- **ArkSIASecretsService** (secrets) - SIA secrets management
    - **ArkSIAVMSecretsService** (VM) - SIA VM secrets services
- **ArkSIASSOService** (SSP) - SIA end-user SSO service
- **ArkSIADatabasesService** (databases) - SIA end-user databases service
- **ArkSIAWorkspacesService** (workspaces) - SIA workspaces management
    - **ArkSIATargetSetsWorkspaceService** (db) - SIA Target Sets workspace management


## Identity service
The Identity (identity) service requires the ArkISPAuth authenticator, and exposes those service classes:

- **ArkIdentityRolesService** - Identity roles service
- **ArkIdentityUsersService** - Identity users service
- **ArkIdentityDirectoriesService** - Identity directories service


## Privilege Cloud service
The Privilege Cloud (pCloud) service requires the ArkISPAuth authenticator, and exposes those service classes:

- **ArkPCloudAccountsService** - Accounts management service
- **ArkPCloudSafesService** - Safes management service


## Connector Manager Service
The Connector Manager (cmgr) service requires the ArkISPAuth authenticator, and exposes those service classes:

- **ArkCmgrService** - Connector Manager service
