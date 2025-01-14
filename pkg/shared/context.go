package shared

import (
	"context"
	"github.com/emirpasic/gods/stacks"
	"github.com/emirpasic/gods/stacks/arraystack"
)

func NewContext() context.Context {
	contextMap := make(map[string]stacks.Stack)
	contextMap["resourceType"] = arraystack.New()
	contextMap["fieldNameReplacer"] = arraystack.New()
	contextMap["conditionNameCounter"] = arraystack.New()
	ctx := context.WithValue(context.Background(), "context", contextMap)
	return ctx
}

func PushResourceType(ctx context.Context, rt string) {
	contextMap := ctx.Value("context").(map[string]stacks.Stack)
	contextMap["resourceType"].Push(rt)
}
