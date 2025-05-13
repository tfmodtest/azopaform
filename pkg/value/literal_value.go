package value

import (
	"strings"

	"github.com/tfmodtest/azopaform/pkg/shared"
)

var _ shared.Rego = LiteralValue{}

type LiteralValue struct {
	Value        string
	ConditionSet shared.Rego
}

func NewLiteralValue(input any, ctx *shared.Context) shared.Rego {
	v := shared.ResolveParameterValue[string](input, ctx)
	v = strings.ReplaceAll(v, "[*]", "[_]")
	return LiteralValue{
		Value: v,
	}
}

func (v LiteralValue) Rego(ctx *shared.Context) (string, error) {
	return v.Value, nil
}
