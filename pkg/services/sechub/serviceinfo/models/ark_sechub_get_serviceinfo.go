package models

// ArkSecHubGetServiceInfo represents the response to get service info in Ark Secrets Hub.
type ArkSecHubGetServiceInfo struct {
	TenantRoleArn string `json:"tenant_role_arn" mapstructure:"tenant_role_arn" flag:"tenant-role-arn" desc:"Role ARN of the Secrets Hub Tenant"`
	TenantPamType string `json:"tenant_pam_type" mapstructure:"tenant_pam_type" flag:"tenant-pam-type" desc:"Tenant PAM Type"`
}
