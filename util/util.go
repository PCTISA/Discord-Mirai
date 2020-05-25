package util

import (
	"os"
	"strings"

	"github.com/CS-5/disgomux"
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

// CheckPermissions takes the user, role(s), and channel IDs and checks them
// against the supplied permissions struct.
// TODO: This should probably be moved as a utility function to the multiplexer?
func CheckPermissions(
	perms *disgomux.CommandPermissions,
	userID string, roleIDs []string, chanID string,
) bool {
	if ArrayContains(perms.UserIDs, userID, true) {
		return true
	}

	for _, id := range roleIDs {
		if ArrayContains(perms.RoleIDs, id, true) {
			return true
		}
	}

	if ArrayContains(perms.ChanIDs, chanID, true) {
		return true
	}
	return false
}
