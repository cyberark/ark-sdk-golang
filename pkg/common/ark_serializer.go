package common

import (
	"encoding/json"
	"github.com/iancoleman/strcase"
	"github.com/mitchellh/mapstructure"
	"io"
)

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

func convertToSnakeCase(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		snakeMap := make(map[string]interface{})
		for key, value := range v {
			snakeKey := strcase.ToSnake(key)
			snakeMap[snakeKey] = convertToSnakeCase(value)
		}
		return snakeMap
	case []interface{}:
		for i, item := range v {
			v[i] = convertToSnakeCase(item)
		}
		return v
	default:
		return v
	}
}

func convertToCamelCase(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		snakeMap := make(map[string]interface{})
		for key, value := range v {
			snakeKey := strcase.ToLowerCamel(key)
			snakeMap[snakeKey] = convertToCamelCase(value)
		}
		return snakeMap
	case []interface{}:
		for i, item := range v {
			v[i] = convertToCamelCase(item)
		}
		return v
	default:
		return v
	}
}

// DeserializeJSONSnake takes an io.ReadCloser response and deserializes it into a map with snake_case keys.
func DeserializeJSONSnake(response io.ReadCloser) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.NewDecoder(response).Decode(&result)
	if err != nil {
		return nil, err
	}
	return convertToSnakeCase(result).(map[string]interface{}), nil
}

// SerializeJSONCamel takes an interface and serializes it into a map with camelCase keys.
func SerializeJSONCamel(item interface{}) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := mapstructure.Decode(item, &result)
	if err != nil {
		return nil, err
	}
	return convertToCamelCase(result).(map[string]interface{}), nil
}
