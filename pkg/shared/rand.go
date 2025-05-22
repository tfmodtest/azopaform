package shared

import "github.com/xyproto/randomstring"

var m = make(map[string]struct{})

func HumanFriendlyEnglishString(length int) string {
	for {
		v := randomstring.HumanFriendlyEnglishString(length)
		if _, ok := m[v]; !ok {
			m[v] = struct{}{}
			return v
		}
	}
}

func HumanFriendlyString(length int) string {
	for {
		v := randomstring.HumanFriendlyString(length)
		if _, ok := m[v]; !ok {
			m[v] = struct{}{}
			return v
		}
	}
}
