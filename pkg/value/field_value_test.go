package value

import (
	"json-rule-finder/pkg/shared"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFieldValue(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedName   string
		expectedRego   string
		resourceType   string
		shouldPushType bool
	}{
		{
			name:         "Simple field name",
			input:        "type",
			expectedName: "type",
			expectedRego: "r.values.type",
		},
		{
			name:         "Kind field name",
			input:        "kind",
			expectedName: "kind",
			expectedRego: "r.values.kind",
		},
		{
			name:           "Field name with single [*]",
			input:          "Microsoft.Web/serverFarms/sku[*].tier",
			expectedName:   "Microsoft.Web/serverFarms/sku[x].tier",
			expectedRego:   "r.values.properties.sku[_].tier",
			resourceType:   "Microsoft.Web/serverFarms",
			shouldPushType: true,
		},
		{
			name:           "Field name with multiple [*]",
			input:          "Microsoft.Network/networkSecurityGroups/securityRules[*]/properties[*].direction",
			expectedName:   "Microsoft.Network/networkSecurityGroups/securityRules[x]/properties[x].direction",
			expectedRego:   "r.values.properties.securityRules[_].properties[_].direction",
			resourceType:   "Microsoft.Network/networkSecurityGroups",
			shouldPushType: true,
		},
		{
			name:         "Field name containing count",
			input:        "Microsoft.Network/networkSecurityGroups/count",
			expectedName: "Microsoft.Network/networkSecurityGroups/count",
			expectedRego: "Microsoft.Network/networkSecurityGroups/count",
		},
		{
			name:         "Field name with count and wildcards",
			input:        "Microsoft.Network/networkSecurityGroups/count[*].value",
			expectedName: "Microsoft.Network/networkSecurityGroups/count[x].value",
			expectedRego: "Microsoft.Network/networkSecurityGroups/count[x].value",
		},
		{
			name:         "Empty string",
			input:        "",
			expectedName: "",
			expectedRego: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := shared.NewContext()
			if tt.shouldPushType && tt.resourceType != "" {
				ctx.PushResourceType(tt.resourceType)
			}

			// Test that NewFieldValue correctly sets the Name field
			result := NewFieldValue(tt.input, ctx)
			fieldValue, ok := result.(FieldValue)
			assert.True(t, ok, "Result should be a FieldValue")
			assert.Equal(t, tt.expectedName, fieldValue.Name)

			// Also test the Rego() method to ensure proper field name processing
			regoResult, err := fieldValue.Rego(ctx)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedRego, regoResult)
		})
	}
}
