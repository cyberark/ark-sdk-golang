package db

import (
	"archive/zip"
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/common"
	"github.com/cyberark/ark-sdk-golang/pkg/common/isp"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	dbmodels "github.com/cyberark/ark-sdk-golang/pkg/services/sia/db/models"
	"github.com/cyberark/ark-sdk-golang/pkg/services/sia/sso"
	ssomodels "github.com/cyberark/ark-sdk-golang/pkg/services/sia/sso/models"
	workspacesdbmodels "github.com/cyberark/ark-sdk-golang/pkg/services/sia/workspaces/db/models"
	"github.com/golang-jwt/jwt/v5"
)

const (
	assetsURL            = "api/adb/guidance/generate"
	defaultSqlcmdTimeout = 60
)

// SIADBServiceConfig is the configuration for the ArkSIADBService.
var SIADBServiceConfig = services.ArkServiceConfig{
	ServiceName:                "sia-db",
	RequiredAuthenticatorNames: []string{"isp"},
	OptionalAuthenticatorNames: []string{},
}

// ArkSIADBService is a struct that implements the ArkService interface and provides functionality for DB service of SIA.
type ArkSIADBService struct {
	services.ArkService
	*services.ArkBaseService
	ispAuth    *auth.ArkISPAuth
	client     *isp.ArkISPServiceClient
	ssoService *sso.ArkSIASSOService
}

// NewArkSIADBService creates a new instance of ArkSIADBService with the provided authenticators.
func NewArkSIADBService(authenticators ...auth.ArkAuth) (*ArkSIADBService, error) {
	dbService := &ArkSIADBService{}
	var dbServiceInterface services.ArkService = dbService
	baseService, err := services.NewArkBaseService(dbServiceInterface, authenticators...)
	if err != nil {
		return nil, err
	}
	ispBaseAuth, err := baseService.Authenticator("isp")
	if err != nil {
		return nil, err
	}
	ispAuth := ispBaseAuth.(*auth.ArkISPAuth)
	client, err := isp.FromISPAuth(ispAuth, "dpa", ".", "", dbService.refreshSIAAuth)
	if err != nil {
		return nil, err
	}
	dbService.client = client
	dbService.ispAuth = ispAuth
	dbService.ArkBaseService = baseService
	dbService.ssoService, err = sso.NewArkSIASSOService(ispAuth)
	if err != nil {
		return nil, err
	}
	return dbService, nil
}

func (s *ArkSIADBService) refreshSIAAuth(client *common.ArkClient) error {
	err := isp.RefreshClient(client, s.ispAuth)
	if err != nil {
		return err
	}
	return nil
}

func (s *ArkSIADBService) proxyAddress(dbType string) (string, error) {
	parsedToken, _, err := new(jwt.Parser).ParseUnverified(s.client.GetToken(), jwt.MapClaims{})
	if err != nil {
		return "", err
	}
	claims := parsedToken.Claims.(jwt.MapClaims)
	return fmt.Sprintf("%s.%s.%s", claims["subdomain"], dbType, claims["platform_domain"]), nil
}

func (s *ArkSIADBService) connectionString(targetAddress string, targetUsername string, networkName string) (string, error) {
	parsedToken, _, err := new(jwt.Parser).ParseUnverified(s.client.GetToken(), jwt.MapClaims{})
	if err != nil {
		return "", err
	}
	claims := parsedToken.Claims.(jwt.MapClaims)
	addressNetwork := targetAddress
	if networkName != "" {
		addressNetwork = fmt.Sprintf("%s#%s", targetAddress, networkName)
	}
	if targetUsername != "" {
		return fmt.Sprintf("%s#%s@%s@%s", claims["unique_name"], claims["subdomain"], targetUsername, addressNetwork), nil
	}
	return fmt.Sprintf("%s#%s@%s", claims["unique_name"], claims["subdomain"], addressNetwork), nil
}

func (s *ArkSIADBService) addToPgPass(username, address, password string) error {
	passFormat := fmt.Sprintf("%s:*:*:%s:%s", address, username, password)
	path := fmt.Sprintf("%s%s.pgpass", os.Getenv("HOME"), string(os.PathSeparator))
	flags := os.O_RDWR | os.O_APPEND
	if _, err := os.Stat(path); os.IsNotExist(err) {
		flags = os.O_RDWR | os.O_CREATE
	}

	file, err := os.OpenFile(path, flags, 0600)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			s.Logger.Warning("Error closing pgpass file: %v", err)
		}
	}(file)
	var lines []string

	scanner := bufio.NewScanner(file)
	found := false
	for scanner.Scan() {
		line := scanner.Text()
		if line == passFormat {
			found = true
		}
		lines = append(lines, line)
	}

	if !found {
		lines = append(lines, passFormat)
	}

	return os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0600)
}

func (s *ArkSIADBService) removeFromPgPass(username, address, password string) error {
	passFormat := fmt.Sprintf("%s:*:*:%s:%s", address, username, password)
	path := fmt.Sprintf("%s%s.pgpass", os.Getenv("HOME"), string(os.PathSeparator))

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	file, err := os.OpenFile(path, os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			s.Logger.Warning("Error closing pgpass file: %v", err)
		}
	}(file)

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != passFormat {
			lines = append(lines, line)
		}
	}
	return os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0600)
}

func (s *ArkSIADBService) createMyLoginCnf(username string, address string, password string) (string, error) {
	tempFile, err := os.CreateTemp("", "mylogin.cnf")
	if err != nil {
		return "", err
	}
	config := fmt.Sprintf("[client]\nuser = '%s'\npassword = '%s'\nhost = '%s'\n", username, password, address)
	if _, err := tempFile.Write([]byte(config)); err != nil {
		tempFile.Close()
		return "", err
	}
	if err := tempFile.Close(); err != nil {
		return "", err
	}
	if err := os.Chmod(tempFile.Name(), 0600); err != nil {
		return "", err
	}
	return tempFile.Name(), nil
}

func (s *ArkSIADBService) execute(commandLine string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd.exe", "/C", commandLine)
	} else {
		cmd = exec.Command("sh", "-c", commandLine)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func (s *ArkSIADBService) generateAssets(
	assetType string,
	connectionMethod string,
	responseFormat string,
	generationHints map[string]interface{},
	includeSSO bool,
	resourceType string,
) (interface{}, error) {
	assetsRequest := map[string]interface{}{
		"asset_type":        assetType,
		"connection_method": connectionMethod,
		"response_format":   responseFormat,
		"generation_hints":  generationHints,
	}
	if includeSSO {
		assetsRequest["include_sso"] = includeSSO
	}
	if resourceType != "" {
		assetsRequest["resource_type"] = resourceType
	}
	response, err := s.client.Post(context.Background(), assetsURL, assetsRequest)
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
		return nil, fmt.Errorf("failed to generate assets - [%d] - [%s]", response.StatusCode, common.SerializeResponseToJSON(response.Body))
	}
	if responseFormat == dbmodels.ResponseFormatRaw {
		respBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %v", err)
		}
		return string(respBytes), nil
	}
	generatedAssets, err := common.DeserializeJSONSnake(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize response body: %v", err)
	}
	return generatedAssets, nil
}

// Psql executes a PostgreSQL command using the provided execution parameters.
func (s *ArkSIADBService) Psql(psqlExecution *dbmodels.ArkSIADBPsqlExecution) error {
	proxyAddress, err := s.proxyAddress("postgres")
	if err != nil {
		return err
	}
	connectionString, err := s.connectionString(psqlExecution.TargetAddress, psqlExecution.TargetUsername, psqlExecution.NetworkName)
	if err != nil {
		return err
	}
	password, err := s.ssoService.ShortLivedPassword(&ssomodels.ArkSIASSOGetShortLivedPassword{
		Service: "DPA-DB",
	})
	if err != nil {
		return err
	}
	executionAction := fmt.Sprintf("%s \"host=%s user=%s\"", psqlExecution.PsqlPath, proxyAddress, connectionString)

	if err := s.addToPgPass(connectionString, proxyAddress, password); err != nil {
		return err
	}
	defer func(s *ArkSIADBService, username, address, password string) {
		err := s.removeFromPgPass(username, address, password)
		if err != nil {
			s.Logger.Warning("Error removing password from pgpass: %v", err)
		}
	}(s, connectionString, proxyAddress, password)
	return s.execute(executionAction)
}

// Mysql executes a MySQL command using the provided execution parameters.
func (s *ArkSIADBService) Mysql(mysqlExecution *dbmodels.ArkSIADBMysqlExecution) error {
	proxyAddress, err := s.proxyAddress("mysql")
	if err != nil {
		return err
	}
	connectionString, err := s.connectionString(mysqlExecution.TargetAddress, mysqlExecution.TargetUsername, mysqlExecution.NetworkName)
	if err != nil {
		return err
	}
	password, err := s.ssoService.ShortLivedPassword(&ssomodels.ArkSIASSOGetShortLivedPassword{
		Service: "DPA-DB",
	})
	if err != nil {
		return err
	}
	tempCnfLogin, err := s.createMyLoginCnf(connectionString, proxyAddress, password)
	if err != nil {
		return err
	}
	executionAction := fmt.Sprintf("%s --defaults-file=%s", mysqlExecution.MysqlPath, tempCnfLogin)
	defer func() {
		err := os.Remove(tempCnfLogin)
		if err != nil {
			s.Logger.Warning("Error removing temporary .mylogin.cnf file: %v", err)
		}
	}()
	return s.execute(executionAction)
}

// Sqlcmd executes a sqlcmd command using the provided execution parameters.
func (s *ArkSIADBService) Sqlcmd(sqlcmdExecution *dbmodels.ArkSIADBSqlcmdExecution) error {
	proxyAddress, err := s.proxyAddress("mssql")
	if err != nil {
		return err
	}
	connectionString, err := s.connectionString(sqlcmdExecution.TargetAddress, sqlcmdExecution.TargetUsername, sqlcmdExecution.NetworkName)
	if err != nil {
		return err
	}
	password, err := s.ssoService.ShortLivedPassword(&ssomodels.ArkSIASSOGetShortLivedPassword{
		Service: "DPA-DB",
	})
	if err != nil {
		return err
	}
	executionAction := fmt.Sprintf("%s -U %s -S %s -l %d -P%s", sqlcmdExecution.SqlcmdPath, connectionString, proxyAddress, defaultSqlcmdTimeout, password)
	return s.execute(executionAction)
}

// GenerateOracleTnsNames generates Oracle TNS names and writes them to the specified folder.
func (s *ArkSIADBService) GenerateOracleTnsNames(generateOracleAssets *dbmodels.ArkSIADBOracleGenerateAssets) error {
	s.Logger.Info("Generating Oracle TNS names")
	assetsData, err := s.generateAssets(
		dbmodels.AssetTypeOracleTNSAssets,
		generateOracleAssets.ConnectionMethod,
		generateOracleAssets.ResponseFormat,
		map[string]interface{}{"folder": generateOracleAssets.Folder},
		generateOracleAssets.IncludeSSO,
		workspacesdbmodels.FamilyTypeOracle,
	)
	if err != nil {
		return err
	}
	if assetsDataMap, ok := assetsData.(map[string]interface{}); ok {
		assetsData = assetsDataMap["generated_assets"]
	}
	decodedAssets, err := base64.StdEncoding.DecodeString(assetsData.(string))
	if err != nil {
		return err
	}
	if _, err := os.Stat(generateOracleAssets.Folder); os.IsNotExist(err) {
		if err := os.MkdirAll(generateOracleAssets.Folder, 0755); err != nil {
			return err
		}
	}
	if !generateOracleAssets.Unzip {
		filePath := filepath.Join(generateOracleAssets.Folder, "oracle_assets.zip")
		if err := os.WriteFile(filePath, decodedAssets, 0644); err != nil {
			return err
		}
	} else {
		zipReader, err := zip.NewReader(bytes.NewReader(decodedAssets), int64(len(decodedAssets)))
		if err != nil {
			return err
		}
		for _, file := range zipReader.File {
			filePath := filepath.Join(generateOracleAssets.Folder, file.Name)
			if file.FileInfo().IsDir() {
				if err := os.MkdirAll(filePath, 0755); err != nil {
					return err
				}
				continue
			}
			outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				return err
			}
			rc, err := file.Open()
			if err != nil {
				outFile.Close()
				return err
			}
			_, err = io.Copy(outFile, rc)
			outFile.Close()
			rc.Close()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// GenerateProxyFullChain generates a proxy full chain asset and writes it to the specified folder.
func (s *ArkSIADBService) GenerateProxyFullChain(generateProxyFullChain *dbmodels.ArkSIADBProxyFullChainGenerateAssets) error {
	s.Logger.Info("Generating proxy full chain")

	assetsData, err := s.generateAssets(
		dbmodels.AssetTypeProxyFullChain,
		generateProxyFullChain.ConnectionMethod,
		generateProxyFullChain.ResponseFormat,
		nil,
		false,
		"",
	)
	if err != nil {
		return err
	}
	if assetsDataMap, ok := assetsData.(map[string]interface{}); ok {
		assetsData = assetsDataMap["generated_assets"]
	}
	if _, err := os.Stat(generateProxyFullChain.Folder); os.IsNotExist(err) {
		if err := os.MkdirAll(generateProxyFullChain.Folder, 0755); err != nil {
			return err
		}
	}

	filePath := filepath.Join(generateProxyFullChain.Folder, "proxy_fullchain.pem")
	if err := os.WriteFile(filePath, []byte(assetsData.(string)), 0644); err != nil {
		return err
	}

	return nil
}

// ServiceConfig returns the service configuration for the ArkSIADBService.
func (s *ArkSIADBService) ServiceConfig() services.ArkServiceConfig {
	return SIADBServiceConfig
}
