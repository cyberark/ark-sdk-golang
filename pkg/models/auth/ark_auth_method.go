package auth

// ArkAuthMethod is a string type that represents the authentication method used in the Ark SDK.
type ArkAuthMethod string

// Authentication methods supported by the Ark SDK.
const (
	Identity            ArkAuthMethod = "identity"
	IdentityServiceUser ArkAuthMethod = "identity_service_user"
	Direct              ArkAuthMethod = "direct"
	Default             ArkAuthMethod = "default"
	Other               ArkAuthMethod = "other"
)

// ArkAuthMethodSettings is an interface that defines the settings for different authentication methods.
type ArkAuthMethodSettings interface{}

// IdentityArkAuthMethodSettings is a struct that represents the settings for the Identity authentication method.
type IdentityArkAuthMethodSettings struct {
	IdentityMFAMethod       string `json:"identity_mfa_method" mapstructure:"identity_mfa_method" validate:"oneof=pf sms email otp" flag:"identity-mfa-method" desc:"MFA Method to use by default [pf, sms, email, otp]"`
	IdentityMFAInteractive  bool   `json:"identity_mfa_interactive" mapstructure:"identity_mfa_interactive" validate:"required" flag:"identity-mfa-interactive" desc:"Allow Interactive MFA"`
	IdentityURL             string `json:"identity_url" mapstructure:"identity_url" flag:"identity-url" desc:"Identity Url"`
	IdentityTenantSubdomain string `json:"identity_tenant_subdomain" mapstructure:"identity_tenant_subdomain" flag:"identity-tenant-subdomain" desc:"Identity Tenant Subdomain"`
}

// IdentityServiceUserArkAuthMethodSettings is a struct that represents the settings for the Identity Service User authentication method.
type IdentityServiceUserArkAuthMethodSettings struct {
	IdentityURL                      string `json:"identity_url" mapstructure:"identity_url" flag:"identity-url" desc:"Identity Url"`
	IdentityTenantSubdomain          string `json:"identity_tenant_subdomain" mapstructure:"identity_tenant_subdomain" flag:"identity-tenant-subdomain" desc:"Identity Tenant Subdomain"`
	IdentityAuthorizationApplication string `json:"identity_authorization_application" mapstructure:"identity_authorization_application" validate:"required" flag:"identity-authorization-application" desc:"Identity Authorization Application" default:"__idaptive_cybr_user_oidc"`
}

// DirectArkAuthMethodSettings is a struct that represents the settings for the Direct authentication method.
type DirectArkAuthMethodSettings struct {
	Endpoint    string `json:"endpoint" mapstructure:"endpoint" flag:"endpoint" desc:"Authentication Endpoint"`
	Interactive bool   `json:"interactive" mapstructure:"interactive" flag:"interactive" desc:"Allow interactiveness"`
}

// DefaultArkAuthMethodSettings is a struct that represents the default settings for the authentication method.
type DefaultArkAuthMethodSettings struct{}

// ArkAuthMethodSettingsMap is a map that associates each ArkAuthMethod with its corresponding settings struct.
var ArkAuthMethodSettingsMap = map[ArkAuthMethod]interface{}{
	Identity:            &IdentityArkAuthMethodSettings{},
	IdentityServiceUser: &IdentityServiceUserArkAuthMethodSettings{},
	Direct:              &DirectArkAuthMethodSettings{},
	Default:             &DefaultArkAuthMethodSettings{},
}

// ArkAuthMethodsDescriptionMap is a map that provides descriptions for each ArkAuthMethod.
var ArkAuthMethodsDescriptionMap = map[ArkAuthMethod]string{
	Identity:            "Identity Personal User",
	IdentityServiceUser: "Identity Service User",
	Direct:              "Direct Endpoint Access",
	Default:             "Default Authenticator Method",
}

// ArkAuthMethodsRequireCredentials is a slice of ArkAuthMethod that require credentials.
var ArkAuthMethodsRequireCredentials = []ArkAuthMethod{
	Identity, IdentityServiceUser, Direct,
}

// ArkAuthMethodSharableCredentials is a slice of ArkAuthMethod that can share credentials.
var ArkAuthMethodSharableCredentials = []ArkAuthMethod{
	Identity,
}
