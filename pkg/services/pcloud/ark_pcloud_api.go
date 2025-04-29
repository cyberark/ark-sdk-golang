package pcloud

import (
	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/services/pcloud/accounts"
	"github.com/cyberark/ark-sdk-golang/pkg/services/pcloud/safes"
)

// ArkPCloudAPI is a struct that provides access to the Ark PCloud API as a wrapped set of services.
type ArkPCloudAPI struct {
	safesService    *safes.ArkPCloudSafesService
	accountsService *accounts.ArkPCloudAccountsService
}

// NewArkPCloudAPI creates a new instance of ArkPCloudAPI with the provided ArkISPAuth.
func NewArkPCloudAPI(ispAuth *auth.ArkISPAuth) (*ArkPCloudAPI, error) {
	var baseIspAuth auth.ArkAuth = ispAuth
	safesService, err := safes.NewArkPCloudSafesService(baseIspAuth)
	if err != nil {
		return nil, err
	}
	accountsService, err := accounts.NewArkPCloudAccountsService(baseIspAuth)
	if err != nil {
		return nil, err
	}
	return &ArkPCloudAPI{
		safesService:    safesService,
		accountsService: accountsService,
	}, nil
}

// Safes returns the Safes service of the ArkPCloudAPI instance.
func (api *ArkPCloudAPI) Safes() *safes.ArkPCloudSafesService {
	return api.safesService
}

// Accounts returns the Accounts service of the ArkPCloudAPI instance.
func (api *ArkPCloudAPI) Accounts() *accounts.ArkPCloudAccountsService {
	return api.accountsService
}
