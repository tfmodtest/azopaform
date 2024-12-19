package pkg

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var _ Rego = stringRego("")

type stringRego string

func (s stringRego) Rego(ctx context.Context) (string, error) {
	return string(s), nil
}

func TestEqualsCondition(t *testing.T) {
	cases := []struct {
		desc     string
		left     Rego
		right    any
		setup    func(ctx context.Context)
		expected string
	}{
		{
			desc:     "string",
			left:     stringRego("left"),
			right:    "right",
			expected: `left == "right"`,
		},
		{
			desc:     "int",
			left:     stringRego("left"),
			right:    1,
			expected: `left == 1`,
		},
		{
			desc:     "bool",
			left:     stringRego("left"),
			right:    true,
			expected: `left == true`,
		},
		{
			desc:  "field equals",
			left:  OperationField("Microsoft.Web/serverFarms/sku.tier"),
			right: "Standard",
			setup: func(ctx context.Context) {
				pushResourceType(ctx, "Microsoft.Web/serverFarms")
			},
			expected: `r.change.after.sku[0].tier == "Standard"`,
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			ctx := NewContext()
			if c.setup != nil {
				c.setup(ctx)
			}
			sut := EqualsCondition{
				condition: condition{
					Subject: c.left,
				},
				Value: c.right,
			}
			actual, err := sut.Rego(ctx)
			require.NoError(t, err)
			assert.Equal(t, c.expected, actual)
		})
	}
}
