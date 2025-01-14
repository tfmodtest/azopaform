package pkg

import (
	"context"
	"fmt"
	"json-rule-finder/pkg/shared"
	"strings"
)

var _ Condition = ContainsCondition{}

type ContainsCondition struct {
	BaseCondition
	Value string
}

func (c ContainsCondition) Rego(ctx context.Context) (string, error) {
	fieldName, err := c.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{shared.RegexExp, "(", "\"", ".*", fmt.Sprint(c.Value), ".*", "\"", ",", "\"", fieldName, "\"", ")"}, ""), nil
}
