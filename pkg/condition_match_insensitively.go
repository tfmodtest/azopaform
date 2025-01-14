package pkg

import (
	"context"
	"fmt"
	"json-rule-finder/pkg/shared"
)

var _ shared.Rego = MatchInsensitivelyCondition{}

type MatchInsensitivelyCondition struct {
	BaseCondition
	Value string
}

func (m MatchInsensitivelyCondition) Rego(context.Context) (string, error) {
	return "", fmt.Errorf("`matchInsensitively` BaseCondition is not supported, yet")
}
