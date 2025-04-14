package pkg

import (
	"fmt"
	"json-rule-finder/pkg/condition"
	"json-rule-finder/pkg/operation"
	"json-rule-finder/pkg/shared"
	"json-rule-finder/pkg/value"
	"os"
	"testing"

	"github.com/open-policy-agent/opa/format"
	"github.com/prashantv/gostub"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const list_field_json = `{
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
          "Disabled"
        ],
        "defaultValue": "Audit"
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

const count_json = `{
  "properties": {
    "displayName": "App Service Environment should have TLS 1.0 and 1.1 disabled",
    "policyType": "BuiltIn",
    "mode": "Indexed",
    "description": "TLS 1.0 and 1.1 are out-of-date protocols that do not support modern cryptographic algorithms. Disabling inbound TLS 1.0 and 1.1 traffic helps secure apps in an App Service Environment.",
    "metadata": {
      "version": "2.0.1",
      "category": "App Service"
    },
    "version": "2.0.1",
    "parameters": {
      "effect": {
        "type": "string",
        "defaultValue": "Deny",
        "allowedValues": [
          "Audit",
          "Deny",
          "Disabled"
        ],
        "metadata": {
          "displayName": "Effect",
          "description": "Enable or disable the execution of the policy"
        }
      }
    },
    "policyRule": {
      "if": {
        "allOf": [
          {
            "field": "type",
            "equals": "Microsoft.Web/hostingEnvironments"
          },
          {
            "field": "kind",
            "like": "ASE*"
          },
          {
            "count": {
              "field": "Microsoft.Web/HostingEnvironments/clusterSettings[*]",
              "where": {
                "allOf": [
                  {
                    "field": "Microsoft.Web/HostingEnvironments/clusterSettings[*].name",
                    "equals": "DisableTls1.0"
                  },
                  {
                    "field": "Microsoft.Web/HostingEnvironments/clusterSettings[*].value",
                    "equals": "1"
                  }
                ]
              }
            },
            "less": 1
          }
        ]
      },
      "then": {
        "effect": "[parameters('effect')]"
      }
    }
  },
  "id": "/providers/Microsoft.Authorization/policyDefinitions/d6545c6b-dd9d-4265-91e6-0b451e2f1c50",
  "name": "d6545c6b-dd9d-4265-91e6-0b451e2f1c50"
}`

const deny_json = `{
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

const nested_json = `{
  "properties": {
    "displayName": "API Management APIs should use only encrypted protocols",
    "policyType": "BuiltIn",
    "mode": "All",
    "description": "To ensure security of data in transit, APIs should be available only through encrypted protocols, like HTTPS or WSS. Avoid using unsecured protocols, such as HTTP or WS.",
    "metadata": {
      "version": "2.0.2",
      "category": "API Management"
    },
    "version": "2.0.2",
    "policyRule": {
      "if": {
        "allOf": [
          {
            "field": "type",
            "equals": "Microsoft.ApiManagement/service/apis"
          },
          {
            "count": {
              "field": "Microsoft.ApiManagement/service/apis/protocols[*]",
              "where": {
                "anyOf": [
                  {
                    "field": "Microsoft.ApiManagement/service/apis/protocols[*]",
                    "equals": "http"
                  },
                  {
                    "field": "Microsoft.ApiManagement/service/apis/protocols[*]",
                    "equals": "ws"
                  }
                ]
              }
            },
            "greaterOrEquals": 1
          }
        ]
      },
      "then": {
        "effect": "[parameters('effect')]"
      }
    },
    "parameters": {
      "effect": {
        "type": "string",
        "metadata": {
          "displayName": "Effect",
          "description": "Enable or disable the execution of the policy"
        },
        "allowedValues": [
          "Audit",
          "Disabled",
          "Deny"
        ],
        "defaultValue": "Audit"
      }
    }
  },
  "id": "/providers/Microsoft.Authorization/policyDefinitions/ee7495e7-3ba7-40b6-bfee-c29e22cc75d4",
  "name": "ee7495e7-3ba7-40b6-bfee-c29e22cc75d4"
}`

func TestBasicTestAzurePolicyToRego(t *testing.T) {
	//rulesJson, err := os.ReadFile("rules.json")
	//require.NoError(t, err)
	//outputJson, err := os.ReadFile("output.json")
	//require.NoError(t, err)

	//	expectedDenyRego := `package main
	//
	//import future.keywords.if
	//import future.keywords.in
	//tfplan := input if {
	//     input.terraform_version
	//} else := input.plan if {
	//     input.plan.terraform_version
	//}
	//
	//r := tfplan.resource_changes[_]
	//
	//warn if {
	// condition3
	//}
	//condition3 if {
	//r.type == "azurerm_service_plan"
	//condition2
	//}
	//condition2 if {
	//not r.change.after.sku[0].tier in ["Basic","Standard","ElasticPremium","Premium","PremiumV2","Premium0V3","PremiumV3","PremiumMV3","Isolated","IsolatedV2","WorkflowStandard"]
	//not r.change.after.sku_name in ["B1","B2","B3","S1","S2","S3","EP1","EP2","EP3","P1","P2","P3","P1V2","P2V2","P3V2","P0V3","P1V3","P2V3","P3V3","P1MV3","P2MV3","P3MV3","P4MV3","P5MV3","I1","I2","I3","I1V2","I2V2","I3V2","I4V2","I5V2","I6V2","WS1","WS2","WS3"]
	//}`
	//
	//	expectedCountRego := `package main
	//
	//import future.keywords.if
	//import future.keywords.in
	//tfplan := input if {
	//     input.terraform_version
	//} else := input.plan if {
	//     input.plan.terraform_version
	//}
	//
	//r := tfplan.resource_changes[_]
	//
	//warn if {
	// condition1
	//}
	//condition1 if {
	//r.type == "azurerm_app_service_environment_v3"
	//regex.match("ASE*",r.kind)
	//count({x|r.change.after.cluster_setting[x];condition1(x)}) < 1
	//}
	//condition1(x) if {
	//condition1(x)
	//}
	//condition1(x) if {
	//r.change.after.cluster_setting[x].name == "DisableTls1.0"
	//r.change.after.cluster_setting[x].value == "1"
	//}`
	//
	//	expectedListRego := `package main
	//
	//import rego.v1
	//
	//r := tfplan.resource_changes[_]
	//
	//warn if {
	// condition1
	//}
	//condition1 if {
	//r.type == "azurerm_healthcare_service"
	//not condition1
	//}
	//condition1 if {
	//r.change.after.cors_configuration[0].allowed_origins[0] != "*"
	//}`

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
				"deny.json": deny_json,
			},
			generatedRegoFileName: "deny.rego",
			input: map[string]any{
				"terraform_version": "1.11.0",
				"resource_changes": []any{
					map[string]any{
						"type": "azurerm_resource_group",
						"change": map[string]any{
							"after": map[string]any{
								"properties": map[string]any{},
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
				"deny.json": deny_json,
			},
			generatedRegoFileName: "deny.rego",
			input: map[string]any{
				"terraform_version": "1.11.0",
				"resource_changes": []any{
					map[string]any{
						"type": "Microsoft.Web/serverFarms",
						"change": map[string]any{
							"after": map[string]any{
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
			deny: true,
		},
		{
			desc:         "deny.json_deny_sku_name",
			inputDirPath: "",
			mockFs: map[string]string{
				"deny.json": deny_json,
			},
			generatedRegoFileName: "deny.rego",
			input: map[string]any{
				"terraform_version": "1.11.0",
				"resource_changes": []any{
					map[string]any{
						"type": "Microsoft.Web/serverFarms",
						"change": map[string]any{
							"after": map[string]any{
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
			deny: true,
		},
		{
			desc:         "deny.json_allow_full_input",
			inputDirPath: "",
			mockFs: map[string]string{
				"deny.json": deny_json,
			},
			generatedRegoFileName: "deny.rego",
			input: map[string]any{
				"terraform_version": "1.11.0",
				"resource_changes": []any{
					map[string]any{
						"type": "Microsoft.Web/serverFarms",
						"change": map[string]any{
							"after": map[string]any{
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
			deny: true,
		},
		{
			desc:         "count_deny",
			inputDirPath: "",
			mockFs: map[string]string{
				"count.json": count_json,
			},
			generatedRegoFileName: "count.rego",
			input: map[string]any{
				"terraform_version": "1.11.0",
				"resource_changes": []any{
					map[string]any{
						"type": "Microsoft.Web/hostingEnvironments",
						"kind": "ASE1",
						"change": map[string]any{
							"after": map[string]any{
								"properties": map[string]any{
									"clusterSettings": []any{
										map[string]any{
											"name":  "DisableTls1.0",
											"value": 1,
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

		////{
		////	desc:         "policy contains nested operations",
		////	inputDirPath: "",
		////	mockFs: map[string]string{
		////		"nested.json": nested_json,
		////	},
		////	expected: map[string]string{
		////		"nested.rego": expectedNestedRego,
		////	},
		////},
		//{
		//	desc:         "policy contains lists with multiple indexes",
		//	inputDirPath: "",
		//	mockFs: map[string]string{
		//		"list.json": list_field_json,
		//	},
		//	expected: map[string]string{
		//		"list.rego": expectedListRego,
		//	},
		//},
	}

	//t.Skip("Skip this test because it's not working on CI")
	//for i := 0; i < 10; i++ {
	for _, c := range cases {
		t.Run(fmt.Sprintf("%s", c.desc), func(t *testing.T) {
			mockFs := fakeFs(c.mockFs)
			counter := 0
			stub := gostub.Stub(&operation.NeoConditionNameGenerator, func(ctx *shared.Context) string {
				newName := fmt.Sprintf("condition%d", counter)
				counter++
				return newName
			}).Stub(&Fs, mockFs)
			defer stub.Reset()
			policyPath := ""
			if len(c.mockFs) == 1 {
				for n, _ := range c.mockFs {
					policyPath = n
				}
			}
			require.NoError(t, AzurePolicyToRego(policyPath, c.inputDirPath, shared.NewContext()))
			content, err := afero.ReadFile(mockFs, c.generatedRegoFileName)
			require.NoError(t, err)
			ctx := shared.NewContext()
			shared.AssertRego(t, "data.main.deny", string(content), c.input, c.deny, ctx)
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

func TestMapEffectToAction(t *testing.T) {
	tests := []struct {
		name          string
		thenBody      *ThenBody
		defaultEffect string
		expected      string
		expectError   bool
	}{
		{
			name: "Effect is deny",
			thenBody: &ThenBody{
				Effect: "deny",
			},
			defaultEffect: "",
			expected:      "deny",
		},
		{
			name: "Effect is [parameters('effect')] and defaultEffect is deny",
			thenBody: &ThenBody{
				Effect: "[parameters('effect')]",
			},
			defaultEffect: "deny",
			expected:      "deny",
		},
		{
			name: "Effect is [parameters('effect')] and defaultEffect is audit",
			thenBody: &ThenBody{
				Effect: "[parameters('effect')]",
			},
			defaultEffect: "audit",
			expected:      "warn",
		},
		{
			name: "Effect is [parameters('effect')] and defaultEffect is disabled",
			thenBody: &ThenBody{
				Effect: "[parameters('effect')]",
			},
			defaultEffect: "disabled",
			expected:      "disabled",
		},
		{
			name: "Effect is empty and defaultEffect is deny",
			thenBody: &ThenBody{
				Effect: "",
			},
			defaultEffect: "deny",
			expected:      "",
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.thenBody.MapEffectToAction(tt.defaultEffect)
			if tt.expectError {
				assert.NotNil(t, err)
				return
			}
			assert.Equal(t, tt.expected, result)
		})
	}
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
	t.Run("LoadRule", func(t *testing.T) {
		fs := prepareMemFs(t)
		counter := 1
		stub := gostub.Stub(&Fs, fs).Stub(&operation.NeoConditionNameGenerator, func(ctx *shared.Context) string {
			defer func() {
				counter++
			}()
			return fmt.Sprintf("condition%d", counter)
		})
		defer stub.Reset()

		ctx := shared.NewContext()
		rule, err := LoadRule(path, ctx)
		require.NoError(t, err)
		err = rule.SaveToDisk()
		require.NoError(t, err)
		file, err := afero.ReadFile(fs, "deny.rego")
		require.NoError(t, err)
		expected := `package main

import future.keywords.if
import future.keywords.in

tfplan := input if {
    input.terraform_version
} else := input.plan if {
    input.plan.terraform_version
}

r := tfplan.resource_changes[_]

warn if {
    condition2
}
condition2 if {
    r.type == "Microsoft.Web/serverFarms"
    condition1
}
condition1 if {
    not r.change.after.properties.sku.tier in ["Basic","Standard","ElasticPremium","Premium","PremiumV2","Premium0V3","PremiumV3","PremiumMV3","Isolated","IsolatedV2","WorkflowStandard"]
    not r.change.after.properties.sku.name in ["B1","B2","B3","S1","S2","S3","EP1","EP2","EP3","P1","P2","P3","P1V2","P2V2","P3V2","P0V3","P1V3","P2V3","P3V3","P1MV3","P2MV3","P3MV3","P4MV3","P5MV3","I1","I2","I3","I1V2","I2V2","I3V2","I4V2","I5V2","I6V2","WS1","WS2","WS3"]
}

`
		formattedExpected, err := format.Source("temp.rego", []byte(expected))
		require.NoError(t, err)
		assert.Equal(t, string(formattedExpected), string(file))
	})
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
						Subject: operation.Count{
							Where: operation.Where{
								Condition: condition.Equals{
									BaseCondition: condition.BaseCondition{
										Subject: &value.FieldValue{
											Name: "Microsoft.Network/networkSecurityGroups/securityRules[x].direction",
										},
									},
									Value: "Inbound",
								},
								ConditionSetName: "condition1",
							},
							CountExp: "count({x|x:=r.change.after.properties.securityRules[_];condition1(x)})",
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
				IfBody: operation.NewNot("condition1", condition.NotEquals{
					BaseCondition: condition.BaseCondition{
						Subject: &value.FieldValue{
							Name: "Microsoft.HealthcareApis/services/corsConfiguration.origins[x]",
						},
					},
					Value: "*",
				}),
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
				IfBody: operation.NewAnyOf("condition1", []shared.Rego{
					operation.NewAnyOf("condition1", []shared.Rego{
						condition.Exists{
							BaseCondition: condition.BaseCondition{
								Subject: &value.FieldValue{
									Name: "Microsoft.Sql/servers/minimalTlsVersion",
								},
							},
							Value: false,
						},
						condition.Less{
							BaseCondition: condition.BaseCondition{
								Subject: &value.FieldValue{
									Name: "Microsoft.Sql/servers/minimalTlsVersion",
								},
							},
							Value: "1.2",
						},
					}),
					operation.NewAllOf("condition1", []shared.Rego{
						condition.Equals{
							BaseCondition: condition.BaseCondition{
								Subject: &value.FieldValue{
									Name: "type",
								},
							},
							Value: "Microsoft.Sql/servers",
						},
						condition.Exists{
							BaseCondition: condition.BaseCondition{
								Subject: &value.FieldValue{
									Name: "Microsoft.Sql/servers/minimalTlsVersion",
								},
							},
							Value: true,
						},
					}),
				}),
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
				IfBody: operation.NewAnyOf("condition1", []shared.Rego{
					condition.Exists{
						BaseCondition: condition.BaseCondition{
							Subject: &value.FieldValue{
								Name: "Microsoft.Sql/servers/minimalTlsVersion",
							},
						},
						Value: false,
					},
					condition.Less{
						BaseCondition: condition.BaseCondition{
							Subject: &value.FieldValue{
								Name: "Microsoft.Sql/servers/minimalTlsVersion",
							},
						},
						Value: "1.2",
					},
				}),
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
				IfBody: operation.NewAllOf("condition1", []shared.Rego{
					condition.Equals{
						BaseCondition: condition.BaseCondition{
							Subject: &value.FieldValue{
								Name: "type",
							},
						},
						Value: "Microsoft.HealthcareApis/services",
					},
					condition.Exists{
						BaseCondition: condition.BaseCondition{
							Subject: &value.FieldValue{
								Name: "Microsoft.HealthcareApis/services/cosmosDbConfiguration.keyVaultKeyUri",
							},
						},
						Value: false,
					},
				}),
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
						Subject: &value.FieldValue{
							Name: "type",
						},
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
						Subject: &value.FieldValue{Name: "type"},
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
						Subject: &value.FieldValue{Name: "type"},
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
						Subject: &value.FieldValue{Name: "type"},
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
						Subject: &value.FieldValue{Name: "type"},
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
						Subject: &value.FieldValue{Name: "type"},
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
						Subject: &value.FieldValue{Name: "type"},
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
						Subject: &value.FieldValue{Name: "type"},
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
						Subject: &value.FieldValue{Name: "type"},
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
						Subject: &value.FieldValue{Name: "type"},
					},
					KeyName: "Microsoft.Web/serverFarms",
				},
			},
		},
		{
			name: "Less",
			input: map[string]any{
				"field": "number",
				"less":  10,
			},
			expected: &PolicyRuleBody{
				IfBody: condition.Less{
					BaseCondition: condition.BaseCondition{
						Subject: &value.FieldValue{Name: "number"},
					},
					Value: 10,
				},
			},
		},
		{
			name: "LessOrEquals",
			input: map[string]any{
				"field":        "number",
				"lessOrEquals": 10,
			},
			expected: &PolicyRuleBody{
				IfBody: condition.LessOrEquals{
					BaseCondition: condition.BaseCondition{
						Subject: &value.FieldValue{Name: "number"},
					},
					Value: 10,
				},
			},
		},
		{
			name: "Greater",
			input: map[string]any{
				"field":   "number",
				"greater": 10,
			},
			expected: &PolicyRuleBody{
				IfBody: condition.Greater{
					BaseCondition: condition.BaseCondition{
						Subject: &value.FieldValue{Name: "number"},
					},
					Value: 10,
				},
			},
		},
		{
			name: "GreaterOrEquals",
			input: map[string]any{
				"field":           "number",
				"greaterOrEquals": 10,
			},
			expected: &PolicyRuleBody{
				IfBody: condition.GreaterOrEquals{
					BaseCondition: condition.BaseCondition{
						Subject: &value.FieldValue{Name: "number"},
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
						Subject: &value.FieldValue{Name: "type"},
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
				stub := gostub.Stub(&operation.NeoConditionNameGenerator, func(ctx *shared.Context) string {
					return "condition1"
				})
				defer stub.Reset()
				result := NewPolicyRuleBody(tt.input, shared.NewContext())
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func prepareMemFs(t *testing.T) afero.Fs {
	fs := afero.NewMemMapFs()
	files := []string{
		"deny.json",
		"output.json",
		"rules.json",
	}
	for _, file := range files {
		content, err := os.ReadFile(file)
		require.NoError(t, err)
		err = afero.WriteFile(fs, file, content, os.ModePerm)
		require.NoError(t, err)
	}
	return fs
}
