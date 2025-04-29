package common

import (
	"os"
	"regexp"
	"strings"
)

// AwsEnv is a string type that represents the AWS environment.
type AwsEnv string

// Constants for AWS environments.
const (
	Prod    AwsEnv = "prod"
	GovProd AwsEnv = "gov-prod"
)

// Constants for environment variables and tenant names.
const (
	DeployEnv          = "DEPLOY_ENV"
	IdentityTenantName = "isp"
)

// RootDomain is a map that associates AWS environments with their respective root domains.
var RootDomain = map[AwsEnv]string{
	Prod:    "cyberark.cloud",
	GovProd: "cyberarkgov.cloud",
}

// IdentityEnvUrls is a map that associates AWS environments with their respective identity environment URLs.
var IdentityEnvUrls = map[AwsEnv]string{
	Prod:    "idaptive.app",
	GovProd: "id.cyberarkgov.cloud",
}

// IdentityTenantNames is a map that associates AWS environments with their respective identity tenant names.
var IdentityTenantNames = map[AwsEnv]string{
	Prod:    IdentityTenantName,
	GovProd: IdentityTenantName,
}

// IdentityGeneratedSuffixPattern is a map that associates AWS environments with their respective identity generated suffix patterns.
var IdentityGeneratedSuffixPattern = map[AwsEnv]string{
	Prod:    `cyberark\.cloud\.\d.*`,
	GovProd: `cyberarkgov\.cloud\.\d.*`,
}

// GetDeployEnv returns the AWS environment based on the DEPLOY_ENV environment variable.
func GetDeployEnv() AwsEnv {
	deployEnv := os.Getenv(DeployEnv)
	if deployEnv == "" {
		return Prod
	}
	return AwsEnv(deployEnv)
}

// CheckIfIdentityGeneratedSuffix checks if the given tenant suffix matches the identity generated suffix pattern for the specified AWS environment.
func CheckIfIdentityGeneratedSuffix(tenantSuffix string, env AwsEnv) bool {
	pattern, exists := IdentityGeneratedSuffixPattern[env]
	if !exists {
		return false
	}
	matched, _ := regexp.MatchString(pattern, tenantSuffix)
	return matched
}

// IsGovCloud checks if the current AWS region is a GovCloud region.
func IsGovCloud() bool {
	regionName := os.Getenv("AWS_REGION")
	if regionName == "" {
		regionName = os.Getenv("AWS_DEFAULT_REGION")
	}
	return strings.HasPrefix(regionName, "us-gov")
}
