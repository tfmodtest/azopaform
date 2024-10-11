package pkg

import (
	"context"
	"fmt"
	"strings"
)

var subjectFactories = map[string]func(input any, ctx context.Context) Rego{
	"field": func(input any, ctx context.Context) Rego {
		name := input.(string)
		name = strings.ReplaceAll(name, "[*]", "[x]")
		return FieldValue{
			Name: name,
		}
	},
	"value": func(input any, ctx context.Context) Rego {
		value := input.(string)
		value = strings.ReplaceAll(value, "[*]", "[x]")
		return Value{
			Value: value,
		}
	},
	"count": func(input any, ctx context.Context) Rego {
		f := operatorFactories[count]
		countConditionSet := f(input, ctx)
		fmt.Printf("countConditionSet: %v\n", countConditionSet)
		return Count{
			Count:        countConditionSet.(CountOperator).CountExp,
			ConditionSet: countConditionSet.(CountOperator).Where,
		}
	},
}

var _ Rego = &FieldValue{}

type FieldValue struct {
	Name string
}

func (f FieldValue) Rego(ctx context.Context) (string, error) {
	processed, _, err := FieldNameProcessor(f.Name, ctx)
	return processed, err
}

var _ Rego = &Value{}

type Value struct {
	Value        string
	ConditionSet Rego
}

func (v Value) Rego(ctx context.Context) (string, error) {
	return v.Value, nil
}

var _ Rego = &Count{}

type Count struct {
	Count        string
	ConditionSet Rego
}

func (c Count) Rego(ctx context.Context) (string, error) {
	return c.Count, nil
}
