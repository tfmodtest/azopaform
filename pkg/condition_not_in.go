package pkg

import (
	"context"
	"strings"
)

var _ Condition = NotInCondition{}

type NotInCondition struct {
	condition
	Values []string
}

func (n NotInCondition) Rego(ctx context.Context) (string, error) {
	fieldName, err := n.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{not, fieldName, "in", SliceConstructor(n.Values)}, " "), nil
}
