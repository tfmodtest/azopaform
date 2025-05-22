package condition

import (
	"strings"

	"github.com/tfmodtest/azopaform/pkg/shared"
)

var _ Condition = Like{}

type Like struct {
	BaseCondition
	Value string
}

func (l Like) Rego(ctx *shared.Context) (string, error) {
	return subjectRego(l.GetSubject(ctx), l.Value, func(subject shared.Rego, value any, ctx *shared.Context) (string, error) {
		fieldName, err := subject.Rego(ctx)
		if err != nil {
			return "", err
		}
		return strings.Join([]string{shared.RegexExp, "(", "\"", value.(string), "\"", ",", "\"", fieldName, "\"", ")"}, ""), nil
	}, ctx)
}
