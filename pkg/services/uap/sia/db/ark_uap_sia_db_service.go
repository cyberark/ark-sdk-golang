package db

import (
	"fmt"
	"reflect"

	"github.com/cyberark/ark-sdk-golang/pkg/auth"
	"github.com/cyberark/ark-sdk-golang/pkg/common"
	commonmodels "github.com/cyberark/ark-sdk-golang/pkg/models/common"
	"github.com/cyberark/ark-sdk-golang/pkg/services"
	uap "github.com/cyberark/ark-sdk-golang/pkg/services/uap/common"
	uapcommonmodels "github.com/cyberark/ark-sdk-golang/pkg/services/uap/common/models"
	uapsiadbmodels "github.com/cyberark/ark-sdk-golang/pkg/services/uap/sia/db/models"
	"github.com/mitchellh/mapstructure"
)

// ArkUAPDBPolicyPage represents a page of SIA DB policies in the UAP service.
type ArkUAPDBPolicyPage = common.ArkPage[uapsiadbmodels.ArkUAPSIADBAccessPolicy]

// ArkUAPSIADBServiceConfig defines the service configuration for ArkUAPSIADBService.
var ArkUAPSIADBServiceConfig = services.ArkServiceConfig{
	ServiceName:                "uap-db",
	RequiredAuthenticatorNames: []string{"isp"},
	OptionalAuthenticatorNames: []string{},
}

// ArkUAPSIADBService represents the UAP SIA DB service.
type ArkUAPSIADBService struct {
	services.ArkService
	*services.ArkBaseService
	baseService *uap.ArkUAPBaseService
}

// NewArkUAPSIADBService creates a new instance of ArkUAPSIADBService with the provided authenticators.
func NewArkUAPSIADBService(authenticators ...auth.ArkAuth) (*ArkUAPSIADBService, error) {
	uapSiaDbService := &ArkUAPSIADBService{}
	var uapSiaDbServiceInterface services.ArkService = uapSiaDbService
	baseService, err := services.NewArkBaseService(uapSiaDbServiceInterface, authenticators...)
	if err != nil {
		return nil, err
	}
	ispBaseAuth, err := baseService.Authenticator("isp")
	if err != nil {
		return nil, err
	}
	ispAuth := ispBaseAuth.(*auth.ArkISPAuth)
	uapSiaDbService.ArkBaseService = baseService
	uapSiaDbService.baseService, err = uap.NewArkUAPBaseService(
		ispAuth,
	)
	if err != nil {
		return nil, err
	}
	return uapSiaDbService, nil
}

func (s *ArkUAPSIADBService) serializeProfile(policy *uapsiadbmodels.ArkUAPSIADBAccessPolicy, policyJSON map[string]interface{}) error {
	// Fill the profiles for the instances
	var err error
	for name := range policy.Targets {
		for idx := range policy.Targets[name].Instances {
			instanceJSON := policyJSON["targets"].(map[string]interface{})[name].(map[string]interface{})["instances"].([]interface{})[idx].(map[string]interface{})
			policy.Targets[name].Instances[idx].ClearProfileFromData(instanceJSON)
			instanceJSON["profile"], err = policy.Targets[name].Instances[idx].SerializeProfile()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *ArkUAPSIADBService) deserializeProfile(policy *uapsiadbmodels.ArkUAPSIADBAccessPolicy, policyJSON map[string]interface{}) error {
	// Fill the profiles for the instances
	var err error
	for name := range policy.Targets {
		for idx := range policy.Targets[name].Instances {
			instanceJSON := policyJSON["targets"].(map[string]interface{})[name].(map[string]interface{})["instances"].([]interface{})[idx].(map[string]interface{})
			if instanceJSON["profile"] != nil {
				err = policy.Targets[name].Instances[idx].DeserializeProfile(instanceJSON["profile"].(map[string]interface{}))
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// AddPolicy adds a new policy with the given information.
func (s *ArkUAPSIADBService) AddPolicy(addPolicy *uapsiadbmodels.ArkUAPSIADBAccessPolicy) (*uapsiadbmodels.ArkUAPSIADBAccessPolicy, error) {
	s.Logger.Info("Adding new policy [%s]", addPolicy.Metadata.Name)
	addPolicy.Metadata.PolicyEntitlement.TargetCategory = commonmodels.CategoryTypeDB
	if addPolicy.Metadata.PolicyTags == nil {
		addPolicy.Metadata.PolicyTags = make([]string, 0)
	}
	policyType := reflect.TypeOf(addPolicy)
	policyJSON, err := common.SerializeJSONCamelSchema(addPolicy, &policyType)
	if err != nil {
		return nil, err
	}
	err = s.serializeProfile(addPolicy, policyJSON)
	if err != nil {
		return nil, err
	}
	policyResp, err := s.baseService.BaseAddPolicy(policyJSON)
	if err != nil {
		return nil, err
	}
	return s.Policy(&uapcommonmodels.ArkUAPGetPolicyRequest{
		PolicyID: policyResp.PolicyID,
	})
}

// Policy retrieves a policy by its ID.
func (s *ArkUAPSIADBService) Policy(policyRequest *uapcommonmodels.ArkUAPGetPolicyRequest) (*uapsiadbmodels.ArkUAPSIADBAccessPolicy, error) {
	s.Logger.Info("Retrieving policy [%s]", policyRequest.PolicyID)
	respType := reflect.TypeOf(uapsiadbmodels.ArkUAPSIADBAccessPolicy{})
	policyJSON, err := s.baseService.BasePolicy(policyRequest.PolicyID, &respType)
	if err != nil {
		return nil, err
	}
	var dbPolicy uapsiadbmodels.ArkUAPSIADBAccessPolicy
	err = mapstructure.Decode(policyJSON, &dbPolicy)
	if err != nil {
		return nil, err
	}
	err = s.deserializeProfile(&dbPolicy, policyJSON)
	if err != nil {
		return nil, err
	}
	return &dbPolicy, nil
}

// UpdatePolicy edits an existing policy with the given information.
func (s *ArkUAPSIADBService) UpdatePolicy(updatePolicy *uapsiadbmodels.ArkUAPSIADBAccessPolicy) (*uapsiadbmodels.ArkUAPSIADBAccessPolicy, error) {
	s.Logger.Info("Updating policy [%s]", updatePolicy.Metadata.PolicyID)
	policyType := reflect.TypeOf(uapsiadbmodels.ArkUAPSIADBAccessPolicy{})
	policyJSON, err := common.SerializeJSONCamelSchema(updatePolicy, &policyType)
	if err != nil {
		return nil, err
	}
	err = s.serializeProfile(updatePolicy, policyJSON)
	if err != nil {
		return nil, err
	}
	err = s.baseService.BaseUpdatePolicy(updatePolicy.Metadata.PolicyID, policyJSON)
	if err != nil {
		return nil, err
	}
	return s.Policy(&uapcommonmodels.ArkUAPGetPolicyRequest{
		PolicyID: updatePolicy.Metadata.PolicyID,
	})
}

// ListPolicies retrieves all policies.
func (s *ArkUAPSIADBService) ListPolicies() (<-chan *ArkUAPDBPolicyPage, error) {
	s.Logger.Info("Listing all policies")
	policyPagesWithType := make(chan *ArkUAPDBPolicyPage)
	go func() {
		filters := uapcommonmodels.NewArkUAPFilters()
		filters.TargetCategory = []string{commonmodels.CategoryTypeDB}
		policyPages, err := s.baseService.BaseListPolicies(filters)
		if err != nil {
			return
		}
		defer close(policyPagesWithType)
		for page := range policyPages {
			dbPolicies := ArkUAPDBPolicyPage{Items: make([]*uapsiadbmodels.ArkUAPSIADBAccessPolicy, len(page.Items))}
			for idx, policy := range page.Items {
				var dbPolicy uapsiadbmodels.ArkUAPSIADBAccessPolicy
				err = mapstructure.Decode(*policy, &dbPolicy)
				if err != nil {
					s.Logger.Error("Failed to decode policy page: %v", err)
					continue
				}
				dbPolicies.Items[idx] = &dbPolicy
			}
			policyPagesWithType <- &dbPolicies
		}
	}()
	return policyPagesWithType, nil
}

// ListPoliciesBy retrieves policies based on the provided filters.
func (s *ArkUAPSIADBService) ListPoliciesBy(filters *uapsiadbmodels.ArkUAPSIADBFilters) (<-chan *ArkUAPDBPolicyPage, error) {
	s.Logger.Info("Listing policies by filter")
	policyPagesWithType := make(chan *ArkUAPDBPolicyPage)
	go func() {
		if filters == nil {
			filters = &uapsiadbmodels.ArkUAPSIADBFilters{
				ArkUAPFilters: *uapcommonmodels.NewArkUAPFilters(),
			}
		}
		filters.TargetCategory = []string{commonmodels.CategoryTypeDB}
		policyPages, err := s.baseService.BaseListPolicies(&filters.ArkUAPFilters)
		if err != nil {
			s.Logger.Error("Failed to list policies by filter: %v", err)
			close(policyPagesWithType)
			return
		}
		defer close(policyPagesWithType)
		for page := range policyPages {
			dbPolicies := ArkUAPDBPolicyPage{Items: make([]*uapsiadbmodels.ArkUAPSIADBAccessPolicy, len(page.Items))}
			for idx, policy := range page.Items {
				var dbPolicy uapsiadbmodels.ArkUAPSIADBAccessPolicy
				err = mapstructure.Decode(*policy, &dbPolicy)
				if err != nil {
					s.Logger.Error("Failed to decode policy page: %v", err)
					continue
				}
				dbPolicies.Items[idx] = &dbPolicy
			}
			policyPagesWithType <- &dbPolicies
		}
	}()
	return policyPagesWithType, nil
}

// DeletePolicy deletes a policy by its ID.
func (s *ArkUAPSIADBService) DeletePolicy(deletePolicy *uapcommonmodels.ArkUAPDeletePolicyRequest) error {
	s.Logger.Info("Deleting policy [%s]", deletePolicy.PolicyID)
	return s.baseService.BaseDeletePolicy(deletePolicy.PolicyID)
}

// PolicyStatus retrieves the status of a policy by its ID or name.
func (s *ArkUAPSIADBService) PolicyStatus(getPolicyStatus *uapcommonmodels.ArkUAPGetPolicyStatus) (string, error) {
	if getPolicyStatus == nil {
		return "", fmt.Errorf("getPolicyStatus cannot be nil")
	}
	if getPolicyStatus.PolicyID == "" && getPolicyStatus.PolicyName == "" {
		return "", fmt.Errorf("either PolicyID or PolicyName must be provided to retrieve policy status")
	}
	s.Logger.Info("Retrieving policy status for ID [%s] and name [%s]", getPolicyStatus.PolicyID, getPolicyStatus.PolicyName)
	respType := reflect.TypeOf(uapsiadbmodels.ArkUAPSIADBAccessPolicy{})
	return s.baseService.BasePolicyStatus(getPolicyStatus.PolicyID, getPolicyStatus.PolicyName, &respType)
}

// PoliciesStats calculates policies statistics.
func (s *ArkUAPSIADBService) PoliciesStats() (*uapcommonmodels.ArkUAPPoliciesStats, error) {
	s.Logger.Info("Calculating policies statistics")
	filters := uapcommonmodels.NewArkUAPFilters()
	filters.TargetCategory = []string{commonmodels.CategoryTypeDB}
	return s.baseService.BasePoliciesStats(filters)
}

// ServiceConfig returns the service configuration for ArkUAPSIADBService.
func (s *ArkUAPSIADBService) ServiceConfig() services.ArkServiceConfig {
	return ArkUAPSIADBServiceConfig
}
