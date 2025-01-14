package condition

import (
	"json-rule-finder/pkg/shared"
	"reflect"
	"strings"
)

var _ Condition = Exists{}

type Exists struct {
	BaseCondition
	Value any
}

func (e Exists) Rego(ctx *shared.Context) (string, error) {
	fieldName, err := e.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	if _, ok := ctx.FieldNameReplacer(); ok {
		fieldName = ReplaceIndex(fieldName)
	}
	if (reflect.TypeOf(e.Value).Kind() == reflect.Bool && e.Value.(bool)) || (reflect.TypeOf(e.Value).Kind() == reflect.String && e.Value.(string) == "true") {
		return fieldName, nil
	} else {
		return strings.Join([]string{shared.Not, fieldName}, " "), nil
	}
}
