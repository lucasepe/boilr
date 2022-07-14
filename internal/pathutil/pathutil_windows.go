//go:build windows
// +build windows

package pathutil

import (
	"os"
	"path/filepath"
	"strings"
)

// Exists returns true if the specified path exists.
func Exists(path string) bool {
	fi, err := os.Lstat(path)
	if fi != nil && fi.Mode()&os.ModeSymlink != 0 {
		_, err = filepath.EvalSymlinks(path)
	}

	return err == nil || os.IsExist(err)
}

// ExpandHome substitutes `%USERPROFILE%` at the start of the specified
// `path` using the provided `home` location.
func ExpandHome(path, home string) string {
	if path == "" || home == "" {
		return path
	}
	if strings.HasPrefix(path, `%USERPROFILE%`) {
		return filepath.Join(home, path[13:])
	}

	return path
}
