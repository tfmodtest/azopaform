package value

import (
	"json-rule-finder/pkg/shared"
	"strings"
)

var _ shared.Rego = &FieldValue{}

type FieldValue struct {
	Name string
}

func NewFieldValue(input any, ctx *shared.Context) shared.Rego {
	name := input.(string)
	name = strings.ReplaceAll(name, "[*]", "[x]")
	return &FieldValue{
		Name: name,
	}
}

func (f FieldValue) Rego(ctx *shared.Context) (string, error) {
	processed, err := shared.FieldNameProcessor(f.Name, ctx)
	return processed, err
}

var _ shared.Rego = &Value{}

type Value struct {
	Value        string
	ConditionSet shared.Rego
}

func NewValue(input any, ctx *shared.Context) shared.Rego {
	v := input.(string)
	v = strings.ReplaceAll(v, "[*]", "[x]")
	return &Value{
		Value: v,
	}
}

func (v Value) Rego(ctx *shared.Context) (string, error) {
	return v.Value, nil
}
