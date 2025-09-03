package actions

import (
	"reflect"
	"strings"
	"testing"

	"github.com/cyberark/ark-sdk-golang/pkg/actions/testutils"
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	"github.com/cyberark/ark-sdk-golang/pkg/profiles"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Mock service for testing method finding
type mockService struct {
	testActionFunc func(*testutils.TestSchema) (interface{}, error)
}

func (m *mockService) TestAction(schema *testutils.TestSchema) (interface{}, error) {
	if m.testActionFunc != nil {
		return m.testActionFunc(schema)
	}
	return map[string]interface{}{"result": "success", "name": schema.Name}, nil
}

func TestNewArkServiceExecAction(t *testing.T) {
	tests := []struct {
		name         string
		setupLoader  func() *profiles.ProfileLoader
		validateFunc func(t *testing.T, action *ArkServiceExecAction)
	}{
		{
			name: "success_creates_action_with_profile_loader",
			setupLoader: func() *profiles.ProfileLoader {
				mock := testutils.NewMockProfileLoader()
				return mock.AsProfileLoader()
			},
			validateFunc: func(t *testing.T, action *ArkServiceExecAction) {
				if action == nil {
					t.Error("Expected non-nil action")
				}
				if action.ArkBaseExecAction == nil {
					t.Error("Expected non-nil ArkBaseExecAction")
				}
			},
		},
		{
			name: "success_creates_action_with_nil_loader",
			setupLoader: func() *profiles.ProfileLoader {
				return nil
			},
			validateFunc: func(t *testing.T, action *ArkServiceExecAction) {
				if action == nil {
					t.Error("Expected non-nil action")
				}
				if action.ArkBaseExecAction == nil {
					t.Error("Expected non-nil ArkBaseExecAction")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			loader := tt.setupLoader()
			action := NewArkServiceExecAction(loader)

			tt.validateFunc(t, action)
		})
	}
}

func TestArkServiceExecAction_isComplexType(t *testing.T) {
	tests := []struct {
		name     string
		field    reflect.StructField
		expected bool
	}{
		{
			name: "success_map_string_struct_is_complex",
			field: reflect.StructField{
				Type: reflect.TypeOf(map[string]testutils.TestComplexType{}),
			},
			expected: true,
		},
		{
			name: "success_slice_struct_is_complex",
			field: reflect.StructField{
				Type: reflect.TypeOf([]testutils.TestComplexType{}),
			},
			expected: true,
		},
		{
			name: "success_array_struct_is_complex",
			field: reflect.StructField{
				Type: reflect.TypeOf([5]testutils.TestComplexType{}),
			},
			expected: true,
		},
		{
			name: "success_string_is_not_complex",
			field: reflect.StructField{
				Type: reflect.TypeOf(""),
			},
			expected: false,
		},
		{
			name: "success_int_is_not_complex",
			field: reflect.StructField{
				Type: reflect.TypeOf(0),
			},
			expected: false,
		},
		{
			name: "success_slice_string_is_not_complex",
			field: reflect.StructField{
				Type: reflect.TypeOf([]string{}),
			},
			expected: false,
		},
		{
			name: "success_map_string_string_is_not_complex",
			field: reflect.StructField{
				Type: reflect.TypeOf(map[string]string{}),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			action := NewArkServiceExecAction(nil)
			result := action.isComplexType(tt.field)

			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestArkServiceExecAction_fillRemainingSchema(t *testing.T) {
	tests := []struct {
		name         string
		schema       interface{}
		validateFunc func(t *testing.T, flags *pflag.FlagSet)
	}{
		{
			name:   "success_adds_complex_type_flags",
			schema: testutils.CreateTestSchema(),
			validateFunc: func(t *testing.T, flags *pflag.FlagSet) {
				// Check that complex_data flag was added for the complex slice field
				flag := flags.Lookup("complex_data")
				if flag == nil {
					t.Error("Expected complex_data flag to be added")
				}
				if flag != nil && !strings.Contains(flag.Usage, "JSON") {
					t.Error("Expected complex type flag to have JSON hint in description")
				}
			},
		},
		{
			name:   "success_handles_empty_schema",
			schema: &struct{}{},
			validateFunc: func(t *testing.T, flags *pflag.FlagSet) {
				// Should not add any flags for empty schema
				if flags.NFlag() > 0 {
					t.Error("Expected no flags for empty schema")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			action := NewArkServiceExecAction(nil)
			flags := pflag.NewFlagSet("test", pflag.ContinueOnError)

			action.fillRemainingSchema(tt.schema, flags)

			tt.validateFunc(t, flags)
		})
	}
}

func TestArkServiceExecAction_fillParsedFlag(t *testing.T) {
	tests := []struct {
		name          string
		schemaElem    reflect.Type
		flags         map[string]interface{}
		key           string
		flag          *pflag.Flag
		expectedError bool
		validateFunc  func(t *testing.T, flags map[string]interface{})
	}{
		{
			name:          "success_parses_json_for_complex_slice",
			schemaElem:    reflect.TypeOf(*testutils.CreateTestSchema()),
			flags:         map[string]interface{}{"complex_data": `[{"id":"1","type":"test"}]`},
			key:           "complex_data",
			flag:          &pflag.Flag{Name: "complex-data"},
			expectedError: false,
			validateFunc: func(t *testing.T, flags map[string]interface{}) {
				value, ok := flags["complex_data"]
				if !ok {
					t.Error("Expected complex_data in flags")
					return
				}
				sliceVal, ok := value.([]map[string]interface{})
				if !ok {
					t.Errorf("Expected []map[string]interface{}, got %T", value)
					return
				}
				if len(sliceVal) != 1 {
					t.Errorf("Expected 1 element, got %d", len(sliceVal))
				}
			},
		},
		{
			name:          "error_invalid_json_for_complex_type",
			schemaElem:    reflect.TypeOf(*testutils.CreateTestSchema()),
			flags:         map[string]interface{}{"complex_data": `invalid json`},
			key:           "complex_data",
			flag:          &pflag.Flag{Name: "complex-data"},
			expectedError: true,
		},
		{
			name:          "error_invalid_choice_value",
			schemaElem:    reflect.TypeOf(*testutils.CreateTestSchema()),
			flags:         map[string]interface{}{"choices": "invalid_choice"},
			key:           "choices",
			flag:          &pflag.Flag{Name: "choices"},
			expectedError: true,
		},
		{
			name:          "success_valid_choice_value",
			schemaElem:    reflect.TypeOf(*testutils.CreateTestSchema()),
			flags:         map[string]interface{}{"choices": "option1"},
			key:           "choices",
			flag:          &pflag.Flag{Name: "choices"},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			action := NewArkServiceExecAction(nil)
			err := action.fillParsedFlag(tt.schemaElem, tt.flags, tt.key, tt.flag)

			if tt.expectedError && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if !tt.expectedError && tt.validateFunc != nil {
				tt.validateFunc(t, tt.flags)
			}
		})
	}
}

func TestArkServiceExecAction_parseFlag(t *testing.T) {
	tests := []struct {
		name          string
		setupCmd      func() *cobra.Command
		setupFlag     func(*cobra.Command) *pflag.Flag
		schema        interface{}
		expectedError bool
		validateFunc  func(t *testing.T, flags map[string]interface{})
	}{
		{
			name: "success_parses_string_flag",
			setupCmd: func() *cobra.Command {
				cmd := &cobra.Command{}
				cmd.Flags().String("test-flag", "default", "test flag")
				return cmd
			},
			setupFlag: func(cmd *cobra.Command) *pflag.Flag {
				flag := cmd.Flags().Lookup("test-flag")
				flag.Value.Set("test-value")
				flag.Changed = true
				return flag
			},
			schema:        testutils.CreateTestSchema(),
			expectedError: false,
			validateFunc: func(t *testing.T, flags map[string]interface{}) {
				if flags["test_flag"] != "test-value" {
					t.Errorf("Expected test-value, got %v", flags["test_flag"])
				}
			},
		},
		{
			name: "success_parses_int_flag",
			setupCmd: func() *cobra.Command {
				cmd := &cobra.Command{}
				cmd.Flags().Int("count", 0, "count flag")
				return cmd
			},
			setupFlag: func(cmd *cobra.Command) *pflag.Flag {
				flag := cmd.Flags().Lookup("count")
				flag.Value.Set("42")
				flag.Changed = true
				return flag
			},
			schema:        testutils.CreateTestSchema(),
			expectedError: false,
			validateFunc: func(t *testing.T, flags map[string]interface{}) {
				if flags["count"] != 42 {
					t.Errorf("Expected 42, got %v", flags["count"])
				}
			},
		},
		{
			name: "success_skips_unchanged_flag",
			setupCmd: func() *cobra.Command {
				cmd := &cobra.Command{}
				cmd.Flags().String("unchanged", "default", "unchanged flag")
				return cmd
			},
			setupFlag: func(cmd *cobra.Command) *pflag.Flag {
				flag := cmd.Flags().Lookup("unchanged")
				flag.Changed = false // Not changed
				return flag
			},
			schema:        testutils.CreateTestSchema(),
			expectedError: false,
			validateFunc: func(t *testing.T, flags map[string]interface{}) {
				if _, exists := flags["unchanged"]; exists {
					t.Error("Expected unchanged flag to be skipped")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			action := NewArkServiceExecAction(nil)
			cmd := tt.setupCmd()
			flag := tt.setupFlag(cmd)
			flags := make(map[string]interface{})

			err := action.parseFlag(flag, cmd, flags, tt.schema)

			if tt.expectedError && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if !tt.expectedError && tt.validateFunc != nil {
				tt.validateFunc(t, flags)
			}
		})
	}
}

func TestArkServiceExecAction_findMethodByName(t *testing.T) {
	tests := []struct {
		name          string
		value         reflect.Value
		methodName    string
		expectedError bool
		expectedFound bool
	}{
		{
			name:          "success_finds_exact_match",
			value:         reflect.ValueOf(&mockService{}),
			methodName:    "TestAction",
			expectedError: false,
			expectedFound: true,
		},
		{
			name:          "success_finds_case_insensitive_match",
			value:         reflect.ValueOf(&mockService{}),
			methodName:    "testaction",
			expectedError: false,
			expectedFound: true,
		},
		{
			name:          "error_method_not_found",
			value:         reflect.ValueOf(&mockService{}),
			methodName:    "NonExistentMethod",
			expectedError: true,
			expectedFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			action := NewArkServiceExecAction(nil)
			method, err := action.findMethodByName(tt.value, tt.methodName)

			if tt.expectedError && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if tt.expectedFound && (method == nil || !method.IsValid()) {
				t.Error("Expected to find method, but didn't")
			}
			if !tt.expectedFound && method != nil && method.IsValid() {
				t.Error("Expected not to find method, but did")
			}
		})
	}
}

func TestArkServiceExecAction_DefineExecAction(t *testing.T) {
	tests := []struct {
		name          string
		setupCmd      func() *cobra.Command
		expectedError bool
		validateFunc  func(t *testing.T, cmd *cobra.Command)
	}{
		{
			name: "success_defines_actions_without_error",
			setupCmd: func() *cobra.Command {
				return &cobra.Command{Use: "test"}
			},
			expectedError: false,
			validateFunc: func(t *testing.T, cmd *cobra.Command) {
				// Check that subcommands were added (depending on services.SupportedServiceActions)
				if cmd.Commands() == nil {
					// This may be expected if no supported service actions are defined
					t.Log("No commands added - this may be expected if services.SupportedServiceActions is empty")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			action := NewArkServiceExecAction(nil)
			cmd := tt.setupCmd()

			err := action.DefineExecAction(cmd)

			if tt.expectedError && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if !tt.expectedError && tt.validateFunc != nil {
				tt.validateFunc(t, cmd)
			}
		})
	}
}

func TestArkServiceExecAction_serializeAndPrintOutput(t *testing.T) {
	tests := []struct {
		name       string
		result     []reflect.Value
		actionName string
		// Note: This function prints to console, so we mainly test that it doesn't panic
		shouldPanic bool
	}{
		{
			name: "success_handles_struct_output",
			result: []reflect.Value{
				reflect.ValueOf(map[string]interface{}{"key": "value"}),
			},
			actionName:  "test-action",
			shouldPanic: false,
		},
		{
			name: "success_handles_int_output",
			result: []reflect.Value{
				reflect.ValueOf(42),
			},
			actionName:  "test-action",
			shouldPanic: false,
		},
		{
			name: "success_handles_string_output",
			result: []reflect.Value{
				reflect.ValueOf("test result"),
			},
			actionName:  "test-action",
			shouldPanic: false,
		},
		{
			name:        "success_handles_empty_result",
			result:      []reflect.Value{},
			actionName:  "test-action",
			shouldPanic: false,
		},
		{
			name: "success_handles_nil_pointer",
			result: []reflect.Value{
				reflect.ValueOf((*string)(nil)),
			},
			actionName:  "test-action",
			shouldPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			action := NewArkServiceExecAction(nil)

			defer func() {
				if r := recover(); r != nil && !tt.shouldPanic {
					t.Errorf("Function panicked unexpectedly: %v", r)
				} else if r == nil && tt.shouldPanic {
					t.Error("Expected function to panic, but it didn't")
				}
			}()

			action.serializeAndPrintOutput(tt.result, tt.actionName)
		})
	}
}

func TestArkServiceExecAction_defineServiceExecAction(t *testing.T) {
	tests := []struct {
		name             string
		actionDef        *actions.ArkServiceActionDefinition
		parentActionsDef []*actions.ArkServiceActionDefinition
		expectedError    bool
		validateFunc     func(t *testing.T, cmd *cobra.Command, err error)
	}{
		{
			name: "success_creates_simple_action",
			actionDef: &actions.ArkServiceActionDefinition{
				ActionName: "test-action",
				Schemas: map[string]interface{}{
					"execute": testutils.CreateTestSchema(),
				},
			},
			parentActionsDef: nil,
			expectedError:    false,
			validateFunc: func(t *testing.T, cmd *cobra.Command, err error) {
				if cmd == nil {
					t.Error("Expected command to be created")
					return
				}
				if cmd.Use != "test-action" {
					t.Errorf("Expected command name 'test-action', got '%s'", cmd.Use)
				}
			},
		},
		{
			name: "success_creates_action_with_parent",
			actionDef: &actions.ArkServiceActionDefinition{
				ActionName: "sub-action",
				Schemas:    map[string]interface{}{"execute": testutils.CreateTestSchema()},
			},
			parentActionsDef: []*actions.ArkServiceActionDefinition{
				{ActionName: "parent-action"},
			},
			expectedError: false,
			validateFunc: func(t *testing.T, cmd *cobra.Command, err error) {
				if cmd == nil {
					t.Error("Expected command to be created")
					return
				}
				if cmd.Use != "sub-action" {
					t.Errorf("Expected command name 'sub-action', got '%s'", cmd.Use)
				}
			},
		},
		{
			name: "success_creates_action_without_schemas",
			actionDef: &actions.ArkServiceActionDefinition{
				ActionName: "no-schema-action",
				Schemas:    map[string]interface{}{},
			},
			parentActionsDef: nil,
			expectedError:    false,
			validateFunc: func(t *testing.T, cmd *cobra.Command, err error) {
				if cmd == nil {
					t.Error("Expected command to be created")
					return
				}
				if len(cmd.Commands()) != 0 {
					t.Error("Expected no subcommands for action without schemas")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			action := NewArkServiceExecAction(nil)
			parentCmd := &cobra.Command{Use: "parent"}

			resultCmd, err := action.defineServiceExecAction(tt.actionDef, parentCmd, tt.parentActionsDef)

			if tt.expectedError && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if tt.validateFunc != nil {
				tt.validateFunc(t, resultCmd, err)
			}
		})
	}
}

// Helper function to create a test action definition
func createTestActionDefinition(name string, hasSubactions bool) *actions.ArkServiceActionDefinition {
	actionDef := &actions.ArkServiceActionDefinition{
		ActionName: name,
		Schemas: map[string]interface{}{
			"execute": testutils.CreateTestSchema(),
		},
	}

	if hasSubactions {
		actionDef.Subactions = []*actions.ArkServiceActionDefinition{
			{
				ActionName: "sub-" + name,
				Schemas: map[string]interface{}{
					"sub-execute": testutils.CreateTestSchema(),
				},
			},
		}
	}

	return actionDef
}

func TestArkServiceExecAction_defineServiceExecActions(t *testing.T) {
	tests := []struct {
		name             string
		actionDef        *actions.ArkServiceActionDefinition
		parentActionsDef []*actions.ArkServiceActionDefinition
		expectedError    bool
		validateFunc     func(t *testing.T, cmd *cobra.Command, err error)
	}{
		{
			name:             "success_defines_action_without_subactions",
			actionDef:        createTestActionDefinition("simple", false),
			parentActionsDef: nil,
			expectedError:    false,
			validateFunc: func(t *testing.T, cmd *cobra.Command, err error) {
				if len(cmd.Commands()) == 0 {
					t.Error("Expected at least one command to be added")
				}
			},
		},
		{
			name:             "success_defines_action_with_subactions",
			actionDef:        createTestActionDefinition("complex", true),
			parentActionsDef: nil,
			expectedError:    false,
			validateFunc: func(t *testing.T, cmd *cobra.Command, err error) {
				if len(cmd.Commands()) == 0 {
					t.Error("Expected at least one command to be added")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			action := NewArkServiceExecAction(nil)
			cmd := &cobra.Command{Use: "parent"}

			err := action.defineServiceExecActions(tt.actionDef, cmd, tt.parentActionsDef)

			if tt.expectedError && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if tt.validateFunc != nil {
				tt.validateFunc(t, cmd, err)
			}
		})
	}
}
