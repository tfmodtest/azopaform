package pkg

import (
	"context"
	"json-rule-finder/pkg/shared"
	"strings"
)

var subjectFactories = map[string]func(input any, ctx context.Context) shared.Rego{
	"field": func(input any, ctx context.Context) shared.Rego {
		name := input.(string)
		name = strings.ReplaceAll(name, "[*]", "[x]")
		return FieldValue{
			Name: name,
		}
	},
	"value": func(input any, ctx context.Context) shared.Rego {
		value := input.(string)
		value = strings.ReplaceAll(value, "[*]", "[x]")
		return Value{
			Value: value,
		}
	},
	"count": func(input any, ctx context.Context) shared.Rego {
		f := operatorFactories[shared.Count_]
		countConditionSet := f(input, ctx)
		//fmt.Printf("countConditionSet: %v\n", countConditionSet)
		return Count{
			Count:        countConditionSet.(CountOperator).CountExp,
			ConditionSet: countConditionSet.(CountOperator).Where,
		}
	},
}

var _ shared.Rego = &FieldValue{}

type FieldValue struct {
	Name string
}

func (f FieldValue) Rego(ctx context.Context) (string, error) {
	processed, _, err := shared.FieldNameProcessor(f.Name, ctx)
	return processed, err
}

var _ shared.Rego = &Value{}

type Value struct {
	Value        string
	ConditionSet shared.Rego
}

func (v Value) Rego(ctx context.Context) (string, error) {
	return v.Value, nil
}

var _ shared.Rego = &Count{}

type Count struct {
	Count        string
	ConditionSet shared.Rego
}

func (c Count) Rego(ctx context.Context) (string, error) {
	return c.Count, nil
}
