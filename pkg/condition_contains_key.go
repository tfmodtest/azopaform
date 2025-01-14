package pkg

import (
	"context"
	"fmt"
	"json-rule-finder/pkg/shared"
)

var _ shared.Rego = ContainsKeyCondition{}

type ContainsKeyCondition struct {
	BaseCondition
	KeyName string
}

func (c ContainsKeyCondition) Rego(context.Context) (string, error) {
	return "", fmt.Errorf("`containsKey` BaseCondition is not supported, yet")
}
