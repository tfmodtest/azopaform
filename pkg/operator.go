package pkg

import (
	"fmt"
	"json-rule-finder/pkg/condition"
	"json-rule-finder/pkg/shared"
	"strings"
)

type Operator interface {
	shared.Rego
	GetConditionSetName() string
}

var otherFactories = make(map[string]func(input any, ctx *shared.Context) shared.Rego)
var operators map[string]func(input any, ctx *shared.Context) shared.Rego

func init() {
	operators = map[string]func(input any, ctx *shared.Context) shared.Rego{
		shared.AllOf: NewAllOf,
		shared.AnyOf: NewAnyOf,
		shared.Not:   NewNot,
	}

	otherFactories[shared.Count] = func(input any, ctx *shared.Context) shared.Rego {
		items := input.(map[string]any)
		var whereBody shared.Rego
		if items[shared.Where] != nil {
			whereMap := items[shared.Where].(map[string]any)
			of := otherFactories[shared.Where]
			whereBody = of(whereMap, ctx)
		}
		fieldName := items[shared.Field]
		if items[shared.Field] == nil {
			fieldName = items[shared.Value]
		}
		countField, _, err := shared.FieldNameProcessor(fieldName.(string), ctx)
		if err != nil {
			countField = items[shared.Field].(string)
			fmt.Printf("error in field name processor: %v\n", err)
		}
		fmt.Printf("count field is %v\n", countField)
		countFieldConverted := replaceIndex(countField)
		var countBody string
		if whereBody != nil {
			countBody = shared.Count + "(" + "{" + "x" + "|" + countFieldConverted + ";" + whereBody.(WhereOperator).ConditionSetName + "(x)" + "}" + ")"
		} else {
			countBody = shared.Count + "(" + "{" + "x" + "|" + countFieldConverted + "}" + ")"
		}
		countBody = strings.Replace(countBody, "*", "x", -1)
		return CountOperator{
			Where:    whereBody,
			CountExp: countBody,
		}
	}
	otherFactories[shared.Where] = func(input any, ctx *shared.Context) shared.Rego {
		itemMap := input.(map[string]any)
		var body []shared.Rego
		var cf func(shared.Rego, any) shared.Rego
		var conditionKey string
		var subjectKey string
		var of func(any, *shared.Context) shared.Rego
		var operatorValue any
		for k, _ := range itemMap {
			if f, ok := condition.ConditionFactory[k]; ok {
				cf = f
				conditionKey = k
				continue
			}
		}
		for k, v := range itemMap {
			k = strings.ToLower(k)
			if f, ok := operators[k]; ok {
				of = f
				operatorValue = v
				continue
			}
			if f, ok := otherFactories[k]; ok {
				of = f
				operatorValue = v
				continue
			}
		}
		if cf != nil {
			for k, _ := range itemMap {
				if k == conditionKey {
					continue
				}
				subjectKey = k
			}
			subject := subjectFactories[subjectKey](itemMap[subjectKey], ctx)
			body = append(body, cf(subject, itemMap[conditionKey]))
		} else if of != nil {
			body = append(body, of(operatorValue, ctx))
			//fmt.Printf("where body is %v\n", body)
		}
		//conditionSetName := conditionNameGenerator(whereConditionLen, charNum)
		conditionSetName, err := NeoConditionNameGenerator(ctx)
		if err != nil {
			return nil
		}
		return WhereOperator{
			Conditions:       body,
			ConditionSetName: conditionSetName,
		}
	}
}

var _ shared.Rego = &CountOperator{}

type CountOperator struct {
	Where    shared.Rego
	CountExp string
}

func (c CountOperator) Rego(ctx *shared.Context) (string, error) {
	var res string
	whereSubset, err := c.Where.Rego(ctx)
	if err != nil {
		return "", err
	}
	if err != nil {
		return "", err
	}
	res = c.CountExp + "\n" + whereSubset
	return res, nil
}

var _ shared.Rego = &WhereOperator{}

var _ Operator = &WhereOperator{}

type WhereOperator struct {
	Conditions       []shared.Rego
	ConditionSetName string
}

func (w WhereOperator) GetConditionSetName() string {
	return w.ConditionSetName
}

func (w WhereOperator) GetConditionSetNameRev() string {
	return strings.Join([]string{"not", w.ConditionSetName}, " ")
}

func (w WhereOperator) Rego(ctx *shared.Context) (string, error) {
	var res string
	var subSets []string

	for _, item := range w.Conditions {
		if _, ok := item.(Operator); ok {
			res += item.(Operator).GetConditionSetName() + "(x)"
			ctx.PushFieldName("x")
			subSet, err := item.Rego(ctx)
			if err != nil {
				return "", err
			}
			subSets = append(subSets, subSet)
		} else {
			ctx.PushFieldName("x")
			condition, err := item.Rego(ctx)
			if err != nil {
				return "", err
			}
			res += condition
		}
	}

	res = w.ConditionSetName + "(x)" + " " + shared.IfCondition + " {\n" + res
	res = res + "\n" + "}"

	for _, subSet := range subSets {
		res += "\n" + subSet
	}

	// add BaseCondition set body at the end
	return res, nil
}

var NeoConditionNameGenerator = func(ctx *shared.Context) (string, error) {
	index, ok := ctx.PopConditionNameCounter()
	if !ok {
		return "", fmt.Errorf("conditionNameStack is empty")
	}
	conditionName := "condition" + fmt.Sprintf("%d", index)
	return conditionName, nil
}
