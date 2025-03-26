package pkg

import (
	"json-rule-finder/pkg/shared"
)

var _ Operation = &AllOf{}

type AllOf struct {
	baseOperator
	Conditions []shared.Rego
}

func NewAllOf(input any, ctx *shared.Context) shared.Rego {
	body, base, err := parseOperationBody(input, ctx)
	if err != nil {
		panic(err)
	}
	return AllOf{
		baseOperator: base,
		Conditions:   body,
	}
}

func (a AllOf) Rego(ctx *shared.Context) (string, error) {
	var res string
	var subSets []string

	res = a.GetConditionSetName() + " " + shared.IfCondition + " {"
	if _, ok := ctx.FieldNameReplacer(); ok {
		res = a.GetConditionSetName() + "(x)" + " " + shared.IfCondition + " {"
	}

	for _, item := range a.Conditions {
		if _, ok := item.(Operation); ok {
			if _, ok := ctx.FieldNameReplacer(); ok {
				res += "\n" + item.(Operation).GetConditionSetName() + "(x)"
			} else {
				res += "\n" + item.(Operation).GetConditionSetName()
			}
			subSet, err := item.Rego(ctx)
			if err != nil {
				return "", err
			}
			subSets = append(subSets, subSet)
			continue
		}

		condition, err := item.Rego(ctx)
		if err != nil {
			return "", err
		}
		res += "\n" + condition
	}

	res = res + "\n}"

	for _, subSet := range subSets {
		res += "\n" + subSet
	}

	// add BaseCondition set body at the end
	return res, nil
}
