package pkg

import (
	"context"
	"fmt"
	"github.com/emirpasic/gods/stacks"
	"strings"
)

var _ Condition = LikeCondition{}

type LikeCondition struct {
	condition
	Value string
}

func (l LikeCondition) Rego(ctx context.Context) (string, error) {
	fieldName, err := l.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
		fieldName = replaceIndex(fieldName)
	}

	return strings.Join([]string{regexExp, "(", "\"", fmt.Sprintf(l.Value), "\"", ",", "\"", fieldName, "\"", ")"}, ""), nil
}
