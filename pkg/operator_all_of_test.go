package pkg

import (
	"fmt"
	"github.com/open-policy-agent/opa/format"
	"testing"

	"github.com/open-policy-agent/opa/rego"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAllOfOperator(t *testing.T) {
	cases := []struct {
		desc             string
		conditions       []Rego
		protocol         string
		port             int
		publicAccessible bool
		allowed          bool
	}{
		{
			desc: "alllow",
			conditions: []Rego{
				&EqualsCondition{
					condition: condition{
						Subject: stringRego(`r.change.after.protocols[x]`),
					},
					Value: "tcp",
				},
				&EqualsCondition{
					condition: condition{
						Subject: stringRego(`r.change.after.port`),
					},
					Value: 22,
				},
			},
			protocol: "tcp",
			port:     22,
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
							Value: "tcp",
						},
						&EqualsCondition{
							condition: condition{
								Subject: stringRego(`r.change.after.port`),
							},
							Value: 22,
						},
					},
					ConditionSetName: "condition1",
				},
				&EqualsCondition{
					condition: condition{
						Subject: stringRego(`r.change.after.public_accessible`),
					},
					Value: false,
				},
			},
			protocol:         "https",
			publicAccessible: false,
			port:             22,
			allowed:          true,
		},
		{
			desc: "disallow",
			conditions: []Rego{
				&EqualsCondition{
					condition: condition{
						Subject: stringRego(`r.change.after.protocols[x]`),
					},
					Value: "tcp",
				},
				&EqualsCondition{
					condition: condition{
						Subject: stringRego(`r.change.after.port`),
					},
					Value: 22,
				},
			},
			protocol: "http",
			port:     22,
			allowed:  false,
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			sut := &AllOf{
				Conditions:       c.conditions,
				ConditionSetName: "condition0",
			}
			ctx := NewContext()
			actual, err := sut.Rego(ctx)
			require.NoError(t, err)
			regoCfg := fmt.Sprintf(testRegoModuleTemplate, actual)
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
								"protocols":         []string{c.protocol},
								"port":              c.port,
								"public_accessible": c.publicAccessible,
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
