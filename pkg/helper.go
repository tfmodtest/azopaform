package pkg

import (
	"regexp"
	"strings"
)

func replaceIndex(str string) string {
	if str == "" {
		return ""
	}

	strArr := strings.Split(str, ".")
	var result []string

	for _, s := range strArr {
		// Check if segment contains any array indexing
		if match, err := regexp.Match(`\[.*\]`, []byte(s)); err == nil && match {
			// Replace all array indexing with [_]
			replaced := regexp.MustCompile(`\[[^\]]*\]`).ReplaceAllString(s, "[_]")
			result = append(result, replaced)
			continue
		}
		result = append(result, s)
	}

	return strings.Join(result, ".")
}
