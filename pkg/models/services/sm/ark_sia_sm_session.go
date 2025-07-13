package sm

// ArkSMSessionStatus represents the status of a session.
type ArkSMSessionStatus string

// Possible session statuses.
const (
	Active ArkSMSessionStatus = "Active"
	Ended  ArkSMSessionStatus = "Ended"
	Failed ArkSMSessionStatus = "Failed"
)

// ArkSMSession represents a session.
type ArkSMSession struct {
	TenantID        string                 `json:"tenant_id,omitempty" mapstructure:"tenant_id,omitempty" desc:"Tenant id of the session"`
	SessionID       string                 `json:"session_id" mapstructure:"session_id" desc:"Session id"`
	SessionStatus   ArkSMSessionStatus     `json:"session_status,omitempty" mapstructure:"session_status,omitempty" desc:"Status of the session"`
	SessionDuration string                 `json:"session_duration,omitempty" mapstructure:"session_duration,omitempty" desc:"Duration of the session in seconds"`
	EndReason       string                 `json:"end_reason,omitempty" mapstructure:"end_reason,omitempty" desc:"End reason for the session"`
	ErrorCode       string                 `json:"error_code,omitempty" mapstructure:"error_code,omitempty" desc:"Error code for the session"`
	ApplicationCode string                 `json:"application_code,omitempty" mapstructure:"application_code,omitempty" desc:"Application code of the session" choices:"SIA,DPA,CSM,PAM,DAP,ITI,UBA,ADM,AUD,ALR,CEM,EPM,SCA,SHSM,CLO,CMS,SMS,PYC,ARS,IDP,ITDR,INTS,MSP,CCE"`
	AccessMethod    string                 `json:"access_method,omitempty" mapstructure:"access_method,omitempty" desc:"Access method of the session" choices:"Vaulted,JIT,Unknown"`
	StartTime       string                 `json:"start_time,omitempty" mapstructure:"start_time,omitempty" desc:"Start time of the session"`
	EndTime         string                 `json:"end_time,omitempty" mapstructure:"end_time,omitempty" desc:"End time of the session"`
	User            string                 `json:"user,omitempty" mapstructure:"user,omitempty" desc:"Username of the session"`
	Source          string                 `json:"source,omitempty" mapstructure:"source,omitempty" desc:"Source of the session (Usually IP)"`
	Target          string                 `json:"target,omitempty" mapstructure:"target,omitempty" desc:"Target of the session (Usually IP/DNS)"`
	TargetUsername  string                 `json:"target_username,omitempty" mapstructure:"target_username,omitempty" desc:"Target username of the session"`
	Protocol        string                 `json:"protocol,omitempty" mapstructure:"protocol,omitempty" desc:"Connection protocol of the session" choices:"SSH,RDP,CLI,CONSOLE,HTTPS,K8S,DB"`
	Platform        string                 `json:"platform,omitempty" mapstructure:"platform,omitempty" desc:"Connection platform of the session"`
	CustomData      map[string]interface{} `json:"custom_data,omitempty" mapstructure:"custom_data,omitempty" desc:"Custom data of the session"`
	IsRecording     bool                   `json:"is_recording,omitempty" mapstructure:"is_recording,omitempty" desc:"Whether the session is recorded or not"`
}

// ArkSMSessions represents a list of sessions with counts.
type ArkSMSessions struct {
	Sessions      []ArkSMSession `json:"sessions" mapstructure:"sessions" desc:"List of the sessions"`
	FilteredCount int            `json:"filtered_count" mapstructure:"filtered_count" desc:"How many sessions were filtered"`
	ReturnedCount int            `json:"returned_count" mapstructure:"returned_count" desc:"How many sessions were returned"`
}
