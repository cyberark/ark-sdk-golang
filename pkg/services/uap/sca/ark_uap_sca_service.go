package sca

import (
	"fmt"
	"reflect"

	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/common"
	commonmodels "github.com/cyberark/ark-sdk-golang/pkg/models/common"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	uap "github.com/cyberark/ark-sdk-golang/pkg/services/uap/common"
	uapcommonmodels "github.com/cyberark/ark-sdk-golang/pkg/services/uap/common/models"
	uapscamodels "github.com/cyberark/ark-sdk-golang/pkg/services/uap/sca/models"
	"github.com/mitchellh/mapstructure"
)

const (
	policyStatusActiveRetryCount = 10
)

// ArkUAPSCAPolicyPage represents a page of SCA policies in the UAP service.
type ArkUAPSCAPolicyPage = common.ArkPage[uapscamodels.ArkUAPSCACloudConsoleAccessPolicy]

// ArkUAPSCAServiceConfig defines the service configuration for ArkUAPSCAService.
var ArkUAPSCAServiceConfig = services.ArkServiceConfig{
	ServiceName:                "uap-sca",
	RequiredAuthenticatorNames: []string{"isp"},
	OptionalAuthenticatorNames: []string{},
}

// ArkUAPSCAService represents the UAP SCA service.
type ArkUAPSCAService struct {
	services.ArkService
	*services.ArkBaseService
	baseService *uap.ArkUAPBaseService
}

// NewArkUAPSCAService creates a new instance of ArkUAPSCAService with the provided authenticators.
func NewArkUAPSCAService(authenticators ...auth.ArkAuth) (*ArkUAPSCAService, error) {
	uapScaService := &ArkUAPSCAService{}
	var uapScaServiceInterface services.ArkService = uapScaService
	baseService, err := services.NewArkBaseService(uapScaServiceInterface, authenticators...)
	if err != nil {
		return nil, err
	}
	ispBaseAuth, err := baseService.Authenticator("isp")
	if err != nil {
		return nil, err
	}
	ispAuth := ispBaseAuth.(*auth.ArkISPAuth)
	uapScaService.ArkBaseService = baseService
	uapScaService.baseService, err = uap.NewArkUAPBaseService(
		ispAuth,
	)
	if err != nil {
		return nil, err
	}
	return uapScaService, nil
}

func (s *ArkUAPSCAService) serializeTargets(policy *uapscamodels.ArkUAPSCACloudConsoleAccessPolicy, policyJSON map[string]interface{}) error {
	var err error
	policy.Targets.ClearTargetsFromData(policyJSON["targets"].(map[string]interface{}))
	policyJSON["targets"], err = policy.Targets.SerializeTargets()
	return err
}

func (s *ArkUAPSCAService) deserializeTargets(policy *uapscamodels.ArkUAPSCACloudConsoleAccessPolicy, policyJSON map[string]interface{}) error {
	return policy.Targets.DeserializeTargets(policyJSON["targets"].(map[string]interface{}))
}

// AddPolicy adds a new policy with the given information.
func (s *ArkUAPSCAService) AddPolicy(addPolicy *uapscamodels.ArkUAPSCACloudConsoleAccessPolicy) (*uapscamodels.ArkUAPSCACloudConsoleAccessPolicy, error) {
	s.Logger.Info("Adding new policy [%s]", addPolicy.Metadata.Name)
	addPolicy.Metadata.PolicyEntitlement.TargetCategory = commonmodels.CategoryTypeCloudConsole
	if addPolicy.Metadata.PolicyTags == nil {
		addPolicy.Metadata.PolicyTags = make([]string, 0)
	}
	policyJSON, err := common.SerializeJSONCamel(addPolicy)
	if err != nil {
		return nil, err
	}
	err = s.serializeTargets(addPolicy, policyJSON)
	if err != nil {
		return nil, err
	}
	policyResp, err := s.baseService.BaseAddPolicy(policyJSON)
	if err != nil {
		return nil, err
	}
	retryCount := 0
	for {
		policy, err := s.Policy(&uapcommonmodels.ArkUAPGetPolicyRequest{
			PolicyID: policyResp.PolicyID,
		})
		if err != nil {
			return nil, err
		}
		if policy.Metadata.Status.Status == uapcommonmodels.StatusTypeActive {
			break
		}
		if policy.Metadata.Status.Status == uapcommonmodels.StatusTypeError {
			return nil, fmt.Errorf("policy [%s] is in error state: %s", policyResp.PolicyID, policy.Metadata.Status.Status)
		}
		if retryCount >= policyStatusActiveRetryCount {
			s.Logger.Warning("Policy [%s] is not active after 10 retries, "+
				"might indicate an issue, moving on regardless", policyResp.PolicyID)
			break
		}
		retryCount++
	}
	return s.Policy(&uapcommonmodels.ArkUAPGetPolicyRequest{
		PolicyID: policyResp.PolicyID,
	})
}

// Policy retrieves a policy by its ID.
func (s *ArkUAPSCAService) Policy(policyRequest *uapcommonmodels.ArkUAPGetPolicyRequest) (*uapscamodels.ArkUAPSCACloudConsoleAccessPolicy, error) {
	s.Logger.Info("Retrieving policy [%s]", policyRequest.PolicyID)
	respType := reflect.TypeOf(uapscamodels.ArkUAPSCACloudConsoleAccessPolicy{})
	policyJSON, err := s.baseService.BasePolicy(policyRequest.PolicyID, &respType)
	if err != nil {
		return nil, err
	}
	var scaPolicy uapscamodels.ArkUAPSCACloudConsoleAccessPolicy
	err = mapstructure.Decode(policyJSON, &scaPolicy)
	if err != nil {
		return nil, err
	}
	err = s.deserializeTargets(&scaPolicy, policyJSON)
	if err != nil {
		return nil, err
	}
	return &scaPolicy, nil
}

// UpdatePolicy edits an existing policy with the given information.
func (s *ArkUAPSCAService) UpdatePolicy(updatePolicy *uapscamodels.ArkUAPSCACloudConsoleAccessPolicy) (*uapscamodels.ArkUAPSCACloudConsoleAccessPolicy, error) {
	s.Logger.Info("Updating policy [%s]", updatePolicy.Metadata.PolicyID)
	policyJSON, err := common.SerializeJSONCamel(updatePolicy)
	if err != nil {
		return nil, err
	}
	err = s.serializeTargets(updatePolicy, policyJSON)
	if err != nil {
		return nil, err
	}
	err = s.baseService.BaseUpdatePolicy(updatePolicy.Metadata.PolicyID, policyJSON)
	if err != nil {
		return nil, err
	}
	retryCount := 0
	for {
		policy, err := s.Policy(&uapcommonmodels.ArkUAPGetPolicyRequest{
			PolicyID: updatePolicy.Metadata.PolicyID,
		})
		if err != nil {
			return nil, err
		}
		if policy.Metadata.Status.Status == uapcommonmodels.StatusTypeActive {
			break
		}
		if policy.Metadata.Status.Status == uapcommonmodels.StatusTypeError {
			return nil, fmt.Errorf("policy [%s] is in error state: %s", updatePolicy.Metadata.PolicyID, policy.Metadata.Status.Status)
		}
		if retryCount >= policyStatusActiveRetryCount {
			s.Logger.Warning("Policy [%s] is not active after 10 retries, "+
				"might indicate an issue, moving on regardless", updatePolicy.Metadata.PolicyID)
			break
		}
		retryCount++
	}
	return s.Policy(&uapcommonmodels.ArkUAPGetPolicyRequest{
		PolicyID: updatePolicy.Metadata.PolicyID,
	})
}

// ListPolicies retrieves all policies.
func (s *ArkUAPSCAService) ListPolicies() (<-chan *ArkUAPSCAPolicyPage, error) {
	s.Logger.Info("Listing all policies")
	policyPagesWithType := make(chan *ArkUAPSCAPolicyPage)
	go func() {
		filters := uapcommonmodels.NewArkUAPFilters()
		filters.TargetCategory = []string{commonmodels.CategoryTypeCloudConsole}
		policyPages, err := s.baseService.BaseListPolicies(filters)
		if err != nil {
			return
		}
		defer close(policyPagesWithType)
		for page := range policyPages {
			scaPolicies := ArkUAPSCAPolicyPage{Items: make([]*uapscamodels.ArkUAPSCACloudConsoleAccessPolicy, len(page.Items))}
			for idx, policy := range page.Items {
				var scaPolicy uapscamodels.ArkUAPSCACloudConsoleAccessPolicy
				err = mapstructure.Decode(*policy, &scaPolicy)
				if err != nil {
					s.Logger.Error("Failed to decode policy page: %v", err)
					continue
				}
				scaPolicies.Items[idx] = &scaPolicy
			}
			policyPagesWithType <- &scaPolicies
		}
	}()
	return policyPagesWithType, nil
}

// ListPoliciesBy retrieves policies based on the provided filters.
func (s *ArkUAPSCAService) ListPoliciesBy(filters *uapscamodels.ArkUAPSCAFilters) (<-chan *ArkUAPSCAPolicyPage, error) {
	s.Logger.Info("Listing policies by filter")
	policyPagesWithType := make(chan *ArkUAPSCAPolicyPage)
	go func() {
		if filters == nil {
			filters = &uapscamodels.ArkUAPSCAFilters{
				ArkUAPFilters: *uapcommonmodels.NewArkUAPFilters(),
			}
		}
		filters.TargetCategory = []string{commonmodels.CategoryTypeCloudConsole}
		policyPages, err := s.baseService.BaseListPolicies(&filters.ArkUAPFilters)
		if err != nil {
			s.Logger.Error("Failed to list policies by filter: %v", err)
			close(policyPagesWithType)
			return
		}
		defer close(policyPagesWithType)
		for page := range policyPages {
			scaPolicies := ArkUAPSCAPolicyPage{Items: make([]*uapscamodels.ArkUAPSCACloudConsoleAccessPolicy, len(page.Items))}
			for idx, policy := range page.Items {
				var scaPolicy uapscamodels.ArkUAPSCACloudConsoleAccessPolicy
				err = mapstructure.Decode(*policy, &scaPolicy)
				if err != nil {
					s.Logger.Error("Failed to decode policy page: %v", err)
					continue
				}
				scaPolicies.Items[idx] = &scaPolicy
			}
			policyPagesWithType <- &scaPolicies
		}
	}()
	return policyPagesWithType, nil
}

// DeletePolicy deletes a policy by its ID.
func (s *ArkUAPSCAService) DeletePolicy(deletePolicy *uapcommonmodels.ArkUAPDeletePolicyRequest) error {
	s.Logger.Info("Deleting policy [%s]", deletePolicy.PolicyID)
	return s.baseService.BaseDeletePolicy(deletePolicy.PolicyID)
}

// PolicyStatus retrieves the status of a policy by its ID or name.
func (s *ArkUAPSCAService) PolicyStatus(getPolicyStatus *uapcommonmodels.ArkUAPGetPolicyStatus) (string, error) {
	if getPolicyStatus == nil {
		return "", fmt.Errorf("getPolicyStatus cannot be nil")
	}
	if getPolicyStatus.PolicyID == "" && getPolicyStatus.PolicyName == "" {
		return "", fmt.Errorf("either PolicyID or PolicyName must be provided to retrieve policy status")
	}
	s.Logger.Info("Retrieving policy status for ID [%s] and name [%s]", getPolicyStatus.PolicyID, getPolicyStatus.PolicyName)
	respType := reflect.TypeOf(uapscamodels.ArkUAPSCACloudConsoleAccessPolicy{})
	return s.baseService.BasePolicyStatus(getPolicyStatus.PolicyID, getPolicyStatus.PolicyName, &respType)
}

// PoliciesStats calculates policies statistics.
func (s *ArkUAPSCAService) PoliciesStats() (*uapcommonmodels.ArkUAPPoliciesStats, error) {
	s.Logger.Info("Calculating policies statistics")
	filters := uapcommonmodels.NewArkUAPFilters()
	filters.TargetCategory = []string{commonmodels.CategoryTypeCloudConsole}
	return s.baseService.BasePoliciesStats(filters)
}

// ServiceConfig returns the service configuration for ArkUAPSCAService.
func (s *ArkUAPSCAService) ServiceConfig() services.ArkServiceConfig {
	return ArkUAPSCAServiceConfig
}
