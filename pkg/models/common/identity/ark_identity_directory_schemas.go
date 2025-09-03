// Package identity provides data structures and types for ARK Identity directory services.
// This package contains models for directory service metadata, query requests and responses,
// and data structures for users, groups, and roles within the ARK Identity system.
// It supports various directory types including Active Directory, Identity (CDS), and FDS.
package identity

import (
	"encoding/json"
)

// Directory type constants represent the supported directory service types.
const (
	// AD represents Active Directory Proxy directory type.
	AD = "AdProxy"
	// Identity represents CDS (Centrify Directory Service) directory type.
	Identity = "CDS"
	// FDS represents FDS (Federated Directory Service) directory type.
	FDS = "FDS"
)

// AllDirectoryTypes contains all supported directory service types.
// This slice includes AD, Identity, and FDS directory types for validation
// and enumeration purposes.
var (
	AllDirectoryTypes = []string{
		AD,
		Identity,
		FDS,
	}
)

// DirectoryServiceMetadata represents metadata information for a directory service.
// This structure contains essential identification information for directory services
// including the service type and unique identifier used for directory operations.
type DirectoryServiceMetadata struct {
	Service              string `json:"Service" mapstructure:"Service"`
	DirectoryServiceUUID string `json:"directoryServiceUuid" mapstructure:"directoryServiceUuid"`
}

// DirectoryServiceRow represents a single row of directory service metadata.
// This structure wraps DirectoryServiceMetadata to match the expected JSON
// structure returned by directory service queries.
type DirectoryServiceRow struct {
	Row DirectoryServiceMetadata `json:"Row" mapstructure:"Row"`
}

// GetDirectoryServicesResult represents the collection of directory services results.
// This structure contains an array of directory service rows returned from
// directory service enumeration queries with validation requiring at least one result.
type GetDirectoryServicesResult struct {
	Results []DirectoryServiceRow `json:"Results" mapstructure:"Results" validate:"min=1"`
}

// GetDirectoryServicesResponse represents the complete response for directory services queries.
// This structure wraps the directory services results in the expected API response format.
type GetDirectoryServicesResponse struct {
	Result GetDirectoryServicesResult `json:"Result" mapstructure:"Result"`
}

// DirectorySearchArgs represents search and pagination arguments for directory queries.
// This structure contains parameters for controlling query behavior including
// pagination, sorting, caching, and result ordering options.
type DirectorySearchArgs struct {
	PageNumber int    `json:"PageNumber,omitempty" mapstructure:"PageNumber,omitempty"`
	PageSize   int    `json:"PageSize,omitempty" mapstructure:"PageSize,omitempty"`
	Limit      int    `json:"Limit,omitempty" mapstructure:"Limit,omitempty"`
	SortBy     string `json:"SortBy,omitempty" mapstructure:"SortBy,omitempty"`
	Caching    int    `json:"Caching,omitempty" mapstructure:"Caching,omitempty"`
	Dir        string `json:"Direction,omitempty" mapstructure:"Direction,omitempty"`
	Ascending  bool   `json:"Ascending,omitempty" mapstructure:"Ascending,omitempty"`
}

// DirectoryServiceQueryRequest represents a comprehensive query request for directory services.
// This structure contains filter criteria for groups, roles, and users along with
// directory service specifications and search arguments for controlling query behavior.
type DirectoryServiceQueryRequest struct {
	DirectoryServices []string            `json:"directoryServices" mapstructure:"directoryServices"`
	Group             string              `json:"group,omitempty" mapstructure:"group,omitempty"`
	Roles             string              `json:"roles,omitempty" mapstructure:"roles,omitempty"`
	User              string              `json:"user,omitempty" mapstructure:"user,omitempty"`
	Args              DirectorySearchArgs `json:"Args" mapstructure:"Args"`
}

// NewDirectoryServiceQueryRequest creates a new DirectoryServiceQueryRequest with optional search filtering.
// It initializes the request with default empty JSON objects for user, roles, and group filters.
// If a search string is provided, it creates appropriate filter criteria for searching across
// display names and system names for groups, role names, and user display names.
//
// Parameters:
//   - searchString: Optional search term to filter results across groups, roles, and users
//
// Returns:
//   - *DirectoryServiceQueryRequest: Initialized request with search filters applied if searchString provided
//
// Example:
//
//	// Create request without search filtering
//	request := NewDirectoryServiceQueryRequest("")
//
//	// Create request with search filtering
//	request := NewDirectoryServiceQueryRequest("admin")
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

// DirectoryServiceQuerySpecificRoleRequest represents a query request targeting a specific role.
// This structure is similar to DirectoryServiceQueryRequest but is specialized for querying
// specific roles by exact name match rather than general search filtering.
type DirectoryServiceQuerySpecificRoleRequest struct {
	DirectoryServices []string            `json:"directoryServices" mapstructure:"directoryServices"`
	Group             string              `json:"group,omitempty" mapstructure:"group,omitempty"`
	Roles             string              `json:"roles,omitempty" mapstructure:"roles,omitempty"`
	User              string              `json:"user,omitempty" mapstructure:"user,omitempty"`
	Args              DirectorySearchArgs `json:"Args" mapstructure:"Args"`
}

// NewDirectoryServiceQuerySpecificRoleRequest creates a new DirectoryServiceQuerySpecificRoleRequest for a specific role.
// It initializes the request with default empty JSON objects and sets up an exact match
// filter for the specified role name if provided.
//
// Parameters:
//   - roleName: The exact name of the role to query for
//
// Returns:
//   - *DirectoryServiceQuerySpecificRoleRequest: Initialized request with role name filter applied if roleName provided
//
// Example:
//
//	// Create request for specific role
//	request := NewDirectoryServiceQuerySpecificRoleRequest("System Administrator")
//
//	// Create request without role filtering
//	request := NewDirectoryServiceQuerySpecificRoleRequest("")
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

// GroupRow represents detailed information about a directory group.
// This structure contains group metadata including display names, service information,
// directory service type, system identifiers, and internal references.
type GroupRow struct {
	DisplayName              string `json:"DisplayName,omitempty" mapstructure:"DisplayName"`
	ServiceInstanceLocalized string `json:"ServiceInstanceLocalized" mapstructure:"ServiceInstanceLocalized"`
	DirectoryServiceType     string `json:"ServiceType" mapstructure:"ServiceType"`
	SystemName               string `json:"SystemName,omitempty" mapstructure:"SystemName"`
	InternalID               string `json:"InternalName,omitempty" mapstructure:"InternalName"`
}

// GroupResult represents a single group result from directory queries.
// This structure wraps GroupRow to match the expected JSON structure
// returned by directory service group queries.
type GroupResult struct {
	Row GroupRow `json:"Row" mapstructure:"Row"`
}

// GroupsResult represents the complete collection of group query results.
// This structure contains an array of group results along with the total
// count of matching groups for pagination purposes.
type GroupsResult struct {
	Results   []GroupResult `json:"Results" mapstructure:"Results"`
	FullCount int           `json:"FullCount,omitempty" mapstructure:"FullCount"`
}

// RoleAdminRight represents administrative rights and permissions for a role.
// This structure defines the scope and service context for role-based
// administrative privileges within the directory system.
type RoleAdminRight struct {
	Path        string `json:"Path" mapstructure:"Path"`
	ServiceName string `json:"ServiceName,omitempty" mapstructure:"ServiceName"`
}

// RoleRow represents detailed information about a directory role.
// This structure contains role metadata including name, unique identifier,
// administrative rights, visibility status, and descriptive information.
type RoleRow struct {
	Name        string           `json:"Name,omitempty" mapstructure:"Name"`
	ID          string           `json:"_ID" mapstructure:"_ID"`
	AdminRights []RoleAdminRight `json:"AdministrativeRights,omitempty" mapstructure:"AdministrativeRights"`
	IsHidden    bool             `json:"IsHidden,omitempty" mapstructure:"IsHidden"`
	Description string           `json:"Description,omitempty" mapstructure:"Description"`
}

// RoleResult represents a single role result from directory queries.
// This structure wraps RoleRow to match the expected JSON structure
// returned by directory service role queries.
type RoleResult struct {
	Row RoleRow `json:"Row" mapstructure:"Row"`
}

// RolesResult represents the complete collection of role query results.
// This structure contains an array of role results along with the total
// count of matching roles for pagination purposes.
type RolesResult struct {
	Results   []RoleResult `json:"Results" mapstructure:"Results"`
	FullCount int          `json:"FullCount,omitempty" mapstructure:"FullCount"`
}

// UserRow represents detailed information about a directory user.
// This structure contains comprehensive user metadata including display information,
// service details, distinguished name, system identifiers, contact information,
// and descriptive data.
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

// UserResult represents a single user result from directory queries.
// This structure wraps UserRow to match the expected JSON structure
// returned by directory service user queries.
type UserResult struct {
	Row UserRow `json:"Row" mapstructure:"Row"`
}

// UsersResult represents the complete collection of user query results.
// This structure contains an array of user results along with the total
// count of matching users for pagination purposes.
type UsersResult struct {
	Results   []UserResult `json:"Results" mapstructure:"Results"`
	FullCount int          `json:"FullCount,omitempty" mapstructure:"FullCount"`
}

// QueryResult represents the comprehensive results from directory service queries.
// This structure aggregates results for groups, roles, and users into a single
// response object, allowing for combined query operations across all entity types.
type QueryResult struct {
	Groups *GroupsResult `json:"Group,omitempty" mapstructure:"Group"`
	Roles  *RolesResult  `json:"Roles,omitempty" mapstructure:"Roles"`
	Users  *UsersResult  `json:"User,omitempty" mapstructure:"User"`
}

// DirectoryServiceQueryResponse represents the complete response for directory service queries.
// This structure wraps the query results in the expected API response format for
// directory service operations involving groups, roles, and users.
type DirectoryServiceQueryResponse struct {
	Result QueryResult `json:"Result" mapstructure:"Result"`
}
