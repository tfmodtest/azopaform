package pkg

import (
	"fmt"
	"testing"

	"github.com/open-policy-agent/opa/format"
	"github.com/open-policy-agent/opa/rego"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnyOfOperator(t *testing.T) {
	template := `
	package main
	
	import rego.v1
	
	default allow := false
	r := input.resource_changes[_]
	
	allow if condition0
	
	%s
	`
	cases := []struct {
		desc       string
		conditions []Rego
		protocol   string
		allowed    bool
	}{
		{
			desc: "alllow left",
			conditions: []Rego{
				&EqualsCondition{
					condition: condition{
						Subject: stringRego(`r.change.after.protocols[x]`),
					},
					Value: "http",
				},
				&EqualsCondition{
					condition: condition{
						Subject: stringRego(`r.change.after.protocols[x]`),
					},
					Value: "ws",
				},
			},
			protocol: "http",
			allowed:  true,
		},
		{
			desc: "alllow right",
			conditions: []Rego{
				&EqualsCondition{
					condition: condition{
						Subject: stringRego(`r.change.after.protocols[x]`),
					},
					Value: "http",
				},
				&EqualsCondition{
					condition: condition{
						Subject: stringRego(`r.change.after.protocols[x]`),
					},
					Value: "ws",
				},
			},
			protocol: "ws",
			allowed:  true,
		},
		{
			desc: "nested operator",
			conditions: []Rego{
				&AnyOf{
					Conditions: []Rego{
						&EqualsCondition{
							condition: condition{
								Subject: stringRego(`r.change.after.protocols[x]`),
							},
							Value: "http",
						},
						&EqualsCondition{
							condition: condition{
								Subject: stringRego(`r.change.after.protocols[x]`),
							},
							Value: "https",
						},
					},
					ConditionSetName: "condition1",
				},
				&EqualsCondition{
					condition: condition{
						Subject: stringRego(`r.change.after.protocols[x]`),
					},
					Value: "ws",
				},
			},
			protocol: "https",
			allowed:  true,
		},
		{
			desc: "disallow",
			conditions: []Rego{
				&EqualsCondition{
					condition: condition{
						Subject: stringRego(`r.change.after.protocols[x]`),
					},
					Value: "http",
				},
				&EqualsCondition{
					condition: condition{
						Subject: stringRego(`r.change.after.protocols[x]`),
					},
					Value: "ws",
				},
			},
			protocol: "tcp",
			allowed:  false,
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			sut := &AnyOf{
				Conditions:       c.conditions,
				ConditionSetName: "condition0",
			}
			ctx := NewContext()
			actual, err := sut.Rego(ctx)
			require.NoError(t, err)
			regoCfg := fmt.Sprintf(template, actual)
			formattedCfg, err := format.Source("test.rego", []byte(regoCfg))
			require.NoError(t, err)
			regoCfg = string(formattedCfg)
			query, err := rego.New(rego.Query("data.main.allow"), rego.Module("test.rego", regoCfg)).PrepareForEval(ctx)
			require.NoError(t, err)
			eval, err := query.Eval(ctx, rego.EvalInput(map[string]any{
				"resource_changes": []map[string]any{
					{
						"type": "azapi_resource",
						"change": map[string]any{
							"after": map[string]any{
								"protocols": []string{c.protocol},
							},
						},
					},
				},
			}))
			require.NoError(t, err)
			assert.Equal(t, c.allowed, eval.Allowed())
		})
	}
}
