package models

// Possible types of principals in UAP.
const (
	PrincipalTypeUser  = "USER"
	PrincipalTypeRole  = "ROLE"
	PrincipalTypeGroup = "GROUP"
)

// ArkUAPPrincipal represents a principal in UAP.
type ArkUAPPrincipal struct {
	ID                  string `json:"id" mapstructure:"id" flag:"id" desc:"The id of the principal" validate:"max=40"`
	Name                string `json:"name" mapstructure:"name" flag:"name" desc:"The name of the principal" validate:"max=100,regexp=^\\w+$"`
	SourceDirectoryName string `json:"source_directory_name,omitempty" mapstructure:"source_directory_name,omitempty" flag:"source-directory-name" desc:"The name of the source directory" validate:"max=50,regexp=^\\w+$"`
	SourceDirectoryID   string `json:"source_directory_id,omitempty" mapstructure:"source_directory_id,omitempty" flag:"source-directory-id" desc:"The id of the source directory"`
	Type                string `json:"type" mapstructure:"type" flag:"type" desc:"The type of the principal user, group or role" choices:"USER,ROLE,GROUP"`
}
