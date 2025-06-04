package shared

import "github.com/xyproto/randomstring"

var m = make(map[string]struct{})

func HumanFriendlyEnglishString(length int) string {
	return generateUniqueRandomString(length, randomstring.HumanFriendlyEnglishString)
}

func HumanFriendlyString(length int) string {
	return generateUniqueRandomString(length, randomstring.HumanFriendlyString)
}

func generateUniqueRandomString(length int, generator func(int) string) string {
	for {
		v := generator(length)
		if _, ok := m[v]; !ok {
			m[v] = struct{}{}
			return v
		}
	}
}
