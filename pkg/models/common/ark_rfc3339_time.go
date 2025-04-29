package common

import (
	"encoding/json"
	"time"
)

// ArkRFC3339Time is a custom time type that represents a time in RFC 3339 format.
type ArkRFC3339Time time.Time

const customTimeFormat = "2006-01-02T15:04:05.999999Z07:00"

// UnmarshalJSON parses the JSON data into the ArkRFC3339Time type.
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

// MarshalJSON converts the ArkRFC3339Time type to JSON format.
func (ct *ArkRFC3339Time) MarshalJSON() ([]byte, error) {
	timeStr := time.Time(*ct).Format(customTimeFormat)
	return json.Marshal(timeStr)
}
