package identity

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/cyberark/ark-sdk-golang/pkg/common"
	commonmodels "github.com/cyberark/ark-sdk-golang/pkg/models/common"
	"github.com/cyberark/ark-sdk-golang/pkg/models/common/identity"
)

const (
	discoveryServiceDomainName = "platform-discovery"
	discoveryTimeout           = 30
)

// DefaultHeaders returns the default headers for HTTP requests to identity.
func DefaultHeaders() map[string]string {
	return map[string]string{
		"Content-Type":         "application/json",
		"X-IDAP-NATIVE-CLIENT": "true",
		"User-Agent":           common.UserAgent(),
		"OobIdPAuth":           "true",
	}
}

// DefaultSystemHeaders returns the default system headers for HTTP requests to identity.
func DefaultSystemHeaders() map[string]string {
	return map[string]string{
		"X-IDAP-NATIVE-CLIENT": "true",
		"User-Agent":           common.UserAgent(),
	}
}

// ResolveTenantFqdnFromTenantSubdomain resolves the tenant's FQDN URL from its subdomain.
// The resolved URL is based on the current working environment, which is provided in the `tenantSubdomain` argument.
func ResolveTenantFqdnFromTenantSubdomain(tenantSubdomain string, env commonmodels.AwsEnv) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), discoveryTimeout*time.Second)
	defer cancel()
	client := common.NewSimpleArkClient(fmt.Sprintf("https://%s.%s", discoveryServiceDomainName, commonmodels.RootDomain[env]))
	client.SetHeaders(map[string]string{
		"Content-Type": "application/json",
	})
	response, err := client.Get(ctx, fmt.Sprintf("api/identity-endpoint/%s", tenantSubdomain), nil)
	if err != nil {
		return "", fmt.Errorf("getting tenant FQDN failed from platform discovery: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.GlobalLogger.Warning("Error closing response body")
		}
	}(response.Body)
	if response.StatusCode == http.StatusOK {
		var parsedResponse identity.TenantEndpointResponse
		if err := json.NewDecoder(response.Body).Decode(&parsedResponse); err != nil {
			return "", fmt.Errorf("getting tenant FQDN failed from platform discovery to be parsed / validated")
		}
		return parsedResponse.Endpoint, nil
	}
	return "", fmt.Errorf("getting tenant FQDN failed from platform discovery [%d] - [%s]", response.StatusCode, response.Status)
}

// ResolveTenantFqdnFromTenantSuffix resolves the tenant's FQDN URL from its suffix.
func ResolveTenantFqdnFromTenantSuffix(tenantSuffix string, identityEnvURL string) (string, error) {
	if identityEnvURL == "" {
		identityEnvURL = commonmodels.IdentityEnvUrls[commonmodels.GetDeployEnv()]
	}
	client := common.NewSimpleArkClient(fmt.Sprintf("https://pod0.%s", identityEnvURL))
	client.SetHeaders(map[string]string{
		"Content-Type":          "application/json",
		"X-IDAP-NATIVE-CLIENT'": "true",
	})
	body := map[string]interface{}{
		"User":                  tenantSuffix,
		"Version":               "1.0",
		"PlatformTokenResponse": true,
		"MfaRequestor":          "DeviceAgent",
	}
	response, err := client.Post(context.Background(), "Security/StartAuthentication", body)
	if err != nil {
		return "", err
	}
	if response.StatusCode == http.StatusOK {
		var parsedResponse identity.TenantFqdnResponse
		if err := json.NewDecoder(response.Body).Decode(&parsedResponse); err != nil {
			return "", fmt.Errorf("getting tenant FQDN failed from identity to be parsed / validated")
		}
		if !parsedResponse.Success || parsedResponse.Result.PodFqdn == "" {
			return "", fmt.Errorf("getting tenant FQDN failed from identity: %s", parsedResponse.Message)
		}
		fqdn := parsedResponse.Result.PodFqdn
		if !strings.HasPrefix(fqdn, "https://") {
			fqdn = "https://" + fqdn
		}
		return fqdn, nil
	}
	return "", fmt.Errorf("getting tenant FQDN failed from identity [%d] - [%s]", response.StatusCode, response.Status)
}
