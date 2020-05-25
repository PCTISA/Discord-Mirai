package util

import (
	"os"
	"strings"
)

/* === Helpers === */

// InitFile opens a file at the specified path. If that file does not exist,
// it creates a new one.
func InitFile(path string) (*os.File, error) {
	file, err := os.Open(path)
	if os.IsNotExist(err) {
		file, err := os.Create(path)

		if err != nil {
			return &os.File{}, err
		}
		return file, err
	}
	return file, err
}

// ArrayContains checks a string array for a given string.
func ArrayContains(array []string, value string, ignoreCase bool) bool {
	for _, e := range array {
		if ignoreCase {
			e = strings.ToLower(e)
		}

		if e == value {
			return true
		}
	}
	return false
}
