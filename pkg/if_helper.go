package pkg

import (
	"github.com/tfmodtest/azopaform/pkg/operation"
	"github.com/tfmodtest/azopaform/pkg/shared"
)

type IfBody map[string]any

var _ shared.Rego = &If{}

type If struct {
	rego shared.Rego
}

func NewIf(body IfBody, ctx *shared.Context) (*If, error) {
	i := &If{}
	var err error
	if i.rego, err = operation.NewOperationOrCondition(body, ctx); err != nil {
		return nil, err
	}
	return i, nil
}

func (i *If) Rego(ctx *shared.Context) (string, error) {
	return i.rego.Rego(ctx)
}

func (i *If) ConditionName(defaultConditionName string) string {
	if operator, ok := i.rego.(operation.Operation); ok {
		return operator.HelperFunctionName()
	}
	return defaultConditionName
}
