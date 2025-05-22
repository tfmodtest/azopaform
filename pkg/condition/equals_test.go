package condition

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tfmodtest/azopaform/pkg/shared"
)

func TestEqualsCondition(t *testing.T) {
	cases := []struct {
		desc  string
		left  shared.Rego
		right any
		allow bool
	}{
		{
			desc:  "string",
			left:  shared.StringRego(`"right"`),
			right: "right",
			allow: true,
		},
		{
			desc:  "string_negative",
			left:  shared.StringRego(`"left"`),
			right: "right",
			allow: false,
		},
		{
			desc:  "int",
			left:  shared.StringRego("1"),
			right: 1,
			allow: true,
		},
		{
			desc:  "int_negative",
			left:  shared.StringRego("1"),
			right: 2,
			allow: false,
		},
		{
			desc:  "bool",
			left:  shared.StringRego("true"),
			right: true,
			allow: true,
		},
		{
			desc:  "bool_negative",
			left:  shared.StringRego("false"),
			right: true,
			allow: false,
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			ctx := shared.NewContext()
			sut := Equals{
				BaseCondition: BaseCondition{
					Subject: c.left,
				},
				Value: c.right,
			}
			actual, err := sut.Rego(ctx)
			require.NoError(t, err)
			cfg := fmt.Sprintf(shared.RegoTestTemplate, actual)
			shared.AssertRegoAllow(t, cfg, nil, c.allow, ctx)
		})
	}
}

func TestEqualsCondition_SpecialCase_type(t *testing.T) {
	ctx := shared.NewContext()
	sut := Equals{
		BaseCondition: BaseCondition{
			Subject: FieldValue{Name: "type"},
		},
		Value: "Microsoft.Network/networkSecurityGroups",
	}
	actual, err := sut.Rego(ctx)
	require.NoError(t, err)
	assert.Equal(t, `is_azure_type(r.values, "Microsoft.Network/networkSecurityGroups")`, actual)
}

func TestEqualsCondition_WithUtilLibraryPackageName(t *testing.T) {
	ctx := shared.NewContextWithOptions(shared.Options{UtilLibraryPackageName: "util"})

	sut := Equals{
		BaseCondition: BaseCondition{
			Subject: FieldValue{Name: "type"},
		},
		Value: "Microsoft.Network/networkSecurityGroups",
	}

	actual, err := sut.Rego(ctx)
	require.NoError(t, err)
	expected := `data.util.is_azure_type(r.values, "Microsoft.Network/networkSecurityGroups")`
	assert.Equal(t, expected, actual)
}

func TestEqualsCondition_SingleQuoteInFieldShouldBeReplacedByDoubleQuote(t *testing.T) {
	ctx := shared.NewContext()
	sut := Equals{
		BaseCondition: BaseCondition{
			Subject: FieldValue{Name: "tags['aks-managed-poolName']"},
		},
		Value: `1`,
	}
	actual, err := sut.Rego(ctx)
	require.NoError(t, err)
	assert.Equal(t, `tags["aks-managed-poolName"] == "1"`, actual)
}

func TestEqualsCondition_WildcardInFieldPathShouldBeEvalAsAllOf(t *testing.T) {
	ctx := shared.NewContext()
	ctx.PushResourceType("Microsoft.Network/networkSecurityGroups/securityRules")

	sut := Equals{
		BaseCondition: BaseCondition{
			Subject: FieldValue{Name: "Microsoft.Network/networkSecurityGroups/securityRules/port[*]"},
		},
		Value: 3,
	}

	actual, err := sut.Rego(ctx)
	require.NoError(t, err)
	cfg := fmt.Sprintf(shared.RegoTestTemplate, actual)
	cfg = cfg + "\n" + `r := input`
	shared.AssertRegoAllow(t, cfg, map[string]any{
		"values": map[string]any{
			"properties": map[string]any{
				"port": []int{
					3,
					3,
				},
			},
		},
	}, true, ctx)
	shared.AssertRegoAllow(t, cfg, map[string]any{
		"values": map[string]any{
			"properties": map[string]any{
				"port": []int{
					3,
					4,
				},
			},
		},
	}, false, ctx)
}
