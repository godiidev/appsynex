// File: pkg/utils/string.go
// Tạo tại: pkg/utils/string.go
// Mục đích: Utility functions for string operations

package utils

import (
	"strings"
	"unicode"
)

// ToSnakeCase converts a string to snake_case
func ToSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if unicode.IsUpper(r) {
			if i > 0 {
				result = append(result, '_')
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

// ToCamelCase converts a string to camelCase
func ToCamelCase(str string) string {
	words := strings.Split(str, "_")
	if len(words) == 0 {
		return ""
	}
	
	result := strings.ToLower(words[0])
	for i := 1; i < len(words); i++ {
		if len(words[i]) > 0 {
			result += strings.ToUpper(string(words[i][0])) + strings.ToLower(words[i][1:])
		}
	}
	return result
}

// Capitalize capitalizes the first letter of a string
func Capitalize(str string) string {
	if len(str) == 0 {
		return str
	}
	return strings.ToUpper(string(str[0])) + strings.ToLower(str[1:])
}

// IsEmpty checks if a string is empty or contains only whitespace
func IsEmpty(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}

// Contains checks if a slice contains a string
func Contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}