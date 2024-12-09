package pkg

import "context"

type IfBody map[string]any

func (i IfBody) condition(ctx context.Context) (*RuleSet, error) {
	return conditionFinder(i, ctx)
}
