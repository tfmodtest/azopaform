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
