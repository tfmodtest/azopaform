package condition

import (
	"context"
	"fmt"
	"json-rule-finder/pkg/shared"
)

var _ shared.Rego = NotMatch{}

type NotMatch struct {
	BaseCondition
	Value string
}

func (n NotMatch) Rego(context.Context) (string, error) {
	return "", fmt.Errorf("`notMatch` condition is not supported, yet")
}
