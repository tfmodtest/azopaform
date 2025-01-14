package pkg

import (
	"context"
	"fmt"
	"github.com/emirpasic/gods/stacks"
	"json-rule-finder/pkg/shared"
	"strings"
)

var _ Condition = NotLikeCondition{}

type NotLikeCondition struct {
	BaseCondition
	Value string
}

func (n NotLikeCondition) Rego(ctx context.Context) (string, error) {
	fieldName, err := n.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
		fieldName = ReplaceIndex(fieldName)
	}
	v := strings.Join([]string{"`", fmt.Sprint(n.Value), "`"}, "")
	return strings.Join([]string{shared.Not, " ", shared.RegexExp, "(", v, ",", fieldName, ")"}, ""), nil
}
