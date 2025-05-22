package condition

import (
	"github.com/tfmodtest/azopaform/pkg/shared"
)

var _ shared.Rego = LiteralValue{}

type LiteralValue struct {
	Value        string
	ConditionSet shared.Rego
}

func NewLiteralValue(input any, ctx *shared.Context) (shared.Rego, error) {
	v, err := shared.ResolveParameterValueAsString(input, ctx)
	if err != nil {
		return nil, err
	}
	return LiteralValue{
		Value: v,
	}, nil
}

func (v LiteralValue) Rego(ctx *shared.Context) (string, error) {
	return v.Value, nil
}
