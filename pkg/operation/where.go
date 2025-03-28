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
	conditionSetName, err := NeoConditionNameGenerator(ctx)
	if err != nil {
		return nil
	}
	return Where{
		Condition:        whereBody,
		ConditionSetName: conditionSetName,
	}
}

func (w Where) GetConditionSetName() string {
	return w.ConditionSetName
}

func (w Where) Rego(ctx *shared.Context) (string, error) {
	var res string
	var subSets []string
	item := w.Condition
	if operation, ok := item.(Operation); ok {
		res += operation.GetConditionSetName() + "(x)"
		ctx.PushFieldName("x")
		subSet, err := item.Rego(ctx)
		if err != nil {
			return "", err
		}
		subSets = append(subSets, subSet)
	} else {
		ctx.PushFieldName("x")
		condition, err := item.Rego(ctx)
		if err != nil {
			return "", err
		}
		res += condition
	}

	res = w.ConditionSetName + "(x)" + " " + shared.IfCondition + " {\n" + res
	res = res + "\n" + "}"

	for _, subSet := range subSets {
		res += "\n" + subSet
	}

	// add BaseCondition set body at the end
	return res, nil
}
