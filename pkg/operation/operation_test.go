package operation

import (
	"testing"

	"github.com/prashantv/gostub"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tfmodtest/azopaform/pkg/condition"
	"github.com/tfmodtest/azopaform/pkg/shared"
	"github.com/tfmodtest/azopaform/pkg/value"
)

func TestParseOperation(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]any
		expected shared.Rego
	}{
		{
			name: "NotOperation",
			input: map[string]any{
				"not": map[string]any{
					"field":     "Microsoft.HealthcareApis/services/corsConfiguration.origins[*]",
					"notEquals": "*",
				},
			},
			expected: NewNot("condition1", condition.NotEquals{
				BaseCondition: condition.BaseCondition{
					Subject: value.FieldValue{
						Name: "Microsoft.HealthcareApis/services/corsConfiguration.origins[*]",
					},
				},
				Value: "*",
			}),
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
			expected: NewAnyOf("condition1", []shared.Rego{
				NewAnyOf("condition1", []shared.Rego{
					condition.Exists{
						BaseCondition: condition.BaseCondition{
							Subject: value.FieldValue{
								Name: "Microsoft.Sql/servers/minimalTlsVersion",
							},
						},
						Value: false,
					},
					condition.Less{
						BaseCondition: condition.BaseCondition{
							Subject: value.FieldValue{
								Name: "Microsoft.Sql/servers/minimalTlsVersion",
							},
						},
						Value: "1.2",
					},
				}),
				NewAllOf("condition1", []shared.Rego{
					condition.Equals{
						BaseCondition: condition.BaseCondition{
							Subject: value.FieldValue{
								Name: "type",
							},
						},
						Value: "Microsoft.Sql/servers",
					},
					condition.Exists{
						BaseCondition: condition.BaseCondition{
							Subject: value.FieldValue{
								Name: "Microsoft.Sql/servers/minimalTlsVersion",
							},
						},
						Value: true,
					},
				}),
			}),
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
			expected: NewAnyOf("condition1", []shared.Rego{
				condition.Exists{
					BaseCondition: condition.BaseCondition{
						Subject: value.FieldValue{
							Name: "Microsoft.Sql/servers/minimalTlsVersion",
						},
					},
					Value: false,
				},
				condition.Less{
					BaseCondition: condition.BaseCondition{
						Subject: value.FieldValue{
							Name: "Microsoft.Sql/servers/minimalTlsVersion",
						},
					},
					Value: "1.2",
				},
			}),
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
			expected: NewAllOf("condition1", []shared.Rego{
				condition.Equals{
					BaseCondition: condition.BaseCondition{
						Subject: value.FieldValue{
							Name: "type",
						},
					},
					Value: "Microsoft.HealthcareApis/services",
				},
				condition.Exists{
					BaseCondition: condition.BaseCondition{
						Subject: value.FieldValue{
							Name: "Microsoft.HealthcareApis/services/cosmosDbConfiguration.keyVaultKeyUri",
						},
					},
					Value: false,
				},
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stub := gostub.Stub(&RandomHelperFunctionNameGenerator, func() string {
				return "condition1"
			})
			defer stub.Reset()
			result, err := NewOperationOrCondition(tt.input, shared.NewContext())
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
