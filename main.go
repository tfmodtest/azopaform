package main

import (
	"encoding/json"
	"errors"
	"flag"
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

var rt string

func main() {
	singlePath := flag.String("path", "", "The path of policy definition file")
	dir := flag.String("dir", "", "The dir which contains policy definitions")
	flag.Parse()
	if err := realMain(*singlePath, *dir); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func realMain(policyPath string, dir string) error {
	//policyPath := testPath
	var paths []string
	keyWordsCollection := make(map[string][]string)
	operators := make(map[string]bool)

	////For batch translation
	if dir != "" {
		res, err := readJsonFilePaths(dir)
		if err != nil {
			fmt.Printf("cannot find files in directory %+v\n", err)
			return err
		}
		paths = res
		for _, path := range paths {
			//words, operatorSet, err := ruleIterator(path)
			rule, err := ruleIterator(path)
			if err != nil {
				fmt.Printf("cannot find rules %+v\n", err)
				return err
			}

			words, operatorSet, err := rule.Properties.listKeyWords()
			for k, v := range operatorSet {
				operators[k] = v
			}
			keyWordsCollection[path] = words
		}
	}

	//Override for hard cases
	//paths := []string{"/Users/jiaweitao/workZone/azure-policy/built-in-policies/policyDefinitions/Key Vault/KeyVault_SoftDeleteMustBeEnabled_Audit.json"}
	if policyPath != "" {
		paths = []string{policyPath}
	}
	for _, path := range paths {
		fmt.Printf("the path is %+v\n", path)

		rule, err := ruleIterator(path)
		if err != nil {
			fmt.Printf("cannot find rules %+v\n", err)
			return err
		}

		conditions := rule.Properties.PolicyRule["if"]
		condition, err := conditionFinder(conditions.(map[string]interface{}))
		if err != nil {
			fmt.Printf("cannot find conditions %+v\n", err)
			return err
		}
		fmt.Printf("the resource type is %+v\n", rt)
		fmt.Printf("the whole condition is %+v\n", *condition)
		fileName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)) + ".rego"
		conditionNames, result, err := condition.RuleSetReader("")
		fmt.Printf("the condition names are %+v\n", conditionNames)
		result = "package main\n\n" + "import rego.v1\n\n" + result
		err = os.WriteFile(fileName, []byte(result), 0644)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	jsonData, err := json.MarshalIndent(keyWordsCollection, "", " ")
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = os.WriteFile("keyWords.json", jsonData, 0644)
	if err != nil {
		fmt.Println(err)
		return err
	}

	jsonSet, err := json.MarshalIndent(operators, "", " ")
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = os.WriteFile("operators.json", jsonSet, 0644)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
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
			//fmt.Printf("the not rule is %+v\n", *rule)

			notRules.RuleSets = append(notRules.RuleSets, *rule)
			notRules.SingleRules = append(notRules.SingleRules, rule.SingleRules...)
			return &notRules, nil
		case where:
			whereConditions := v.(map[string]interface{})
			//fmt.Printf("the where conditions are %+v\n", whereConditions)
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

	if fieldName == typeOfResource {
		switch operatorValue.(type) {
		case string:
			rt = operatorValue.(string)
			v, err := ResourceTypeParser(rt)
			if err != nil {
				fmt.Printf("cannot find resource type %+v\n", err)
				return nil, err
			}
			singleRule.Operator.Value = v
		case []interface{}:
			var res []string
			for _, v := range operatorValue.([]interface{}) {
				parsedType, err := ResourceTypeParser(v.(string))
				if err != nil {
					fmt.Printf("cannot find resource type %+v\n", err)
					return nil, err
				}
				res = append(res, parsedType)
			}
			singleRule.Operator.Value = res
		}
	}

	var singleRules []SingleRule
	return &RuleSet{
		Flag:        "single",
		SingleRules: append(singleRules, singleRule),
	}, nil
}
