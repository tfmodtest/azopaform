package pkg

import "context"

var subjectFactories = map[string]func(input any) Rego{
	"field": func(input any) Rego {
		return &FieldValue{
			Name: input.(string),
		}
	},
}

var _ Rego = &FieldValue{}

type FieldValue struct {
	Name string
}

func (f FieldValue) Rego(ctx context.Context) (string, error) {
	//TODO implement me
	panic("implement me")
}
