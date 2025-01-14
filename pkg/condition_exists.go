package pkg

import (
	"context"
	"github.com/emirpasic/gods/stacks"
	"reflect"
	"strings"
)

var _ Condition = ExistsCondition{}

type ExistsCondition struct {
	condition
	Value any
}

func (e ExistsCondition) Rego(ctx context.Context) (string, error) {
	fieldName, err := e.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
		fieldName = ReplaceIndex(fieldName)
	}
	if (reflect.TypeOf(e.Value).Kind() == reflect.Bool && e.Value.(bool)) || (reflect.TypeOf(e.Value).Kind() == reflect.String && e.Value.(string) == "true") {
		return fieldName, nil
	} else {
		return strings.Join([]string{not, fieldName}, " "), nil
	}
}
