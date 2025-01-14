package pkg

import (
	"fmt"
	"json-rule-finder/pkg/shared"
	"testing"

	"github.com/open-policy-agent/opa/format"
	"github.com/open-policy-agent/opa/rego"
	"github.com/stretchr/testify/require"
)

func TestNotOperator(t *testing.T) {
	cases := []struct {
		desc      string
		condition shared.Rego
		protocol  string
		allowed   bool
	}{
		{
			desc: "alllow",
			condition: &EqualsCondition{
				BaseCondition: BaseCondition{
					Subject: shared.StringRego(`r.change.after.protocols[x]`),
				},
				Value: "http",
			},
			protocol: "tcp",
			allowed:  true,
		},
		{
			desc: "disallow",
			condition: &EqualsCondition{
				BaseCondition: BaseCondition{
					Subject: shared.StringRego(`r.change.after.protocols[x]`),
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
				Conditions: []shared.Rego{
					&AnyOf{
						Conditions: []shared.Rego{
							&EqualsCondition{
								BaseCondition: BaseCondition{
									Subject: shared.StringRego(`r.change.after.protocols[x]`),
								},
								Value: "http",
							},
							&EqualsCondition{
								BaseCondition: BaseCondition{
									Subject: shared.StringRego(`r.change.after.protocols[x]`),
								},
								Value: "https",
							},
						},
						ConditionSetName: "condition_any_of_1",
					},
					&EqualsCondition{
						BaseCondition: BaseCondition{
							Subject: shared.StringRego(`r.change.after.protocols[x]`),
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
				Conditions: []shared.Rego{
					&AnyOf{
						Conditions: []shared.Rego{
							&EqualsCondition{
								BaseCondition: BaseCondition{
									Subject: shared.StringRego(`r.change.after.protocols[x]`),
								},
								Value: "http",
							},
							&EqualsCondition{
								BaseCondition: BaseCondition{
									Subject: shared.StringRego(`r.change.after.protocols[x]`),
								},
								Value: "https",
							},
						},
						ConditionSetName: "condition_any_of",
					},
					&EqualsCondition{
						BaseCondition: BaseCondition{
							Subject: shared.StringRego(`r.change.after.protocols[x]`),
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
			assertRegoAllow(t, regoCfg, func() *rego.EvalOption {
				input := rego.EvalInput(map[string]any{
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
				})
				return &input
			}(), c.allowed, ctx)
		})
	}
}
