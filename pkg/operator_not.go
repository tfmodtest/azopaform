package pkg

import (
	"fmt"
	"json-rule-finder/pkg/condition"
	"json-rule-finder/pkg/shared"
	"strings"
)

var _ Operation = &NotOperator{}

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
	var subject shared.Rego
	var body shared.Rego
	for k, _ := range itemMap {
		if f, ok := condition.ConditionFactory[strings.ToLower(k)]; ok {
			cf = f
			conditionKey = k
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
		subject = NewSubject(subjectKey, itemMap[subjectKey], ctx)
		body = cf(subject, itemMap[conditionKey])
	}
	for k, v := range itemMap {
		if operation := NewOperation(strings.ToLower(k), v, ctx); operation != nil {
			body = operation
			break
		}
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
	body, ok := n.Body.(Operation)
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
	if _, ok := ctx.FieldNameReplacer(); ok {
		bodyRes = body.GetConditionSetName() + "(x)"
	} else {
		bodyRes = body.GetConditionSetName()
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
