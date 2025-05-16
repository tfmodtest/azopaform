package shared

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFunctionCall(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		wantFunc       *FunctionCall
		wantIsFunction bool
	}{
		{
			name:           "Empty string",
			input:          "",
			wantFunc:       nil,
			wantIsFunction: false,
		},
		{
			name:           "Not a function call - no brackets",
			input:          "parameters('test')",
			wantFunc:       nil,
			wantIsFunction: false,
		},
		{
			name:           "Not a function call - missing opening bracket",
			input:          "parameters('test')]",
			wantFunc:       nil,
			wantIsFunction: false,
		},
		{
			name:           "Not a function call - missing closing bracket",
			input:          "[parameters('test')",
			wantFunc:       nil,
			wantIsFunction: false,
		},
		{
			name:           "Not a function call - no parentheses",
			input:          "[parameters]",
			wantFunc:       nil,
			wantIsFunction: false,
		},
		{
			name:           "Not a function call - missing closing parenthesis",
			input:          "[parameters('test']",
			wantFunc:       nil,
			wantIsFunction: false,
		},
		{
			name:           "Valid function call - no parameters",
			input:          "[resourceGroup()]",
			wantFunc:       &FunctionCall{Name: "resourceGroup", Parameters: []any{}},
			wantIsFunction: true,
		},
		{
			name:           "Valid function call - single quoted parameter",
			input:          "[parameters('storageAccountName')]",
			wantFunc:       &FunctionCall{Name: "parameters", Parameters: []any{"storageAccountName"}},
			wantIsFunction: true,
		},
		{
			name:           "Valid function call - multiple parameters",
			input:          "[concat('prefix-', 'suffix')]",
			wantFunc:       &FunctionCall{Name: "concat", Parameters: []any{"prefix-", "suffix"}},
			wantIsFunction: true,
		},
		{
			name:           "Valid function call - unquoted parameter",
			input:          "[length(variables)]",
			wantFunc:       &FunctionCall{Name: "length", Parameters: []any{"variables"}},
			wantIsFunction: true,
		},
		{
			name:           "Valid function call - multiple mixed parameters",
			input:          "[add(1, parameter)]",
			wantFunc:       &FunctionCall{Name: "add", Parameters: []any{"1", "parameter"}},
			wantIsFunction: true,
		},
		{
			name:           "Valid function call - extra whitespace",
			input:          "[  parameters ( 'test' ) ]",
			wantFunc:       &FunctionCall{Name: "parameters", Parameters: []any{"test"}},
			wantIsFunction: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFunc, gotIsFunction := ParseFunctionCall(tt.input)
			assert.Equal(t, tt.wantIsFunction, gotIsFunction)
			assert.Equal(t, tt.wantFunc, gotFunc)
		})
	}
}
