package pkg

import (
	"context"
	"fmt"
	"json-rule-finder/pkg/shared"
)

var _ shared.Rego = NotMatchInsensitivelyCondition{}

func (n NotMatchInsensitivelyCondition) Rego(context.Context) (string, error) {
	return "", fmt.Errorf("`notMatchInsensitively` BaseCondition is not supported, yet")
}

type NotMatchInsensitivelyCondition struct {
	BaseCondition
	Value string
}
