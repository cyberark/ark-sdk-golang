package common

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

// TestStruct represents a test structure for schema-aware testing
type TestStruct struct {
	FirstName string                 `json:"first_name"`
	LastName  string                 `json:"last_name"`
	Age       int                    `json:"age"`
	Metadata  map[string]interface{} `json:"metadata"`
}

// EmbeddedStruct represents a test structure with embedded fields
type EmbeddedStruct struct {
	TestStruct `mapstructure:",squash"`
	Email      string `json:"email"`
}

func TestSerializeResponseToJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "success_valid_json_object",
			input:    `{"name": "John", "age": 30}`,
			expected: `{"age":30,"name":"John"}`,
		},
		{
			name:     "success_valid_json_array",
			input:    `[1, 2, 3]`,
			expected: `[1, 2, 3]`,
		},
		{
			name:     "error_invalid_json",
			input:    `{invalid json}`,
			expected: `{invalid json}`,
		},
		{
			name:     "success_empty_string",
			input:    ``,
			expected: ``,
		},
		{
			name:     "success_plain_text",
			input:    `plain text response`,
			expected: `plain text response`,
		},
		{
			name:     "success_json_string",
			input:    `"hello world"`,
			expected: `"hello world"`,
		},
		{
			name:     "success_json_number",
			input:    `42`,
			expected: `42`,
		},
		{
			name:     "success_json_boolean",
			input:    `true`,
			expected: `true`,
		},
		{
			name:     "success_json_null",
			input:    `null`,
			expected: `null`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := io.NopCloser(strings.NewReader(tt.input))
			result := SerializeResponseToJSON(response)
			if result != tt.expected {
				t.Errorf("SerializeResponseToJSON() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestConvertToSnakeCase(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		schema   *reflect.Type
		expected interface{}
	}{
		{
			name: "simple_map_with_camel_case_keys",
			input: map[string]interface{}{
				"firstName": "John",
				"lastName":  "Doe",
				"userAge":   30,
			},
			schema: nil,
			expected: map[string]interface{}{
				"first_name": "John",
				"last_name":  "Doe",
				"user_age":   30,
			},
		},
		{
			name: "nested_map",
			input: map[string]interface{}{
				"userInfo": map[string]interface{}{
					"firstName": "John",
					"lastName":  "Doe",
				},
				"accountDetails": map[string]interface{}{
					"accountId": "123",
					"isActive":  true,
				},
			},
			schema: nil,
			expected: map[string]interface{}{
				"user_info": map[string]interface{}{
					"first_name": "John",
					"last_name":  "Doe",
				},
				"account_details": map[string]interface{}{
					"account_id": "123",
					"is_active":  true,
				},
			},
		},
		{
			name: "array_of_maps",
			input: []interface{}{
				map[string]interface{}{
					"firstName": "John",
					"lastName":  "Doe",
				},
				map[string]interface{}{
					"firstName": "Jane",
					"lastName":  "Smith",
				},
			},
			schema: nil,
			expected: []interface{}{
				map[string]interface{}{
					"first_name": "John",
					"last_name":  "Doe",
				},
				map[string]interface{}{
					"first_name": "Jane",
					"last_name":  "Smith",
				},
			},
		},
		{
			name:     "primitive_types_unchanged",
			input:    "hello world",
			schema:   nil,
			expected: "hello world",
		},
		{
			name:     "number_unchanged",
			input:    42,
			schema:   nil,
			expected: 42,
		},
		{
			name:     "boolean_unchanged",
			input:    true,
			schema:   nil,
			expected: true,
		},
		{
			name:     "nil_input",
			input:    nil,
			schema:   nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertToSnakeCase(tt.input, tt.schema)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ConvertToSnakeCase() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestConvertToCamelCase(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		schema   *reflect.Type
		expected interface{}
	}{
		{
			name: "simple_map_with_snake_case_keys",
			input: map[string]interface{}{
				"first_name": "John",
				"last_name":  "Doe",
				"user_age":   30,
			},
			schema: nil,
			expected: map[string]interface{}{
				"firstName": "John",
				"lastName":  "Doe",
				"userAge":   30,
			},
		},
		{
			name: "nested_map",
			input: map[string]interface{}{
				"user_info": map[string]interface{}{
					"first_name": "John",
					"last_name":  "Doe",
				},
				"account_details": map[string]interface{}{
					"account_id": "123",
					"is_active":  true,
				},
			},
			schema: nil,
			expected: map[string]interface{}{
				"userInfo": map[string]interface{}{
					"firstName": "John",
					"lastName":  "Doe",
				},
				"accountDetails": map[string]interface{}{
					"accountId": "123",
					"isActive":  true,
				},
			},
		},
		{
			name: "array_of_maps",
			input: []interface{}{
				map[string]interface{}{
					"first_name": "John",
					"last_name":  "Doe",
				},
				map[string]interface{}{
					"first_name": "Jane",
					"last_name":  "Smith",
				},
			},
			schema: nil,
			expected: []interface{}{
				map[string]interface{}{
					"firstName": "John",
					"lastName":  "Doe",
				},
				map[string]interface{}{
					"firstName": "Jane",
					"lastName":  "Smith",
				},
			},
		},
		{
			name:     "primitive_types_unchanged",
			input:    "hello world",
			schema:   nil,
			expected: "hello world",
		},
		{
			name:     "number_unchanged",
			input:    42,
			schema:   nil,
			expected: 42,
		},
		{
			name:     "boolean_unchanged",
			input:    true,
			schema:   nil,
			expected: true,
		},
		{
			name:     "nil_input",
			input:    nil,
			schema:   nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertToCamelCase(tt.input, tt.schema)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ConvertToCamelCase() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDeserializeJSONSnake(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    interface{}
		expectError bool
	}{
		{
			name:  "simple_json_object_with_camel_case_keys",
			input: `{"firstName": "John", "lastName": "Doe", "userAge": 30}`,
			expected: map[string]interface{}{
				"first_name": "John",
				"last_name":  "Doe",
				"user_age":   float64(30), // JSON numbers are float64
			},
			expectError: false,
		},
		{
			name:  "nested_json_object",
			input: `{"userInfo": {"firstName": "John", "lastName": "Doe"}, "isActive": true}`,
			expected: map[string]interface{}{
				"user_info": map[string]interface{}{
					"first_name": "John",
					"last_name":  "Doe",
				},
				"is_active": true,
			},
			expectError: false,
		},
		{
			name:  "json_array",
			input: `[{"firstName": "John"}, {"firstName": "Jane"}]`,
			expected: []interface{}{
				map[string]interface{}{
					"first_name": "John",
				},
				map[string]interface{}{
					"first_name": "Jane",
				},
			},
			expectError: false,
		},
		{
			name:        "invalid_json",
			input:       `{invalid json}`,
			expected:    nil,
			expectError: true,
		},
		{
			name:        "empty_json",
			input:       `{}`,
			expected:    map[string]interface{}{},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := io.NopCloser(strings.NewReader(tt.input))
			result, err := DeserializeJSONSnake(response)

			if tt.expectError {
				if err == nil {
					t.Errorf("DeserializeJSONSnake() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("DeserializeJSONSnake() unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("DeserializeJSONSnake() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSerializeJSONCamel(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expected    map[string]interface{}
		expectError bool
	}{
		{
			name: "simple_struct",
			input: TestStruct{
				FirstName: "John",
				LastName:  "Doe",
				Age:       30,
				Metadata: map[string]interface{}{
					"key1": "value1",
				},
			},
			expected: map[string]interface{}{
				"firstName": "John",
				"lastName":  "Doe",
				"age":       float64(30),
				"metadata": map[string]interface{}{
					"key1": "value1",
				},
			},
			expectError: false,
		},
		{
			name: "map_with_snake_case_keys",
			input: map[string]interface{}{
				"first_name": "John",
				"last_name":  "Doe",
				"user_age":   30,
			},
			expected: map[string]interface{}{
				"firstName": "John",
				"lastName":  "Doe",
				"userAge":   float64(30),
			},
			expectError: false,
		},
		{
			name: "nested_structure",
			input: map[string]interface{}{
				"user_info": map[string]interface{}{
					"first_name": "John",
					"last_name":  "Doe",
				},
				"account_status": "active",
			},
			expected: map[string]interface{}{
				"userInfo": map[string]interface{}{
					"firstName": "John",
					"lastName":  "Doe",
				},
				"accountStatus": "active",
			},
			expectError: false,
		},
		{
			name:        "invalid_input_channel",
			input:       make(chan int),
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SerializeJSONCamel(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("SerializeJSONCamel() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("SerializeJSONCamel() unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("SerializeJSONCamel() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDeserializeJSONSnakeSchema(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		schema      *reflect.Type
		expected    interface{}
		expectError bool
	}{
		{
			name:   "json_without_schema",
			input:  `{"firstName": "John", "lastName": "Doe"}`,
			schema: nil,
			expected: map[string]interface{}{
				"first_name": "John",
				"last_name":  "Doe",
			},
			expectError: false,
		},
		{
			name:        "invalid_json",
			input:       `{invalid json}`,
			schema:      nil,
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := io.NopCloser(strings.NewReader(tt.input))
			result, err := DeserializeJSONSnakeSchema(response, tt.schema)

			if tt.expectError {
				if err == nil {
					t.Errorf("DeserializeJSONSnakeSchema() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("DeserializeJSONSnakeSchema() unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("DeserializeJSONSnakeSchema() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSerializeJSONCamelSchema(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		schema      *reflect.Type
		expected    map[string]interface{}
		expectError bool
	}{
		{
			name: "map_without_schema",
			input: map[string]interface{}{
				"first_name": "John",
				"last_name":  "Doe",
			},
			schema: nil,
			expected: map[string]interface{}{
				"firstName": "John",
				"lastName":  "Doe",
			},
			expectError: false,
		},
		{
			name:        "invalid_input_with_schema",
			input:       make(chan int),
			schema:      func() *reflect.Type { t := reflect.TypeOf(TestStruct{}); return &t }(),
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SerializeJSONCamelSchema(tt.input, tt.schema)

			if tt.expectError {
				if err == nil {
					t.Errorf("SerializeJSONCamelSchema() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("SerializeJSONCamelSchema() unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("SerializeJSONCamelSchema() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestResolveFieldsSquashed(t *testing.T) {
	tests := []struct {
		name     string
		schema   reflect.Type
		expected int // expected number of fields
	}{
		{
			name:     "simple_struct",
			schema:   reflect.TypeOf(TestStruct{}),
			expected: 4, // FirstName, LastName, Age, Metadata
		},
		{
			name:     "embedded_struct",
			schema:   reflect.TypeOf(EmbeddedStruct{}),
			expected: 5, // FirstName, LastName, Age, Metadata (from embedded), Email
		},
		{
			name:     "pointer_to_struct",
			schema:   reflect.TypeOf(&TestStruct{}),
			expected: 4, // FirstName, LastName, Age, Metadata
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := resolveFieldsSquashed(tt.schema)
			if len(result) != tt.expected {
				t.Errorf("resolveFieldsSquashed() returned %d fields, want %d", len(result), tt.expected)
			}
		})
	}
}

func TestFindFieldByName(t *testing.T) {
	tests := []struct {
		name       string
		schema     reflect.Type
		fieldName  string
		expectNil  bool
		expectName string
	}{
		{
			name:       "find_existing_field",
			schema:     reflect.TypeOf(TestStruct{}),
			fieldName:  "first_name",
			expectNil:  false,
			expectName: "FirstName",
		},
		{
			name:       "find_field_with_different_case",
			schema:     reflect.TypeOf(TestStruct{}),
			fieldName:  "firstname",
			expectNil:  false,
			expectName: "FirstName",
		},
		{
			name:      "field_not_found",
			schema:    reflect.TypeOf(TestStruct{}),
			fieldName: "nonexistent_field",
			expectNil: true,
		},
		{
			name:       "find_field_in_pointer_type",
			schema:     reflect.TypeOf(&TestStruct{}),
			fieldName:  "last_name",
			expectNil:  false,
			expectName: "LastName",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findFieldByName(tt.schema, tt.fieldName)
			if tt.expectNil {
				if result != nil {
					t.Errorf("findFieldByName() expected nil but got %v", result)
				}
			} else {
				if result == nil {
					t.Errorf("findFieldByName() expected field but got nil")
				} else if result.Name != tt.expectName {
					t.Errorf("findFieldByName() expected field name %s but got %s", tt.expectName, result.Name)
				}
			}
		})
	}
}

// Test error handling for SerializeResponseToJSON
func TestSerializeResponseToJSON_ErrorHandling(t *testing.T) {
	// Test with a ReadCloser that returns an error
	errorReader := &errorReadCloser{}
	result := SerializeResponseToJSON(errorReader)
	if result != "" {
		t.Errorf("SerializeResponseToJSON() with error reader expected empty string, got %v", result)
	}
}

// errorReadCloser is a helper type that always returns an error when reading
type errorReadCloser struct{}

func (e *errorReadCloser) Read(p []byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}

func (e *errorReadCloser) Close() error {
	return nil
}
