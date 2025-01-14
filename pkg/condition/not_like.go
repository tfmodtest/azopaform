package condition

import (
	"fmt"
	"json-rule-finder/pkg/shared"
	"strings"
)

var _ Condition = NotLike{}

type NotLike struct {
	BaseCondition
	Value string
}

func (n NotLike) Rego(ctx *shared.Context) (string, error) {
	fieldName, err := n.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	if _, ok := ctx.FieldNameReplacer(); ok {
		fieldName = ReplaceIndex(fieldName)
	}
	v := strings.Join([]string{"`", fmt.Sprint(n.Value), "`"}, "")
	return strings.Join([]string{shared.Not, " ", shared.RegexExp, "(", v, ",", fieldName, ")"}, ""), nil
}
