package pkg

import (
	"context"
	"json-rule-finder/pkg/shared"
	"strings"
)

var _ Condition = NotInCondition{}

type NotInCondition struct {
	BaseCondition
	Values []string
}

func (n NotInCondition) Rego(ctx context.Context) (string, error) {
	fieldName, err := n.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{shared.Not, fieldName, "in", shared.SliceConstructor(n.Values)}, " "), nil
}
