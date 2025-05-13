package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tfmodtest/azopaform/pkg/shared"
)

func Test_RuleNameShouldBeDisplayNameInSnakeCase(t *testing.T) {
	rule := &Rule{
		Properties: &PolicyRuleModel{
			PolicyRule: &PolicyRuleBody{
				If: map[string]any{
					"value":  "1",
					"equals": 1,
				},
				Then: &ThenBody{Effect: "deny"},
			},
			DisplayName: "This is a test",
		},
	}
	err := rule.Parse(shared.NewContextWithOptions(shared.Options{
		GenerateRuleName: true,
	}))
	require.NoError(t, err)
	assert.Equal(t, "this_is_a_test", rule.Name)
}

func Test_RuleNameShouldFollowAction(t *testing.T) {
	rule := &Rule{
		Properties: &PolicyRuleModel{
			PolicyRule: &PolicyRuleBody{
				If: map[string]any{
					"value":  "1",
					"equals": 1,
				},
				Then: &ThenBody{Effect: "deny"},
			},
			DisplayName: "rule_1",
		},
	}
	err := rule.Parse(shared.NewContextWithOptions(shared.Options{
		GenerateRuleName: true,
	}))
	require.NoError(t, err)
	assert.Contains(t, rule.result, "deny_rule_1")
}

func Test_ParseParameters_FullParameters(t *testing.T) {
	// Arrange
	rule := newRule()
	input := map[string]any{
		"properties": map[string]any{
			"parameters": map[string]any{
				"allowedClusterAutoUpgradeChannels": map[string]any{
					"type": "Array",
					"defaultValue": []any{
						"rapid", "stable", "patch",
					},
					"metadata": map[string]any{
						"displayName":  "Allowed Cluster Auto-upgrade Channels",
						"description":  "Cluster auto-upgrade channels viewed as complaint.",
						"portalReview": true,
					},
				},
			},
		},
	}

	// Act
	rule.ParseParameters(input)

	// Assert
	require.NotNil(t, rule.Properties.Parameters)
	require.NotNil(t, rule.Properties.Parameters.Parameters)

	param := rule.Properties.Parameters.Parameters["allowedClusterAutoUpgradeChannels"]
	require.NotNil(t, param)
	assert.Equal(t, "allowedClusterAutoUpgradeChannels", param.Name)
	assert.Equal(t, PolicyRuleParameterType("Array"), param.Type)
	assert.Equal(t, []any{"rapid", "stable", "patch"}, param.DefaultValue)

	require.NotNil(t, param.MetaData)
	assert.Equal(t, "Allowed Cluster Auto-upgrade Channels", param.MetaData.DisplayName)
	assert.Equal(t, "Cluster auto-upgrade channels viewed as complaint.", param.MetaData.Description)
	assert.False(t, param.MetaData.Deprecated)
}

func Test_ParseParameters_MissingMetadata(t *testing.T) {
	// Arrange
	rule := newRule()
	input := map[string]any{
		"properties": map[string]any{
			"parameters": map[string]any{
				"simpleParam": map[string]any{
					"type":         "String",
					"defaultValue": "default",
				},
			},
		},
	}

	// Act
	rule.ParseParameters(input)

	// Assert
	require.NotNil(t, rule.Properties.Parameters)
	param := rule.Properties.Parameters.Parameters["simpleParam"]
	require.NotNil(t, param)
	assert.Equal(t, "simpleParam", param.Name)
	assert.Equal(t, PolicyRuleParameterType("String"), param.Type)
	assert.Equal(t, "default", param.DefaultValue)
	assert.Nil(t, param.MetaData)
}

func Test_ParseParameters_MissingProperties(t *testing.T) {
	// Arrange
	rule := newRule()
	input := map[string]any{
		"someOtherField": "value",
	}

	// Act
	rule.ParseParameters(input)

	// Assert
	assert.Empty(t, rule.Properties.Parameters.Parameters)
}

func Test_ParseParameters_MissingParametersSection(t *testing.T) {
	// Arrange
	rule := newRule()
	input := map[string]any{
		"properties": map[string]any{
			"someOtherField": "value",
		},
	}

	// Act
	rule.ParseParameters(input)

	// Assert
	assert.Empty(t, rule.Properties.Parameters.Parameters)
}
