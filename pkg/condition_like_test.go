package pkg

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"json-rule-finder/pkg/shared"
	"testing"
)

func TestLikeCondition(t *testing.T) {
	cases := []struct {
		desc  string
		left  shared.Rego
		right string
		setup func(ctx context.Context)
		allow bool
	}{
		{
			desc:  "like",
			left:  shared.StringRego("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.ApiManagement/service/instance1/apis"),
			right: ".*/apis",
			allow: true,
		},
		{
			desc:  "not_like",
			left:  shared.StringRego("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.ApiManagement/service/instance1/apis"),
			right: ".*/apis2",
			allow: false,
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			ctx := NewContext()
			if c.setup != nil {
				c.setup(ctx)
			}
			sut := LikeCondition{
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
