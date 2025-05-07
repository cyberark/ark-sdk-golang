package actions

import (
	"encoding/json"
	"fmt"
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
	"reflect"
	"slices"
	"strings"
)

// ArkServiceExecAction is a struct that implements the ArkExecAction interface for executing service actions.
type ArkServiceExecAction struct {
	ArkExecAction
	*ArkBaseExecAction
}

// NewArkServiceExecAction creates a new instance of ArkServiceExecAction.
func NewArkServiceExecAction(profilesLoader *profiles.ProfileLoader) *ArkServiceExecAction {
	action := &ArkServiceExecAction{}
	var actionInterface ArkExecAction = action
	baseAction := NewArkBaseExecAction(&actionInterface, "ArkServiceExecAction", profilesLoader)
	action.ArkBaseExecAction = baseAction
	return action
}

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
		val, err := cmd.Flags().GetIntSlice(f.Name)
		if err == nil {
			flags[key] = val
		}
	case "stringToString":
		val, err := cmd.Flags().GetStringToString(f.Name)
		if err == nil {
			flags[key] = val
		}
	default:
		flags[key] = f.Value.String()
	}
	schemaElem := reflect.TypeOf(schema).Elem()
	for i := 0; i < schemaElem.NumField(); i++ {
		field := schemaElem.Field(i)
		if strings.HasPrefix(field.Tag.Get("mapstructure"), key) {
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
				default:
					return fmt.Errorf("unexpected type for flag %s: %T", f.Name, flags[key])
				}
			}
		}
	}
	return nil
}

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
		} else {
			args.PrintSuccess(res.Interface())
			shouldPrintGenericResult = false
		}
	}
	if len(result) == 0 || shouldPrintGenericResult {
		args.PrintSuccess(fmt.Sprintf("%s finished successfully", strings.Replace(strings.Title(actionName), "-", " ", -1)))
	}
}

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

// DefineExecAction defines the exec action for the ArkServiceExecAction.
func (s *ArkServiceExecAction) DefineExecAction(cmd *cobra.Command) error {
	for _, actionDef := range services.SupportedServiceActions {
		err := s.defineServiceExecActions(actionDef, cmd, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

// RunExecAction runs the exec action for the ArkServiceExecAction.
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
