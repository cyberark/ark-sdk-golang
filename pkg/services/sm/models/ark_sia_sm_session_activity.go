package models

// ArkSMSessionActivity represents a single session audit activity.
type ArkSMSessionActivity struct {
	UUID            string `json:"uuid" mapstructure:"uuid" desc:"ID of the audit"`
	TenantID        string `json:"tenant_id" mapstructure:"tenant_id" desc:"Tenant id of the audit"`
	Timestamp       string `json:"timestamp" mapstructure:"timestamp" desc:"Time of the audit"`
	Username        string `json:"username" mapstructure:"username" desc:"Username of the audit"`
	ApplicationCode string `json:"application_code" mapstructure:"application_code" desc:"Application code of the audit" choices:"SIA,DPA,CSM,PAM,DAP,ITI,UBA,ADM,AUD,ALR,CEM,EPM,SCA,SHSM,CLO,CMS,SMS,PYC,ARS,IDP,ITDR,INTS,MSP,CCE"`
	Action          string `json:"action" mapstructure:"action" desc:"Action performed for the audit"`
	UserID          string `json:"user_id" mapstructure:"user_id" desc:"Id of the user who performed the audit"`
	Source          string `json:"source" mapstructure:"source" desc:"Source of the audit"`
	ActionType      string `json:"action_type" mapstructure:"action_type" desc:"Type of action for the audit"`
	AuditCode       string `json:"audit_code,omitempty" mapstructure:"audit_code,omitempty" desc:"Audit code of the audit"`
	Command         string `json:"command,omitempty" mapstructure:"command,omitempty" desc:"Command performed as part of the audit"`
	Target          string `json:"target,omitempty" mapstructure:"target,omitempty" desc:"Target of the audit"`
	ServiceName     string `json:"service_name,omitempty" mapstructure:"service_name,omitempty" desc:"Service name of the audit"`
	SessionID       string `json:"session_id,omitempty" mapstructure:"session_id,omitempty" desc:"Session id of the audit if related to a session"`
	Message         string `json:"message,omitempty" mapstructure:"message,omitempty" desc:"Message of the audit"`
}

// ArkSMSessionActivities holds a list of session activities and related counts.
type ArkSMSessionActivities struct {
	Activities    []ArkSMSessionActivity `json:"activities" mapstructure:"activities" desc:"List of the session activities"`
	FilteredCount int                    `json:"filtered_count" mapstructure:"filtered_count" desc:"How many session activities were filtered"`
	ReturnedCount int                    `json:"returned_count" mapstructure:"returned_count" desc:"How many session activities were returned"`
}
