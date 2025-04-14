package operation

import (
	"fmt"
	"json-rule-finder/pkg/condition"
	"json-rule-finder/pkg/shared"
)

var _ Operation = &AnyOf{}

type AnyOf struct {
	baseOperation
	Conditions []shared.Rego
}

func NewAnyOf(conditionSetName string, conditions []shared.Rego) AnyOf {
	return AnyOf{
		baseOperation: baseOperation{
			helperFunctionName: conditionSetName,
		},
		Conditions: conditions,
	}
}

func ParseAnyOf(input any, ctx *shared.Context) shared.Rego {
	body, conditionSetName, err := parseOperationBody(input, ctx)
	if err != nil {
		panic(err)
	}
	return AnyOf{
		Conditions:    body,
		baseOperation: baseOperation{helperFunctionName: conditionSetName},
	}
}

func (a AnyOf) Rego(ctx *shared.Context) (string, error) {
	var res string
	head := a.HelperFunctionName()
	if _, ok := ctx.FieldNameReplacer(); ok {
		head = a.HelperFunctionName() + "(x)"
	}
	for _, item := range a.Conditions {
		if res != "" {
			res = res + "\n"
		}
		if _, ok := item.(Operation); ok {
			if _, ok := ctx.FieldNameReplacer(); ok {
				res += head + " if {" + item.(Operation).HelperFunctionName() + "(x)}"
			} else {
				res += head + " if {" + item.(Operation).HelperFunctionName() + "}"
			}
			subSet, err := item.Rego(ctx)
			if err != nil {
				return "", err
			}
			ctx.EnqueueHelperFunction(subSet)
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
	return res, nil
}
