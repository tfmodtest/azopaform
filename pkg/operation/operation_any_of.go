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
			conditionSetName: conditionSetName,
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
		baseOperation: baseOperation{conditionSetName: conditionSetName},
	}
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
		if _, ok := item.(Operation); ok {
			if _, ok := ctx.FieldNameReplacer(); ok {
				res += head + " if {" + item.(Operation).GetConditionSetName() + "(x)}"
			} else {
				res += head + " if {" + item.(Operation).GetConditionSetName() + "}"
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
