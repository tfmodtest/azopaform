package pkg

import "context"

var subjectFactories = map[string]func(input any) Rego{
	"field": func(input any) Rego {
		return FieldValue{
			Name: input.(string),
		}
	},
	"value": func(input any) Rego {
		return Value{
			Value: input.(string),
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
	Value string
}

func (v Value) Rego(ctx context.Context) (string, error) {
	return v.Value, nil
}
