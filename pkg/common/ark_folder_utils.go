package common

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// ExpandFolder expands the given folder path by replacing environment variables and user home directory.
func ExpandFolder(folder string) string {
	folderPath := os.ExpandEnv(folder)
	if !strings.HasSuffix(folderPath, "/") {
		folderPath += "/"
	}
	if len(folderPath) > 0 && folderPath[0] == '~' {
		usr, err := user.Current()
		if err != nil {
			return ""
		}
		return filepath.Join(usr.HomeDir, folderPath[1:])
	}
	return folderPath
}
