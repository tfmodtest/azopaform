package condition

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tfmodtest/azopaform/pkg/shared"
)

func TestNotEqualsCondition(t *testing.T) {
	cases := []struct {
		desc  string
		left  shared.Rego
		right any
		allow bool
	}{
		{
			desc:  "string_negative",
			left:  shared.StringRego(`"right"`),
			right: "right",
			allow: false,
		},
		{
			desc:  "string",
			left:  shared.StringRego(`"left"`),
			right: "right",
			allow: true,
		},
		{
			desc:  "int_negative",
			left:  shared.StringRego("1"),
			right: 1,
			allow: false,
		},
		{
			desc:  "int",
			left:  shared.StringRego("1"),
			right: 2,
			allow: true,
		},
		{
			desc:  "bool_negative",
			left:  shared.StringRego("true"),
			right: true,
			allow: false,
		},
		{
			desc:  "bool",
			left:  shared.StringRego("false"),
			right: true,
			allow: true,
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			ctx := shared.NewContext()
			sut := NotEquals{
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
