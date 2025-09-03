package common

import (
	"os"
	"testing"
)

func TestGetDeployEnv(t *testing.T) {
	tests := []struct {
		name        string
		envVarValue string
		expectedEnv AwsEnv
	}{
		{
			name:        "returns_prod_when_env_var_not_set",
			envVarValue: "",
			expectedEnv: Prod,
		},
		{
			name:        "returns_prod_when_env_var_set_to_prod",
			envVarValue: "prod",
			expectedEnv: Prod,
		},
		{
			name:        "returns_gov_prod_when_env_var_set_to_gov_prod",
			envVarValue: "gov-prod",
			expectedEnv: GovProd,
		},
		{
			name:        "returns_unknown_when_env_var_set_to_unknown",
			envVarValue: "unknown",
			expectedEnv: AwsEnv("unknown"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv(DeployEnv, tt.envVarValue)
			defer os.Unsetenv(DeployEnv)

			result := GetDeployEnv()
			if result != tt.expectedEnv {
				t.Errorf("Expected %v, got %v", tt.expectedEnv, result)
			}
		})
	}
}

func TestCheckIfIdentityGeneratedSuffix(t *testing.T) {
	tests := []struct {
		name         string
		tenantSuffix string
		env          AwsEnv
		expected     bool
	}{
		{
			name:         "matches_pattern_in_prod",
			tenantSuffix: "cyberark.cloud.12345",
			env:          Prod,
			expected:     true,
		},
		{
			name:         "does_not_match_pattern_in_prod",
			tenantSuffix: "ab123",
			env:          Prod,
			expected:     false,
		},
		{
			name:         "matches_pattern_in_gov",
			tenantSuffix: "cyberarkgov.cloud.67890",
			env:          GovProd,
			expected:     true,
		},
		{
			name:         "does_not_match_pattern_in_gov",
			tenantSuffix: "gov-ab123",
			env:          GovProd,
			expected:     false,
		},
		{
			name:         "unknown_env_returns_false",
			tenantSuffix: "abc12345",
			env:          AwsEnv("unknown"),
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckIfIdentityGeneratedSuffix(tt.tenantSuffix, tt.env)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestIsGovCloud(t *testing.T) {
	tests := []struct {
		name           string
		envVarName     string
		envVarValue    string
		expectedResult bool
	}{
		{
			name:           "true_when_AWS_REGION_starts_with_us-gov",
			envVarName:     "AWS_REGION",
			envVarValue:    "us-gov-1",
			expectedResult: true,
		},
		{
			name:           "true_when_AWS_DEFAULT_REGION_starts_with_us-gov",
			envVarName:     "AWS_DEFAULT_REGION",
			envVarValue:    "us-gov-1",
			expectedResult: true,
		},
		{
			name:           "false_when_AWS_REGION_starts_with_eu-west",
			envVarName:     "AWS_REGION",
			envVarValue:    "eu-west-1",
			expectedResult: false,
		},
		{
			name:           "false_when_AWS_DEFAULT_REGION_starts_with_eu-west",
			envVarName:     "AWS_REGION",
			envVarValue:    "eu-west-1",
			expectedResult: false,
		},
		{
			name:           "false_when_AWS_REGION_empty",
			envVarName:     "AWS_REGION",
			envVarValue:    "",
			expectedResult: false,
		},
		{
			name:           "false_when_AWS_DEFAULT_REGION_empty",
			envVarName:     "AWS_DEFAULT_REGION",
			envVarValue:    "",
			expectedResult: false,
		},
		{
			name:           "false_when_no_env_vars_set",
			envVarName:     "",
			envVarValue:    "",
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv(tt.envVarName, tt.envVarValue)
			defer os.Unsetenv(tt.envVarName)

			result := IsGovCloud()

			if result != tt.expectedResult {
				t.Errorf("Expected IsGovCloud() to return %v, got %v", tt.expectedResult, result)
			}
		})
	}
}
