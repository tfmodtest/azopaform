package pkg

import (
	"context"
	"fmt"
	"testing"

	"github.com/open-policy-agent/opa/rego"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ Rego = stringRego("")

type stringRego string

func (s stringRego) Rego(ctx context.Context) (string, error) {
	return string(s), nil
}

func TestEqualsCondition(t *testing.T) {
	cases := []struct {
		desc  string
		left  Rego
		right any
		setup func(ctx context.Context)
		allow bool
	}{
		{
			desc:  "string",
			left:  stringRego(`"right"`),
			right: "right",
			allow: true,
		},
		{
			desc:  "string_negative",
			left:  stringRego(`"left"`),
			right: "right",
			allow: false,
		},
		{
			desc:  "int",
			left:  stringRego("1"),
			right: 1,
			allow: true,
		},
		{
			desc:  "int_negative",
			left:  stringRego("1"),
			right: 2,
			allow: false,
		},
		{
			desc:  "bool",
			left:  stringRego("true"),
			right: true,
			allow: true,
		},
		{
			desc:  "bool_negative",
			left:  stringRego("false"),
			right: true,
			allow: false,
		},
	}
	template := `package main

import rego.v1

default allow := false
allow if %s`
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			ctx := NewContext()
			if c.setup != nil {
				c.setup(ctx)
			}
			sut := EqualsCondition{
				condition: condition{
					Subject: c.left,
				},
				Value: c.right,
			}
			actual, err := sut.Rego(ctx)
			require.NoError(t, err)
			cfg := fmt.Sprintf(template, actual)
			eval := rego.New(rego.Query("data.main.allow"), rego.Module("test.rego", cfg))
			result, err := eval.Eval(ctx)
			require.NoError(t, err)
			assert.Equal(t, c.allow, result.Allowed())
		})
	}
}
