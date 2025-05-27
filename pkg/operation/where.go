package operation

import (
	"github.com/tfmodtest/azopaform/pkg/shared"
)

var _ Operation = &Where{}

type Where struct {
	Condition shared.Rego
	//helperFunctionName string
	//rego               string
	helperFunctionName string
}

func NewWhere(input any, ctx *shared.Context) (Operation, error) {
	whereBody, err := NewOperationOrCondition(input.(map[string]any), ctx)
	if err != nil {
		return nil, err
	}
	return Where{
		helperFunctionName: RandomHelperFunctionNameGenerator(),
		Condition:          whereBody,
	}, nil
	//rego, err := whereBody.Rego(ctx)
	//if err != nil {
	//	return nil, err
	//}
	//return Where{
	//	rego:               rego,
	//}, nil
}

func (w Where) HelperFunctionName() string {
	return w.helperFunctionName
}

func (w Where) Rego(ctx *shared.Context) (string, error) {
	var res string
	item := w.Condition
	if operation, ok := item.(Operation); ok {
		res += operation.HelperFunctionName() + "(r, x)"
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

	res = w.helperFunctionName + "(r, x)" + " " + shared.IfCondition + " {\n" + res
	res = res + "\n" + "}"

	// add BaseCondition set body at the end
	return res, nil
}
