package pkg

import (
	"context"
	"fmt"
	"json-rule-finder/pkg/shared"
	"strings"
)

var _ Condition = NotContainsCondition{}

type NotContainsCondition struct {
	BaseCondition
	Value string
}

func (n NotContainsCondition) Rego(ctx context.Context) (string, error) {
	fieldName, err := n.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	//if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
	//	fieldName = ReplaceIndex(fieldName)
	//}
	//v := strings.Join([]string{"`", fmt.Sprint(n.Value), "`"}, "")
	return strings.Join([]string{shared.Not, " ", shared.RegexExp, "(", "\"", ".*", fmt.Sprint(n.Value), ".*", "\"", ",", "\"", fieldName, "\"", ")"}, ""), nil
}
