package pkg

import (
	"context"
	"strings"
)

var _ Condition = InCondition{}

type InCondition struct {
	condition
	Values []string
}

func (i InCondition) Rego(ctx context.Context) (string, error) {
	fieldName, err := i.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{"some", fieldName, "in", SliceConstructor(i.Values)}, " "), nil
}
