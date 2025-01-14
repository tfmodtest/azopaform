package condition

import (
	"fmt"
	"json-rule-finder/pkg/shared"
	"strings"
)

var _ Condition = GreaterOrEquals{}

type GreaterOrEquals struct {
	BaseCondition
	Value any
}

func (g GreaterOrEquals) Rego(ctx *shared.Context) (string, error) {
	fieldName, err := g.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	if _, ok := ctx.FieldNameReplacer(); ok {
		fieldName = ReplaceIndex(fieldName)
	}
	return strings.Join([]string{fieldName, ">=", fmt.Sprint(g.Value)}, " "), nil
}
