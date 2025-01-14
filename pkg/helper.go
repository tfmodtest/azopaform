package pkg

import (
	"regexp"
	"strings"
)

func replaceIndex(str string) string {
	strArr := strings.Split(str, ".")
	var result string
	for _, s := range strArr {
		if match, err := regexp.Match("[.*[0-9]+]|[*]", []byte(s)); err == nil && match {
			var newSegment string
			for _, c := range s {
				if c == '[' {
					break
				}
				newSegment += string(c)
			}
			newSegment += "[x]"
			result += newSegment + "."
		} else {
			result += s + "."
		}
	}
	return strings.TrimSuffix(result, ".")
}
