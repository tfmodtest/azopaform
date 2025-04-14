package shared

import (
	"context"
	"github.com/emirpasic/gods/stacks"
	"github.com/emirpasic/gods/stacks/arraystack"
	"strings"
)

type Context struct {
	context.Context
	resourceTypeStack      stacks.Stack
	fieldNameReplacerStack stacks.Stack
	helperFuncs            []string
}

func NewContext() *Context {
	return &Context{
		Context:                context.Background(),
		resourceTypeStack:      arraystack.New(),
		fieldNameReplacerStack: arraystack.New(),
	}
}

func (c *Context) PushFieldName(name string) {
	c.fieldNameReplacerStack.Push(name)
}

func (c *Context) PopFieldName() {
	c.fieldNameReplacerStack.Pop()
}

func (c *Context) InHelperFunction(parameterName string, action func() error) error {
	c.PushFieldName(parameterName)
	defer c.PopFieldName()
	return action()
}

func (c *Context) currentResourceType() (string, bool) {
	value, ok := c.resourceTypeStack.Peek()
	if !ok {
		return "", false
	}
	return value.(string), true
}

func (c *Context) FieldNameReplacer() (string, bool) {
	value, ok := c.fieldNameReplacerStack.Peek()
	if !ok {
		return "", false
	}
	return value.(string), true
}

func (c *Context) PushResourceType(rt string) {
	c.resourceTypeStack.Push(rt)
}

func (c *Context) EnqueueHelperFunction(funcDec string) {
	c.helperFuncs = append(c.helperFuncs, funcDec)
}

func (c *Context) HelperFunctionsRego() string {
	sb := new(strings.Builder)
	for _, helperFunc := range c.helperFuncs {
		sb.WriteString(helperFunc)
		sb.WriteString("\n")
	}
	return sb.String()
}
