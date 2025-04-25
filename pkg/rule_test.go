package pkg

import (
	"testing"

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
	err := rule.Parse(shared.NewContext())
	require.NoError(t, err)
	require.Equal(t, "this_is_a_test", rule.Name)
}
