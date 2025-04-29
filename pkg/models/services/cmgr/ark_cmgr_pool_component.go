package cmgr

// Possible values for the Type field in ArkCmgrPoolComponent
const (
	PlatformConnector = "PLATFORM_CONNECTOR"
	AccessConnector   = "ACCESS_CONNECTOR"
)

// ArkCmgrPoolComponent is a struct representing a component in the Ark CMGR service.
type ArkCmgrPoolComponent struct {
	ID         string `json:"id" mapstructure:"id" flag:"id" desc:"ID of the component"`
	Type       string `json:"type" mapstructure:"type" flag:"type" desc:"Type of the component" choices:"PLATFORM_CONNECTOR,ACCESS_CONNECTOR"`
	ExternalID string `json:"external_id" mapstructure:"external_id" flag:"external-id" desc:"External identifier of the component"`
	PoolID     string `json:"pool_id,omitempty" mapstructure:"pool_id,omitempty" flag:"pool-id" desc:"Pool id of the pool holding the component"`
	PoolName   string `json:"pool_name,omitempty" mapstructure:"pool_name,omitempty" flag:"pool-name" desc:"Name of the pool holding the component"`
	CreatedAt  string `json:"created_at,omitempty" mapstructure:"created_at,omitempty" flag:"created-at" desc:"The creation time of the component"`
	UpdatedAt  string `json:"updated_at,omitempty" mapstructure:"updated_at,omitempty" flag:"updated-at" desc:"The last update time of the component"`
}
