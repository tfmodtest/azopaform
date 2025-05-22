package condition

import (
	"fmt"
	"strings"

	"github.com/tfmodtest/azopaform/pkg/shared"
)

var _ Condition = NotLike{}

type NotLike struct {
	BaseCondition
	Value string
}

func (n NotLike) Rego(ctx *shared.Context) (string, error) {
	return subjectRego(n.GetSubject(ctx), n.Value, func(subject shared.Rego, value any, ctx *shared.Context) (string, error) {
		fieldName, err := subject.Rego(ctx)
		if err != nil {
			return "", err
		}
		v := strings.Join([]string{"`", fmt.Sprint(value), "`"}, "")
		return strings.Join([]string{shared.Not, " ", shared.RegexExp, "(", v, ",", fieldName, ")"}, ""), nil
	}, ctx)
}
