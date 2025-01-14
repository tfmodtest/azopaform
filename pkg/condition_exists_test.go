package pkg

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"json-rule-finder/pkg/shared"
	"testing"
)

func TestExistsCondition(t *testing.T) {
	cases := []struct {
		desc  string
		left  shared.Rego
		right string
		setup func(ctx context.Context)
		allow bool
	}{
		{
			desc:  "exists",
			left:  shared.StringRego(`"ingress[0].transport"`),
			right: "true",
			allow: true,
		},
		{
			desc:  "not_exists",
			left:  shared.StringRego(`"ingress[0].transport"`),
			right: "false",
			allow: false,
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			ctx := NewContext()
			if c.setup != nil {
				c.setup(ctx)
			}
			sut := ExistsCondition{
				BaseCondition: BaseCondition{
					Subject: c.left,
				},
				Value: c.right,
			}
			actual, err := sut.Rego(ctx)
			require.NoError(t, err)
			cfg := fmt.Sprintf(shared.RegoTestTemplate, actual)
			shared.AssertRegoAllow(t, cfg, nil, c.allow, ctx)
		})
	}
}
