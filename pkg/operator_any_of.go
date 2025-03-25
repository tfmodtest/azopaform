package pkg

import (
	"fmt"
	"json-rule-finder/pkg/condition"
	"json-rule-finder/pkg/shared"
	"reflect"
	"strings"
)

var _ Operator = &AnyOf{}

type AnyOf struct {
	baseOperator
	Conditions []shared.Rego
}

func NewAnyOf(input any, ctx *shared.Context) shared.Rego {
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
				ctx.PushResourceType(v.(string))
			}
		}
		if cf != nil {
			for k, _ := range itemMap {
				if k == conditionKey {
					continue
				}
				subjectKey = k
			}
			subjectItem := itemMap[subjectKey]
			if subjectKey == shared.Field && subjectItem == shared.TypeOfResource {
				rawType := itemMap[conditionKey]
				translatedType, err := ResourceTypeParser(rawType.(string))
				if err != nil {
					return nil
				}
				body = append(body, cf(NewSubject(subjectKey, subjectItem, ctx), translatedType))
				continue
			}
			subject := NewSubject(subjectKey, subjectItem, ctx)
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
	//conditionName := conditionNameGenerator(orConditionLen, charNum)
	conditionName, err := NeoConditionNameGenerator(ctx)
	if err != nil {
		return nil
	}
	return AnyOf{
		Conditions: body,
		baseOperator: baseOperator{
			conditionSetName: conditionName,
		},
	}
}

func (a AnyOf) Rego(ctx *shared.Context) (string, error) {
	var res string
	var subSets []string
	head := a.GetConditionSetName()
	if _, ok := ctx.FieldNameReplacer(); ok {
		head = a.GetConditionSetName() + "(x)"
	}
	for _, item := range a.Conditions {
		if res != "" {
			res = res + "\n"
		}
		if _, ok := item.(Operator); ok {
			if reflect.TypeOf(item) != reflect.TypeOf(WhereOperator{}) {
				if _, ok := ctx.FieldNameReplacer(); ok {
					res += head + " if {" + item.(Operator).GetConditionSetName() + "(x)}"
				} else {
					res += head + " if {" + item.(Operator).GetConditionSetName() + "}"
				}
			}
			subSet, err := item.Rego(ctx)
			if err != nil {
				return "", err
			}
			subSets = append(subSets, subSet)
			continue
		}

		if _, ok := item.(condition.Condition); ok {
			condition, err := item.(condition.Condition).Rego(ctx)
			if err != nil {
				return "", err
			}
			res += fmt.Sprintf("%s if {%s}", head, condition)
		}
	}

	for _, subSet := range subSets {
		res += "\n" + subSet
	}
	return res, nil
}
