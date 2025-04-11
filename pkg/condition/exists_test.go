package condition

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"json-rule-finder/pkg/shared"
	"testing"
)

func TestExistsCondition(t *testing.T) {
	cases := []struct {
		desc  string
		left  shared.Rego
		right any
		allow bool
	}{
		{
			desc:  "exists",
			left:  shared.StringRego(`input.body.sku`),
			right: "true",
			allow: true,
		},
		{
			desc:  "exists_bool",
			left:  shared.StringRego(`input.body.sku`),
			right: true,
			allow: true,
		},
		{
			desc:  "not_exists",
			left:  shared.StringRego(`input.body.not_exist`),
			right: "true",
			allow: false,
		},
		{
			desc:  "not_exists_bool",
			left:  shared.StringRego(`input.body.not_exist`),
			right: true,
			allow: false,
		},
		{
			desc:  "exists_negative",
			left:  shared.StringRego(`input.body.sku`),
			right: "false",
			allow: false,
		},
		{
			desc:  "exists_negative_bool",
			left:  shared.StringRego(`input.body.sku`),
			right: false,
			allow: false,
		},
		{
			desc:  "not_exists_negative",
			left:  shared.StringRego(`input.body.not_exist`),
			right: "false",
			allow: true,
		},
		{
			desc:  "not_exists_negative_bool",
			left:  shared.StringRego(`input.body.not_exist`),
			right: false,
			allow: true,
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			ctx := shared.NewContext()
			sut := Exists{
				BaseCondition: BaseCondition{
					Subject: c.left,
				},
				Value: c.right,
			}
			actual, err := sut.Rego(ctx)
			require.NoError(t, err)
			cfg := fmt.Sprintf(shared.RegoTestTemplate, actual)
			shared.AssertRegoAllow(t, cfg, map[string]any{
				"body": map[string]any{
					"sku": map[string]any{
						"name": "GP_Gen5_2",
						"tier": "GP_Gen5",
					},
				},
			}, c.allow, ctx)
		})
	}
}
