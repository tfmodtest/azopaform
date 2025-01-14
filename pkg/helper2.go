package pkg

import (
	"strings"
)

func ReplaceIndex(str string) string {
	return strings.Replace(str, "[*]", "[_]", -1)
}
