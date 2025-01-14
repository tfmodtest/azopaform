package condition

import (
	"fmt"
	"json-rule-finder/pkg/shared"
	"reflect"
	"strings"

	"github.com/emirpasic/gods/stacks"
)

var _ Condition = Equals{}

type Equals struct {
	BaseCondition
	Value any
}

// Rego For conditions under 'where' operator, "[[0-9]+]" should be replaced with "[x]"
func (e Equals) Rego(ctx *shared.Context) (string, error) {
	fieldName, err := e.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
		fieldName = ReplaceIndex(fieldName)
	}
	var v string
	if reflect.TypeOf(e.Value).Kind() == reflect.String {
		v = strings.Join([]string{"\"", fmt.Sprint(e.Value), "\""}, "")
	} else if reflect.TypeOf(e.Value).Kind() == reflect.Bool {
		v = fmt.Sprint(e.Value)
	} else {
		v = fmt.Sprint(e.Value)
	}
	return strings.Join([]string{fieldName, "==", v}, " "), nil
}

func (e Equals) GetReverseRego(ctx *shared.Context) (string, error) {
	fieldName, err := e.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	if stack, ok := ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack); ok && stack.Size() > 0 {
		fieldName = ReplaceIndex(fieldName)
	}
	var v string
	if reflect.TypeOf(e.Value).Kind() == reflect.String {
		v = strings.Join([]string{"\"", fmt.Sprint(e.Value), "\""}, "")
	} else if reflect.TypeOf(e.Value).Kind() == reflect.Bool {
		v = fmt.Sprint(e.Value)
	} else {
		v = fmt.Sprint(e.Value)
	}
	return strings.Join([]string{fieldName, "!=", v}, " "), nil
}
