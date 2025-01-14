package condition

import (
	"context"
	"github.com/emirpasic/gods/stacks"
	"json-rule-finder/pkg/shared"
	"reflect"
	"strings"
)

var _ Condition = Exists{}

type Exists struct {
	BaseCondition
	Value any
}

func (e Exists) Rego(ctx context.Context) (string, error) {
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
		return strings.Join([]string{shared.Not, fieldName}, " "), nil
	}
}
