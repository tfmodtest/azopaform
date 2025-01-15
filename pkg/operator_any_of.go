package pkg

import (
	"fmt"
	"json-rule-finder/pkg/condition"
	"json-rule-finder/pkg/shared"
	"reflect"
)

var _ Operator = &AnyOf{}

type AnyOf struct {
	baseOperator
	Conditions []shared.Rego
}

func (a AnyOf) Rego(ctx *shared.Context) (string, error) {
	var res string
	var subSets []string
	head := a.GetConditionSetName()
	if _, ok := ctx.FieldNameReplacer(); ok {
		head = a.GetConditionSetName() + "(x)"
	}
	for _, item := range a.Conditions {
		if res != "" {
			res = res + "\n"
		}
		if _, ok := item.(Operator); ok {
			if reflect.TypeOf(item) != reflect.TypeOf(WhereOperator{}) {
				if _, ok := ctx.FieldNameReplacer(); ok {
					res += head + " if {" + item.(Operator).GetConditionSetName() + "(x)}"
				} else {
					res += head + " if {" + item.(Operator).GetConditionSetName() + "}"
				}
			}
			subSet, err := item.Rego(ctx)
			if err != nil {
				return "", err
			}
			subSets = append(subSets, subSet)
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

	for _, subSet := range subSets {
		res += "\n" + subSet
	}
	return res, nil
}
