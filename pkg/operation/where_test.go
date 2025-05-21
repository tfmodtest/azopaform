package operation

import (
	"testing"

	"github.com/prashantv/gostub"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tfmodtest/azopaform/pkg/condition"
	"github.com/tfmodtest/azopaform/pkg/shared"
)

func TestNewWhere(t *testing.T) {
	stub := gostub.Stub(&RandomHelperFunctionNameGenerator, func() string {
		return "condition1"
	})
	defer stub.Reset()
	where, err := NewWhere(map[string]any{
		"field":  "Microsoft.Network/networkSecurityGroups/securityRules[*].direction",
		"equals": "Inbound",
	}, shared.NewContext())
	require.NoError(t, err)
	expected := Where{
		Condition: condition.Equals{
			BaseCondition: condition.BaseCondition{
				Subject: condition.FieldValue{
					Name: "Microsoft.Network/networkSecurityGroups/securityRules[*].direction",
				},
			},
			Value: "Inbound",
		},
		helperFunctionName: "condition1",
	}
	assert.Equal(t, expected, where)
}
