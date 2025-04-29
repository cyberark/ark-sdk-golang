package access

type ArkSiaTestConnectorReachability struct {
	ConnectorID           string `json:"connector_id" mapstructure:"connector_id" flag:"connector-id" desc:"The id of the connector to test" validate:"required"`
	TargetHostname        string `json:"target_hostname" mapstructure:"target_hostname" flag:"target-hostname" desc:"Target hostname to test the connector against"`
	TargetPort            int    `json:"target_port" mapstructure:"target_port" flag:"target-port" desc:"Target port to test the connector against" default:"22"`
	CheckBackendEndpoints bool   `json:"check_backend_endpoints" mapstructure:"check_backend_endpoints" flag:"check-backend-endpoints" desc:"Whether to check the backend endpoints as well"`
}
type ArkSiaTargetElement struct {
	TargetIP     string `json:"target_ip" mapstructure:"target_ip"`
	TargetPort   int    `json:"target_port" mapstructure:"target_port"`
	LatencyMlsec int    `json:"latency_mlsec" mapstructure:"latency_mlsec"`
	Status       string `json:"status" mapstructure:"status"`
	Description  string `json:"description" mapstructure:"description"`
}

type ArkSiaBackendEndpoint struct {
	BackendConnectorAddress string `json:"backend_connector_endpoint" mapstructure:"backend_connector_endpoint"`
	LatencyMlsec            int    `json:"latency_mlsec" mapstructure:"latency_mlsec"`
	Status                  string `json:"status" mapstructure:"status"`
	Description             string `json:"description" mapstructure:"description"`
}

type ArkSiaReachabilityTestResponse struct {
	Targets  []ArkSiaTargetElement   `json:"targets" mapstructure:"targets"`
	Backends []ArkSiaBackendEndpoint `json:"backends" mapstructure:"backends"`
}
