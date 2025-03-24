package pkg

import (
	"fmt"
	"json-rule-finder/pkg/condition"
	"json-rule-finder/pkg/shared"
	"reflect"
	"strings"
)

var _ Operator = &AllOf{}

type AllOf struct {
	baseOperator
	Conditions []shared.Rego
}

func NewAllOf(input any, ctx *shared.Context) shared.Rego {
	items := input.([]any)
	var body []shared.Rego
	for _, item := range items {
		itemMap := item.(map[string]any)
		var cf func(shared.Rego, any) shared.Rego
		var conditionKey string
		var subjectKey string
		var of func(any, *shared.Context) shared.Rego
		var operatorValue any
		var containsTypeOfResource bool
		for k, v := range itemMap {
			if k == shared.Field && v == shared.TypeOfResource {
				containsTypeOfResource = true
			}
			if f, ok := condition.ConditionFactory[strings.ToLower(k)]; ok {
				cf = f
				conditionKey = k
				continue
			}
		}
		for k, v := range itemMap {
			if f, ok := operators[strings.ToLower(k)]; ok {
				of = f
				operatorValue = v
				continue
			}
		}
		if containsTypeOfResource {
			for k, v := range itemMap {
				if k == shared.Field && v == shared.TypeOfResource {
					continue
				}
				if reflect.TypeOf(v).Kind() == reflect.String {
					ctx.PushResourceType(v.(string))
				}
			}
		}
		if cf != nil {
			for k, _ := range itemMap {
				if k == conditionKey {
					continue
				}
				subjectKey = k
				fmt.Printf("subject key is %v\n", subjectKey)
			}
			subjectItem := itemMap[subjectKey]
			if subjectKey == shared.Field && subjectItem == shared.TypeOfResource {
				if reflect.TypeOf(itemMap[conditionKey]).Kind() == reflect.String {
					rawType := itemMap[conditionKey]
					translatedType, err := ResourceTypeParser(rawType.(string))
					if err != nil {
						fmt.Printf("error in resource type parser: %v\n", err)
						return nil
					}
					body = append(body, cf(subjectFactories[subjectKey](subjectItem, ctx), translatedType))
					continue
				} else {
					rawTypes := itemMap[conditionKey]
					var translatedTypes []string
					for _, rawType := range rawTypes.([]interface{}) {
						translatedType, err := ResourceTypeParser(rawType.(string))
						if err != nil {
							fmt.Printf("error in resource type parser: %v\n", err)
							return nil
						}
						translatedTypes = append(translatedTypes, translatedType)
					}
					body = append(body, cf(subjectFactories[subjectKey](subjectItem, ctx), translatedTypes))
				}
			}
			subject := subjectFactories[subjectKey](subjectItem, ctx)
			if reflect.TypeOf(subject) == reflect.TypeOf(Count{}) {
				body = append(body, cf(subject, itemMap[conditionKey]))
				body = append(body, subject.(Count).ConditionSet)
			} else {
				body = append(body, cf(subject, itemMap[conditionKey]))
			}
		} else if of != nil {
			body = append(body, of(operatorValue, ctx))
		}
	}
	//conditionSetName := conditionNameGenerator(andConditionLen, charNum)
	conditionSetName, err := NeoConditionNameGenerator(ctx)
	if err != nil {
		fmt.Printf("error in BaseCondition name generator: %v\n", err)
		return nil
	}
	return AllOf{
		baseOperator: baseOperator{
			conditionSetName: conditionSetName,
		},
		Conditions: body,
	}
}

func (a AllOf) Rego(ctx *shared.Context) (string, error) {
	var res string
	var subSets []string

	res = a.GetConditionSetName() + " " + shared.IfCondition + " {"
	if _, ok := ctx.FieldNameReplacer(); ok {
		res = a.GetConditionSetName() + "(x)" + " " + shared.IfCondition + " {"
	}

	for _, item := range a.Conditions {
		if _, ok := item.(Operator); ok {
			if reflect.TypeOf(item) != reflect.TypeOf(WhereOperator{}) {
				if _, ok := ctx.FieldNameReplacer(); ok {
					res += "\n" + item.(Operator).GetConditionSetName() + "(x)"
				} else {
					res += "\n" + item.(Operator).GetConditionSetName()
				}
			}
			subSet, err := item.Rego(ctx)
			if err != nil {
				return "", err
			}
			subSets = append(subSets, subSet)
			continue
		}

		condition, err := item.Rego(ctx)
		if err != nil {
			return "", err
		}
		res += "\n" + condition
	}

	res = res + "\n}"

	for _, subSet := range subSets {
		res += "\n" + subSet
	}

	// add BaseCondition set body at the end
	return res, nil
}
