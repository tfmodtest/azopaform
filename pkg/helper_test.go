package pkg

import (
	"fmt"
	"testing"

	"github.com/emirpasic/gods/stacks"
	"github.com/open-policy-agent/opa/rego"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReplaceIndex_WithWildcardIndex(t *testing.T) {
	cases := []struct {
		desc    string
		subject string
		value   string
		allowed bool
	}{
		{
			desc:    "map wildcard allowed",
			subject: `input.protocols[*]`,
			value:   "http",
			allowed: true,
		},
		{
			desc:    "map wildcard disallowed",
			subject: `input.protocols[*]`,
			value:   "tcp",
			allowed: false,
		},
		{
			desc:    "array wildcard allowed",
			subject: `input.network_acl[*].ports[*]`,
			value:   "80",
			allowed: true,
		},
		{
			desc:    "array wildcard disallowed",
			subject: `input.network_acl[*].ports[*]`,
			value:   "443",
			allowed: false,
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			sut := EqualsCondition{
				condition: condition{
					Subject: stringRego(c.subject),
				},
				Value: c.value,
			}
			ctx := NewContext()
			stack := ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"]
			stack.Push("x")
			assertion, err := sut.Rego(ctx)
			require.NoError(t, err)
			regoCode := fmt.Sprintf(`package main

import rego.v1

default allow := false
allow if %s`, assertion)
			eval, err := rego.New(rego.Query("data.main.allow"), rego.Module("test.rego", regoCode)).PrepareForEval(ctx)
			require.NoError(t, err)
			var result rego.ResultSet
			i := map[string]any{
				"protocols": []string{"http", "https"},
				"network_acl": map[string]any{

					"tcp": map[string]any{
						"ports": []string{"22"},
					},
					"http": map[string]any{
						"ports": []string{"80"},
					},
				},
			}
			input := rego.EvalInput(i)
			result, err = eval.Eval(ctx, input)

			require.NoError(t, err)
			assert.Equal(t, c.allowed, result.Allowed())
		})
	}
}

func p[T any](v T) *T {
	return &v
}
