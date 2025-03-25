package pkg

import (
	"json-rule-finder/pkg/shared"
	"strings"
)

func NewSubject(subjectKey string, body any, ctx *shared.Context) shared.Rego {
	return subjectFactories[subjectKey](body, ctx)
}

var subjectFactories map[string]func(input any, ctx *shared.Context) shared.Rego

func init() {
	subjectFactories = map[string]func(input any, ctx *shared.Context) shared.Rego{
		"field": func(input any, ctx *shared.Context) shared.Rego {
			name := input.(string)
			name = strings.ReplaceAll(name, "[*]", "[x]")
			return FieldValue{
				Name: name,
			}
		},
		"value": func(input any, ctx *shared.Context) shared.Rego {
			value := input.(string)
			value = strings.ReplaceAll(value, "[*]", "[x]")
			return Value{
				Value: value,
			}
		},
		"count": func(input any, ctx *shared.Context) shared.Rego {
			countConditionSet := NewCountOperator(input, ctx)
			return Count{
				Count:        countConditionSet.CountExp,
				ConditionSet: countConditionSet.Where,
			}
		},
	}
}

var _ shared.Rego = &FieldValue{}

type FieldValue struct {
	Name string
}

func (f FieldValue) Rego(ctx *shared.Context) (string, error) {
	processed, _, err := shared.FieldNameProcessor(f.Name, ctx)
	return processed, err
}

var _ shared.Rego = &Value{}

type Value struct {
	Value        string
	ConditionSet shared.Rego
}

func (v Value) Rego(ctx *shared.Context) (string, error) {
	return v.Value, nil
}
