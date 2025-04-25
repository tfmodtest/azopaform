package condition

import (
	"fmt"
	"github.com/tfmodtest/azopaform/pkg/shared"
	"strings"
)

var _ Condition = NotLike{}

type NotLike struct {
	BaseCondition
	Value string
}

func (n NotLike) Rego(ctx *shared.Context) (string, error) {
	fieldName, err := n.GetSubject(ctx).Rego(ctx)
	if err != nil {
		return "", err
	}
	v := strings.Join([]string{"`", fmt.Sprint(n.Value), "`"}, "")
	return strings.Join([]string{shared.Not, " ", shared.RegexExp, "(", v, ",", fieldName, ")"}, ""), nil
}
