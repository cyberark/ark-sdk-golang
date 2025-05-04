package services

import "github.com/cyberark/ark-sdk-golang/pkg/models/actions"

// SupportedServiceActions is a list of supported service actions.
var SupportedServiceActions = []*actions.ArkServiceActionDefinition{
	SIAActions,
	CmgrActions,
	PCloudActions,
	IdentityActions,
}
