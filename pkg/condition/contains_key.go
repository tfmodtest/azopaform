package condition

import (
	"fmt"
	"json-rule-finder/pkg/shared"
)

var _ shared.Rego = ContainsKey{}

type ContainsKey struct {
	BaseCondition
	KeyName string
}

func (c ContainsKey) Rego(*shared.Context) (string, error) {
	return "", fmt.Errorf("`containsKey` BaseCondition is not supported, yet")
}
