package condition

import (
	"fmt"
	"json-rule-finder/pkg/shared"
)

var _ Condition = Exists{}

type Exists struct {
	BaseCondition
	Value any
}

func (e Exists) Rego(ctx *shared.Context) (string, error) {
	fieldName, err := e.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	if _, ok := ctx.FieldNameReplacer(); ok {
		fieldName = ReplaceIndex(fieldName)
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
