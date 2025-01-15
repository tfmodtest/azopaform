package pkg

import (
	"json-rule-finder/pkg/shared"
	"reflect"
)

var _ Operator = &AllOf{}

type AllOf struct {
	baseOperator
	Conditions []shared.Rego
}

func (a AllOf) Rego(ctx *shared.Context) (string, error) {
	var res string
	var subSets []string

	res = a.GetConditionSetName() + " " + shared.IfCondition + " {"
	if _, ok := ctx.FieldNameReplacer(); ok {
		res = a.GetConditionSetName() + "(x)" + " " + shared.IfCondition + " {"
	}

	for _, item := range a.Conditions {
		if _, ok := item.(Operator); ok {
			if reflect.TypeOf(item) != reflect.TypeOf(WhereOperator{}) {
				if _, ok := ctx.FieldNameReplacer(); ok {
					res += "\n" + item.(Operator).GetConditionSetName() + "(x)"
				} else {
					res += "\n" + item.(Operator).GetConditionSetName()
				}
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
