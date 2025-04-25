package condition

import (
	"fmt"
	"github.com/tfmodtest/azopaform/pkg/shared"
	"strings"
)

var _ Condition = NotContains{}

type NotContains struct {
	BaseCondition
	Value string
}

func (n NotContains) Rego(ctx *shared.Context) (string, error) {
	fieldName, err := n.GetSubject(ctx).Rego(ctx)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{shared.Not, " ", shared.RegexExp, "(", "\"", ".*", fmt.Sprint(n.Value), ".*", "\"", ",", "\"", fieldName, "\"", ")"}, ""), nil
}
