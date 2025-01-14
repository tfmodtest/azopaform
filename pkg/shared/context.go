package shared

import (
	"context"
	"github.com/emirpasic/gods/stacks"
	"github.com/emirpasic/gods/stacks/arraystack"
)

type Context struct {
	context.Context
	resourceTypeStack      stacks.Stack
	fieldNameReplacerStack stacks.Stack
	conditionNameCounter   stacks.Stack
}

func NewContext() *Context {
	//contextMap := make(map[string]stacks.Stack)
	//contextMap["resourceType"] = arraystack.New()
	//contextMap["fieldNameReplacer"] = arraystack.New()
	//contextMap["conditionNameCounter"] = arraystack.New()
	//ctx := context.WithValue(context.Background(), "context", contextMap)
	return &Context{
		Context:                context.Background(),
		resourceTypeStack:      arraystack.New(),
		fieldNameReplacerStack: arraystack.New(),
		conditionNameCounter:   arraystack.New(),
	}
}

func (c *Context) PushFieldName(name string) {
	c.fieldNameReplacerStack.Push(name)
}

func (c *Context) currentResourceType() (string, bool) {
	value, ok := c.resourceTypeStack.Peek()
	if !ok {
		return "", false
	}
	return value.(string), true
}

func (c *Context) PopConditionNameCounter() (int, bool) {
	value, ok := c.conditionNameCounter.Pop()
	if !ok {
		return -1, false
	}
	return value.(int), true
}

func (c *Context) ClearConditionNameCounter() {
	c.conditionNameCounter.Clear()
}

func (c *Context) PushConditionNameCounter(counter int) {
	c.conditionNameCounter.Push(counter)
}

func (c *Context) FieldNameReplacer() (string, bool) {
	value, ok := c.fieldNameReplacerStack.Peek()
	if !ok {
		return "", false
	}
	return value.(string), true
}

func (c *Context) PushResourceType(rt string) {
	//contextMap := c.Context.Value("context").(map[string]stacks.Stack)
	//contextMap["resourceType"].Push(rt)
	c.resourceTypeStack.Push(rt)
}
