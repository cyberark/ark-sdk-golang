package common

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// ExpandFolder expands the given folder path by replacing environment variables and user home directory.
//
// This function performs path expansion by handling environment variables ($VAR || ${VAR}) and tilde (~)
// expansion for user home directories. It ensures the returned path ends with a trailing
// slash for consistency. If the path starts with ~/ or is exactly ~/, it replaces the
// tilde portion with the current user's home directory.
//
// The function processes the input in the following order:
// 1. Expands environment variables using os.ExpandEnv
// 2. Adds trailing slash if not present
// 3. Performs tilde expansion if path starts with ~/ or is exactly ~
//
// Note: The function has some quirks in its current implementation:
// - Returns empty string if user lookup fails during tilde expansion
//
// Parameters:
//   - folder: The folder path to expand (can contain environment variables and ~/)
//
// Returns the expanded folder path with trailing slash, or empty string if user lookup fails.
//
// Example:
//
//	expanded := ExpandFolder("~/Documents")        // Returns "/Users/username/Documents" (no trailing slash due to filepath.Join)
//	expanded := ExpandFolder("$HOME/config")       // Returns "/Users/username/config/"
//	expanded := ExpandFolder("/absolute/path")     // Returns "/absolute/path/"
//	expanded := ExpandFolder("~")                  // Returns "/Users/username/"
func ExpandFolder(folder string) string {
	folderPath := os.ExpandEnv(folder)
	if folderPath == "~" || strings.HasPrefix(folderPath, "~/") {
		usr, err := user.Current()
		if err != nil {
			return ""
		}
		folderPath = filepath.Join(usr.HomeDir, folderPath[1:])
	}
	if !strings.HasSuffix(folderPath, "/") {
		folderPath += "/"
	}
	return folderPath
}
