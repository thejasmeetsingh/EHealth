package utils

import "strings"

func ReplacePlaceholders(htmlContent []byte, placeholders map[string]string) []byte {
	for placeholder, dynamicContent := range placeholders {
		htmlContent = []byte(strings.Replace(string(htmlContent), placeholder, dynamicContent, -1))
	}

	return htmlContent
}
