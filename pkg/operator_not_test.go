package pkg

import (
	"fmt"
	"testing"

	"github.com/open-policy-agent/opa/format"
	"github.com/open-policy-agent/opa/rego"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNotOperator(t *testing.T) {
	cases := []struct {
		desc      string
		condition Rego
		protocol  string
		allowed   bool
	}{
		{
			desc: "alllow",
			condition: &EqualsCondition{
				condition: condition{
					Subject: stringRego(`r.change.after.protocols[x]`),
				},
				Value: "http",
			},
			protocol: "tcp",
			allowed:  true,
		},
		{
			desc: "disallow",
			condition: &EqualsCondition{
				condition: condition{
					Subject: stringRego(`r.change.after.protocols[x]`),
				},
				Value: "http",
			},
			protocol: "http",
			allowed:  false,
		},
		{
			desc: "nested operator disallow",
			condition: &AnyOf{
				ConditionSetName: "condition_any_of_0",
				Conditions: []Rego{
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
						ConditionSetName: "condition_any_of_1",
					},
					&EqualsCondition{
						condition: condition{
							Subject: stringRego(`r.change.after.protocols[x]`),
						},
						Value: "ws",
					},
				},
			},
			protocol: "ws",
			allowed:  false,
		},
		{
			desc: "nested operator allow",
			condition: &AllOf{
				ConditionSetName: "condition_all_of",
				Conditions: []Rego{
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
						ConditionSetName: "condition_any_of",
					},
					&EqualsCondition{
						condition: condition{
							Subject: stringRego(`r.change.after.protocols[x]`),
						},
						Value: "ws",
					},
				},
			},
			protocol: "tcp",
			allowed:  true,
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			sut := &NotOperator{
				Body:             c.condition,
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
