package condition

import (
	"github.com/tfmodtest/azopaform/pkg/shared"
	"strings"
)

var _ Condition = NotIn{}

type NotIn struct {
	BaseCondition
	Values []string
}

func (n NotIn) Rego(ctx *shared.Context) (string, error) {
	fieldName, err := n.GetSubject(ctx).Rego(ctx)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{shared.Not, fieldName, "in", shared.SliceConstructor(n.Values)}, " "), nil
}
