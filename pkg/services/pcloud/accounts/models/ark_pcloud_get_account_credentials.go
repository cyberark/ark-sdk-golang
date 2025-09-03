package models

// Possible values for ActionType
const (
	Show    = "show"
	Copy    = "copy"
	Connect = "connect"
)

// ArkPCloudGetAccountCredentials represents the details required to retrieve account credentials.
type ArkPCloudGetAccountCredentials struct {
	AccountID           string `json:"account_id" mapstructure:"account_id" desc:"The id of the account to retrieve the credentials for" flag:"account-id" validate:"required"`
	Reason              string `json:"reason,omitempty" mapstructure:"reason,omitempty" desc:"Reason for retrieving the credentials" flag:"reason"`
	TicketingSystemName string `json:"ticketing_system_name,omitempty" mapstructure:"ticketing_system_name,omitempty" desc:"Ticketing system name to use for retrieval of the credentials" flag:"ticketing-system-name"`
	TicketID            string `json:"ticket_id,omitempty" mapstructure:"ticket_id,omitempty" desc:"Ticket id allowing retrieval of the credentials" flag:"ticket-id"`
	Version             string `json:"version,omitempty" mapstructure:"version,omitempty" desc:"Version of the credentials to retrieve" flag:"version"`
	ActionType          string `json:"action_type" mapstructure:"action_type" desc:"Action type of the retrieval (show,copy,connect)" flag:"action-type" default:"show" choices:"show,copy,connect"`
	Machine             string `json:"machine,omitempty" mapstructure:"machine,omitempty" desc:"The address of the remote machine to connect to with the credentials" flag:"machine"`
}
