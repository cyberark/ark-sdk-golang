package models

// ArkSIAAccessConnectorID is a struct that represents the connector ID for Ark SIA Access.
type ArkSIAAccessConnectorID struct {
	ConnectorID string `json:"connector_id" mapstructure:"connector_id" flag:"connector-id" desc:"The connector ID" validate:"required"`
}
