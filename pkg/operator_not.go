package pkg

import (
	"context"
	"fmt"
	"github.com/emirpasic/gods/stacks"
	"reflect"
)

var _ Rego = &NotOperator{}

var _ Operator = &NotOperator{}

type NotOperator struct {
	Body             Rego
	ConditionSetName string
}

func (n NotOperator) GetConditionSetName() string {
	return n.ConditionSetName
}

func (n NotOperator) Rego(ctx context.Context) (string, error) {
	body, ok := n.Body.(Operator)
	if !ok {
		body = &AllOf{
			Conditions: []Rego{
				n.Body,
			},
			ConditionSetName: fmt.Sprintf("%s_%s", n.ConditionSetName, "negation"),
		}
	}
	bodyRes, subSet, err := n.encapsulateHelperOperator(body, ctx)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`%s if {
  %s
}

%s`, n.ConditionSetName, bodyRes, subSet), nil
}

func (n NotOperator) encapsulateHelperOperator(body Operator, ctx context.Context) (string, string, error) {
	var res string
	if reflect.TypeOf(body) != reflect.TypeOf(WhereOperator{}) {
		if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
			res = "not " + body.GetConditionSetName() + "(x)"
		} else {
			res = "not " + body.GetConditionSetName()
		}
	}

	subSet, err := body.Rego(ctx)
	if err != nil {
		return "", "", err
	}
	return res, subSet, nil
}
