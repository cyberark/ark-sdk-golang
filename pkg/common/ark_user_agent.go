package common

import (
	"fmt"
	browser "github.com/EDDYCJY/fake-useragent"
)

// UserAgent returns the user agent string for the Ark SDK in Python.
func UserAgent() string {
	return browser.Chrome() + fmt.Sprintf(" Ark-SDK-Golang/%s", ArkVersion())
}
