package condition

import (
	"fmt"
	"strings"

	"github.com/tfmodtest/azopaform/pkg/shared"
)

var _ Condition = Greater{}

type Greater struct {
	BaseCondition
	Value any
}

func (g Greater) Rego(ctx *shared.Context) (string, error) {
	fieldName, err := g.GetSubject(ctx).Rego(ctx)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{fieldName, ">", fmt.Sprint(g.Value)}, " "), nil
}
