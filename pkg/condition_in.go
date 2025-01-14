package pkg

import (
	"context"
	"json-rule-finder/pkg/shared"
	"strings"
)

var _ Condition = InCondition{}

type InCondition struct {
	BaseCondition
	Values []string
}

func (i InCondition) Rego(ctx context.Context) (string, error) {
	fieldName, err := i.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{"some", fieldName, "in", shared.SliceConstructor(i.Values)}, " "), nil
}
