package main

import (
	randomnames "github.com/random-names/go"
	"log"
	"os"
)

// RuleSetReader takes a RuleSet and returns a string that can be used in a Rego file
func (ruleSet RuleSet) RuleSetReader() (string, string, error) {
	var result string
	var conditionNames []string

	switch ruleSet.Flag {
	case allOf:
		var subsetResult string
		result = result + "{"
		for _, singleRule := range ruleSet.SingleRules {
			switch singleRule.Operator.Name {
			case equals:
				fieldName := singleRule.Field.(string)
				result = result + "\n" + fieldName + " == " + singleRule.Operator.Value.(string) + "\n"
			case notEquals:
				fieldName := singleRule.Field.(string)
				result = result + "\n" + fieldName + " != " + singleRule.Operator.Value.(string) + "\n"
			}
		}
		for _, thisSet := range ruleSet.RuleSets {
			subsetName, subRule, err := thisSet.RuleSetReader()
			if err != nil {
				return "", "", err
			}
			subsetResult = subRule
			conditionNames = append(conditionNames, subsetName)
			result = result + "\n" + subsetName + "\n"
		}

		result = result + "}"
		result = result + "\n" + subsetResult + "\n"

		conditionName, _ := randomnames.GetRandomName("condition", nil)

		result = conditionName + " " + result

		return conditionName, result, nil
	case anyOf:
		var subsetResult string
		for _, singleRule := range ruleSet.SingleRules {
			switch singleRule.Operator.Name {
			case equals:
				fieldName := singleRule.Field.(string)
				result = result + "\n" + fieldName + " != " + singleRule.Operator.Value.(string) + "\n"
			case notEquals:
				fieldName := singleRule.Field.(string)
				result = result + "\n" + fieldName + " == " + singleRule.Operator.Value.(string) + "\n"
			}
		}
		for _, thisSet := range ruleSet.RuleSets {
			subsetName, subsetResult, err := thisSet.RuleSetReader()
			if err != nil {
				return "", "", err
			}
			conditionNames = append(conditionNames, subsetName)
			result = result + "\n" + "not " + subsetResult + "\n"
		}
		result = result + "}"
		result = result + "\n" + subsetResult + "\n"

		conditionName, _ := randomnames.GetRandomName("condition", nil)

		result = conditionName + " " + result

		return conditionName, result, nil
	}
	return "", result, nil
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

//
//func (ruleSet RuleSet) RuleSetParser() (string, error) {
//
//}
