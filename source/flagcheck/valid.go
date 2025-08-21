package flagcheck

import (
	"regexp"
	"strings"
)

func IsValidBucketName(name string) bool {
	if len(name) < 3 || len(name) > 63 {
		return false
	}

	if !regexp.MustCompile(`^[a-z0-9.-]+$`).MatchString(name) {
		return false
	}

	if strings.HasPrefix(name, "-") || strings.HasSuffix(name, "-") {
		return false
	}

	if strings.Contains(name, "..") || strings.Contains(name, "--") {
		return false
	}

	ipAddressRegex := `^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`

	if regexp.MustCompile(ipAddressRegex).MatchString(name) {
		return false
	}

	return true
}

func ValidateObjectName(fileName string) bool {
	matched, _ := regexp.MatchString(`^[0-9A-Za-z\!\-\.\*\_\(\)]+$`, fileName) // "file"
	return matched
}
