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
		result = strings.Join([]string{result, "{"}, " ")
		result = result + "\n"
		result = strings.Join([]string{result, res}, " ")
		if len(result) != 0 {
			result = result + "\n" + "}"
		}

		conditionName := RandStringFromCharSet(singleConditionLen, charNum)
		conditionNames = append(conditionNames, conditionName)
		result = strings.Join([]string{conditionName, ifCondition, result}, " ")

		return conditionNames, result, nil
	case allOf:
		var subsetResult string
		if len(ruleSet.SingleRules) != 0 {
			result = strings.Join([]string{result, "{"}, "")
			for _, singleRule := range ruleSet.SingleRules {
				//fmt.Printf("here is a singleRule with operator %+v\n", singleRule.Operator.Name)
				switch singleRule.Operator.Name {
				case equals:
					fieldName, condition := FieldNameProcessor(singleRule.Field)
					if fieldNameReplacer != "" {
						fieldName = fieldNameReplacer
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
					fieldName, condition := FieldNameProcessor(singleRule.Field)
					if fieldNameReplacer != "" {
						fieldName = fieldNameReplacer
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
						fieldName = fieldNameReplacer
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
				case less:
					fieldName, condition := FieldNameProcessor(singleRule.Field)
					if fieldNameReplacer != "" {
						fieldName = fieldNameReplacer
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
				case exists:
					result = result + "\n"
					if reflect.String == reflect.TypeOf(singleRule.Operator.Value).Kind() {
						if strings.EqualFold(singleRule.Operator.Value.(string), "true") {
							fieldName := singleRule.Field.(string)
							result = strings.Join([]string{result, fieldName}, " ")
						} else {
							fieldName := singleRule.Field.(string)
							result = strings.Join([]string{result, fieldName, not}, " ")
						}
					} else if reflect.Bool == reflect.TypeOf(singleRule.Operator.Value).Kind() {
						if singleRule.Operator.Value.(bool) {
							fieldName := singleRule.Field.(string)
							result = strings.Join([]string{result, fieldName}, " ")
						} else {
							fieldName := singleRule.Field.(string)
							result = strings.Join([]string{result, fieldName, not}, " ")
						}
					}
				case contains:
					result = result + "\n"
					result = strings.Join([]string{result, "", regexExp, "(", "\"", ".*", fmt.Sprint(singleRule.Operator.Value), ".*", "\"", ",", fmt.Sprint(singleRule.Field), ")"}, "")
				case like:
					fieldName := singleRule.Field.(string)
					if fieldNameReplacer != "" {
						fieldName = fieldNameReplacer
					}
					result = result + "\n"
					result = strings.Join([]string{result, " ", regexExp, "(", "\"", fmt.Sprint(singleRule.Operator.Value), "\"", ",", fieldName, ")"}, "")
				case where:
					//fmt.Printf("here is a where case %+v\n", singleRule)
					fieldName := singleRule.Field.(string)
					if fieldNameReplacer != "" {
						fieldName = fieldNameReplacer
					}
					var exper string
					switch singleRule.FieldOperation {
					case count:
						operator := singleRule.Operator.Value.(RuleSet)
						subsetNames, subRule, err := operator.RuleSetReader("")
						if err != nil {
							return []string{}, "", err
						}
						exper = count + "(" + "{" + "x" + "|" + fieldName + "[x]" + ";" + subsetNames[0] + "}" + ")"
						result = strings.Join([]string{result, " ", exper, fmt.Sprint(singleRule.Operator.Value)}, " ")
						if len(subsetResult) != 0 {
							subsetResult = strings.Join([]string{subsetResult, subRule}, "")
						} else {
							subsetResult = subRule
						}
					}
				case in:
					fieldName, condition := FieldNameProcessor(singleRule.Field)
					if fieldNameReplacer != "" {
						fieldName = fieldNameReplacer
					}
					result = result + "\n"
					result = strings.Join([]string{result, "some", fieldName, "in", fmt.Sprint(singleRule.Operator.Value)}, " ")
					if condition != "" {
						if len(subsetResult) != 0 {
							subsetResult = strings.Join([]string{subsetResult, condition}, "")
						} else {
							subsetResult = condition
						}
					}
				//TODO: notIn case is incorrectly addressed, should think of a way to espress "not in" in rego
				case notIn:
					fieldName, condition := FieldNameProcessor(singleRule.Field)
					if fieldNameReplacer != "" {
						fieldName = fieldNameReplacer
					}
					result = result + "\n"
					result = strings.Join([]string{result, "some", fieldName, "in", fmt.Sprint(singleRule.Operator.Value)}, " ")
					if condition != "" {
						if len(subsetResult) != 0 {
							subsetResult = strings.Join([]string{subsetResult, condition}, "")
						} else {
							subsetResult = condition
						}
					}
				case not:
					operatorValue := singleRule.Operator.Value.(map[string]interface{})
					//notRule := SingleRule{}
					//err := mapstructure.Decode(operatorValue, &notRule)
					//if err != nil {
					//	return []string{}, "", err
					//}
					//fmt.Printf("[WARN]here is a not case %+v\n", notRule)
					res, err := conditionFinder(operatorValue)
					if err != nil {
						return []string{}, "", err
					}
					//fmt.Printf("[WARN]here is a not case %+v\n", res)

					subsetNames, subRules, err := res.RuleSetReader("")
					if err != nil {
						return []string{}, "", err
					}
					if len(subsetResult) != 0 {
						subsetResult = strings.Join([]string{subsetResult, subRules}, "")
					} else {
						subsetResult = subRules
					}

					for _, subnetName := range subsetNames {
						result = result + "\n"
						if len(subnetName) == andConditionLen || len(subnetName) == singleConditionLen {
							result = strings.Join([]string{result, not, subnetName}, " ")
						} else if len(subnetName) == orConditionLen {
							result = strings.Join([]string{result, subnetName}, " ")
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
				switch singleRule.Operator.Name {
				case equals:
					fieldName, condition := FieldNameProcessor(singleRule.Field)
					if fieldNameReplacer != "" {
						fieldName = fieldNameReplacer
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
						fieldName = fieldNameReplacer
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
						fieldName = fieldNameReplacer
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
						fieldName = fieldNameReplacer
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
				case exists:
					result = result + "\n"
					if reflect.String == reflect.TypeOf(singleRule.Operator.Value).Kind() {
						if strings.EqualFold(singleRule.Operator.Value.(string), "true") {
							fieldName := singleRule.Field.(string)
							result = strings.Join([]string{result, fieldName, not}, " ")
						} else {
							fieldName := singleRule.Field.(string)
							result = strings.Join([]string{result, fieldName}, " ")
						}
					} else if reflect.Bool == reflect.TypeOf(singleRule.Operator.Value).Kind() {
						if singleRule.Operator.Value.(bool) {
							fieldName := singleRule.Field.(string)
							result = strings.Join([]string{result, fieldName, not}, " ")
						} else {
							fieldName := singleRule.Field.(string)
							result = strings.Join([]string{result, fieldName}, " ")
						}
					}
				case contains:
					result = result + "\n"
					result = strings.Join([]string{result, "", not, " ", regexExp, "(", "\"", ".*", fmt.Sprint(singleRule.Operator.Value), ".*", "\"", ",", fmt.Sprint(singleRule.Field), ")"}, "")
				case like:
					fieldName := singleRule.Field.(string)
					result = result + "\n"
					result = strings.Join([]string{result, " ", not, " ", regexExp, "(", fmt.Sprint(singleRule.Operator.Value), ",", fieldName, ")"}, "")
				case where:
					fieldName := singleRule.Field.(string)
					var exper string
					switch singleRule.FieldOperation {
					case count:
						operator := singleRule.Operator.Value.(RuleSet)
						subsetNames, subRule, err := operator.RuleSetReader("x")
						if err != nil {
							return []string{}, "", err
						}
						exper = count + "(" + "{" + "x" + "|" + fieldName + "[x]" + ";" + subsetNames[0] + "}" + ")"
						result = strings.Join([]string{result, " ", exper, fmt.Sprint(singleRule.Operator.Value)}, " ")
						if len(subsetResult) != 0 {
							subsetResult = strings.Join([]string{subsetResult, subRule}, "")
						} else {
							subsetResult = subRule
						}
					}
				//TODO: notIn case is incorrectly addressed, should think of a way to espress "not in" in rego
				case in:
					fieldName, condition := FieldNameProcessor(singleRule.Field)
					if fieldNameReplacer != "" {
						fieldName = fieldNameReplacer
					}
					result = result + "\n"
					result = strings.Join([]string{result, "some", fieldName, "in", fmt.Sprint(singleRule.Operator.Value)}, " ")
					if condition != "" {
						if len(subsetResult) != 0 {
							subsetResult = strings.Join([]string{subsetResult, condition}, "")
						} else {
							subsetResult = condition
						}
					}
				case notIn:
					fieldName, condition := FieldNameProcessor(singleRule.Field)
					if fieldNameReplacer != "" {
						fieldName = fieldNameReplacer
					}
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
			subsetNames, subRules, err := thisSet.RuleSetReader("")
			if err != nil {
				return []string{}, "", err
			}
			if len(subsetResult) != 0 {
				subsetResult = strings.Join([]string{subsetResult, subRules}, "")
			} else {
				subsetResult = subRules
			}
			for _, subnetName := range subsetNames {
				result = result + "\n"
				if len(subnetName) == andConditionLen {
					result = strings.Join([]string{result, subnetName}, " ")
				} else if len(subnetName) == orConditionLen {
					result = strings.Join([]string{result, not, subnetName}, " ")
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
				switch singleRule.Operator.Name {
				case equals:
					fieldName, condition := FieldNameProcessor(singleRule.Field)
					if fieldNameReplacer != "" {
						fieldName = fieldNameReplacer
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
					fieldName, condition := FieldNameProcessor(singleRule.Field)
					if fieldNameReplacer != "" {
						fieldName = fieldNameReplacer
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
				case exists:
					result = result + "\n"
					if strings.EqualFold(singleRule.Operator.Value.(string), "true") {
						fieldName := singleRule.Field.(string)
						result = strings.Join([]string{result, fieldName}, " ")
					} else {
						fieldName := singleRule.Field.(string)
						result = strings.Join([]string{result, fieldName, not}, " ")
					}
				case contains:
					result = result + "\n"
					result = strings.Join([]string{result, "", regexExp, "(", "\"", ".*", fmt.Sprint(singleRule.Operator.Value), ".*", "\"", ",", fmt.Sprint(singleRule.Field), ")"}, "")
				case like:
					fieldName := singleRule.Field.(string)
					result = result + "\n"
					result = strings.Join([]string{result, " ", regexExp, "(", fmt.Sprint(singleRule.Operator.Value), ",", fieldName, ")"}, "")
				}
			}
		}

		for _, thisSet := range ruleSet.RuleSets {
			if len(result) == 0 {
				result = strings.Join([]string{result, "{"}, "")
			}
			subsetNames, subRules, err := thisSet.RuleSetReader("x")
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
				} else {
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
	}
	return conditionNames, result, nil
}

func (singleRule SingleRule) SingleRuleReader() (string, string, error) {
	var result string
	var rules string
	switch singleRule.Operator.Name {
	case equals:
		fieldName := singleRule.Field.(string)
		result = strings.Join([]string{result, fieldName, "==", fmt.Sprint(singleRule.Operator.Value)}, " ")
	case notEquals:
		fieldName := singleRule.Field.(string)
		result = strings.Join([]string{result, fieldName, "!=", fmt.Sprint(singleRule.Operator.Value)}, " ")
	case exists:
		if strings.EqualFold(singleRule.Operator.Value.(string), "true") {
			fieldName := singleRule.Field.(string)
			result = strings.Join([]string{result, fieldName}, " ")
		} else {
			fieldName := singleRule.Field.(string)
			result = strings.Join([]string{result, fieldName, not}, " ")
		}
	case like:
		fieldName := singleRule.Field.(string)
		result = strings.Join([]string{result, " ", regexExp, "(", fmt.Sprint(singleRule.Operator.Value), ",", fieldName, ")"}, "")
	case where:
		fieldName := singleRule.Field.(string)
		operator := singleRule.Operator.Value.(RuleSet)
		subsetNames, subRule, err := operator.RuleSetReader("x")
		if err != nil {
			return "", "", err
		}
		rules = subRule
		exp := count + "(" + "{" + "x" + "|" + fieldName + "[x]" + ";" + subsetNames[0] + "}" + ")"
		result = strings.Join([]string{result, " ", exp}, " ")
	case in:
		fieldName := singleRule.Field.(string)
		result = strings.Join([]string{result, "some", fieldName, "in", fmt.Sprint(singleRule.Operator.Value)}, " ")
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
