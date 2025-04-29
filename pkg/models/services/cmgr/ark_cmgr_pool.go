package cmgr

// Possible pool types
const (
	Platform = "PLATFORM"
	Access   = "ACCESS"
)

// ArkCmgrPool is a struct representing a pool in the Ark CMGR service.
type ArkCmgrPool struct {
	ID                 string         `json:"id" mapstructure:"id" flag:"id" desc:"ID of the pool"`
	Name               string         `json:"name" mapstructure:"name" flag:"name" desc:"Name of the pool"`
	Description        string         `json:"description,omitempty" mapstructure:"description,omitempty" flag:"description" desc:"Description of the pool"`
	AssignedNetworkIDs []string       `json:"assigned_network_ids" mapstructure:"assigned_network_ids" flag:"assigned-network-ids" desc:"Assigned networks of the pool"`
	IdentifiersCount   int            `json:"identifiers_count,omitempty" mapstructure:"identifiers_count,omitempty" flag:"identifiers-count" desc:"Count of identifiers on the pool"`
	ComponentsCount    map[string]int `json:"components_count,omitempty" mapstructure:"components_count,omitempty" flag:"components-count" desc:"Count of components on the pool"`
	CreatedAt          string         `json:"created_at" mapstructure:"created_at" flag:"created-at" desc:"The creation time of the pool"`
	UpdatedAt          string         `json:"updated_at" mapstructure:"updated_at" flag:"updated-at" desc:"The last update time of the pool"`
}
