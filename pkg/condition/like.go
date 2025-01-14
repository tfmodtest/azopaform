package condition

import (
	"fmt"
	"json-rule-finder/pkg/shared"
	"strings"
)

var _ Condition = Like{}

type Like struct {
	BaseCondition
	Value string
}

func (l Like) Rego(ctx *shared.Context) (string, error) {
	fieldName, err := l.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	if _, ok := ctx.FieldNameReplacer(); ok {
		fieldName = ReplaceIndex(fieldName)
	}

	return strings.Join([]string{shared.RegexExp, "(", "\"", fmt.Sprintf(l.Value), "\"", ",", "\"", fieldName, "\"", ")"}, ""), nil
}
