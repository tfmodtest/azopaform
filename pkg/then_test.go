package pkg_test

import (
	"github.com/tfmodtest/azopaform/pkg"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapEffectToAction(t *testing.T) {
	tests := []struct {
		name          string
		thenBody      *pkg.ThenBody
		defaultEffect string
		expected      string
		expectError   bool
	}{
		{
			name: "Effect is deny",
			thenBody: &pkg.ThenBody{
				Effect: "deny",
			},
			defaultEffect: "",
			expected:      "deny",
		},
		{
			name: "Effect is [parameters('effect')] and defaultEffect is deny",
			thenBody: &pkg.ThenBody{
				Effect: "[parameters('effect')]",
			},
			defaultEffect: "deny",
			expected:      "deny",
		},
		{
			name: "Effect is [parameters('effect')] and defaultEffect is audit",
			thenBody: &pkg.ThenBody{
				Effect: "[parameters('effect')]",
			},
			defaultEffect: "audit",
			expected:      "warn",
		},
		{
			name: "Effect is [parameters('effect')] and defaultEffect is disabled",
			thenBody: &pkg.ThenBody{
				Effect: "[parameters('effect')]",
			},
			defaultEffect: "disabled",
			expected:      "deny",
		},
		{
			name: "Effect is empty and defaultEffect is deny",
			thenBody: &pkg.ThenBody{
				Effect: "",
			},
			defaultEffect: "deny",
			expected:      "",
			expectError:   true,
		},
		{
			name: "Effect is [parameters('effect')] and defaultEffect is Modify",
			thenBody: &pkg.ThenBody{
				Effect: "[parameters('effect')]",
			},
			defaultEffect: "Modify",
			expected:      "deny",
		},
		{
			name: "Effect is [parameters('effect')] and defaultEffect is Modify",
			thenBody: &pkg.ThenBody{
				Effect: "[parameters('effect')]",
			},
			defaultEffect: "DeployIfNotExists",
			expected:      "deployifnotexists",
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

func TestThenBody_Action(t *testing.T) {
	tests := []struct {
		name               string
		thenBody           *pkg.ThenBody
		ruleName           string
		result             string
		helperFunctionName string
		rule               *pkg.Rule
		expected           string
		expectError        bool
	}{
		{
			name: "deny action with helper function",
			thenBody: &pkg.ThenBody{
				Effect: "[parameters('effect')]",
			},
			ruleName:           "test_rule",
			result:             "input.value == true",
			helperFunctionName: "helper_func",
			rule: &pkg.Rule{
				Properties: &pkg.PolicyRuleModel{
					Parameters: &pkg.PolicyRuleParameters{
						Effect: &pkg.EffectBody{
							DefaultValue: "deny",
						},
					},
				},
			},
			expected:    "deny_test_rule if {\n  res := resource(input, \"azapi_resource\")[_]\n helper_func(res)\n}\ninput.value == true",
			expectError: false,
		},
		{
			name: "warn action without helper function",
			thenBody: &pkg.ThenBody{
				Effect: "[parameters('effect')]",
			},
			ruleName:           "warn_rule",
			result:             "input.value == false",
			helperFunctionName: "",
			rule: &pkg.Rule{
				Properties: &pkg.PolicyRuleModel{
					Parameters: &pkg.PolicyRuleParameters{
						Effect: &pkg.EffectBody{
							DefaultValue: "audit",
						},
					},
				},
			},
			expected:    "warn_warn_rule if {\n input.value == false\n}\n",
			expectError: false,
		},
		{
			name: "DeployIfNotExists action without helper function",
			thenBody: &pkg.ThenBody{
				Effect: "[parameters('effect')]",
			},
			ruleName:           "rule1",
			result:             "input.value == false",
			helperFunctionName: "",
			rule: &pkg.Rule{
				Properties: &pkg.PolicyRuleModel{
					Parameters: &pkg.PolicyRuleParameters{
						Effect: &pkg.EffectBody{
							DefaultValue: "deployIfNotExists",
						},
					},
				},
			},
			expected:    "deny_rule1 if {\n not input.value == false\n}\n",
			expectError: false,
		},
		{
			name: "DeployIfNotExists action with helper function",
			thenBody: &pkg.ThenBody{
				Effect: "[parameters('effect')]",
			},
			ruleName:           "test_rule",
			result:             "input.value == true",
			helperFunctionName: "helper_func",
			rule: &pkg.Rule{
				Properties: &pkg.PolicyRuleModel{
					Parameters: &pkg.PolicyRuleParameters{
						Effect: &pkg.EffectBody{
							DefaultValue: "DeployIfNotExists",
						},
					},
				},
			},
			expected:    "deny_test_rule if {\n  res := resource(input, \"azapi_resource\")[_]\n not helper_func(res)\n}\ninput.value == true",
			expectError: false,
		},
		{
			name: "unexpected effect error",
			thenBody: &pkg.ThenBody{
				Effect: "invalid_effect",
			},
			ruleName:           "invalid_rule",
			result:             "input.value == false",
			helperFunctionName: "",
			rule: &pkg.Rule{
				Properties: &pkg.PolicyRuleModel{
					Parameters: &pkg.PolicyRuleParameters{
						Effect: &pkg.EffectBody{
							DefaultValue: "deny",
						},
					},
				},
			},
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := tt.thenBody.Action(tt.ruleName, tt.result, tt.helperFunctionName, tt.rule)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, actual)
			}
		})
	}
}
