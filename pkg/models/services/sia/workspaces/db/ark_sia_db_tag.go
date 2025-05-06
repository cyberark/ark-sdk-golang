package db

// ArkSIADBTag represents a tag associated with a database in a workspace.
type ArkSIADBTag struct {
	Key   string `json:"key" mapstructure:"key" flag:"key" desc:"Key of the tag, for example environment" validate:"required"`
	Value string `json:"value" mapstructure:"value" flag:"value" desc:"Value of the tag, for example production" validate:"required"`
}

// ArkSIADBTagList represents a list of tags associated with a database in a workspace.
type ArkSIADBTagList struct {
	Tags  []ArkSIADBTag `json:"tags" mapstructure:"tags" flag:"tags" desc:"List of tags"`
	Count int           `json:"count" mapstructure:"count" flag:"count" desc:"The amount of tags listed"`
}
