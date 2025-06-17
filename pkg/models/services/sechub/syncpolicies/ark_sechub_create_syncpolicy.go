package syncpolicies

// ArkSechubCreateSyncPolicy represents a sync policy for the sechub service.
// It includes information about the policy name, description, source, target, filter, and transformation.
type ArkSechubCreateSyncPolicy struct {
	Name           string                            `json:"name" mapstructure:"name" description:"Name of the sync policy" validate:"required"`
	Description    string                            `json:"description,omitempty" mapstructure:"description,omitempty" description:"Description of the sync policy"`
	Source         ArkSechubSyncPolicyStore          `json:"source" mapstructure:"source" description:"Source store configuration"`
	Target         ArkSechubSyncPolicyStore          `json:"target" mapstructure:"target" description:"Target store configuration"`
	Filter         ArkSechubSyncPolicyFilter         `json:"filter" mapstructure:"filter" description:"Filter for selecting items to sync"`
	Transformation ArkSechubSyncPolicyTransformation `json:"transformation" mapstructure:"transformation" description:"Transformation to apply during sync"`
}

// ArkSechubSyncPolicyStore represents a store configuration with an ID.
type ArkSechubSyncPolicyStore struct {
	ID string `json:"id" mapstructure:"id" description:"Unique identifier of the store"`
}

// ArkSechubSyncPolicyFilter represents a filter for selecting items to sync.
type ArkSechubSyncPolicyFilter struct {
	ID   string                        `json:"id,omitempty" mapstructure:"id,omitempty" description:"The unique filter identifier. If used, Data and Type are not required."`
	Data ArkSechubSyncPolicyFilterData `json:"data,omitempty" mapstructure:"data,omitempty" description:"Filter data"`
	Type string                        `json:"type,omitempty" mapstructure:"type,omitempty" description:"Type of filter - Allowed Value: 'PAM_SAFE'"`
}

// ArkSechubSyncPolicyFilterData represents filter data for the sync policy.
type ArkSechubSyncPolicyFilterData struct {
	SafeName string `json:"safe_name" mapstructure:"safe_name" description:"Name of the safe to filter" validate:"required"`
}

// ArkSechubSyncPolicyTransformation represents transformation configuration for the sync policy.
type ArkSechubSyncPolicyTransformation struct {
	Predefined string `json:"predefined" mapstructure:"predefined" description:"Predefined transformation to apply - Allowed Values: 'password_only_plain_text, default'" default:"default"`
}
