package pkg

import (
	"json-rule-finder/pkg/operation"
	"json-rule-finder/pkg/shared"
)

type IfBody map[string]any

var _ shared.Rego = &If{}

type If struct {
	//body IfBody
	rego shared.Rego
}

func NewIf(body IfBody, ctx *shared.Context) *If {
	i := &If{}
	i.rego = operation.NewOperationOrCondition(body, ctx)
	return i
}

func (i *If) Rego(ctx *shared.Context) (string, error) {
	return i.rego.Rego(ctx)
}

func (i *If) ConditionName(defaultConditionName string) string {
	if operator, ok := i.rego.(operation.Operation); ok {
		return operator.GetConditionSetName()
	}
	return defaultConditionName
}
