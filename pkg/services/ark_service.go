package services

import (
	"fmt"
	"slices"

	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/common"
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
)

// ArkServiceConfig defines the configuration for an Ark service.
type ArkServiceConfig struct {
	ServiceName                string
	RequiredAuthenticatorNames []string
	OptionalAuthenticatorNames []string
	ActionsConfigurations      map[actions.ArkServiceActionType][]actions.ArkServiceActionDefinition
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

var (
	serviceRegistry  = make(map[string]ArkServiceConfig)
	topLevelServices []string
)

// Register registers a new Ark service configuration.
func Register(serviceConfig ArkServiceConfig, topLevel bool) error {
	if _, exists := serviceRegistry[serviceConfig.ServiceName]; exists {
		return fmt.Errorf("service %s already registered", serviceConfig.ServiceName)
	}
	serviceRegistry[serviceConfig.ServiceName] = serviceConfig
	if topLevel {
		topLevelServices = append(topLevelServices, serviceConfig.ServiceName)
	}
	return nil
}

// GetServiceConfig retrieves the Ark service configuration by service name.
func GetServiceConfig(serviceName string) (ArkServiceConfig, error) {
	if config, exists := serviceRegistry[serviceName]; exists {
		return config, nil
	}
	return ArkServiceConfig{}, fmt.Errorf("service %s not registered", serviceName)
}

// AllServiceConfigs returns a slice of all registered Ark service configurations.
func AllServiceConfigs() []ArkServiceConfig {
	configs := make([]ArkServiceConfig, 0, len(serviceRegistry))
	for _, config := range serviceRegistry {
		configs = append(configs, config)
	}
	return configs
}

// TopLevelServiceConfigs returns a slice of all registered top-level Ark service configurations.
func TopLevelServiceConfigs() []ArkServiceConfig {
	configs := make([]ArkServiceConfig, 0, len(topLevelServices))
	for _, serviceName := range topLevelServices {
		if config, exists := serviceRegistry[serviceName]; exists {
			configs = append(configs, config)
		}
	}
	return configs
}
