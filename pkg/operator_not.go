package pkg

import (
	"context"
	"fmt"
	"github.com/emirpasic/gods/stacks"
	"json-rule-finder/pkg/shared"
	"reflect"
)

var _ Operator = &NotOperator{}

type NotOperator struct {
	Body             shared.Rego
	ConditionSetName string
}

func (n NotOperator) GetConditionSetName() string {
	return n.ConditionSetName
}

func (n NotOperator) Rego(ctx context.Context) (string, error) {
	body, ok := n.Body.(Operator)
	if !ok {
		body = &AllOf{
			Conditions: []shared.Rego{
				n.Body,
			},
			ConditionSetName: fmt.Sprintf("%s_%s", n.ConditionSetName, "negation"),
		}
	}
	var bodyRes string
	if reflect.TypeOf(body) != reflect.TypeOf(WhereOperator{}) {
		if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
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

%s`, n.ConditionSetName, bodyRes, subSet), nil
}
