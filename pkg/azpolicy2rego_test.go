package pkg

import (
	"fmt"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"

	"github.com/prashantv/gostub"
	"github.com/stretchr/testify/require"
)

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
        "defaultValue": "Audit",
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
        "defaultValue": "Audit"
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
            "allOf": [
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
	rulesJson, err := os.ReadFile("rules.json")
	require.NoError(t, err)
	outputJson, err := os.ReadFile("output.json")
	require.NoError(t, err)

	expectedDenyRego := `package main

import rego.v1

r := tfplan.resource_changes[_]

warn if {
 aaaaa
}
aaaaa if {
 type == azurerm_service_plan
 aaaaa
}
aaaaa if {
 not r.change.after.sku[0].tier in ["Basic","Standard","ElasticPremium","Premium","PremiumV2","Premium0V3","PremiumV3","PremiumMV3","Isolated","IsolatedV2","WorkflowStandard"]
 not r.change.after.sku_name in ["B1","B2","B3","S1","S2","S3","EP1","EP2","EP3","P1","P2","P3","P1V2","P2V2","P3V2","P0V3","P1V3","P2V3","P3V3","P1MV3","P2MV3","P3MV3","P4MV3","P5MV3","I1","I2","I3","I1V2","I2V2","I3V2","I4V2","I5V2","I6V2","WS1","WS2","WS3"]
}
`

	expectedCountRego := `package main

import rego.v1

r := tfplan.resource_changes[_]

warn if {
 aaaaa
}
aaaaa if {
 type == azurerm_app_service_environment_v3
 regex.match("ASE*",kind)
 count({x|r.change.after.properties.Microsoft.Web.HostingEnvironments.clusterSettings[x];aaaaaaaaa(x)}) < 1
}
aaaaaaaaa(x) if {
 aaaaa(x)
}
aaaaa(x) if {
 r.change.after.properties.Microsoft.Web.HostingEnvironments.clusterSettings[x].name == DisableTls1.0
 r.change.after.properties.Microsoft.Web.HostingEnvironments.clusterSettings[x].value == 1
}
`

	cases := []struct {
		desc         string
		inputDirPath string
		mockFs       map[string]string
		expected     map[string]string
	}{
		{
			desc:         "deny.json",
			inputDirPath: "",
			mockFs: map[string]string{
				"deny.json": deny_json,
			},
			expected: map[string]string{
				"deny.rego": expectedDenyRego,
			},
		},
		{
			desc:         "dirPath",
			inputDirPath: "/config",
			mockFs: map[string]string{
				"/config/deny.json": deny_json,
			},
			expected: map[string]string{
				"deny.rego": expectedDenyRego,
			},
		},
		{
			desc:         "multiple json files in dirPath",
			inputDirPath: "/config",
			mockFs: map[string]string{
				"/config/deny1.json": deny_json,
				"/config/deny2.json": deny_json,
			},
			expected: map[string]string{
				"deny1.rego": expectedDenyRego,
				"deny2.rego": expectedDenyRego,
			},
		},
		{
			desc:         "json files in grandson's folders",
			inputDirPath: "/config",
			mockFs: map[string]string{
				"/config/deny1/deny1.json": deny_json,
				"/config/deny2/deny2.json": deny_json,
			},
			expected: map[string]string{
				"deny1.rego": expectedDenyRego,
				"deny2.rego": expectedDenyRego,
			},
		},
		{
			desc:         "policy contains count operator",
			inputDirPath: "",
			mockFs: map[string]string{
				"count.json": count_json,
			},
			expected: map[string]string{
				"count.rego": expectedCountRego,
			},
		},
	}

	//for i := 0; i < 10; i++ {
	for _, c := range cases {
		t.Run(fmt.Sprintf("%s", c.desc), func(t *testing.T) {
			files := map[string]string{
				"rules.json":  string(rulesJson),
				"output.json": string(outputJson),
			}
			for n, f := range c.mockFs {
				files[n] = f
			}
			mockFs := fakeFs(files)
			stub := gostub.Stub(&RandIntRange, func(min int, max int) int {
				return 0
			}).Stub(&Fs, mockFs)
			defer stub.Reset()
			policyPath := ""
			if len(c.mockFs) == 1 {
				for n, _ := range c.mockFs {
					policyPath = n
				}
			}
			require.NoError(t, AzurePolicyToRego(policyPath, c.inputDirPath, NewContext()))
			for n, expected := range c.expected {
				content, err := afero.ReadFile(mockFs, n)
				require.NoError(t, err)
				assert.Equal(t, expected, string(content))
			}
		})
	}
	//}
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
