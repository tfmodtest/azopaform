package value

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tfmodtest/azopaform/pkg/shared"
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
			expectedValue: `"test"`,
		},
		{
			name:          "Empty string",
			input:         "",
			expectedValue: `""`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := shared.NewContext()

			// Test creation of LiteralValue
			result, err := NewLiteralValue(tt.input, ctx)
			require.NoError(t, err)

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
