package operation

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tfmodtest/azopaform/pkg/shared"
)

func TestCount(t *testing.T) {
	//t.Skip("skip count for now")
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
				"field": "Microsoft.Network/networkSecurityGroups/zones[*]",
			},
			input: map[string]any{
				"values": map[string]any{
					"properties": map[string]any{
						"zones": []any{
							1, 2, 3,
						},
					},
				},
			},
			query:    "c",
			expected: 3,
		},
		{
			name: "with where",
			unparsed: map[string]any{
				"field": "Microsoft.Network/networkSecurityGroups/zones[*]",
				"where": map[string]any{
					"field":           "Microsoft.Network/networkSecurityGroups/zones[*]",
					"greaterOrEquals": 3,
				},
			},
			input: map[string]any{
				"values": map[string]any{
					"properties": map[string]any{
						"zones": []any{
							2, 3, 4,
						},
					},
				},
			},
			query:    "c",
			expected: 2,
		},
		{
			name: "with where and longer field",
			unparsed: map[string]any{
				"field": "Microsoft.Network/networkSecurityGroups/zones[*]",
				"where": map[string]any{
					"field":  "Microsoft.Network/networkSecurityGroups/zones[*].name",
					"equals": "test",
				},
			},
			input: map[string]any{
				"values": map[string]any{
					"properties": map[string]any{
						"zones": []map[string]any{
							{
								"name": "test",
							},
							{
								"name": "test2",
							},
						},
					},
				},
			},
			query:    "c",
			expected: 1,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := shared.NewContext()
			ctx.PushResourceType("Microsoft.Network/networkSecurityGroups")
			sut, err := NewCount(c.unparsed, ctx)
			require.NoError(t, err)
			exp, err := sut.Rego(ctx)
			require.NoError(t, err)
			jsCount, ok := shared.EvaluateRego(t, fmt.Sprintf("data.main.%s", c.query), fmt.Sprintf(`
	package main
	
	import rego.v1
	r := input
	c := %s

%s
`, exp, ctx.HelperFunctionsRego()), c.input, ctx).(json.Number)
			require.True(t, ok)
			count, err := strconv.Atoi(jsCount.String())
			require.NoError(t, err)
			assert.Equal(t, c.expected, count)
		})
	}
}
