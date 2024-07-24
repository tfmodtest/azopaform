package main

import "math/rand"

const ifCondition = "if"
const regexExp = "regex.match"
const charNum = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const singleConditionLen = 3
const andConditionLen = 5
const orConditionLen = 7
const whereConditionLen = 9

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

const deny = "deny"
const disabled = "disabled"
const allow = "allow"
const audit = "audit"
const warn = "warn"

func RandStringFromCharSet(strlen int, charSet string) string {
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = charSet[RandIntRange(0, len(charSet))]
	}
	return string(result)
}

func RandIntRange(min int, max int) int {
	return rand.Intn(max-min) + min
}
