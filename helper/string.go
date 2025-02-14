package helper

import "strings"

func WithDefaultString(value string, defaultValue string) string {
	if strings.TrimSpace(value) == "" {
		return defaultValue
	}

	return value
}

func StringToBool(value string) bool {
	value = strings.ToLower(strings.TrimSpace(value))
	return value == "true" || value == "1"
}
