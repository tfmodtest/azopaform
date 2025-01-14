package condition

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"json-rule-finder/pkg/shared"
	"testing"
)

func TestInCondition(t *testing.T) {
	cases := []struct {
		desc  string
		left  shared.Rego
		right []string
		setup func(ctx context.Context)
		allow bool
	}{
		{
			desc:  "in",
			left:  shared.StringRego(`"right"`),
			right: []string{"right", "left"},
			allow: true,
		},
		{
			desc:  "in_negative",
			left:  shared.StringRego(`"left"`),
			right: []string{"right", "middle"},
			allow: false,
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			ctx := shared.NewContext()
			if c.setup != nil {
				c.setup(ctx)
			}
			sut := In{
				BaseCondition: BaseCondition{
					Subject: c.left,
				},
				Values: c.right,
			}
			actual, err := sut.Rego(ctx)
			require.NoError(t, err)
			cfg := fmt.Sprintf(shared.RegoTestTemplate, actual)
			shared.AssertRegoAllow(t, cfg, nil, c.allow, ctx)
		})
	}
}
