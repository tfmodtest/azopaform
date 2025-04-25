package condition

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tfmodtest/azopaform/pkg/shared"
)

func TestContainsCondition(t *testing.T) {
	cases := []struct {
		desc  string
		left  shared.Rego
		right string
		allow bool
	}{
		{
			desc:  "contains_negative",
			left:  shared.StringRego("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.RecoveryServices/vaults/vault1"),
			right: "Microsoft.Web/sites",
			allow: false,
		},
		{
			desc:  "contains",
			left:  shared.StringRego("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/sites/site1/slots/slot1"),
			right: "Microsoft.Web/sites",
			allow: true,
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			ctx := shared.NewContext()

			sut := Contains{
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
