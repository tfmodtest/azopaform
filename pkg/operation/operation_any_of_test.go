package operation

import (
	"fmt"
	"testing"

	"github.com/open-policy-agent/opa/v1/format"
	"github.com/stretchr/testify/require"
	"github.com/tfmodtest/azopaform/pkg/condition"
	"github.com/tfmodtest/azopaform/pkg/shared"
)

const testRegoModuleTemplate = `
	package main
	
	import rego.v1

    %s
	
	default allow := false
	
	allow if condition0
	
	%s

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
						Subject: shared.StringRego(`r.values.protocols[x]`),
					},
					Value: "http",
				},
				&condition.Equals{
					BaseCondition: condition.BaseCondition{
						Subject: shared.StringRego(`r.values.protocols[x]`),
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
						Subject: shared.StringRego(`r.values.protocols[x]`),
					},
					Value: "http",
				},
				&condition.Equals{
					BaseCondition: condition.BaseCondition{
						Subject: shared.StringRego(`r.values.protocols[x]`),
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
						helperFunctionName: "condition1",
					},
				},
				&condition.Equals{
					BaseCondition: condition.BaseCondition{
						Subject: shared.StringRego(`r.values.protocols[x]`),
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
						Subject: shared.StringRego(`r.values.protocols[x]`),
					},
					Value: "http",
				},
				&condition.Equals{
					BaseCondition: condition.BaseCondition{
						Subject: shared.StringRego(`r.values.protocols[x]`),
					},
					Value: "ws",
				},
			},
			protocol: "tcp",
			allowed:  false,
		},
		{
			desc: "condition mixed with operation",
			conditions: []shared.Rego{
				&condition.Equals{
					BaseCondition: condition.BaseCondition{
						Subject: shared.StringRego(`r.values.protocols[x]`),
					},
					Value: "http",
				},
				&AllOf{
					baseOperation: newBaseOperation(),
					Conditions: []shared.Rego{
						&condition.Equals{
							BaseCondition: condition.BaseCondition{
								Subject: shared.StringRego(`r.values.name`),
							},
							Value: "test",
						},
						&condition.Equals{
							BaseCondition: condition.BaseCondition{
								Subject: shared.StringRego(`r.values.location`),
							},
							Value: "eastus",
						},
					},
				},
			},
			protocol: "tcp",
			allowed:  true,
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			sut := &AnyOf{
				Conditions: c.conditions,
				baseOperation: baseOperation{
					helperFunctionName: "condition0",
				},
			}
			ctx := shared.NewContext()
			actual, err := sut.Rego(ctx)
			require.NoError(t, err)
			regoCfg := fmt.Sprintf(testRegoModuleTemplate, shared.UtilsRego, actual, ctx.HelperFunctionsRego())
			formattedCfg, err := format.Source("test.rego", []byte(regoCfg))
			require.NoError(t, err)
			regoCfg = string(formattedCfg)
			shared.AssertRegoAllow(t, regoCfg, map[string]any{
				"resource_changes": []map[string]any{
					{
						"address": "azapi_resource.this",
						"mode":    "managed",
						"type":    "azapi_resource",
						"change": map[string]any{
							"after": map[string]any{
								"protocols": []string{c.protocol},
								"name":      "test",
								"location":  "eastus",
							},
						},
					},
				},
			}, c.allowed, ctx)
		})
	}
}
