package condition

import (
	"fmt"
	"strings"

	"github.com/tfmodtest/azopaform/pkg/shared"
)

var _ Condition = NotContains{}

type NotContains struct {
	BaseCondition
	Value string
}

func (n NotContains) Rego(ctx *shared.Context) (string, error) {
	return subjectRego(n.GetSubject(ctx), n.Value, func(subject shared.Rego, value any, ctx *shared.Context) (string, error) {
		fieldName, err := subject.Rego(ctx)
		if err != nil {
			return "", err
		}
		return strings.Join([]string{shared.Not, " ", shared.RegexExp, "(", "\"", ".*", fmt.Sprint(value), ".*", "\"", ",", "\"", fieldName, "\"", ")"}, ""), nil
	}, ctx)
}
