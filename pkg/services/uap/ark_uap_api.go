package uap

import (
	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/services/uap/sca"
	"github.com/cyberark/ark-sdk-golang/pkg/services/uap/sia/db"
)

// ArkUAPAPI provides a unified API for accessing various UAP services, including SCA and SIA DB services.
type ArkUAPAPI struct {
	uap *ArkUAPService
	sca *sca.ArkUAPSCAService
	db  *db.ArkUAPSIADBService
	vm  *vm.ArkUAPSIAVMService
}

// NewArkUAPAPI creates a new instance of ArkUAPAPI with the provided ArkISPAuth.
func NewArkUAPAPI(ispAuth *auth.ArkISPAuth) (*ArkUAPAPI, error) {
	var baseIspAuth auth.ArkAuth = ispAuth
	uapService, err := NewArkUAPService(baseIspAuth)
	if err != nil {
		return nil, err
	}
	scaService, err := sca.NewArkUAPSCAService(baseIspAuth)
	if err != nil {
		return nil, err
	}
	dbService, err := db.NewArkUAPSIADBService(baseIspAuth)
	if err != nil {
		return nil, err
	}
	vmService, err := vm.NewArkUAPSIAVMService(baseIspAuth)
	if err != nil {
		return nil, err
	}
	return &ArkUAPAPI{
		uap: uapService,
		sca: scaService,
		db:  dbService,
		vm:  vmService,
	}, nil
}

// Uap returns the ArkUAPService instance from the ArkUAPAPI.
func (api *ArkUAPAPI) Uap() *ArkUAPService {
	return api.uap
}

// Sca returns the ArkUAPSCAService instance from the ArkUAPAPI.
func (api *ArkUAPAPI) Sca() *sca.ArkUAPSCAService {
	return api.sca
}

// Db returns the ArkUAPSIADBService instance from the ArkUAPAPI.
func (api *ArkUAPAPI) Db() *db.ArkUAPSIADBService {
	return api.db
}

// VM returns the ArkUAPSIAVMService instance from the ArkUAPAPI.
func (api *ArkUAPAPI) VM() *vm.ArkUAPSIAVMService {
	return api.vm
}
