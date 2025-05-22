package condition

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/tfmodtest/azopaform/pkg/shared"
)

var _ Condition = NotEquals{}

type NotEquals struct {
	BaseCondition
	Value any
}

func (n NotEquals) Rego(ctx *shared.Context) (string, error) {
	return subjectRego(n.GetSubject(ctx), n.Value, func(subject shared.Rego, value any, ctx *shared.Context) (string, error) {
		fieldName, err := subject.Rego(ctx)
		if err != nil {
			return "", err
		}
		var v string
		if reflect.TypeOf(value).Kind() == reflect.String {
			v = strings.Join([]string{"\"", fmt.Sprint(value), "\""}, "")
		} else if reflect.TypeOf(value).Kind() == reflect.Bool {
			v = fmt.Sprint(value)
		} else {
			v = fmt.Sprint(value)
		}
		return strings.Join([]string{fieldName, "!=", v}, " "), nil
	}, ctx)
}
