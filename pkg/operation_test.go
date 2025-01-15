package pkg

import (
	"github.com/prashantv/gostub"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"json-rule-finder/pkg/condition"
	"json-rule-finder/pkg/shared"
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
				Conditions: []shared.Rego{
					AllOf{
						Conditions: []shared.Rego{
							condition.Equals{
								BaseCondition: condition.BaseCondition{
									Subject: OperationField("Microsoft.Web/serverFarms/sku.tier"),
								},
								Value: "Standard",
							},
							condition.Exists{
								BaseCondition: condition.BaseCondition{
									Subject: OperationField("Microsoft.Web/serverFarms/sku.name"),
								},
								Value: true,
							},
							condition.Equals{
								BaseCondition: condition.BaseCondition{
									Subject: OperationField("Microsoft.Web/serverFarms/sku.size"),
								},
								Value: "P1v3",
							},
						},
						baseOperator: baseOperator{
							conditionSetName: "aaaaa",
						},
					},
				},
				ConditionSetName: "aaaaaaaaa",
			},
			expected: "aaaaaaaaa(x) if {\naaaaa(x)\n}\naaaaa(x) if {\nr.change.after.sku[x].tier == \"Standard\"\nr.change.after.sku_name\nr.change.after.sku[x].size == \"P1v3\"\n}",
		},
		{
			name: "WhereOperator",
			operation: WhereOperator{
				Conditions: []shared.Rego{
					condition.Equals{
						BaseCondition: condition.BaseCondition{
							Subject: OperationField("type"),
						},
						Value: "azurerm_app_service_plan",
					},
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
									Subject: OperationField("type"),
								},
								Value: "azurerm_app_service_plan",
							},
							condition.Exists{
								BaseCondition: condition.BaseCondition{
									Subject: OperationField("Microsoft.Web/serverFarms/sku.name"),
								},
								Value: true,
							},
						},
						baseOperator: baseOperator{
							conditionSetName: "aaaaa",
						},
					},
					AnyOf{
						Conditions: []shared.Rego{
							condition.Equals{
								BaseCondition: condition.BaseCondition{
									Subject: OperationField("Microsoft.Web/serverFarms/sku.tier"),
								},
								Value: "Standard",
							},
							condition.Equals{
								BaseCondition: condition.BaseCondition{
									Subject: OperationField("Microsoft.Web/serverFarms/sku.tier"),
								},
								Value: "Basic",
							},
						},
						baseOperator: baseOperator{
							conditionSetName: "aaaaaaa",
						},
					},
				},
				baseOperator: baseOperator{
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
									Subject: OperationField("type"),
								},
								Value: "azurerm_app_service_plan",
							},
							condition.Equals{
								BaseCondition: condition.BaseCondition{
									Subject: OperationField("type"),
								},
								Value: "azurerm_app_service_environment",
							},
						},
						baseOperator: baseOperator{
							conditionSetName: "aaaaaaa",
						},
					},
					AnyOf{
						Conditions: []shared.Rego{
							condition.Equals{
								BaseCondition: condition.BaseCondition{
									Subject: OperationField("Microsoft.Web/serverFarms/sku.tier"),
								},
								Value: "Standard",
							},
							condition.Equals{
								BaseCondition: condition.BaseCondition{
									Subject: OperationField("Microsoft.Web/serverFarms/sku.tier"),
								},
								Value: "Basic",
							},
						},
						baseOperator: baseOperator{
							conditionSetName: "aaaaaaa",
						},
					},
				},
				baseOperator: baseOperator{
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
							Subject: OperationField("type"),
						},
						Value: "azurerm_app_service_plan",
					},
					condition.Exists{
						BaseCondition: condition.BaseCondition{
							Subject: OperationField("Microsoft.Web/serverFarms/sku.name"),
						},
						Value: true,
					},
				},
				baseOperator: baseOperator{
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
							Subject: OperationField("Microsoft.Web/serverFarms/sku.tier"),
						},
						Value: "Standard",
					},
					condition.In{
						BaseCondition: condition.BaseCondition{
							Subject: OperationField("Microsoft.Web/serverFarms/sku.tier"),
						},
						Values: []string{"Basic", "Premium"},
					},
				},
				baseOperator: baseOperator{
					conditionSetName: "aaaaaaa",
				},
			},
			expected: "aaaaaaa if {\nr.change.after.sku[0].tier != \"Standard\"\nnot r.change.after.sku[0].tier in [\"Basic\",\"Premium\"]\n}",
		},
		{
			name: "NotOperator",
			operation: NotOperator{
				Body: condition.Equals{
					BaseCondition: condition.BaseCondition{
						Subject: FieldValue{
							Name: "Microsoft.Web/serverFarms/sku.tier",
						},
					},
					Value: "Standard",
				},
				ConditionSetName: "aaa",
			},
			expected: "aaa if {\nr.change.after.sku[0].tier == \"Standard\"\n}",
		},
		{
			name: "Equals",
			operation: condition.Equals{
				BaseCondition: condition.BaseCondition{
					Subject: OperationField("Microsoft.Web/serverFarms/sku.tier"),
				},
				Value: "Standard",
			},
			expected: "r.change.after.sku[0].tier == \"Standard\"",
		},
		{
			name: "NotEquals",
			operation: condition.NotEquals{
				BaseCondition: condition.BaseCondition{
					Subject: OperationField("Microsoft.Web/serverFarms/sku.tier"),
				},
				Value: "Standard",
			},
			expected: "r.change.after.sku[0].tier != \"Standard\"",
		},
		{
			name: "Like",
			operation: condition.Like{
				BaseCondition: condition.BaseCondition{
					Subject: OperationField("Microsoft.Web/serverFarms/sku.tier"),
				},
				Value: `^[^@]+@[^@]+\.[^@]+$`,
			},
			expected: "regex.match(\"^[^@]+@[^@]+\\.[^@]+$\",r.change.after.sku[0].tier)",
		},
		{
			name: "NotLike",
			operation: condition.NotLike{
				BaseCondition: condition.BaseCondition{
					Subject: OperationField("Microsoft.Web/serverFarms/sku.tier"),
				},
				Value: `^[^@]+@[^@]+\.[^@]+$`,
			},
			expected: "not regex.match(`^[^@]+@[^@]+\\.[^@]+$`,r.change.after.sku[0].tier)",
		},
		{
			name: "In",
			operation: condition.In{
				BaseCondition: condition.BaseCondition{
					Subject: OperationField("Microsoft.Web/serverFarms/sku.tier"),
				},
				Values: []string{"Basic", "Standard", "Premium"},
			},
			expected: "some r.change.after.sku[0].tier in [\"Basic\",\"Standard\",\"Premium\"]",
		},
		{
			name: "NotIn",
			operation: condition.NotIn{
				BaseCondition: condition.BaseCondition{
					Subject: OperationField("Microsoft.Web/serverFarms/sku.tier"),
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

func TestNewPolicyRuleBody(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]any
		expected *PolicyRuleBody
	}{
		{
			name: "CountOperation",
			input: map[string]any{
				"count": map[string]any{
					"field": "Microsoft.Network/networkSecurityGroups/securityRules[*]",
					"where": map[string]any{
						"field":  "Microsoft.Network/networkSecurityGroups/securityRules[*].direction",
						"equals": "Inbound",
					},
				},
				"greater": 0,
			},
			expected: &PolicyRuleBody{
				IfBody: condition.Greater{
					BaseCondition: condition.BaseCondition{
						Subject: CountOperator{
							Where: WhereOperator{
								Conditions: []shared.Rego{
									condition.Equals{
										BaseCondition: condition.BaseCondition{
											Subject: FieldValue{
												Name: "Microsoft.Network/networkSecurityGroups/securityRules[x].direction",
											},
										},
										Value: "Inbound",
									},
								},
								ConditionSetName: "condition1",
							},
							CountExp: "count({x|Microsoft.Network/networkSecurityGroups/securityRules[x];condition1(x)})",
						},
					},
					Value: 0,
				},
			},
		},
		{
			name: "NotOperation",
			input: map[string]any{
				"not": map[string]any{
					"field":     "Microsoft.HealthcareApis/services/corsConfiguration.origins[*]",
					"notEquals": "*",
				},
			},
			expected: &PolicyRuleBody{
				IfBody: NotOperator{
					Body: condition.NotEquals{
						BaseCondition: condition.BaseCondition{
							Subject: FieldValue{
								Name: "Microsoft.HealthcareApis/services/corsConfiguration.origins[x]",
							},
						},
						Value: "*",
					},
					ConditionSetName: "condition1",
				},
			},
		},
		{
			name: "NestedAnyOfOperation",
			input: map[string]any{
				"anyof": []any{
					map[string]any{
						"anyof": []any{
							map[string]any{
								"field":  "Microsoft.Sql/servers/minimalTlsVersion",
								"exists": false,
							},
							map[string]any{
								"field": "Microsoft.Sql/servers/minimalTlsVersion",
								"less":  "1.2",
							},
						},
					},
					map[string]any{
						"allof": []any{
							map[string]any{
								"field":  "type",
								"equals": "Microsoft.Sql/servers",
							},
							map[string]any{
								"field":  "Microsoft.Sql/servers/minimalTlsVersion",
								"exists": true,
							},
						},
					},
				},
			},
			expected: &PolicyRuleBody{
				IfBody: AnyOf{
					Conditions: []shared.Rego{
						AnyOf{
							Conditions: []shared.Rego{
								condition.Exists{
									BaseCondition: condition.BaseCondition{
										Subject: FieldValue{
											Name: "Microsoft.Sql/servers/minimalTlsVersion",
										},
									},
									Value: false,
								},
								condition.Less{
									BaseCondition: condition.BaseCondition{
										Subject: FieldValue{
											Name: "Microsoft.Sql/servers/minimalTlsVersion",
										},
									},
									Value: "1.2",
								},
							},
							baseOperator: baseOperator{
								conditionSetName: "condition1",
							},
						},
						AllOf{
							Conditions: []shared.Rego{
								condition.Equals{
									BaseCondition: condition.BaseCondition{
										Subject: FieldValue{
											Name: "type",
										},
									},
									Value: "azurerm_mssql_server",
								},
								condition.Exists{
									BaseCondition: condition.BaseCondition{
										Subject: FieldValue{
											Name: "Microsoft.Sql/servers/minimalTlsVersion",
										},
									},
									Value: true,
								},
							},
							baseOperator: baseOperator{
								conditionSetName: "condition1",
							},
						},
					},
					baseOperator: baseOperator{
						conditionSetName: "condition1",
					},
				},
			},
		},
		{
			name: "AnyOfOperation",
			input: map[string]any{
				"anyof": []any{
					map[string]any{
						"field":  "Microsoft.Sql/servers/minimalTlsVersion",
						"exists": false,
					},
					map[string]any{
						"field": "Microsoft.Sql/servers/minimalTlsVersion",
						"less":  "1.2",
					},
				},
			},
			expected: &PolicyRuleBody{
				IfBody: AnyOf{
					Conditions: []shared.Rego{
						condition.Exists{
							BaseCondition: condition.BaseCondition{
								Subject: FieldValue{
									Name: "Microsoft.Sql/servers/minimalTlsVersion",
								},
							},
							Value: false,
						},
						condition.Less{
							BaseCondition: condition.BaseCondition{
								Subject: FieldValue{
									Name: "Microsoft.Sql/servers/minimalTlsVersion",
								},
							},
							Value: "1.2",
						},
					},
					baseOperator: baseOperator{
						conditionSetName: "condition1",
					},
				},
			},
		},
		{
			name: "AllOfOperation",
			input: map[string]any{
				"allof": []any{
					map[string]any{
						"field":  "type",
						"equals": "Microsoft.HealthcareApis/services",
					},
					map[string]any{
						"field":  "Microsoft.HealthcareApis/services/cosmosDbConfiguration.keyVaultKeyUri",
						"exists": false,
					},
				},
			},
			expected: &PolicyRuleBody{
				IfBody: AllOf{
					Conditions: []shared.Rego{
						condition.Equals{
							BaseCondition: condition.BaseCondition{
								Subject: FieldValue{
									Name: "type",
								},
							},
							Value: "azurerm_healthcare_service",
						},
						condition.Exists{
							BaseCondition: condition.BaseCondition{
								Subject: FieldValue{
									Name: "Microsoft.HealthcareApis/services/cosmosDbConfiguration.keyVaultKeyUri",
								},
							},
							Value: false,
						},
					},
					baseOperator: baseOperator{
						conditionSetName: "condition1",
					},
				},
			},
		},
		{
			name: "Equals",
			input: map[string]any{
				"field":  "type",
				"equals": "Microsoft.Web/serverFarms",
			},
			expected: &PolicyRuleBody{
				IfBody: condition.Equals{
					BaseCondition: condition.BaseCondition{
						Subject: OperationField("type"),
					},
					Value: "Microsoft.Web/serverFarms",
				},
			},
		},
		{
			name: "NotEquals",
			input: map[string]any{
				"field":     "type",
				"notEquals": "Microsoft.Web/serverFarms",
			},
			expected: &PolicyRuleBody{
				IfBody: condition.NotEquals{
					BaseCondition: condition.BaseCondition{
						Subject: OperationField("type"),
					},
					Value: "Microsoft.Web/serverFarms",
				},
			},
		},
		{
			name: "Like",
			input: map[string]any{
				"field": "type",
				"like":  "Microsoft.Web/serverFarms",
			},
			expected: &PolicyRuleBody{
				IfBody: condition.Like{
					BaseCondition: condition.BaseCondition{
						Subject: OperationField("type"),
					},
					Value: "Microsoft.Web/serverFarms",
				},
			},
		},
		{
			name: "NotLike",
			input: map[string]any{
				"field":   "type",
				"notLike": "Microsoft.Web/serverFarms",
			},
			expected: &PolicyRuleBody{
				IfBody: condition.NotLike{
					BaseCondition: condition.BaseCondition{
						Subject: OperationField("type"),
					},
					Value: "Microsoft.Web/serverFarms",
				},
			},
		},
		{
			name: "In",
			input: map[string]any{
				"field": "type",
				"in":    []any{"Microsoft.Web/serverFarms", "Microsoft.Compute/virtualMachines"},
			},
			expected: &PolicyRuleBody{
				IfBody: condition.In{
					BaseCondition: condition.BaseCondition{
						Subject: OperationField("type"),
					},
					Values: []string{"Microsoft.Web/serverFarms", "Microsoft.Compute/virtualMachines"},
				},
			},
		},
		{
			name: "NotIn",
			input: map[string]any{
				"field": "type",
				"notIn": []any{"Microsoft.Web/serverFarms", "Microsoft.Compute/virtualMachines"},
			},
			expected: &PolicyRuleBody{
				IfBody: condition.NotIn{
					BaseCondition: condition.BaseCondition{
						Subject: OperationField("type"),
					},
					Values: []string{"Microsoft.Web/serverFarms", "Microsoft.Compute/virtualMachines"},
				},
			},
		},
		{
			name: "Contains",
			input: map[string]any{
				"field":    "type",
				"contains": "Microsoft.Web/serverFarms",
			},
			expected: &PolicyRuleBody{
				IfBody: condition.Contains{
					BaseCondition: condition.BaseCondition{
						Subject: OperationField("type"),
					},
					Value: "Microsoft.Web/serverFarms",
				},
			},
		},
		{
			name: "NotContains",
			input: map[string]any{
				"field":       "type",
				"notContains": "Microsoft.Web/serverFarms",
			},
			expected: &PolicyRuleBody{
				IfBody: condition.NotContains{
					BaseCondition: condition.BaseCondition{
						Subject: OperationField("type"),
					},
					Value: "Microsoft.Web/serverFarms",
				},
			},
		},
		{
			name: "ContainsKey",
			input: map[string]any{
				"field":       "type",
				"containsKey": "Microsoft.Web/serverFarms",
			},
			expected: &PolicyRuleBody{
				IfBody: condition.ContainsKey{
					BaseCondition: condition.BaseCondition{
						Subject: OperationField("type"),
					},
					KeyName: "Microsoft.Web/serverFarms",
				},
			},
		},
		{
			name: "NotContainsKey",
			input: map[string]any{
				"field":          "type",
				"notContainsKey": "Microsoft.Web/serverFarms",
			},
			expected: &PolicyRuleBody{
				IfBody: condition.NotContainsKey{
					BaseCondition: condition.BaseCondition{
						Subject: OperationField("type"),
					},
					KeyName: "Microsoft.Web/serverFarms",
				},
			},
		},
		{
			name: "Less",
			input: map[string]any{
				"field": "type",
				"less":  10,
			},
			expected: &PolicyRuleBody{
				IfBody: condition.Less{
					BaseCondition: condition.BaseCondition{
						Subject: OperationField("type"),
					},
					Value: 10,
				},
			},
		},
		{
			name: "LessOrEquals",
			input: map[string]any{
				"field":        "type",
				"lessOrEquals": 10,
			},
			expected: &PolicyRuleBody{
				IfBody: condition.LessOrEquals{
					BaseCondition: condition.BaseCondition{
						Subject: OperationField("type"),
					},
					Value: 10,
				},
			},
		},
		{
			name: "Greater",
			input: map[string]any{
				"field":   "type",
				"greater": 10,
			},
			expected: &PolicyRuleBody{
				IfBody: condition.Greater{
					BaseCondition: condition.BaseCondition{
						Subject: OperationField("type"),
					},
					Value: 10,
				},
			},
		},
		{
			name: "GreaterOrEquals",
			input: map[string]any{
				"field":           "type",
				"greaterOrEquals": 10,
			},
			expected: &PolicyRuleBody{
				IfBody: condition.GreaterOrEquals{
					BaseCondition: condition.BaseCondition{
						Subject: OperationField("type"),
					},
					Value: 10,
				},
			},
		},
		{
			name: "Exists",
			input: map[string]any{
				"field":  "type",
				"exists": true,
			},
			expected: &PolicyRuleBody{
				IfBody: condition.Exists{
					BaseCondition: condition.BaseCondition{
						Subject: OperationField("type"),
					},
					Value: true,
				},
			},
		},
		{
			name: "Unknown BaseCondition",
			input: map[string]any{
				"unknown": "value",
			},
			expected: nil, // Expecting a panic
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expected == nil {
				assert.Panics(t, func() {
					NewPolicyRuleBody(tt.input, shared.NewContext())
				})
			} else {
				stub := gostub.Stub(&NeoConditionNameGenerator, func(ctx *shared.Context) (string, error) {
					return "condition1", nil
				})
				defer stub.Reset()
				result := NewPolicyRuleBody(tt.input, shared.NewContext())
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
