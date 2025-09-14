package k8s

import (
	"context"
	"errors"
	"fmt"

	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/common"
	"github.com/cyberark/ark-sdk-golang/pkg/common/isp"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	k8smodels "github.com/cyberark/ark-sdk-golang/pkg/services/sia/k8s/models"

	"io"
	"net/http"
	"os"
	"path/filepath"
)

const (
	kubeConfigGenerationURL = "/api/k8s/kube-config"
)

// DefaultKubeConfigFolderPath is the default folder path for kubeconfig files.
const (
	DefaultKubeConfigFolderPath = "~/.kube"
)

// ArkSIAK8SService is a struct that implements the ArkService interface and provides functionality for K8S service of SIA.
type ArkSIAK8SService struct {
	services.ArkService
	*services.ArkBaseService
	ispAuth *auth.ArkISPAuth
	client  *isp.ArkISPServiceClient
}

// NewArkSIAK8SService creates a new instance of ArkSIAK8SService with the provided authenticators.
func NewArkSIAK8SService(authenticators ...auth.ArkAuth) (*ArkSIAK8SService, error) {
	k8sService := &ArkSIAK8SService{}
	var k8sServiceInterface services.ArkService = k8sService
	baseService, err := services.NewArkBaseService(k8sServiceInterface, authenticators...)
	if err != nil {
		return nil, err
	}
	ispBaseAuth, err := baseService.Authenticator("isp")
	if err != nil {
		return nil, err
	}
	ispAuth := ispBaseAuth.(*auth.ArkISPAuth)
	client, err := isp.FromISPAuth(ispAuth, "dpa", ".", "", k8sService.refreshSIAAuth)
	if err != nil {
		return nil, err
	}
	k8sService.client = client
	k8sService.ispAuth = ispAuth
	k8sService.ArkBaseService = baseService
	return k8sService, nil
}

func (s *ArkSIAK8SService) refreshSIAAuth(client *common.ArkClient) error {
	err := isp.RefreshClient(client, s.ispAuth)
	if err != nil {
		return err
	}
	return nil
}

// GenerateKubeconfig generates a kubeconfig file for the SIA K8S service and saves it to the specified folder.
func (s *ArkSIAK8SService) GenerateKubeconfig(generateKubeConfig *k8smodels.ArkSIAK8SGenerateKubeconfig) (string, error) {
	s.Logger.Info("Getting kubeconfig")
	response, err := s.client.Get(context.Background(), kubeConfigGenerationURL, nil)
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
		return "", fmt.Errorf("failed to get kubeconfig - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	folderPath := generateKubeConfig.Folder
	if folderPath == "" {
		folderPath = DefaultKubeConfigFolderPath
	}
	folderPath = common.ExpandFolder(folderPath)
	if folderPath == "" {
		return "", errors.New("folder parameter is required")
	}
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		err := os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			return "", err
		}
	}
	baseName := "config"
	fullPath := filepath.Join(folderPath, baseName)
	resp, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	err = os.WriteFile(fullPath, resp, 0644)
	if err != nil {
		return "", err
	}
	return fullPath, nil
}

// ServiceConfig returns the service configuration for the ArkSIAK8SService.
func (s *ArkSIAK8SService) ServiceConfig() services.ArkServiceConfig {
	return ServiceConfig
}
