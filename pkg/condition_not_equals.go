package pkg

import (
	"context"
	"fmt"
	"github.com/emirpasic/gods/stacks"
	"reflect"
	"strings"
)

var _ Condition = NotEqualsCondition{}

type NotEqualsCondition struct {
	condition
	Value any
}

func (n NotEqualsCondition) Rego(ctx context.Context) (string, error) {
	fieldName, err := n.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
		fieldName = replaceIndex(fieldName)
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
