package shared

import (
	"context"
	"strings"

	"github.com/emirpasic/gods/stacks"
	"github.com/emirpasic/gods/stacks/arraystack"
)

type Context struct {
	context.Context
	option                 Options
	resourceTypeStack      stacks.Stack
	fieldNameReplacerStack stacks.Stack
	helperFunctions        []string
	GetParameterFunc       func(string) (any, bool)
}

func NewContext() *Context {
	return &Context{
		Context:                context.Background(),
		resourceTypeStack:      arraystack.New(),
		fieldNameReplacerStack: arraystack.New(),
	}
}

func NewContextWithOptions(option Options) *Context {
	ctx := NewContext()
	ctx.option = option
	return ctx
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
	c.helperFunctions = append(c.helperFunctions, funcDec)
}

func (c *Context) HelperFunctionsRego() string {
	sb := new(strings.Builder)
	for _, helperFunc := range c.helperFunctions {
		sb.WriteString(helperFunc)
		sb.WriteString("\n")
	}
	return sb.String()
}

func (c *Context) PackageName() string {
	return getOrDefault(c.option.PackageName, "main")
}

func (c *Context) UtilRegoFileName() string {
	return getOrDefault(c.option.UtilRegoFileName, "util.rego")
}

func (c *Context) GenerateRuleName() bool {
	return c.option.GenerateRuleName
}

func (c *Context) UtilLibraryPackageName() string {
	return c.option.UtilLibraryPackageName
}

func getOrDefault[T comparable](value, defaultValue T) T {
	var defaultTValue T
	if value == defaultTValue {
		return defaultValue
	}
	return value
}
