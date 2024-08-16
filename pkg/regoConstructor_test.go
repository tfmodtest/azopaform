package pkg

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFieldNameProcessor(t *testing.T) {
	ctx := NewContext()
	pushResourceType(ctx, "Microsoft.Web/serverFarms")
	subject, rules, err := FieldNameProcessor("Microsoft.Web/serverFarms/sku.tier", ctx)
	require.NoError(t, err)
	assert.Equal(t, "", rules)
	assert.Equal(t, "r.change.after.sku[0].tier", subject)
}
