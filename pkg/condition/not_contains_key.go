package condition

import (
	"fmt"
	"json-rule-finder/pkg/shared"
)

type NotContainsKey struct {
	BaseCondition
	KeyName string
}

var _ shared.Rego = NotContainsKey{}

func (n NotContainsKey) Rego(*shared.Context) (string, error) {
	return "", fmt.Errorf("`notContainsKey` BaseCondition is not supported, yet")
}
