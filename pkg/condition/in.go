package condition

import (
	"context"
	"json-rule-finder/pkg/shared"
	"strings"
)

var _ Condition = In{}

type In struct {
	BaseCondition
	Values []string
}

func (i In) Rego(ctx context.Context) (string, error) {
	fieldName, err := i.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{"some", fieldName, "in", shared.SliceConstructor(i.Values)}, " "), nil
}
