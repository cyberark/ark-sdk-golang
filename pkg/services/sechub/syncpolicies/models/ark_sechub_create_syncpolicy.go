package models

// ArkSechubCreateSyncPolicy represents a sync policy for the sechub service.
// It includes information about the policy name, desc, source, target, filter, and transformation.
type ArkSechubCreateSyncPolicy struct {
	Name           string                            `json:"name" mapstructure:"name" desc:"Name of the sync policy" flag:"name" validate:"required"`
	Description    string                            `json:"desc,omitempty" mapstructure:"desc,omitempty" desc:"Description of the sync policy" flag:"desc,omitempty"`
	Source         ArkSechubSyncPolicyStore          `json:"source" mapstructure:"source" desc:"Source store configuration" flag:"source"`
	Target         ArkSechubSyncPolicyStore          `json:"target" mapstructure:"target" desc:"Target store configuration" flag:"target"`
	Filter         ArkSechubSyncPolicyFilter         `json:"filter" mapstructure:"filter" desc:"Filter for selecting items to sync" flag:"filter"`
	Transformation ArkSechubSyncPolicyTransformation `json:"transformation" mapstructure:"transformation" desc:"Transformation to apply during sync" flag:"transformation"`
}

// ArkSechubSyncPolicyStore represents a store configuration with an ID.
type ArkSechubSyncPolicyStore struct {
	ID string `json:"id" mapstructure:"id" desc:"Unique identifier of the store" flag:"id" validate:"required"`
}

// ArkSechubSyncPolicyFilter represents a filter for selecting items to sync.
type ArkSechubSyncPolicyFilter struct {
	ID   string                        `json:"id,omitempty" mapstructure:"id,omitempty" desc:"The unique filter identifier. If used, Data and Type are not required." flag:"id,omitempty"`
	Data ArkSechubSyncPolicyFilterData `json:"data,omitzero" mapstructure:"data,omitzero" desc:"Filter data" flag:"data,omitempty"`
	Type string                        `json:"type,omitempty" mapstructure:"type,omitempty" desc:"Type of filter (PAM_SAFE)" flag:"type,omitempty" choices:"PAM_SAFE"`
}

// ArkSechubSyncPolicyFilterData represents filter data for the sync policy.
type ArkSechubSyncPolicyFilterData struct {
	SafeName string `json:"safe_name" mapstructure:"safe_name" desc:"Name of the safe to filter" validate:"required" flag:"safe-name"`
}

// ArkSechubSyncPolicyTransformation represents transformation configuration for the sync policy.
type ArkSechubSyncPolicyTransformation struct {
	Predefined string `json:"predefined" mapstructure:"predefined" desc:"Predefined transformation to apply (password_only_plain_text,default)" flag:"predefined" choices:"password_only_plain_text,default"`
}
