// Package common provides JSON serialization and deserialization utilities for the ARK SDK.
//
// This package handles conversion between different JSON key naming conventions (camelCase
// and snake_case) and provides utilities for serializing/deserializing JSON data with
// optional schema validation. It supports both simple conversions and schema-aware
// transformations that preserve specific field mappings based on struct tags.
package common

import (
	"encoding/json"
	"io"
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
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
//
// SerializeResponseToJSON reads all data from the provided io.ReadCloser, attempts to
// parse it as JSON, and returns a properly formatted JSON string. If the input data
// is not valid JSON, it returns the original data as a string. This function is useful
// for normalizing response data into a consistent JSON format.
//
// Parameters:
//   - response: The io.ReadCloser containing the response data to serialize
//
// Returns a JSON string representation of the response data, or the original data
// as a string if JSON parsing fails.
//
// Example:
//
//	jsonStr := SerializeResponseToJSON(httpResponse.Body)
//	fmt.Println(jsonStr) // Outputs properly formatted JSON
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
//
// ConvertToSnakeCase recursively processes data structures, converting all string keys
// from camelCase to snake_case format. It supports nested maps, slices, and arbitrary
// data types. When a schema is provided, it uses struct field information to determine
// whether specific fields should be converted or preserved as-is (e.g., for map fields
// with string keys that should not be converted).
//
// Parameters:
//   - data: The data structure to convert (supports maps, slices, and primitive types)
//   - schema: Optional reflect.Type pointer for schema-aware conversion (nil for simple conversion)
//
// Returns the converted data structure with snake_case keys.
//
// Example:
//
//	input := map[string]interface{}{"firstName": "John", "lastName": "Doe"}
//	result := ConvertToSnakeCase(input, nil)
//	// result: map[string]interface{}{"first_name": "John", "last_name": "Doe"}
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
//
// ConvertToCamelCase recursively processes data structures, converting all string keys
// from snake_case to camelCase format. It supports nested maps, slices, and arbitrary
// data types. When a schema is provided, it uses struct field information to determine
// whether specific fields should be converted or preserved as-is (e.g., for map fields
// with string keys that should not be converted).
//
// Parameters:
//   - data: The data structure to convert (supports maps, slices, and primitive types)
//   - schema: Optional reflect.Type pointer for schema-aware conversion (nil for simple conversion)
//
// Returns the converted data structure with camelCase keys.
//
// Example:
//
//	input := map[string]interface{}{"first_name": "John", "last_name": "Doe"}
//	result := ConvertToCamelCase(input, nil)
//	// result: map[string]interface{}{"firstName": "John", "lastName": "Doe"}
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
//
// DeserializeJSONSnake reads JSON data from the provided io.ReadCloser, parses it,
// and converts all keys from the original format to snake_case. This function is
// useful for normalizing JSON responses that may have keys in camelCase or other
// formats to a consistent snake_case format.
//
// Parameters:
//   - response: The io.ReadCloser containing JSON data to deserialize
//
// Returns the deserialized data with snake_case keys and any error encountered
// during JSON decoding.
//
// Example:
//
//	data, err := DeserializeJSONSnake(httpResponse.Body)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	// data contains the JSON with snake_case keys
func DeserializeJSONSnake(response io.ReadCloser) (interface{}, error) {
	var result interface{}
	err := json.NewDecoder(response).Decode(&result)
	if err != nil {
		return nil, err
	}
	return ConvertToSnakeCase(result, nil).(interface{}), nil
}

// SerializeJSONCamel takes an interface and serializes it into a map with camelCase keys.
//
// SerializeJSONCamel converts the provided data structure to JSON, then parses it
// back into a map with all keys converted to camelCase format. This function is
// useful for preparing data for APIs or systems that expect camelCase key naming
// conventions.
//
// Parameters:
//   - item: The data structure to serialize (must be JSON-serializable)
//
// Returns a map with camelCase keys and any error encountered during JSON
// marshaling or unmarshaling.
//
// Example:
//
//	input := struct{ FirstName string }{FirstName: "John"}
//	result, err := SerializeJSONCamel(input)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	// result: map[string]interface{}{"firstName": "John"}
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
//
// DeserializeJSONSnakeSchema reads JSON data from the provided io.ReadCloser, parses it,
// and converts all keys from the original format to snake_case using schema-aware
// conversion. The schema parameter allows for more intelligent key conversion by
// considering struct field types and tags, which helps preserve certain fields
// (like maps with string keys) that should not be converted.
//
// Parameters:
//   - response: The io.ReadCloser containing JSON data to deserialize
//   - schema: Pointer to reflect.Type for schema-aware conversion (can be nil)
//
// Returns the deserialized data with snake_case keys and any error encountered
// during JSON decoding.
//
// Example:
//
//	var schemaType reflect.Type = reflect.TypeOf(MyStruct{})
//	data, err := DeserializeJSONSnakeSchema(httpResponse.Body, &schemaType)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	// data contains the JSON with schema-aware snake_case keys
func DeserializeJSONSnakeSchema(response io.ReadCloser, schema *reflect.Type) (interface{}, error) {
	var result interface{}
	err := json.NewDecoder(response).Decode(&result)
	if err != nil {
		return nil, err
	}
	return ConvertToSnakeCase(result, schema).(interface{}), nil
}

// SerializeJSONCamelSchema takes an interface and serializes it into a map with camelCase keys.
//
// SerializeJSONCamelSchema converts the provided data structure to JSON, then parses it
// back into a map with all keys converted to camelCase format using schema-aware
// conversion. The schema parameter allows for more intelligent key conversion by
// considering struct field types and tags, which helps preserve certain fields
// (like maps with string keys) that should not be converted.
//
// Parameters:
//   - item: The data structure to serialize (must be JSON-serializable)
//   - schema: Pointer to reflect.Type for schema-aware conversion (can be nil)
//
// Returns a map with camelCase keys and any error encountered during JSON
// marshaling or unmarshaling.
//
// Example:
//
//	var schemaType reflect.Type = reflect.TypeOf(MyStruct{})
//	input := struct{ FirstName string }{FirstName: "John"}
//	result, err := SerializeJSONCamelSchema(input, &schemaType)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	// result: map[string]interface{}{"firstName": "John"}
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
