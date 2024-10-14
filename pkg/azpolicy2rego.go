package pkg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/emirpasic/gods/stacks/arraystack"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/emirpasic/gods/stacks"
	"github.com/spf13/afero"
)

type Rule struct {
	Properties *PolicyRuleModel
	Id         string
	Name       string
}

func NewRule(input map[string]any, ctx context.Context) *Rule {
	return &Rule{
		Id:         input["id"].(string),
		Name:       input["name"].(string),
		Properties: NewPolicyRuleModel(input["properties"].(map[string]any), ctx),
	}
}

func NewPolicyRuleModel(input map[string]any, ctx context.Context) *PolicyRuleModel {
	return &PolicyRuleModel{
		DisplayName: input["displayName"].(string),
		PolicyType:  input["policyType"].(string),
		Mode:        input["mode"].(string),
		Description: input["description"].(string),
		Version:     input["version"].(string),
		Metadata:    NewPolicyRuleMetaData(input["metadata"].(map[string]any)),
		PolicyRule:  NewPolicyRuleBody(input["policyRule"].(map[string]any), ctx),
		Parameters:  NewPolicyRuleParameters(input["parameters"].(map[string]any)),
	}
}

func NewPolicyRuleBody(input map[string]any, ctx context.Context) *PolicyRuleBody {
	//ifBody := input["if"]

	conditionMap := input
	var subject Rego
	var creator func(subject Rego, input any) Rego
	var cv any
	for key, conditionValue := range conditionMap {
		key = strings.ToLower(key)
		if key == count {
			operationFactory, ok := operatorFactories[key]
			if !ok {
				panic(fmt.Sprintf("unknown operation: %s", key))
			}
			//fmt.Printf("the condition value is %v\n", conditionValue)
			conditionSet := operationFactory(conditionValue, ctx)
			subject = conditionSet
			continue
		}
		if key == allOf {
			operationFactory, ok := operatorFactories[key]
			if !ok {
				panic(fmt.Sprintf("unknown operation: %s", key))
			}
			conditionSet := operationFactory(conditionValue, ctx)
			return &PolicyRuleBody{
				Then:   nil,
				IfBody: conditionSet,
			}
		}
		if key == anyOf {
			operationFactory, ok := operatorFactories[key]
			if !ok {
				panic(fmt.Sprintf("unknown operation: %s", key))
			}
			conditionSet := operationFactory(conditionValue, ctx)
			return &PolicyRuleBody{
				Then:   nil,
				IfBody: conditionSet,
			}
		}
		if key == not {
			operationFactory, ok := operatorFactories[key]
			if !ok {
				panic(fmt.Sprintf("unknown operation: %s", key))
			}
			conditionSet := operationFactory(conditionValue, ctx)
			return &PolicyRuleBody{
				Then:   nil,
				IfBody: conditionSet,
			}
		}
		if key == field {
			if conditionValue == typeOfResource {
				pushResourceType(context.Background(), conditionValue.(string))
			}
			subject = OperationField(conditionValue.(string))
			continue
		}
		if key == value {
			subject = OperationValue(conditionValue.(string))
			continue
		}
		factory, ok := conditionFactory[key]
		if !ok {
			panic(fmt.Sprintf("unknown condition: %s", key))
		}
		creator = factory
		cv = conditionValue
	}
	return &PolicyRuleBody{
		Then:   nil,
		IfBody: creator(subject, cv),
	}
}

func NewPolicyRuleMetaData(input map[string]any) *PolicyRuleMetaData {
	return &PolicyRuleMetaData{
		Version:  input["version"].(string),
		Category: input["category"].(string),
	}
}

type PolicyRuleMetaData struct {
	Version  string
	Category string
}

type PolicyRuleModel struct {
	PolicyRule  *PolicyRuleBody
	Parameters  *PolicyRuleParameters
	DisplayName string
	PolicyType  string
	Mode        string
	Description string
	Version     string
	Metadata    *PolicyRuleMetaData
}

type PolicyRuleBody struct {
	Then   *ThenBody
	If     map[string]any `json:"if,omitempty"`
	IfBody Rego
}

type IfBody map[string]any

func (i IfBody) condition(ctx context.Context) (*RuleSet, error) {
	return conditionFinder(i, ctx)
}

func (p *PolicyRuleBody) GetThen() *ThenBody {
	if p == nil {
		return nil
	}
	return p.Then
}

func (p *PolicyRuleBody) GetIf() IfBody {
	if p == nil {
		return nil
	}
	return p.If
}

type ThenBody struct {
	Effect string `json:"effect,omitempty"`
}

func (t *ThenBody) GetEffect() string {
	if t == nil {
		return ""
	}
	return t.Effect
}

func (t *ThenBody) MapEffectToAction(defaultEffect string) (string, error) {
	effect := t.GetEffect()
	if effect == "" {
		return "", fmt.Errorf("unexpected input, effect is %s, defaultEffect is %s", effect, defaultEffect)
	}
	effect = strings.ToLower(effect)
	if effect == deny {
		return deny, nil
	}
	if effect != "[parameters('effect')]" {
		return "", fmt.Errorf("unexpected input, effect is %s, defaultEffect is %s", effect, defaultEffect)
	}
	defaultEffect = strings.ToLower(defaultEffect)
	if defaultEffect == audit {
		return warn, nil
	}
	if defaultEffect == deny || defaultEffect == disabled {
		return defaultEffect, nil
	}
	return "", fmt.Errorf("unexpected input, effect is %s, defaultEffect is %s", effect, defaultEffect)
}

type PolicyRuleParameterType string

const (
	PolicyRuleParameterTypeString   PolicyRuleParameterType = "string"
	PolicyRuleParameterTypeArray    PolicyRuleParameterType = "array"
	PolicyRuleParameterTypeObject   PolicyRuleParameterType = "object"
	PolicyRuleParameterTypeBool     PolicyRuleParameterType = "boolean"
	PolicyRuleParameterTypeInteger  PolicyRuleParameterType = "integer"
	PolicyRuleParameterTypeFloat    PolicyRuleParameterType = "float"
	PolicyRuleParameterTypeDateTime PolicyRuleParameterType = "dateTime"
)

type PolicyRuleParameterMetaData struct {
	Description string
	DisplayName string
	Deprecated  bool
}

func NewPolicyRuleParameterMetaData(input map[string]any) *PolicyRuleParameterMetaData {
	if input == nil {
		return nil
	}
	return &PolicyRuleParameterMetaData{
		Deprecated:  input["deprecated"].(bool),
		Description: input["description"].(string),
	}
}

type PolicyRuleParameter struct {
	Name         string
	Type         PolicyRuleParameterType
	DefaultValue any
	MetaData     *PolicyRuleParameterMetaData
}

func NewPolicyRuleParameter(input map[string]any) *PolicyRuleParameter {
	if input == nil {
		return nil
	}
	r := &PolicyRuleParameter{
		Name:         input["name"].(string),
		Type:         input["type"].(PolicyRuleParameterType),
		DefaultValue: input["defaultValue"],
	}
	if metaData, ok := input["metadata"].(map[string]any); ok {
		r.MetaData = NewPolicyRuleParameterMetaData(metaData)
	}
	return r
}

type PolicyRuleParameters struct {
	Effect     *EffectBody
	Parameters map[string]*PolicyRuleParameter
}

func NewPolicyRuleParameters(input map[string]any) *PolicyRuleParameters {
	if input == nil {
		return nil
	}
	parameters := make(map[string]*PolicyRuleParameter)
	for k, v := range input {
		i, ok := v.(map[string]any)
		if !ok {
			continue
		}
		parameters[k] = NewPolicyRuleParameter(i)
	}
	return &PolicyRuleParameters{
		Parameters: parameters,
	}
}

func (p *PolicyRuleParameters) GetEffect() *EffectBody {
	if p == nil {
		return nil
	}
	return p.Effect
}

type EffectBody struct {
	DefaultValue string `json:"defaultValue"`
}

func (e *EffectBody) GetDefaultValue() string {
	if e == nil {
		return ""
	}
	return e.DefaultValue
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

func AzurePolicyToRego(policyPath string, dir string, ctx context.Context) error {
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
		err = NeoAzPolicy2Rego(path, ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewContext() context.Context {
	contextMap := make(map[string]stacks.Stack)
	contextMap["resourceType"] = arraystack.New()
	contextMap["fieldNameReplacer"] = arraystack.New()
	ctx := context.WithValue(context.Background(), "context", contextMap)
	return ctx
}

func NeoAzPolicy2Rego(path string, ctx context.Context) error {
	var action string
	var conditionName string
	fmt.Printf("the path is %+v\n", path)

	rule, err := ruleIterator(path)
	if err != nil {
		fmt.Printf("cannot find rules %+v\n", err)
		return err
	}

	then := rule.Properties.PolicyRule.GetThen()
	action, err = then.MapEffectToAction(rule.Properties.Parameters.GetEffect().GetDefaultValue())
	if err != nil {
		return err
	}
	fmt.Printf("the effect is %+v\n", action)
	condition := rule.Properties.PolicyRule.GetIf()
	fmt.Printf("the condition is %+v\n", condition)
	//_, err = currentResourceType(ctx)
	//if err != nil {
	//	return err
	//}
	ruleBody := NewPolicyRuleBody(condition, ctx)
	fmt.Printf("the rule body is %+v\n", ruleBody)
	result, err := ruleBody.IfBody.Rego(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("the result is %s", result)
	fileName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)) + ".rego"
	switch reflect.TypeOf(ruleBody.IfBody) {
	case reflect.TypeOf(AllOf{}):
		conditionName = ruleBody.IfBody.(AllOf).ConditionSetName
	case reflect.TypeOf(AnyOf{}):
		conditionName = ruleBody.IfBody.(AnyOf).ConditionSetName
	case reflect.TypeOf(NotOperator{}):
		conditionName = ruleBody.IfBody.(NotOperator).ConditionSetName
	default:
		conditionName = result
	}
	if action == disabled {
		result = "default allow := true\n\n" + result
	} else if action == deny {
		top := "deny if {\n" + " " + conditionName + "\n}\n"
		result = top + result
	} else if action == warn {
		top := "warn if {\n" + " " + conditionName + "\n}\n"
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

func azPolicy2Rego(path string, ctx context.Context) error {
	var action string
	fmt.Printf("the path is %+v\n", path)

	rule, err := ruleIterator(path)
	if err != nil {
		fmt.Printf("cannot find rules %+v\n", err)
		return err
	}

	then := rule.Properties.PolicyRule.GetThen()
	action, err = then.MapEffectToAction(rule.Properties.Parameters.GetEffect().GetDefaultValue())
	if err != nil {
		return err
	}
	fmt.Printf("the effect is %+v\n", action)
	condition, err := rule.Properties.PolicyRule.GetIf().condition(ctx)
	if err != nil {
		fmt.Printf("cannot find conditions %+v\n", err)
		return err
	}
	rt, err := currentResourceType(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("the resource type is %+v\n", rt)
	fmt.Printf("the whole condition is %+v\n", *condition)
	fileName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)) + ".rego"
	conditionNames, result, err := condition.RuleSetReader("", ctx)
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

func currentResourceType(ctx context.Context) (string, error) {
	resourceTypeStack := ctx.Value("context").(map[string]stacks.Stack)["resourceType"]
	if resourceTypeStack == nil {
		return "", fmt.Errorf("cannot find the resource type in the context")
	}
	resourceType, ok := resourceTypeStack.Peek()
	if !ok {
		return "", fmt.Errorf("cannot find the resource type in the context")
	}
	rt, ok := resourceType.(string)
	if !ok {
		return "", fmt.Errorf("cannot convert the resource type to string")
	}
	return rt, nil
}

func pushResourceType(ctx context.Context, rt string) {
	contextMap := ctx.Value("context").(map[string]stacks.Stack)
	contextMap["resourceType"].Push(rt)
}

//func popResourceType(ctx context.Context) {
//	resourceTypeStack := ctx.Value("resourceType").(stacks.Stack)
//	resourceTypeStack.Pop()
//}

func jsonFiles(dir string) ([]string, error) {
	res, err := readJsonFilePaths(dir)
	if err != nil {
		fmt.Printf("cannot find files in directory %+v\n", err)
		return nil, err
	}
	return res, nil
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
			continue
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
	m := make(map[string]any)
	json.Unmarshal(data, &m)

	return &rule, nil
}

// In this function, the input is a single condition or a set of conditions(allof, anyof)
func conditionFinder(conditions map[string]interface{}, ctx context.Context) (*RuleSet, error) {
	if conditions == nil {
		return nil, errors.New("cannot find conditions")
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
				rule, err := conditionFinder(condition.(map[string]interface{}), ctx)
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
				rule, err := conditionFinder(condition.(map[string]interface{}), ctx)
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
			rule, err := conditionFinder(notCondition, ctx)
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
			rule, err := conditionFinder(whereConditions, ctx)
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
			rule, err := conditionFinder(countConditions, ctx)
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
			rt := operatorValue.(string)
			pushResourceType(ctx, rt)
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
