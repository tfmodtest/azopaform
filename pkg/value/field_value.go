package value

import (
	"json-rule-finder/pkg/shared"
	"strings"
)

var _ shared.Rego = FieldValue{}

type FieldValue struct {
	Name string
}

func NewFieldValue(input any, ctx *shared.Context) shared.Rego {
	name := input.(string)
	name = strings.ReplaceAll(name, "[*]", "[x]")
	return FieldValue{
		Name: name,
	}
}

func (f FieldValue) Rego(ctx *shared.Context) (string, error) {
	processed, err := shared.FieldNameProcessor(f.Name, ctx)
	return processed, err
}
