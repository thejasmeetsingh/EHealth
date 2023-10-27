package utils

import "strings"

// Replace the placeholders in given HTML content with the dynamic data
func ReplacePlaceholders(htmlContent []byte, placeholders map[string]string) []byte {
	for placeholder, dynamicContent := range placeholders {
		htmlContent = []byte(strings.Replace(string(htmlContent), placeholder, dynamicContent, -1))
	}

	return htmlContent
}
