package condition

import (
	"context"
	"fmt"
	"json-rule-finder/pkg/shared"
)

var _ shared.Rego = NotMatchInsensitively{}

func (n NotMatchInsensitively) Rego(context.Context) (string, error) {
	return "", fmt.Errorf("`notMatchInsensitively` condition is not supported, yet")
}

type NotMatchInsensitively struct {
	BaseCondition
	Value string
}
