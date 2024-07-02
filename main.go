package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type Rule struct {
	Properties PolicyRuleModel
	Id         string
	Name       string
}

type PolicyRuleModel struct {
	PolicyRule map[string]interface{}
}

type RuleSet struct {
	Flag        string
	SingleRules []SingleRule
	RuleSets    []RuleSet
}

type SingleRule struct {
	Field          any
	FieldOperation string
	Operator       OperatorModel
}

type OperatorModel struct {
	Name  string
	Value any
}

type CountOperatorModel[T any] struct {
	FieldName string
	Condition []RuleSet
	Operator  OperatorModel
}

const path = "/home/jiawei/workZone/azure-policy/built-in-policies/policyDefinitions"
const testPath = "/Users/jiaweitao/workZone/azure-policy/built-in-policies/policyDefinitions/Key Vault"

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

func main() {
	//policyPath := testPath

	keyWordsCollection := make(map[string][]string)
	operators := make(map[string]bool)

	////For batch translation
	//paths, err := readJsonFilePaths(policyPath)
	//if err != nil {
	//	fmt.Printf("cannot find files in directory %+v\n", err)
	//	return
	//}
	//for _, path := range paths {
	//	//words, operatorSet, err := ruleIterator(path)
	//	rule, err := ruleIterator(path)
	//	if err != nil {
	//		fmt.Printf("cannot find rules %+v\n", err)
	//		return
	//	}
	//
	//	words, operatorSet, err := rule.Properties.listKeyWords()
	//	for k, v := range operatorSet {
	//		operators[k] = v
	//	}
	//	keyWordsCollection[path] = words
	//}

	//Override for hard cases
	paths := []string{"/Users/jiaweitao/workZone/azure-policy/built-in-policies/policyDefinitions/Key Vault/KeyVault_SoftDeleteMustBeEnabled_Audit.json"}

	for _, path := range paths {
		rule, err := ruleIterator(path)
		if err != nil {
			fmt.Printf("cannot find rules %+v\n", err)
			return
		}

		conditions := rule.Properties.PolicyRule["if"]
		condition, err := conditionFinder(conditions.(map[string]interface{}))
		if err != nil {
			fmt.Printf("cannot find conditions %+v\n", err)
			return
		}
		fmt.Printf("the whole condition is %+v\n", *condition)
		fileName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)) + ".rego"
		conditionNames, result, err := condition.RuleSetReader("")
		fmt.Printf("the condition names are %+v\n", conditionNames)
		result = "package main\n\n" + "import rego.v1\n\n" + result
		err = os.WriteFile(fileName, []byte(result), 0644)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	jsonData, err := json.MarshalIndent(keyWordsCollection, "", " ")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = os.WriteFile("keyWords.json", jsonData, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	jsonSet, err := json.MarshalIndent(operators, "", " ")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = os.WriteFile("operators.json", jsonSet, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (policyRule PolicyRuleModel) listKeyWords() ([]string, map[string]bool, error) {
	var words []string
	operatorSet := make(map[string]bool)
	for k, v := range policyRule.PolicyRule {
		words = append(words, k)
		if k == "if" && reflect.TypeOf(v) == reflect.TypeOf(map[string]interface{}{}) {
			result, error := findAllOperators(v.(map[string]interface{}))
			//operatorSet = result
			if error != nil {
				fmt.Printf("cannot find operators %+v\n", error)
				return nil, nil, error
			}
			for key, value := range result {
				operatorSet[key] = value
			}
		}
	}
	return words, operatorSet, nil
}

func findAllOperators(entries map[string]interface{}) (map[string]bool, error) {
	operatorSet := make(map[string]bool)
	//fmt.Printf("the entry type is %+v\n", reflect.TypeOf(entries))
	for k, v := range entries {
		operatorSet[k] = true
		if reflect.TypeOf(v) != reflect.TypeOf("") {
			if reflect.TypeOf(v) == reflect.TypeOf([]interface{}{}) {
				for _, value := range v.([]interface{}) {
					if reflect.TypeOf(value) == reflect.TypeOf(map[string]interface{}{}) {
						subSet, error := findAllOperators(value.(map[string]interface{}))
						if error != nil {
							fmt.Printf("cannot find operators %+v\n", error)
							return nil, error
						}
						for key, value := range subSet {
							operatorSet[key] = value
						}
					}
				}
				continue
			}
			if reflect.TypeOf(v) == reflect.TypeOf(map[string]interface{}{}) {
				subSet, error := findAllOperators(v.(map[string]interface{}))
				if error != nil {
					fmt.Printf("cannot find operators %+v\n", error)
					return nil, error
				}
				for key, value := range subSet {
					operatorSet[key] = value
				}
				continue
			}
		}
	}
	return operatorSet, nil
}

func readJsonFilePaths(path string) ([]string, error) {
	var filePaths []string
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("cannot find files in directory %+v\n", err)
		return nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			subFilePaths, err := readJsonFilePaths(path + "/" + entry.Name())
			if err != nil {
				return nil, err
			}
			filePaths = append(filePaths, subFilePaths...)
		}
		if filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		fileName := entry.Name()
		filePath := path + "/" + fileName
		filePaths = append(filePaths, filePath)
	}
	return filePaths, nil
}

func ruleIterator(path string) (*Rule, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("unable to read file content %+v\n\n", err)
		return nil, err
		//return nil, nil, err
	}

	var rule Rule
	err = json.Unmarshal(data, &rule)
	if err != nil {
		fmt.Printf("unable to unmarshal json %+v\n\n", err)
		return nil, err
	}

	return &rule, nil
}

// In this function, the input is a single condition or a set of conditions(allof, anyof)
func conditionFinder(conditions map[string]interface{}) (*RuleSet, error) {
	if conditions == nil {
		err := errors.New("cannot find conditions")
		return nil, err
	}

	var fieldName any
	var operatorName string
	var operatorValue any
	//var operatorValueType reflect.Type
	andRules := RuleSet{
		Flag: allOf,
	}
	orRules := RuleSet{
		Flag: anyOf,
	}
	whereRules := RuleSet{
		Flag: where,
	}
	notRules := RuleSet{
		Flag: not,
	}

	for k, v := range conditions {
		switch strings.ToLower(k) {
		case allOf:
			allOfConditions := v.([]interface{})
			for _, condition := range allOfConditions {
				rule, err := conditionFinder(condition.(map[string]interface{}))
				if err != nil {
					fmt.Printf("cannot find AND conditions %+v\n", err)
					return nil, err
				}
				if rule.Flag == allOf || rule.Flag == anyOf || rule.Flag == where || rule.Flag == count || rule.Flag == not {
					andRules.RuleSets = append(andRules.RuleSets, *rule)
				} else {
					andRules.SingleRules = append(andRules.SingleRules, rule.SingleRules...)
				}
			}
			return &andRules, nil
		case anyOf:
			anyOfConditions := v.([]interface{})
			for _, condition := range anyOfConditions {
				rule, err := conditionFinder(condition.(map[string]interface{}))
				if err != nil {
					fmt.Printf("cannot find OR conditions %+v\n", err)
					return nil, err
				}
				if rule.Flag == allOf || rule.Flag == anyOf || rule.Flag == where || rule.Flag == count || rule.Flag == not {
					orRules.RuleSets = append(orRules.RuleSets, *rule)
				} else {
					orRules.SingleRules = append(orRules.SingleRules, rule.SingleRules...)
				}
			}
			return &orRules, nil
		case not:
			notCondition := v.(map[string]interface{})
			rule, err := conditionFinder(notCondition)
			if err != nil {
				fmt.Printf("cannot find NOT conditions %+v\n", err)
				return nil, err
			}
			fmt.Printf("the not rule is %+v\n", *rule)

			notRules.RuleSets = append(notRules.RuleSets, *rule)
			notRules.SingleRules = append(notRules.SingleRules, rule.SingleRules...)
			return &notRules, nil
		case where:
			whereConditions := v.(map[string]interface{})
			fmt.Printf("the where conditions are %+v\n", whereConditions)
			rule, err := conditionFinder(whereConditions)
			if err != nil {
				fmt.Printf("cannot find WHERE conditions %+v\n", err)
				return nil, err
			}
			if rule.Flag == allOf || rule.Flag == anyOf || rule.Flag == where || rule.Flag == count || rule.Flag == not {
				whereRules.RuleSets = append(whereRules.RuleSets, *rule)
			} else {
				whereRules.SingleRules = append(whereRules.SingleRules, rule.SingleRules...)
			}

			operatorName = where
			operatorValue = whereRules
		case count:
			countConditions := v.(map[string]interface{})
			rule, err := conditionFinder(countConditions)
			if err != nil {
				fmt.Printf("cannot find COUNT conditions %+v\n", err)
				return nil, err
			}

			countSingleRule := SingleRule{}
			countSingleRule.Field = rule.SingleRules[0].Field
			countSingleRule.FieldOperation = count
			countSingleRule.Operator = rule.SingleRules[0].Operator
			fieldName = countSingleRule
		case field:
			fieldName = v.(string)
		case value:
			fieldName = v.(string)
		default:
			operatorName = k
			operatorValue = v
		}
	}
	operator := OperatorModel{
		Name:  operatorName,
		Value: operatorValue,
	}

	singleRule := SingleRule{
		Field:    fieldName,
		Operator: operator,
	}

	var singleRules []SingleRule
	return &RuleSet{
		Flag:        "single",
		SingleRules: append(singleRules, singleRule),
	}, nil
}
