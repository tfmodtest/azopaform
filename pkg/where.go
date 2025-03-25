package pkg

import (
	"json-rule-finder/pkg/condition"
	"json-rule-finder/pkg/shared"
	"strings"
)

var _ shared.Rego = &WhereOperator{}

var _ Operator = &WhereOperator{}

type WhereOperator struct {
	Conditions       []shared.Rego
	ConditionSetName string
}

func NewWhere(input any, ctx *shared.Context) shared.Rego {
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
		if f, ok := valueFactories[k]; ok {
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
		subject := NewSubject(subjectKey, itemMap[subjectKey], ctx)
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
