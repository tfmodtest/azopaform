package value

import (
	"github.com/tfmodtest/azopaform/pkg/shared"
)

var _ shared.Rego = FieldValue{}

type FieldValue struct {
	Name string
}

func NewFieldValue(input any, ctx *shared.Context) shared.Rego {
	return FieldValue{
		Name: input.(string),
	}
}

func (f FieldValue) Rego(ctx *shared.Context) (string, error) {
	return shared.FieldNameProcessor(f.Name, ctx)
}
