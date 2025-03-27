package pkg

import (
	"fmt"
	"json-rule-finder/pkg/condition"
	"json-rule-finder/pkg/shared"
	"testing"

	"github.com/open-policy-agent/opa/format"
	"github.com/open-policy-agent/opa/rego"
	"github.com/stretchr/testify/require"
)

const testRegoModuleTemplate = `
	package main
	
	import rego.v1
	
	default allow := false
	r := input.resource_changes[_]
	
	allow if condition0
	
	%s
	`

func TestAnyOfOperator(t *testing.T) {
	cases := []struct {
		desc       string
		conditions []shared.Rego
		protocol   string
		allowed    bool
	}{
		{
			desc: "alllow left",
			conditions: []shared.Rego{
				&condition.Equals{
					BaseCondition: condition.BaseCondition{
						Subject: shared.StringRego(`r.change.after.protocols[x]`),
					},
					Value: "http",
				},
				&condition.Equals{
					BaseCondition: condition.BaseCondition{
						Subject: shared.StringRego(`r.change.after.protocols[x]`),
					},
					Value: "ws",
				},
			},
			protocol: "http",
			allowed:  true,
		},
		{
			desc: "alllow right",
			conditions: []shared.Rego{
				&condition.Equals{
					BaseCondition: condition.BaseCondition{
						Subject: shared.StringRego(`r.change.after.protocols[x]`),
					},
					Value: "http",
				},
				&condition.Equals{
					BaseCondition: condition.BaseCondition{
						Subject: shared.StringRego(`r.change.after.protocols[x]`),
					},
					Value: "ws",
				},
			},
			protocol: "ws",
			allowed:  true,
		},
		{
			desc: "nested operator",
			conditions: []shared.Rego{
				&AnyOf{
					Conditions: []shared.Rego{
						&condition.Equals{
							BaseCondition: condition.BaseCondition{
								Subject: shared.StringRego(`r.change.after.protocols[x]`),
							},
							Value: "http",
						},
						&condition.Equals{
							BaseCondition: condition.BaseCondition{
								Subject: shared.StringRego(`r.change.after.protocols[x]`),
							},
							Value: "https",
						},
					},
					baseOperation: baseOperation{
						conditionSetName: "condition1",
					},
				},
				&condition.Equals{
					BaseCondition: condition.BaseCondition{
						Subject: shared.StringRego(`r.change.after.protocols[x]`),
					},
					Value: "ws",
				},
			},
			protocol: "https",
			allowed:  true,
		},
		{
			desc: "disallow",
			conditions: []shared.Rego{
				&condition.Equals{
					BaseCondition: condition.BaseCondition{
						Subject: shared.StringRego(`r.change.after.protocols[x]`),
					},
					Value: "http",
				},
				&condition.Equals{
					BaseCondition: condition.BaseCondition{
						Subject: shared.StringRego(`r.change.after.protocols[x]`),
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
				Conditions: c.conditions,
				baseOperation: baseOperation{
					conditionSetName: "condition0",
				},
			}
			ctx := shared.NewContext()
			actual, err := sut.Rego(ctx)
			require.NoError(t, err)
			regoCfg := fmt.Sprintf(testRegoModuleTemplate, actual)
			formattedCfg, err := format.Source("test.rego", []byte(regoCfg))
			require.NoError(t, err)
			regoCfg = string(formattedCfg)
			shared.AssertRegoAllow(t, regoCfg, func() *rego.EvalOption {
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
