package actions

// Package actions provides constants and configuration settings for the Ark CLI actions.
var (
	ConfigurationIgnoredDefinitionKeys = []string{
		"auth-profiles",
		"auth-method",
		"auth-method-settings",
		"user-param-name",
		"password-param-name",
		"identity-mfa-interactive",
	}

	ConfigurationAuthenticatorIgnoredDefinitionKeys = map[string][]string{
		"isp": {"identity-application", "identity-tenant-url"},
	}

	ConfigurationIgnoredInteractiveKeys = []string{
		"raw",
		"silent",
		"verbose",
		"profile_name",
		"auth_profiles",
		"auth_method",
		"auth_method_settings",
		"interactive",
	}

	ConfigurationAuthenticatorIgnoredInteractiveKeys = map[string][]string{
		"isp": {
			"identity_application",
			"identity_application_id",
			"identity_authorization_application",
			"identity_tenant_url",
		},
	}

	ConfigurationAllowedEmptyValues = []string{
		"isp_identity_url",
		"isp_identity_tenant_subdomain",
	}

	ConfigurationAuthenticatorsDefaults = map[string]string{}

	ConfigurationOverrideAliases = map[string]string{
		"region": "Region",
		"env":    "Environment",
	}
)
