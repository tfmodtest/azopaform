package pkg

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"json-rule-finder/pkg/shared"
	"testing"
)

func TestNotContainsCondition(t *testing.T) {
	cases := []struct {
		desc  string
		left  shared.Rego
		right string
		setup func(ctx context.Context)
		allow bool
	}{
		{
			desc:  "not_contains",
			left:  shared.StringRego("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.RecoveryServices/vaults/vault1"),
			right: "Microsoft.Web/sites",
			allow: true,
		},
		{
			desc:  "not_contains_negative",
			left:  shared.StringRego("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/sites/site1/slots/slot1"),
			right: "Microsoft.Web/sites",
			allow: false,
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			ctx := NewContext()
			if c.setup != nil {
				c.setup(ctx)
			}

			sut := NotContainsCondition{
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
