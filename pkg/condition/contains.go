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
	return subjectRego(c.GetSubject(ctx), c.Value, func(subject shared.Rego, value any, ctx *shared.Context) (string, error) {
		fieldName, err := subject.Rego(ctx)
		if err != nil {
			return "", err
		}
		return strings.Join([]string{shared.RegexExp, "(", "\"", ".*", fmt.Sprint(value), ".*", "\"", ",", "\"", fieldName, "\"", ")"}, ""), nil
	}, ctx)
}
