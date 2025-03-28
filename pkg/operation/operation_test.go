package operation

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"json-rule-finder/pkg/condition"
	"json-rule-finder/pkg/shared"
	"json-rule-finder/pkg/value"
	"testing"
)

func TestOperations(t *testing.T) {
	tests := []struct {
		name      string
		operation shared.Rego
		expected  string
	}{
		{
			name: "NestedWhereOperator",
			operation: WhereOperator{
				Condition: AllOf{
					Conditions: []shared.Rego{
						condition.Equals{
							BaseCondition: condition.BaseCondition{
								Subject: &value.FieldValue{Name: "Microsoft.Web/serverFarms/sku.tier"},
							},
							Value: "Standard",
						},
						condition.Exists{
							BaseCondition: condition.BaseCondition{
								Subject: &value.FieldValue{Name: "Microsoft.Web/serverFarms/sku.tier"},
							},
							Value: true,
						},
						condition.Equals{
							BaseCondition: condition.BaseCondition{
								Subject: &value.FieldValue{Name: "Microsoft.Web/serverFarms/sku.size"},
							},
							Value: "P1v3",
						},
					},
					baseOperation: baseOperation{
						conditionSetName: "aaaaa",
					},
				},
				ConditionSetName: "aaaaaaaaa",
			},
			expected: "aaaaaaaaa(x) if {\naaaaa(x)\n}\naaaaa(x) if {\nr.change.after.sku[x].tier == \"Standard\"\nr.change.after.sku_name\nr.change.after.sku[x].size == \"P1v3\"\n}",
		},
		{
			name: "WhereOperator",
			operation: WhereOperator{
				Condition: condition.Equals{
					BaseCondition: condition.BaseCondition{
						Subject: &value.FieldValue{Name: "type"},
					},
					Value: "azurerm_app_service_plan",
				},
				ConditionSetName: "aaaaaaaaa",
			},
			expected: "aaaaaaaaa(x) if {\nr.type == \"azurerm_app_service_plan\"\n}",
		},
		{
			name: "NestedAllOfOperator",
			operation: AllOf{
				Conditions: []shared.Rego{
					AllOf{
						Conditions: []shared.Rego{
							condition.Equals{
								BaseCondition: condition.BaseCondition{
									Subject: &value.FieldValue{Name: "type"},
								},
								Value: "azurerm_app_service_plan",
							},
							condition.Exists{
								BaseCondition: condition.BaseCondition{
									Subject: &value.FieldValue{Name: "Microsoft.Web/serverFarms/sku.name"},
								},
								Value: true,
							},
						},
						baseOperation: baseOperation{
							conditionSetName: "aaaaa",
						},
					},
					AnyOf{
						Conditions: []shared.Rego{
							condition.Equals{
								BaseCondition: condition.BaseCondition{
									Subject: &value.FieldValue{Name: "Microsoft.Web/serverFarms/sku.tier"},
								},
								Value: "Standard",
							},
							condition.Equals{
								BaseCondition: condition.BaseCondition{
									Subject: &value.FieldValue{Name: "Microsoft.Web/serverFarms/sku.tier"},
								},
								Value: "Basic",
							},
						},
						baseOperation: baseOperation{
							conditionSetName: "aaaaaaa",
						},
					},
				},
				baseOperation: baseOperation{
					conditionSetName: "aaaaa",
				},
			},
			expected: "aaaaa if {\naaaaa\nnot aaaaaaa\n}\naaaaa if {\nr.type == \"azurerm_app_service_plan\"\nr.change.after.sku_name\n}\naaaaaaa if {\nr.change.after.sku[0].tier != \"Standard\"\nr.change.after.sku[0].tier != \"Basic\"\n}",
		},
		{
			name: "NestedAnyOfOperator",
			operation: AnyOf{
				Conditions: []shared.Rego{
					AnyOf{
						Conditions: []shared.Rego{
							condition.Equals{
								BaseCondition: condition.BaseCondition{
									Subject: &value.FieldValue{Name: "type"},
								},
								Value: "azurerm_app_service_plan",
							},
							condition.Equals{
								BaseCondition: condition.BaseCondition{
									Subject: &value.FieldValue{Name: "type"},
								},
								Value: "azurerm_app_service_environment",
							},
						},
						baseOperation: baseOperation{
							conditionSetName: "aaaaaaa",
						},
					},
					AnyOf{
						Conditions: []shared.Rego{
							condition.Equals{
								BaseCondition: condition.BaseCondition{
									Subject: &value.FieldValue{Name: "Microsoft.Web/serverFarms/sku.tier"},
								},
								Value: "Standard",
							},
							condition.Equals{
								BaseCondition: condition.BaseCondition{
									Subject: &value.FieldValue{Name: "Microsoft.Web/serverFarms/sku.tier"},
								},
								Value: "Basic",
							},
						},
						baseOperation: baseOperation{
							conditionSetName: "aaaaaaa",
						},
					},
				},
				baseOperation: baseOperation{
					conditionSetName: "aaaaaaa",
				},
			},
			expected: "aaaaaaa if {\naaaaaaa\naaaaaaa\n}\naaaaaaa if {\nr.type != \"azurerm_app_service_plan\"\nr.type != \"azurerm_app_service_environment\"\n}\naaaaaaa if {\nr.change.after.sku[0].tier != \"Standard\"\nr.change.after.sku[0].tier != \"Basic\"\n}",
		},
		{
			name: "AllOfOperator",
			operation: AllOf{
				Conditions: []shared.Rego{
					condition.Equals{
						BaseCondition: condition.BaseCondition{
							Subject: &value.FieldValue{Name: "type"},
						},
						Value: "azurerm_app_service_plan",
					},
					condition.Exists{
						BaseCondition: condition.BaseCondition{
							Subject: &value.FieldValue{Name: "Microsoft.Web/serverFarms/sku.name"},
						},
						Value: true,
					},
				},
				baseOperation: baseOperation{
					conditionSetName: "aaaaa",
				},
			},
			expected: "aaaaa if {\nr.type == \"azurerm_app_service_plan\"\nr.change.after.sku_name\n}",
		},
		{
			name: "AnyOfOperator",
			operation: AnyOf{
				Conditions: []shared.Rego{
					condition.Equals{
						BaseCondition: condition.BaseCondition{
							Subject: &value.FieldValue{Name: "Microsoft.Web/serverFarms/sku.tier"},
						},
						Value: "Standard",
					},
					condition.In{
						BaseCondition: condition.BaseCondition{
							Subject: &value.FieldValue{Name: "Microsoft.Web/serverFarms/sku.tier"},
						},
						Values: []string{"Basic", "Premium"},
					},
				},
				baseOperation: baseOperation{
					conditionSetName: "aaaaaaa",
				},
			},
			expected: "aaaaaaa if {\nr.change.after.sku[0].tier != \"Standard\"\nnot r.change.after.sku[0].tier in [\"Basic\",\"Premium\"]\n}",
		},
		{
			name: "NotOperation",
			operation: Not{
				Body: condition.Equals{
					BaseCondition: condition.BaseCondition{
						Subject: value.FieldValue{
							Name: "Microsoft.Web/serverFarms/sku.tier",
						},
					},
					Value: "Standard",
				},
				baseOperation: baseOperation{
					conditionSetName: "aaa",
				},
			},
			expected: "aaa if {\nr.change.after.sku[0].tier == \"Standard\"\n}",
		},
		{
			name: "Equals",
			operation: condition.Equals{
				BaseCondition: condition.BaseCondition{
					Subject: &value.FieldValue{Name: "Microsoft.Web/serverFarms/sku.tier"},
				},
				Value: "Standard",
			},
			expected: "r.change.after.sku[0].tier == \"Standard\"",
		},
		{
			name: "NotEquals",
			operation: condition.NotEquals{
				BaseCondition: condition.BaseCondition{
					Subject: &value.FieldValue{Name: "Microsoft.Web/serverFarms/sku.tier"},
				},
				Value: "Standard",
			},
			expected: "r.change.after.sku[0].tier != \"Standard\"",
		},
		{
			name: "Like",
			operation: condition.Like{
				BaseCondition: condition.BaseCondition{
					Subject: &value.FieldValue{Name: "Microsoft.Web/serverFarms/sku.tier"},
				},
				Value: `^[^@]+@[^@]+\.[^@]+$`,
			},
			expected: "regex.match(\"^[^@]+@[^@]+\\.[^@]+$\",r.change.after.sku[0].tier)",
		},
		{
			name: "NotLike",
			operation: condition.NotLike{
				BaseCondition: condition.BaseCondition{
					Subject: &value.FieldValue{Name: "Microsoft.Web/serverFarms/sku.tier"},
				},
				Value: `^[^@]+@[^@]+\.[^@]+$`,
			},
			expected: "not regex.match(`^[^@]+@[^@]+\\.[^@]+$`,r.change.after.sku[0].tier)",
		},
		{
			name: "In",
			operation: condition.In{
				BaseCondition: condition.BaseCondition{
					Subject: &value.FieldValue{Name: "Microsoft.Web/serverFarms/sku.tier"},
				},
				Values: []string{"Basic", "Standard", "Premium"},
			},
			expected: "some r.change.after.sku[0].tier in [\"Basic\",\"Standard\",\"Premium\"]",
		},
		{
			name: "NotIn",
			operation: condition.NotIn{
				BaseCondition: condition.BaseCondition{
					Subject: &value.FieldValue{Name: "Microsoft.Web/serverFarms/sku.tier"},
				},
				Values: []string{"Basic", "Standard", "Premium"},
			},
			expected: "not r.change.after.sku[0].tier in [\"Basic\",\"Standard\",\"Premium\"]",
		},
	}
	t.Skipf("skipping")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := shared.NewContext()
			ctx.PushResourceType("Microsoft.Web/serverFarms")
			actual, err := tt.operation.Rego(ctx)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
