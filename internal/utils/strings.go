package utils

import "strings"

// ReplaceMultiple replaces multiple strings in a given string
func ReplaceMultiple(str string, replacements map[string]string) string {
	for oldValue, newValue := range replacements {
		str = strings.ReplaceAll(str, oldValue, newValue)
	}
	return str
}
