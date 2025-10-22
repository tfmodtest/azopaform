package pkg

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/open-policy-agent/opa/v1/format"
	"github.com/prashantv/gostub"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tfmodtest/azopaform/pkg/condition"
	"github.com/tfmodtest/azopaform/pkg/operation"
	"github.com/tfmodtest/azopaform/pkg/shared"
)

const listFieldJson = `{
  "properties": {
    "displayName": "CORS should not allow every domain to access your API for FHIR",
    "policyType": "BuiltIn",
    "mode": "Indexed",
    "description": "Cross-Origin Resource Sharing (CORS) should not allow all domains to access your API for FHIR. To protect your API for FHIR, remove access for all domains and explicitly define the domains allowed to connect.",
    "metadata": {
      "version": "1.1.0",
      "category": "API for FHIR"
    },
    "version": "1.1.0",
    "parameters": {
      "effect": {
        "type": "String",
        "metadata": {
          "displayName": "Effect",
          "description": "Enable or disable the execution of the policy"
        },
        "allowedValues": [
          "audit",
          "Audit",
          "disabled",
          "Disabled",
		  "Deny"
        ],
        "defaultValue": "Deny"
      }
    },
    "policyRule": {
      "if": {
        "allOf": [
          {
            "field": "type",
            "equals": "Microsoft.HealthcareApis/services"
          },
          {
            "not": {
              "field": "Microsoft.HealthcareApis/services/corsConfiguration.origins[*]",
              "notEquals": "*"
            }
          }
        ]
      },
      "then": {
        "effect": "[parameters('effect')]"
      }
    }
  },
  "id": "/providers/Microsoft.Authorization/policyDefinitions/0fea8f8a-4169-495d-8307-30ec335f387d",
  "name": "0fea8f8a-4169-495d-8307-30ec335f387d"
}`

const denyJson = `{
  "properties": {
    "displayName": "App Service apps should use a SKU that supports private link",
    "policyType": "BuiltIn",
    "mode": "Indexed",
    "description": "With supported SKUs, Azure Private Link lets you connect your virtual network to Azure services without a public IP address at the source or destination. The Private Link platform handles the connectivity between the consumer and services over the Azure backbone network. By mapping private endpoints to apps, you can reduce data leakage risks. Learn more about private links at: https://aka.ms/private-link.",
    "metadata": {
      "version": "4.1.0",
      "category": "App Service"
    },
    "version": "4.1.0",
    "parameters": {
      "effect": {
        "type": "String",
        "metadata": {
          "displayName": "Effect",
          "description": "Enable or disable the execution of the policy"
        },
        "allowedValues": [
          "Audit",
          "Deny",
          "Disabled"
        ],
        "defaultValue": "Deny"
      }
    },
    "policyRule": {
      "if": {
        "allOf": [
          {
            "field": "type",
            "equals": "Microsoft.Web/serverFarms"
          },
          {
            "anyOf": [
              {
                "field": "Microsoft.Web/serverFarms/sku.tier",
                "notIn": [
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
                  "WorkflowStandard"
                ]
              },
              {
                "field": "Microsoft.Web/serverFarms/sku.name",
                "notIn": [
                  "B1",
                  "B2",
                  "B3",
                  "S1",
                  "S2",
                  "S3",
                  "EP1",
                  "EP2",
                  "EP3",
                  "P1",
                  "P2",
                  "P3",
                  "P1V2",
                  "P2V2",
                  "P3V2",
                  "P0V3",
                  "P1V3",
                  "P2V3",
                  "P3V3",
                  "P1MV3",
                  "P2MV3",
                  "P3MV3",
                  "P4MV3",
                  "P5MV3",
                  "I1",
                  "I2",
                  "I3",
                  "I1V2",
                  "I2V2",
                  "I3V2",
                  "I4V2",
                  "I5V2",
                  "I6V2",
                  "WS1",
                  "WS2",
                  "WS3"
                ]
              }
            ]
          }
        ]
      },
      "then": {
        "effect": "[parameters('effect')]"
      }
    }
  },
  "id": "/providers/Microsoft.Authorization/policyDefinitions/546fe8d2-368d-4029-a418-6af48a7f61e5",
  "name": "546fe8d2-368d-4029-a418-6af48a7f61e5"
}`

func TestBasicTestAzurePolicyToRego(t *testing.T) {
	cases := []struct {
		desc                  string
		inputDirPath          string
		mockFs                map[string]string
		input                 map[string]any
		generatedRegoFileName string
		deny                  bool
	}{
		{
			desc:         "deny.json_allow_type_mismatch",
			inputDirPath: "",
			mockFs: map[string]string{
				"deny.json": denyJson,
			},
			generatedRegoFileName: "deny.rego",
			input: map[string]any{
				"terraform_version": "1.11.0",
				"resource_changes": []any{
					map[string]any{
						"type": "azurerm_resource_group",
						"change": map[string]any{
							"after": map[string]any{
								"body": map[string]any{
									"properties": map[string]any{},
								},
							},
						},
					},
				},
			},
			deny: false,
		},
		{
			desc:         "deny.json_deny_sku_tier",
			inputDirPath: "",
			mockFs: map[string]string{
				"deny.json": denyJson,
			},
			generatedRegoFileName: "deny.rego",
			input: map[string]any{
				"terraform_version": "1.11.0",
				"resource_changes": []any{
					map[string]any{
						"address": "azapi_resource.this",
						"mode":    "managed",
						"type":    "azapi_resource",
						"change": map[string]any{
							"after": map[string]any{
								"type": "Microsoft.Web/serverFarms@2024-04-01",
								"body": map[string]any{
									"properties": map[string]any{
										"sku": map[string]any{
											"tier": "Basic",
											"name": "NotB1",
										},
									},
								},
							},
						},
					},
				},
			},
			deny: true,
		},
		{
			desc:         "deny.json_deny_sku_name",
			inputDirPath: "",
			mockFs: map[string]string{
				"deny.json": denyJson,
			},
			generatedRegoFileName: "deny.rego",
			input: map[string]any{
				"terraform_version": "1.11.0",
				"resource_changes": []any{
					map[string]any{
						"address": "azapi_resource.this",
						"mode":    "managed",
						"type":    "azapi_resource",
						"change": map[string]any{
							"after": map[string]any{
								"type": "Microsoft.Web/serverFarms@2024-04-01",
								"body": map[string]any{
									"properties": map[string]any{
										"sku": map[string]any{
											"tier": "NotBasic",
											"name": "B1",
										},
									},
								},
							},
						},
					},
				},
			},
			deny: true,
		},
		{
			desc:         "deny.json_allow_full_input",
			inputDirPath: "",
			mockFs: map[string]string{
				"deny.json": denyJson,
			},
			generatedRegoFileName: "deny.rego",
			input: map[string]any{
				"terraform_version": "1.11.0",
				"resource_changes": []any{
					map[string]any{
						"address": "azapi_resource.this",
						"mode":    "managed",
						"type":    "azapi_resource",
						"change": map[string]any{
							"after": map[string]any{
								"type": "Microsoft.Web/serverFarms@2024-04-01",
								"body": map[string]any{
									"properties": map[string]any{
										"sku": map[string]any{
											"tier": "NotBasic",
											"name": "NotB1",
										},
									},
								},
							},
						},
					},
				},
			},
			deny: true,
		},
		{
			desc:         "policy contains lists with multiple indexes",
			inputDirPath: "",
			mockFs: map[string]string{
				"list.json": listFieldJson,
			},
			generatedRegoFileName: "list.rego",
			input: map[string]any{
				"terraform_version": "1.11.0",
				"resource_changes": []any{
					map[string]any{
						"address": "azapi_resource.this",
						"mode":    "managed",
						"type":    "azapi_resource",
						"change": map[string]any{
							"after": map[string]any{
								"type": "Microsoft.HealthcareApis/services@2024-04-01",
								"body": map[string]any{
									"properties": map[string]any{
										"corsConfiguration": map[string]any{
											"origins": []string{
												"*",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			deny: true,
		},
		{
			desc:         "policy contains lists with multiple indexes negative",
			inputDirPath: "",
			mockFs: map[string]string{
				"list.json": listFieldJson,
			},
			generatedRegoFileName: "list.rego",
			input: map[string]any{
				"terraform_version": "1.11.0",
				"resource_changes": []any{
					map[string]any{
						"address": "azapi_resource.this",
						"mode":    "managed",
						"type":    "azapi_resource",
						"change": map[string]any{
							"after": map[string]any{
								"type": "Microsoft.HealthcareApis/services@2024-04-01",
								"body": map[string]any{
									"properties": map[string]any{
										"corsConfiguration": map[string]any{
											"origins": []string{
												"http://*.example.com",
												"http://*.example2.com",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			deny: false,
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			fs := fakeFs(c.mockFs)
			stub := gostub.Stub(&Fs, fs)
			defer stub.Reset()
			policyPath := ""
			if len(c.mockFs) == 1 {
				for n := range c.mockFs {
					policyPath = n
				}
			}
			require.NoError(t, AzurePolicyToRego(policyPath, c.inputDirPath, shared.NewContext()))
			content, err := afero.ReadFile(fs, c.generatedRegoFileName)
			require.NoError(t, err)
			generated := string(content) + "\n" + shared.UtilsRego
			ctx := shared.NewContext()
			shared.AssertRego(t, "data.main.deny", generated, c.input, c.deny, ctx)
		})
	}
}

func fakeFs(files map[string]string) afero.Fs {
	fs := afero.NewMemMapFs()
	for n, content := range files {
		_ = afero.WriteFile(fs, n, []byte(content), 0644)
	}
	return fs
}

func TestRule_SaveToDisk(t *testing.T) {
	rule := &Rule{
		Name:       "test",
		Properties: &PolicyRuleModel{},
		Id:         "",
		path:       "tmp.json",
		result:     "hello",
	}
	t.Run("SaveToDisk", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		stub := gostub.Stub(&Fs, fs)
		defer stub.Reset()
		err := rule.SaveToDisk()
		require.NoError(t, err)
		file, err := afero.ReadFile(fs, "tmp.rego")
		require.NoError(t, err)
		assert.Equal(t, "hello", string(file))
	})
}

func TestNeoAzPolicy2Rego(t *testing.T) {
	path := "deny.json"
	t.Run("loadRule", func(t *testing.T) {
		fs := prepareMemFs(t)
		counter := 1
		stub := gostub.Stub(&Fs, fs).Stub(&operation.RandomHelperFunctionNameGenerator, func() string {
			defer func() {
				counter++
			}()
			return fmt.Sprintf("condition%d", counter)
		})
		defer stub.Reset()

		ctx := shared.NewContext()
		rule, err := loadRule(path, ctx)
		require.NoError(t, err)
		err = rule.SaveToDisk()
		require.NoError(t, err)
		file, err := afero.ReadFile(fs, "deny.rego")
		require.NoError(t, err)
		expected := `package main
import rego.v1

warn if {
    res := resource(input, "azapi_resource")[_]
    condition2(res)
}
condition2(r) if {
    is_azure_type(r.values, "Microsoft.Web/serverFarms")
    condition1(r)
}
condition1(r) if {
    not r.values.body.properties.sku.tier in ["Basic","Standard","ElasticPremium","Premium","PremiumV2","Premium0V3","PremiumV3","PremiumMV3","Isolated","IsolatedV2","WorkflowStandard"]
    not r.values.body.properties.sku.name in ["B1","B2","B3","S1","S2","S3","EP1","EP2","EP3","P1","P2","P3","P1V2","P2V2","P3V2","P0V3","P1V3","P2V3","P3V3","P1MV3","P2MV3","P3MV3","P4MV3","P5MV3","I1","I2","I3","I1V2","I2V2","I3V2","I4V2","I5V2","I6V2","WS1","WS2","WS3"]
}
`
		formattedExpected, err := format.Source("temp.rego", []byte(expected))
		require.NoError(t, err)
		assert.Equal(t, string(formattedExpected), string(file))
	})
}

func TestAzPolicy2Rego_customizePackageName(t *testing.T) {
	path := "deny.json"
	fs := prepareMemFs(t)
	stub := gostub.Stub(&Fs, fs)
	defer stub.Reset()
	ctx := shared.NewContextWithOptions(shared.Options{PackageName: "customized"})
	rule, err := loadRule(path, ctx)
	require.NoError(t, err)
	regoCode, err := rule.Rego(ctx)
	require.NoError(t, err)
	formattedExpected, err := format.Source("temp.rego", []byte(regoCode))
	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(string(formattedExpected), "package customized"))
}

func TestAzPolicy2Rego_customizeUtilRegoFileName(t *testing.T) {
	path := "deny.json"
	fs := prepareMemFs(t)
	stub := gostub.Stub(&Fs, fs)
	defer stub.Reset()
	customizedRegoFileName := "customized.rego"
	err := AzurePolicyToRego(path, "", shared.NewContextWithOptions(shared.Options{
		UtilRegoFileName: customizedRegoFileName,
	}))
	require.NoError(t, err)
	readfile, err := afero.ReadFile(fs, customizedRegoFileName)
	require.NoError(t, err)
	assert.Contains(t, string(readfile), shared.UtilsRego)
}

func TestAzPolicy2Rego_WithUtilLibraryPackageName_NoUtilFileGenerated(t *testing.T) {
	fs := prepareMemFs(t)
	stub := gostub.Stub(&Fs, fs)
	defer stub.Reset()

	ctx := shared.NewContextWithOptions(shared.Options{
		UtilLibraryPackageName: "util",
	})

	err := AzurePolicyToRego("deny.json", "", ctx)
	require.NoError(t, err)

	matches, err := afero.Glob(fs, "*.rego")
	require.NoError(t, err)
	// Only deny.json exist
	assert.Len(t, matches, 1)
	assert.Equal(t, "deny.rego", matches[0])
}

func TestParseCondition(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]any
		expected shared.Rego
	}{
		{
			name: "Equals",
			input: map[string]any{
				"field":  "type",
				"equals": "Microsoft.Web/serverFarms",
			},
			expected: condition.Equals{
				BaseCondition: condition.BaseCondition{
					Subject: condition.FieldValue{
						Name: "type",
					},
				},
				Value: "Microsoft.Web/serverFarms",
			},
		},
		{
			name: "NotEquals",
			input: map[string]any{
				"field":     "type",
				"notEquals": "Microsoft.Web/serverFarms",
			},
			expected: condition.NotEquals{
				BaseCondition: condition.BaseCondition{
					Subject: condition.FieldValue{Name: "type"},
				},
				Value: "Microsoft.Web/serverFarms",
			},
		},
		{
			name: "Like",
			input: map[string]any{
				"field": "type",
				"like":  "Microsoft.Web/serverFarms",
			},
			expected: condition.Like{
				BaseCondition: condition.BaseCondition{
					Subject: condition.FieldValue{Name: "type"},
				},
				Value: "Microsoft.Web/serverFarms",
			},
		},
		{
			name: "NotLike",
			input: map[string]any{
				"field":   "type",
				"notLike": "Microsoft.Web/serverFarms",
			},
			expected: condition.NotLike{
				BaseCondition: condition.BaseCondition{
					Subject: condition.FieldValue{Name: "type"},
				},
				Value: "Microsoft.Web/serverFarms",
			},
		},
		{
			name: "In",
			input: map[string]any{
				"field": "type",
				"in":    []any{"Microsoft.Web/serverFarms", "Microsoft.Compute/virtualMachines"},
			},
			expected: condition.In{
				BaseCondition: condition.BaseCondition{
					Subject: condition.FieldValue{Name: "type"},
				},
				Values: []string{"Microsoft.Web/serverFarms", "Microsoft.Compute/virtualMachines"},
			},
		},
		{
			name: "NotIn",
			input: map[string]any{
				"field": "type",
				"notIn": []any{"Microsoft.Web/serverFarms", "Microsoft.Compute/virtualMachines"},
			},
			expected: condition.NotIn{
				BaseCondition: condition.BaseCondition{
					Subject: condition.FieldValue{Name: "type"},
				},
				Values: []string{"Microsoft.Web/serverFarms", "Microsoft.Compute/virtualMachines"},
			},
		},
		{
			name: "Contains",
			input: map[string]any{
				"field":    "type",
				"contains": "Microsoft.Web/serverFarms",
			},
			expected: condition.Contains{
				BaseCondition: condition.BaseCondition{
					Subject: condition.FieldValue{Name: "type"},
				},
				Value: "Microsoft.Web/serverFarms",
			},
		},
		{
			name: "NotContains",
			input: map[string]any{
				"field":       "type",
				"notContains": "Microsoft.Web/serverFarms",
			},
			expected: condition.NotContains{
				BaseCondition: condition.BaseCondition{
					Subject: condition.FieldValue{Name: "type"},
				},
				Value: "Microsoft.Web/serverFarms",
			},
		},
		{
			name: "ContainsKey",
			input: map[string]any{
				"field":       "type",
				"containsKey": "Microsoft.Web/serverFarms",
			},
			expected: condition.ContainsKey{
				BaseCondition: condition.BaseCondition{
					Subject: condition.FieldValue{Name: "type"},
				},
				KeyName: "Microsoft.Web/serverFarms",
			},
		},
		{
			name: "NotContainsKey",
			input: map[string]any{
				"field":          "type",
				"notContainsKey": "Microsoft.Web/serverFarms",
			},
			expected: condition.NotContainsKey{
				BaseCondition: condition.BaseCondition{
					Subject: condition.FieldValue{Name: "type"},
				},
				KeyName: "Microsoft.Web/serverFarms",
			},
		},
		{
			name: "Less",
			input: map[string]any{
				"field": "number",
				"less":  10,
			},
			expected: condition.Less{
				BaseCondition: condition.BaseCondition{
					Subject: condition.FieldValue{Name: "number"},
				},
				Value: 10,
			},
		},
		{
			name: "LessOrEquals",
			input: map[string]any{
				"field":        "number",
				"lessOrEquals": 10,
			},
			expected: condition.LessOrEquals{
				BaseCondition: condition.BaseCondition{
					Subject: condition.FieldValue{Name: "number"},
				},
				Value: 10,
			},
		},
		{
			name: "Greater",
			input: map[string]any{
				"field":   "number",
				"greater": 10,
			},
			expected: condition.Greater{
				BaseCondition: condition.BaseCondition{
					Subject: condition.FieldValue{Name: "number"},
				},
				Value: 10,
			},
		},
		{
			name: "GreaterOrEquals",
			input: map[string]any{
				"field":           "number",
				"greaterOrEquals": 10,
			},
			expected: condition.GreaterOrEquals{
				BaseCondition: condition.BaseCondition{
					Subject: condition.FieldValue{Name: "number"},
				},
				Value: 10,
			},
		},
		{
			name: "Exists",
			input: map[string]any{
				"field":  "type",
				"exists": true,
			},
			expected: condition.Exists{
				BaseCondition: condition.BaseCondition{
					Subject: condition.FieldValue{Name: "type"},
				},
				Value: true,
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
			ctx := shared.NewContext()
			if tt.expected == nil {
				_, err := NewPolicyRuleBody(tt.input).GetIf(ctx)
				assert.NotNil(t, err)
			} else {
				result := NewPolicyRuleBody(tt.input)
				r, err := result.GetIf(ctx)
				require.NoError(t, err)
				assert.Equal(t, tt.expected, r.rego)
			}
		})
	}
}

func prepareMemFs(t *testing.T) afero.Fs {
	fs := afero.NewMemMapFs()
	files := []string{
		"deny.json",
	}
	for _, file := range files {
		content, err := os.ReadFile(fmt.Sprintf("test-fixtures/%s", file))
		require.NoError(t, err)
		err = afero.WriteFile(fs, file, content, os.ModePerm)
		require.NoError(t, err)
	}
	return fs
}
