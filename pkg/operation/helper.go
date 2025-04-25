package operation

import (
	"regexp"
	"strings"
)

var matchPattern = regexp.MustCompile(`\[.*\]`)
var replacePattern = regexp.MustCompile(`\[[^\]]*\]`)

func replaceIndex(str string) string {
	if str == "" {
		return ""
	}

	strArr := strings.Split(str, ".")
	var result []string

	for _, s := range strArr {
		// Check if segment contains any array indexing
		if match := matchPattern.MatchString(s); match {
			// Replace all array indexing with [_]
			replaced := replacePattern.ReplaceAllString(s, "[_]")
			result = append(result, replaced)
			continue
		}
		result = append(result, s)
	}

	return strings.Join(result, ".")
}
