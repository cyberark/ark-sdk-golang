package cmgr

// ArkCmgrPoolsCommonFilter is a struct representing the common filter for pools in the Ark CMGR service.
type ArkCmgrPoolsCommonFilter struct {
	Projection string `json:"projection" mapstructure:"projection" flag:"projection" desc:"Type of projection for the response" default:"BASIC" choices:"BASIC,EXTENDED"`
	Sort       string `json:"sort,omitempty" mapstructure:"sort,omitempty" flag:"sort" desc:"Sort by given parameter"`
	Filter     string `json:"filter,omitempty" mapstructure:"filter,omitempty" flag:"filter" desc:"Filter parameters"`
	Order      string `json:"order,omitempty" mapstructure:"order,omitempty" flag:"order" desc:"Response sort order" choices:"ASC,DESC" default:"ASC"`
	PageSize   int    `json:"page_size,omitempty" mapstructure:"page_size,omitempty" flag:"page-size" desc:"Size of page"`
}
