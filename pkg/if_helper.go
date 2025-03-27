package pkg

import (
	"json-rule-finder/pkg/shared"
)

type IfBody map[string]any

func (i *If) Rego(ctx *shared.Context) (string, error) {
	if i.rego != nil {
		return i.rego.Rego(ctx)
	}
	i.rego = NewOperationOrCondition(i.body, ctx)
	return i.rego.Rego(ctx)
}

func (i *If) ConditionName(defaultConditionName string) string {
	if operator, ok := i.rego.(Operation); ok {
		return operator.GetConditionSetName()
	}
	return defaultConditionName
}

var _ shared.Rego = &If{}

type If struct {
	body IfBody
	rego shared.Rego
}
