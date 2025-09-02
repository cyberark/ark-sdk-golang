// Package common provides common data models and utilities for the ARK SDK.
// This package contains shared data types and helper functions that are used
// across different components of the ARK SDK for consistent data handling.
package common

import (
	"encoding/json"
	"time"
)

// ArkRFC3339Time is a custom time type that represents a time in RFC 3339 format.
// This type provides JSON marshaling and unmarshaling capabilities for time values
// that need to be serialized in RFC 3339 format with microsecond precision.
// It wraps the standard time.Time type and implements the json.Marshaler and
// json.Unmarshaler interfaces for proper JSON handling.
//
// Example usage:
//
//	var arkTime ArkRFC3339Time
//	err := json.Unmarshal([]byte(`"2023-01-01T12:00:00.123456Z"`), &arkTime)
//	if err != nil {
//		// handle error
//	}
type ArkRFC3339Time time.Time

// customTimeFormat defines the RFC 3339 time format with microsecond precision
// used for JSON marshaling and unmarshaling operations.
// Format: "2006-01-02T15:04:05.999999Z07:00"
const customTimeFormat = "2006-01-02T15:04:05.999999Z07:00"

// UnmarshalJSON implements the json.Unmarshaler interface for ArkRFC3339Time.
// It parses JSON data containing a time string in RFC 3339 format with microsecond
// precision and converts it to an ArkRFC3339Time value.
//
// The method handles both quoted and unquoted JSON strings, automatically removing
// surrounding quotes if present. It uses the customTimeFormat constant to parse
// the time string with the expected RFC 3339 format.
//
// Parameters:
//   - data: JSON byte data containing the time string to parse
//
// Returns:
//   - error: nil if parsing succeeds, otherwise an error describing the parse failure
//
// Example JSON input: "2023-01-01T12:00:00.123456Z"
func (ct *ArkRFC3339Time) UnmarshalJSON(data []byte) error {
	str := string(data)
	if str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}
	t, err := time.Parse(customTimeFormat, str)
	if err != nil {
		return err
	}
	*ct = ArkRFC3339Time(t)
	return nil
}

// MarshalJSON implements the json.Marshaler interface for ArkRFC3339Time.
// It converts an ArkRFC3339Time value to JSON format by formatting the underlying
// time value as an RFC 3339 string with microsecond precision.
//
// The method formats the time using the customTimeFormat constant and returns
// the result as a JSON-encoded string value.
//
// Returns:
//   - []byte: JSON-encoded byte array containing the formatted time string
//   - error: nil if marshaling succeeds, otherwise an error from json.Marshal
//
// Example output: "2023-01-01T12:00:00.123456Z"
func (ct *ArkRFC3339Time) MarshalJSON() ([]byte, error) {
	timeStr := time.Time(*ct).Format(customTimeFormat)
	return json.Marshal(timeStr)
}
