package services

import (
	"fmt"
	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/common"
	"slices"
)

// ArkServiceConfig defines the configuration for an Ark service.
type ArkServiceConfig struct {
	ServiceName                string
	RequiredAuthenticatorNames []string
	OptionalAuthenticatorNames []string
}

// ArkService is an interface that defines the methods for an Ark service.
type ArkService interface {
	ServiceConfig() ArkServiceConfig
}

// ArkBaseService is a struct that implements the ArkService interface and provides base functionality for Ark services.
type ArkBaseService struct {
	Service        ArkService
	Logger         *common.ArkLogger
	authenticators []auth.ArkAuth
}

// NewArkBaseService creates a new instance of ArkBaseService with the provided service and authenticators.
func NewArkBaseService(service ArkService, authenticators ...auth.ArkAuth) (*ArkBaseService, error) {
	baseService := &ArkBaseService{
		Service:        service,
		Logger:         common.GetLogger("ArkBaseService", common.Unknown),
		authenticators: make([]auth.ArkAuth, 0),
	}

	for _, authenticator := range authenticators {
		baseService.authenticators = append(baseService.authenticators, authenticator)
	}

	var givenAuthNames []string
	for _, authenticator := range baseService.authenticators {
		givenAuthNames = append(givenAuthNames, authenticator.AuthenticatorName())
	}

	config := service.ServiceConfig()
	for _, requiredAuth := range config.RequiredAuthenticatorNames {
		if !slices.Contains(givenAuthNames, requiredAuth) {
			return nil, fmt.Errorf("%s missing required authenticators for service", config.ServiceName)
		}
	}

	return baseService, nil
}

// Authenticators returns the list of authenticators for the ArkBaseService.
func (s *ArkBaseService) Authenticators() []auth.ArkAuth {
	return s.authenticators
}

// Authenticator returns the authenticator with the specified name from the ArkBaseService.
func (s *ArkBaseService) Authenticator(authName string) (auth.ArkAuth, error) {
	for _, authenticator := range s.authenticators {
		if authenticator.AuthenticatorName() == authName {
			return authenticator, nil
		}
	}
	return nil, fmt.Errorf("%s Failed to find authenticator %s", s.Service.ServiceConfig().ServiceName, authName)
}

// HasAuthenticator checks if the ArkBaseService has an authenticator with the specified name.
func (s *ArkBaseService) HasAuthenticator(authName string) bool {
	for _, authenticator := range s.authenticators {
		if authenticator.AuthenticatorName() == authName {
			return true
		}
	}
	return false
}
