package condition

import (
	"fmt"
	"json-rule-finder/pkg/shared"
	"strings"
)

var _ Condition = Less{}

type Less struct {
	BaseCondition
	Value any
}

func (l Less) Rego(ctx *shared.Context) (string, error) {
	fieldName, err := l.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	if _, ok := ctx.FieldNameReplacer(); ok {
		fieldName = ReplaceIndex(fieldName)
	}
	return strings.Join([]string{fieldName, "<", fmt.Sprint(l.Value)}, " "), nil
}
