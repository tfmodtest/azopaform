package pkg

import (
	"context"
	"fmt"
	"github.com/emirpasic/gods/stacks"
	"strings"
)

var _ Condition = LessCondition{}

type LessCondition struct {
	BaseCondition
	Value any
}

func (l LessCondition) Rego(ctx context.Context) (string, error) {
	fieldName, err := l.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
		fieldName = ReplaceIndex(fieldName)
	}
	return strings.Join([]string{fieldName, "<", fmt.Sprint(l.Value)}, " "), nil
}
