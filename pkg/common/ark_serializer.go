package common

import (
	"encoding/json"
	"github.com/iancoleman/strcase"
	"io"
	"reflect"
	"strings"
)

func resolveFieldsSquashed(schema reflect.Type) []reflect.StructField {
	var fields []reflect.StructField
	if schema.Kind() == reflect.Ptr {
		schema = schema.Elem()
	}
	for i := 0; i < schema.NumField(); i++ {
		field := schema.Field(i)
		if field.Tag.Get("mapstructure") == ",squash" {
			nestedFields := resolveFieldsSquashed(field.Type)
			fields = append(fields, nestedFields...)
			continue
		}
		if field.PkgPath != "" { // unexported field
			continue
		}
		fields = append(fields, field)
	}
	return fields
}

func findFieldByName(schema reflect.Type, name string) *reflect.StructField {
	flagNameTitled := strings.Replace(strings.Replace(strings.Title(name), "-", "", -1), "_", "", -1)
	if schema.Kind() == reflect.Ptr {
		schema = schema.Elem()
	}
	field, ok := schema.FieldByName(flagNameTitled)
	if ok {
		return &field
	}
	actualFields := resolveFieldsSquashed(schema)
	for i := 0; i < len(actualFields); i++ {
		possibleField := actualFields[i]
		if strings.EqualFold(possibleField.Name, flagNameTitled) {
			return &possibleField
		}
	}
	return nil
}

// SerializeResponseToJSON takes an io.ReadCloser response and serializes it to a JSON string.
func SerializeResponseToJSON(response io.ReadCloser) string {
	data, err := io.ReadAll(response)
	if err != nil {
		return ""
	}
	jsonMap := make(map[string]interface{})
	err = json.Unmarshal(data, &jsonMap)
	if err != nil {
		return string(data)
	}
	jsonData, err := json.Marshal(jsonMap)
	if err != nil {
		return string(data)
	}
	return string(jsonData)
}

// ConvertToSnakeCase converts a map with camelCase keys to snake_case keys.
func ConvertToSnakeCase(data interface{}, schema *reflect.Type) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		snakeMap := make(map[string]interface{})
		for key, value := range v {
			snakeKey := strcase.ToSnake(key)
			var innerFieldType *reflect.Type
			if schema != nil {
				innerField := findFieldByName(*schema, key)
				if innerField == nil || (innerField.Type.Kind() == reflect.Map && innerField.Type.Key().Kind() == reflect.String) {
					snakeKey = key
				}
				if innerField != nil {
					if innerField.Type.Kind() == reflect.Slice || innerField.Type.Kind() == reflect.Array || innerField.Type.Kind() == reflect.Map {
						elem := innerField.Type.Elem()
						innerFieldType = &elem
					} else {
						innerFieldType = &innerField.Type
					}
				}
			}
			snakeMap[snakeKey] = ConvertToSnakeCase(value, innerFieldType)
		}
		return snakeMap
	case []interface{}:
		for i, item := range v {
			v[i] = ConvertToSnakeCase(item, schema)
		}
		return v
	default:
		return v
	}
}

// ConvertToCamelCase converts a map with snake_case keys to camelCase keys.
func ConvertToCamelCase(data interface{}, schema *reflect.Type) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		camelMap := make(map[string]interface{})
		for key, value := range v {
			camelKey := strcase.ToLowerCamel(key)
			var innerFieldType *reflect.Type
			if schema != nil {
				innerField := findFieldByName(*schema, key)
				if innerField == nil || (innerField.Type.Kind() == reflect.Map && innerField.Type.Key().Kind() == reflect.String) {
					camelKey = key
				}
				if innerField != nil {
					if innerField.Type.Kind() == reflect.Slice || innerField.Type.Kind() == reflect.Array || innerField.Type.Kind() == reflect.Map {
						elem := innerField.Type.Elem()
						innerFieldType = &elem
					} else {
						innerFieldType = &innerField.Type
					}
				}
			}
			camelMap[camelKey] = ConvertToCamelCase(value, innerFieldType)
		}
		return camelMap
	case []interface{}:
		for i, item := range v {
			v[i] = ConvertToCamelCase(item, schema)
		}
		return v
	default:
		return v
	}
}

// DeserializeJSONSnake takes an io.ReadCloser response and deserializes it into a map with snake_case keys.
func DeserializeJSONSnake(response io.ReadCloser) (interface{}, error) {
	var result interface{}
	err := json.NewDecoder(response).Decode(&result)
	if err != nil {
		return nil, err
	}
	return ConvertToSnakeCase(result, nil).(interface{}), nil
}

// SerializeJSONCamel takes an interface and serializes it into a map with camelCase keys.
func SerializeJSONCamel(item interface{}) (map[string]interface{}, error) {
	resultBytes, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(resultBytes, &result)
	if err != nil {
		return nil, err
	}
	return ConvertToCamelCase(result, nil).(map[string]interface{}), nil
}

// DeserializeJSONSnakeSchema takes an io.ReadCloser response and deserializes it into a map with snake_case keys.
func DeserializeJSONSnakeSchema(response io.ReadCloser, schema *reflect.Type) (interface{}, error) {
	var result interface{}
	err := json.NewDecoder(response).Decode(&result)
	if err != nil {
		return nil, err
	}
	return ConvertToSnakeCase(result, schema).(interface{}), nil
}

// SerializeJSONCamelSchema takes an interface and serializes it into a map with camelCase keys.
func SerializeJSONCamelSchema(item interface{}, schema *reflect.Type) (map[string]interface{}, error) {
	resultBytes, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(resultBytes, &result)
	if err != nil {
		return nil, err
	}
	return ConvertToCamelCase(result, schema).(map[string]interface{}), nil
}
