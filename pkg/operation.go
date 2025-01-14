package pkg

import (
	"json-rule-finder/pkg/shared"
)

type OperationValue string
type OperationField string

func (o OperationValue) Rego(ctx *shared.Context) (string, error) {
	processed, _, err := shared.FieldNameProcessor(string(o), ctx)
	return processed, err
}

func (o OperationField) Rego(ctx *shared.Context) (string, error) {
	processed, _, err := shared.FieldNameProcessor(string(o), ctx)
	return processed, err
}
