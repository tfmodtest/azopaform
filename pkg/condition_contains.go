package pkg

import (
	"context"
	"fmt"
	"strings"
)

var _ Condition = ContainsCondition{}

type ContainsCondition struct {
	condition
	Value string
}

func (c ContainsCondition) Rego(ctx context.Context) (string, error) {
	fieldName, err := c.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{regexExp, "(", "\"", ".*", fmt.Sprint(c.Value), ".*", "\"", ",", "\"", fieldName, "\"", ")"}, ""), nil
}
