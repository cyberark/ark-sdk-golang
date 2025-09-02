package models

// ArkSIADeleteConnector represents the request to delete a connector in Ark SIA.
type ArkSIADeleteConnector struct {
	ConnectorID string `json:"connector_id" mapstructure:"connector_id" flag:"connector-id" desc:"The connector ID to delete" validate:"required"`
	RetryCount  int    `json:"retry_count" mapstructure:"retry_count" flag:"retry-count" desc:"Number of times to retry the deletion if it fails" default:"10"`
	RetryDelay  int    `json:"retry_delay" mapstructure:"retry_delay" flag:"retry-delay" desc:"Delay in seconds between retries" default:"5"`
}
