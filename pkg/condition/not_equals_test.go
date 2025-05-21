package condition

import (
	"fmt"
	"github.com/stretchr/testify/assert"
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

func TestNotEqualsCondition_Parameter(t *testing.T) {
	ctx := shared.NewContext()
	ctx.GetParameterFunc = func(key string) (any, bool, error) {
		if key == "param1" {
			return "value1", true, nil
		}
		return nil, false, nil
	}
	rhs := shared.StringRego("[parameters('param1')]")
	subject, _ := NewFieldValue("Microsoft.Storage/storageAccounts/networkAcls.bypass", ctx)
	sut := NotEquals{
		BaseCondition: BaseCondition{
			Subject: subject,
		},
		Value: rhs,
	}
	actual, err := sut.Rego(ctx)
	require.NoError(t, err)
	assert.NotEqual(t, "", actual)
}
