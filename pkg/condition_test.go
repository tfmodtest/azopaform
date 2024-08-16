package pkg

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNotInCondition(t *testing.T) {
	sut := NotInOperation{
		operation: operation{
			Subject: OperationField("Microsoft.Web/serverFarms/sku.tier"),
		},
		Values: []string{
			"Basic",
			"Standard",
			"ElasticPremium",
			"Premium",
			"PremiumV2",
			"Premium0V3",
			"PremiumV3",
			"PremiumMV3",
			"Isolated",
			"IsolatedV2",
			"WorkflowStandard",
		},
	}
	ctx := NewContext()
	pushResourceType(ctx, "Microsoft.Web/serverFarms")
	actual, err := sut.Rego(ctx)
	require.NoError(t, err)
	assert.Equal(t, `not r.change.after.sku[0].tier in ["Basic","Standard","ElasticPremium","Premium","PremiumV2","Premium0V3","PremiumV3","PremiumMV3","Isolated","IsolatedV2","WorkflowStandard"]`, actual)
}

func TestInCondition(t *testing.T) {
	sut := InOperation{
		operation: operation{
			Subject: OperationField("Microsoft.Web/serverFarms/sku.tier"),
		},
		Values: []string{
			"Basic",
			"Standard",
			"ElasticPremium",
			"Premium",
			"PremiumV2",
			"Premium0V3",
			"PremiumV3",
			"PremiumMV3",
			"Isolated",
			"IsolatedV2",
			"WorkflowStandard",
		},
	}
	ctx := NewContext()
	pushResourceType(ctx, "Microsoft.Web/serverFarms")
	actual, err := sut.Rego(ctx)
	require.NoError(t, err)
	assert.Equal(t, `some r.change.after.sku[0].tier in ["Basic","Standard","ElasticPremium","Premium","PremiumV2","Premium0V3","PremiumV3","PremiumMV3","Isolated","IsolatedV2","WorkflowStandard"]`, actual)
}

func TestLikeCondition(t *testing.T) {
	sut := LikeOperation{
		operation: operation{
			Subject: OperationField("Microsoft.Web/serverFarms/sku.tier"),
		},
		Value: `^[^@]+@[^@]+\.[^@]+$`,
	}
	ctx := NewContext()
	pushResourceType(ctx, "Microsoft.Web/serverFarms")
	actual, err := sut.Rego(ctx)
	require.NoError(t, err)
	assert.Equal(t, "regex.match(^[^@]+@[^@]+\\.[^@]+$,r.change.after.sku[0].tier)", actual)
}

func TestNotLikeCondition(t *testing.T) {
	sut := NotLikeOperation{
		operation: operation{
			Subject: OperationField("Microsoft.Web/serverFarms/sku.tier"),
		},
		Value: `^[^@]+@[^@]+\.[^@]+$`,
	}
	ctx := NewContext()
	pushResourceType(ctx, "Microsoft.Web/serverFarms")
	actual, err := sut.Rego(ctx)
	require.NoError(t, err)
	assert.Equal(t, "not regex.match(^[^@]+@[^@]+\\.[^@]+$,r.change.after.sku[0].tier)", actual)
}

func TestEqualsCondition(t *testing.T) {
	sut := EqualsOperation{
		operation: operation{
			Subject: OperationField("Microsoft.Web/serverFarms/sku.tier"),
		},
		Value: "Standard",
	}
	ctx := NewContext()
	pushResourceType(ctx, "Microsoft.Web/serverFarms")
	actual, err := sut.Rego(ctx)
	require.NoError(t, err)
	assert.Equal(t, "r.change.after.sku[0].tier == Standard", actual)
}

func TestNotEqualsCondition(t *testing.T) {
	sut := NotEqualsOperation{
		operation: operation{
			Subject: OperationField("Microsoft.Web/serverFarms/sku.tier"),
		},
		Value: "Standard",
	}
	ctx := NewContext()
	pushResourceType(ctx, "Microsoft.Web/serverFarms")
	actual, err := sut.Rego(ctx)
	require.NoError(t, err)
	assert.Equal(t, "r.change.after.sku[0].tier != Standard", actual)
}
