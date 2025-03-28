package operation

import (
	"regexp"
	"strings"
)

var pattern = regexp.MustCompile(`\[[^\]]*\]`)

func replaceIndex(str string) string {
	if str == "" {
		return ""
	}

	strArr := strings.Split(str, ".")
	var result []string

	for _, s := range strArr {
		// Check if segment contains any array indexing
		if match, err := regexp.MatchString(`\[.*\]`, s); err == nil && match {
			// Replace all array indexing with [_]
			replaced := pattern.ReplaceAllString(s, "[_]")
			result = append(result, replaced)
			continue
		}
		result = append(result, s)
	}

	return strings.Join(result, ".")
}
