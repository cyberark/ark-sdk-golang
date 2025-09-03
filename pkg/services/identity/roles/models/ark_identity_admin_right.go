package models

// Possible admin rights.
const (
	AdminPortalLogin           = "/lib/rights/adminportallogin.json"
	ApplicationManagement      = "/lib/rights/appman.json"
	FederationManagement       = "/lib/rights/fedman.json"
	IdentityVerification       = "/lib/rights/identityverification.json"
	MFAUnlock                  = "/lib/rights/mfaunlock.json"
	RadiusManagement           = "/lib/rights/radiusman.json"
	RoleManagement             = "/lib/rights/roleman.json"
	ReadOnlyRoleManagement     = "/lib/rights/roroleman.json"
	ReadOnlyUserManagement     = "/lib/rights/rouserman.json"
	RegisterAndAdminConnectors = "/lib/rights/proxycode.json"
	ReportManagement           = "/lib/rights/reportman.json"
	SharedCredentials          = "/lib/rights/sharedcredentials.json"
	UserManagement             = "/lib/rights/dsman.json"
	ReadOnlySysAdmin           = "/lib/rights/monitor.json"
	SystemEnrollment           = "/lib/rights/agentmanjoin.json"
	VaultManagement            = "/lib/rights/vaultman.json"

	Audit               = "ServiceRight/auditShowTile"
	JIT                 = "ServiceRight/dpaShowTile"
	DPA                 = "ServiceRight/dpaShowTile"
	SIA                 = "ServiceRight/dpaShowTile"
	CloudOnboarding     = "ServiceRight/cloudonboardingShowTile"
	ConnectorManagement = "ServiceRight/connectormanagementShowTile"
	PrivilegeCloud      = "ServiceRight/passwordvaultShowTile"
	SessionMonitoring   = "ServiceRight/smShowTile"
	AccessToResources   = "ServiceRight/accesstoresourcesShowTile"
)
