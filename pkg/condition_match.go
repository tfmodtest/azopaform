package pkg

import (
	"context"
	"fmt"
	"json-rule-finder/pkg/shared"
)

var _ shared.Rego = MatchCondition{}

type MatchCondition struct {
	BaseCondition
	Value string
}

func (m MatchCondition) Rego(ctx context.Context) (string, error) {
	return "", fmt.Errorf("`match` BaseCondition is not supported, yet")
}
