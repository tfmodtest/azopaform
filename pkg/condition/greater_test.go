package condition

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tfmodtest/azopaform/pkg/shared"
)

func TestGreaterCondition(t *testing.T) {
	cases := []struct {
		desc  string
		left  shared.Rego
		right string
		allow bool
	}{
		{
			desc:  "greater",
			left:  shared.StringRego("2"),
			right: "1",
			allow: true,
		},
		{
			desc:  "less",
			left:  shared.StringRego("1"),
			right: "2",
			allow: false,
		},
		{
			desc:  "equal",
			left:  shared.StringRego("1"),
			right: "1",
			allow: false,
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			ctx := shared.NewContext()
			sut := Greater{
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

func TestGreaterCondition_WithoutWildcardFieldPath(t *testing.T) {
	ctx := shared.NewContext()
	ctx.PushResourceType("Microsoft.Network/networkSecurityGroups/securityRules")

	sut := Greater{
		BaseCondition: BaseCondition{
			Subject: FieldValue{Name: "Microsoft.Network/networkSecurityGroups/securityRules/port"},
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
				"port": 4,
			},
		},
	}, true, ctx)
	shared.AssertRegoAllow(t, cfg, map[string]any{
		"values": map[string]any{
			"properties": map[string]any{
				"port": 3,
			},
		},
	}, false, ctx)
}

func TestGreaterCondition_WildcardInFieldPathShouldBeEvalAsAllOf(t *testing.T) {
	ctx := shared.NewContext()
	ctx.PushResourceType("Microsoft.Network/networkSecurityGroups/securityRules")

	sut := Greater{
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
					5,
					4,
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

func TestGreaterCondition_MultipleWildcardInFieldPathShouldBeEvalAsAllOf(t *testing.T) {
	cases := []struct {
		desc   string
		path   string
		inputs map[string]any
		allow  bool
	}{
		{
			desc: "multiple wildcard in field path",
			path: "Microsoft.Network/networkSecurityGroups/securityRules/profile[*]/port[*]",
			inputs: map[string]any{
				"values": map[string]any{
					"properties": map[string]any{
						"profile": []map[string]any{
							{
								"port": []int{
									5,
									4,
								},
							},
							{
								"port": []int{
									5,
									4,
								},
							},
						},
					},
				},
			},
			allow: true,
		},
		{
			desc: "multiple wildcard in field path disallow",
			path: "Microsoft.Network/networkSecurityGroups/securityRules/profile[*]/port[*]",
			inputs: map[string]any{
				"values": map[string]any{
					"properties": map[string]any{
						"profile": []map[string]any{
							{
								"port": []int{
									3,
									4,
								},
							},
							{
								"port": []int{
									5,
									4,
								},
							},
						},
					},
				},
			},
			allow: false,
		},
		{
			desc: "multiple wildcard spanning across fixed path segments",
			path: "Microsoft.Network/networkSecurityGroups/securityRules/profile[*]/network/port[*]",
			inputs: map[string]any{
				"values": map[string]any{
					"properties": map[string]any{
						"profile": []map[string]any{
							{
								"network": map[string]any{
									"port": []int{
										5,
										4,
									},
								},
							},
							{
								"network": map[string]any{
									"port": []int{
										5,
										4,
									},
								},
							},
						},
					},
				},
			},
			allow: true,
		},
		{
			desc: "multiple wildcard spanning across fixed path segments disallow",
			path: "Microsoft.Network/networkSecurityGroups/securityRules/profile[*]/network/port[*]",
			inputs: map[string]any{
				"values": map[string]any{
					"properties": map[string]any{
						"profile": []map[string]any{
							{
								"network": map[string]any{
									"port": []int{
										3,
										4,
									},
								},
							},
							{
								"network": map[string]any{
									"port": []int{
										5,
										4,
									},
								},
							},
						},
					},
				},
			},
			allow: false,
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			ctx := shared.NewContext()
			ctx.PushResourceType("Microsoft.Network/networkSecurityGroups/securityRules")

			sut := Greater{
				BaseCondition: BaseCondition{
					Subject: FieldValue{Name: c.path},
				},
				Value: 3,
			}

			actual, err := sut.Rego(ctx)
			require.NoError(t, err)
			cfg := fmt.Sprintf(shared.RegoTestTemplate, actual)
			cfg = cfg + "\n" + `r := input`
			shared.AssertRegoAllow(t, cfg, c.inputs, c.allow, ctx)
		})
	}
}
