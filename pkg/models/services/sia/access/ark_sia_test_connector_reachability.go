package access

// ArkSIATestConnectorReachability represents the schema for testing connector reachability.
type ArkSIATestConnectorReachability struct {
	ConnectorID           string `json:"connector_id" mapstructure:"connector_id" flag:"connector-id" desc:"The id of the connector to test" validate:"required"`
	TargetHostname        string `json:"target_hostname" mapstructure:"target_hostname" flag:"target-hostname" desc:"Target hostname to test the connector against"`
	TargetPort            int    `json:"target_port" mapstructure:"target_port" flag:"target-port" desc:"Target port to test the connector against" default:"22"`
	CheckBackendEndpoints bool   `json:"check_backend_endpoints" mapstructure:"check_backend_endpoints" flag:"check-backend-endpoints" desc:"Whether to check the backend endpoints as well"`
}

// ArkSIATargetElement represents the schema for a target element in the reachability test response.
type ArkSIATargetElement struct {
	TargetIP     string `json:"target_ip" mapstructure:"target_ip"`
	TargetPort   int    `json:"target_port" mapstructure:"target_port"`
	LatencyMlsec int    `json:"latency_mlsec" mapstructure:"latency_mlsec"`
	Status       string `json:"status" mapstructure:"status"`
	Description  string `json:"description" mapstructure:"description"`
}

// ArkSIABackendEndpoint represents the schema for a backend endpoint in the reachability test response.
type ArkSIABackendEndpoint struct {
	BackendConnectorAddress string `json:"backend_connector_endpoint" mapstructure:"backend_connector_endpoint"`
	LatencyMlsec            int    `json:"latency_mlsec" mapstructure:"latency_mlsec"`
	Status                  string `json:"status" mapstructure:"status"`
	Description             string `json:"description" mapstructure:"description"`
}

// ArkSIAReachabilityTestResponse represents the response for the reachability test.
type ArkSIAReachabilityTestResponse struct {
	Targets  []ArkSIATargetElement   `json:"targets" mapstructure:"targets"`
	Backends []ArkSIABackendEndpoint `json:"backends" mapstructure:"backends"`
}
