package condition

import (
	"fmt"
	"github.com/emirpasic/gods/stacks"
	"json-rule-finder/pkg/shared"
	"strings"
)

var _ Condition = LessOrEquals{}

type LessOrEquals struct {
	BaseCondition
	Value any
}

func (l LessOrEquals) Rego(ctx *shared.Context) (string, error) {
	fieldName, err := l.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
		fieldName = ReplaceIndex(fieldName)
	}
	return strings.Join([]string{fieldName, "<=", fmt.Sprint(l.Value)}, " "), nil
}
