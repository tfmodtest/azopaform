package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
)

// RuleSetReader takes a RuleSet and returns a string that can be used in a Rego file
func (ruleSet RuleSet) RuleSetReader(fieldNameReplacer string) ([]string, string, error) {
	var result string
	var conditionNames []string

	switch ruleSet.Flag {
	case single:
		singleRule := ruleSet.SingleRules[0]
		res, _, err := singleRule.SingleRuleReader()
		if err != nil {
			return []string{}, "", err
		}
		result = strings.Join([]string{result, "{"}, "")
		result = result + "\n"
		result = strings.Join([]string{result, res}, " ")
		if len(result) != 0 {
			result = result + "\n" + "}"
		}

		conditionName := RandStringFromCharSet(singleConditionLen, charNum)
		conditionNames = append(conditionNames, conditionName)
		result = strings.Join([]string{conditionName, ifCondition, result}, " ")
		result = result + "\n"

		return conditionNames, result, nil
	case not:
		singleRule := ruleSet.SingleRules[0]
		res, _, err := singleRule.SingleRuleReader()
		if err != nil {
			return []string{}, "", err
		}
		result = strings.Join([]string{result, "{"}, "")
		result = result + "\n"
		result = strings.Join([]string{result, not, res}, " ")
		if len(result) != 0 {
			result = result + "\n" + "}"
		}

		conditionName := RandStringFromCharSet(singleConditionLen, charNum)
		conditionNames = append(conditionNames, conditionName)
		result = strings.Join([]string{conditionName, ifCondition, result}, " ")
		result = result + "\n"

		return conditionNames, result, nil
	case allOf:
		var subsetResult string
		if len(ruleSet.SingleRules) != 0 {
			result = strings.Join([]string{result, "{"}, "")
			for _, singleRule := range ruleSet.SingleRules {
				//fmt.Printf("here is a singleRule with operator %+v\n", singleRule.Operator.Name)
				switch strings.ToLower(singleRule.Operator.Name) {
				case equals:
					fieldName, condition := FieldNameProcessor(singleRule.Field)
					if fieldNameReplacer != "" {
						if fieldName[len(fieldName)-3:] == "[*]" {
							fieldName = fieldName[:len(fieldName)-3]
						}
						fieldName = strings.Join([]string{fieldName, "[", fieldNameReplacer, "]"}, "")
					}

					result = result + "\n"
					result = strings.Join([]string{result, fieldName, "==", fmt.Sprint(singleRule.Operator.Value)}, " ")
					if condition != "" {
						if len(subsetResult) != 0 {
							subsetResult = strings.Join([]string{subsetResult, condition}, "")
						} else {
							subsetResult = condition
						}
					}
				case notEquals:
					fmt.Printf("here is a notEquals case %+v\n", singleRule.Operator.Value)
					fieldName, condition := FieldNameProcessor(singleRule.Field)
					if fieldNameReplacer != "" {
						if fieldName[len(fieldName)-3:] == "[*]" {
							fieldName = fieldName[:len(fieldName)-3]
						}
						fieldName = strings.Join([]string{fieldName, "[", fieldNameReplacer, "]"}, "")
					}

					result = result + "\n"
					result = strings.Join([]string{result, fieldName, "!=", fmt.Sprint(singleRule.Operator.Value)}, " ")
					if condition != "" {
						if len(subsetResult) != 0 {
							subsetResult = strings.Join([]string{subsetResult, condition}, "")
						} else {
							subsetResult = condition
						}
					}
				case greaterOrEquals:
					fieldName, condition := FieldNameProcessor(singleRule.Field)
					if fieldNameReplacer != "" {
						if fieldName[len(fieldName)-3:] == "[*]" {
							fieldName = fieldName[:len(fieldName)-3]
						}
						fieldName = strings.Join([]string{fieldName, "[", fieldNameReplacer, "]"}, "")
					}

					result = result + "\n"
					result = strings.Join([]string{result, fieldName, ">=", fmt.Sprint(singleRule.Operator.Value)}, " ")
					if condition != "" {
						if len(subsetResult) != 0 {
							subsetResult = strings.Join([]string{subsetResult, condition}, "")
						} else {
							subsetResult = condition
						}
					}
				case lessOrEquals:
					fieldName, condition := FieldNameProcessor(singleRule.Field)
					fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)

					result = result + "\n"
					result = strings.Join([]string{result, fieldName, "<=", fmt.Sprint(singleRule.Operator.Value)}, " ")
					if condition != "" {
						if len(subsetResult) != 0 {
							subsetResult = strings.Join([]string{subsetResult, condition}, "")
						} else {
							subsetResult = condition
						}
					}
				case less:
					fieldName, condition := FieldNameProcessor(singleRule.Field)
					if fieldNameReplacer != "" {
						if fieldName[len(fieldName)-3:] == "[*]" {
							fieldName = fieldName[:len(fieldName)-3]
						}
						fieldName = strings.Join([]string{fieldName, "[", fieldNameReplacer, "]"}, "")
					}

					result = result + "\n"
					result = strings.Join([]string{result, fieldName, "<", fmt.Sprint(singleRule.Operator.Value)}, " ")
					if condition != "" {
						if len(subsetResult) != 0 {
							subsetResult = strings.Join([]string{subsetResult, condition}, "")
						} else {
							subsetResult = condition
						}
					}
				case greater:
					fieldName, condition := FieldNameProcessor(singleRule.Field)
					fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)

					result = result + "\n"
					result = strings.Join([]string{result, fieldName, ">", fmt.Sprint(singleRule.Operator.Value)}, " ")
					if condition != "" {
						if len(subsetResult) != 0 {
							subsetResult = strings.Join([]string{subsetResult, condition}, "")
						} else {
							subsetResult = condition
						}
					}
				case exists:
					result = result + "\n"
					if reflect.String == reflect.TypeOf(singleRule.Operator.Value).Kind() {
						if strings.EqualFold(singleRule.Operator.Value.(string), "true") {
							fieldName := singleRule.Field.(string)
							fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)
							result = strings.Join([]string{result, fieldName}, " ")
						} else {
							fieldName := singleRule.Field.(string)
							fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)
							result = strings.Join([]string{result, not, fieldName}, " ")
						}
					} else if reflect.Bool == reflect.TypeOf(singleRule.Operator.Value).Kind() {
						if singleRule.Operator.Value.(bool) {
							fieldName := singleRule.Field.(string)
							fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)
							result = strings.Join([]string{result, fieldName}, " ")
						} else {
							fieldName := singleRule.Field.(string)
							fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)
							result = strings.Join([]string{result, not, fieldName}, " ")
						}
					}
				case contains:
					//fmt.Printf("Field name replacer is %s\n", fieldNameReplacer)
					result = result + "\n"
					fieldName := fmt.Sprint(singleRule.Field)
					if fieldNameReplacer != "" {
						if fieldName[len(fieldName)-3:] == "[*]" {
							fieldName = fieldName[:len(fieldName)-3]
						}
						fieldName = strings.Join([]string{fieldName, "[", fieldNameReplacer, "]"}, "")
					}

					result = strings.Join([]string{result, " ", regexExp, "(", "\"", ".*", fmt.Sprint(singleRule.Operator.Value), ".*", "\"", ",", fieldName, ")"}, "")
				case notContains:
					fieldName := fmt.Sprint(singleRule.Field)
					if fieldNameReplacer != "" {
						if fieldName[len(fieldName)-3:] == "[*]" {
							fieldName = fieldName[:len(fieldName)-3]
						}
						fieldName = strings.Join([]string{fieldName, "[", fieldNameReplacer, "]"}, "")
					}

					result = result + "\n"
					result = strings.Join([]string{result, " ", not, " ", regexExp, "(", "\"", ".*", fmt.Sprint(singleRule.Operator.Value), ".*", "\"", ",", fmt.Sprint(singleRule.Field), ")"}, "")
				case like:
					fieldName := singleRule.Field.(string)
					if fieldNameReplacer != "" {
						if fieldName[len(fieldName)-3:] == "[*]" {
							fieldName = fieldName[:len(fieldName)-3]
						}
						fieldName = strings.Join([]string{fieldName, "[", fieldNameReplacer, "]"}, "")
					}

					result = result + "\n"
					result = strings.Join([]string{result, " ", regexExp, "(", "\"", fmt.Sprint(singleRule.Operator.Value), "\"", ",", fieldName, ")"}, "")
				case notLike:
					fieldName := singleRule.Field.(string)
					if fieldNameReplacer != "" {
						if fieldName[len(fieldName)-3:] == "[*]" {
							fieldName = fieldName[:len(fieldName)-3]
						}
						fieldName = strings.Join([]string{fieldName, "[", fieldNameReplacer, "]"}, "")
					}

					result = result + "\n"
					result = strings.Join([]string{result, " ", not, " ", regexExp, "(", "\"", fmt.Sprint(singleRule.Operator.Value), "\"", ",", fieldName, ")"}, "")
				case where:
					//fmt.Printf("here is a where case %+v\n", singleRule)
					fieldName := singleRule.Field.(string)
					if fieldNameReplacer != "" {
						if fieldName[len(fieldName)-3:] == "[*]" {
							fieldName = fieldName[:len(fieldName)-3]
						}
						fieldName = strings.Join([]string{fieldName, "[", fieldNameReplacer, "]"}, "")
					}

					var exper string
					switch singleRule.FieldOperation {
					case count:
						operator := singleRule.Operator.Value.(RuleSet)
						subsetNames, subRule, err := operator.RuleSetReader(fieldName)
						if err != nil {
							return []string{}, "", err
						}
						//fmt.Printf("The field name is %s", fieldName)
						if string(fieldName[len(fieldName)-3:]) == "[*]" {
							//fmt.Printf("here is a fieldname %+v\n", fieldName)
							fieldName = fieldName[:len(fieldName)-3] + "[x]"
						}
						exper = count + "(" + "{" + "x" + "|" + fieldName + ";" + subsetNames[0] + "}" + ")"
						result = strings.Join([]string{result, " ", exper, fmt.Sprint(singleRule.Operator.Value)}, "")
						if len(subsetResult) != 0 {
							subsetResult = strings.Join([]string{subsetResult, subRule}, "")
						} else {
							subsetResult = subRule
						}
					}
				case in:
					fieldName, condition := FieldNameProcessor(singleRule.Field)
					// Find the common substring, replace it with the fieldNameReplacer with suffix [x]
					if fieldNameReplacer != "" {
						if fieldName[len(fieldName)-3:] == "[*]" {
							fieldName = fieldName[:len(fieldName)-3]
						}
						fieldName = strings.Join([]string{fieldName, "[", fieldNameReplacer, "]"}, "")
					}

					result = result + "\n"
					result = strings.Join([]string{result, "some", fieldName, "in", SliceConstructor(singleRule.Operator.Value)}, " ")
					if condition != "" {
						if len(subsetResult) != 0 {
							subsetResult = strings.Join([]string{subsetResult, condition}, "")
						} else {
							subsetResult = condition
						}
					}
				//TODO: notIn case is incorrectly addressed, should think of a way to express "not in" in rego
				case notIn:
					fieldName, condition := FieldNameProcessor(singleRule.Field)
					if fieldNameReplacer != "" {
						if strings.Contains(fieldName, fieldNameReplacer) {
							//fmt.Printf("it really contains...%s\n", fieldName)
							if fieldNameReplacer[len(fieldNameReplacer)-3:] == "[*]" {
								//fmt.Printf("it really needs to replace...%s\n", fieldName)
								newFieldNameReplacer := fieldNameReplacer[:len(fieldNameReplacer)-3] + "[x]"
								//fmt.Printf("let's replace this %s with %s\n", fieldNameReplacer, newFieldNameReplacer)
								fieldName = strings.Replace(fieldName, fieldNameReplacer, newFieldNameReplacer, 1)
								//fmt.Printf("after replacing %s\n", fieldName)
							}
						}
					}

					result = result + "\n"
					result = strings.Join([]string{result, "not", fieldName, "in", SliceConstructor(singleRule.Operator.Value)}, " ")
					if condition != "" {
						if len(subsetResult) != 0 {
							subsetResult = strings.Join([]string{subsetResult, condition}, "")
						} else {
							subsetResult = condition
						}
					}
				}
			}
		}

		for _, thisSet := range ruleSet.RuleSets {
			if len(result) == 0 {
				result = strings.Join([]string{result, "{"}, "")
			}
			subsetNames, subRule, err := thisSet.RuleSetReader("")
			if err != nil {
				return []string{}, "", err
			}
			if len(subsetResult) != 0 {
				subsetResult = strings.Join([]string{subsetResult, subRule}, "")
			} else {
				subsetResult = subRule
			}
			//conditionNames = append(conditionNames, subsetNames...)
			for _, subnetName := range subsetNames {
				result = result + "\n"
				if len(subnetName) == andConditionLen {
					result = strings.Join([]string{result, subnetName}, " ")
				} else if len(subnetName) == orConditionLen {
					result = strings.Join([]string{result, not, subnetName}, " ")
				} else if len(subnetName) == singleConditionLen {
					result = strings.Join([]string{result, subnetName}, " ")
				}
			}
		}
		if len(result) != 0 {
			result = result + "\n" + "}"
		}

		result = result + "\n" + subsetResult

		if len(ruleSet.SingleRules) != 0 || len(ruleSet.RuleSets) != 0 {
			conditionName := RandStringFromCharSet(andConditionLen, charNum)
			if fieldNameReplacer != "" {
				conditionName = conditionName + "(x)"
			}
			conditionNames = append(conditionNames, conditionName)
			result = strings.Join([]string{conditionName, ifCondition, result}, " ")
		}

		return conditionNames, result, nil
	case anyOf:
		var subsetResult string
		if len(ruleSet.SingleRules) != 0 {
			result = strings.Join([]string{result, "{"}, "")
			for _, singleRule := range ruleSet.SingleRules {
				//fmt.Printf("here is a singleRule with operator %+v\n", singleRule.Operator.Name)
				switch strings.ToLower(singleRule.Operator.Name) {
				case equals:
					fieldName, condition := FieldNameProcessor(singleRule.Field)
					if fieldNameReplacer != "" {
						if fieldName[len(fieldName)-3:] == "[*]" {
							fieldName = fieldName[:len(fieldName)-3]
						}
						fieldName = strings.Join([]string{fieldName, "[", fieldNameReplacer, "]"}, "")
					}
					result = result + "\n"
					result = strings.Join([]string{result, fieldName, "!=", fmt.Sprint(singleRule.Operator.Value)}, " ")
					if condition != "" {
						if len(subsetResult) != 0 {
							subsetResult = strings.Join([]string{subsetResult, condition}, "")
						} else {
							subsetResult = condition
						}
					}
				case notEquals:
					fieldName, condition := FieldNameProcessor(singleRule.Field)
					if fieldNameReplacer != "" {
						if fieldName[len(fieldName)-3:] == "[*]" {
							fieldName = fieldName[:len(fieldName)-3]
						}
						fieldName = strings.Join([]string{fieldName, "[", fieldNameReplacer, "]"}, "")
					}
					result = result + "\n"
					result = strings.Join([]string{result, fieldName, "==", fmt.Sprint(singleRule.Operator.Value)}, " ")
					if condition != "" {
						if len(subsetResult) != 0 {
							subsetResult = strings.Join([]string{subsetResult, condition}, "")
						} else {
							subsetResult = condition
						}
					}
				case greaterOrEquals:
					fieldName, condition := FieldNameProcessor(singleRule.Field)
					if fieldNameReplacer != "" {
						if fieldName[len(fieldName)-3:] == "[*]" {
							fieldName = fieldName[:len(fieldName)-3]
						}
						fieldName = strings.Join([]string{fieldName, "[", fieldNameReplacer, "]"}, "")
					}
					result = result + "\n"
					result = strings.Join([]string{result, fieldName, "<", fmt.Sprint(singleRule.Operator.Value)}, " ")
					if condition != "" {
						if len(subsetResult) != 0 {
							subsetResult = strings.Join([]string{subsetResult, condition}, "")
						} else {
							subsetResult = condition
						}
					}
				case less:
					fieldName, condition := FieldNameProcessor(singleRule.Field)
					if fieldNameReplacer != "" {
						if fieldName[len(fieldName)-3:] == "[*]" {
							fieldName = fieldName[:len(fieldName)-3]
						}
						fieldName = strings.Join([]string{fieldName, "[", fieldNameReplacer, "]"}, "")
					}
					result = result + "\n"
					result = strings.Join([]string{result, fieldName, ">=", fmt.Sprint(singleRule.Operator.Value)}, " ")
					if condition != "" {
						if len(subsetResult) != 0 {
							subsetResult = strings.Join([]string{subsetResult, condition}, "")
						} else {
							subsetResult = condition
						}
					}
				case greater:
					fieldName, condition := FieldNameProcessor(singleRule.Field)
					if fieldNameReplacer != "" {
						if fieldName[len(fieldName)-3:] == "[*]" {
							fieldName = fieldName[:len(fieldName)-3]
						}
						fieldName = strings.Join([]string{fieldName, "[", fieldNameReplacer, "]"}, "")
					}
					result = result + "\n"
					result = strings.Join([]string{result, fieldName, "<=", fmt.Sprint(singleRule.Operator.Value)}, " ")
					if condition != "" {
						if len(subsetResult) != 0 {
							subsetResult = strings.Join([]string{subsetResult, condition}, "")
						} else {
							subsetResult = condition
						}
					}
				case exists:
					result = result + "\n"
					if reflect.String == reflect.TypeOf(singleRule.Operator.Value).Kind() {
						if strings.EqualFold(singleRule.Operator.Value.(string), "true") {
							fieldName := singleRule.Field.(string)
							fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)
							result = strings.Join([]string{result, not, fieldName}, " ")
						} else {
							fieldName := singleRule.Field.(string)
							fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)
							result = strings.Join([]string{result, fieldName}, " ")
						}
					} else if reflect.Bool == reflect.TypeOf(singleRule.Operator.Value).Kind() {
						if singleRule.Operator.Value.(bool) {
							fieldName := singleRule.Field.(string)
							fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)
							result = strings.Join([]string{result, not, fieldName}, " ")
						} else {
							fieldName := singleRule.Field.(string)
							fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)
							result = strings.Join([]string{result, fieldName}, " ")
						}
					}
				case contains:
					result = result + "\n"
					fieldName := fmt.Sprint(singleRule.Field)
					fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)
					result = strings.Join([]string{result, " ", not, " ", regexExp, "(", "\"", ".*", fmt.Sprint(singleRule.Operator.Value), ".*", "\"", ",", fieldName, ")"}, "")
				case notContains:
					result = result + "\n"
					fieldName := fmt.Sprint(singleRule.Field)
					fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)
					result = strings.Join([]string{result, " ", regexExp, "(", "\"", ".*", fmt.Sprint(singleRule.Operator.Value), ".*", "\"", ",", fieldName, ")"}, "")
				case like:
					fieldName := singleRule.Field.(string)
					fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)
					result = result + "\n"
					result = strings.Join([]string{result, " ", not, " ", regexExp, "(", "\"", fmt.Sprint(singleRule.Operator.Value), "\"", ",", fieldName, ")"}, "")
				case notLike:
					fieldName := singleRule.Field.(string)
					fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)
					result = result + "\n"
					result = strings.Join([]string{result, regexExp, "(", "\"", fmt.Sprint(singleRule.Operator.Value), "\"", ",", fieldName, ")"}, "")
				case where:
					fieldName := singleRule.Field.(string)
					fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)
					var exper string
					switch singleRule.FieldOperation {
					case count:
						operator := singleRule.Operator.Value.(RuleSet)
						subsetNames, subRule, err := operator.RuleSetReader("x")
						if err != nil {
							return []string{}, "", err
						}
						if string(fieldName[len(fieldName)-3:]) == "[*]" {
							fmt.Printf("here is a fieldname %+v\n", fieldName)
							fieldName = fieldName[:len(fieldName)-3] + "[x]"
						}
						exper = count + "(" + "{" + "x" + "|" + fieldName + ";" + subsetNames[0] + "}" + ")"
						result = strings.Join([]string{result, " ", exper, fmt.Sprint(singleRule.Operator.Value)}, "")
						if len(subsetResult) != 0 {
							subsetResult = strings.Join([]string{subsetResult, subRule}, "")
						} else {
							subsetResult = subRule
						}
					}
				//TODO: notIn case is incorrectly addressed, should think of a way to express "not in" in rego
				case in:
					fieldName, condition := FieldNameProcessor(singleRule.Field)
					fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)
					result = result + "\n"
					result = strings.Join([]string{result, "not", fieldName, "in", SliceConstructor(singleRule.Operator.Value)}, " ")
					if condition != "" {
						if len(subsetResult) != 0 {
							subsetResult = strings.Join([]string{subsetResult, condition}, "")
						} else {
							subsetResult = condition
						}
					}
				case notIn:
					fieldName, condition := FieldNameProcessor(singleRule.Field)
					fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)
					result = result + "\n"
					result = strings.Join([]string{result, "some", fieldName, "in", SliceConstructor(singleRule.Operator.Value)}, " ")
					if condition != "" {
						if len(subsetResult) != 0 {
							subsetResult = strings.Join([]string{subsetResult, condition}, "")
						} else {
							subsetResult = condition
						}
					}
				}

			}
		}

		for _, thisSet := range ruleSet.RuleSets {
			if len(result) == 0 {
				result = strings.Join([]string{result, "{"}, "")
			}
			subsetNames, subRules, err := thisSet.RuleSetReader("")
			if err != nil {
				return []string{}, "", err
			}
			if len(subsetResult) != 0 {
				subsetResult = strings.Join([]string{subsetResult, subRules}, "")
			} else {
				subsetResult = subRules
			}
			//The subrule set also need to be the opposite
			for _, subnetName := range subsetNames {
				result = result + "\n"
				if len(subnetName) == andConditionLen {
					result = strings.Join([]string{result, not, subnetName}, " ")
				} else if len(subnetName) == orConditionLen {
					result = strings.Join([]string{result, subnetName}, " ")
				} else if len(subnetName) == singleConditionLen {
					result = strings.Join([]string{result, subnetName}, " ")
				}
			}
		}
		if len(result) != 0 {
			result = result + "\n" + "}"
		}

		result = result + "\n" + subsetResult

		if len(ruleSet.SingleRules) != 0 || len(ruleSet.RuleSets) != 0 {
			conditionName := RandStringFromCharSet(orConditionLen, charNum)
			if fieldNameReplacer != "" {
				conditionName = conditionName + "(x)"
			}
			conditionNames = append(conditionNames, conditionName)
			result = strings.Join([]string{conditionName, ifCondition, result}, " ")
		}
		//fmt.Printf("here is a subresult2 %s", result)
		//fmt.Printf("conditions are %+v\n", conditionNames)
		//fmt.Printf("here is a subresult2 %s", result)
		return conditionNames, result, nil
	case where:
		var subsetResult string
		if len(ruleSet.SingleRules) != 0 {
			result = strings.Join([]string{result, "{"}, "")
			for _, singleRule := range ruleSet.SingleRules {
				switch strings.ToLower(singleRule.Operator.Name) {
				case equals:
					fieldName, condition := FieldNameProcessor(singleRule.Field)
					fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)

					result = result + "\n"
					result = strings.Join([]string{result, fieldName, "==", fmt.Sprint(singleRule.Operator.Value)}, " ")
					if condition != "" {
						if len(subsetResult) != 0 {
							subsetResult = strings.Join([]string{subsetResult, condition}, "")
						} else {
							subsetResult = condition
						}
					}
				case notEquals:
					fieldName, condition := FieldNameProcessor(singleRule.Field)
					fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)
					result = result + "\n"
					result = strings.Join([]string{result, fieldName, "!=", fmt.Sprint(singleRule.Operator.Value)}, " ")
					if condition != "" {
						if len(subsetResult) != 0 {
							subsetResult = strings.Join([]string{subsetResult, condition}, "")
						} else {
							subsetResult = condition
						}
					}
				case exists:
					result = result + "\n"
					if strings.EqualFold(singleRule.Operator.Value.(string), "true") {
						fieldName := singleRule.Field.(string)
						fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)
						result = strings.Join([]string{result, fieldName}, " ")
					} else {
						fieldName := singleRule.Field.(string)
						fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)
						result = strings.Join([]string{result, fieldName, not}, " ")
					}
				case contains:
					fieldName := fmt.Sprint(singleRule.Field)
					fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)
					result = result + "\n"
					result = strings.Join([]string{result, " ", regexExp, "(", "\"", ".*", fmt.Sprint(singleRule.Operator.Value), ".*", "\"", ",", fieldName, ")"}, "")
				case notContains:
					fieldName := fmt.Sprint(singleRule.Field)
					fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)
					result = result + "\n"
					result = strings.Join([]string{result, " ", not, " ", regexExp, "(", "\"", ".*", fmt.Sprint(singleRule.Operator.Value), ".*", "\"", ",", fieldName, ")"}, "")
				case like:
					fieldName := singleRule.Field.(string)
					fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)
					result = result + "\n"
					result = strings.Join([]string{result, " ", regexExp, "(", fmt.Sprint(singleRule.Operator.Value), ",", fieldName, ")"}, "")
				case notLike:
					fieldName := singleRule.Field.(string)
					fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)
					result = result + "\n"
					result = strings.Join([]string{result, " ", not, " ", regexExp, "(", fmt.Sprint(singleRule.Operator.Value), ",", fieldName, ")"}, "")
				case in:
					fieldName, condition := FieldNameProcessor(singleRule.Field)
					fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)
					result = result + "\n"
					result = strings.Join([]string{result, "some", fieldName, "in", fmt.Sprint(singleRule.Operator.Value)}, " ")
					if condition != "" {
						if len(subsetResult) != 0 {
							subsetResult = strings.Join([]string{subsetResult, condition}, "")
						} else {
							subsetResult = condition
						}
					}
				}
			}
		}

		for _, thisSet := range ruleSet.RuleSets {
			if len(result) == 0 {
				result = strings.Join([]string{result, "{"}, "")
			}
			subsetNames, subRules, err := thisSet.RuleSetReader(fieldNameReplacer)
			if err != nil {
				return []string{}, "", err
			}
			if len(subsetResult) != 0 {
				subsetResult = strings.Join([]string{subsetResult, subRules}, "")
			} else {
				subsetResult = subRules
			}
			//conditionNames = append(conditionNames, subsetNames...)
			for _, subnetName := range subsetNames {
				result = result + "\n"
				if len(subnetName) == andConditionLen {
					result = strings.Join([]string{result, subnetName}, " ")
				} else if len(subnetName) == orConditionLen {
					result = strings.Join([]string{result, not, subnetName}, " ")
				} else if len(subnetName) == singleConditionLen {
					result = strings.Join([]string{result, subnetName}, " ")
				}
			}
		}

		if len(result) != 0 {
			result = result + "\n" + "}"
		}

		result = result + "\n" + subsetResult

		if len(ruleSet.SingleRules) != 0 || len(ruleSet.RuleSets) != 0 {
			conditionName := RandStringFromCharSet(whereConditionLen, charNum)
			if fieldNameReplacer != "" {
				conditionName = conditionName + "(x)"
			}
			conditionNames = append(conditionNames, conditionName)
			result = strings.Join([]string{conditionName, ifCondition, result}, " ")
		}
		return conditionNames, result, nil
		//case not:
		//	var subsetResult string
		//	if len(ruleSet.SingleRules) != 0 {
		//		result = strings.Join([]string{result, "{"}, "")
		//		for _, singleRule := range ruleSet.SingleRules {
		//			//fmt.Printf("here is a singleRule with operator %+v\n", singleRule.Operator.Name)
		//			switch strings.ToLower(singleRule.Operator.Name) {
		//			case equals:
		//				fieldName, condition := FieldNameProcessor(singleRule.Field)
		//				fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)
		//
		//				result = result + "\n"
		//				result = strings.Join([]string{result, "not", fieldName, "==", fmt.Sprint(singleRule.Operator.Value)}, " ")
		//				if condition != "" {
		//					if len(subsetResult) != 0 {
		//						subsetResult = strings.Join([]string{subsetResult, condition}, "")
		//					} else {
		//						subsetResult = condition
		//					}
		//				}
		//			case notEquals:
		//				fmt.Printf("here is a notEquals case %+v\n", singleRule.Operator.Value)
		//				fieldName, condition := FieldNameProcessor(singleRule.Field)
		//				fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)
		//
		//				result = result + "\n"
		//				result = strings.Join([]string{result, "not", fieldName, "!=", fmt.Sprint(singleRule.Operator.Value)}, " ")
		//				if condition != "" {
		//					if len(subsetResult) != 0 {
		//						subsetResult = strings.Join([]string{subsetResult, condition}, "")
		//					} else {
		//						subsetResult = condition
		//					}
		//				}
		//			case greaterOrEquals:
		//				fieldName, condition := FieldNameProcessor(singleRule.Field)
		//				fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)
		//
		//				result = result + "\n"
		//				result = strings.Join([]string{result, "not", fieldName, ">=", fmt.Sprint(singleRule.Operator.Value)}, " ")
		//				if condition != "" {
		//					if len(subsetResult) != 0 {
		//						subsetResult = strings.Join([]string{subsetResult, condition}, "")
		//					} else {
		//						subsetResult = condition
		//					}
		//				}
		//			case lessOrEquals:
		//				fieldName, condition := FieldNameProcessor(singleRule.Field)
		//				fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)
		//
		//				result = result + "\n"
		//				result = strings.Join([]string{result, "not", fieldName, "<=", fmt.Sprint(singleRule.Operator.Value)}, " ")
		//				if condition != "" {
		//					if len(subsetResult) != 0 {
		//						subsetResult = strings.Join([]string{subsetResult, condition}, "")
		//					} else {
		//						subsetResult = condition
		//					}
		//				}
		//			case less:
		//				fieldName, condition := FieldNameProcessor(singleRule.Field)
		//				fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)
		//
		//				result = result + "\n"
		//				result = strings.Join([]string{result, "not", fieldName, "<", fmt.Sprint(singleRule.Operator.Value)}, " ")
		//				if condition != "" {
		//					if len(subsetResult) != 0 {
		//						subsetResult = strings.Join([]string{subsetResult, condition}, "")
		//					} else {
		//						subsetResult = condition
		//					}
		//				}
		//			case exists:
		//				result = result + "\n"
		//				if reflect.String == reflect.TypeOf(singleRule.Operator.Value).Kind() {
		//					if strings.EqualFold(singleRule.Operator.Value.(string), "true") {
		//						fieldName := singleRule.Field.(string)
		//						fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)
		//						result = strings.Join([]string{result, not, fieldName}, " ")
		//					} else {
		//						fieldName := singleRule.Field.(string)
		//						fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)
		//						result = strings.Join([]string{result, fieldName}, " ")
		//					}
		//				} else if reflect.Bool == reflect.TypeOf(singleRule.Operator.Value).Kind() {
		//					if singleRule.Operator.Value.(bool) {
		//						fieldName := singleRule.Field.(string)
		//						fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)
		//						result = strings.Join([]string{result, not, fieldName}, " ")
		//					} else {
		//						fieldName := singleRule.Field.(string)
		//						fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)
		//						result = strings.Join([]string{result, fieldName}, " ")
		//					}
		//				}
		//			case contains:
		//				//fmt.Printf("Field name replacer is %s\n", fieldNameReplacer)
		//				result = result + "\n"
		//				fieldName := fmt.Sprint(singleRule.Field)
		//				fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)
		//				result = strings.Join([]string{result, " ", not, " ", regexExp, "(", "\"", ".*", fmt.Sprint(singleRule.Operator.Value), ".*", "\"", ",", fieldName, ")"}, "")
		//			case notContains:
		//				fieldName := fmt.Sprint(singleRule.Field)
		//				fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)
		//				result = result + "\n"
		//				result = strings.Join([]string{result, " ", regexExp, "(", "\"", ".*", fmt.Sprint(singleRule.Operator.Value), ".*", "\"", ",", fmt.Sprint(singleRule.Field), ")"}, "")
		//			case like:
		//				fieldName := singleRule.Field.(string)
		//				fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)
		//
		//				result = result + "\n"
		//				result = strings.Join([]string{result, " ", not, " ", regexExp, "(", "\"", fmt.Sprint(singleRule.Operator.Value), "\"", ",", fieldName, ")"}, "")
		//			case notLike:
		//				fieldName := singleRule.Field.(string)
		//				fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)
		//
		//				result = result + "\n"
		//				result = strings.Join([]string{result, " ", regexExp, "(", "\"", fmt.Sprint(singleRule.Operator.Value), "\"", ",", fieldName, ")"}, "")
		//			case where:
		//				//fmt.Printf("here is a where case %+v\n", singleRule)
		//				fieldName := singleRule.Field.(string)
		//				if fieldNameReplacer != "" {
		//					if strings.Contains(fieldName, fieldNameReplacer) {
		//						//fmt.Printf("it really contains...%s\n", fieldName)
		//						if fieldNameReplacer[len(fieldNameReplacer)-3:] == "[*]" {
		//							//fmt.Printf("it really needs to replace...%s\n", fieldName)
		//							newFieldNameReplacer := fieldNameReplacer[:len(fieldNameReplacer)-3] + "[x]"
		//							//fmt.Printf("let's replace this %s with %s\n", fieldNameReplacer, newFieldNameReplacer)
		//							fieldName = strings.Replace(fieldName, fieldNameReplacer, newFieldNameReplacer, 1)
		//							//fmt.Printf("after replacing %s\n", fieldName)
		//						}
		//					}
		//				}
		//
		//				var exper string
		//				switch singleRule.FieldOperation {
		//				case count:
		//					operator := singleRule.Operator.Value.(RuleSet)
		//					subsetNames, subRule, err := operator.RuleSetReader(fieldName)
		//					if err != nil {
		//						return []string{}, "", err
		//					}
		//					//fmt.Printf("The field name is %s", fieldName)
		//					if string(fieldName[len(fieldName)-3:]) == "[*]" {
		//						//fmt.Printf("here is a fieldname %+v\n", fieldName)
		//						fieldName = fieldName[:len(fieldName)-3] + "[x]"
		//					}
		//					exper = count + "(" + "{" + "x" + "|" + fieldName + ";" + subsetNames[0] + "}" + ")"
		//					result = strings.Join([]string{result, " ", exper, fmt.Sprint(singleRule.Operator.Value)}, "")
		//					if len(subsetResult) != 0 {
		//						subsetResult = strings.Join([]string{subsetResult, subRule}, "")
		//					} else {
		//						subsetResult = subRule
		//					}
		//				}
		//			case in:
		//				fieldName, condition := FieldNameProcessor(singleRule.Field)
		//				// Find the common substring, replace it with the fieldNameReplacer with suffix [x]
		//				fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)
		//
		//				result = result + "\n"
		//				result = strings.Join([]string{result, "not", fieldName, "in", fmt.Sprint(singleRule.Operator.Value)}, " ")
		//				if condition != "" {
		//					if len(subsetResult) != 0 {
		//						subsetResult = strings.Join([]string{subsetResult, condition}, "")
		//					} else {
		//						subsetResult = condition
		//					}
		//				}
		//			//TODO: notIn case is incorrectly addressed, should think of a way to express "not in" in rego
		//			case notIn:
		//				fieldName, condition := FieldNameProcessor(singleRule.Field)
		//				fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)
		//
		//				result = result + "\n"
		//				result = strings.Join([]string{result, "some", fieldName, "in", SliceConstructor(singleRule.Operator.Value)}, " ")
		//				if condition != "" {
		//					if len(subsetResult) != 0 {
		//						subsetResult = strings.Join([]string{subsetResult, condition}, "")
		//					} else {
		//						subsetResult = condition
		//					}
		//				}
		//			}
		//		}
		//	}
		//
		//	for _, thisSet := range ruleSet.RuleSets {
		//		if len(result) == 0 {
		//			result = strings.Join([]string{result, "{"}, "")
		//		}
		//		subsetNames, subRule, err := thisSet.RuleSetReader("")
		//		if err != nil {
		//			return []string{}, "", err
		//		}
		//		if len(subsetResult) != 0 {
		//			subsetResult = strings.Join([]string{subsetResult, subRule}, "")
		//		} else {
		//			subsetResult = subRule
		//		}
		//		//conditionNames = append(conditionNames, subsetNames...)
		//		for _, subnetName := range subsetNames {
		//			result = result + "\n"
		//			if len(subnetName) == andConditionLen {
		//				result = strings.Join([]string{result, subnetName}, " ")
		//			} else if len(subnetName) == orConditionLen {
		//				result = strings.Join([]string{result, not, subnetName}, " ")
		//			}
		//		}
		//	}
		//	if len(result) != 0 {
		//		result = result + "\n" + "}"
		//	}
		//
		//	result = result + "\n" + subsetResult
		//
		//	if len(ruleSet.SingleRules) != 0 || len(ruleSet.RuleSets) != 0 {
		//		conditionName := RandStringFromCharSet(andConditionLen, charNum)
		//		if fieldNameReplacer != "" {
		//			conditionName = conditionName + "(x)"
		//		}
		//		conditionNames = append(conditionNames, conditionName)
		//		result = strings.Join([]string{conditionName, ifCondition, result}, " ")
		//	}
		//
		//	return conditionNames, result, nil
	}
	return conditionNames, result, nil
}

func (singleRule SingleRule) SingleRuleReader() (string, string, error) {
	var result string
	var rules string

	//The same as "where" case below without operator name/value, without adding conditions suffix/prefix
	if singleRule.Operator.Name == "" {
		fieldName := singleRule.Field.(string)
		if string(fieldName[len(fieldName)-3:]) == "[*]" {
			fieldName = fieldName[:len(fieldName)-3]
		}
		exp := count + "(" + fieldName + ")"
		result = strings.Join([]string{result, exp}, "")
	}
	switch strings.ToLower(singleRule.Operator.Name) {
	case equals:
		//fmt.Printf("print this single rule field ! %+v\n", singleRule.Field)
		fieldName, _ := FieldNameProcessor(singleRule.Field)
		//fmt.Printf("the conditions from the field processor is %s\n", condition)

		//fieldName := singleRule.Field.(string)
		result = strings.Join([]string{fieldName, "==", fmt.Sprint(singleRule.Operator.Value)}, " ")
	case notEquals:
		fieldName, _ := FieldNameProcessor(singleRule.Field)
		result = strings.Join([]string{fieldName, "!=", fmt.Sprint(singleRule.Operator.Value)}, " ")
	case exists:
		if strings.EqualFold(singleRule.Operator.Value.(string), "true") {
			fieldName := singleRule.Field.(string)
			result = fieldName
		} else {
			fieldName := singleRule.Field.(string)
			result = strings.Join([]string{not, fieldName}, " ")
		}
	case like:
		fieldName, _ := FieldNameProcessor(singleRule.Field)
		result = strings.Join([]string{regexExp, "(", fmt.Sprint(singleRule.Operator.Value), ",", fieldName, ")"}, "")
	case notLike:
		fieldName, _ := FieldNameProcessor(singleRule.Field)
		result = strings.Join([]string{not, " ", regexExp, "(", fmt.Sprint(singleRule.Operator.Value), ",", fieldName, ")"}, "")
	case where:
		//fmt.Printf("here is a where case %+v\n", singleRule)
		var subNames []string
		fieldName, _ := FieldNameProcessor(singleRule.Field)
		switch singleRule.Operator.Value.(type) {
		case SingleRule:
			//fmt.Printf("here is a singlerule case %+v\n", singleRule.Operator.Value)
			operator := singleRule.Operator.Value.(SingleRule)
			operatorSet := RuleSet{
				Flag:        "allOf",
				SingleRules: []SingleRule{operator},
				RuleSets:    nil,
			}
			subsetNames, subRule, err := operatorSet.RuleSetReader("x")
			if err != nil {
				return "", "", err
			}
			subNames = subsetNames
			rules = subRule
		case RuleSet:
			//fmt.Printf("rule set hit. the field name is %s\n", fieldName)
			//fmt.Printf("rule set hit. The whole rule set is %+v\n", singleRule.Operator.Value)
			operator := singleRule.Operator.Value.(RuleSet)
			subsetNames, subRule, err := operator.RuleSetReader(fieldName)
			if err != nil {
				return "", "", err
			}
			subNames = subsetNames
			rules = subRule
		}

		//fmt.Printf("The rules are %+v\n", rules)
		if string(fieldName[len(fieldName)-3:]) == "[*]" {
			fieldName = fieldName[:len(fieldName)-3] + "[x]"
		}
		exp := count + "(" + "{" + "x" + "|" + fieldName + ";" + subNames[0] + "}" + ")"
		result = exp
	case in:
		fieldName, _ := FieldNameProcessor(singleRule.Field)
		result = strings.Join([]string{"some", fieldName, "in", fmt.Sprint(singleRule.Operator.Value)}, " ")
	}
	return result, rules, nil
}

func FieldNameProcessor(fieldName interface{}) (string, string) {
	var result string
	var rules string
	switch fieldName.(type) {
	case string:
		result = fieldName.(string)
	case SingleRule:
		fmt.Printf("the field name is %+v\n", fieldName)
		res, singleRule, err := fieldName.(SingleRule).SingleRuleReader()
		if err != nil {
			return "", ""
		}
		result = res
		rules = singleRule
	}

	return result, rules
}

func RegoWriter(fileName string, condition string) error {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(condition); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func SliceConstructor(input any) string {
	var array []string
	var res string
	fmt.Printf("the input type is %+v\n", reflect.TypeOf(input))
	switch input.(type) {
	case []interface{}:
		for _, v := range input.([]interface{}) {
			array = append(array, fmt.Sprint(v))
		}
	case []string:
		for _, v := range input.([]string) {
			array = append(array, fmt.Sprint(v))
		}
	case string:
		array = append(array, fmt.Sprint(input))
	}

	res = strings.Join(array, ",")
	res = strings.Join([]string{"[", res, "]"}, "")
	return res
}

func FieldNameReplacer(fieldName string, replacer string) string {
	if replacer != "" {
		if strings.Contains(fieldName, replacer) {
			//fmt.Printf("it really contains...%s\n", fieldName)
			if replacer[len(replacer)-3:] == "[*]" {
				//fmt.Printf("it really needs to replace...%s\n", fieldName)
				newFieldNameReplacer := replacer[:len(replacer)-3] + "[x]"
				//fmt.Printf("let's replace this %s with %s\n", fieldNameReplacer, newFieldNameReplacer)
				fieldName = strings.Replace(fieldName, replacer, newFieldNameReplacer, 1)
				//fmt.Printf("after replacing %s\n", fieldName)
			}
		}
	}

	return fieldName
}
