package main

import (
	"fmt"
	"os"

	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	authmodels "github.com/cyberark/ark-sdk-golang/pkg/models/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/services/cmgr"
	cmgrmodels "github.com/cyberark/ark-sdk-golang/pkg/services/cmgr/models"
	"github.com/cyberark/ark-sdk-golang/pkg/services/sia"
	accessmodels "github.com/cyberark/ark-sdk-golang/pkg/services/sia/access/models"
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

	// Install a connector on the pool above
	siaAPI, err := sia.NewArkSIAAPI(ispAuth.(*auth.ArkISPAuth))
	if err != nil {
		panic(err)
	}
	connectorID, err := siaAPI.Access().InstallConnector(
		&accessmodels.ArkSIAInstallConnector{
			ConnectorType:   "ON-PREMISE",
			ConnectorOS:     "linux",
			ConnectorPoolID: pool.PoolID,
			TargetMachine:   "1.1.1.1",
			Username:        "root",
			PrivateKeyPath:  "/path/to/key.pem",
		},
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Installed connector: %s\n", connectorID)
}
