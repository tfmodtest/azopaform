package operation

import (
	"fmt"
	"json-rule-finder/pkg/condition"
	"json-rule-finder/pkg/shared"
	"testing"

	"github.com/open-policy-agent/opa/format"
	"github.com/stretchr/testify/require"
)

func TestAllOfOperator(t *testing.T) {
	cases := []struct {
		desc             string
		conditions       []shared.Rego
		protocol         string
		port             int
		publicAccessible bool
		allowed          bool
	}{
		{
			desc: "alllow",
			conditions: []shared.Rego{
				&condition.Equals{
					BaseCondition: condition.BaseCondition{
						Subject: shared.StringRego(`r.values.protocols[x]`),
					},
					Value: "tcp",
				},
				&condition.Equals{
					BaseCondition: condition.BaseCondition{
						Subject: shared.StringRego(`r.values.port`),
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
			conditions: []shared.Rego{
				&AnyOf{
					Conditions: []shared.Rego{
						&condition.Equals{
							BaseCondition: condition.BaseCondition{
								Subject: shared.StringRego(`r.values.protocols[x]`),
							},
							Value: "tcp",
						},
						&condition.Equals{
							BaseCondition: condition.BaseCondition{
								Subject: shared.StringRego(`r.values.port`),
							},
							Value: 22,
						},
					},
					baseOperation: baseOperation{
						helperFunctionName: "condition1",
					},
				},
				&condition.Equals{
					BaseCondition: condition.BaseCondition{
						Subject: shared.StringRego(`r.values.public_accessible`),
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
			conditions: []shared.Rego{
				&condition.Equals{
					BaseCondition: condition.BaseCondition{
						Subject: shared.StringRego(`r.values.protocols[x]`),
					},
					Value: "tcp",
				},
				&condition.Equals{
					BaseCondition: condition.BaseCondition{
						Subject: shared.StringRego(`r.values.port`),
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
								"protocols":         []string{c.protocol},
								"port":              c.port,
								"public_accessible": c.publicAccessible,
							},
						},
					},
				},
			}, c.allowed, ctx)
		})
	}
}
