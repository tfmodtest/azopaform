package shared

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFieldNameProcessor(t *testing.T) {
	ctx := NewContext()
	ctx.PushResourceType("Microsoft.Network/networkSecurityGroups")

	rego, err := FieldNameProcessor("Microsoft.Network/networkSecurityGroups/securityRules[*].direction", ctx)
	require.NoError(t, err)
	assert.Equal(t, "r.change.after.properties.securityRules[_].direction", rego)
}
