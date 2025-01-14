package pkg

import (
	"context"
	"fmt"
	"json-rule-finder/pkg/shared"
)

var _ shared.Rego = NotMatchCondition{}

type NotMatchCondition struct {
	BaseCondition
	Value string
}

func (n NotMatchCondition) Rego(context.Context) (string, error) {
	return "", fmt.Errorf("`notMatch` BaseCondition is not supported, yet")
}
