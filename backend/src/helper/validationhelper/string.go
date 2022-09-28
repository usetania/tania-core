package validationhelper

import (
	"regexp"
)

// IsNumeric check if the string contains only numbers. Empty string is valid.
func IsNumeric(val string) bool {
	if val == "" {
		return true
	}

	regexNumeric := regexp.MustCompile("^[0-9]+$")

	return regexNumeric.MatchString(val)
}

// IsFloat check if the string is a float.
func IsFloat(val string) bool {
	regexFloat := regexp.MustCompile(`^(?:[-+]?(?:[0-9]+))?(?:\.[0-9]*)?(?:[eE][\+\-]?(?:[0-9]+))?$`)

	return regexFloat.MatchString(val)
}

// IsAlpha check if the string contains only letters. Empty string is valid.
func IsAlpha(val string) bool {
	regexAlpha := regexp.MustCompile("^[a-zA-Z]+$")

	return regexAlpha.MatchString(val)
}

// IsAlphanumeric check if the string contains only letters and numbers. Empty string is valid.
func IsAlphanumeric(val string) bool {
	if val == "" {
		return true
	}

	regexAlphanumeric := regexp.MustCompile("^[a-zA-Z0-9]+$")

	return regexAlphanumeric.MatchString(val)
}

// IsAlphanumeric check if the string contains only letters, numbers, space, hypens, and underscore.
// Only allow letters and numbers at the start and the end.
// Empty string is valid.
func IsAlphanumSpaceHyphenUnderscore(val string) bool {
	if val == "" {
		return true
	}

	regex := regexp.MustCompile("^[a-zA-Z0-9]+[a-zA-Z0-9-_ ]*[a-zA-Z0-9]$")

	return regex.MatchString(val)
}
