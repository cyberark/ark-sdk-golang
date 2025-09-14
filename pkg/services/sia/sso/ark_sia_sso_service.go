package sso

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/common"
	"github.com/cyberark/ark-sdk-golang/pkg/common/isp"
	"github.com/cyberark/ark-sdk-golang/pkg/common/keyring"
	authmodels "github.com/cyberark/ark-sdk-golang/pkg/models/auth"
	commonmodels "github.com/cyberark/ark-sdk-golang/pkg/models/common"
	"github.com/cyberark/ark-sdk-golang/pkg/profiles"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	ssomodels "github.com/cyberark/ark-sdk-golang/pkg/services/sia/sso/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mitchellh/mapstructure"
)

const (
	acquireSsoTokenURL = "/api/adb/sso/acquire"
	tokenSsoInfoURL    = "/api/adb/sso/info"
	sshSsoKeyURL       = "/api/ssh/sso/key"
)

// DefaultSSHFolderPath is the default folder path for SSH keys.
const (
	DefaultSSHFolderPath = "~/.ssh"
)

// ArkSIASSOService is a struct that implements the ArkService interface and provides functionality for SSO services of SIA.
type ArkSIASSOService struct {
	services.ArkService
	*services.ArkBaseService
	ispAuth      *auth.ArkISPAuth
	cacheKeyring *keyring.ArkKeyring
	client       *isp.ArkISPServiceClient
}

// NewArkSIASSOService creates a new instance of ArkSIASSOService with the provided authenticators.
func NewArkSIASSOService(authenticators ...auth.ArkAuth) (*ArkSIASSOService, error) {
	ssoService := &ArkSIASSOService{
		cacheKeyring: keyring.NewArkKeyring(ServiceConfig.ServiceName),
	}
	var ssoServiceInterface services.ArkService = ssoService
	baseService, err := services.NewArkBaseService(ssoServiceInterface, authenticators...)
	if err != nil {
		return nil, err
	}
	ispBaseAuth, err := baseService.Authenticator("isp")
	if err != nil {
		return nil, err
	}
	ispAuth := ispBaseAuth.(*auth.ArkISPAuth)
	client, err := isp.FromISPAuth(ispAuth, "dpa", ".", "", ssoService.refreshSIAAuth)
	if err != nil {
		return nil, err
	}
	ssoService.client = client
	ssoService.ispAuth = ispAuth
	ssoService.ArkBaseService = baseService
	return ssoService, nil
}

func (s *ArkSIASSOService) refreshSIAAuth(client *common.ArkClient) error {
	err := isp.RefreshClient(client, s.ispAuth)
	if err != nil {
		return err
	}
	return nil
}

func (s *ArkSIASSOService) loadFromCache(tokenType string) (*ssomodels.ArkSIASSOAcquireTokenResponse, error) {
	parsedToken, _, err := new(jwt.Parser).ParseUnverified(s.client.GetToken(), jwt.MapClaims{})
	if err != nil {
		return nil, err
	}
	claims := parsedToken.Claims.(jwt.MapClaims)
	defaultProfile, err := (*profiles.DefaultProfilesLoader()).LoadDefaultProfile()
	if err != nil {
		return nil, err
	}
	token, err := s.cacheKeyring.LoadToken(
		defaultProfile,
		fmt.Sprintf("%s_%s_sia_sso_short_lived_%s", claims["tenant_id"], claims["unique_name"], tokenType),
		false,
	)
	if err != nil {
		return nil, err
	}
	if token != nil {
		var response ssomodels.ArkSIASSOAcquireTokenResponse
		err := json.Unmarshal([]byte(token.Token), &response)
		if err != nil {
			return nil, err
		}
		return &response, nil
	}
	return nil, nil
}

func (s *ArkSIASSOService) saveToCache(result *ssomodels.ArkSIASSOAcquireTokenResponse, tokenType string) error {
	parsedToken, _, err := new(jwt.Parser).ParseUnverified(s.client.GetToken(), jwt.MapClaims{})
	if err != nil {
		return err
	}
	claims := parsedToken.Claims.(jwt.MapClaims)
	defaultProfile, err := (*profiles.DefaultProfilesLoader()).LoadDefaultProfile()
	if err != nil {
		return err
	}
	expiresAt, err := time.Parse(time.RFC3339, result.Metadata["expires_at"].(string))
	if err != nil {
		return err
	}
	createdAt, err := time.Parse(time.RFC3339, result.Metadata["created_at"].(string))
	if err != nil {
		return err
	}
	expiresIn := time.Now().Add(expiresAt.Sub(createdAt))
	marshaledToken, err := json.Marshal(result)
	if err != nil {
		return err
	}

	token := &authmodels.ArkToken{
		Token:     string(marshaledToken),
		TokenType: authmodels.Token,
		ExpiresIn: commonmodels.ArkRFC3339Time(expiresIn),
	}
	return s.cacheKeyring.SaveToken(
		defaultProfile,
		token,
		fmt.Sprintf("%s_%s_sia_sso_short_lived_%s", claims["tenant_id"], claims["unique_name"], tokenType),
		false,
	)
}

func (s *ArkSIASSOService) outputClientCertificate(folder string, outputFormat string, result *ssomodels.ArkSIASSOAcquireTokenResponse) error {
	folderPath := common.ExpandFolder(folder)
	parsedToken, _, err := new(jwt.Parser).ParseUnverified(s.client.GetToken(), jwt.MapClaims{})
	if err != nil {
		return err
	}
	claims := parsedToken.Claims.(jwt.MapClaims)
	baseName := strings.Split(claims["unique_name"].(string), "@")[0]
	clientCertificate := result.Token["client_certificate"].(string)
	privateKey := result.Token["private_key"].(string)

	switch outputFormat {
	case ssomodels.Raw:
		fmt.Printf("client-certificate-data: %s\n", clientCertificate)
		fmt.Printf("client-key-data: %s\n", privateKey)
	case ssomodels.Base64:
		fmt.Printf("client-certificate-data: %s\n", base64.StdEncoding.EncodeToString([]byte(clientCertificate)))
		fmt.Printf("client-key-data: %s\n", base64.StdEncoding.EncodeToString([]byte(privateKey)))
	case ssomodels.File:
		if folderPath == "" {
			return errors.New("folder parameter is required if format is FILE")
		}
		if _, err := os.Stat(folderPath); os.IsNotExist(err) {
			err := os.MkdirAll(folderPath, os.ModePerm)
			if err != nil {
				return err
			}
		}
		err = os.WriteFile(filepath.Join(folderPath, baseName+"_client_cert.crt"), []byte(clientCertificate), 0644)
		if err != nil {
			return err
		}
		err = os.WriteFile(filepath.Join(folderPath, baseName+"_client_key.pem"), []byte(privateKey), 0644)
		if err != nil {
			return err
		}
	case ssomodels.SingleFile:
		if folderPath == "" {
			return errors.New("folder parameter is required if format is SINGLE_FILE")
		}
		if _, err = os.Stat(folderPath); os.IsNotExist(err) {
			err = os.MkdirAll(folderPath, os.ModePerm)
			if err != nil {
				return err
			}
		}
		err = os.WriteFile(filepath.Join(folderPath, baseName+"_client_cert.pem"), []byte(clientCertificate+"\n"+privateKey), 0644)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown format [%v]", outputFormat)
	}
	return nil
}

func (s *ArkSIASSOService) saveOracleSSOWallet(folder string, unzipWallet bool, result *ssomodels.ArkSIASSOAcquireTokenResponse) error {
	folderPath := common.ExpandFolder(folder)
	wallet, err := base64.StdEncoding.DecodeString(result.Token["wallet"].(string))
	if err != nil {
		return err
	}
	if _, err = os.Stat(folderPath); os.IsNotExist(err) {
		err = os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	if !unzipWallet {
		parsedToken, _, err := new(jwt.Parser).ParseUnverified(s.client.GetToken(), jwt.MapClaims{})
		if err != nil {
			return err
		}
		claims := parsedToken.Claims.(jwt.MapClaims)
		baseName := strings.Split(claims["unique_name"].(string), "@")[0]
		err = os.WriteFile(filepath.Join(folderPath, baseName+"_wallet.zip"), wallet, 0644)
		if err != nil {
			return err
		}
	} else {
		// Unzip the wallet
		walletBytes := bytes.NewReader(wallet)
		zipReader, err := zip.NewReader(walletBytes, int64(len(wallet)))
		if err != nil {
			return err
		}
		for _, file := range zipReader.File {
			filePath := filepath.Join(folderPath, file.Name)
			if file.FileInfo().IsDir() {
				err := os.MkdirAll(filePath, os.ModePerm)
				if err != nil {
					return err
				}
				continue
			}
			if err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
				return err
			}
			destFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				return err
			}
			srcFile, err := file.Open()
			if err != nil {
				return err
			}
			_, err = io.Copy(destFile, srcFile)
			_ = destFile.Close()
			_ = srcFile.Close()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *ArkSIASSOService) saveOraclePEMWallet(folder string, result *ssomodels.ArkSIASSOAcquireTokenResponse) error {
	folderPath := common.ExpandFolder(folder)
	parsedToken, _, err := new(jwt.Parser).ParseUnverified(s.client.GetToken(), jwt.MapClaims{})
	if err != nil {
		return err
	}
	claims := parsedToken.Claims.(jwt.MapClaims)
	baseName := strings.Split(claims["unique_name"].(string), "@")[0]
	pemWallet, err := base64.StdEncoding.DecodeString(result.Token["pem_wallet"].(string))
	if err != nil {
		return err
	}
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		err := os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	err = os.WriteFile(filepath.Join(folderPath, baseName+"_ewallet.pem"), pemWallet, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (s *ArkSIASSOService) saveRDPFile(getShortLivedRDPFile *ssomodels.ArkSIASSOGetShortLivedRDPFile, result *ssomodels.ArkSIASSOAcquireTokenResponse) error {
	folderPath := common.ExpandFolder(getShortLivedRDPFile.Folder)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		err := os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	filename := fmt.Sprintf("sia _a %s", getShortLivedRDPFile.TargetAddress)
	if getShortLivedRDPFile.TargetDomain != "" {
		filename += fmt.Sprintf(" _d %s", getShortLivedRDPFile.TargetDomain)
	}
	return os.WriteFile(filepath.Join(folderPath, filename+".rdp"), []byte(result.Token["text"].(string)), 0644)
}

// ShortLivedPassword generates a short-lived password token for the user to connect.
func (s *ArkSIASSOService) ShortLivedPassword(getShortLivedPassword *ssomodels.ArkSIASSOGetShortLivedPassword) (string, error) {
	s.Logger.Info("Generating short lived password token")
	if getShortLivedPassword.AllowCaching {
		result, err := s.loadFromCache("password")
		if err == nil && result != nil {
			return result.Token["key"].(string), nil
		}
	}
	response, err := s.client.Post(context.Background(), acquireSsoTokenURL, map[string]interface{}{
		"token_type": "password",
		"service":    getShortLivedPassword.Service,
	})
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.GlobalLogger.Warning("Error closing response body")
		}
	}(response.Body)
	if response.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("failed to generate short lived password - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	var result ssomodels.ArkSIASSOAcquireTokenResponse
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return "", err
	}
	if key, ok := result.Token["key"].(string); ok {
		if getShortLivedPassword.AllowCaching {
			_ = s.saveToCache(&result, "password")
		}
		return key, nil
	}
	return "", fmt.Errorf("failed to generate short lived password - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
}

// ShortLivedClientCertificate generates a short-lived client certificate for the user to connect.
func (s *ArkSIASSOService) ShortLivedClientCertificate(getShortLivedClientCertificate *ssomodels.ArkSIASSOGetShortLivedClientCertificate) error {
	s.Logger.Info("Generating short lived client certificate")
	if getShortLivedClientCertificate.AllowCaching {
		result, err := s.loadFromCache("client_certificate")
		if err == nil && result != nil {
			return s.outputClientCertificate(getShortLivedClientCertificate.Folder, getShortLivedClientCertificate.OutputFormat, result)
		}
	}
	response, err := s.client.Post(context.Background(), acquireSsoTokenURL, map[string]interface{}{
		"token_type": "client_certificate",
		"service":    getShortLivedClientCertificate.Service,
	})
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.GlobalLogger.Warning("Error closing response body")
		}
	}(response.Body)
	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to generate short lived client certificate - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	var result ssomodels.ArkSIASSOAcquireTokenResponse
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return err
	}
	if _, ok := result.Token["client_certificate"].(string); ok {
		if _, ok := result.Token["private_key"].(string); ok {
			if getShortLivedClientCertificate.AllowCaching {
				_ = s.saveToCache(&result, "client_certificate")
			}
			return s.outputClientCertificate(getShortLivedClientCertificate.Folder, getShortLivedClientCertificate.OutputFormat, &result)
		}
	}
	return fmt.Errorf("failed to generate short lived client certificate - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
}

// ShortLivedOracleWallet generates a short-lived oracle wallet for the user to connect to oracle databases.
func (s *ArkSIASSOService) ShortLivedOracleWallet(getShortLivedOracleWallet *ssomodels.ArkSIASSOGetShortLivedOracleWallet) error {
	s.Logger.Info("Generating short lived oracle wallet")
	if getShortLivedOracleWallet.AllowCaching {
		result, err := s.loadFromCache("oracle_wallet")
		if err == nil && result != nil {
			if getShortLivedOracleWallet.WalletType == ssomodels.SSO {
				return s.saveOracleSSOWallet(getShortLivedOracleWallet.Folder, getShortLivedOracleWallet.UnzipWallet, result)
			}
			if getShortLivedOracleWallet.WalletType == ssomodels.PEM {
				return s.saveOraclePEMWallet(getShortLivedOracleWallet.Folder, result)
			}
		}
	}
	response, err := s.client.Post(context.Background(), acquireSsoTokenURL, map[string]interface{}{
		"token_type": "oracle_wallet",
		"service":    "DPA-DB",
		"token_parameters": map[string]interface{}{
			"walletType": getShortLivedOracleWallet.WalletType,
		},
	})
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.GlobalLogger.Warning("Error closing response body")
		}
	}(response.Body)
	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to generate short lived oracle wallet - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	var result ssomodels.ArkSIASSOAcquireTokenResponse
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return err
	}
	if _, ok := result.Token["wallet"].(string); ok && getShortLivedOracleWallet.WalletType == ssomodels.SSO {
		if getShortLivedOracleWallet.AllowCaching {
			_ = s.saveToCache(&result, "oracle_wallet")
		}
		return s.saveOracleSSOWallet(getShortLivedOracleWallet.Folder, getShortLivedOracleWallet.UnzipWallet, &result)
	}
	if _, ok := result.Token["pem_wallet"].(string); ok && getShortLivedOracleWallet.WalletType == ssomodels.PEM {
		if getShortLivedOracleWallet.AllowCaching {
			_ = s.saveToCache(&result, "oracle_wallet")
		}
		return s.saveOraclePEMWallet(getShortLivedOracleWallet.Folder, &result)
	}
	return fmt.Errorf("failed to generate short lived oracle wallet - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
}

// ShortLivedRdpFile generates a short-lived RDP file for the user to connect to remote desktops.
func (s *ArkSIASSOService) ShortLivedRdpFile(getShortLivedRDPFile *ssomodels.ArkSIASSOGetShortLivedRDPFile) error {
	s.Logger.Info("Generating short lived rdp file")
	if getShortLivedRDPFile.AllowCaching {
		result, err := s.loadFromCache("rdp_file")
		if err == nil && result != nil {
			return s.saveRDPFile(getShortLivedRDPFile, result)
		}
	}
	tokenParameters := map[string]interface{}{
		"targetAddress":      getShortLivedRDPFile.TargetAddress,
		"targetDomain":       getShortLivedRDPFile.TargetDomain,
		"targetUser":         getShortLivedRDPFile.TargetUser,
		"elevatedPrivileges": getShortLivedRDPFile.ElevatedPrivileges,
	}
	response, err := s.client.Post(context.Background(), acquireSsoTokenURL, map[string]interface{}{
		"token_type":            "rdp_file",
		"service":               "DPA-RDP",
		"token_parameters":      tokenParameters,
		"token_response_format": "extended",
	})
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.GlobalLogger.Warning("Error closing response body")
		}
	}(response.Body)
	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to generate short lived rdp file - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	var result ssomodels.ArkSIASSOAcquireTokenResponse
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return err
	}
	if _, ok := result.Token["text"].(string); ok {
		if getShortLivedRDPFile.AllowCaching {
			_ = s.saveToCache(&result, "rdp_file")
		}
		return s.saveRDPFile(getShortLivedRDPFile, &result)
	}
	return fmt.Errorf("failed to generate short rdp file - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
}

// ShortLivedSSHKey generates a short-lived SSH key for the user to connect to remote servers.
func (s *ArkSIASSOService) ShortLivedSSHKey(getSSHKey *ssomodels.ArkSIASSOGetSSHKey) (string, error) {
	s.Logger.Info("Getting short lived ssh sso key")
	response, err := s.client.Get(context.Background(), sshSsoKeyURL, nil)
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
		return "", fmt.Errorf("failed to get short lived ssh sso key - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	folderPath := getSSHKey.Folder
	if folderPath == "" {
		folderPath = DefaultSSHFolderPath
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
	parsedToken, _, err := new(jwt.Parser).ParseUnverified(s.client.GetToken(), jwt.MapClaims{})
	if err != nil {
		return "", err
	}
	claims := parsedToken.Claims.(jwt.MapClaims)
	baseName := fmt.Sprintf("sia_ssh_key_%s.pem", strings.Split(claims["unique_name"].(string), "@")[0])
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

// ShortLivedTokenInfo retrieves information about a short-lived token.
func (s *ArkSIASSOService) ShortLivedTokenInfo(getTokenInfo *ssomodels.ArkSIASSOGetTokenInfo) (*ssomodels.ArkSIASSOTokenInfo, error) {
	s.Logger.Info("Getting short lived token info")
	getTokenInfoParams := map[string]string{}
	_ = mapstructure.Decode(getTokenInfo, &getTokenInfoParams)
	response, err := s.client.Get(context.Background(), tokenSsoInfoURL, getTokenInfoParams)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.GlobalLogger.Warning("Error closing response body")
		}
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get short lived token info - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	var tokenInfo ssomodels.ArkSIASSOTokenInfo
	err = json.NewDecoder(response.Body).Decode(&tokenInfo)
	if err != nil {
		s.Logger.Error("Failed to parse get short lived token info response [%s] - [%s]", err.Error(), common.SerializeResponseToJSON(response.Body))
		return nil, fmt.Errorf("failed to parse get short lived token info response [%s]", err.Error())
	}
	return &tokenInfo, nil
}

// ServiceConfig returns the service configuration for the ArkSIASSOService.
func (s *ArkSIASSOService) ServiceConfig() services.ArkServiceConfig {
	return ServiceConfig
}
