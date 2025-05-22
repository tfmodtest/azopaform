package condition

import (
	"strings"

	"github.com/tfmodtest/azopaform/pkg/shared"
)

var _ Condition = NotIn{}

type NotIn struct {
	BaseCondition
	Values []string
}

func (n NotIn) Rego(ctx *shared.Context) (string, error) {
	return subjectRego(n.GetSubject(ctx), n.Values, func(subject shared.Rego, value any, context *shared.Context) (string, error) {
		fieldName, err := subject.Rego(ctx)
		if err != nil {
			return "", err
		}
		return strings.Join([]string{shared.Not, fieldName, "in", shared.SliceConstructor(value.([]string))}, " "), nil
	}, ctx)
}
