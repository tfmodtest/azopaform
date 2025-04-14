package operation

import (
	"json-rule-finder/pkg/shared"
)

var _ Operation = &Where{}

type Where struct {
	Condition        shared.Rego
	ConditionSetName string
}

func NewWhere(input any, ctx *shared.Context) Operation {
	whereBody := NewOperationOrCondition(input.(map[string]any), ctx)
	conditionSetName := NeoConditionNameGenerator(ctx)
	return Where{
		Condition:        whereBody,
		ConditionSetName: conditionSetName,
	}
}

func (w Where) HelperFunctionName() string {
	return w.ConditionSetName
}

func (w Where) Rego(ctx *shared.Context) (string, error) {
	var res string
	item := w.Condition
	if operation, ok := item.(Operation); ok {
		res += operation.HelperFunctionName() + "(x)"
		if err := ctx.InHelperFunction("x", func() error {
			helperFunctionBody, err := item.Rego(ctx)
			if err != nil {
				return err
			}
			ctx.EnqueueHelperFunction(helperFunctionBody)
			return nil
		}); err != nil {
			return "", err
		}
	} else {
		if err := ctx.InHelperFunction("x", func() error {
			condition, err := item.Rego(ctx)
			if err != nil {
				return err
			}
			res += condition
			return nil
		}); err != nil {
			return "", err
		}
	}

	res = w.ConditionSetName + "(x)" + " " + shared.IfCondition + " {\n" + res
	res = res + "\n" + "}"

	// add BaseCondition set body at the end
	return res, nil
}
