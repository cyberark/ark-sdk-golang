package vm

import (
	"errors"

	"github.com/cyberark/ark-sdk-golang/pkg/models/common"
)

// Possible operators for FQDN and IP rules in VM access policies.
const (
	VMFQDNOperatorExactly  = "EXACTLY"
	VMFQDNOperatorWildcard = "WILDCARD"
	VMFQDNOperatorPrefix   = "PREFIX"
	VMFQDNOperatorSuffix   = "SUFFIX"
	VMFQDNOperatorContains = "CONTAINS"
)

// Possible operators for IP rules in VM access policies.
const (
	VMIPOperatorExactly  = "EXACTLY"
	VMIPOperatorWildcard = "WILDCARD"
)

// ArkUAPSIAVMKeyValTag defines a key/value pair used to match a given tag or label on a VM resource.
type ArkUAPSIAVMKeyValTag struct {
	Key   string   `json:"key" mapstructure:"key" flag:"key" validate:"min=1"`
	Value []string `json:"value,omitempty" mapstructure:"value" flag:"value"`
}

// Serialize converts the key-value tag to a map.
func (t *ArkUAPSIAVMKeyValTag) Serialize() map[string]interface{} {
	return map[string]interface{}{
		"key":   t.Key,
		"value": t.Value,
	}
}

// Deserialize populates the key-value tag from a map.
func (t *ArkUAPSIAVMKeyValTag) Deserialize(data map[string]interface{}) error {
	if key, ok := data["key"].(string); ok {
		t.Key = key
	} else {
		return errors.New("key must be a string")
	}
	if value, ok := data["value"].([]interface{}); ok {
		for _, v := range value {
			if str, ok := v.(string); ok {
				t.Value = append(t.Value, str)
			} else {
				return errors.New("value must be a list of strings")
			}
		}
	}
	return nil
}

// ArkUAPSIAVMFQDNRule defines a specific FQDN rule used to match a given DNS record.
type ArkUAPSIAVMFQDNRule struct {
	Operator            string `json:"operator" mapstructure:"operator" flag:"operator" choices:"EXACTLY,WILDCARD,PREFIX,SUFFIX,CONTAINS" desc:"The operator to use for matching the FQDN. Valid values are EXACTLY, WILDCARD, PREFIX, SUFFIX, and CONTAINS."`
	ComputernamePattern string `json:"computername_pattern" mapstructure:"computername_pattern" flag:"computername-pattern" validate:"max=300" desc:"The pattern to match against the computer name. This can be a full FQDN or a partial match."`
	Domain              string `json:"domain,omitempty" mapstructure:"domain" flag:"domain" validate:"max=1000" desc:"The domain to match against the FQDN. This is optional and can be used to further restrict the match."`
}

// Serialize converts the FQDN rule to a map.
func (r *ArkUAPSIAVMFQDNRule) Serialize() map[string]interface{} {
	return map[string]interface{}{
		"operator":            r.Operator,
		"computernamePattern": r.ComputernamePattern,
		"domain":              r.Domain,
	}
}

// Deserialize populates the FQDN rule from a map.
func (r *ArkUAPSIAVMFQDNRule) Deserialize(data map[string]interface{}) error {
	if operator, ok := data["operator"].(string); ok {
		r.Operator = operator
	} else {
		return errors.New("operator must be a string")
	}
	if pattern, ok := data["computername_pattern"].(string); ok {
		r.ComputernamePattern = pattern
	} else {
		return errors.New("computername_pattern must be a string")
	}
	if domain, ok := data["domain"].(string); ok {
		r.Domain = domain
	}
	return nil
}

// ArkUAPSIAVMIPRule defines a specific logical name rule used to match a given IP+logical name.
type ArkUAPSIAVMIPRule struct {
	Operator    string   `json:"operator" mapstructure:"operator" flag:"operator" choices:"EXACTLY,WILDCARD" desc:"The operator to use for matching the IP addresses. Valid values are EXACTLY and WILDCARD."`
	IPAddresses []string `json:"ip_addresses" mapstructure:"ip_addresses" flag:"ip-addresses" validate:"max=1000" desc:"The list of IP addresses to match against. This can include both IPv4 and IPv6 addresses."`
	LogicalName string   `json:"logical_name" mapstructure:"logical_name" flag:"logical-name" validate:"min=1,max=256" desc:"The logical name of the network."`
}

// Serialize converts the IP rule to a map.
func (r *ArkUAPSIAVMIPRule) Serialize() map[string]interface{} {
	return map[string]interface{}{
		"operator":    r.Operator,
		"ipAddresses": r.IPAddresses,
		"logicalName": r.LogicalName,
	}
}

// Deserialize populates the IP rule from a map.
func (r *ArkUAPSIAVMIPRule) Deserialize(data map[string]interface{}) error {
	if operator, ok := data["operator"].(string); ok {
		r.Operator = operator
	} else {
		return errors.New("operator must be a string")
	}
	if ips, ok := data["ip_addresses"].([]interface{}); ok {
		for _, ip := range ips {
			if str, ok := ip.(string); ok {
				r.IPAddresses = append(r.IPAddresses, str)
			} else {
				return errors.New("ip_addresses must be a list of strings")
			}
		}
	}
	if logicalName, ok := data["logical_name"].(string); ok {
		r.LogicalName = logicalName
	} else {
		return errors.New("logical_name must be a string")
	}
	return nil
}

// ArkUAPSIAVMAWSResource represents the AWS resources for a virtual machine access policy.
type ArkUAPSIAVMAWSResource struct {
	Regions    []string               `json:"regions" mapstructure:"regions" flag:"regions" desc:"The AWS regions where the resources are located."`
	Tags       []ArkUAPSIAVMKeyValTag `json:"tags" mapstructure:"tags" flag:"tags" desc:"The tags used to match AWS resources. This is a list of key-value pairs."`
	VPCIDs     []string               `json:"vpc_ids" mapstructure:"vpc_ids" flag:"vpc-ids" desc:"The VPC IDs where the resources are located. This is a list of VPC identifiers."`
	AccountIDs []string               `json:"account_ids" mapstructure:"account_ids" flag:"account-ids" desc:"The AWS account IDs where the resources are located. This is a list of account identifiers."`
}

// Serialize converts the AWS resource to a map.
func (r *ArkUAPSIAVMAWSResource) Serialize() map[string]interface{} {
	tags := make([]map[string]interface{}, len(r.Tags))
	for i, tag := range r.Tags {
		tags[i] = tag.Serialize()
	}
	return map[string]interface{}{
		"regions":    r.Regions,
		"tags":       tags,
		"vpcIds":     r.VPCIDs,
		"accountIds": r.AccountIDs,
	}
}

// Deserialize populates the AWS resource from a map.
func (r *ArkUAPSIAVMAWSResource) Deserialize(data map[string]interface{}) error {
	if regions, ok := data["regions"].([]interface{}); ok {
		for _, region := range regions {
			if str, ok := region.(string); ok {
				r.Regions = append(r.Regions, str)
			} else {
				return errors.New("regions must be a list of strings")
			}
		}
	}
	if tags, ok := data["tags"].([]interface{}); ok {
		for _, tag := range tags {
			if tagMap, ok := tag.(map[string]interface{}); ok {
				var t ArkUAPSIAVMKeyValTag
				if err := t.Deserialize(tagMap); err != nil {
					return err
				}
				r.Tags = append(r.Tags, t)
			} else {
				return errors.New("tags must be a list of maps")
			}
		}
	}
	if vpcIDs, ok := data["vpc_ids"].([]interface{}); ok {
		for _, vpc := range vpcIDs {
			if str, ok := vpc.(string); ok {
				r.VPCIDs = append(r.VPCIDs, str)
			} else {
				return errors.New("vpc_ids must be a list of strings")
			}
		}
	}
	if accountIDs, ok := data["account_ids"].([]interface{}); ok {
		for _, account := range accountIDs {
			if str, ok := account.(string); ok {
				r.AccountIDs = append(r.AccountIDs, str)
			} else {
				return errors.New("account_ids must be a list of strings")
			}
		}
	}
	return nil
}

// ArkUAPSIAVMAzureResource represents the Azure resources for a virtual machine access policy.
type ArkUAPSIAVMAzureResource struct {
	Regions        []string               `json:"regions" mapstructure:"regions" flag:"regions" desc:"The Azure regions where the resources are located."`
	Tags           []ArkUAPSIAVMKeyValTag `json:"tags" mapstructure:"tags" flag:"tags" desc:"The tags used to match Azure resources. This is a list of key-value pairs."`
	ResourceGroups []string               `json:"resource_groups" mapstructure:"resource_groups" flag:"resource-groups" desc:"The Azure resource groups where the resources are located. This is a list of resource group names."`
	VNetIDs        []string               `json:"vnet_ids" mapstructure:"vnet_ids" flag:"vnet-ids" desc:"The Azure Virtual Network IDs where the resources are located. This is a list of VNet identifiers."`
	Subscriptions  []string               `json:"subscriptions" mapstructure:"subscriptions" flag:"subscriptions" desc:"The Azure subscription IDs where the resources are located. This is a list of subscription identifiers."`
}

// Serialize converts the Azure resource to a map.
func (r *ArkUAPSIAVMAzureResource) Serialize() map[string]interface{} {
	tags := make([]map[string]interface{}, len(r.Tags))
	for i, tag := range r.Tags {
		tags[i] = tag.Serialize()
	}
	return map[string]interface{}{
		"regions":        r.Regions,
		"tags":           tags,
		"resourceGroups": r.ResourceGroups,
		"vnetIds":        r.VNetIDs,
		"subscriptions":  r.Subscriptions,
	}
}

// Deserialize populates the Azure resource from a map.
func (r *ArkUAPSIAVMAzureResource) Deserialize(data map[string]interface{}) error {
	if regions, ok := data["regions"].([]interface{}); ok {
		for _, region := range regions {
			if str, ok := region.(string); ok {
				r.Regions = append(r.Regions, str)
			} else {
				return errors.New("regions must be a list of strings")
			}
		}
	}
	if tags, ok := data["tags"].([]interface{}); ok {
		for _, tag := range tags {
			if tagMap, ok := tag.(map[string]interface{}); ok {
				var t ArkUAPSIAVMKeyValTag
				if err := t.Deserialize(tagMap); err != nil {
					return err
				}
				r.Tags = append(r.Tags, t)
			} else {
				return errors.New("tags must be a list of maps")
			}
		}
	}
	if resourceGroups, ok := data["resource_groups"].([]interface{}); ok {
		for _, group := range resourceGroups {
			if str, ok := group.(string); ok {
				r.ResourceGroups = append(r.ResourceGroups, str)
			} else {
				return errors.New("resource_groups must be a list of strings")
			}
		}
	}
	if vnetIDs, ok := data["vnet_ids"].([]interface{}); ok {
		for _, vnet := range vnetIDs {
			if str, ok := vnet.(string); ok {
				r.VNetIDs = append(r.VNetIDs, str)
			} else {
				return errors.New("vnet_ids must be a list of strings")
			}
		}
	}
	if subscriptions, ok := data["subscriptions"].([]interface{}); ok {
		for _, sub := range subscriptions {
			if str, ok := sub.(string); ok {
				r.Subscriptions = append(r.Subscriptions, str)
			} else {
				return errors.New("subscriptions must be a list of strings")
			}
		}
	}
	return nil
}

// ArkUAPSIAVMGCPResource represents the GCP resources for a virtual machine access policy.
type ArkUAPSIAVMGCPResource struct {
	Regions  []string               `json:"regions" mapstructure:"regions" flag:"regions" desc:"The GCP regions where the resources are located."`
	Labels   []ArkUAPSIAVMKeyValTag `json:"labels" mapstructure:"labels" flag:"labels" desc:"The labels used to match GCP resources. This is a list of key-value pairs."`
	VPCIDs   []string               `json:"vpc_ids" mapstructure:"vpc_ids" flag:"vpc-ids" desc:"The GCP VPC IDs where the resources are located. This is a list of VPC identifiers."`
	Projects []string               `json:"projects" mapstructure:"projects" flag:"projects" desc:"The GCP project IDs where the resources are located. This is a list of project identifiers."`
}

// Serialize converts the GCP resource to a map.
func (r *ArkUAPSIAVMGCPResource) Serialize() map[string]interface{} {
	labels := make([]map[string]interface{}, len(r.Labels))
	for i, label := range r.Labels {
		labels[i] = label.Serialize()
	}
	return map[string]interface{}{
		"regions":  r.Regions,
		"labels":   labels,
		"vpcIds":   r.VPCIDs,
		"projects": r.Projects,
	}
}

// Deserialize populates the GCP resource from a map.
func (r *ArkUAPSIAVMGCPResource) Deserialize(data map[string]interface{}) error {
	if regions, ok := data["regions"].([]interface{}); ok {
		for _, region := range regions {
			if str, ok := region.(string); ok {
				r.Regions = append(r.Regions, str)
			} else {
				return errors.New("regions must be a list of strings")
			}
		}
	}
	if labels, ok := data["labels"].([]interface{}); ok {
		for _, label := range labels {
			if labelMap, ok := label.(map[string]interface{}); ok {
				var l ArkUAPSIAVMKeyValTag
				if err := l.Deserialize(labelMap); err != nil {
					return err
				}
				r.Labels = append(r.Labels, l)
			} else {
				return errors.New("labels must be a list of maps")
			}
		}
	}
	if vpcIDs, ok := data["vpc_ids"].([]interface{}); ok {
		for _, vpc := range vpcIDs {
			if str, ok := vpc.(string); ok {
				r.VPCIDs = append(r.VPCIDs, str)
			} else {
				return errors.New("vpc_ids must be a list of strings")
			}
		}
	}
	if projects, ok := data["projects"].([]interface{}); ok {
		for _, project := range projects {
			if str, ok := project.(string); ok {
				r.Projects = append(r.Projects, str)
			} else {
				return errors.New("projects must be a list of strings")
			}
		}
	}
	return nil
}

// ArkUAPSIAVMOnPremResource represents the on-premises resources for a virtual machine access policy.
type ArkUAPSIAVMOnPremResource struct {
	FQDNRules []ArkUAPSIAVMFQDNRule `json:"fqdn_rules,omitempty" mapstructure:"fqdn_rules" flag:"fqdn-rules" desc:"The FQDN rules used to match DNS records. This is a list of FQDN rules."`
	IPRules   []ArkUAPSIAVMIPRule   `json:"ip_rules,omitempty" mapstructure:"ip_rules" flag:"ip-rules" desc:"The IP rules used to match IP addresses and logical names. This is a list of IP rules."`
}

// Serialize converts the on-prem resource to a map.
func (r *ArkUAPSIAVMOnPremResource) Serialize() map[string]interface{} {
	fqdnRules := make([]map[string]interface{}, len(r.FQDNRules))
	for i, rule := range r.FQDNRules {
		fqdnRules[i] = rule.Serialize()
	}
	ipRules := make([]map[string]interface{}, len(r.IPRules))
	for i, rule := range r.IPRules {
		ipRules[i] = rule.Serialize()
	}
	return map[string]interface{}{
		"fqdnRules": fqdnRules,
		"ipRules":   ipRules,
	}
}

// Deserialize populates the on-prem resource from a map.
func (r *ArkUAPSIAVMOnPremResource) Deserialize(data map[string]interface{}) error {
	if fqdnRules, ok := data["fqdn_rules"].([]interface{}); ok {
		for _, rule := range fqdnRules {
			if ruleMap, ok := rule.(map[string]interface{}); ok {
				var fqdnRule ArkUAPSIAVMFQDNRule
				if err := fqdnRule.Deserialize(ruleMap); err != nil {
					return err
				}
				r.FQDNRules = append(r.FQDNRules, fqdnRule)
			} else {
				return errors.New("fqdn_rules must be a list of maps")
			}
		}
	}
	if ipRules, ok := data["ip_rules"].([]interface{}); ok {
		for _, rule := range ipRules {
			if ruleMap, ok := rule.(map[string]interface{}); ok {
				var ipRule ArkUAPSIAVMIPRule
				if err := ipRule.Deserialize(ruleMap); err != nil {
					return err
				}
				r.IPRules = append(r.IPRules, ipRule)
			} else {
				return errors.New("ip_rules must be a list of maps")
			}
		}
	}
	return nil
}

// ArkUAPSIAVMPlatformTargets represents the targets for a virtual machine access policy.
type ArkUAPSIAVMPlatformTargets struct {
	AWSResource    *ArkUAPSIAVMAWSResource    `json:"aws_resource,omitempty" mapstructure:"aws_resource" flag:"aws-resource" desc:"The AWS resources for the VM access policy. This includes regions, tags, VPC IDs, and account IDs."`
	AzureResource  *ArkUAPSIAVMAzureResource  `json:"azure_resource,omitempty" mapstructure:"azure_resource" flag:"azure-resource" desc:"The Azure resources for the VM access policy. This includes regions, tags, resource groups, VNet IDs, and subscriptions."`
	GCPResource    *ArkUAPSIAVMGCPResource    `json:"gcp_resource,omitempty" mapstructure:"gcp_resource" flag:"gcp-resource" desc:"The GCP resources for the VM access policy. This includes regions, labels, VPC IDs, and project IDs."`
	OnPremResource *ArkUAPSIAVMOnPremResource `json:"onprem_resource,omitempty" mapstructure:"onprem_resource" flag:"onprem-resource" desc:"The on-premises resources for the VM access policy. This includes FQDN rules and IP rules."`
}

// Serialize converts the platform targets to a map.
func (t *ArkUAPSIAVMPlatformTargets) Serialize(workspace string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	if workspace == common.WorkspaceTypeAWS && t.AWSResource != nil {
		data[common.WorkspaceTypeAWS] = t.AWSResource.Serialize()
	} else if workspace == common.WorkspaceTypeAzure && t.AzureResource != nil {
		data[common.WorkspaceTypeAzure] = t.AzureResource.Serialize()
	} else if workspace == common.WorkspaceTypeGCP && t.GCPResource != nil {
		data[common.WorkspaceTypeGCP] = t.GCPResource.Serialize()
	} else if workspace == common.WorkspaceTypeFQDNIP && t.OnPremResource != nil {
		data[common.WorkspaceTypeFQDNIP] = t.OnPremResource.Serialize()
	} else {
		return nil, errors.New("unsupported workspace type or missing resource")
	}
	return data, nil
}

// Deserialize populates the platform targets from a map.
func (t *ArkUAPSIAVMPlatformTargets) Deserialize(data map[string]interface{}, workspace string) error {
	if workspace == common.WorkspaceTypeAWS {
		if awsData, ok := data[common.WorkspaceTypeAWS].(map[string]interface{}); ok {
			t.AWSResource = &ArkUAPSIAVMAWSResource{}
			if err := t.AWSResource.Deserialize(awsData); err != nil {
				return err
			}
		} else {
			return errors.New("missing AWS resource data")
		}
	} else if workspace == common.WorkspaceTypeAzure {
		if azureData, ok := data[common.WorkspaceTypeAzure].(map[string]interface{}); ok {
			t.AzureResource = &ArkUAPSIAVMAzureResource{}
			if err := t.AzureResource.Deserialize(azureData); err != nil {
				return err
			}
		} else {
			return errors.New("missing Azure resource data")
		}
	} else if workspace == common.WorkspaceTypeGCP {
		if gcpData, ok := data[common.WorkspaceTypeGCP].(map[string]interface{}); ok {
			t.GCPResource = &ArkUAPSIAVMGCPResource{}
			if err := t.GCPResource.Deserialize(gcpData); err != nil {
				return err
			}
		} else {
			return errors.New("missing GCP resource data")
		}
	} else if workspace == common.WorkspaceTypeFQDNIP {
		if onPremData, ok := data[common.WorkspaceTypeFQDNIP].(map[string]interface{}); ok {
			t.OnPremResource = &ArkUAPSIAVMOnPremResource{}
			if err := t.OnPremResource.Deserialize(onPremData); err != nil {
				return err
			}
		} else {
			return errors.New("missing FQDN/IP resource data")
		}
	} else {
		return errors.New("unsupported workspace type")
	}
	return nil
}
