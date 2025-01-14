package pkg

import (
	"context"
	"json-rule-finder/pkg/shared"
)

type OperationValue string
type OperationField string

func (o OperationValue) Rego(ctx context.Context) (string, error) {
	processed, _, err := shared.FieldNameProcessor(string(o), ctx)
	return processed, err
}

func (o OperationField) Rego(ctx context.Context) (string, error) {
	processed, _, err := shared.FieldNameProcessor(string(o), ctx)
	return processed, err
}
