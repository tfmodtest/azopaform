package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/spf13/afero"
)

type Rule struct {
	Properties PolicyRuleModel
	Id         string
	Name       string
}

type PolicyRuleModel struct {
	PolicyRule map[string]interface{}
	Parameters map[string]interface{}
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

var Fs = afero.NewOsFs()

var rt string
var action string

func AzurePolicyToRego(policyPath string, dir string) error {
	//policyPath := testPath
	var paths []string
	var err error

	////For batch translation
	if dir != "" {
		paths, err = jsonFiles(dir)
		if err != nil {
			return err
		}
	}

	//Override for hard cases
	//paths := []string{"/Users/jiaweitao/workZone/azure-policy/built-in-policies/policyDefinitions/Key Vault/KeyVault_SoftDeleteMustBeEnabled_Audit.json"}
	if policyPath != "" {
		paths = []string{policyPath}
	}
	for _, path := range paths {
		err = azPolicy2Rego(path)
		if err != nil {
			return err
		}
	}
	return nil
}

func azPolicy2Rego(path string) error {
	fmt.Printf("the path is %+v\n", path)

	rule, err := ruleIterator(path)
	if err != nil {
		fmt.Printf("cannot find rules %+v\n", err)
		return err
	}

	effectParams := rule.Properties.Parameters["effect"].(map[string]interface{})
	then := rule.Properties.PolicyRule["then"].(map[string]interface{})
	if effect := then["effect"]; effect != nil {
		effect = strings.ToLower(effect.(string))
		if effect == deny {
			action = deny
		} else if effect == "[parameters('effect')]" {
			defaultEffect := effectParams["defaultValue"].(string)
			defaultEffect = strings.ToLower(defaultEffect)
			if defaultEffect == deny {
				action = deny
			} else if defaultEffect == audit {
				action = warn
			} else if defaultEffect == disabled {
				action = disabled
			}
		}
	}
	fmt.Printf("the effect is %+v\n", action)
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
	if action == disabled {
		result = "default allow := true\n\n" + result
	} else if action == deny {
		top := "deny if {\n" + " " + conditionNames[0] + "\n}\n"
		result = top + result
	} else if action == warn {
		top := "warn if {\n" + " " + conditionNames[0] + "\n}\n"
		result = top + result
	}

	result = "package main\n\n" + "import rego.v1\n\n" + "r := tfplan.resource_changes[_]\n\n" + result

	err = afero.WriteFile(Fs, fileName, []byte(result), 0644)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func jsonFiles(dir string) ([]string, error) {
	res, err := readJsonFilePaths(dir)
	if err != nil {
		fmt.Printf("cannot find files in directory %+v\n", err)
		return nil, err
	}
	return res, nil
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

	entries, err := afero.ReadDir(Fs, path)
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

	data, err := afero.ReadFile(Fs, path)
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
