package condition

import (
	"fmt"
	"strings"

	"github.com/tfmodtest/azopaform/pkg/shared"
)

var _ Condition = Contains{}

type Contains struct {
	BaseCondition
	Value string
}

func (c Contains) Rego(ctx *shared.Context) (string, error) {
	fieldName, err := c.GetSubject(ctx).Rego(ctx)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{shared.RegexExp, "(", "\"", ".*", fmt.Sprint(c.Value), ".*", "\"", ",", "\"", fieldName, "\"", ")"}, ""), nil
}
