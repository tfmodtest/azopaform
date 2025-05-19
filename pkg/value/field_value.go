package value

import (
	"github.com/tfmodtest/azopaform/pkg/shared"
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
	return FieldValue{
		Name: v,
	}, nil
}

func (f FieldValue) Rego(ctx *shared.Context) (string, error) {
	return shared.FieldNameProcessor(f.Name, ctx)
}
