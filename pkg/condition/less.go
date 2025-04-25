package condition

import (
	"fmt"
	"strings"

	"github.com/tfmodtest/azopaform/pkg/shared"
)

var _ Condition = Less{}

type Less struct {
	BaseCondition
	Value any
}

func (l Less) Rego(ctx *shared.Context) (string, error) {
	fieldName, err := l.GetSubject(ctx).Rego(ctx)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{fieldName, "<", fmt.Sprint(l.Value)}, " "), nil
}
