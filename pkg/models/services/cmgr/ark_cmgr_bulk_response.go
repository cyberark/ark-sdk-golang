package cmgr

// ArkCmgrBulkResponse is a struct representing the response of a bulk request in the Ark CMGR service.
type ArkCmgrBulkResponse struct {
	Body       map[string]interface{} `json:"body,omitempty" mapstructure:"body,omitempty" flag:"body" desc:"Response body of the request"`
	StatusCode int                    `json:"status_code" mapstructure:"status_code" flag:"status-code" desc:"Status code of the response"`
}

// ArkCmgrBulkResponses is a struct representing the responses of a bulk request in the Ark CMGR service.
type ArkCmgrBulkResponses struct {
	Responses map[string]ArkCmgrBulkResponse `json:"responses" mapstructure:"responses" flag:"responses" desc:"Responses of the bulk request"`
}
