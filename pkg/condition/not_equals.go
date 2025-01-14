package condition

import (
	"fmt"
	"github.com/emirpasic/gods/stacks"
	"json-rule-finder/pkg/shared"
	"reflect"
	"strings"
)

var _ Condition = NotEquals{}

type NotEquals struct {
	BaseCondition
	Value any
}

func (n NotEquals) Rego(ctx *shared.Context) (string, error) {
	fieldName, err := n.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
		fieldName = ReplaceIndex(fieldName)
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
