package pkg

import (
	"json-rule-finder/pkg/shared"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOperationField_Rego(t *testing.T) {
	sut := OperationField("Microsoft.Network/networkSecurityGroups/securityRules[*].direction")
	ctx := shared.NewContext()
	ctx.PushResourceType("Microsoft.Network/networkSecurityGroups")
	rego, err := sut.Rego(ctx)
	require.NoError(t, err)
	assert.Equal(t, "r.change.after.securityRules[_].direction", rego)
}
