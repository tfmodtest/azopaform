package condition

import (
	"context"
	"json-rule-finder/pkg/shared"
	"strings"
)

var _ Condition = NotIn{}

type NotIn struct {
	BaseCondition
	Values []string
}

func (n NotIn) Rego(ctx context.Context) (string, error) {
	fieldName, err := n.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{shared.Not, fieldName, "in", shared.SliceConstructor(n.Values)}, " "), nil
}
