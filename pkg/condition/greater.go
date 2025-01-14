package condition

import (
	"fmt"
	"github.com/emirpasic/gods/stacks"
	"json-rule-finder/pkg/shared"
	"strings"
)

var _ Condition = Greater{}

type Greater struct {
	BaseCondition
	Value any
}

func (g Greater) Rego(ctx *shared.Context) (string, error) {
	fieldName, err := g.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
		fieldName = ReplaceIndex(fieldName)
	}
	return strings.Join([]string{fieldName, ">", fmt.Sprint(g.Value)}, " "), nil
}
