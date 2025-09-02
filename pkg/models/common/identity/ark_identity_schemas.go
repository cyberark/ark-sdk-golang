// Package identity provides data structures and types for ARK Identity API operations.
// This package contains models for authentication, authorization, and identity management
// operations including authentication challenges, tokens, tenant information, and
// various response structures used in ARK Identity service interactions.
package identity

import (
	"strings"
)

// BaseIdentityAPIResponse represents the common response structure from the ARK Identity API.
// This structure contains standard fields that are present in most Identity API responses
// including success status, error information, and diagnostic details for troubleshooting
// API interactions.
type BaseIdentityAPIResponse struct {
	Success   bool   `json:"Success" validate:"required"`
	Exception string `json:"Exception"`
	ErrorCode string `json:"ErrorCode"`
	Message   string `json:"Message"`
	ErrorID   string `json:"ErrorID"`
}

// PodFqdnResult represents the result containing Pod Fully Qualified Domain Name information.
// This structure contains the Pod FQDN which is used to identify the specific Identity
// service instance and extract tenant information for multi-tenant operations.
type PodFqdnResult struct {
	PodFqdn string `json:"PodFqdn" validate:"required,min=2"`
}

// GetTenantID extracts the tenant identifier from the Pod FQDN.
// It parses the PodFqdn field by splitting on the first dot character and
// returns the leftmost component, which represents the tenant ID in the
// ARK Identity service naming convention.
//
// Returns:
//   - string: The tenant ID extracted from the Pod FQDN, or empty string if PodFqdn is empty
//
// Example:
//
//	podResult := &PodFqdnResult{PodFqdn: "tenant123.example.com"}
//	tenantID := podResult.GetTenantID() // Returns "tenant123"
func (p *PodFqdnResult) GetTenantID() string {
	return strings.Split(p.PodFqdn, ".")[0]
}

// AdvanceAuthResult represents the result of a successful advanced authentication request.
// This structure contains authentication tokens, user information, and session details
// returned when authentication is completed successfully through the ARK Identity system.
type AdvanceAuthResult struct {
	DisplayName   string `json:"DisplayName" validate:"omitempty,min=2"`
	Auth          string `json:"Auth" validate:"required,min=2"`
	Summary       string `json:"Summary" validate:"required,min=2"`
	Token         string `json:"Token" validate:"omitempty,min=2"`
	RefreshToken  string `json:"RefreshToken" validate:"omitempty,min=2"`
	TokenLifetime int    `json:"TokenLifetime"`
	CustomerID    string `json:"CustomerID"`
	UserID        string `json:"UserId"`
	PodFqdn       string `json:"PodFqdn"`
}

// AdvanceAuthMidResult represents the result of an in-progress advanced authentication request.
// This structure contains intermediate state information when authentication is still
// being processed and polling is required to check for completion.
type AdvanceAuthMidResult struct {
	Summary            string `json:"Summary" validate:"required,min=2"`
	GeneratedAuthValue string `json:"GeneratedAuthValue"`
}

// Mechanism represents an authentication mechanism within an authentication challenge.
// This structure defines the properties and prompts for a specific authentication
// method that can be used to complete an authentication challenge in the ARK Identity system.
type Mechanism struct {
	AnswerType       string `json:"AnswerType" validate:"required,min=2"`
	Name             string `json:"Name" validate:"required,min=2"`
	PromptMechChosen string `json:"PromptMechChosen" validate:"required,min=2"`
	PromptSelectMech string `json:"PromptSelectMech" validate:"omitempty,min=2"`
	MechanismID      string `json:"MechanismId" validate:"required,min=2"`
}

// Challenge represents an authentication challenge containing available mechanisms.
// This structure groups one or more authentication mechanisms that can be used
// to satisfy an authentication requirement in the multi-factor authentication flow.
type Challenge struct {
	Mechanisms []Mechanism `json:"Mechanisms" validate:"required,dive,required"`
}

// StartAuthResult represents the result of initiating an authentication request.
// This structure contains authentication challenges, session information, and
// Identity Provider (IdP) redirect details for starting the authentication process.
type StartAuthResult struct {
	Challenges            []Challenge `json:"Challenges" validate:"omitempty,dive,required"`
	SessionID             string      `json:"SessionId" validate:"omitempty,min=2"`
	IdpRedirectURL        string      `json:"IdpRedirectUrl"`
	IdpLoginSessionID     string      `json:"IdpLoginSessionId"`
	IdpRedirectShortURL   string      `json:"IdpRedirectShortUrl"`
	IdpShortURLID         string      `json:"IdpShortUrlId"`
	IdpOobAuthPinRequired bool        `json:"IdpOobAuthPinRequired"`
	TenantID              string      `json:"TenantId"`
}

// IdpAuthStatusResult represents the result of an Identity Provider authentication status check.
// This structure contains the current state of an IdP authentication session
// along with token information when authentication is completed.
type IdpAuthStatusResult struct {
	AuthLevel     string `json:"AuthLevel"`
	DisplayName   string `json:"DisplayName"`
	Auth          string `json:"Auth"`
	UserID        string `json:"UserId"`
	State         string `json:"State"`
	TokenLifetime int    `json:"TokenLifetime"`
	Token         string `json:"Token"`
	RefreshToken  string `json:"RefreshToken"`
	EmailAddress  string `json:"EmailAddress"`
	UserDirectory string `json:"UserDirectory"`
	PodFqdn       string `json:"PodFqdn"`
	User          string `json:"User"`
	CustomerID    string `json:"CustomerID"`
	Forest        string `json:"Forest"`
	SystemID      string `json:"SystemID"`
	SourceDsType  string `json:"SourceDsType"`
	Summary       string `json:"Summary"`
}

// TenantFqdnResponse represents the complete response for tenant FQDN requests.
// This structure combines the base API response with Pod FQDN result data
// for tenant identification and service endpoint discovery operations.
type TenantFqdnResponse struct {
	BaseIdentityAPIResponse
	Result PodFqdnResult `json:"Result"`
}

// AdvanceAuthMidResponse represents the complete response for in-progress authentication requests.
// This structure combines the base API response with intermediate authentication
// state information for polling-based authentication flows.
type AdvanceAuthMidResponse struct {
	BaseIdentityAPIResponse
	Result AdvanceAuthMidResult `json:"Result"`
}

// AdvanceAuthResponse represents the complete response for successful authentication requests.
// This structure combines the base API response with authentication tokens and
// user information when authentication has been completed successfully.
type AdvanceAuthResponse struct {
	BaseIdentityAPIResponse
	Result AdvanceAuthResult `json:"Result"`
}

// StartAuthResponse represents the complete response for authentication initiation requests.
// This structure combines the base API response with authentication challenges
// and session information for starting the authentication process.
type StartAuthResponse struct {
	BaseIdentityAPIResponse
	Result StartAuthResult `json:"Result"`
}

// GetTenantSuffixResult represents the complete response for tenant suffix requests.
// This structure combines the base API response with a flexible result map
// containing tenant-specific configuration and suffix information.
type GetTenantSuffixResult struct {
	BaseIdentityAPIResponse
	Result map[string]interface{} `json:"Result"`
}

// IdpAuthStatusResponse represents the complete response for IdP authentication status requests.
// This structure combines the base API response with Identity Provider authentication
// status and token information for federated authentication flows.
type IdpAuthStatusResponse struct {
	BaseIdentityAPIResponse
	Result IdpAuthStatusResult `json:"Result"`
}

// TenantEndpointResponse represents the response containing tenant endpoint information.
// This structure provides the endpoint URL for accessing tenant-specific services
// in the ARK Identity system.
type TenantEndpointResponse struct {
	Endpoint string `json:"endpoint"`
}
