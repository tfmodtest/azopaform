package pkg

import (
	"context"
	"fmt"
	"json-rule-finder/pkg/shared"
)

type NotContainsKeyCondition struct {
	BaseCondition
	KeyName string
}

var _ shared.Rego = NotContainsKeyCondition{}

func (n NotContainsKeyCondition) Rego(context.Context) (string, error) {
	return "", fmt.Errorf("`notContainsKey` BaseCondition is not supported, yet")
}
