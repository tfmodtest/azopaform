package operation

import (
	"fmt"
	"json-rule-finder/pkg/condition"
	"json-rule-finder/pkg/shared"
	"testing"

	"github.com/open-policy-agent/opa/format"
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
			condition: &condition.Equals{
				BaseCondition: condition.BaseCondition{
					Subject: shared.StringRego(`r.values.protocols[x]`),
				},
				Value: "http",
			},
			protocol: "tcp",
			allowed:  true,
		},
		{
			desc: "disallow",
			condition: &condition.Equals{
				BaseCondition: condition.BaseCondition{
					Subject: shared.StringRego(`r.values.protocols[x]`),
				},
				Value: "http",
			},
			protocol: "http",
			allowed:  false,
		},
		{
			desc: "nested operator disallow",
			condition: &AnyOf{
				baseOperation: baseOperation{
					helperFunctionName: "condition_any_of_0",
				},
				Conditions: []shared.Rego{
					&AnyOf{
						Conditions: []shared.Rego{
							&condition.Equals{
								BaseCondition: condition.BaseCondition{
									Subject: shared.StringRego(`r.values.protocols[x]`),
								},
								Value: "http",
							},
							&condition.Equals{
								BaseCondition: condition.BaseCondition{
									Subject: shared.StringRego(`r.values.protocols[x]`),
								},
								Value: "https",
							},
						},
						baseOperation: baseOperation{
							helperFunctionName: "condition_any_of_1",
						},
					},
					&condition.Equals{
						BaseCondition: condition.BaseCondition{
							Subject: shared.StringRego(`r.values.protocols[x]`),
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
				baseOperation: baseOperation{
					helperFunctionName: "condition_all_of",
				},
				Conditions: []shared.Rego{
					&AnyOf{
						Conditions: []shared.Rego{
							&condition.Equals{
								BaseCondition: condition.BaseCondition{
									Subject: shared.StringRego(`r.values.protocols[x]`),
								},
								Value: "http",
							},
							&condition.Equals{
								BaseCondition: condition.BaseCondition{
									Subject: shared.StringRego(`r.values.protocols[x]`),
								},
								Value: "https",
							},
						},
						baseOperation: baseOperation{
							helperFunctionName: "condition_any_of",
						},
					},
					&condition.Equals{
						BaseCondition: condition.BaseCondition{
							Subject: shared.StringRego(`r.values.protocols[x]`),
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
			sut := &Not{
				Body: c.condition,
				baseOperation: baseOperation{
					helperFunctionName: "condition0",
				},
			}
			ctx := shared.NewContext()
			actual, err := sut.Rego(ctx)
			require.NoError(t, err)
			regoCfg := fmt.Sprintf(testRegoModuleTemplate, shared.UTILS_REGO, actual, ctx.HelperFunctionsRego())
			formattedCfg, err := format.Source("test.rego", []byte(regoCfg))
			require.NoError(t, err)
			regoCfg = string(formattedCfg)
			shared.AssertRegoAllow(t, regoCfg, map[string]any{
				"resource_changes": []map[string]any{
					{
						"mode":    "managed",
						"address": "azapi_resource.this",
						"type":    "azapi_resource",
						"change": map[string]any{
							"after": map[string]any{
								"protocols": []string{c.protocol},
							},
						},
					},
				},
			}, c.allowed, ctx)
		})
	}
}
