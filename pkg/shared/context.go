package shared

import (
	"context"
	"github.com/emirpasic/gods/stacks"
	"github.com/emirpasic/gods/stacks/arraystack"
)

type Context struct {
	context.Context
}

func NewContext() *Context {
	contextMap := make(map[string]stacks.Stack)
	contextMap["resourceType"] = arraystack.New()
	contextMap["fieldNameReplacer"] = arraystack.New()
	contextMap["conditionNameCounter"] = arraystack.New()
	ctx := context.WithValue(context.Background(), "context", contextMap)
	return &Context{
		Context: ctx,
	}
}

func (c *Context) PushResourceType(rt string) {
	contextMap := c.Context.Value("context").(map[string]stacks.Stack)
	contextMap["resourceType"].Push(rt)
}
