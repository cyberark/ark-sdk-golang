package models

import (
	"errors"
	"strings"
)

// Constants for validation
const (
	DBProfileMaximumEntities     = 50
	LDAPGroupMaxNameLength       = 50
	DBRoleMaxLength              = 50
	DBUserMaxLength              = 256
	DatabaseNameMaxLength        = 256
	UAPSiaDBTargetsMaxItemsCount = 1000
)

// ArkUAPSIAAuthProfile is an interface for authentication profiles in UAP SIA DB.
type ArkUAPSIAAuthProfile interface {
	Serialize() map[string]interface{}
	Deserialize(data map[string]interface{}) error
}

// ArkUAPSIADBLDAPAuthProfile represents an LDAP authentication profile.
type ArkUAPSIADBLDAPAuthProfile struct {
	AssignGroups []string `json:"assign_groups" validate:"min=1,max=50" mapstructure:"assign_groups" flag:"assign-groups" desc:"List of groups to assign to the user"`
}

// Serialize converts the LDAP authentication profile to a map.
func (s *ArkUAPSIADBLDAPAuthProfile) Serialize() map[string]interface{} {
	return map[string]interface{}{
		"assignGroups": s.AssignGroups,
	}
}

// Deserialize populates the LDAP authentication profile from a map.
func (s *ArkUAPSIADBLDAPAuthProfile) Deserialize(data map[string]interface{}) error {
	s.AssignGroups = []string{}
	if assignGroups, ok := data["assign_groups"].([]interface{}); ok {
		for _, group := range assignGroups {
			if groupStr, ok := group.(string); ok && len(groupStr) <= LDAPGroupMaxNameLength {
				s.AssignGroups = append(s.AssignGroups, groupStr)
			} else {
				return errors.New("invalid group name in assign_groups")
			}
		}
	} else {
		return errors.New("assign_groups must be an array of strings")
	}
	return nil
}

// ArkUAPSIADBDBAuthProfile represents a local DB authentication profile.
type ArkUAPSIADBDBAuthProfile struct {
	Roles []string `json:"roles" validate:"min=1,max=50" mapstructure:"roles" flag:"roles" desc:"List of roles assigned to the user"`
}

// Serialize converts the local DB authentication profile to a map.
func (s *ArkUAPSIADBDBAuthProfile) Serialize() map[string]interface{} {
	return map[string]interface{}{
		"roles": s.Roles,
	}
}

// Deserialize populates the local DB authentication profile from a map.
func (s *ArkUAPSIADBDBAuthProfile) Deserialize(data map[string]interface{}) error {
	s.Roles = []string{}
	if roles, ok := data["roles"].([]interface{}); ok {
		for _, role := range roles {
			if roleStr, ok := role.(string); ok && len(roleStr) <= DBRoleMaxLength {
				s.Roles = append(s.Roles, roleStr)
			} else {
				return errors.New("invalid role name in roles")
			}
		}
	} else {
		return errors.New("roles must be an array of strings")
	}
	return nil
}

// ArkUAPSIADBOracleAuthProfile represents an Oracle DB authentication profile.
type ArkUAPSIADBOracleAuthProfile struct {
	Roles       []string `json:"roles" validate:"min=1,max=50" mapstructure:"roles" flag:"roles" desc:"List of roles assigned to the user"`
	DbaRole     bool     `json:"dba_role" mapstructure:"dba_role" flag:"dba-role" desc:"Indicates if the user has DBA role"`
	SysdbaRole  bool     `json:"sysdba_role" mapstructure:"sysdba_role" flag:"sysdba-role" desc:"Indicates if the user has SYSDBA role"`
	SysoperRole bool     `json:"sysoper_role" mapstructure:"sysoper_role" flag:"sysoper-role" desc:"Indicates if the user has SYSOPER role"`
}

// Serialize converts the Oracle DB authentication profile to a map.
func (s *ArkUAPSIADBOracleAuthProfile) Serialize() map[string]interface{} {
	return map[string]interface{}{
		"roles":       s.Roles,
		"dbaRole":     s.DbaRole,
		"sysdbaRole":  s.SysdbaRole,
		"sysoperRole": s.SysoperRole,
	}
}

// Deserialize populates the Oracle DB authentication profile from a map.
func (s *ArkUAPSIADBOracleAuthProfile) Deserialize(data map[string]interface{}) error {
	s.Roles = []string{}
	if roles, ok := data["roles"].([]interface{}); ok {
		for _, role := range roles {
			if roleStr, ok := role.(string); ok && len(roleStr) <= DBRoleMaxLength {
				s.Roles = append(s.Roles, roleStr)
			} else {
				return errors.New("invalid role name in roles")
			}
		}
	} else {
		return errors.New("roles must be an array of strings")
	}

	if dbaRole, ok := data["dba_role"].(bool); ok {
		s.DbaRole = dbaRole
	}

	if sysdbaRole, ok := data["sysdba_role"].(bool); ok {
		s.SysdbaRole = sysdbaRole
	}

	if sysoperRole, ok := data["sysoper_role"].(bool); ok {
		s.SysoperRole = sysoperRole
	}

	return nil
}

// ArkUAPSIADBMongoAuthProfile represents a MongoDB authentication profile.
type ArkUAPSIADBMongoAuthProfile struct {
	GlobalBuiltinRoles   []string            `json:"global_builtin_roles" validate:"max=50" mapstructure:"global_builtin_roles" flag:"global-builtin-roles" desc:"List of global built-in roles"`
	DatabaseBuiltinRoles map[string][]string `json:"database_builtin_roles" validate:"max=1000" mapstructure:"database_builtin_roles" flag:"database-builtin-roles" desc:"Map of database names to their built-in roles"`
	DatabaseCustomRoles  map[string][]string `json:"database_custom_roles" validate:"max=1000" mapstructure:"database_custom_roles" flag:"database-custom-roles" desc:"Map of database names to their custom roles"`
}

// Serialize converts the MongoDB authentication profile to a map.
func (s *ArkUAPSIADBMongoAuthProfile) Serialize() map[string]interface{} {
	return map[string]interface{}{
		"globalBuiltinRoles":   s.GlobalBuiltinRoles,
		"databaseBuiltinRoles": s.DatabaseBuiltinRoles,
		"databaseCustomRoles":  s.DatabaseCustomRoles,
	}
}

// Deserialize populates the MongoDB authentication profile from a map.
func (s *ArkUAPSIADBMongoAuthProfile) Deserialize(data map[string]interface{}) error {
	s.GlobalBuiltinRoles = []string{}
	s.DatabaseBuiltinRoles = make(map[string][]string)
	s.DatabaseCustomRoles = make(map[string][]string)
	if globalRoles, ok := data["global_builtin_roles"].([]interface{}); ok {
		for _, role := range globalRoles {
			if roleStr, ok := role.(string); ok && len(roleStr) <= DBRoleMaxLength {
				s.GlobalBuiltinRoles = append(s.GlobalBuiltinRoles, roleStr)
			} else {
				return errors.New("invalid global role name in global_builtin_roles")
			}
		}
	} else {
		return errors.New("global_builtin_roles must be an array of strings")
	}
	if dbBuiltinRoles, ok := data["database_builtin_roles"].(map[string]interface{}); ok {
		for db, roles := range dbBuiltinRoles {
			if len(db) > DatabaseNameMaxLength {
				return errors.New("database name exceeds maximum length in database_builtin_roles")
			}
			if roleList, ok := roles.([]interface{}); ok {
				for _, role := range roleList {
					if roleStr, ok := role.(string); ok && len(roleStr) <= DBRoleMaxLength {
						if _, ok := s.DatabaseBuiltinRoles[db]; !ok {
							s.DatabaseBuiltinRoles[db] = []string{}
						}
						s.DatabaseBuiltinRoles[db] = append(s.DatabaseBuiltinRoles[db], roleStr)
					} else {
						return errors.New("invalid database builtin role name in database_builtin_roles")
					}
				}
			} else {
				return errors.New("database_builtin_roles must be a map of arrays of strings")
			}
		}
	}

	if dbCustomRoles, ok := data["database_custom_roles"].(map[string]interface{}); ok {
		for db, roles := range dbCustomRoles {
			if len(db) > DatabaseNameMaxLength {
				return errors.New("database name exceeds maximum length in database_custom_roles")
			}
			if roleList, ok := roles.([]interface{}); ok {
				for _, role := range roleList {
					if roleStr, ok := role.(string); ok && len(roleStr) <= DBRoleMaxLength {
						if _, ok := s.DatabaseCustomRoles[db]; !ok {
							s.DatabaseCustomRoles[db] = []string{}
						}
						s.DatabaseCustomRoles[db] = append(s.DatabaseCustomRoles[db], roleStr)
					} else {
						return errors.New("invalid database custom role name in database_custom_roles")
					}
				}
			} else {
				return errors.New("database_custom_roles must be a map of arrays of strings")
			}
		}
	}

	return nil
}

// ValidateGlobalRoles validates global roles for MongoDB profiles.
func (s *ArkUAPSIADBMongoAuthProfile) ValidateGlobalRoles() error {
	if len(s.DatabaseBuiltinRoles) == 0 && len(s.DatabaseCustomRoles) == 0 {
		if len(s.GlobalBuiltinRoles) == 0 {
			return errors.New("at least one global role must be defined when no databases are specified")
		}
	}
	return nil
}

// ValidateDatabasesRolesLogic validates database roles logic for MongoDB profiles.
func (s *ArkUAPSIADBMongoAuthProfile) ValidateDatabasesRolesLogic() error {
	allChildrenDBNames := make(map[string]bool)
	for db := range s.DatabaseBuiltinRoles {
		allChildrenDBNames[db] = true
	}
	for db := range s.DatabaseCustomRoles {
		allChildrenDBNames[db] = true
	}

	var missingRoleDBs []string
	for db := range allChildrenDBNames {
		if len(s.DatabaseBuiltinRoles[db]) == 0 && len(s.DatabaseCustomRoles[db]) == 0 {
			missingRoleDBs = append(missingRoleDBs, db)
		}
	}

	if len(missingRoleDBs) > 0 {
		return errors.New("the following databases are missing roles: " + strings.Join(missingRoleDBs, ", "))
	}
	return nil
}

// ArkUAPSIADBSqlServerAuthProfile represents a SQL Server authentication profile.
type ArkUAPSIADBSqlServerAuthProfile struct {
	GlobalBuiltinRoles   []string            `json:"global_builtin_roles" validate:"max=50" mapstructure:"global_builtin_roles" flag:"global-builtin-roles" desc:"List of global built-in roles"`
	GlobalCustomRoles    []string            `json:"global_custom_roles" validate:"max=50" mapstructure:"global_custom_roles" flag:"global-custom-roles" desc:"List of global custom roles"`
	DatabaseBuiltinRoles map[string][]string `json:"database_builtin_roles" validate:"max=1000" mapstructure:"database_builtin_roles" flag:"database-builtin-roles" desc:"Map of database names to their built-in roles"`
	DatabaseCustomRoles  map[string][]string `json:"database_custom_roles" validate:"max=1000" mapstructure:"database_custom_roles" flag:"database-custom-roles" desc:"Map of database names to their custom roles"`
}

// Serialize converts the SQL Server authentication profile to a map.
func (s *ArkUAPSIADBSqlServerAuthProfile) Serialize() map[string]interface{} {
	return map[string]interface{}{
		"globalBuiltinRoles":   s.GlobalBuiltinRoles,
		"globalCustomRoles":    s.GlobalCustomRoles,
		"databaseBuiltinRoles": s.DatabaseBuiltinRoles,
		"databaseCustomRoles":  s.DatabaseCustomRoles,
	}
}

// Deserialize populates the SQL Server authentication profile from a map.
func (s *ArkUAPSIADBSqlServerAuthProfile) Deserialize(data map[string]interface{}) error {
	s.GlobalBuiltinRoles = []string{}
	s.GlobalCustomRoles = []string{}
	s.DatabaseBuiltinRoles = make(map[string][]string)
	s.DatabaseCustomRoles = make(map[string][]string)
	if globalRoles, ok := data["global_builtin_roles"].([]interface{}); ok {
		for _, role := range globalRoles {
			if roleStr, ok := role.(string); ok && len(roleStr) <= DBRoleMaxLength {
				s.GlobalBuiltinRoles = append(s.GlobalBuiltinRoles, roleStr)
			} else {
				return errors.New("invalid global role name in global_builtin_roles")
			}
		}
	} else {
		return errors.New("global_builtin_roles must be an array of strings")
	}

	if globalCustomRoles, ok := data["global_custom_roles"].([]interface{}); ok {
		for _, role := range globalCustomRoles {
			if roleStr, ok := role.(string); ok && len(roleStr) <= DBRoleMaxLength {
				s.GlobalCustomRoles = append(s.GlobalCustomRoles, roleStr)
			} else {
				return errors.New("invalid global custom role name in global_custom_roles")
			}
		}
	} else {
		return errors.New("global_custom_roles must be an array of strings")
	}

	if dbBuiltinRoles, ok := data["database_builtin_roles"].(map[string]interface{}); ok {
		for db, roles := range dbBuiltinRoles {
			if len(db) > DatabaseNameMaxLength {
				return errors.New("database name exceeds maximum length in database_builtin_roles")
			}
			if roleList, ok := roles.([]interface{}); ok {
				for _, role := range roleList {
					if roleStr, ok := role.(string); ok && len(roleStr) <= DBRoleMaxLength {
						if _, ok := s.DatabaseBuiltinRoles[db]; !ok {
							s.DatabaseBuiltinRoles[db] = []string{}
						}
						s.DatabaseBuiltinRoles[db] = append(s.DatabaseBuiltinRoles[db], roleStr)
					} else {
						return errors.New("invalid database builtin role name in database_builtin_roles")
					}
				}
			} else {
				return errors.New("database_builtin_roles must be a map of arrays of strings")
			}
		}
	}

	if dbCustomRoles, ok := data["database_custom_roles"].(map[string]interface{}); ok {
		for db, roles := range dbCustomRoles {
			if len(db) > DatabaseNameMaxLength {
				return errors.New("database name exceeds maximum length in database_custom_roles")
			}
			if roleList, ok := roles.([]interface{}); ok {
				for _, role := range roleList {
					if roleStr, ok := role.(string); ok && len(roleStr) <= DBRoleMaxLength {
						if _, ok := s.DatabaseCustomRoles[db]; !ok {
							s.DatabaseCustomRoles[db] = []string{}
						}
						s.DatabaseCustomRoles[db] = append(s.DatabaseCustomRoles[db], roleStr)
					} else {
						return errors.New("invalid database custom role name in database_custom_roles")
					}
				}
			} else {
				return errors.New("database_custom_roles must be a map of arrays of strings")
			}
		}
	}
	return nil
}

// ArkUAPSIADBRDSIAMUserAuthProfile represents an RDS IAM User authentication profile.
type ArkUAPSIADBRDSIAMUserAuthProfile struct {
	DBUser string `json:"db_user" validate:"min=1,max=256" mapstructure:"db_user" flag:"db-user" desc:"The database user for RDS IAM User authentication"`
}

// Serialize converts the RDS IAM User authentication profile to a map.
func (s *ArkUAPSIADBRDSIAMUserAuthProfile) Serialize() map[string]interface{} {
	return map[string]interface{}{
		"dbUser": s.DBUser,
	}
}

// Deserialize populates the RDS IAM User authentication profile from a map.
func (s *ArkUAPSIADBRDSIAMUserAuthProfile) Deserialize(data map[string]interface{}) error {
	if dbUser, ok := data["db_user"].(string); ok && len(dbUser) <= DBUserMaxLength {
		s.DBUser = dbUser
	} else {
		return errors.New("db_user must be a non-empty string with a maximum length of 256 characters")
	}
	return nil
}

// ArkUAPSIADBDBAuthMethodToProfile maps authentication methods to profile types.
var ArkUAPSIADBDBAuthMethodToProfile = map[string]interface{}{
	AuthMethodLDAPAuth:       ArkUAPSIADBLDAPAuthProfile{},
	AuthMethodDBAuth:         ArkUAPSIADBDBAuthProfile{},
	AuthMethodOracleAuth:     ArkUAPSIADBOracleAuthProfile{},
	AuthMethodMongoAuth:      ArkUAPSIADBMongoAuthProfile{},
	AuthMethodSQLServerAuth:  ArkUAPSIADBSqlServerAuthProfile{},
	AuthMethodRDSIAMUserAuth: ArkUAPSIADBRDSIAMUserAuthProfile{},
}
