package sshca

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/common"
	"github.com/cyberark/ark-sdk-golang/pkg/common/isp"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	sshcamodels "github.com/cyberark/ark-sdk-golang/pkg/services/sia/sshca/models"
)

const (
	generateNewCAKeyURL        = "api/public-keys/rotation/generate-new"
	deactivatePreviousCAKeyURL = "api/public-keys/rotation/deactivate-previous"
	reactivatePreviousCAKeyURL = "api/public-keys/rotation/reactivate-previous"
	publicKeyURL               = "api/public-keys"
	publicKeyScriptURL         = "api/public-keys/scripts"
)

// SIASSHCAServiceConfig is the configuration for the ArkSIASSHCAService.
var SIASSHCAServiceConfig = services.ArkServiceConfig{
	ServiceName:                "sia-ssh-ca",
	RequiredAuthenticatorNames: []string{"isp"},
	OptionalAuthenticatorNames: []string{},
}

// ArkSIASSHCAService is a struct that implements the ArkService interface and provides functionality for SSH CA of SIA.
type ArkSIASSHCAService struct {
	services.ArkService
	*services.ArkBaseService
	ispAuth *auth.ArkISPAuth
	client  *isp.ArkISPServiceClient
}

// NewArkSIASSHCAService creates a new instance of ArkSIASSHCAService with the provided authenticators.
func NewArkSIASSHCAService(authenticators ...auth.ArkAuth) (*ArkSIASSHCAService, error) {
	sshCaService := &ArkSIASSHCAService{}
	var sshCaServiceInterface services.ArkService = sshCaService
	baseService, err := services.NewArkBaseService(sshCaServiceInterface, authenticators...)
	if err != nil {
		return nil, err
	}
	ispBaseAuth, err := baseService.Authenticator("isp")
	if err != nil {
		return nil, err
	}
	ispAuth := ispBaseAuth.(*auth.ArkISPAuth)
	client, err := isp.FromISPAuth(ispAuth, "dpa", ".", "", sshCaService.refreshSIAAuth)
	if err != nil {
		return nil, err
	}
	sshCaService.client = client
	sshCaService.ispAuth = ispAuth
	sshCaService.ArkBaseService = baseService
	return sshCaService, nil
}

func (s *ArkSIASSHCAService) refreshSIAAuth(client *common.ArkClient) error {
	err := isp.RefreshClient(client, s.ispAuth)
	if err != nil {
		return err
	}
	return nil
}

// GenerateNewCA generates a new CA key version.
func (s *ArkSIASSHCAService) GenerateNewCA() error {
	s.Logger.Info("Generate new CA key version")
	response, err := s.client.Post(context.Background(), generateNewCAKeyURL, nil)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf("Failed to generate new CA key  - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	return nil
}

// DeactivatePreviousCa Deactivate previous CA key version.
func (s *ArkSIASSHCAService) DeactivatePreviousCa() error {
	s.Logger.Info("Deactivate previous CA key version")
	response, err := s.client.Post(context.Background(), deactivatePreviousCAKeyURL, nil)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to deactivate previous CA key  - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	return nil
}

// ReactivatePreviousCa Deactivate previous CA key version.
func (s *ArkSIASSHCAService) ReactivatePreviousCa() error {
	s.Logger.Info("Reactivate previous CA key version")
	response, err := s.client.Post(context.Background(), reactivatePreviousCAKeyURL, nil)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to reactivate previous CA key  - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	return nil
}

// PublicKey retrieves the public key for the SSH CA.
func (s *ArkSIASSHCAService) PublicKey(getPublicKey *sshcamodels.ArkSIAGetSSHPublicKey) (string, error) {
	s.Logger.Info("Getting public key")
	response, err := s.client.Get(context.Background(), publicKeyURL, nil)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.GlobalLogger.Warning("Error closing response body")
		}
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Failed to get public key  - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	publicKey, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	if getPublicKey != nil && getPublicKey.OutputFile != "" {
		file, err := os.Create(getPublicKey.OutputFile)
		if err != nil {
			return "", err
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				s.Logger.Warning("Error closing output file")
			}
		}(file)
		_, err = file.Write(publicKey)
		if err != nil {
			return "", err
		}
	}
	return string(publicKey), nil
}

// PublicKeyScript retrieves the public key script for the SSH CA.
func (s *ArkSIASSHCAService) PublicKeyScript(getPublicKey *sshcamodels.ArkSIAGetSSHPublicKey) (string, error) {
	s.Logger.Info("Getting public key script")
	response, err := s.client.Get(context.Background(), publicKeyScriptURL, nil)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.GlobalLogger.Warning("Error closing response body")
		}
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Failed to get public key script  - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	publicKeyScript, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	if getPublicKey != nil && getPublicKey.OutputFile != "" {
		file, err := os.Create(getPublicKey.OutputFile)
		if err != nil {
			return "", err
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				s.Logger.Warning("Error closing output file")
			}
		}(file)
		_, err = file.Write(publicKeyScript)
		if err != nil {
			return "", err
		}
	}
	return string(publicKeyScript), nil
}

// ServiceConfig returns the service configuration for the ArkSIASSHCAService.
func (s *ArkSIASSHCAService) ServiceConfig() services.ArkServiceConfig {
	return SIASSHCAServiceConfig
}
