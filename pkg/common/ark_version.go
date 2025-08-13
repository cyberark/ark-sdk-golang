package common

var arkVersion = "0.0.0"

// SetArkVersion sets the version of the Ark SDK.
func SetArkVersion(version string) {
	if version != "" {
		arkVersion = version
	}
}

// ArkVersion returns the current version of the Ark SDK.
func ArkVersion() string {
	return arkVersion
}
