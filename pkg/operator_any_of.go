package pkg

import (
	"context"
	"fmt"
	"json-rule-finder/pkg/condition"
	"json-rule-finder/pkg/shared"
	"reflect"

	"github.com/emirpasic/gods/stacks"
)

var _ Operator = &AnyOf{}

type AnyOf struct {
	Conditions       []shared.Rego
	ConditionSetName string
}

func (a AnyOf) GetConditionSetName() string {
	return a.ConditionSetName
}

func (a AnyOf) Rego(ctx context.Context) (string, error) {
	var res string
	var subSets []string
	head := a.ConditionSetName
	if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
		head = a.ConditionSetName + "(x)"
	}
	for _, item := range a.Conditions {
		if res != "" {
			res = res + "\n"
		}
		if _, ok := item.(Operator); ok {
			if reflect.TypeOf(item) != reflect.TypeOf(WhereOperator{}) {
				if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
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
