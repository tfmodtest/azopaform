package condition

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tfmodtest/azopaform/pkg/shared"
)

func TestLikeCondition(t *testing.T) {
	cases := []struct {
		desc  string
		left  shared.Rego
		right string
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
			ctx := shared.NewContext()
			sut := Like{
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
