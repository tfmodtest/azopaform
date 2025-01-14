package pkg

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"json-rule-finder/pkg/shared"
	"testing"
)

func TestFieldNameProcessor(t *testing.T) {
	ctx := shared.NewContext()
	shared.PushResourceType(ctx, "Microsoft.Web/serverFarms")
	subject, rules, err := shared.FieldNameProcessor("Microsoft.Web/serverFarms/sku.tier", ctx)
	require.NoError(t, err)
	assert.Equal(t, "", rules)
	assert.Equal(t, "r.change.after.sku[0].tier", subject)
}
