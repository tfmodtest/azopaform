package shared

import (
	"fmt"
	"strings"
)

var functions = map[string]*FunctionCall{
	"parameters": {
		Name: "parameters",
		call: func(params []any, c *Context) (any, error) {
			if len(params) != 1 {
				return nil, fmt.Errorf("parameters function expects exactly one argument")
			}
			paramName, ok := params[0].(string)
			if !ok {
				return nil, fmt.Errorf("parameter name must be a string")
			}

			if c.GetParameterFunc != nil {
				value, ok, err := c.GetParameterFunc(paramName)
				if err != nil {
					return nil, err
				}
				if !ok {
					return nil, fmt.Errorf("parameter %s not found", paramName)
				}
				return value, nil
			}
			return nil, fmt.Errorf("GetParameterFunc not set")
		},
	},
}

type FunctionCall struct {
	Name       string
	Parameters []any
	call       func(params []any, c *Context) (any, error)
}

// ParseFunctionCall extracts function name and parameters from a string if it's a function call
func ParseFunctionCall(str string) (*FunctionCall, bool) {
	// Function calls are wrapped in square brackets like [functionName(param1, param2, ...)]
	if !strings.HasPrefix(str, "[") || !strings.HasSuffix(str, "]") {
		return nil, false
	}

	// Remove the outer brackets
	content := str[1 : len(str)-1]

	// Find the opening parenthesis for function arguments
	openParenIndex := strings.Index(content, "(")
	if openParenIndex == -1 {
		return nil, false // Not a function call if no opening parenthesis
	}

	// Find the closing parenthesis
	closeParenIndex := strings.LastIndex(content, ")")
	if closeParenIndex == -1 || closeParenIndex < openParenIndex {
		return nil, false // Not a valid function call if no proper closing parenthesis
	}

	// Extract function name
	functionName := strings.TrimSpace(content[:openParenIndex])

	// Extract parameters string
	paramsStr := content[openParenIndex+1 : closeParenIndex]

	// Parse parameters (simple implementation for basic cases)
	parameters := make([]any, 0)
	if paramsStr != "" {
		// This simple parsing works for basic cases like 'paramName'
		// For complex cases with nested quotes, comma separators in string literals, etc.,
		// a more sophisticated parser would be needed
		paramStrings := strings.Split(paramsStr, ",")
		for _, paramStr := range paramStrings {
			paramStr = strings.TrimSpace(paramStr)
			// Handle quoted strings
			if strings.HasPrefix(paramStr, "'") && strings.HasSuffix(paramStr, "'") {
				parameters = append(parameters, paramStr[1:len(paramStr)-1])
			} else {
				parameters = append(parameters, paramStr)
			}
		}
	}

	return &FunctionCall{
		Name:       functionName,
		Parameters: parameters,
	}, true
}

// EvaluateFunctionCall handles different function evaluations based on function name
func EvaluateFunctionCall(call *FunctionCall, c *Context) (any, error) {
	if function, exists := functions[call.Name]; exists {
		return function.call(call.Parameters, c)
	}
	return nil, fmt.Errorf("unsupported function %s", call.Name)
}
