package pkg

import (
	"context"
	"github.com/open-policy-agent/opa/rego"
	"github.com/prashantv/gostub"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const conditionRegoTemplate = `package main

import rego.v1

default allow := false
allow if %s`

func TestOperations(t *testing.T) {
	tests := []struct {
		name      string
		operation Rego
		expected  string
	}{
		{
			name: "NestedWhereOperator",
			operation: WhereOperator{
				Conditions: []Rego{
					AllOf{
						Conditions: []Rego{
							EqualsCondition{
								condition: condition{
									Subject: OperationField("Microsoft.Web/serverFarms/sku.tier"),
								},
								Value: "Standard",
							},
							ExistsCondition{
								condition: condition{
									Subject: OperationField("Microsoft.Web/serverFarms/sku.name"),
								},
								Value: true,
							},
							EqualsCondition{
								condition: condition{
									Subject: OperationField("Microsoft.Web/serverFarms/sku.size"),
								},
								Value: "P1v3",
							},
						},
						ConditionSetName: "aaaaa",
					},
				},
				ConditionSetName: "aaaaaaaaa",
			},
			expected: "aaaaaaaaa(x) if {\naaaaa(x)\n}\naaaaa(x) if {\nr.change.after.sku[x].tier == \"Standard\"\nr.change.after.sku_name\nr.change.after.sku[x].size == \"P1v3\"\n}",
		},
		{
			name: "WhereOperator",
			operation: WhereOperator{
				Conditions: []Rego{
					EqualsCondition{
						condition: condition{
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
				Conditions: []Rego{
					AllOf{
						Conditions: []Rego{
							EqualsCondition{
								condition: condition{
									Subject: OperationField("type"),
								},
								Value: "azurerm_app_service_plan",
							},
							ExistsCondition{
								condition: condition{
									Subject: OperationField("Microsoft.Web/serverFarms/sku.name"),
								},
								Value: true,
							},
						},
						ConditionSetName: "aaaaa",
					},
					AnyOf{
						Conditions: []Rego{
							EqualsCondition{
								condition: condition{
									Subject: OperationField("Microsoft.Web/serverFarms/sku.tier"),
								},
								Value: "Standard",
							},
							EqualsCondition{
								condition: condition{
									Subject: OperationField("Microsoft.Web/serverFarms/sku.tier"),
								},
								Value: "Basic",
							},
						},
						ConditionSetName: "aaaaaaa",
					},
				},
				ConditionSetName: "aaaaa",
			},
			expected: "aaaaa if {\naaaaa\nnot aaaaaaa\n}\naaaaa if {\nr.type == \"azurerm_app_service_plan\"\nr.change.after.sku_name\n}\naaaaaaa if {\nr.change.after.sku[0].tier != \"Standard\"\nr.change.after.sku[0].tier != \"Basic\"\n}",
		},
		{
			name: "NestedAnyOfOperator",
			operation: AnyOf{
				Conditions: []Rego{
					AnyOf{
						Conditions: []Rego{
							EqualsCondition{
								condition: condition{
									Subject: OperationField("type"),
								},
								Value: "azurerm_app_service_plan",
							},
							EqualsCondition{
								condition: condition{
									Subject: OperationField("type"),
								},
								Value: "azurerm_app_service_environment",
							},
						},
						ConditionSetName: "aaaaaaa",
					},
					AnyOf{
						Conditions: []Rego{
							EqualsCondition{
								condition: condition{
									Subject: OperationField("Microsoft.Web/serverFarms/sku.tier"),
								},
								Value: "Standard",
							},
							EqualsCondition{
								condition: condition{
									Subject: OperationField("Microsoft.Web/serverFarms/sku.tier"),
								},
								Value: "Basic",
							},
						},
						ConditionSetName: "aaaaaaa",
					},
				},
				ConditionSetName: "aaaaaaa",
			},
			expected: "aaaaaaa if {\naaaaaaa\naaaaaaa\n}\naaaaaaa if {\nr.type != \"azurerm_app_service_plan\"\nr.type != \"azurerm_app_service_environment\"\n}\naaaaaaa if {\nr.change.after.sku[0].tier != \"Standard\"\nr.change.after.sku[0].tier != \"Basic\"\n}",
		},
		{
			name: "AllOfOperator",
			operation: AllOf{
				Conditions: []Rego{
					EqualsCondition{
						condition: condition{
							Subject: OperationField("type"),
						},
						Value: "azurerm_app_service_plan",
					},
					ExistsCondition{
						condition: condition{
							Subject: OperationField("Microsoft.Web/serverFarms/sku.name"),
						},
						Value: true,
					},
				},
				ConditionSetName: "aaaaa",
			},
			expected: "aaaaa if {\nr.type == \"azurerm_app_service_plan\"\nr.change.after.sku_name\n}",
		},
		{
			name: "AnyOfOperator",
			operation: AnyOf{
				Conditions: []Rego{
					EqualsCondition{
						condition: condition{
							Subject: OperationField("Microsoft.Web/serverFarms/sku.tier"),
						},
						Value: "Standard",
					},
					InCondition{
						condition: condition{
							Subject: OperationField("Microsoft.Web/serverFarms/sku.tier"),
						},
						Values: []string{"Basic", "Premium"},
					},
				},
				ConditionSetName: "aaaaaaa",
			},
			expected: "aaaaaaa if {\nr.change.after.sku[0].tier != \"Standard\"\nnot r.change.after.sku[0].tier in [\"Basic\",\"Premium\"]\n}",
		},
		{
			name: "NotOperator",
			operation: NotOperator{
				Body: EqualsCondition{
					condition: condition{
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
			name: "EqualsCondition",
			operation: EqualsCondition{
				condition: condition{
					Subject: OperationField("Microsoft.Web/serverFarms/sku.tier"),
				},
				Value: "Standard",
			},
			expected: "r.change.after.sku[0].tier == \"Standard\"",
		},
		{
			name: "NotEqualsCondition",
			operation: NotEqualsCondition{
				condition: condition{
					Subject: OperationField("Microsoft.Web/serverFarms/sku.tier"),
				},
				Value: "Standard",
			},
			expected: "r.change.after.sku[0].tier != \"Standard\"",
		},
		{
			name: "LikeCondition",
			operation: LikeCondition{
				condition: condition{
					Subject: OperationField("Microsoft.Web/serverFarms/sku.tier"),
				},
				Value: `^[^@]+@[^@]+\.[^@]+$`,
			},
			expected: "regex.match(\"^[^@]+@[^@]+\\.[^@]+$\",r.change.after.sku[0].tier)",
		},
		{
			name: "NotLikeCondition",
			operation: NotLikeCondition{
				condition: condition{
					Subject: OperationField("Microsoft.Web/serverFarms/sku.tier"),
				},
				Value: `^[^@]+@[^@]+\.[^@]+$`,
			},
			expected: "not regex.match(`^[^@]+@[^@]+\\.[^@]+$`,r.change.after.sku[0].tier)",
		},
		{
			name: "InCondition",
			operation: InCondition{
				condition: condition{
					Subject: OperationField("Microsoft.Web/serverFarms/sku.tier"),
				},
				Values: []string{"Basic", "Standard", "Premium"},
			},
			expected: "some r.change.after.sku[0].tier in [\"Basic\",\"Standard\",\"Premium\"]",
		},
		{
			name: "NotInCondition",
			operation: NotInCondition{
				condition: condition{
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
			ctx := NewContext()
			pushResourceType(ctx, "Microsoft.Web/serverFarms")
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
				IfBody: GreaterCondition{
					condition: condition{
						Subject: CountOperator{
							Where: WhereOperator{
								Conditions: []Rego{
									EqualsCondition{
										condition: condition{
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
					Body: NotEqualsCondition{
						condition: condition{
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
					Conditions: []Rego{
						AnyOf{
							Conditions: []Rego{
								ExistsCondition{
									condition: condition{
										Subject: FieldValue{
											Name: "Microsoft.Sql/servers/minimalTlsVersion",
										},
									},
									Value: false,
								},
								LessCondition{
									condition: condition{
										Subject: FieldValue{
											Name: "Microsoft.Sql/servers/minimalTlsVersion",
										},
									},
									Value: "1.2",
								},
							},
							ConditionSetName: "condition1",
						},
						AllOf{
							Conditions: []Rego{
								EqualsCondition{
									condition: condition{
										Subject: FieldValue{
											Name: "type",
										},
									},
									Value: "azurerm_mssql_server",
								},
								ExistsCondition{
									condition: condition{
										Subject: FieldValue{
											Name: "Microsoft.Sql/servers/minimalTlsVersion",
										},
									},
									Value: true,
								},
							},
							ConditionSetName: "condition1",
						},
					},
					ConditionSetName: "condition1",
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
					Conditions: []Rego{
						ExistsCondition{
							condition: condition{
								Subject: FieldValue{
									Name: "Microsoft.Sql/servers/minimalTlsVersion",
								},
							},
							Value: false,
						},
						LessCondition{
							condition: condition{
								Subject: FieldValue{
									Name: "Microsoft.Sql/servers/minimalTlsVersion",
								},
							},
							Value: "1.2",
						},
					},
					ConditionSetName: "condition1",
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
					Conditions: []Rego{
						EqualsCondition{
							condition: condition{
								Subject: FieldValue{
									Name: "type",
								},
							},
							Value: "azurerm_healthcare_service",
						},
						ExistsCondition{
							condition: condition{
								Subject: FieldValue{
									Name: "Microsoft.HealthcareApis/services/cosmosDbConfiguration.keyVaultKeyUri",
								},
							},
							Value: false,
						},
					},
					ConditionSetName: "condition1",
				},
			},
		},
		{
			name: "EqualsCondition",
			input: map[string]any{
				"field":  "type",
				"equals": "Microsoft.Web/serverFarms",
			},
			expected: &PolicyRuleBody{
				IfBody: EqualsCondition{
					condition: condition{
						Subject: OperationField("type"),
					},
					Value: "Microsoft.Web/serverFarms",
				},
			},
		},
		{
			name: "NotEqualsCondition",
			input: map[string]any{
				"field":     "type",
				"notEquals": "Microsoft.Web/serverFarms",
			},
			expected: &PolicyRuleBody{
				IfBody: NotEqualsCondition{
					condition: condition{
						Subject: OperationField("type"),
					},
					Value: "Microsoft.Web/serverFarms",
				},
			},
		},
		{
			name: "LikeCondition",
			input: map[string]any{
				"field": "type",
				"like":  "Microsoft.Web/serverFarms",
			},
			expected: &PolicyRuleBody{
				IfBody: LikeCondition{
					condition: condition{
						Subject: OperationField("type"),
					},
					Value: "Microsoft.Web/serverFarms",
				},
			},
		},
		{
			name: "NotLikeCondition",
			input: map[string]any{
				"field":   "type",
				"notLike": "Microsoft.Web/serverFarms",
			},
			expected: &PolicyRuleBody{
				IfBody: NotLikeCondition{
					condition: condition{
						Subject: OperationField("type"),
					},
					Value: "Microsoft.Web/serverFarms",
				},
			},
		},
		{
			name: "InCondition",
			input: map[string]any{
				"field": "type",
				"in":    []any{"Microsoft.Web/serverFarms", "Microsoft.Compute/virtualMachines"},
			},
			expected: &PolicyRuleBody{
				IfBody: InCondition{
					condition: condition{
						Subject: OperationField("type"),
					},
					Values: []string{"Microsoft.Web/serverFarms", "Microsoft.Compute/virtualMachines"},
				},
			},
		},
		{
			name: "NotInCondition",
			input: map[string]any{
				"field": "type",
				"notIn": []any{"Microsoft.Web/serverFarms", "Microsoft.Compute/virtualMachines"},
			},
			expected: &PolicyRuleBody{
				IfBody: NotInCondition{
					condition: condition{
						Subject: OperationField("type"),
					},
					Values: []string{"Microsoft.Web/serverFarms", "Microsoft.Compute/virtualMachines"},
				},
			},
		},
		{
			name: "ContainsCondition",
			input: map[string]any{
				"field":    "type",
				"contains": "Microsoft.Web/serverFarms",
			},
			expected: &PolicyRuleBody{
				IfBody: ContainsCondition{
					condition: condition{
						Subject: OperationField("type"),
					},
					Value: "Microsoft.Web/serverFarms",
				},
			},
		},
		{
			name: "NotContainsCondition",
			input: map[string]any{
				"field":       "type",
				"notContains": "Microsoft.Web/serverFarms",
			},
			expected: &PolicyRuleBody{
				IfBody: NotContainsCondition{
					condition: condition{
						Subject: OperationField("type"),
					},
					Value: "Microsoft.Web/serverFarms",
				},
			},
		},
		{
			name: "ContainsKeyCondition",
			input: map[string]any{
				"field":       "type",
				"containsKey": "Microsoft.Web/serverFarms",
			},
			expected: &PolicyRuleBody{
				IfBody: ContainsKeyCondition{
					condition: condition{
						Subject: OperationField("type"),
					},
					KeyName: "Microsoft.Web/serverFarms",
				},
			},
		},
		{
			name: "NotContainsKeyCondition",
			input: map[string]any{
				"field":          "type",
				"notContainsKey": "Microsoft.Web/serverFarms",
			},
			expected: &PolicyRuleBody{
				IfBody: NotContainsKeyCondition{
					condition: condition{
						Subject: OperationField("type"),
					},
					KeyName: "Microsoft.Web/serverFarms",
				},
			},
		},
		{
			name: "LessCondition",
			input: map[string]any{
				"field": "type",
				"less":  10,
			},
			expected: &PolicyRuleBody{
				IfBody: LessCondition{
					condition: condition{
						Subject: OperationField("type"),
					},
					Value: 10,
				},
			},
		},
		{
			name: "LessOrEqualsCondition",
			input: map[string]any{
				"field":        "type",
				"lessOrEquals": 10,
			},
			expected: &PolicyRuleBody{
				IfBody: LessOrEqualsCondition{
					condition: condition{
						Subject: OperationField("type"),
					},
					Value: 10,
				},
			},
		},
		{
			name: "GreaterCondition",
			input: map[string]any{
				"field":   "type",
				"greater": 10,
			},
			expected: &PolicyRuleBody{
				IfBody: GreaterCondition{
					condition: condition{
						Subject: OperationField("type"),
					},
					Value: 10,
				},
			},
		},
		{
			name: "GreaterOrEqualsCondition",
			input: map[string]any{
				"field":           "type",
				"greaterOrEquals": 10,
			},
			expected: &PolicyRuleBody{
				IfBody: GreaterOrEqualsCondition{
					condition: condition{
						Subject: OperationField("type"),
					},
					Value: 10,
				},
			},
		},
		{
			name: "ExistsCondition",
			input: map[string]any{
				"field":  "type",
				"exists": true,
			},
			expected: &PolicyRuleBody{
				IfBody: ExistsCondition{
					condition: condition{
						Subject: OperationField("type"),
					},
					Value: true,
				},
			},
		},
		{
			name: "Unknown condition",
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
					NewPolicyRuleBody(tt.input, context.Context(nil))
				})
			} else {
				stub := gostub.Stub(&NeoConditionNameGenerator, func(ctx context.Context) (string, error) {
					return "condition1", nil
				})
				defer stub.Reset()
				result := NewPolicyRuleBody(tt.input, NewContext())
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func assertRegoAllow(t *testing.T, cfg string, input *rego.EvalOption, allowed bool, ctx context.Context) {
	eval, err := rego.New(rego.Query("data.main.allow"), rego.Module("test.rego", cfg)).PrepareForEval(ctx)
	require.NoError(t, err)
	var result rego.ResultSet
	if input == nil {
		result, err = eval.Eval(ctx)
	} else {
		result, err = eval.Eval(ctx, *input)
	}
	require.NoError(t, err)
	assert.Equal(t, allowed, result.Allowed())
}
