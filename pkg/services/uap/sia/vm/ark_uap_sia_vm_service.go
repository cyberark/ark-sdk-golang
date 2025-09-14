package vm

import (
	"fmt"

	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/common"
	commonmodels "github.com/cyberark/ark-sdk-golang/pkg/models/common"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	uap "github.com/cyberark/ark-sdk-golang/pkg/services/uap/common"
	uapcommonmodels "github.com/cyberark/ark-sdk-golang/pkg/services/uap/common/models"
	uapsiavmmodels "github.com/cyberark/ark-sdk-golang/pkg/services/uap/sia/vm/models"
	"github.com/mitchellh/mapstructure"

	"reflect"
)

// ArkUAPVMPolicyPage represents a page of SIA VM policies in the UAP service.
type ArkUAPVMPolicyPage = common.ArkPage[uapsiavmmodels.ArkUAPSIAVMAccessPolicy]

// ArkUAPSIAVMService represents the UAP SIA VM service.
type ArkUAPSIAVMService struct {
	services.ArkService
	*services.ArkBaseService
	baseService *uap.ArkUAPBaseService
}

// NewArkUAPSIAVMService creates a new instance of ArkUAPSIAVMService with the provided authenticators.
func NewArkUAPSIAVMService(authenticators ...auth.ArkAuth) (*ArkUAPSIAVMService, error) {
	uapSiaVMService := &ArkUAPSIAVMService{}
	var uapSiaVMServiceInterface services.ArkService = uapSiaVMService
	baseService, err := services.NewArkBaseService(uapSiaVMServiceInterface, authenticators...)
	if err != nil {
		return nil, err
	}
	ispBaseAuth, err := baseService.Authenticator("isp")
	if err != nil {
		return nil, err
	}
	ispAuth := ispBaseAuth.(*auth.ArkISPAuth)
	uapSiaVMService.ArkBaseService = baseService
	uapSiaVMService.baseService, err = uap.NewArkUAPBaseService(
		ispAuth,
	)
	if err != nil {
		return nil, err
	}
	return uapSiaVMService, nil
}

// AddPolicy adds a new policy with the given information.
func (s *ArkUAPSIAVMService) AddPolicy(addPolicy *uapsiavmmodels.ArkUAPSIAVMAccessPolicy) (*uapsiavmmodels.ArkUAPSIAVMAccessPolicy, error) {
	s.Logger.Info("Adding new policy [%s]", addPolicy.Metadata.Name)
	addPolicy.Metadata.PolicyEntitlement.TargetCategory = commonmodels.CategoryTypeVM
	if addPolicy.Metadata.PolicyTags == nil {
		addPolicy.Metadata.PolicyTags = make([]string, 0)
	}
	policyType := reflect.TypeOf(addPolicy)
	addPolicySerialized, err := addPolicy.Serialize()
	if err != nil {
		return nil, err
	}
	addPolicyJSON := common.ConvertToCamelCase(addPolicySerialized, &policyType)
	if err != nil {
		return nil, err
	}
	policyResp, err := s.baseService.BaseAddPolicy(addPolicyJSON.(map[string]interface{}))
	if err != nil {
		return nil, err
	}
	return s.Policy(&uapcommonmodels.ArkUAPGetPolicyRequest{
		PolicyID: policyResp.PolicyID,
	})
}

// Policy retrieves a policy by its ID.
func (s *ArkUAPSIAVMService) Policy(policyRequest *uapcommonmodels.ArkUAPGetPolicyRequest) (*uapsiavmmodels.ArkUAPSIAVMAccessPolicy, error) {
	s.Logger.Info("Retrieving policy [%s]", policyRequest.PolicyID)
	respType := reflect.TypeOf(uapsiavmmodels.ArkUAPSIAVMAccessPolicy{})
	policyJSON, err := s.baseService.BasePolicy(policyRequest.PolicyID, &respType)
	if err != nil {
		return nil, err
	}
	policyJSONSnake := common.ConvertToSnakeCase(policyJSON, &respType)
	var vmPolicy uapsiavmmodels.ArkUAPSIAVMAccessPolicy
	err = vmPolicy.Deserialize(policyJSONSnake.(map[string]interface{}))
	if err != nil {
		return nil, err
	}
	return &vmPolicy, nil
}

// UpdatePolicy edits an existing policy with the given information.
func (s *ArkUAPSIAVMService) UpdatePolicy(updatePolicy *uapsiavmmodels.ArkUAPSIAVMAccessPolicy) (*uapsiavmmodels.ArkUAPSIAVMAccessPolicy, error) {
	s.Logger.Info("Updating policy [%s]", updatePolicy.Metadata.PolicyID)
	policyType := reflect.TypeOf(uapsiavmmodels.ArkUAPSIAVMAccessPolicy{})
	updatePolicySerialized, err := updatePolicy.Serialize()
	if err != nil {
		return nil, err
	}
	updatePolicyJSON := common.ConvertToCamelCase(updatePolicySerialized, &policyType)
	if err != nil {
		return nil, err
	}
	err = s.baseService.BaseUpdatePolicy(updatePolicy.Metadata.PolicyID, updatePolicyJSON.(map[string]interface{}))
	if err != nil {
		return nil, err
	}
	return s.Policy(&uapcommonmodels.ArkUAPGetPolicyRequest{
		PolicyID: updatePolicy.Metadata.PolicyID,
	})
}

// ListPolicies retrieves all policies.
func (s *ArkUAPSIAVMService) ListPolicies() (<-chan *ArkUAPVMPolicyPage, error) {
	s.Logger.Info("Listing all policies")
	policyPagesWithType := make(chan *ArkUAPVMPolicyPage)
	go func() {
		filters := uapcommonmodels.NewArkUAPFilters()
		filters.TargetCategory = []string{commonmodels.CategoryTypeVM}
		policyPages, err := s.baseService.BaseListPolicies(filters)
		if err != nil {
			return
		}
		defer close(policyPagesWithType)
		for page := range policyPages {
			vmPolicies := ArkUAPVMPolicyPage{Items: make([]*uapsiavmmodels.ArkUAPSIAVMAccessPolicy, len(page.Items))}
			for idx, policy := range page.Items {
				var vmPolicy uapsiavmmodels.ArkUAPSIAVMAccessPolicy
				err = mapstructure.Decode(*policy, &vmPolicy)
				if err != nil {
					s.Logger.Error("Failed to decode policy page: %v", err)
					continue
				}
				vmPolicies.Items[idx] = &vmPolicy
			}
			policyPagesWithType <- &vmPolicies
		}
	}()
	return policyPagesWithType, nil
}

// ListPoliciesBy retrieves policies based on the provided filters.
func (s *ArkUAPSIAVMService) ListPoliciesBy(filters *uapsiavmmodels.ArkUAPSIAVMFilters) (<-chan *ArkUAPVMPolicyPage, error) {
	s.Logger.Info("Listing policies by filter")
	policyPagesWithType := make(chan *ArkUAPVMPolicyPage)
	go func() {
		if filters == nil {
			filters = &uapsiavmmodels.ArkUAPSIAVMFilters{
				ArkUAPFilters: *uapcommonmodels.NewArkUAPFilters(),
			}
		}
		filters.TargetCategory = []string{commonmodels.CategoryTypeVM}
		policyPages, err := s.baseService.BaseListPolicies(&filters.ArkUAPFilters)
		if err != nil {
			s.Logger.Error("Failed to list policies by filter: %v", err)
			close(policyPagesWithType)
			return
		}
		defer close(policyPagesWithType)
		for page := range policyPages {
			vmPolicies := ArkUAPVMPolicyPage{Items: make([]*uapsiavmmodels.ArkUAPSIAVMAccessPolicy, len(page.Items))}
			for idx, policy := range page.Items {
				var vmPolicy uapsiavmmodels.ArkUAPSIAVMAccessPolicy
				err = mapstructure.Decode(*policy, &vmPolicy)
				if err != nil {
					s.Logger.Error("Failed to decode policy page: %v", err)
					continue
				}
				vmPolicies.Items[idx] = &vmPolicy
			}
			policyPagesWithType <- &vmPolicies
		}
	}()
	return policyPagesWithType, nil
}

// DeletePolicy deletes a policy by its ID.
func (s *ArkUAPSIAVMService) DeletePolicy(deletePolicy *uapcommonmodels.ArkUAPDeletePolicyRequest) error {
	s.Logger.Info("Deleting policy [%s]", deletePolicy.PolicyID)
	return s.baseService.BaseDeletePolicy(deletePolicy.PolicyID)
}

// PolicyStatus retrieves the status of a policy by its ID or name.
func (s *ArkUAPSIAVMService) PolicyStatus(getPolicyStatus *uapcommonmodels.ArkUAPGetPolicyStatus) (string, error) {
	if getPolicyStatus == nil {
		return "", fmt.Errorf("getPolicyStatus cannot be nil")
	}
	if getPolicyStatus.PolicyID == "" && getPolicyStatus.PolicyName == "" {
		return "", fmt.Errorf("either PolicyID or PolicyName must be provided to retrieve policy status")
	}
	s.Logger.Info("Retrieving policy status for ID [%s] and name [%s]", getPolicyStatus.PolicyID, getPolicyStatus.PolicyName)
	respType := reflect.TypeOf(uapsiavmmodels.ArkUAPSIAVMAccessPolicy{})
	return s.baseService.BasePolicyStatus(getPolicyStatus.PolicyID, getPolicyStatus.PolicyName, &respType)
}

// PoliciesStats calculates policies statistics.
func (s *ArkUAPSIAVMService) PoliciesStats() (*uapcommonmodels.ArkUAPPoliciesStats, error) {
	s.Logger.Info("Calculating policies statistics")
	filters := uapcommonmodels.NewArkUAPFilters()
	filters.TargetCategory = []string{commonmodels.CategoryTypeVM}
	return s.baseService.BasePoliciesStats(filters)
}

// ServiceConfig returns the service configuration for ArkUAPSIAVMService.
func (s *ArkUAPSIAVMService) ServiceConfig() services.ArkServiceConfig {
	return ServiceConfig
}
