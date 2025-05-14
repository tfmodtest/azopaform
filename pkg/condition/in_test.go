package condition

import (
	"github.com/tfmodtest/azopaform/pkg/value"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tfmodtest/azopaform/pkg/shared"
)

func TestInCondition(t *testing.T) {
	cases := []struct {
		desc  string
		left  shared.Rego
		right []string
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
			sut := In{
				BaseCondition: BaseCondition{
					Subject: c.left,
				},
				Values: c.right,
			}
			actual, err := sut.Rego(ctx)
			require.NoError(t, err)
			cfg := shared.WithUtilFunctions(actual)
			shared.AssertRegoAllow(t, cfg, nil, c.allow, ctx)
		})
	}
}

func TestInConditionWithResourceType(t *testing.T) {
	cases := []struct {
		desc         string
		resourceType string
		right        []string
		allow        bool
	}{
		{
			desc:         "in",
			resourceType: "Microsoft.Compute/virtualMachines@2023-03-01",
			right:        []string{"Microsoft.Compute/virtualMachines", "Microsoft.Network/virtualNetworks"},
			allow:        true,
		},
		{
			desc:         "in_negative",
			resourceType: "\"Microsoft.DocumentDB/databaseAccounts@2024-12-01-preview\"",
			right:        []string{"Microsoft.Compute/virtualMachines", "Microsoft.Network/virtualNetworks"},
			allow:        false,
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			ctx := shared.NewContext()
			sut := In{
				BaseCondition: BaseCondition{
					Subject: value.NewFieldValue("type", ctx),
				},
				Values: c.right,
			}
			actual, err := sut.Rego(ctx)
			require.NoError(t, err)
			cfg := shared.WithUtilFunctions("r := input\n" + actual)
			shared.AssertRegoAllow(t, cfg, map[string]any{"values": map[string]any{
				"type": c.resourceType,
			}}, c.allow, ctx)
		})
	}
}
