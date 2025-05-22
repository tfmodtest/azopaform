package condition

import (
	"fmt"
	"strings"

	"github.com/tfmodtest/azopaform/pkg/shared"
)

var _ Condition = GreaterOrEquals{}

type GreaterOrEquals struct {
	BaseCondition
	Value any
}

func (g GreaterOrEquals) Rego(ctx *shared.Context) (string, error) {
	return subjectRego(g.GetSubject(ctx), g.Value, func(subject shared.Rego, value any, ctx *shared.Context) (string, error) {
		fieldName, err := subject.Rego(ctx)
		if err != nil {
			return "", err
		}
		return strings.Join([]string{fieldName, ">=", fmt.Sprint(value)}, " "), nil
	}, ctx)
}
