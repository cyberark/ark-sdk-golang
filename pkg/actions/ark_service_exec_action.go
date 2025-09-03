package actions

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"slices"
	"strings"

	"github.com/cyberark/ark-sdk-golang/pkg/common"

	"github.com/cyberark/ark-sdk-golang/pkg/cli"
	"github.com/cyberark/ark-sdk-golang/pkg/common/args"
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions"
	"github.com/cyberark/ark-sdk-golang/pkg/models/actions/services"
	"github.com/cyberark/ark-sdk-golang/pkg/profiles"
	"github.com/mitchellh/mapstructure"
	"github.com/octago/sflags"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// ArkServiceExecAction is a struct that implements the ArkExecAction interface for executing service actions.
//
// ArkServiceExecAction provides functionality for dynamically executing service actions
// based on service action definitions. It handles the parsing of command-line flags,
// schema validation, method invocation, and output serialization for service operations.
//
// The action supports:
//   - Dynamic command generation from service action definitions
//   - Complex type parsing for JSON and array inputs
//   - Flag validation with choices and required field checking
//   - Method reflection and invocation on service APIs
//   - Multiple output format serialization (JSON, primitive types, channels)
//   - Request file input support for complex payloads
type ArkServiceExecAction struct {
	// ArkExecAction interface for execution capabilities
	ArkExecAction
	// ArkBaseExecAction provides common execution functionality
	*ArkBaseExecAction
}

// NewArkServiceExecAction creates a new instance of ArkServiceExecAction.
//
// NewArkServiceExecAction initializes a new ArkServiceExecAction with the provided
// profile loader and embedded ArkBaseExecAction for common execution functionality.
// The action is configured with reflection-based method invocation capabilities
// for dynamic service action execution.
//
// Parameters:
//   - profilesLoader: A pointer to a ProfileLoader for handling profile operations
//
// Returns a new ArkServiceExecAction instance ready for defining and executing
// service commands.
//
// Example:
//
//	loader := profiles.NewProfileLoader()
//	serviceExecAction := NewArkServiceExecAction(loader)
//	serviceExecAction.DefineExecAction(rootCmd)
func NewArkServiceExecAction(profilesLoader *profiles.ProfileLoader) *ArkServiceExecAction {
	action := &ArkServiceExecAction{}
	var actionInterface ArkExecAction = action
	baseAction := NewArkBaseExecAction(&actionInterface, "ArkServiceExecAction", profilesLoader)
	action.ArkBaseExecAction = baseAction
	return action
}

// isComplexType determines if a struct field represents a complex type requiring JSON parsing.
//
// isComplexType checks if the field is a map[string]struct or slice/array of structs,
// which require special handling during flag parsing as they need to be parsed from
// JSON strings rather than simple flag values.
//
// Parameters:
//   - field: The reflect.StructField to check for complexity
//
// Returns true if the field is a complex type (map[string]struct or []struct),
// false otherwise.
func (s *ArkServiceExecAction) isComplexType(field reflect.StructField) bool {
	if field.Type.Kind() == reflect.Map && field.Type.Key().Kind() == reflect.String && field.Type.Elem().Kind() == reflect.Struct {
		return true
	}
	if (field.Type.Kind() == reflect.Slice || field.Type.Kind() == reflect.Array) && field.Type.Elem().Kind() == reflect.Struct {
		return true
	}
	return false
}

// fillRemainingSchema adds flags for complex types and squashed struct fields.
//
// fillRemainingSchema processes a schema struct and adds command-line flags for
// complex types (maps and slices of structs) that require JSON parsing, and
// recursively processes squashed struct fields to flatten their fields into
// the parent command's flag set.
//
// Parameters:
//   - schema: The schema interface to process for flag generation
//   - flags: The pflag.FlagSet to add the generated flags to
//
// The function handles:
//   - Complex types by adding string flags with JSON parsing hints
//   - Squashed struct fields by recursively processing embedded structs
//   - Flag naming from struct tags (flag, mapstructure, or field name)
//   - Description enhancement for complex types
func (s *ArkServiceExecAction) fillRemainingSchema(schema interface{}, flags *pflag.FlagSet) {
	schemaType := reflect.TypeOf(schema).Elem()
	for i := 0; i < schemaType.NumField(); i++ {
		field := schemaType.Field(i)

		// Skip unexported fields
		if field.PkgPath != "" && !field.Anonymous {
			continue
		}

		if s.isComplexType(field) {
			flagName := field.Tag.Get("flag")
			if flagName == "" {
				flagName = field.Tag.Get("mapstructure")
			}
			if flagName == "" {
				flagName = field.Name
			}
			desc := field.Tag.Get("desc")
			if desc != "" {
				desc += " (This is a complex type and will be parsed as JSON or array of JSONs)"
			}
			flags.String(flagName, field.Tag.Get("default"), desc)
		}
		if field.Tag.Get("mapstructure") == ",squash" {
			// If the field is a struct with the `squash` tag, we need to add its fields as flags
			subSchema := reflect.New(field.Type).Interface()
			s.fillRemainingSchema(subSchema, flags)
		}
	}
}

// defineServiceExecAction creates a cobra command for a service action definition.
//
// defineServiceExecAction processes a service action definition and creates the
// corresponding cobra command with subcommands for each schema. It handles flag
// generation, validation, and default value assignment based on struct tags.
//
// Parameters:
//   - actionDef: The service action definition to process
//   - cmd: The parent cobra command to add the action command to
//   - parentActionsDef: Slice of parent action definitions for nested actions
//
// Returns the created action command and any error encountered during processing.
//
// The function handles:
//   - Command creation with proper naming from action definitions
//   - Schema-based subcommand generation
//   - Flag parsing using sflags library
//   - Required field marking based on validation tags
//   - Default value assignment from struct tags
func (s *ArkServiceExecAction) defineServiceExecAction(
	actionDef *actions.ArkServiceActionDefinition,
	cmd *cobra.Command,
	parentActionsDef []*actions.ArkServiceActionDefinition,
) (*cobra.Command, error) {
	actionCmd := &cobra.Command{
		Use: actionDef.ActionName,
	}

	actionDest := actionDef.ActionName
	if len(parentActionsDef) > 0 {
		for _, p := range parentActionsDef {
			actionDest += "_" + p.ActionName
		}
		actionDest += "_" + actionDef.ActionName
	}

	if len(actionDef.Schemas) > 0 {
		for actionName, schema := range actionDef.Schemas {
			subCmd := &cobra.Command{
				Use: actionName,
				Run: func(cmd *cobra.Command, args []string) {
					if help, _ := cmd.Flags().GetBool("help"); help {
						_ = cmd.Help()
						return
					}
					s.runExecAction(cmd, args)
				},
			}
			if schema != nil {
				flags, err := sflags.ParseStruct(schema)
				if err != nil {
					s.logger.Error("Error parsing flags to ArkAuthMethod settings %v", err)
					return nil, err
				}
				gpflag.GenerateTo(flags, subCmd.Flags())
				s.fillRemainingSchema(schema, subCmd.Flags())
				reflectedSchema := reflect.TypeOf(schema).Elem()
				// We find the field by the flag name
				// There might be a misalignment between the flag name and the field name case wise
				// So we first try to find the field by the flag name
				// And then try to find it with ignore case
				for _, flag := range flags {
					flagNameTitled := strings.Replace(strings.Title(flag.Name), "-", "", -1)
					field, ok := reflectedSchema.FieldByName(flagNameTitled)
					if !ok {
						fieldFound := false
						for i := 0; i < reflectedSchema.NumField(); i++ {
							possibleField := reflectedSchema.Field(i)
							if strings.EqualFold(possibleField.Name, flagNameTitled) {
								field = possibleField
								fieldFound = true
								break
							}
						}
						if !fieldFound {
							continue
						}
					}
					if strings.Contains(field.Tag.Get("validate"), "required") {
						err = subCmd.MarkFlagRequired(flag.Name)
						if err != nil {
							return nil, err
						}
					}
					if field.Tag.Get("default") != "" {
						subCmd.Flag(flag.Name).DefValue = field.Tag.Get("default")
						err = subCmd.Flag(flag.Name).Value.Set(field.Tag.Get("default"))
						if err != nil {
							return nil, err
						}
					}
				}
			}
			actionCmd.AddCommand(subCmd)
		}
	}

	cmd.AddCommand(actionCmd)
	return actionCmd, nil
}

// defineServiceExecActions recursively defines service execution actions and their subactions.
//
// defineServiceExecActions processes a service action definition and its nested
// subactions, creating a hierarchy of cobra commands. It recursively processes
// subactions to build a complete command tree structure.
//
// Parameters:
//   - actionDef: The service action definition to process
//   - cmd: The parent cobra command to add actions to
//   - parentActionsDef: Slice of parent action definitions for context
//
// Returns an error if any action definition processing fails.
//
// The function handles:
//   - Primary action definition processing through defineServiceExecAction
//   - Recursive subaction processing for nested command structures
//   - Error propagation from nested action creation
func (s *ArkServiceExecAction) defineServiceExecActions(
	actionDef *actions.ArkServiceActionDefinition,
	cmd *cobra.Command,
	parentActionsDef []*actions.ArkServiceActionDefinition,
) error {
	actionSubparsers, err := s.defineServiceExecAction(actionDef, cmd, parentActionsDef)
	if err != nil {
		return err
	}
	if len(actionDef.Subactions) > 0 {
		for _, subaction := range actionDef.Subactions {
			err = s.defineServiceExecActions(subaction, actionSubparsers, append(parentActionsDef, actionDef))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// fillParsedFlag processes complex flag values and validates choices for schema fields.
//
// fillParsedFlag handles the parsing of complex types (JSON objects and arrays)
// from string flag values and validates that field values match defined choices
// constraints. It processes mapstructure tags to find matching schema fields
// and applies appropriate transformations and validations.
//
// Parameters:
//   - schemaElem: The reflect.Type of the schema struct to process
//   - flags: Map of flag names to values being processed
//   - key: The specific flag key being processed
//   - f: The pflag.Flag being processed for error reporting
//
// Returns an error if JSON parsing fails or choice validation fails.
//
// The function handles:
//   - JSON unmarshaling for complex map and slice types
//   - Choice validation for string, string slice, and map types
//   - Recursive processing of squashed struct fields
//   - Error reporting with context about which flag failed
func (s *ArkServiceExecAction) fillParsedFlag(schemaElem reflect.Type, flags map[string]interface{}, key string, f *pflag.Flag) error {
	for i := 0; i < schemaElem.NumField(); i++ {
		field := schemaElem.Field(i)
		if strings.HasPrefix(field.Tag.Get("mapstructure"), key) {
			if s.isComplexType(field) {
				if field.Type.Kind() == reflect.Map && field.Type.Key().Kind() == reflect.String && field.Type.Elem().Kind() == reflect.Struct {
					var mapJSON map[string]interface{}
					err := json.Unmarshal([]byte(flags[key].(string)), &mapJSON)
					if err != nil {
						return err
					}
					flags[key] = mapJSON
				} else {
					var sliceJSON []map[string]interface{}
					err := json.Unmarshal([]byte(flags[key].(string)), &sliceJSON)
					if err != nil {
						return err
					}
					flags[key] = sliceJSON
				}
			}
			if field.Tag.Get("choices") != "" {
				choices := strings.Split(field.Tag.Get("choices"), ",")
				switch v := flags[key].(type) {
				case string:
					if !slices.Contains(choices, v) {
						return fmt.Errorf("invalid value for flag %s: %s, valid choices are: %s", f.Name, v, strings.Join(choices, ", "))
					}
				case []string:
					for _, item := range v {
						if !slices.Contains(choices, item) {
							return fmt.Errorf("invalid value for flag %s: %s, valid choices are: %s", f.Name, item, strings.Join(choices, ", "))
						}
					}
				case map[string]any:
					for fieldKey := range v {
						if !slices.Contains(choices, fieldKey) {
							return fmt.Errorf("invalid key for flag %s: %s, valid choices are: %s", f.Name, fieldKey, strings.Join(choices, ", "))
						}
					}
				default:
					return fmt.Errorf("unexpected type for flag %s: %T", f.Name, flags[key])
				}
			}
		} else if field.Tag.Get("mapstructure") == ",squash" {
			err := s.fillParsedFlag(field.Type, flags, key, f)
			if err != nil {
				return err
			}
			continue
		}
	}
	return nil
}

// parseFlag extracts and converts flag values to appropriate types for schema processing.
//
// parseFlag handles the extraction of flag values from cobra commands and converts
// them to the appropriate Go types for later processing by mapstructure. It supports
// all common Go primitive types and collections, then applies complex type processing
// and choice validation through fillParsedFlag.
//
// Parameters:
//   - f: The pflag.Flag to process
//   - cmd: The cobra.Command containing the flag values
//   - flags: Map to store the parsed flag values
//   - schema: The schema interface for validation and complex type processing
//
// Returns an error if flag parsing or validation fails.
//
// The function handles:
//   - Type-specific flag value extraction (bool, int variants, float variants, slices, maps)
//   - Conversion of flag names from kebab-case to snake_case
//   - Delegation to fillParsedFlag for complex type processing and validation
//   - Skipping unchanged flags to avoid unnecessary processing
func (s *ArkServiceExecAction) parseFlag(f *pflag.Flag, cmd *cobra.Command, flags map[string]interface{}, schema interface{}) error {
	if !f.Changed {
		return nil
	}
	key := strings.ReplaceAll(f.Name, "-", "_")
	switch f.Value.Type() {
	case "bool":
		val, err := cmd.Flags().GetBool(f.Name)
		if err == nil {
			flags[key] = val
		}
	case "int":
		val, err := cmd.Flags().GetInt(f.Name)
		if err == nil {
			flags[key] = val
		}
	case "int8":
		val, err := cmd.Flags().GetInt8(f.Name)
		if err == nil {
			flags[key] = val
		}
	case "int16":
		val, err := cmd.Flags().GetInt16(f.Name)
		if err == nil {
			flags[key] = val
		}
	case "int32":
		val, err := cmd.Flags().GetInt32(f.Name)
		if err == nil {
			flags[key] = val
		}
	case "int64":
		val, err := cmd.Flags().GetInt64(f.Name)
		if err == nil {
			flags[key] = val
		}
	case "uint":
		val, err := cmd.Flags().GetUint(f.Name)
		if err == nil {
			flags[key] = val
		}
	case "uint8":
		val, err := cmd.Flags().GetUint8(f.Name)
		if err == nil {
			flags[key] = val
		}
	case "uint16":
		val, err := cmd.Flags().GetUint16(f.Name)
		if err == nil {
			flags[key] = val
		}
	case "uint32":
		val, err := cmd.Flags().GetUint32(f.Name)
		if err == nil {
			flags[key] = val
		}
	case "uint64":
		val, err := cmd.Flags().GetUint64(f.Name)
		if err == nil {
			flags[key] = val
		}
	case "float32":
		val, err := cmd.Flags().GetFloat32(f.Name)
		if err == nil {
			flags[key] = val
		}
	case "float64":
		val, err := cmd.Flags().GetFloat64(f.Name)
		if err == nil {
			flags[key] = val
		}
	case "stringSlice":
	case "[]string":
		val, err := cmd.Flags().GetStringSlice(f.Name)
		if err == nil {
			flags[key] = val
		}
	case "stringArray":
		val, err := cmd.Flags().GetStringArray(f.Name)
		if err == nil {
			flags[key] = val
		}
	case "intSlice":
	case "[]int":
		val, err := cmd.Flags().GetIntSlice(f.Name)
		if err == nil {
			flags[key] = val
		}
	case "stringToString":
	case "map[string]string":
		val, err := cmd.Flags().GetStringToString(f.Name)
		if err == nil {
			flags[key] = val
		}
	default:
		flags[key] = f.Value.String()
	}
	schemaElem := reflect.TypeOf(schema).Elem()
	return s.fillParsedFlag(schemaElem, flags, key, f)
}

// serializeAndPrintOutput formats and displays the results of service action execution.
//
// serializeAndPrintOutput processes the reflection values returned from service
// method execution and formats them appropriately for console output. It handles
// various result types including structs, maps, arrays, channels, and primitive types.
//
// Parameters:
//   - result: Slice of reflect.Value containing the method execution results
//   - actionName: The name of the action being executed (for generic success messages)
//
// The function handles:
//   - JSON serialization for complex types (structs, maps, arrays, slices)
//   - Channel processing for paginated results with Items field extraction
//   - Integer formatting for numeric results
//   - Generic success messages when no specific output is available
//   - Error handling for JSON serialization failures with fallback output
func (s *ArkServiceExecAction) serializeAndPrintOutput(result []reflect.Value, actionName string) {
	shouldPrintGenericResult := true
	for _, res := range result {
		if res.Kind() == reflect.Ptr && res.IsNil() {
			continue
		}
		if res.Kind() == reflect.Interface && res.Type().Implements(reflect.TypeOf((*error)(nil)).Elem()) {
			continue
		}
		if res.Kind() == reflect.Ptr {
			res = res.Elem()
		}
		if res.Kind() == reflect.Struct || res.Kind() == reflect.Map || res.Kind() == reflect.Array || res.Kind() == reflect.Slice {
			jsonData, err := json.MarshalIndent(res.Interface(), "", "  ")
			if err != nil {
				s.logger.Warning("error serializing result to JSON: %v", err)
				args.PrintSuccess(res.Interface())
			} else {
				args.PrintSuccess(string(jsonData))
			}
			shouldPrintGenericResult = false
		} else if res.Kind() == reflect.Chan {
			items := make([]interface{}, 0)
			for {
				pageValue, ok := res.Recv()
				if !ok {
					break
				}
				if !pageValue.IsValid() {
					continue
				}
				if pageValue.Kind() == reflect.Ptr {
					pageValue = pageValue.Elem()
				}
				itemsField := pageValue.FieldByName("Items")
				if !itemsField.IsValid() || itemsField.Kind() != reflect.Slice {
					items = append(items, pageValue.Interface())
					continue
				}
				for i := 0; i < itemsField.Len(); i++ {
					items = append(items, itemsField.Index(i).Interface())
				}
			}
			jsonData, err := json.MarshalIndent(items, "", "  ")
			if err != nil {
				s.logger.Warning("error serializing result to JSON: %v", err)
				args.PrintSuccess(items)
			} else {
				args.PrintSuccess(string(jsonData))
			}
			shouldPrintGenericResult = false
		} else if res.Kind() == reflect.Int {
			args.PrintSuccess(fmt.Sprintf("%d", res.Int()))
			shouldPrintGenericResult = false
		} else {
			args.PrintSuccess(res.Interface())
			shouldPrintGenericResult = false
		}
	}
	if len(result) == 0 || shouldPrintGenericResult {
		args.PrintSuccess(fmt.Sprintf("%s finished successfully", strings.Replace(strings.Title(actionName), "-", " ", -1)))
	}
}

// findMethodByName locates a method on a reflect.Value using case-insensitive matching.
//
// findMethodByName searches for a method by name on the provided reflection value,
// first attempting an exact match and then falling back to case-insensitive matching
// if the exact match fails. This provides flexibility for method name variations.
//
// Parameters:
//   - value: The reflect.Value to search for methods on
//   - methodName: The name of the method to find
//
// Returns a pointer to the reflect.Value representing the method and any error
// encountered during the search.
//
// The function handles:
//   - Exact method name matching first
//   - Case-insensitive fallback matching through all available methods
//   - Error reporting when no matching method is found
func (s *ArkServiceExecAction) findMethodByName(value reflect.Value, methodName string) (*reflect.Value, error) {
	actionMethod := value.MethodByName(methodName)
	if !actionMethod.IsValid() {
		for i := 0; i < value.NumMethod(); i++ {
			method := value.Type().Method(i)
			if strings.EqualFold(method.Name, methodName) {
				actionMethod = value.MethodByName(method.Name)
				break
			}
		}
		if !actionMethod.IsValid() {
			return nil, fmt.Errorf("method %s not found", methodName)
		}
	}
	return &actionMethod, nil
}

// DefineExecAction defines the execution actions for all supported service operations.
//
// DefineExecAction processes all supported service action definitions and creates
// the corresponding cobra command hierarchy for service execution. It iterates through
// the available service actions and creates the complete command structure for
// dynamic service operation execution.
//
// Parameters:
//   - cmd: The parent cobra command to add service execution commands to
//
// Returns an error if any service action definition processing fails.
//
// The function handles:
//   - Processing all supported service actions from the services package
//   - Creating command hierarchies for each service action through defineServiceExecActions
//   - Error propagation from nested action processing
//
// Example:
//
//	err := serviceExecAction.DefineExecAction(rootCmd)
//	// This adds all supported service commands to rootCmd
func (s *ArkServiceExecAction) DefineExecAction(cmd *cobra.Command) error {
	for _, actionDef := range services.SupportedServiceActions {
		err := s.defineServiceExecActions(actionDef, cmd, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

// RunExecAction executes a service action using reflection-based method invocation.
//
// RunExecAction processes the command hierarchy to determine the target service and action,
// then uses reflection to locate and invoke the appropriate method on the API service.
// It handles flag parsing, schema validation, method resolution, and output formatting
// for dynamic service action execution.
//
// Parameters:
//   - api: The ArkCLIAPI instance containing the service methods
//   - cmd: The cobra command being executed
//   - execCmd: The parent execution command for context
//   - execArgs: Command line arguments for the execution
//
// Returns an error if service resolution, method invocation, or parameter processing fails.
//
// The function handles:
//   - Service path resolution from command hierarchy
//   - Method name transformation and case-insensitive lookup
//   - Schema resolution from service action definitions
//   - Flag parsing and validation against schema constraints
//   - Request file input for complex payloads
//   - Method invocation with appropriate parameters
//   - Result serialization and output formatting
//
// Example:
//
//	err := serviceExecAction.RunExecAction(api, cmd, execCmd, args)
//	// Executes the service method and displays formatted output
func (s *ArkServiceExecAction) RunExecAction(api *cli.ArkCLIAPI, cmd *cobra.Command, execCmd *cobra.Command, execArgs []string) error {
	serviceParts := make([]string, 0)
	for currentCmd := cmd.Parent(); currentCmd != execCmd; currentCmd = currentCmd.Parent() {
		serviceParts = append([]string{currentCmd.Name()}, serviceParts...)
	}
	actionName := cmd.Name()
	actionNameTitled := strings.Replace(strings.Title(actionName), "-", "", -1)
	serviceNameTitled := ""
	for _, part := range serviceParts {
		serviceNameTitled += strings.Title(part)
	}
	serviceNameTitled = strings.Replace(strings.Title(serviceNameTitled), "-", "", -1)
	// First, resolve the action method
	serviceMethod, err := s.findMethodByName(reflect.ValueOf(api), serviceNameTitled)
	if err != nil {
		return err
	}
	serviceErr := serviceMethod.Call(nil)
	service := serviceErr[0]
	if len(serviceErr) > 1 {
		if err, ok := serviceErr[1].Interface().(error); ok && err != nil {
			return err
		}
	}
	actionMethod, err := s.findMethodByName(reflect.ValueOf(service.Interface()), actionNameTitled)
	if err != nil {
		return err
	}

	// Resolve the action schema
	var actionSchemaDef *actions.ArkServiceActionDefinition = nil
	for _, servicePart := range serviceParts {
		if actionSchemaDef != nil {
			for _, actionDef := range actionSchemaDef.Subactions {
				if actionDef.ActionName == servicePart {
					actionSchemaDef = actionDef
					break
				}
			}
		} else {
			for _, actionDef := range services.SupportedServiceActions {
				if actionDef.ActionName == servicePart {
					actionSchemaDef = actionDef
					break
				}
			}
			if actionSchemaDef == nil {
				return fmt.Errorf("action %s not found in service %s", actionName, serviceNameTitled)
			}
		}
	}
	actionSchema, ok := actionSchemaDef.Schemas[actionName]
	if !ok {
		return fmt.Errorf("action %s not supported", actionName)
	}
	var result []reflect.Value
	if actionSchema != nil {
		flags := map[string]interface{}{}
		if requestFile, err := execCmd.PersistentFlags().GetString("request-file"); err == nil && requestFile != "" {
			fileContent, err := os.ReadFile(requestFile)
			if err != nil {
				return err
			}
			var data map[string]interface{}
			err = json.Unmarshal(fileContent, &data)
			if err != nil {
				return err
			}
			schemaType := reflect.ValueOf(actionSchema).Type()
			flags = common.ConvertToSnakeCase(data, &schemaType).(map[string]interface{})
		}
		err = nil
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			err = s.parseFlag(f, cmd, flags, actionSchema)
		})
		if err != nil {
			return err
		}
		err = mapstructure.Decode(flags, actionSchema)
		if err != nil {
			return err
		}
		actionArgs := []reflect.Value{reflect.ValueOf(actionSchema)}
		result = actionMethod.Call(actionArgs)
	} else {
		var actionArgs []reflect.Value
		result = actionMethod.Call(actionArgs)
	}
	for _, res := range result {
		if err, ok := res.Interface().(error); ok && err != nil {
			return err
		}
	}

	s.serializeAndPrintOutput(result, actionName)

	return nil
}
