package pkg

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGreaterCondition(t *testing.T) {
	cases := []struct {
		desc  string
		left  Rego
		right string
		setup func(ctx context.Context)
		allow bool
	}{
		{
			desc:  "greater",
			left:  stringRego("2"),
			right: "1",
			allow: true,
		},
		{
			desc:  "less",
			left:  stringRego("1"),
			right: "2",
			allow: false,
		},
		{
			desc:  "equal",
			left:  stringRego("1"),
			right: "1",
			allow: false,
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			ctx := NewContext()
			if c.setup != nil {
				c.setup(ctx)
			}
			sut := GreaterCondition{
				condition: condition{
					Subject: c.left,
				},
				Value: c.right,
			}
			actual, err := sut.Rego(ctx)
			require.NoError(t, err)
			cfg := fmt.Sprintf(conditionRegoTemplate, actual)
			assertRegoAllow(t, cfg, nil, c.allow, ctx)
		})
	}
}
