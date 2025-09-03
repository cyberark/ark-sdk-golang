package models

// ArkSIAGetConnectorSetupScript represents the setup script details for getting a connector.
type ArkSIAGetConnectorSetupScript struct {
	ConnectorType   string `json:"connector_type" mapstructure:"connector_type" flag:"connector-type" desc:"The type of the platform for the connector to be installed in (ON-PREMISE,AWS,AZURE,GCP)" default:"ON-PREMISE" choices:"ON-PREMISE,AWS,AZURE,GCP"`
	ConnectorOS     string `json:"connector_os" mapstructure:"connector_os" flag:"connector-os" desc:"The type of the operating system for the connector to be installed on (linux,windows)" default:"linux" choices:"linux,windows"`
	ConnectorPoolID string `json:"connector_pool_id" mapstructure:"connector_pool_id" flag:"connector-pool-id" desc:"The connector pool which the connector will be part of, if not given, the connector will be assigned to the default one" validate:"required"`
}
