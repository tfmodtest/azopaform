package condition

import (
	"fmt"

	"github.com/tfmodtest/azopaform/pkg/shared"
)

var _ shared.Rego = NotMatchInsensitively{}

func (n NotMatchInsensitively) Rego(*shared.Context) (string, error) {
	return "", fmt.Errorf("`notMatchInsensitively` condition is not supported, yet")
}

type NotMatchInsensitively struct {
	BaseCondition
	Value string
}
