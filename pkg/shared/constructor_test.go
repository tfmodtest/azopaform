package shared

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFieldNameProcessor(t *testing.T) {
	ctx := NewContext()
	ctx.PushResourceType("Microsoft.Network/networkSecurityGroups")

	rego, err := FieldNameProcessor("Microsoft.Network/networkSecurityGroups/securityRules[*].direction", ctx)
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("%s.properties.securityRules[_].direction", ResourcePathPrefix), rego)
}
