package condition

import (
	"fmt"

	"github.com/tfmodtest/azopaform/pkg/shared"
)

var _ Condition = Exists{}

type Exists struct {
	BaseCondition
	Value any
}

func (e Exists) Rego(ctx *shared.Context) (string, error) {
	fieldName, err := e.GetSubject(ctx).Rego(ctx)
	if err != nil {
		return "", err
	}
	assertion := fmt.Sprintf("%s == %s", fieldName, fieldName)
	var expected, isBool bool
	expectedStr, isString := e.Value.(string)
	if isString {
		expected = expectedStr == "true"
	}
	_, isBool = e.Value.(bool)
	if isBool {
		expected = e.Value.(bool)
	}
	if !expected {
		assertion = "not " + assertion
	}
	return assertion, nil
}
