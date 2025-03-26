package pkg

import (
	"fmt"
	"json-rule-finder/pkg/shared"
)

var _ Operation = &NotOperator{}

type NotOperator struct {
	baseOperator
	Body shared.Rego
}

func NewNot(input any, ctx *shared.Context) shared.Rego {
	itemMap := input.(map[string]any)
	body := NewOperationOrCondition(itemMap, ctx)
	conditionSetName, err := NeoConditionNameGenerator(ctx)
	if err != nil {
		panic(err)
	}
	return NotOperator{
		Body: body,
		baseOperator: baseOperator{
			conditionSetName: conditionSetName,
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
