package condition

import (
	"context"
	"fmt"
	"github.com/emirpasic/gods/stacks"
	"json-rule-finder/pkg/shared"
	"strings"
)

var _ Condition = Like{}

type Like struct {
	BaseCondition
	Value string
}

func (l Like) Rego(ctx context.Context) (string, error) {
	fieldName, err := l.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
		fieldName = ReplaceIndex(fieldName)
	}

	return strings.Join([]string{shared.RegexExp, "(", "\"", fmt.Sprintf(l.Value), "\"", ",", "\"", fieldName, "\"", ")"}, ""), nil
}
