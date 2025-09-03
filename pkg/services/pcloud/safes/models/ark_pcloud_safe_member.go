package models

// Possible safe member types
const (
	User  = "User"
	Group = "Group"
	Role  = "Role"
)

// Possible safe member permission sets
const (
	ConnectOnly     = "connect_only"
	ReadOnly        = "read_only"
	Approver        = "approver"
	AccountsManager = "accounts_manager"
	Full            = "full"
	Custom          = "custom"
)

// ArkPCloudSafeMemberPermissions represents the permissions of a safe member.
type ArkPCloudSafeMemberPermissions struct {
	UseAccounts                            bool `json:"use_accounts" mapstructure:"use_accounts" desc:"Use accounts permission" default:"false"`
	RetrieveAccounts                       bool `json:"retrieve_accounts" mapstructure:"retrieve_accounts" desc:"Retrieve accounts permission" default:"false"`
	ListAccounts                           bool `json:"list_accounts" mapstructure:"list_accounts" desc:"List accounts permission" default:"false"`
	AddAccounts                            bool `json:"add_accounts" mapstructure:"add_accounts" desc:"Add accounts permission" default:"false"`
	UpdateAccountContent                   bool `json:"update_account_content" mapstructure:"update_account_content" desc:"Update account content permission" default:"false"`
	UpdateAccountProperties                bool `json:"update_account_properties" mapstructure:"update_account_properties" desc:"Update account properties permission" default:"false"`
	InitiateCPMAccountManagementOperations bool `json:"initiate_cpm_account_management_operations" mapstructure:"initiate_cpm_account_management_operations" desc:"Initiate CPM account management operations permission" default:"false"`
	SpecifyNextAccountContent              bool `json:"specify_next_account_content" mapstructure:"specify_next_account_content" desc:"Specify next account content permissions" default:"false"`
	RenameAccounts                         bool `json:"rename_accounts" mapstructure:"rename_accounts" desc:"Rename accounts permission" default:"false"`
	DeleteAccounts                         bool `json:"delete_accounts" mapstructure:"delete_accounts" desc:"Delete accounts permission" default:"false"`
	UnlockAccounts                         bool `json:"unlock_accounts" mapstructure:"unlock_accounts" desc:"Unlock accounts permission" default:"false"`
	ManageSafe                             bool `json:"manage_safe" mapstructure:"manage_safe" desc:"Manage safe permission" default:"false"`
	ManageSafeMembers                      bool `json:"manage_safe_members" mapstructure:"manage_safe_members" desc:"Manage safe members" default:"false"`
	BackupSafe                             bool `json:"backup_safe" mapstructure:"backup_safe" desc:"Backup safe permission" default:"false"`
	ViewAuditLog                           bool `json:"view_audit_log" mapstructure:"view_audit_log" desc:"View audit log permission" default:"false"`
	ViewSafeMembers                        bool `json:"view_safe_members" mapstructure:"view_safe_members" desc:"View safe members permission" default:"false"`
	AccessWithoutConfirmation              bool `json:"access_without_confirmation" mapstructure:"access_without_confirmation" desc:"Access without confirmation permission" default:"false"`
	CreateFolders                          bool `json:"create_folders" mapstructure:"create_folders" desc:"Create folders permission" default:"false"`
	DeleteFolders                          bool `json:"delete_folders" mapstructure:"delete_folders" desc:"Delete folders permission" default:"false"`
	MoveAccountsAndFolders                 bool `json:"move_accounts_and_folders" mapstructure:"move_accounts_and_folders" desc:"Move accounts and folders permission" default:"false"`
	RequestsAuthorizationLevel1            bool `json:"requests_authorization_level1" mapstructure:"requests_authorization_level1" desc:"Request authorization level 1 permission" default:"false"`
	RequestsAuthorizationLevel2            bool `json:"requests_authorization_level2" mapstructure:"requests_authorization_level2" desc:"Request authorization level 2 permission" default:"false"`
}

// ArkPCloudSafeMember represents a safe member with its details and permissions.
type ArkPCloudSafeMember struct {
	SafeID                     string                         `json:"safe_id" mapstructure:"safe_id" desc:"Safe url identifier"`
	SafeName                   string                         `json:"safe_name" mapstructure:"safe_name" desc:"Name of the safe of the member"`
	SafeNumber                 int                            `json:"safe_number" mapstructure:"safe_number" desc:"Number id of the safe"`
	MemberID                   interface{}                    `json:"member_id" mapstructure:"member_id" desc:"Member id"`
	MemberName                 string                         `json:"member_name" mapstructure:"member_name" desc:"Name of the member of the safe"`
	MemberType                 string                         `json:"member_type" mapstructure:"member_type" desc:"Type of the member of the safe (User,Group,Role)" choices:"User,Group,Role"`
	MembershipExpirationDate   int                            `json:"membership_expiration_date,omitempty" mapstructure:"membership_expiration_date" desc:"Expiration date of the member on the safe"`
	IsExpiredMembershipEnabled bool                           `json:"is_expired_membership_enabled,omitempty" mapstructure:"is_expired_membership_enabled" desc:"Whether expired membership is enabled or not"`
	IsPredefinedUser           bool                           `json:"is_predefined_user" mapstructure:"is_predefined_user" desc:"Whether this is a predefined user or not"`
	IsReadOnly                 bool                           `json:"is_read_only" mapstructure:"is_read_only" desc:"Whether this member is read only"`
	Permissions                ArkPCloudSafeMemberPermissions `json:"permissions" mapstructure:"permissions" desc:"Permissions of the safe member"`
	PermissionSet              string                         `json:"permission_set" mapstructure:"permission_set" desc:"Permission set type the permissions are set to (connect_only,read_only,approver,accounts_manager,full,custom)" default:"custom" choices:"connect_only,read_only,approver,accounts_manager,full,custom"`
}
