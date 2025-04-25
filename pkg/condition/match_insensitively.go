package condition

import (
	"fmt"

	"github.com/tfmodtest/azopaform/pkg/shared"
)

var _ shared.Rego = MatchInsensitivelyCondition{}

type MatchInsensitivelyCondition struct {
	BaseCondition
	Value string
}

func (m MatchInsensitivelyCondition) Rego(*shared.Context) (string, error) {
	return "", fmt.Errorf("`matchInsensitively` condition is not supported, yet")
}
