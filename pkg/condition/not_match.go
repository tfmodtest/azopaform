package condition

import (
	"fmt"

	"github.com/tfmodtest/azopaform/pkg/shared"
)

var _ shared.Rego = NotMatch{}

type NotMatch struct {
	BaseCondition
	Value string
}

func (n NotMatch) Rego(*shared.Context) (string, error) {
	return "", fmt.Errorf("`notMatch` condition is not supported, yet")
}
