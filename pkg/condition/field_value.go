package condition

import (
	"github.com/tfmodtest/azopaform/pkg/shared"
	"strings"
)

var _ shared.Rego = FieldValue{}

type FieldValue struct {
	Name string
}

func NewFieldValue(input any, ctx *shared.Context) (shared.Rego, error) {
	v, err := shared.ResolveParameterValueAsString(input, ctx)
	if err != nil {
		return nil, err
	}
	if v == ctx.GetCountFieldName() {
		v = strings.TrimPrefix(v, ctx.GetCountFieldName())
		v = shared.VarInCountWhere + v
	}
	return FieldValue{
		Name: v,
	}, nil
}

func (f FieldValue) Rego(ctx *shared.Context) (string, error) {
	return shared.FieldNameProcessor(f.Name, ctx)
}
