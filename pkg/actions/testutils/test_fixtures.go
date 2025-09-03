package testutils

// TestSchema represents a common test schema structure for service action testing.
//
// TestSchema provides a standardized struct with various field types and
// validation tags that can be used across different action tests. It includes
// common patterns like required fields, choice validation, and complex types.
type TestSchema struct {
	Name        string            `mapstructure:"name" validate:"required"`
	Count       int               `mapstructure:"count"`
	Tags        []string          `mapstructure:"tags"`
	Metadata    map[string]string `mapstructure:"metadata"`
	ComplexData []TestComplexType `mapstructure:"complex_data" desc:"Complex data field"`
	Choices     string            `mapstructure:"choices" choices:"option1,option2,option3"`
}

// TestComplexType represents a nested struct for complex type testing.
//
// TestComplexType provides a simple nested structure that can be used
// to test complex type parsing and serialization in service actions.
type TestComplexType struct {
	ID   string `mapstructure:"id"`
	Type string `mapstructure:"type"`
}

// CreateTestSchema creates a TestSchema instance with default test values.
//
// CreateTestSchema generates a TestSchema with sensible default values
// for testing purposes. This provides a consistent starting point for
// tests that need a populated schema.
//
// Returns a pointer to a TestSchema with default test values.
//
// Example:
//
//	schema := CreateTestSchema()
//	// schema.Name == "test-schema"
func CreateTestSchema() *TestSchema {
	return &TestSchema{
		Name:     "test-schema",
		Count:    42,
		Tags:     []string{"tag1", "tag2"},
		Metadata: map[string]string{"key": "value"},
		ComplexData: []TestComplexType{
			{ID: "1", Type: "test"},
		},
		Choices: "option1",
	}
}

// CreateTestSchemaWithName creates a TestSchema with a specific name.
//
// CreateTestSchemaWithName generates a TestSchema with the provided name
// and default values for other fields. This is useful when tests need
// to verify name-specific behavior.
//
// Parameters:
//   - name: The name to set in the schema
//
// Returns a pointer to a TestSchema with the specified name.
//
// Example:
//
//	schema := CreateTestSchemaWithName("custom-name")
//	// schema.Name == "custom-name"
func CreateTestSchemaWithName(name string) *TestSchema {
	schema := CreateTestSchema()
	schema.Name = name
	return schema
}
