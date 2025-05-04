package identity

import (
	"encoding/json"
)

// Possible directory types
const (
	AD       = "AdProxy"
	Identity = "CDS"
	FDS      = "FDS"
)

// AllDirectoryTypes is a list of all supported directory types.
var (
	AllDirectoryTypes = []string{
		AD,
		Identity,
		FDS,
	}
)

// DirectoryServiceMetadata represents metadata for a directory service.
type DirectoryServiceMetadata struct {
	Service              string `json:"Service" mapstructure:"Service"`
	DirectoryServiceUUID string `json:"directoryServiceUuid" mapstructure:"directoryServiceUuid"`
}

// DirectoryServiceRow represents a row containing directory service metadata.
type DirectoryServiceRow struct {
	Row DirectoryServiceMetadata `json:"Row" mapstructure:"Row"`
}

// GetDirectoryServicesResult represents the result of a directory services query.
type GetDirectoryServicesResult struct {
	Results []DirectoryServiceRow `json:"Results" mapstructure:"Results" validate:"min=1"`
}

// GetDirectoryServicesResponse represents the response for a directory services query.
type GetDirectoryServicesResponse struct {
	Result GetDirectoryServicesResult `json:"Result" mapstructure:"Result"`
}

// DirectorySearchArgs represents search arguments for directory queries.
type DirectorySearchArgs struct {
	PageNumber int    `json:"PageNumber,omitempty" mapstructure:"PageNumber,omitempty"`
	PageSize   int    `json:"PageSize,omitempty" mapstructure:"PageSize,omitempty"`
	Limit      int    `json:"Limit,omitempty" mapstructure:"Limit,omitempty"`
	SortBy     string `json:"SortBy,omitempty" mapstructure:"SortBy,omitempty"`
	Caching    int    `json:"Caching,omitempty" mapstructure:"Caching,omitempty"`
	Dir        string `json:"Direction,omitempty" mapstructure:"Direction,omitempty"`
	Ascending  bool   `json:"Ascending,omitempty" mapstructure:"Ascending,omitempty"`
}

// DirectoryServiceQueryRequest represents a query request for directory services.
type DirectoryServiceQueryRequest struct {
	DirectoryServices []string            `json:"directoryServices" mapstructure:"directoryServices"`
	Group             string              `json:"group,omitempty" mapstructure:"group,omitempty"`
	Roles             string              `json:"roles,omitempty" mapstructure:"roles,omitempty"`
	User              string              `json:"user,omitempty" mapstructure:"user,omitempty"`
	Args              DirectorySearchArgs `json:"Args" mapstructure:"Args"`
}

// NewDirectoryServiceQueryRequest initializes a DirectoryServiceQueryRequest with optional search string.
func NewDirectoryServiceQueryRequest(searchString string) *DirectoryServiceQueryRequest {
	request := &DirectoryServiceQueryRequest{}
	request.User = "{}"
	request.Roles = "{}"
	request.Group = "{}"
	if searchString != "" {
		groupFilter := map[string]interface{}{
			"_or": []map[string]interface{}{
				{"DisplayName": map[string]string{"_like": searchString}},
				{"SystemName": map[string]string{"_like": searchString}},
			},
		}
		rolesFilter := map[string]interface{}{
			"Name": map[string]interface{}{
				"_like": map[string]interface{}{
					"value":      searchString,
					"ignoreCase": true,
				},
			},
		}
		usersFilter := map[string]interface{}{
			"DisplayName": map[string]string{"_like": searchString},
		}
		grp, _ := json.Marshal(groupFilter)
		roles, _ := json.Marshal(rolesFilter)
		users, _ := json.Marshal(usersFilter)
		request.Group = string(grp)
		request.Roles = string(roles)
		request.User = string(users)
	}
	return request
}

// DirectoryServiceQuerySpecificRoleRequest represents a query request for a specific role.
type DirectoryServiceQuerySpecificRoleRequest struct {
	DirectoryServices []string            `json:"directoryServices" mapstructure:"directoryServices"`
	Group             string              `json:"group,omitempty" mapstructure:"group,omitempty"`
	Roles             string              `json:"roles,omitempty" mapstructure:"roles,omitempty"`
	User              string              `json:"user,omitempty" mapstructure:"user,omitempty"`
	Args              DirectorySearchArgs `json:"Args" mapstructure:"Args"`
}

// NewDirectoryServiceQuerySpecificRoleRequest initializes a DirectoryServiceQuerySpecificRoleRequest with a specific role name.
func NewDirectoryServiceQuerySpecificRoleRequest(roleName string) *DirectoryServiceQuerySpecificRoleRequest {
	request := &DirectoryServiceQuerySpecificRoleRequest{}
	request.User = "{}"
	request.Roles = "{}"
	request.Group = "{}"
	if roleName != "" {
		request.Roles = `{"Name":{"_eq":"` + roleName + `"}}`
	}
	return request
}

// GroupRow represents a row containing group information.
type GroupRow struct {
	DisplayName              string `json:"DisplayName,omitempty" mapstructure:"DisplayName"`
	ServiceInstanceLocalized string `json:"ServiceInstanceLocalized" mapstructure:"ServiceInstanceLocalized"`
	DirectoryServiceType     string `json:"ServiceType" mapstructure:"ServiceType"`
	SystemName               string `json:"SystemName,omitempty" mapstructure:"SystemName"`
	InternalID               string `json:"InternalName,omitempty" mapstructure:"InternalName"`
}

// GroupResult represents a result containing a group row.
type GroupResult struct {
	Row GroupRow `json:"Row" mapstructure:"Row"`
}

// GroupsResult represents the results of a group query.
type GroupsResult struct {
	Results   []GroupResult `json:"Results" mapstructure:"Results"`
	FullCount int           `json:"FullCount,omitempty" mapstructure:"FullCount"`
}

// RoleAdminRight represents administrative rights for a role.
type RoleAdminRight struct {
	Path        string `json:"Path" mapstructure:"Path"`
	ServiceName string `json:"ServiceName,omitempty" mapstructure:"ServiceName"`
}

// RoleRow represents a row containing role information.
type RoleRow struct {
	Name        string           `json:"Name,omitempty" mapstructure:"Name"`
	ID          string           `json:"_ID" mapstructure:"_ID"`
	AdminRights []RoleAdminRight `json:"AdministrativeRights,omitempty" mapstructure:"AdministrativeRights"`
	IsHidden    bool             `json:"IsHidden,omitempty" mapstructure:"IsHidden"`
	Description string           `json:"Description,omitempty" mapstructure:"Description"`
}

// RoleResult represents a result containing a role row.
type RoleResult struct {
	Row RoleRow `json:"Row" mapstructure:"Row"`
}

// RolesResult represents the results of a role query.
type RolesResult struct {
	Results   []RoleResult `json:"Results" mapstructure:"Results"`
	FullCount int          `json:"FullCount,omitempty" mapstructure:"FullCount"`
}

// UserRow represents a row containing user information.
type UserRow struct {
	DisplayName              string `json:"DisplayName,omitempty" mapstructure:"DisplayName"`
	ServiceInstanceLocalized string `json:"ServiceInstanceLocalized" mapstructure:"ServiceInstanceLocalized"`
	DistinguishedName        string `json:"DistinguishedName" mapstructure:"DistinguishedName"`
	SystemName               string `json:"SystemName,omitempty" mapstructure:"SystemName"`
	DirectoryServiceType     string `json:"ServiceType" mapstructure:"ServiceType"`
	Email                    string `json:"EMail,omitempty" mapstructure:"EMail"`
	InternalID               string `json:"InternalName,omitempty" mapstructure:"InternalName"`
	Description              string `json:"Description,omitempty" mapstructure:"Description"`
}

// UserResult represents a result containing a user row.
type UserResult struct {
	Row UserRow `json:"Row" mapstructure:"Row"`
}

// UsersResult represents the results of a user query.
type UsersResult struct {
	Results   []UserResult `json:"Results" mapstructure:"Results"`
	FullCount int          `json:"FullCount,omitempty" mapstructure:"FullCount"`
}

// QueryResult represents the combined results of group, role, and user queries.
type QueryResult struct {
	Groups *GroupsResult `json:"Group,omitempty" mapstructure:"Group"`
	Roles  *RolesResult  `json:"Roles,omitempty" mapstructure:"Roles"`
	Users  *UsersResult  `json:"User,omitempty" mapstructure:"User"`
}

// DirectoryServiceQueryResponse represents the response for a directory service query.
type DirectoryServiceQueryResponse struct {
	Result QueryResult `json:"Result" mapstructure:"Result"`
}
