package models

import (
	"errors"
)

// ArkSCAAzureWorkspaceType represents the workspace types in Azure.
const (
	AzureWSTypeDirectory       = "directory"
	AzureWSTypeSubscription    = "subscription"
	AzureWSTypeResourceGroup   = "resource_group"
	AzureWSTypeResource        = "resource"
	AzureWSTypeManagementGroup = "management_group"
)

// ArkSCAGCPWorkspaceType represents the workspace types in GCP.
const (
	GCPWSTypeOrganization = "gcp_organization"
	GCPWSTypeFolder       = "folder"
	GCPWSTypeProject      = "project"
)

// ArkSCAGcpRoleType represents the role types in GCP.
const (
	GCPRoleTypePreDefined = iota
	GCPRoleTypeCustom
	GCPRoleTypeBasic
)

// ArkUAPSCABaseTarget defines the interface for serializing and deserializing target structures.
type ArkUAPSCABaseTarget interface {
	Serialize() (map[string]interface{}, error)
	Deserialize(data map[string]interface{}) error
}

// ArkUAPSCATarget represents the base target structure.
type ArkUAPSCATarget struct {
	RoleID        string `json:"role_id" mapstructure:"role_id" flag:"role-id" desc:"The role id of the target"`
	WorkspaceID   string `json:"workspace_id" mapstructure:"workspace_id" flag:"workspace-id" desc:"The workspace id of the target"`
	RoleName      string `json:"role_name,omitempty" mapstructure:"role_name,omitempty" flag:"role-name" desc:"The role name of the target"`
	WorkspaceName string `json:"workspace_name,omitempty" mapstructure:"workspace_name,omitempty" flag:"workspace-name" desc:"The workspace name of the target"`
}

// ArkUAPSCAOrgTarget represents a target with an organization ID.
type ArkUAPSCAOrgTarget struct {
	ArkUAPSCATarget `mapstructure:",squash"`
	OrgID           string `json:"org_id" mapstructure:"org_id" flag:"org-id" desc:"The organization id of the cloud target"`
}

// ArkUAPSCAAWSAccountTarget represents an AWS account target.
type ArkUAPSCAAWSAccountTarget struct {
	ArkUAPSCATarget `mapstructure:",squash"`
}

// Serialize serializes the ArkUAPSCAAWSAccountTarget to a map.
func (s *ArkUAPSCAAWSAccountTarget) Serialize() (map[string]interface{}, error) {
	return map[string]interface{}{
		"roleId":        s.RoleID,
		"workspaceId":   s.WorkspaceID,
		"roleName":      s.RoleName,
		"workspaceName": s.WorkspaceName,
	}, nil
}

// Deserialize deserializes the map into the ArkUAPSCAAWSAccountTarget.
func (s *ArkUAPSCAAWSAccountTarget) Deserialize(data map[string]interface{}) error {
	if roleID, ok := data["role_id"].(string); ok {
		s.RoleID = roleID
	}
	if workspaceID, ok := data["workspace_id"].(string); ok {
		s.WorkspaceID = workspaceID
	}
	if roleName, ok := data["role_name"].(string); ok {
		s.RoleName = roleName
	}
	if workspaceName, ok := data["workspace_name"].(string); ok {
		s.WorkspaceName = workspaceName
	}
	return nil
}

// ArkUAPSCAAWSOrganizationTarget represents an AWS organization target.
type ArkUAPSCAAWSOrganizationTarget struct {
	ArkUAPSCAOrgTarget `mapstructure:",squash"`
}

// Serialize serializes the ArkUAPSCAAWSOrganizationTarget to a map.
func (s *ArkUAPSCAAWSOrganizationTarget) Serialize() (map[string]interface{}, error) {
	return map[string]interface{}{
		"roleId":        s.RoleID,
		"workspaceId":   s.WorkspaceID,
		"roleName":      s.RoleName,
		"workspaceName": s.WorkspaceName,
		"orgId":         s.OrgID,
	}, nil
}

// Deserialize deserializes the map into the ArkUAPSCAAWSOrganizationTarget.
func (s *ArkUAPSCAAWSOrganizationTarget) Deserialize(data map[string]interface{}) error {
	if roleID, ok := data["role_id"].(string); ok {
		s.RoleID = roleID
	}
	if workspaceID, ok := data["workspace_id"].(string); ok {
		s.WorkspaceID = workspaceID
	}
	if roleName, ok := data["role_name"].(string); ok {
		s.RoleName = roleName
	}
	if workspaceName, ok := data["workspace_name"].(string); ok {
		s.WorkspaceName = workspaceName
	}
	if orgID, ok := data["org_id"].(string); ok {
		s.OrgID = orgID
	}
	return nil
}

// ArkUAPSCAAzureTarget represents an Azure target.
type ArkUAPSCAAzureTarget struct {
	ArkUAPSCAOrgTarget `mapstructure:",squash"`
	WorkspaceType      string `json:"workspace_type" mapstructure:"workspace_type" flag:"workspace-type" desc:"The type of the workspace in Azure" choices:"directory,subscription,resource_group,resource,management_group"`
	RoleType           int    `json:"role_type,omitempty" mapstructure:"role_type,omitempty" flag:"role-type" desc:"The type of the role in Azure"`
}

// Serialize serializes the ArkUAPSCAAzureTarget to a map.
func (s *ArkUAPSCAAzureTarget) Serialize() (map[string]interface{}, error) {
	return map[string]interface{}{
		"roleId":        s.RoleID,
		"workspaceId":   s.WorkspaceID,
		"roleName":      s.RoleName,
		"workspaceName": s.WorkspaceName,
		"orgId":         s.OrgID,
		"workspaceType": s.WorkspaceType,
		"roleType":      s.RoleType,
	}, nil
}

// Deserialize deserializes the map into the ArkUAPSCAAzureTarget.
func (s *ArkUAPSCAAzureTarget) Deserialize(data map[string]interface{}) error {
	if roleID, ok := data["role_id"].(string); ok {
		s.RoleID = roleID
	}
	if workspaceID, ok := data["workspace_id"].(string); ok {
		s.WorkspaceID = workspaceID
	}
	if roleName, ok := data["role_name"].(string); ok {
		s.RoleName = roleName
	}
	if workspaceName, ok := data["workspace_name"].(string); ok {
		s.WorkspaceName = workspaceName
	}
	if orgID, ok := data["org_id"].(string); ok {
		s.OrgID = orgID
	}
	if workspaceType, ok := data["workspace_type"].(string); ok {
		s.WorkspaceType = workspaceType
	}
	if roleType, ok := data["role_type"].(int); ok {
		s.RoleType = roleType
	}
	return nil
}

// ArkUAPSCAGCPTarget represents a GCP target.
type ArkUAPSCAGCPTarget struct {
	ArkUAPSCAOrgTarget `mapstructure:",squash"`
	WorkspaceType      string `json:"workspace_type" mapstructure:"workspace_type" flag:"workspace-type" desc:"The type of the workspace in GCP" choices:"gcp_organization,folder,project"`
	RolePackage        string `json:"role_package,omitempty" mapstructure:"role_package,omitempty" flag:"role-package" desc:"The role package of the target"`
	RoleType           int    `json:"role_type,omitempty" mapstructure:"role_type,omitempty" flag:"role-type" desc:"The type of the role in GCP"`
}

// Serialize serializes the ArkUAPSCAGCPTarget to a map.
func (s *ArkUAPSCAGCPTarget) Serialize() (map[string]interface{}, error) {
	return map[string]interface{}{
		"roleId":        s.RoleID,
		"workspaceId":   s.WorkspaceID,
		"roleName":      s.RoleName,
		"workspaceName": s.WorkspaceName,
		"orgId":         s.OrgID,
		"workspaceType": s.WorkspaceType,
		"rolePackage":   s.RolePackage,
		"roleType":      s.RoleType,
	}, nil
}

// Deserialize deserializes the map into the ArkUAPSCAGCPTarget.
func (s *ArkUAPSCAGCPTarget) Deserialize(data map[string]interface{}) error {
	if roleID, ok := data["role_id"].(string); ok {
		s.RoleID = roleID
	}
	if workspaceID, ok := data["workspace_id"].(string); ok {
		s.WorkspaceID = workspaceID
	}
	if roleName, ok := data["role_name"].(string); ok {
		s.RoleName = roleName
	}
	if workspaceName, ok := data["workspace_name"].(string); ok {
		s.WorkspaceName = workspaceName
	}
	if orgID, ok := data["org_id"].(string); ok {
		s.OrgID = orgID
	}
	if workspaceType, ok := data["workspace_type"].(string); ok {
		s.WorkspaceType = workspaceType
	}
	if rolePackage, ok := data["role_package"].(string); ok {
		s.RolePackage = rolePackage
	}
	if roleType, ok := data["role_type"].(int); ok {
		s.RoleType = roleType
	}
	return nil
}

// ArkUAPSCACloudConsoleTarget represents a cloud console target.
type ArkUAPSCACloudConsoleTarget struct {
	AwsAccountTargets      []ArkUAPSCAAWSAccountTarget      `json:"aws_account_targets,omitempty" mapstructure:"aws_account_targets,omitempty" flag:"aws-account-targets" desc:"List of AWS account targets"`
	AwsOrganizationTargets []ArkUAPSCAAWSOrganizationTarget `json:"aws_organization_targets,omitempty" mapstructure:"aws_organization_targets,omitempty" flag:"aws-organization-targets" desc:"List of AWS organization targets"`
	AzureTargets           []ArkUAPSCAAzureTarget           `json:"azure_targets,omitempty" mapstructure:"azure_targets,omitempty" flag:"azure-targets" desc:"List of Azure targets"`
	GcpTargets             []ArkUAPSCAGCPTarget             `json:"gcp_targets,omitempty" mapstructure:"gcp_targets,omitempty" flag:"gcp-targets" desc:"List of GCP targets"`
}

// SerializeTargets serializes the ArkUAPSCACloudConsoleTarget to a map.
func (s *ArkUAPSCACloudConsoleTarget) SerializeTargets() (map[string]interface{}, error) {
	targets := make([]interface{}, 0)

	for _, target := range s.AwsAccountTargets {
		data, err := target.Serialize()
		if err != nil {
			return nil, err
		}
		targets = append(targets, data)
	}

	for _, target := range s.AwsOrganizationTargets {
		data, err := target.Serialize()
		if err != nil {
			return nil, err
		}
		targets = append(targets, data)
	}

	for _, target := range s.AzureTargets {
		data, err := target.Serialize()
		if err != nil {
			return nil, err
		}
		targets = append(targets, data)
	}

	for _, target := range s.GcpTargets {
		data, err := target.Serialize()
		if err != nil {
			return nil, err
		}
		targets = append(targets, data)
	}

	return map[string]interface{}{
		"targets": targets,
	}, nil
}

// DeserializeTargets deserializes the map into the ArkUAPSCACloudConsoleTarget.
func (s *ArkUAPSCACloudConsoleTarget) DeserializeTargets(data map[string]interface{}) error {
	targetsData, ok := data["targets"].([]interface{})
	if !ok {
		return errors.New("invalid targets data format")
	}
	for _, targetData := range targetsData {
		targetMap, ok := targetData.(map[string]interface{})
		if !ok {
			return errors.New("invalid target data format")
		}
		if workspaceType, ok := targetMap["workspace_type"].(string); ok {
			switch workspaceType {
			case AzureWSTypeDirectory, AzureWSTypeSubscription, AzureWSTypeResourceGroup, AzureWSTypeResource, AzureWSTypeManagementGroup:
				var target ArkUAPSCAAzureTarget
				if err := target.Deserialize(targetMap); err != nil {
					return err
				}
				s.AzureTargets = append(s.AzureTargets, target)
			case GCPWSTypeOrganization, GCPWSTypeFolder, GCPWSTypeProject:
				var target ArkUAPSCAGCPTarget
				if err := target.Deserialize(targetMap); err != nil {
					return err
				}
				s.GcpTargets = append(s.GcpTargets, target)
			default:
				return errors.New("unknown workspace type in cloud console targets")
			}
		} else {
			if _, ok := targetMap["org_id"]; ok {
				var target ArkUAPSCAAWSOrganizationTarget
				if err := target.Deserialize(targetMap); err != nil {
					return err
				}
				s.AwsOrganizationTargets = append(s.AwsOrganizationTargets, target)
			} else if _, ok := targetMap["workspace_id"]; ok {
				var target ArkUAPSCAAWSAccountTarget
				if err := target.Deserialize(targetMap); err != nil {
					return err
				}
				s.AwsAccountTargets = append(s.AwsAccountTargets, target)
			} else {
				return errors.New("unknown target type in cloud console targets")
			}
		}
	}
	return nil
}

// ClearTargetsFromData clears the target data from the provided map.
func (s *ArkUAPSCACloudConsoleTarget) ClearTargetsFromData(data map[string]interface{}) {
	if _, ok := data["aws_account_targets"]; ok {
		delete(data, "aws_account_targets")
	}
	if _, ok := data["awsAccountTargets"]; ok {
		delete(data, "awsAccountTargets")
	}
	if _, ok := data["aws_organization_targets"]; ok {
		delete(data, "aws_organization_targets")
	}
	if _, ok := data["awsOrganizationTargets"]; ok {
		delete(data, "awsOrganizationTargets")
	}
	if _, ok := data["azure_targets"]; ok {
		delete(data, "azure_targets")
	}
	if _, ok := data["azureTargets"]; ok {
		delete(data, "azureTargets")
	}
	if _, ok := data["gcp_targets"]; ok {
		delete(data, "gcp_targets")
	}
	if _, ok := data["gcpTargets"]; ok {
		delete(data, "gcpTargets")
	}
}
