package uap

import (
	"fmt"

	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/common"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	uap "github.com/cyberark/ark-sdk-golang/pkg/services/uap/common"
	uapcommonmodels "github.com/cyberark/ark-sdk-golang/pkg/services/uap/common/models"
	"github.com/mitchellh/mapstructure"

	"reflect"
)

// ArkUAPPolicyPage represents a page of common policies in the UAP service.
type ArkUAPPolicyPage = common.ArkPage[uapcommonmodels.ArkUAPCommonAccessPolicy]

// ArkUAPService represents the UAP service.
type ArkUAPService struct {
	services.ArkService
	*services.ArkBaseService
	baseService *uap.ArkUAPBaseService
}

// NewArkUAPService creates a new instance of ArkUAPService with the provided authenticators.
func NewArkUAPService(authenticators ...auth.ArkAuth) (*ArkUAPService, error) {
	uapService := &ArkUAPService{}
	var uapServiceInterface services.ArkService = uapService
	baseService, err := services.NewArkBaseService(uapServiceInterface, authenticators...)
	if err != nil {
		return nil, err
	}
	ispBaseAuth, err := baseService.Authenticator("isp")
	if err != nil {
		return nil, err
	}
	ispAuth := ispBaseAuth.(*auth.ArkISPAuth)
	uapService.ArkBaseService = baseService
	uapService.baseService, err = uap.NewArkUAPBaseService(
		ispAuth,
	)
	if err != nil {
		return nil, err
	}
	return uapService, nil
}

// ListPolicies retrieves all policies.
func (s *ArkUAPService) ListPolicies() (<-chan *ArkUAPPolicyPage, error) {
	s.Logger.Info("Listing all policies")
	policyPagesWithType := make(chan *ArkUAPPolicyPage)
	go func() {
		filters := uapcommonmodels.NewArkUAPFilters()
		policyPages, err := s.baseService.BaseListPolicies(filters)
		if err != nil {
			return
		}
		defer close(policyPagesWithType)
		for page := range policyPages {
			policies := ArkUAPPolicyPage{Items: make([]*uapcommonmodels.ArkUAPCommonAccessPolicy, len(page.Items))}
			for idx, policy := range page.Items {
				var commonPolicy uapcommonmodels.ArkUAPCommonAccessPolicy
				err = mapstructure.Decode(*policy, &commonPolicy)
				if err != nil {
					s.Logger.Error("Failed to decode policy page: %v", err)
					continue
				}
				policies.Items[idx] = &commonPolicy
			}
			policyPagesWithType <- &policies
		}
	}()
	return policyPagesWithType, nil
}

// ListPoliciesBy retrieves policies based on the provided filters.
func (s *ArkUAPService) ListPoliciesBy(filters *uapcommonmodels.ArkUAPFilters) (<-chan *ArkUAPPolicyPage, error) {
	s.Logger.Info("Listing policies by filter")
	policyPagesWithType := make(chan *ArkUAPPolicyPage)
	go func() {
		if filters == nil {
			filters = uapcommonmodels.NewArkUAPFilters()
		}
		policyPages, err := s.baseService.BaseListPolicies(filters)
		if err != nil {
			s.Logger.Error("Failed to list policies by filter: %v", err)
			close(policyPagesWithType)
			return
		}
		defer close(policyPagesWithType)
		for page := range policyPages {
			policies := ArkUAPPolicyPage{Items: make([]*uapcommonmodels.ArkUAPCommonAccessPolicy, len(page.Items))}
			for idx, policy := range page.Items {
				var commonPolicy uapcommonmodels.ArkUAPCommonAccessPolicy
				err = mapstructure.Decode(*policy, &commonPolicy)
				if err != nil {
					s.Logger.Error("Failed to decode policy page: %v", err)
					continue
				}
				policies.Items[idx] = &commonPolicy
			}
			policyPagesWithType <- &policies
		}
	}()
	return policyPagesWithType, nil
}

// PolicyStatus retrieves the status of a policy by its ID or name.
func (s *ArkUAPService) PolicyStatus(getPolicyStatus *uapcommonmodels.ArkUAPGetPolicyStatus) (string, error) {
	if getPolicyStatus == nil {
		return "", fmt.Errorf("getPolicyStatus cannot be nil")
	}
	if getPolicyStatus.PolicyID == "" && getPolicyStatus.PolicyName == "" {
		return "", fmt.Errorf("either PolicyID or PolicyName must be provided to retrieve policy status")
	}
	s.Logger.Info("Retrieving policy status for ID [%s] and name [%s]", getPolicyStatus.PolicyID, getPolicyStatus.PolicyName)
	respType := reflect.TypeOf(uapcommonmodels.ArkUAPCommonAccessPolicy{})
	return s.baseService.BasePolicyStatus(getPolicyStatus.PolicyID, getPolicyStatus.PolicyName, &respType)
}

// PoliciesStats retrieves statistics for all policies.
func (s *ArkUAPService) PoliciesStats() (*uapcommonmodels.ArkUAPPoliciesStats, error) {
	s.Logger.Info("Retrieving policies statistics")
	filters := uapcommonmodels.NewArkUAPFilters()
	stats, err := s.baseService.BasePoliciesStats(filters)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve policies statistics: %w", err)
	}
	return stats, nil
}

// ServiceConfig returns the service configuration for ArkUAPSCAService.
func (s *ArkUAPService) ServiceConfig() services.ArkServiceConfig {
	return ServiceConfig
}
