package pkg

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/magodo/aztfq/aztfq"
)

// RuleSetReader takes a RuleSet and returns a string that can be used in a Rego file
func (ruleSet RuleSet) RuleSetReader(fieldNameReplacer string, ctx context.Context) ([]string, string, error) {
	var result string
	var conditionNames []string

	switch ruleSet.Flag {
	case single:
		singleRule := ruleSet.SingleRules[0]
		res, _, err := singleRule.SingleRuleReader(ctx)
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
		res, _, err := singleRule.SingleRuleReader(ctx)
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
				switch strings.ToLower(singleRule.Operator.Name) {
				case equals:
					fieldName, condition, _ := FieldNameProcessor(singleRule.Field, ctx)
					if fieldNameReplacer != "" && fieldName[len(fieldName)-3:] == "[*]" {
						fieldName = strings.Join([]string{fieldName[:len(fieldName)-3], "[", fieldNameReplacer, "]"}, "")
					} else {
						fieldName = strings.Replace(fieldName, "*", "x", -1)
					}

					result = result + "\n"
					operatorValue := fmt.Sprint(singleRule.Operator.Value)
					if operatorValue == "" {
						operatorValue = "\"\""
					}
					result = strings.Join([]string{result, fieldName, "==", operatorValue}, " ")
					if condition != "" {
						if len(subsetResult) != 0 {
							subsetResult = strings.Join([]string{subsetResult, condition}, "")
						} else {
							subsetResult = condition
						}
					}
				case notEquals:
					fieldName, condition, _ := FieldNameProcessor(singleRule.Field, ctx)
					if fieldNameReplacer != "" && fieldName[len(fieldName)-3:] == "[*]" {
						fieldName = strings.Join([]string{fieldName[:len(fieldName)-3], "[", fieldNameReplacer, "]"}, "")
					}

					result = result + "\n"
					operatorValue := fmt.Sprint(singleRule.Operator.Value)
					if operatorValue == "" {
						operatorValue = "\"\""
					}
					result = strings.Join([]string{result, fieldName, "!=", operatorValue}, " ")
					if condition != "" {
						if len(subsetResult) != 0 {
							subsetResult = strings.Join([]string{subsetResult, condition}, "")
						} else {
							subsetResult = condition
						}
					}
				case greaterOrEquals:
					fieldName, condition, _ := FieldNameProcessor(singleRule.Field, ctx)
					if fieldNameReplacer != "" && fieldName[len(fieldName)-3:] == "[*]" {
						fieldName = strings.Join([]string{fieldName[:len(fieldName)-3], "[", fieldNameReplacer, "]"}, "")
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
					fieldName, condition, _ := FieldNameProcessor(singleRule.Field, ctx)
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
					fieldName, condition, _ := FieldNameProcessor(singleRule.Field, ctx)
					if fieldNameReplacer != "" && fieldName[len(fieldName)-3:] == "[*]" {
						fieldName = strings.Join([]string{fieldName[:len(fieldName)-3], "[", fieldNameReplacer, "]"}, "")
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
					fieldName, condition, _ := FieldNameProcessor(singleRule.Field, ctx)
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
							fieldName, condition, _ := FieldNameProcessor(singleRule.Field, ctx)
							fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)
							result = strings.Join([]string{result, fieldName}, " ")
							if condition != "" {
								if len(subsetResult) != 0 {
									subsetResult = strings.Join([]string{subsetResult, condition}, "")
								} else {
									subsetResult = condition
								}
							}
						} else {
							fieldName, condition, _ := FieldNameProcessor(singleRule.Field, ctx)
							fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)
							result = strings.Join([]string{result, not, fieldName}, " ")
							if condition != "" {
								if len(subsetResult) != 0 {
									subsetResult = strings.Join([]string{subsetResult, condition}, "")
								} else {
									subsetResult = condition
								}
							}
						}
					} else if reflect.Bool == reflect.TypeOf(singleRule.Operator.Value).Kind() {
						if singleRule.Operator.Value.(bool) {
							fieldName, condition, _ := FieldNameProcessor(singleRule.Field, ctx)
							fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)
							result = strings.Join([]string{result, fieldName}, " ")
							if condition != "" {
								if len(subsetResult) != 0 {
									subsetResult = strings.Join([]string{subsetResult, condition}, "")
								} else {
									subsetResult = condition
								}
							}
						} else {
							fieldName, condition, _ := FieldNameProcessor(singleRule.Field, ctx)
							fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)
							result = strings.Join([]string{result, not, fieldName}, " ")
							if condition != "" {
								if len(subsetResult) != 0 {
									subsetResult = strings.Join([]string{subsetResult, condition}, "")
								} else {
									subsetResult = condition
								}
							}
						}
					}
				case contains:
					//fmt.Printf("Field name replacer is %s\n", fieldNameReplacer)
					result = result + "\n"
					fieldName := fmt.Sprint(singleRule.Field)
					fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)

					result = strings.Join([]string{result, " ", regexExp, "(", "\"", ".*", fmt.Sprint(singleRule.Operator.Value), ".*", "\"", ",", fieldName, ")"}, "")
				case notContains:
					result = result + "\n"
					fieldName := fmt.Sprint(singleRule.Field)
					fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)

					result = strings.Join([]string{result, " ", not, " ", regexExp, "(", "\"", ".*", fmt.Sprint(singleRule.Operator.Value), ".*", "\"", ",", fmt.Sprint(singleRule.Field), ")"}, "")
				case like:
					fieldName := singleRule.Field.(string)
					fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)

					result = result + "\n"
					result = strings.Join([]string{result, " ", regexExp, "(", "\"", fmt.Sprint(singleRule.Operator.Value), "\"", ",", fieldName, ")"}, "")
				case notLike:
					fieldName := singleRule.Field.(string)
					fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)

					result = result + "\n"
					result = strings.Join([]string{result, " ", not, " ", regexExp, "(", "\"", fmt.Sprint(singleRule.Operator.Value), "\"", ",", fieldName, ")"}, "")
				case where:
					fieldName := singleRule.Field.(string)
					fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)

					var exper string
					switch singleRule.FieldOperation {
					case count:
						operator := singleRule.Operator.Value.(RuleSet)
						subsetNames, subRule, err := operator.RuleSetReader(fieldName, ctx)
						if err != nil {
							return []string{}, "", err
						}
						//fmt.Printf("The field name is %s", fieldName)
						if string(fieldName[len(fieldName)-3:]) == "[*]" {
							//fmt.Printf("here is a fieldname %+v\n", fieldName)
							fieldName, _, _ = FieldNameProcessor(fieldName[:len(fieldName)-3], ctx)
							fieldName = fieldName + "[x]"
						} else {
							fieldName, _, _ = FieldNameProcessor(fieldName, ctx)
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
					fieldName, condition, _ := FieldNameProcessor(singleRule.Field, ctx)
					// Find the common substring, replace it with the fieldNameReplacer with suffix [x]
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
				//TODO: notIn case is incorrectly addressed, should think of a way to express "not in" in rego
				case notIn:
					fieldName, condition, _ := FieldNameProcessor(singleRule.Field, ctx)
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
				}
			}
		}

		for _, thisSet := range ruleSet.RuleSets {
			if len(result) == 0 {
				result = strings.Join([]string{result, "{"}, "")
			}
			subsetNames, subRule, err := thisSet.RuleSetReader("", ctx)
			if err != nil {
				return []string{}, "", err
			}
			if len(subsetResult) != 0 {
				subsetResult = strings.Join([]string{subsetResult, subRule}, "")
			} else {
				subsetResult = subRule
			}
			//conditionNames = append(conditionNames, subsetNames...)
			for _, subsetName := range subsetNames {
				result = result + "\n"
				if len(subsetName) == andConditionLen || len(subsetName) == andConditionLenPlusX {
					result = strings.Join([]string{result, subsetName}, " ")
				} else if len(subsetName) == orConditionLen || len(subsetName) == orConditionLenPlusX {
					result = strings.Join([]string{result, not, subsetName}, " ")
				} else if len(subsetName) == singleConditionLen || len(subsetName) == singleConditionLenPlusX {
					result = strings.Join([]string{result, subsetName}, " ")
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
					fieldName, condition, _ := FieldNameProcessor(singleRule.Field, ctx)
					if fieldNameReplacer != "" && fieldName[len(fieldName)-3:] == "[*]" {
						fieldName = strings.Join([]string{fieldName[:len(fieldName)-3], "[", fieldNameReplacer, "]"}, "")
					}
					result = result + "\n"
					operatorValue := fmt.Sprint(singleRule.Operator.Value)
					if operatorValue == "" {
						operatorValue = "\"\""
					}
					result = strings.Join([]string{result, fieldName, "!=", operatorValue}, " ")
					if condition != "" {
						if len(subsetResult) != 0 {
							subsetResult = strings.Join([]string{subsetResult, condition}, "")
						} else {
							subsetResult = condition
						}
					}
				case notEquals:
					fieldName, condition, _ := FieldNameProcessor(singleRule.Field, ctx)
					if fieldNameReplacer != "" && fieldName[len(fieldName)-3:] == "[*]" {
						fieldName = strings.Join([]string{fieldName[:len(fieldName)-3], "[", fieldNameReplacer, "]"}, "")
					}
					result = result + "\n"
					operatorValue := fmt.Sprint(singleRule.Operator.Value)
					if operatorValue == "" {
						operatorValue = "\"\""
					}
					result = strings.Join([]string{result, fieldName, "==", operatorValue}, " ")
					if condition != "" {
						if len(subsetResult) != 0 {
							subsetResult = strings.Join([]string{subsetResult, condition}, "")
						} else {
							subsetResult = condition
						}
					}
				case greaterOrEquals:
					fieldName, condition, _ := FieldNameProcessor(singleRule.Field, ctx)
					if fieldNameReplacer != "" && fieldName[len(fieldName)-3:] == "[*]" {
						fieldName = strings.Join([]string{fieldName[:len(fieldName)-3], "[", fieldNameReplacer, "]"}, "")
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
					fieldName, condition, _ := FieldNameProcessor(singleRule.Field, ctx)
					if fieldNameReplacer != "" && fieldName[len(fieldName)-3:] == "[*]" {
						fieldName = strings.Join([]string{fieldName[:len(fieldName)-3], "[", fieldNameReplacer, "]"}, "")
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
					fieldName, condition, _ := FieldNameProcessor(singleRule.Field, ctx)
					if fieldNameReplacer != "" && fieldName[len(fieldName)-3:] == "[*]" {
						fieldName = strings.Join([]string{fieldName[:len(fieldName)-3], "[", fieldNameReplacer, "]"}, "")
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
							fieldName, condition, _ := FieldNameProcessor(singleRule.Field, ctx)
							fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)
							result = strings.Join([]string{result, not, fieldName}, " ")
							if condition != "" {
								if len(subsetResult) != 0 {
									subsetResult = strings.Join([]string{subsetResult, condition}, "")
								} else {
									subsetResult = condition
								}
							}
						} else {
							fieldName, condition, _ := FieldNameProcessor(singleRule.Field, ctx)
							fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)
							result = strings.Join([]string{result, fieldName}, " ")
							if condition != "" {
								if len(subsetResult) != 0 {
									subsetResult = strings.Join([]string{subsetResult, condition}, "")
								} else {
									subsetResult = condition
								}
							}
						}
					} else if reflect.Bool == reflect.TypeOf(singleRule.Operator.Value).Kind() {
						if singleRule.Operator.Value.(bool) {
							fieldName, condition, _ := FieldNameProcessor(singleRule.Field, ctx)
							fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)
							result = strings.Join([]string{result, not, fieldName}, " ")
							if condition != "" {
								if len(subsetResult) != 0 {
									subsetResult = strings.Join([]string{subsetResult, condition}, "")
								} else {
									subsetResult = condition
								}
							}
						} else {
							fieldName, condition, _ := FieldNameProcessor(singleRule.Field, ctx)
							//fmt.Printf("after processing the field name is %s\n", fieldName)
							fieldName = FieldNameReplacer(fieldName, fieldNameReplacer)
							result = strings.Join([]string{result, fieldName}, " ")
							if condition != "" {
								if len(subsetResult) != 0 {
									subsetResult = strings.Join([]string{subsetResult, condition}, "")
								} else {
									subsetResult = condition
								}
							}
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
						subsetNames, subRule, err := operator.RuleSetReader("x", ctx)
						if err != nil {
							return []string{}, "", err
						}
						if string(fieldName[len(fieldName)-3:]) == "[*]" {
							//fmt.Printf("here is a fieldname %+v\n", fieldName)
							fieldName, _, _ = FieldNameProcessor(fieldName[:len(fieldName)-3], ctx)
							fieldName = fieldName + "[x]"
						} else {
							fieldName, _, _ = FieldNameProcessor(fieldName, ctx)
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
					fieldName, condition, _ := FieldNameProcessor(singleRule.Field, ctx)
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
					fieldName, condition, _ := FieldNameProcessor(singleRule.Field, ctx)
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
			subsetNames, subRules, err := thisSet.RuleSetReader("", ctx)
			if err != nil {
				return []string{}, "", err
			}
			if len(subsetResult) != 0 {
				subsetResult = strings.Join([]string{subsetResult, subRules}, "")
			} else {
				subsetResult = subRules
			}
			//The subrule set also need to be the opposite
			for _, subsetName := range subsetNames {
				result = result + "\n"
				if len(subsetName) == andConditionLen || len(subsetName) == andConditionLenPlusX {
					result = strings.Join([]string{result, not, subsetName}, " ")
				} else if len(subsetName) == orConditionLen || len(subsetName) == orConditionLenPlusX {
					result = strings.Join([]string{result, subsetName}, " ")
				} else if len(subsetName) == singleConditionLen || len(subsetName) == singleConditionLenPlusX {
					result = strings.Join([]string{result, subsetName}, " ")
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
					fieldName, condition, _ := FieldNameProcessor(singleRule.Field, ctx)
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
					fieldName, condition, _ := FieldNameProcessor(singleRule.Field, ctx)
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
					fieldName, condition, _ := FieldNameProcessor(singleRule.Field, ctx)
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
				case notIn:
					fieldName, condition, _ := FieldNameProcessor(singleRule.Field, ctx)
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
				}
			}
		}

		for _, thisSet := range ruleSet.RuleSets {
			//fmt.Printf("The rule set is %+v\n", thisSet)
			if len(result) == 0 {
				result = strings.Join([]string{result, "{"}, "")
			}
			subsetNames, subRules, err := thisSet.RuleSetReader(fieldNameReplacer, ctx)
			if err != nil {
				return []string{}, "", err
			}
			if len(subsetResult) != 0 {
				subsetResult = strings.Join([]string{subsetResult, subRules}, "")
			} else {
				subsetResult = subRules
			}
			//conditionNames = append(conditionNames, subsetNames...)
			for _, subsetName := range subsetNames {
				result = result + "\n"
				if len(subsetName) == andConditionLen || len(subsetName) == andConditionLenPlusX {
					result = strings.Join([]string{result, subsetName}, " ")
				} else if len(subsetName) == orConditionLen || len(subsetName) == orConditionLenPlusX {
					result = strings.Join([]string{result, not, subsetName}, " ")
				} else if len(subsetName) == singleConditionLen || len(subsetName) == singleConditionLenPlusX {
					result = strings.Join([]string{result, subsetName}, " ")
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
	}
	return conditionNames, result, nil
}

var predicateFactory = map[string]Predicate{
	"":        EmptyPredicate{},
	equals:    EqualsPredicate{},
	notEquals: NotEqualsPredicate{},
	exists:    ExistsPredicate{},
	like:      LikePredicate{},
	notLike:   NotLikePredicate{},
	where:     WherePredicate{},
	in:        InPredicate{},
	notIn:     NotInPredicate{},
}

func (singleRule SingleRule) SingleRuleReader(ctx context.Context) (string, string, error) {
	predicate := predicateFactory[strings.ToLower(singleRule.Operator.Name)]
	return predicate.Evaluate(singleRule, ctx)
}

func FieldNameProcessor(fieldName interface{}, ctx context.Context) (string, string, error) {
	var result string
	var rules string
	switch fn := fieldName.(type) {
	case string:
		if fn == typeOfResource || fn == kindOfResource {
			return fn, "", nil
		}
		if strings.Contains(fn, "count") {
			return fn, "", nil
		}
		rt, err := currentResourceType(ctx)
		if err != nil {
			return "", "", err
		}
		res, err := FieldNameParser(fn, rt, "")
		if err != nil {
			return "", "", err
		}
		//fmt.Printf("before mapping %s\n", res)
		result = TFNameMapping(res)
		//fmt.Printf("after mapping %s\n", result)
	case SingleRule:
		//fmt.Printf("the field name is %+v\n", fieldName)
		res, singleRule, err := fieldName.(SingleRule).SingleRuleReader(ctx)
		if err != nil {
			return "", "", err
		}
		result = res
		rules = singleRule
	}

	return result, rules, nil
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
	//fmt.Printf("the input type is %+v\n", reflect.TypeOf(input))
	switch input.(type) {
	case []interface{}:
		for _, v := range input.([]interface{}) {
			array = append(array, "\""+fmt.Sprint(v)+"\"")
		}
	case []string:
		for _, v := range input.([]string) {
			array = append(array, "\""+fmt.Sprint(v)+"\"")
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
			if replacer[len(replacer)-3:] == "[*]" {
				newFieldNameReplacer := replacer[:len(replacer)-3] + "[x]"
				fieldName = strings.Replace(fieldName, replacer, newFieldNameReplacer, 1)
			}
		}
	}

	return fieldName
}

type LookupTable aztfq.LookupTable

func (t LookupTable) QueryProperty(resourceType, apiVersion, propertyAddress string) ([]aztfq.TFResult, bool) {
	m, ok := t.QueryResource(resourceType, apiVersion)
	if !ok {
		return nil, false
	}
	r, ok := m[propertyAddress]
	return r, ok
}

func (t LookupTable) QueryParentProperty(resourceType, apiVersion, propertyAddress string) string {
	var result string
	m, ok := t.QueryResource(resourceType, apiVersion)
	if !ok {
		return ""
	}
	_, ok = m[propertyAddress]
	if !ok {
		for k, v := range m {
			if strings.HasPrefix(k, propertyAddress) {
				childAddr := v[0].PropertyAddr
				addrArray := strings.Split(childAddr, "/")
				for i := len(addrArray) - 1; i >= 0; i-- {
					if _, err := strconv.Atoi(addrArray[i]); err == nil {
						continue
					}
					result = strings.Join(addrArray[:i], "/")
					break
				}
			}
		}
	}
	return result
}

func (t LookupTable) QueryResource(resourceType, apiVersion string) (map[string][]aztfq.TFResult, bool) {
	l2, ok := t[resourceType]
	if !ok {
		return nil, false
	}
	l3, ok := l2[apiVersion]
	if !ok {
		return nil, false
	}
	return l3, true
}

var lookupTable = func() LookupTable {
	b, err := os.ReadFile("output.json")
	if err != nil {
		panic(err.Error())
	}
	t, err := aztfq.BuildLookupTable(b, nil)
	if err != nil {
		panic(err.Error())
	}
	return LookupTable(t)
}()

func FieldNameParser(fieldNameRaw, resourceType, version string) (string, error) {
	if fieldNameRaw == typeOfResource {
		return fieldNameRaw, nil
	}
	//if strings.Contains(fieldNameRaw, "count") {
	//	return fieldNameRaw, nil
	//}
	if strings.HasPrefix(strings.ToLower(fieldNameRaw), strings.ToLower(resourceType)) {
		rtLen := len(resourceType)
		fieldNameRaw = fieldNameRaw[rtLen:]
	}
	//some attributes has "properties/" in the middle of the path after the list name, need to address this case
	prop := fieldNameRaw
	prop = strings.Replace(prop, ".", "/", -1)
	prop = strings.Replace(prop, "[x]", "/*", -1)
	prop = strings.Replace(prop, "[*]", "/*", -1)
	prop = strings.TrimPrefix(prop, "/")
	//fmt.Printf("the prop is %s\n", prop)
	upperRt := strings.ToUpper(resourceType)
	if results, ok := lookupTable.QueryProperty(upperRt, version, prop); ok {
		return results[0].PropertyAddr, nil
	}
	prop = "properties/" + prop
	if results, ok := lookupTable.QueryProperty(upperRt, version, prop); ok {
		return results[0].PropertyAddr, nil
	}
	prop = strings.Replace(prop, "*/", "*/properties/", -1)
	if results, ok := lookupTable.QueryProperty(upperRt, version, prop); ok {
		return results[0].PropertyAddr, nil
	}

	parentPropAddr := lookupTable.QueryParentProperty(upperRt, version, prop)
	if parentPropAddr != "" {
		return parentPropAddr, nil
	}

	fmt.Printf("cannot find the property %s in the lookup table\n", prop)
	prop = strings.Replace(prop, "properties/", "", -1)
	prop = ToSnakeCase(prop)
	return prop, nil
}

func ToSnakeCase(str string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func ResourceTypeParser(resourceType string) (string, error) {
	upperRt := strings.ToUpper(resourceType)
	ttt, ok := lookupTable.QueryResource(upperRt, "")
	if !ok || len(ttt) == 0 {
		return "", fmt.Errorf("cannot find the resource type %s in the lookup table", resourceType)
	}
	var result string
	for _, v := range ttt {
		result = v[0].ResourceType
		break
	}
	// The `azurerm_app_service_plan` resource has been superseded by the `azurerm_service_plan` resource.
	if result == "azurerm_app_service_plan" {
		result = "azurerm_service_plan"
	} else if result == "azurerm_app_service_environment" {
		result = "azurerm_app_service_environment_v3"
	}
	return result, nil
}

func TFNameMapping(fieldName string) string {
	var result string
	attributes := strings.Split(fieldName, "/")
	for _, v := range attributes {
		if v == "" {
			continue
		}
		next := result + "." + v
		if _, err := strconv.Atoi(v); err == nil {
			next = result + "[" + v + "]"
		}
		result = next
	}
	result = "r.change.after" + result

	return result
}
