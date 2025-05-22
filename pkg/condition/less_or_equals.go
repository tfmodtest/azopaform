package condition

import (
	"fmt"
	"strings"

	"github.com/tfmodtest/azopaform/pkg/shared"
)

var _ Condition = LessOrEquals{}

type LessOrEquals struct {
	BaseCondition
	Value any
}

func (l LessOrEquals) Rego(ctx *shared.Context) (string, error) {
	return subjectRego(l.GetSubject(ctx), l.Value, func(subject shared.Rego, value any, ctx *shared.Context) (string, error) {
		fieldName, err := subject.Rego(ctx)
		if err != nil {
			return "", err
		}
		return strings.Join([]string{fieldName, "<=", fmt.Sprint(value)}, " "), nil
	}, ctx)
}
