package main

import "math/rand"

const ifCondition = "if"
const regexExp = "regex.match"
const charNum = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const singleConditionLen = 3
const andConditionLen = 5
const orConditionLen = 7
const whereConditionLen = 9

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
