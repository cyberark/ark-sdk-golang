// Package common provides shared utilities and types for the ARK SDK.
//
// This package handles user agent string generation for HTTP requests made by the
// ARK SDK, combining browser user agent strings with SDK version information to
// provide proper identification in network communications.
package common

import (
	"fmt"

	browser "github.com/EDDYCJY/fake-useragent"
)

// UserAgent returns the user agent string for the Ark SDK in Golang.
//
// UserAgent generates a composite user agent string by combining a Chrome browser
// user agent (obtained from the fake-useragent library) with the current ARK SDK
// version. This provides proper identification for HTTP requests made by the SDK
// while maintaining compatibility with web services that expect browser-like
// user agents.
//
// Returns a formatted user agent string in the format:
// "{Chrome User Agent} Ark-SDK-Golang/{version}"
//
// Example:
//
//	userAgent := UserAgent()
//	// userAgent might be:
//	// "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Ark-SDK-Golang/1.2.3"
func UserAgent() string {
	return browser.Chrome() + fmt.Sprintf(" Ark-SDK-Golang/%s", ArkVersion())
}
