package configuration

// ArkSecHubSyncSettinsgs represents the sync settings within the configuration settings.
type ArkSecHubSyncSettings struct {
	SecretValidity int `json:"secret_validity" mapstructure:"secret_validity" desc:"Secret Validity period in seconds" flag:"secret-validity" validate:"required"`
}

type ArkSecHubSecretsSource struct {
	// The tenant PAM type Exxample: PCLOUD_NON_UM, PAM_SELF_HOSTED
	PAMType string `json:"pam_type" mapstructure:"pam_type" desc:"PAM Type for Secrets Source" flag:"pam-type"`
}

type ArkSecHubAuthenticationIdentities struct {
	// Identities defines the authentication identities for the secrets hub.
	AWS ArkSecHubAuthenticationIdentitiesAWS `json:"aws" mapstructure:"aws" desc:"AWS Authentication Identities for Secrets Hub" flag:"aws-authentication-identities"`
}

type ArkSecHubAuthenticationIdentitiesAWS struct {
	// The Secrets Hub tenant role ARN
	TenantRoleARN string `json:"tenant_role_arn" mapstructure:"tenant_role_arn" desc:"The Secrets Hub tenant role ARN" flag:"tenant-role-arn"`
}

// ArkSecHubGetConfiguration represnets the response when requesting configuraiton settings.
type ArkSecHubGetConfiguration struct {
	SyncSettings             ArkSecHubSyncSettings             `json:"sync_settings" mapstructure:"sync_settings" desc:"Sync Settings for Secrets Hub" flag:"sync-settings"`
	SecretsSource            ArkSecHubSecretsSource            `json:"secrets_source" mapstructure:"secrets_source" desc:"Secrets Source for Secrets Hub" flag:"secrets-source"`
	AuthenticationIdentities ArkSecHubAuthenticationIdentities `json:"authentication_identities" mapstructure:"authentication_identities" desc:"Authentication Identities for Secrets Hub" flag:"authentication-identities"`
}
