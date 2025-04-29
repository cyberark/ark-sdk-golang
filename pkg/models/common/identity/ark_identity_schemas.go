package identity

import (
	"strings"
)

// BaseIdentityAPIResponse is a struct that represents the base response from the Identity API.
type BaseIdentityAPIResponse struct {
	Success   bool   `json:"Success" validate:"required"`
	Exception string `json:"Exception"`
	ErrorCode string `json:"ErrorCode"`
	Message   string `json:"Message"`
	ErrorID   string `json:"ErrorID"`
}

// PodFqdnResult is a struct that represents the result of a Pod FQDN request.
type PodFqdnResult struct {
	PodFqdn string `json:"PodFqdn" validate:"required,min=2"`
}

// GetTenantID extracts the tenant ID from the Pod FQDN.
func (p *PodFqdnResult) GetTenantID() string {
	return strings.Split(p.PodFqdn, ".")[0]
}

// AdvanceAuthResult is a struct that represents the result of an advanced authentication successful request.
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

// AdvanceAuthMidResult is a struct that represents the result of an advanced authentication polling / not finished request.
type AdvanceAuthMidResult struct {
	Summary            string `json:"Summary" validate:"required,min=2"`
	GeneratedAuthValue string `json:"GeneratedAuthValue"`
}

// Mechanism is a struct that represents a mechanism in the authentication process as part of a challenge.
type Mechanism struct {
	AnswerType       string `json:"AnswerType" validate:"required,min=2"`
	Name             string `json:"Name" validate:"required,min=2"`
	PromptMechChosen string `json:"PromptMechChosen" validate:"required,min=2"`
	PromptSelectMech string `json:"PromptSelectMech" validate:"omitempty,min=2"`
	MechanismID      string `json:"MechanismId" validate:"required,min=2"`
}

// Challenge is a struct that represents a challenge in the authentication process.
type Challenge struct {
	Mechanisms []Mechanism `json:"Mechanisms" validate:"required,dive,required"`
}

// StartAuthResult is a struct that represents the result of a start authentication request.
type StartAuthResult struct {
	Challenges          []Challenge `json:"Challenges" validate:"omitempty,dive,required"`
	SessionID           string      `json:"SessionId" validate:"omitempty,min=2"`
	IdpRedirectURL      string      `json:"IdpRedirectUrl"`
	IdpLoginSessionID   string      `json:"IdpLoginSessionId"`
	IdpRedirectShortURL string      `json:"IdpRedirectShortUrl"`
	IdpShortURLID       string      `json:"IdpShortUrlId"`
	TenantID            string      `json:"TenantId"`
}

// IdpAuthStatusResult is a struct that represents the result of an IdP authentication status request.
type IdpAuthStatusResult struct {
	State         string `json:"State" validate:"required"`
	TokenLifetime int    `json:"TokenLifetime"`
	Token         string `json:"Token"`
	RefreshToken  string `json:"RefreshToken"`
}

// TenantFqdnResponse is a struct that represents the response from the Identity API for tenant FQDN.
type TenantFqdnResponse struct {
	BaseIdentityAPIResponse
	Result PodFqdnResult `json:"Result"`
}

// AdvanceAuthMidResponse is a struct that represents the response from the Identity API for advanced authentication polling / not finished.
type AdvanceAuthMidResponse struct {
	BaseIdentityAPIResponse
	Result AdvanceAuthMidResult `json:"Result"`
}

// AdvanceAuthResponse is a struct that represents the response from the Identity API for advanced authentication successful.
type AdvanceAuthResponse struct {
	BaseIdentityAPIResponse
	Result AdvanceAuthResult `json:"Result"`
}

// StartAuthResponse is a struct that represents the response from the Identity API for starting authentication.
type StartAuthResponse struct {
	BaseIdentityAPIResponse
	Result StartAuthResult `json:"Result"`
}

// GetTenantSuffixResult is a struct that represents the response from the Identity API for getting tenant suffix.
type GetTenantSuffixResult struct {
	BaseIdentityAPIResponse
	Result map[string]interface{} `json:"Result"`
}

// IdpAuthStatusResponse is a struct that represents the response from the Identity API for IdP authentication status.
type IdpAuthStatusResponse struct {
	BaseIdentityAPIResponse
	Result IdpAuthStatusResult `json:"Result"`
}

// TenantEndpointResponse is a struct that represents the response from the Identity API for tenant endpoint.
type TenantEndpointResponse struct {
	Endpoint string `json:"endpoint"`
}
