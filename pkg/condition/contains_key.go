package condition

import (
	"context"
	"fmt"
	"json-rule-finder/pkg/shared"
)

var _ shared.Rego = ContainsKey{}

type ContainsKey struct {
	BaseCondition
	KeyName string
}

func (c ContainsKey) Rego(context.Context) (string, error) {
	return "", fmt.Errorf("`containsKey` BaseCondition is not supported, yet")
}
