package condition

import (
	"fmt"
	"json-rule-finder/pkg/shared"
	"strings"
)

var _ Condition = Contains{}

type Contains struct {
	BaseCondition
	Value string
}

func (c Contains) Rego(ctx *shared.Context) (string, error) {
	fieldName, err := c.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{shared.RegexExp, "(", "\"", ".*", fmt.Sprint(c.Value), ".*", "\"", ",", "\"", fieldName, "\"", ")"}, ""), nil
}
