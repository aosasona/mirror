package helper

import "strings"

func WithDefaultString(value string, defaultValue string) string {
	if value == "" {
		return defaultValue
	}

	return value
}

func StringToBool(value string) bool {
	value = strings.ToLower(strings.TrimSpace(value))
	return value == "true" || value == "1"
}
