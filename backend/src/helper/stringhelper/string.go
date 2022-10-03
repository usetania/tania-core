package stringhelper

import (
	"strings"
)

// Join joins strings.
func Join(values ...string) string {
	return strings.Join(values, "")
}
