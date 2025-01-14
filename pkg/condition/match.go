package condition

import (
	"context"
	"fmt"
	"json-rule-finder/pkg/shared"
)

var _ shared.Rego = Match{}

type Match struct {
	BaseCondition
	Value string
}

func (m Match) Rego(ctx context.Context) (string, error) {
	return "", fmt.Errorf("`match` BaseCondition is not supported, yet")
}
