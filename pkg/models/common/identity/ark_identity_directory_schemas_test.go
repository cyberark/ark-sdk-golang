package identity

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestNewDirectoryServiceQueryRequest(t *testing.T) {
	tests := []struct {
		name           string
		searchString   string
		expectedResult *DirectoryServiceQueryRequest
		validateFunc   func(t *testing.T, result *DirectoryServiceQueryRequest)
	}{
		{
			name:         "success_empty_search_string",
			searchString: "",
			expectedResult: &DirectoryServiceQueryRequest{
				User:  "{}",
				Roles: "{}",
				Group: "{}",
			},
		},
		{
			name:         "success_with_search_string",
			searchString: "admin",
			validateFunc: func(t *testing.T, result *DirectoryServiceQueryRequest) {
				if result.User == "{}" {
					t.Error("Expected user filter to be set with search string")
				}
				if result.Roles == "{}" {
					t.Error("Expected roles filter to be set with search string")
				}
				if result.Group == "{}" {
					t.Error("Expected group filter to be set with search string")
				}

				// Verify that filters contain the search string
				var userFilter map[string]interface{}
				err := json.Unmarshal([]byte(result.User), &userFilter)
				if err != nil {
					t.Errorf("Failed to unmarshal user filter: %v", err)
				}

				var rolesFilter map[string]interface{}
				err = json.Unmarshal([]byte(result.Roles), &rolesFilter)
				if err != nil {
					t.Errorf("Failed to unmarshal roles filter: %v", err)
				}

				var groupFilter map[string]interface{}
				err = json.Unmarshal([]byte(result.Group), &groupFilter)
				if err != nil {
					t.Errorf("Failed to unmarshal group filter: %v", err)
				}
			},
		},
		{
			name:         "success_special_characters_search",
			searchString: "test@domain.com",
			validateFunc: func(t *testing.T, result *DirectoryServiceQueryRequest) {
				if result.User == "{}" || result.Roles == "{}" || result.Group == "{}" {
					t.Error("Expected all filters to be set with special character search string")
				}
			},
		},
		{
			name:         "success_unicode_search",
			searchString: "üñîçödé",
			validateFunc: func(t *testing.T, result *DirectoryServiceQueryRequest) {
				if result.User == "{}" || result.Roles == "{}" || result.Group == "{}" {
					t.Error("Expected all filters to be set with unicode search string")
				}
			},
		},
		{
			name:         "edge_case_single_character",
			searchString: "a",
			validateFunc: func(t *testing.T, result *DirectoryServiceQueryRequest) {
				if result.User == "{}" || result.Roles == "{}" || result.Group == "{}" {
					t.Error("Expected all filters to be set with single character search")
				}
			},
		},
		{
			name:         "edge_case_long_search_string",
			searchString: "very_long_search_string_that_exceeds_normal_length_expectations_for_testing_purposes",
			validateFunc: func(t *testing.T, result *DirectoryServiceQueryRequest) {
				if result.User == "{}" || result.Roles == "{}" || result.Group == "{}" {
					t.Error("Expected all filters to be set with long search string")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := NewDirectoryServiceQueryRequest(tt.searchString)

			// Validate that result is not nil
			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			// Use custom validation if provided
			if tt.validateFunc != nil {
				tt.validateFunc(t, result)
				return
			}

			// Validate expected result if provided
			if tt.expectedResult != nil {
				if result.User != tt.expectedResult.User {
					t.Errorf("Expected User '%s', got '%s'", tt.expectedResult.User, result.User)
				}
				if result.Roles != tt.expectedResult.Roles {
					t.Errorf("Expected Roles '%s', got '%s'", tt.expectedResult.Roles, result.Roles)
				}
				if result.Group != tt.expectedResult.Group {
					t.Errorf("Expected Group '%s', got '%s'", tt.expectedResult.Group, result.Group)
				}
			}
		})
	}
}

func TestNewDirectoryServiceQuerySpecificRoleRequest(t *testing.T) {
	tests := []struct {
		name           string
		roleName       string
		expectedResult *DirectoryServiceQuerySpecificRoleRequest
		validateFunc   func(t *testing.T, result *DirectoryServiceQuerySpecificRoleRequest)
	}{
		{
			name:     "success_empty_role_name",
			roleName: "",
			expectedResult: &DirectoryServiceQuerySpecificRoleRequest{
				User:  "{}",
				Roles: "{}",
				Group: "{}",
			},
		},
		{
			name:     "success_with_role_name",
			roleName: "System Administrator",
			expectedResult: &DirectoryServiceQuerySpecificRoleRequest{
				User:  "{}",
				Roles: `{"Name":{"_eq":"System Administrator"}}`,
				Group: "{}",
			},
		},
		{
			name:     "success_role_with_spaces",
			roleName: "Database Admin",
			expectedResult: &DirectoryServiceQuerySpecificRoleRequest{
				User:  "{}",
				Roles: `{"Name":{"_eq":"Database Admin"}}`,
				Group: "{}",
			},
		},
		{
			name:     "success_role_with_special_characters",
			roleName: "Admin@Domain.com",
			expectedResult: &DirectoryServiceQuerySpecificRoleRequest{
				User:  "{}",
				Roles: `{"Name":{"_eq":"Admin@Domain.com"}}`,
				Group: "{}",
			},
		},
		{
			name:     "edge_case_single_character_role",
			roleName: "A",
			expectedResult: &DirectoryServiceQuerySpecificRoleRequest{
				User:  "{}",
				Roles: `{"Name":{"_eq":"A"}}`,
				Group: "{}",
			},
		},
		{
			name:     "edge_case_unicode_role_name",
			roleName: "Administrâtör",
			expectedResult: &DirectoryServiceQuerySpecificRoleRequest{
				User:  "{}",
				Roles: `{"Name":{"_eq":"Administrâtör"}}`,
				Group: "{}",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := NewDirectoryServiceQuerySpecificRoleRequest(tt.roleName)

			// Validate that result is not nil
			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			// Use custom validation if provided
			if tt.validateFunc != nil {
				tt.validateFunc(t, result)
				return
			}

			// Validate expected result if provided
			if tt.expectedResult != nil {
				if result.User != tt.expectedResult.User {
					t.Errorf("Expected User '%s', got '%s'", tt.expectedResult.User, result.User)
				}
				if result.Roles != tt.expectedResult.Roles {
					t.Errorf("Expected Roles '%s', got '%s'", tt.expectedResult.Roles, result.Roles)
				}
				if result.Group != tt.expectedResult.Group {
					t.Errorf("Expected Group '%s', got '%s'", tt.expectedResult.Group, result.Group)
				}
			}
		})
	}
}

// TestDirectoryTypes_Constants tests the directory type constants
func TestDirectoryTypes_Constants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{
			name:     "ad_constant_value",
			constant: AD,
			expected: "AdProxy",
		},
		{
			name:     "identity_constant_value",
			constant: Identity,
			expected: "CDS",
		},
		{
			name:     "fds_constant_value",
			constant: FDS,
			expected: "FDS",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.constant != tt.expected {
				t.Errorf("Expected constant value '%s', got '%s'", tt.expected, tt.constant)
			}
		})
	}
}

// TestAllDirectoryTypes_Variable tests the AllDirectoryTypes variable
func TestAllDirectoryTypes_Variable(t *testing.T) {
	t.Run("all_directory_types_content", func(t *testing.T) {
		expected := []string{"AdProxy", "CDS", "FDS"}

		if len(AllDirectoryTypes) != len(expected) {
			t.Errorf("Expected %d directory types, got %d", len(expected), len(AllDirectoryTypes))
		}

		for i, expectedType := range expected {
			if i >= len(AllDirectoryTypes) {
				t.Errorf("Missing expected directory type at index %d: %s", i, expectedType)
				continue
			}
			if AllDirectoryTypes[i] != expectedType {
				t.Errorf("Expected directory type at index %d to be '%s', got '%s'", i, expectedType, AllDirectoryTypes[i])
			}
		}
	})

	t.Run("all_directory_types_matches_constants", func(t *testing.T) {
		expectedTypes := map[string]bool{
			AD:       true,
			Identity: true,
			FDS:      true,
		}

		for _, dirType := range AllDirectoryTypes {
			if !expectedTypes[dirType] {
				t.Errorf("AllDirectoryTypes contains unexpected directory type: %s", dirType)
			}
		}

		if len(AllDirectoryTypes) != len(expectedTypes) {
			t.Errorf("AllDirectoryTypes length (%d) doesn't match expected constants count (%d)", len(AllDirectoryTypes), len(expectedTypes))
		}
	})
}

// TestStructJSONMarshaling tests JSON marshaling/unmarshaling for structs
func TestStructJSONMarshaling(t *testing.T) {
	tests := []struct {
		name       string
		structData interface{}
		jsonData   string
	}{
		{
			name: "directory_service_metadata_marshaling",
			structData: DirectoryServiceMetadata{
				Service:              "TestService",
				DirectoryServiceUUID: "uuid-123-456",
			},
			jsonData: `{"Service":"TestService","directoryServiceUuid":"uuid-123-456"}`,
		},
		{
			name: "directory_search_args_marshaling",
			structData: DirectorySearchArgs{
				PageNumber: 1,
				PageSize:   50,
				Limit:      100,
				SortBy:     "name",
				Caching:    1,
				Dir:        "asc",
				Ascending:  true,
			},
			jsonData: `{"PageNumber":1,"PageSize":50,"Limit":100,"SortBy":"name","Caching":1,"Direction":"asc","Ascending":true}`,
		},
		{
			name: "group_row_marshaling",
			structData: GroupRow{
				DisplayName:              "Test Group",
				ServiceInstanceLocalized: "Test Service",
				DirectoryServiceType:     "CDS",
				SystemName:               "testgroup",
				InternalID:               "internal-123",
			},
			jsonData: `{"DisplayName":"Test Group","ServiceInstanceLocalized":"Test Service","ServiceType":"CDS","SystemName":"testgroup","InternalName":"internal-123"}`,
		},
		{
			name: "user_row_marshaling",
			structData: UserRow{
				DisplayName:              "Test User",
				ServiceInstanceLocalized: "Test Service",
				DistinguishedName:        "CN=Test User,OU=Users,DC=example,DC=com",
				SystemName:               "testuser",
				DirectoryServiceType:     "CDS",
				Email:                    "test@example.com",
				InternalID:               "user-123",
				Description:              "Test user account",
			},
			jsonData: `{"DisplayName":"Test User","ServiceInstanceLocalized":"Test Service","DistinguishedName":"CN=Test User,OU=Users,DC=example,DC=com","SystemName":"testuser","ServiceType":"CDS","EMail":"test@example.com","InternalName":"user-123","Description":"Test user account"}`,
		},
		{
			name: "role_row_marshaling",
			structData: RoleRow{
				Name: "Administrator",
				ID:   "role-123",
				AdminRights: []RoleAdminRight{
					{
						Path:        "/admin",
						ServiceName: "TestService",
					},
				},
				IsHidden:    false,
				Description: "Administrator role",
			},
			jsonData: `{"Name":"Administrator","_ID":"role-123","AdministrativeRights":[{"Path":"/admin","ServiceName":"TestService"}],"Description":"Administrator role"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Test marshaling
			marshaledData, err := json.Marshal(tt.structData)
			if err != nil {
				t.Fatalf("Failed to marshal struct: %v", err)
			}

			// Compare JSON (order-independent comparison)
			var expected, actual interface{}
			if err := json.Unmarshal([]byte(tt.jsonData), &expected); err != nil {
				t.Fatalf("Failed to unmarshal expected JSON: %v", err)
			}
			if err := json.Unmarshal(marshaledData, &actual); err != nil {
				t.Fatalf("Failed to unmarshal actual JSON: %v", err)
			}

			if !reflect.DeepEqual(expected, actual) {
				t.Errorf("JSON mismatch.\nExpected: %s\nActual: %s", tt.jsonData, string(marshaledData))
			}

			// Test unmarshaling back to struct
			structType := reflect.TypeOf(tt.structData)
			newStruct := reflect.New(structType).Interface()

			if err := json.Unmarshal([]byte(tt.jsonData), newStruct); err != nil {
				t.Fatalf("Failed to unmarshal JSON to struct: %v", err)
			}

			// Compare structs
			actualStruct := reflect.ValueOf(newStruct).Elem().Interface()
			if !reflect.DeepEqual(tt.structData, actualStruct) {
				t.Errorf("Struct mismatch after round-trip.\nExpected: %+v\nActual: %+v", tt.structData, actualStruct)
			}
		})
	}
}

// TestStructDefaults tests struct field defaults and empty values
func TestStructDefaults(t *testing.T) {
	tests := []struct {
		name         string
		createStruct func() interface{}
		validateFunc func(t *testing.T, s interface{})
	}{
		{
			name: "directory_service_metadata_defaults",
			createStruct: func() interface{} {
				return DirectoryServiceMetadata{}
			},
			validateFunc: func(t *testing.T, s interface{}) {
				metadata := s.(DirectoryServiceMetadata)
				if metadata.Service != "" {
					t.Errorf("Expected empty Service, got '%s'", metadata.Service)
				}
				if metadata.DirectoryServiceUUID != "" {
					t.Errorf("Expected empty DirectoryServiceUUID, got '%s'", metadata.DirectoryServiceUUID)
				}
			},
		},
		{
			name: "directory_search_args_defaults",
			createStruct: func() interface{} {
				return DirectorySearchArgs{}
			},
			validateFunc: func(t *testing.T, s interface{}) {
				args := s.(DirectorySearchArgs)
				if args.PageNumber != 0 {
					t.Errorf("Expected PageNumber 0, got %d", args.PageNumber)
				}
				if args.Ascending != false {
					t.Errorf("Expected Ascending false, got %t", args.Ascending)
				}
			},
		},
		{
			name: "group_row_defaults",
			createStruct: func() interface{} {
				return GroupRow{}
			},
			validateFunc: func(t *testing.T, s interface{}) {
				row := s.(GroupRow)
				if row.DisplayName != "" {
					t.Errorf("Expected empty DisplayName, got '%s'", row.DisplayName)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			structInstance := tt.createStruct()
			tt.validateFunc(t, structInstance)
		})
	}
}
