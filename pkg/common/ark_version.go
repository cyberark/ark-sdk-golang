// Package common provides shared utilities and types for the ARK SDK.
//
// This package handles version information storage and retrieval for the ARK SDK,
// providing a centralized way to manage and access the current SDK version throughout
// the application lifecycle.
package common

var arkVersion = "0.0.0"

// SetArkVersion sets the version of the Ark SDK.
//
// SetArkVersion updates the global arkVersion variable with the provided version
// string. If an empty string is provided, the version remains unchanged. This
// function is typically called during application initialization to set the
// correct SDK version.
//
// Parameters:
//   - version: The version string to set (empty string is ignored)
//
// Example:
//
//	SetArkVersion("1.2.3")
//	fmt.Println(ArkVersion()) // Outputs: 1.2.3
//
//	SetArkVersion("") // No change
//	fmt.Println(ArkVersion()) // Still outputs: 1.2.3
func SetArkVersion(version string) {
	if version != "" {
		arkVersion = version
	}
}

// ArkVersion returns the current version of the Ark SDK.
//
// ArkVersion retrieves the currently stored SDK version string. The default
// version is "0.0.0" if no version has been explicitly set using SetArkVersion.
//
// Returns the current SDK version string.
//
// Example:
//
//	version := ArkVersion()
//	fmt.Printf("Current SDK version: %s\n", version)
func ArkVersion() string {
	return arkVersion
}
