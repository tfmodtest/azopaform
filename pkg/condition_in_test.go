package pkg

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInCondition(t *testing.T) {
	cases := []struct {
		desc  string
		left  Rego
		right []string
		setup func(ctx context.Context)
		allow bool
	}{
		{
			desc:  "in",
			left:  stringRego(`"right"`),
			right: []string{"right", "left"},
			allow: true,
		},
		{
			desc:  "in_negative",
			left:  stringRego(`"left"`),
			right: []string{"right", "middle"},
			allow: false,
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			ctx := NewContext()
			if c.setup != nil {
				c.setup(ctx)
			}
			sut := InCondition{
				condition: condition{
					Subject: c.left,
				},
				Values: c.right,
			}
			actual, err := sut.Rego(ctx)
			require.NoError(t, err)
			cfg := fmt.Sprintf(conditionRegoTemplate, actual)
			assertRegoAllow(t, cfg, nil, c.allow, ctx)
		})
	}
}
