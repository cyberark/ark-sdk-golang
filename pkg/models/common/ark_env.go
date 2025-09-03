// Package common provides common models and utilities for AWS environment management
// and configuration handling within the ARK SDK. This package contains environment
// type definitions, environment detection utilities, and mapping configurations
// for different AWS environments including production and government cloud deployments.
package common

import (
	"os"
	"regexp"
	"strings"
)

// AwsEnv represents the AWS environment type used throughout the ARK SDK.
//
// This type is used to distinguish between different AWS deployment environments
// such as production and government cloud environments. It provides type safety
// when working with environment-specific configurations and mappings.
type AwsEnv string

// Supported AWS environments for ARK SDK deployments.
//
// These constants define the available AWS environments that the ARK SDK
// can operate within. Each environment has specific configurations and
// endpoint mappings defined in the associated maps below.
const (
	// Prod represents the standard AWS production environment.
	Prod AwsEnv = "prod"
	// GovProd represents the AWS GovCloud production environment.
	GovProd AwsEnv = "gov-prod"
)

// Environment variable and tenant configuration constants.
//
// These constants define the standard environment variables and default
// values used for environment detection and tenant configuration across
// different AWS environments.
const (
	// DeployEnv is the environment variable name used to determine the current deployment environment.
	DeployEnv = "DEPLOY_ENV"
	// IdentityTenantName is the default tenant name used for identity services.
	IdentityTenantName = "isp"
)

// RootDomain maps AWS environments to their respective root domain names.
//
// This mapping provides the base domain for each AWS environment, which is used
// to construct various service endpoints and URLs throughout the ARK SDK.
// The root domains differ between standard AWS and GovCloud environments.
var RootDomain = map[AwsEnv]string{
	Prod:    "cyberark.cloud",
	GovProd: "cyberarkgov.cloud",
}

// IdentityEnvUrls maps AWS environments to their respective identity service URLs.
//
// This mapping provides the identity service endpoints for each AWS environment.
// These URLs are used for authentication and identity management operations
// and vary between standard AWS and GovCloud deployments.
var IdentityEnvUrls = map[AwsEnv]string{
	Prod:    "idaptive.app",
	GovProd: "id.cyberarkgov.cloud",
}

// IdentityTenantNames maps AWS environments to their respective identity tenant names.
//
// This mapping provides the default tenant names used for identity services
// in each AWS environment. Currently, both environments use the same default
// tenant name, but this mapping allows for environment-specific customization.
var IdentityTenantNames = map[AwsEnv]string{
	Prod:    IdentityTenantName,
	GovProd: IdentityTenantName,
}

// IdentityGeneratedSuffixPattern maps AWS environments to their respective regex patterns.
//
// These patterns are used to validate and identify auto-generated identity suffixes
// for each AWS environment. The patterns help distinguish between different
// environment-specific tenant naming conventions and ensure proper tenant routing.
var IdentityGeneratedSuffixPattern = map[AwsEnv]string{
	Prod:    `cyberark\.cloud\.\d.*`,
	GovProd: `cyberarkgov\.cloud\.\d.*`,
}

// GetDeployEnv returns the current AWS environment based on the DEPLOY_ENV environment variable.
//
// This function reads the DEPLOY_ENV environment variable to determine the current
// deployment environment. If the environment variable is not set or is empty,
// it defaults to the production environment for backward compatibility.
//
// Returns the AwsEnv corresponding to the current deployment environment.
//
// Example:
//
//	// Set environment variable
//	os.Setenv("DEPLOY_ENV", "gov-prod")
//	env := GetDeployEnv()
//	if env == GovProd {
//	    // Handle GovCloud-specific logic
//	}
//
//	// Default behavior when not set
//	os.Unsetenv("DEPLOY_ENV")
//	env = GetDeployEnv() // Returns Prod
func GetDeployEnv() AwsEnv {
	deployEnv := os.Getenv(DeployEnv)
	if deployEnv == "" {
		return Prod
	}
	return AwsEnv(deployEnv)
}

// CheckIfIdentityGeneratedSuffix validates if a tenant suffix matches the environment-specific pattern.
//
// This function checks whether the provided tenant suffix matches the expected
// pattern for auto-generated identity suffixes in the specified AWS environment.
// It uses regex patterns defined in IdentityGeneratedSuffixPattern to perform
// the validation, helping to ensure proper tenant routing and identification.
//
// Parameters:
//   - tenantSuffix: The tenant suffix string to validate against the pattern
//   - env: The AWS environment to check the pattern against
//
// Returns true if the tenant suffix matches the environment's pattern, false otherwise.
// Returns false if the environment is not recognized or the pattern match fails.
//
// Example:
//
//	// Check production environment suffix
//	isValid := CheckIfIdentityGeneratedSuffix("cyberark.cloud.123", Prod)
//	if isValid {
//	    // Handle auto-generated tenant
//	}
//
//	// Check GovCloud environment suffix
//	isValid = CheckIfIdentityGeneratedSuffix("cyberarkgov.cloud.456", GovProd)
func CheckIfIdentityGeneratedSuffix(tenantSuffix string, env AwsEnv) bool {
	pattern, exists := IdentityGeneratedSuffixPattern[env]
	if !exists {
		return false
	}
	matched, _ := regexp.MatchString(pattern, tenantSuffix)
	return matched
}

// IsGovCloud determines if the current AWS region is a government cloud region.
//
// This function checks the AWS region environment variables to determine if the
// current deployment is running in an AWS GovCloud region. It first checks the
// AWS_REGION environment variable, and if that's not set, falls back to checking
// AWS_DEFAULT_REGION. GovCloud regions are identified by the "us-gov" prefix.
//
// Returns true if the current region is a GovCloud region, false otherwise.
// Returns false if no region environment variables are set.
//
// Example:
//
//	// Set GovCloud region
//	os.Setenv("AWS_REGION", "us-gov-west-1")
//	if IsGovCloud() {
//	    // Configure for GovCloud deployment
//	    env := GovProd
//	}
//
//	// Standard AWS region
//	os.Setenv("AWS_REGION", "us-east-1")
//	if !IsGovCloud() {
//	    // Configure for standard AWS deployment
//	    env := Prod
//	}
func IsGovCloud() bool {
	regionName := os.Getenv("AWS_REGION")
	if regionName == "" {
		regionName = os.Getenv("AWS_DEFAULT_REGION")
	}
	return strings.HasPrefix(regionName, "us-gov")
}
