package operation

import (
	"json-rule-finder/pkg/shared"
)

var _ Operation = &AllOf{}

type AllOf struct {
	baseOperation
	Conditions []shared.Rego
}

func NewAllOf(conditionSetName string, conditions []shared.Rego) AllOf {
	return AllOf{
		baseOperation: baseOperation{
			helperFunctionName: conditionSetName,
		},
		Conditions: conditions,
	}
}

func ParseAllOf(input any, ctx *shared.Context) shared.Rego {
	body, err := parseOperationBody(input, ctx)
	if err != nil {
		panic(err)
	}
	return AllOf{
		baseOperation: newBaseOperation(),
		Conditions:    body,
	}
}

func (a AllOf) Rego(ctx *shared.Context) (string, error) {
	var res string

	res = a.HelperFunctionName() + " " + shared.IfCondition + " {"
	if _, ok := ctx.FieldNameReplacer(); ok {
		res = a.HelperFunctionName() + "(x)" + " " + shared.IfCondition + " {"
	}

	for _, item := range a.Conditions {
		if _, ok := item.(Operation); ok {
			if _, ok := ctx.FieldNameReplacer(); ok {
				res += "\n" + item.(Operation).HelperFunctionName() + "(x)"
			} else {
				res += "\n" + item.(Operation).HelperFunctionName()
			}
			subFunction, err := item.Rego(ctx)
			if err != nil {
				return "", err
			}
			ctx.EnqueueHelperFunction(subFunction)
			continue
		}

		condition, err := item.Rego(ctx)
		if err != nil {
			return "", err
		}
		res += "\n" + condition
	}

	res = res + "\n}"

	// add BaseCondition set body at the end
	return res, nil
}
