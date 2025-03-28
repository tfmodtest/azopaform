package operation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplaceIndex(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple string without indexes",
			input:    "properties.name",
			expected: "properties.name",
		},
		{
			name:     "String with numeric index",
			input:    "properties.array[0].name",
			expected: "properties.array[_].name",
		},
		{
			name:     "String with multiple numeric indexes",
			input:    "properties.array[1].items[2].value",
			expected: "properties.array[_].items[_].value",
		},
		{
			name:     "String with wildcard",
			input:    "properties.array[*].name",
			expected: "properties.array[_].name",
		},
		{
			name:     "String with mixed wildcards and indexes",
			input:    "properties.array[*].items[3].values[*]",
			expected: "properties.array[_].items[_].values[_]",
		},
		{
			name:     "String with trailing index",
			input:    "properties.values[0]",
			expected: "properties.values[_]",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "String with only index",
			input:    "[0]",
			expected: "[_]",
		},
		{
			name:     "String with multiple consecutive indexes",
			input:    "items[0][1].value",
			expected: "items[_][_].value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceIndex(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
