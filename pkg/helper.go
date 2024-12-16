package pkg

import (
	"math/rand"
	"regexp"
	"strings"
)

const ifCondition = "if"
const regexExp = "regex.match"
const charNum = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const singleConditionLen = 3
const andConditionLen = 5
const orConditionLen = 7
const whereConditionLen = 9
const singleConditionLenPlusX = 6
const andConditionLenPlusX = 8
const orConditionLenPlusX = 10
const whereConditionLenPlusX = 12

const allOf = "allof"
const anyOf = "anyof"
const single = "single"
const count = "count"
const contains = "contains"
const notContains = "notcontains"
const containsKey = "containskey"
const equals = "equals"
const less = "less"
const greater = "greater"
const notMatch = "notmatch"
const in = "in"
const notIn = "notin"
const name = "name"
const exists = "exists"
const like = "like"
const notLike = "notlike"
const not = "not"
const notEquals = "notequals"
const greaterOrEquals = "greaterorequals"
const lessOrEquals = "lessorequals"
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

const concat = "concat"
const array = "array"
const empty = "empty"
const parameters = "parameters"

const resourcePrefix = "r."

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
