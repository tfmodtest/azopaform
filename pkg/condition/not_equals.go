package condition

import (
	"fmt"
	"github.com/tfmodtest/azopaform/pkg/shared"
	"reflect"
	"strings"
)

var _ Condition = NotEquals{}

type NotEquals struct {
	BaseCondition
	Value any
}

func (n NotEquals) Rego(ctx *shared.Context) (string, error) {
	fieldName, err := n.GetSubject(ctx).Rego(ctx)
	if err != nil {
		return "", err
	}
	var v string
	if reflect.TypeOf(n.Value).Kind() == reflect.String {
		v = strings.Join([]string{"\"", fmt.Sprint(n.Value), "\""}, "")
	} else if reflect.TypeOf(n.Value).Kind() == reflect.Bool {
		v = fmt.Sprint(n.Value)
	} else {
		v = fmt.Sprint(n.Value)
	}
	return strings.Join([]string{fieldName, "!=", v}, " "), nil
}
