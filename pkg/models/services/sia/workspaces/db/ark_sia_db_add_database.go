package db

// ArkSIADBAddDatabase represents the request to add a database.
type ArkSIADBAddDatabase struct {
	Name                                 string            `json:"name" mapstructure:"name" flag:"name" desc:"Name of the database, often referenced in policies and other APIs" validate:"required"`
	NetworkName                          string            `json:"network_name" mapstructure:"network_name" flag:"network-name" desc:"Name of the network the database resides in, defaulted to on premises" default:"ON-PREMISE"`
	Platform                             string            `json:"platform" mapstructure:"platform" flag:"platform" desc:"Platform of the database, as in, where it resides, defaulted to on premises" default:"ON-PREMISE" choices:"AWS,AZURE,GCP,ON-PREMISE,ATLAS"`
	AuthDatabase                         string            `json:"auth_database" mapstructure:"auth_database" flag:"auth-database" desc:"Authentication database used, most commonly used with mongodb" default:"admin"`
	Services                             []string          `json:"services" mapstructure:"services" flag:"services" desc:"Services related to the database, most commonly used with oracle/sql-server"`
	Domain                               string            `json:"domain,omitempty" mapstructure:"domain,omitempty" flag:"domain" desc:"The domain the DB resides in"`
	DomainControllerName                 string            `json:"domain_controller_name,omitempty" mapstructure:"domain_controller_name,omitempty" flag:"domain-controller-name" desc:"Domain controller name associated to this database"`
	DomainControllerNetbios              string            `json:"domain_controller_netbios,omitempty" mapstructure:"domain_controller_netbios,omitempty" flag:"domain-controller-netbios" desc:"Domain controller netbios associated to this database"`
	DomainControllerUseLDAPS             bool              `json:"domain_controller_use_ldaps" mapstructure:"domain_controller_use_ldaps" flag:"domain-controller-use-ldaps" desc:"Whether to work with LDAP secure or not" default:"false"`
	DomainControllerEnableCertValidation bool              `json:"domain_controller_enable_certificate_validation" mapstructure:"domain_controller_enable_certificate_validation" flag:"domain-controller-enable-certificate-validation" desc:"Whether to enforce certificate validation on TLS comm to the DC" default:"true"`
	DomainControllerLDAPSCertificate     string            `json:"domain_controller_ldaps_certificate,omitempty" mapstructure:"domain_controller_ldaps_certificate,omitempty" flag:"domain-controller-ldaps-certificate" desc:"Certificate id to use for the domain controller TLS comm"`
	Account                              string            `json:"account,omitempty" mapstructure:"account,omitempty" flag:"account" desc:"Account to be used for provider based databases such as snowflake or atlas"`
	ProviderEngine                       string            `json:"provider_engine" mapstructure:"provider_engine" flag:"provider-engine" desc:"Provider engine, will be later deduced to the identifier of the provider"`
	EnableCertificateValidation          bool              `json:"enable_certificate_validation" mapstructure:"enable_certificate_validation" flag:"enable-certificate-validation" desc:"Whether to enable and enforce certificate validation" default:"true"`
	Certificate                          string            `json:"certificate,omitempty" mapstructure:"certificate,omitempty" flag:"certificate" desc:"Certificate id used for this database that resides in the certificates service"`
	ReadWriteEndpoint                    string            `json:"read_write_endpoint,omitempty" mapstructure:"read_write_endpoint,omitempty" flag:"read-write-endpoint" desc:"Read write endpoint of the database"`
	ReadOnlyEndpoint                     string            `json:"read_only_endpoint,omitempty" mapstructure:"read_only_endpoint,omitempty" flag:"read-only-endpoint" desc:"Optionally, a read only endpoint of the database"`
	Port                                 int               `json:"port,omitempty" mapstructure:"port,omitempty" flag:"port" desc:"Port of the database, if not given, the default one will be used" validate:"omitempty,min=1,max=65535"`
	SecretID                             string            `json:"secret_id,omitempty" mapstructure:"secret_id,omitempty" flag:"secret-id" desc:"Secret identifier stored in the secret service related to this database"`
	Tags                                 map[string]string `json:"tags,omitempty" mapstructure:"tags,omitempty" flag:"tags" desc:"Tags for the database"`
	ConfiguredAuthMethodType             string            `json:"configured_auth_method_type,omitempty" mapstructure:"configured_auth_method_type,omitempty" flag:"configured-auth-method-type" desc:"The target configured auth method type" choices:"ad_ephemeral_user,local_ephemeral_user,rds_iam_authentication,atlas_ephemeral_user"`
	Region                               string            `json:"region,omitempty" mapstructure:"region,omitempty" flag:"region" desc:"Region of the database, most commonly used with IAM authentication"`
}
