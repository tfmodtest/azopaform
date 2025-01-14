package pkg

import (
	"context"
	"fmt"
	"json-rule-finder/pkg/shared"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEqualsCondition(t *testing.T) {
	cases := []struct {
		desc  string
		left  shared.Rego
		right any
		setup func(ctx context.Context)
		allow bool
	}{
		{
			desc:  "string",
			left:  shared.StringRego(`"right"`),
			right: "right",
			allow: true,
		},
		{
			desc:  "string_negative",
			left:  shared.StringRego(`"left"`),
			right: "right",
			allow: false,
		},
		{
			desc:  "int",
			left:  shared.StringRego("1"),
			right: 1,
			allow: true,
		},
		{
			desc:  "int_negative",
			left:  shared.StringRego("1"),
			right: 2,
			allow: false,
		},
		{
			desc:  "bool",
			left:  shared.StringRego("true"),
			right: true,
			allow: true,
		},
		{
			desc:  "bool_negative",
			left:  shared.StringRego("false"),
			right: true,
			allow: false,
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			ctx := NewContext()
			if c.setup != nil {
				c.setup(ctx)
			}
			sut := EqualsCondition{
				BaseCondition: BaseCondition{
					Subject: c.left,
				},
				Value: c.right,
			}
			actual, err := sut.Rego(ctx)
			require.NoError(t, err)
			cfg := fmt.Sprintf(conditionRegoTemplate, actual)
			shared.AssertRegoAllow(t, cfg, nil, c.allow, ctx)
		})
	}
}
