package pkg

import (
	"fmt"
	"json-rule-finder/pkg/condition"
	"json-rule-finder/pkg/shared"
	"reflect"
	"strings"
)

var _ Operation = &AllOf{}

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
		for k, _ := range itemMap {
			if f, ok := condition.ConditionFactory[strings.ToLower(k)]; ok {
				cf = f
				conditionKey = k
				continue
			}
		}
		if v, ok := itemMap[shared.Field]; ok && v == shared.TypeOfResource {
			resourceType, ok := itemMap["equals"].(string)
			if !ok {
				panic("resource type without value")
			}
			ctx.PushResourceType(resourceType)
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
					body = append(body, cf(NewSubject(subjectKey, subjectItem, ctx), rawType))
					continue
				} else {
					rawTypes := itemMap[conditionKey]
					var translatedTypes []string
					for _, rawType := range rawTypes.([]interface{}) {
						translatedTypes = append(translatedTypes, rawType.(string))
					}
					body = append(body, cf(NewSubject(subjectKey, subjectItem, ctx), translatedTypes))
				}
			}
			subject := NewSubject(subjectKey, subjectItem, ctx)
			if reflect.TypeOf(subject) == reflect.TypeOf(Count{}) {
				body = append(body, cf(subject, itemMap[conditionKey]))
				body = append(body, subject.(Count).ConditionSet)
			} else {
				body = append(body, cf(subject, itemMap[conditionKey]))
			}
		}
		for k, v := range itemMap {
			if operation := NewOperation(strings.ToLower(k), v, ctx); operation != nil {
				body = append(body, operation)
			}
			break
		}
	}
	conditionSetName, err := NeoConditionNameGenerator(ctx)
	if err != nil {
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
		if _, ok := item.(Operation); ok {
			if _, ok := ctx.FieldNameReplacer(); ok {
				res += "\n" + item.(Operation).GetConditionSetName() + "(x)"
			} else {
				res += "\n" + item.(Operation).GetConditionSetName()
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
