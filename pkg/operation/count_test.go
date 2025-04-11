package operation

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"json-rule-finder/pkg/shared"
	"strconv"
	"testing"
)

func TestCount(t *testing.T) {
	cases := []struct {
		name     string
		unparsed map[string]any
		input    map[string]any
		query    string
		expected int
	}{
		{
			name: "simple",
			unparsed: map[string]any{
				"value": "input.items",
			},
			input: map[string]any{
				"items": []any{
					1, 2, 3,
				},
			},
			query:    "c",
			expected: 3,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := shared.NewContext()
			sut := NewCount(c.unparsed, ctx)
			exp, err := sut.Rego(ctx)
			require.NoError(t, err)
			jsCount, ok := shared.EvaluateRego(t, fmt.Sprintf("data.main.%s", c.query), fmt.Sprintf(`
	package main
	
	import rego.v1
	
	c := %s
`, exp), c.input, ctx).(json.Number)
			require.True(t, ok)
			count, err := strconv.Atoi(jsCount.String())
			require.NoError(t, err)
			assert.Equal(t, c.expected, count)
		})
	}
}
