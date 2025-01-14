package pkg

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"json-rule-finder/pkg/shared"
	"testing"
)

func TestLessCondition(t *testing.T) {
	cases := []struct {
		desc  string
		left  shared.Rego
		right string
		setup func(ctx context.Context)
		allow bool
	}{
		{
			desc:  "less",
			left:  shared.StringRego("1"),
			right: "2",
			allow: true,
		},
		{
			desc:  "greater",
			left:  shared.StringRego("2"),
			right: "1",
			allow: false,
		},
		{
			desc:  "equals",
			left:  shared.StringRego("1"),
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
			sut := LessCondition{
				BaseCondition: BaseCondition{
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
