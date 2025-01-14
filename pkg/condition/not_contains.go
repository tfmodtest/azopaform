package condition

import (
	"fmt"
	"json-rule-finder/pkg/shared"
	"strings"
)

var _ Condition = NotContains{}

type NotContains struct {
	BaseCondition
	Value string
}

func (n NotContains) Rego(ctx *shared.Context) (string, error) {
	fieldName, err := n.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{shared.Not, " ", shared.RegexExp, "(", "\"", ".*", fmt.Sprint(n.Value), ".*", "\"", ",", "\"", fieldName, "\"", ")"}, ""), nil
}
