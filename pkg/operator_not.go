package pkg

import (
	"fmt"
	"json-rule-finder/pkg/condition"
	"json-rule-finder/pkg/shared"
	"reflect"
	"strings"
)

var _ Operator = &NotOperator{}

type NotOperator struct {
	baseOperator
	Body shared.Rego
}

func NewNot(input any, ctx *shared.Context) shared.Rego {
	itemMap := input.(map[string]any)
	fmt.Printf("item map is %v\n", itemMap)
	var cf func(shared.Rego, any) shared.Rego
	var conditionKey string
	var subjectKey string
	var of func(any, *shared.Context) shared.Rego
	var subject shared.Rego
	var body shared.Rego
	var operatorValue any
	for k, _ := range itemMap {
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
	if cf != nil {
		for k, _ := range itemMap {
			if k == conditionKey {
				continue
			}
			subjectKey = k
		}
		fmt.Printf("subject key is %v\n", subjectKey)
		subject = subjectFactories[subjectKey](itemMap[subjectKey], ctx)
		body = cf(subject, itemMap[conditionKey])
	} else if of != nil {
		body = of(operatorValue, ctx)
	}
	conditionName, err := NeoConditionNameGenerator(ctx)
	if err != nil {
		return nil
	}
	return NotOperator{
		Body: body,
		baseOperator: baseOperator{
			conditionSetName: conditionName,
		},
	}
}

func (n NotOperator) Rego(ctx *shared.Context) (string, error) {
	body, ok := n.Body.(Operator)
	if !ok {
		body = &AllOf{
			Conditions: []shared.Rego{
				n.Body,
			},
			baseOperator: baseOperator{
				conditionSetName: fmt.Sprintf("%s_%s", n.GetConditionSetName(), "negation"),
			},
		}
	}
	var bodyRes string
	if reflect.TypeOf(body) != reflect.TypeOf(WhereOperator{}) {
		if _, ok := ctx.FieldNameReplacer(); ok {
			bodyRes = body.GetConditionSetName() + "(x)"
		} else {
			bodyRes = body.GetConditionSetName()
		}
	}
	subSet, err := body.Rego(ctx)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`%s if {
  not %s
}

%s`, n.GetConditionSetName(), bodyRes, subSet), nil
}
