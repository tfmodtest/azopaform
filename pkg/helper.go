package pkg

import (
	"math/rand"
	"regexp"
	"strings"
)

const ifCondition = "if"
const regexExp = "regex.match"

const allOf = "allof"
const anyOf = "anyof"
const count = "count"
const name = "name"
const not = "not"
const field = "field"
const value = "value"
const where = "where"
const typeOfResource = "type"
const kindOfResource = "kind"

const deny = "deny"
const disabled = "disabled"
const allow = "allow"
const audit = "audit"
const warn = "warn"

var RandIntRange = func(min int, max int) int {
	return rand.Intn(max-min) + min
}

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
