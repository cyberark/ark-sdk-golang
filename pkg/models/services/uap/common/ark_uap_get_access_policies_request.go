package common

import "strings"

// Possible ArkUAPFilterOperator
const (
	ArkUAPFilterOperatorEQ       = "eq"
	ArkUAPFilterOperatorContains = "contains"
	ArkUAPFilterOperatorOR       = "or"
	ArkUAPFilterOperatorAND      = "and"
)

// ArkUAPDefaultLimitSize defines the default limit size for UAP access policies.
const (
	ArkUAPDefaultLimitSize = 50
)

// filterOperators maps field names to their corresponding filter operators.
var filterOperators = map[string]string{
	"locationType":   ArkUAPFilterOperatorEQ,
	"policyType":     ArkUAPFilterOperatorEQ,
	"targetCategory": ArkUAPFilterOperatorEQ,
	"policyTags":     ArkUAPFilterOperatorEQ,
	"status":         ArkUAPFilterOperatorEQ,
	"identities":     ArkUAPFilterOperatorContains,
}

// mapAliasToFieldName maps alias names to their corresponding field names.
var mapAliasToFieldName = map[string]string{
	"locationType":   "LocationType",
	"policyType":     "PolicyType",
	"targetCategory": "TargetCategory",
	"policyTags":     "PolicyTags",
	"status":         "Status",
	"identities":     "Identities",
}

// ArkUAPGetQueryParams represents the query parameters for retrieving access policies.
type ArkUAPGetQueryParams struct {
	Filter               string `json:"filter,omitempty" mapstructure:"filter,omitempty" flag:"filter" desc:"The filter query to apply on the policies"`
	ShowEditablePolicies bool   `json:"show_editable_policies,omitempty" mapstructure:"show_editable_policies,omitempty" flag:"show-editable-policies" desc:"Show editable policies"`
	Q                    string `json:"q,omitempty" mapstructure:"q,omitempty" flag:"q" desc:"Free text search on policy name or description"`
	NextToken            string `json:"next_token,omitempty" mapstructure:"next_token,omitempty" flag:"next-token" desc:"The next token for pagination"`
	Limit                int    `json:"limit" mapstructure:"limit" flag:"limit" desc:"The maximum number of policies to return in the response"`
}

// ArkUAPFilters represents the filters for UAP policies.
type ArkUAPFilters struct {
	LocationType         []string `json:"location_type,omitempty" mapstructure:"location_type,omitempty" flag:"location-type" desc:"List of wanted location types for the policies"`
	TargetCategory       []string `json:"target_category,omitempty" mapstructure:"target_category,omitempty" flag:"target-category" desc:"List of wanted target categories for the policies"`
	PolicyType           []string `json:"policy_type,omitempty" mapstructure:"policy_type,omitempty" flag:"policy-type" desc:"List of wanted policy types for the policies"`
	PolicyTags           []string `json:"policy_tags,omitempty" mapstructure:"policy_tags,omitempty" flag:"policy-tags" desc:"List of wanted policy tags for the policies"`
	Identities           []string `json:"identities,omitempty" mapstructure:"identities,omitempty" flag:"identities" desc:"List of identities to filter the policies by"`
	Status               []string `json:"status,omitempty" mapstructure:"status,omitempty" flag:"status" desc:"List of wanted policy statuses for the policies"`
	TextSearch           string   `json:"text_search,omitempty" mapstructure:"text_search,omitempty" flag:"text-search" desc:"Text search filter to apply on the policies names and descriptions"`
	ShowEditablePolicies bool     `json:"show_editable_policies,omitempty" mapstructure:"show_editable_policies,omitempty" flag:"show-editable-policies" desc:"Whether to show editable policies or not" default:"true"`
	MaxPages             int      `json:"max_pages" mapstructure:"max_pages" flag:"max-pages" desc:"The maximum number of pages for pagination, default is 1000000" default:"1000000"`
}

// NewArkUAPFilters initializes a new instance of ArkUAPFilters with default values.
func NewArkUAPFilters() *ArkUAPFilters {
	return &ArkUAPFilters{
		LocationType:         []string{},
		TargetCategory:       []string{},
		PolicyType:           []string{},
		PolicyTags:           []string{},
		Identities:           []string{},
		Status:               []string{},
		TextSearch:           "",
		ShowEditablePolicies: true,
		MaxPages:             1000000,
	}
}

// BuildFilterQueryFromFilters constructs a filter query string from the provided filters.
func (filters *ArkUAPFilters) BuildFilterQueryFromFilters() string {
	var clauses []string

	for fieldName, operator := range filterOperators {
		alias := mapAliasToFieldName[fieldName]
		var values []string

		switch alias {
		case "LocationType":
			values = filters.LocationType
		case "PolicyType":
			values = filters.PolicyType
		case "TargetCategory":
			values = filters.TargetCategory
		case "PolicyTags":
			values = filters.PolicyTags
		case "Status":
			values = filters.Status
		case "Identities":
			values = filters.Identities
		}

		if len(values) > 0 {
			var itemClauses []string
			for _, v := range values {
				itemClauses = append(itemClauses, "("+fieldName+" "+operator+" '"+v+"')")
			}
			joined := strings.Join(itemClauses, " "+ArkUAPFilterOperatorOR+" ")
			if len(itemClauses) > 1 {
				clauses = append(clauses, "("+joined+")")
			} else {
				clauses = append(clauses, joined)
			}
		}
	}

	if len(clauses) > 1 {
		return "(" + strings.Join(clauses, " and ") + ")"
	} else if len(clauses) == 1 {
		return clauses[0]
	}
	return ""
}

// ArkUAPGetAccessPoliciesRequest represents the request to get access policies.
type ArkUAPGetAccessPoliciesRequest struct {
	Filters   *ArkUAPFilters `json:"filters,omitempty" mapstructure:"filters,omitempty" flag:"filters" desc:"The filter query to apply on the policies"`
	Limit     int            `json:"limit" mapstructure:"limit" flag:"limit" desc:"The maximum number of policies to return in the response" default:"50"`
	NextToken string         `json:"next_token,omitempty" mapstructure:"next_token,omitempty" flag:"next-token" desc:"The next token for pagination"`
}

// BuildGetQueryParams constructs the query parameters for retrieving access policies.
func (request *ArkUAPGetAccessPoliciesRequest) BuildGetQueryParams() ArkUAPGetQueryParams {
	queryParams := ArkUAPGetQueryParams{
		Limit: request.Limit,
	}
	if queryParams.Limit <= 0 {
		queryParams.Limit = ArkUAPDefaultLimitSize
	}

	if request.Filters == nil {
		return queryParams
	}

	localFilters := request.Filters

	if localFilters.TextSearch != "" {
		queryParams.Q = localFilters.TextSearch
	}

	filterQuery := localFilters.BuildFilterQueryFromFilters()
	if filterQuery != "" {
		queryParams.Filter = filterQuery
	}

	if request.NextToken != "" {
		queryParams.NextToken = request.NextToken
	}

	queryParams.ShowEditablePolicies = localFilters.ShowEditablePolicies

	return queryParams
}
