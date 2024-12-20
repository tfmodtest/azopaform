package pkg

import (
	"context"
	"github.com/emirpasic/gods/stacks"
	"reflect"
)

var _ Rego = &AllOf{}

var _ Operator = &AllOf{}

type AllOf struct {
	Conditions       []Rego
	ConditionSetName string
}

func (a AllOf) GetConditionSetName() string {
	return a.ConditionSetName
}

func (a AllOf) Rego(ctx context.Context) (string, error) {
	var res string
	var subSets []string

	if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
		res = a.ConditionSetName + "(x)" + " " + ifCondition + " {"
	} else {
		res = a.ConditionSetName + " " + ifCondition + " {"
	}

	for _, item := range a.Conditions {
		if _, ok := item.(Operator); ok {
			if reflect.TypeOf(item) != reflect.TypeOf(WhereOperator{}) {
				if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
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

	// add condition set body at the end
	return res, nil
}
