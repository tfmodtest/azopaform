package value

import (
	"json-rule-finder/pkg/shared"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLiteralValue(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedValue string
	}{
		{
			name:          "Simple string value",
			input:         "test",
			expectedValue: "test",
		},
		{
			name:          "String with wildcards",
			input:         "array[*].value",
			expectedValue: "array[_].value",
		},
		{
			name:          "String with multiple wildcards",
			input:         "array[*].items[*].name",
			expectedValue: "array[_].items[_].name",
		},
		{
			name:          "Empty string",
			input:         "",
			expectedValue: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := shared.NewContext()

			// Test creation of LiteralValue
			result := NewLiteralValue(tt.input, ctx)

			literalValue, ok := result.(LiteralValue)
			assert.True(t, ok, "Result should be a LiteralValue")
			assert.Equal(t, tt.expectedValue, literalValue.Value)

			// Test Rego() method
			regoValue, err := literalValue.Rego(ctx)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedValue, regoValue)
		})
	}
}
