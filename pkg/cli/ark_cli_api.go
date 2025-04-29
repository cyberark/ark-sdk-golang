package cli

import (
	api "github.com/cyberark/ark-sdk-golang/pkg"
	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/models"
)

// ArkCLIAPI is a struct that represents the Ark CLI API client.
type ArkCLIAPI struct {
	api.ArkAPI
}

// NewArkCLIAPI creates a new instance of ArkCLIAPI.
func NewArkCLIAPI(authenticators []auth.ArkAuth, profile *models.ArkProfile) (*ArkCLIAPI, error) {
	arkAPI, err := api.NewArkAPI(authenticators, profile)
	if err != nil {
		return nil, err
	}
	return &ArkCLIAPI{
		ArkAPI: *arkAPI,
	}, nil
}
