package shared

import (
	"context"
	"encoding/json"
	"fmt"
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
	GetParameterFunc       func(string) (any, bool, error)
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

func ResolveParameterValue[T any](input any, c *Context) (T, error) {
	str, ok := input.(string)
	if !ok {
		return input.(T), nil
	}
	if funcCall, ok := ParseFunctionCall(str); ok {
		value, err := EvaluateFunctionCall(funcCall, c)
		if err != nil {
			return func() T {
				var defaultT T
				return defaultT
			}(), err
		}
		return value.(T), nil
	}

	// Return original input if not a parameter reference or parameter not found
	return input.(T), nil
}

func ResolveParameterValueAsString(input any, c *Context) (string, error) {
	value, err := ResolveParameterValue[any](input, c)
	if err != nil {
		return "", err
	}
	strValue, ok := value.(string)
	if ok {
		return strValue, nil
	}
	// Handle other types
	switch v := value.(type) {
	case bool:
		return fmt.Sprintf("%t", v), nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v), nil
	case float32, float64:
		return fmt.Sprintf("%g", v), nil
	case []interface{}, map[string]interface{}:
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		return string(jsonBytes), nil
	default:
		// For any other type, convert to string representation
		return fmt.Sprintf("%v", v), nil
	}
}

func getOrDefault[T comparable](value, defaultValue T) T {
	var defaultTValue T
	if value == defaultTValue {
		return defaultValue
	}
	return value
}
