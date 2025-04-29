package isp

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/common"
	commonmodels "github.com/cyberark/ark-sdk-golang/pkg/models/common"
	"net/url"
	"os"
	"strings"
)

// ArkISPServiceClient is a struct that represents a client for the Ark ISP service.
type ArkISPServiceClient struct {
	*common.ArkClient
	tenantEnv commonmodels.AwsEnv
}

// NewArkISPServiceClient creates a new instance of ArkISPServiceClient.
func NewArkISPServiceClient(
	serviceName string,
	tenantSubdomain string,
	baseTenantURL string,
	tenantEnv commonmodels.AwsEnv,
	token string,
	authHeaderName string,
	separator string,
	basePath string,
	refreshConnectionCallback func(*common.ArkClient) error,
) (*ArkISPServiceClient, error) {
	if tenantEnv == "" {
		tenantEnv = commonmodels.AwsEnv(os.Getenv("DEPLOY_ENV"))
		if tenantEnv == "" {
			tenantEnv = commonmodels.Prod
		}
	}

	serviceURL, err := resolveServiceURL(serviceName, tenantSubdomain, baseTenantURL, tenantEnv, token, separator)
	if err != nil {
		return nil, err
	}
	if basePath != "" {
		serviceURL = fmt.Sprintf("%s/%s", serviceURL, basePath)
	}

	client := common.NewArkClient(serviceURL, token, "Bearer", authHeaderName, refreshConnectionCallback)
	client.SetHeader("Origin", serviceURL)
	client.SetHeader("Referer", serviceURL)
	client.SetHeader("Content-Type", "application/json")

	return &ArkISPServiceClient{
		ArkClient: client,
		tenantEnv: tenantEnv,
	}, nil
}

func resolveServiceURL(
	serviceName string,
	tenantSubdomain string,
	baseTenantURL string,
	tenantEnv commonmodels.AwsEnv,
	token string,
	separator string,
) (string, error) {
	if tenantEnv == "" {
		tenantEnv = commonmodels.AwsEnv(os.Getenv("DEPLOY_ENV"))
		if tenantEnv == "" {
			tenantEnv = commonmodels.Prod
		}
	}

	platformDomain := commonmodels.RootDomain[tenantEnv]
	var tenantChosenSubdomain string

	if token != "" {
		parsedToken, _, err := new(jwt.Parser).ParseUnverified(token, jwt.MapClaims{})
		if err != nil {
			return "", err
		}
		claims := parsedToken.Claims.(jwt.MapClaims)
		if subdomain, ok := claims["subdomain"].(string); ok {
			tenantChosenSubdomain = subdomain
		}
		if platformTokenDomain, ok := claims["platform_domain"].(string); ok {
			platformDomain = platformTokenDomain
			if strings.HasPrefix(platformDomain, "shell.") && serviceName != "" {
				platformDomain = strings.TrimPrefix(platformDomain, "shell.")
			}
			for env, domain := range commonmodels.RootDomain {
				if domain == platformDomain {
					tenantEnv = env
					break
				}
			}
		}
	}

	if tenantChosenSubdomain == "" && tenantSubdomain != "" {
		tenantChosenSubdomain = tenantSubdomain
	}

	if tenantChosenSubdomain == "" && baseTenantURL != "" {
		if !strings.HasPrefix(baseTenantURL, "https://") {
			baseTenantURL = "https://" + baseTenantURL
		}
		parsedURL, err := url.Parse(baseTenantURL)
		if err != nil {
			return "", err
		}
		tenantChosenSubdomain = strings.Split(parsedURL.Host, ".")[0]
	}

	if tenantChosenSubdomain == "" {
		parsedToken, _, err := new(jwt.Parser).ParseUnverified(token, jwt.MapClaims{})
		if err != nil {
			return "", err
		}
		claims := parsedToken.Claims.(jwt.MapClaims)
		if uniqueName, ok := claims["unique_name"].(string); ok {
			fullDomain := strings.Split(uniqueName, "@")
			if len(fullDomain) > 1 {
				domainPart := fullDomain[1]
				for env, domain := range commonmodels.RootDomain {
					if strings.Contains(domainPart, domain) {
						tenantChosenSubdomain = strings.Split(domainPart, ".")[0]
						platformDomain = domain
						tenantEnv = env
						break
					}
				}
			}
		}
	}

	if tenantChosenSubdomain == "" {
		return "", fmt.Errorf("failed to resolve tenant subdomain")
	}

	var baseURL string
	if serviceName != "" {
		baseURL = fmt.Sprintf("https://%s%s%s.%s", tenantChosenSubdomain, separator, serviceName, platformDomain)
	} else {
		baseURL = fmt.Sprintf("https://%s.%s", tenantChosenSubdomain, platformDomain)
	}

	return baseURL, nil
}

// TenantEnv returns the tenant environment of the ArkISPServiceClient.
func (client *ArkISPServiceClient) TenantEnv() commonmodels.AwsEnv {
	return client.tenantEnv
}

// TenantID returns the tenant ID from the JWT token of the ArkISPServiceClient.
func (client *ArkISPServiceClient) TenantID() (string, error) {
	if client.ArkClient.GetToken() != "" {
		parsedToken, _, err := new(jwt.Parser).ParseUnverified(client.ArkClient.GetToken(), jwt.MapClaims{})
		if err != nil {
			return "", err
		}
		claims := parsedToken.Claims.(jwt.MapClaims)
		return claims["tenant_id"].(string), nil
	}
	return "", fmt.Errorf("failed to retrieve tenant id")
}

// FromISPAuth creates a new ArkISPServiceClient from an ArkISPAuth instance.
func FromISPAuth(ispAuth *auth.ArkISPAuth, serviceName string, separator string, basePath string, refreshConnectionCallback func(*common.ArkClient) error) (*ArkISPServiceClient, error) {
	var tenantEnv commonmodels.AwsEnv
	var baseTenantURL string
	if ispAuth.Token.Username != "" {
		for env, domain := range commonmodels.RootDomain {
			if strings.Contains(ispAuth.Token.Username, domain) && strings.Contains(ispAuth.Token.Username, "@") {
				baseTenantURL = strings.Split(ispAuth.Token.Username, "@")[1]
				tenantEnv = env
				break
			}
		}
	}
	if tenantEnv == "" && ispAuth.Token.Metadata["env"] != "" {
		tenantEnv = commonmodels.AwsEnv(ispAuth.Token.Metadata["env"].(string))
	}
	if tenantEnv == "" {
		tenantEnv = commonmodels.AwsEnv(os.Getenv("DEPLOY_ENV"))
		if tenantEnv == "" {
			tenantEnv = commonmodels.Prod
		}
	}

	cookieJar := make(map[string]string)
	if cookies, ok := ispAuth.Token.Metadata["cookies"]; ok {
		decoded, _ := base64.StdEncoding.DecodeString(cookies.(string))
		_ = json.Unmarshal(decoded, &cookieJar)
	}

	return NewArkISPServiceClient(serviceName, "", baseTenantURL, tenantEnv, ispAuth.Token.Token, "Authorization", separator, basePath, refreshConnectionCallback)
}

// RefreshClient refreshes the ArkISPServiceClient with the latest authentication token and cookies.
func RefreshClient(client *common.ArkClient, ispAuth *auth.ArkISPAuth) error {
	token, err := ispAuth.LoadAuthentication(ispAuth.ActiveProfile, true)
	if err != nil {
		return err
	}
	if token != nil {
		client.UpdateToken(token.Token, client.GetTokenType())
		cookieJar := make(map[string]string)
		if cookies, ok := token.Metadata["cookies"]; ok {
			decoded, _ := base64.StdEncoding.DecodeString(cookies.(string))
			_ = json.Unmarshal(decoded, &cookieJar)
		}
		client.UpdateCookies(cookieJar)
	}
	return nil
}
