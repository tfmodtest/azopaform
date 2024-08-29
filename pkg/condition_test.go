package pkg

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNotInCondition(t *testing.T) {
	sut := NotInOperation{
		operation: operation{
			Subject: OperationField("Microsoft.Web/serverFarms/sku.tier"),
		},
		Values: []string{
			"Basic",
			"Standard",
			"ElasticPremium",
			"Premium",
			"PremiumV2",
			"Premium0V3",
			"PremiumV3",
			"PremiumMV3",
			"Isolated",
			"IsolatedV2",
			"WorkflowStandard",
		},
	}
	ctx := NewContext()
	pushResourceType(ctx, "Microsoft.Web/serverFarms")
	actual, err := sut.Rego(ctx)
	require.NoError(t, err)
	assert.Equal(t, `not r.change.after.sku[0].tier in ["Basic","Standard","ElasticPremium","Premium","PremiumV2","Premium0V3","PremiumV3","PremiumMV3","Isolated","IsolatedV2","WorkflowStandard"]`, actual)
}

func TestInCondition(t *testing.T) {
	sut := InOperation{
		operation: operation{
			Subject: OperationField("Microsoft.Web/serverFarms/sku.tier"),
		},
		Values: []string{
			"Basic",
			"Standard",
			"ElasticPremium",
			"Premium",
			"PremiumV2",
			"Premium0V3",
			"PremiumV3",
			"PremiumMV3",
			"Isolated",
			"IsolatedV2",
			"WorkflowStandard",
		},
	}
	ctx := NewContext()
	pushResourceType(ctx, "Microsoft.Web/serverFarms")
	actual, err := sut.Rego(ctx)
	require.NoError(t, err)
	assert.Equal(t, `some r.change.after.sku[0].tier in ["Basic","Standard","ElasticPremium","Premium","PremiumV2","Premium0V3","PremiumV3","PremiumMV3","Isolated","IsolatedV2","WorkflowStandard"]`, actual)
}

func TestLikeCondition(t *testing.T) {
	sut := LikeOperation{
		operation: operation{
			Subject: OperationField("Microsoft.Web/serverFarms/sku.tier"),
		},
		Value: `^[^@]+@[^@]+\.[^@]+$`,
	}
	ctx := NewContext()
	pushResourceType(ctx, "Microsoft.Web/serverFarms")
	actual, err := sut.Rego(ctx)
	require.NoError(t, err)
	assert.Equal(t, "regex.match(^[^@]+@[^@]+\\.[^@]+$,r.change.after.sku[0].tier)", actual)
}

func TestNotLikeCondition(t *testing.T) {
	sut := NotLikeOperation{
		operation: operation{
			Subject: OperationField("Microsoft.Web/serverFarms/sku.tier"),
		},
		Value: `^[^@]+@[^@]+\.[^@]+$`,
	}
	ctx := NewContext()
	pushResourceType(ctx, "Microsoft.Web/serverFarms")
	actual, err := sut.Rego(ctx)
	require.NoError(t, err)
	assert.Equal(t, "not regex.match(^[^@]+@[^@]+\\.[^@]+$,r.change.after.sku[0].tier)", actual)
}

func TestEqualsCondition(t *testing.T) {
	sut := EqualsOperation{
		operation: operation{
			Subject: OperationField("Microsoft.Web/serverFarms/sku.tier"),
		},
		Value: "Standard",
	}
	ctx := NewContext()
	pushResourceType(ctx, "Microsoft.Web/serverFarms")
	actual, err := sut.Rego(ctx)
	require.NoError(t, err)
	assert.Equal(t, "r.change.after.sku[0].tier == Standard", actual)
}

func TestNotEqualsCondition(t *testing.T) {
	sut := NotEqualsOperation{
		operation: operation{
			Subject: OperationField("Microsoft.Web/serverFarms/sku.tier"),
		},
		Value: "Standard",
	}
	ctx := NewContext()
	pushResourceType(ctx, "Microsoft.Web/serverFarms")
	actual, err := sut.Rego(ctx)
	require.NoError(t, err)
	assert.Equal(t, "r.change.after.sku[0].tier != Standard", actual)
}

func TestOperations(t *testing.T) {
	tests := []struct {
		name      string
		operation Rego
		expected  string
	}{
		{
			name: "EqualsOperation",
			operation: EqualsOperation{
				operation: operation{
					Subject: OperationField("Microsoft.Web/serverFarms/sku.tier"),
				},
				Value: "Standard",
			},
			expected: "r.change.after.sku[0].tier == Standard",
		},
		{
			name: "NotEqualsOperation",
			operation: NotEqualsOperation{
				operation: operation{
					Subject: OperationField("Microsoft.Web/serverFarms/sku.tier"),
				},
				Value: "Standard",
			},
			expected: "r.change.after.sku[0].tier != Standard",
		},
		{
			name: "LikeOperation",
			operation: LikeOperation{
				operation: operation{
					Subject: OperationField("Microsoft.Web/serverFarms/sku.tier"),
				},
				Value: `^[^@]+@[^@]+\.[^@]+$`,
			},
			expected: "regex.match(^[^@]+@[^@]+\\.[^@]+$,r.change.after.sku[0].tier)",
		},
		{
			name: "NotLikeOperation",
			operation: NotLikeOperation{
				operation: operation{
					Subject: OperationField("Microsoft.Web/serverFarms/sku.tier"),
				},
				Value: `^[^@]+@[^@]+\.[^@]+$`,
			},
			expected: "not regex.match(^[^@]+@[^@]+\\.[^@]+$,r.change.after.sku[0].tier)",
		},
		{
			name: "InOperation",
			operation: InOperation{
				operation: operation{
					Subject: OperationField("Microsoft.Web/serverFarms/sku.tier"),
				},
				Values: []string{"Basic", "Standard", "Premium"},
			},
			expected: "some r.change.after.sku[0].tier in [\"Basic\",\"Standard\",\"Premium\"]",
		},
		{
			name: "NotInOperation",
			operation: NotInOperation{
				operation: operation{
					Subject: OperationField("Microsoft.Web/serverFarms/sku.tier"),
				},
				Values: []string{"Basic", "Standard", "Premium"},
			},
			expected: "not r.change.after.sku[0].tier in [\"Basic\",\"Standard\",\"Premium\"]",
		},
	}

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
					EqualsOperation{
						operation: operation{
							Subject: FieldValue{
								Name: "type",
							},
						},
						Value: "Microsoft.HealthcareApis/services",
					},
					ExistsOperation{
						operation: operation{
							Subject: FieldValue{
								Name: "Microsoft.HealthcareApis/services/cosmosDbConfiguration.keyVaultKeyUri",
							},
						},
						Value: false,
					},
				},
			},
		},
		{
			name: "EqualsOperation",
			input: map[string]any{
				"field":  "type",
				"equals": "Microsoft.Web/serverFarms",
			},
			expected: &PolicyRuleBody{
				IfBody: EqualsOperation{
					operation: operation{
						Subject: OperationField("type"),
					},
					Value: "Microsoft.Web/serverFarms",
				},
			},
		},
		{
			name: "NotEqualsOperation",
			input: map[string]any{
				"field":     "type",
				"notEquals": "Microsoft.Web/serverFarms",
			},
			expected: &PolicyRuleBody{
				IfBody: NotEqualsOperation{
					operation: operation{
						Subject: OperationField("type"),
					},
					Value: "Microsoft.Web/serverFarms",
				},
			},
		},
		{
			name: "LikeOperation",
			input: map[string]any{
				"field": "type",
				"like":  "Microsoft.Web/serverFarms",
			},
			expected: &PolicyRuleBody{
				IfBody: LikeOperation{
					operation: operation{
						Subject: OperationField("type"),
					},
					Value: "Microsoft.Web/serverFarms",
				},
			},
		},
		{
			name: "NotLikeOperation",
			input: map[string]any{
				"field":   "type",
				"notLike": "Microsoft.Web/serverFarms",
			},
			expected: &PolicyRuleBody{
				IfBody: NotLikeOperation{
					operation: operation{
						Subject: OperationField("type"),
					},
					Value: "Microsoft.Web/serverFarms",
				},
			},
		},
		{
			name: "InOperation",
			input: map[string]any{
				"field": "type",
				"in":    []any{"Microsoft.Web/serverFarms", "Microsoft.Compute/virtualMachines"},
			},
			expected: &PolicyRuleBody{
				IfBody: InOperation{
					operation: operation{
						Subject: OperationField("type"),
					},
					Values: []string{"Microsoft.Web/serverFarms", "Microsoft.Compute/virtualMachines"},
				},
			},
		},
		{
			name: "NotInOperation",
			input: map[string]any{
				"field": "type",
				"notIn": []any{"Microsoft.Web/serverFarms", "Microsoft.Compute/virtualMachines"},
			},
			expected: &PolicyRuleBody{
				IfBody: NotInOperation{
					operation: operation{
						Subject: OperationField("type"),
					},
					Values: []string{"Microsoft.Web/serverFarms", "Microsoft.Compute/virtualMachines"},
				},
			},
		},
		{
			name: "ContainsOperation",
			input: map[string]any{
				"field":    "type",
				"contains": "Microsoft.Web/serverFarms",
			},
			expected: &PolicyRuleBody{
				IfBody: ContainsOperation{
					operation: operation{
						Subject: OperationField("type"),
					},
					Value: "Microsoft.Web/serverFarms",
				},
			},
		},
		{
			name: "NotContainsOperation",
			input: map[string]any{
				"field":       "type",
				"notContains": "Microsoft.Web/serverFarms",
			},
			expected: &PolicyRuleBody{
				IfBody: NotContainsOperation{
					operation: operation{
						Subject: OperationField("type"),
					},
					Value: "Microsoft.Web/serverFarms",
				},
			},
		},
		{
			name: "ContainsKeyOperation",
			input: map[string]any{
				"field":       "type",
				"containsKey": "Microsoft.Web/serverFarms",
			},
			expected: &PolicyRuleBody{
				IfBody: ContainsKeyOperation{
					operation: operation{
						Subject: OperationField("type"),
					},
					KeyName: "Microsoft.Web/serverFarms",
				},
			},
		},
		{
			name: "NotContainsKeyOperation",
			input: map[string]any{
				"field":          "type",
				"notContainsKey": "Microsoft.Web/serverFarms",
			},
			expected: &PolicyRuleBody{
				IfBody: NotContainsKeyOperation{
					operation: operation{
						Subject: OperationField("type"),
					},
					KeyName: "Microsoft.Web/serverFarms",
				},
			},
		},
		{
			name: "LessOperation",
			input: map[string]any{
				"field": "type",
				"less":  10,
			},
			expected: &PolicyRuleBody{
				IfBody: LessOperation{
					operation: operation{
						Subject: OperationField("type"),
					},
					Value: 10,
				},
			},
		},
		{
			name: "LessOrEqualsOperation",
			input: map[string]any{
				"field":        "type",
				"lessOrEquals": 10,
			},
			expected: &PolicyRuleBody{
				IfBody: LessOrEqualsOperation{
					operation: operation{
						Subject: OperationField("type"),
					},
					Value: 10,
				},
			},
		},
		{
			name: "GreaterOperation",
			input: map[string]any{
				"field":   "type",
				"greater": 10,
			},
			expected: &PolicyRuleBody{
				IfBody: GreaterOperation{
					operation: operation{
						Subject: OperationField("type"),
					},
					Value: 10,
				},
			},
		},
		{
			name: "GreaterOrEqualsOperation",
			input: map[string]any{
				"field":           "type",
				"greaterOrEquals": 10,
			},
			expected: &PolicyRuleBody{
				IfBody: GreaterOrEqualsOperation{
					operation: operation{
						Subject: OperationField("type"),
					},
					Value: 10,
				},
			},
		},
		{
			name: "ExistsOperation",
			input: map[string]any{
				"field":  "type",
				"exists": true,
			},
			expected: &PolicyRuleBody{
				IfBody: ExistsOperation{
					operation: operation{
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
					NewPolicyRuleBody(tt.input)
				})
			} else {
				result := NewPolicyRuleBody(map[string]any{
					"if": tt.input,
				})
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
