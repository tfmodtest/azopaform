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
	return subjectRego(e.GetSubject(ctx), e.Value, func(subject shared.Rego, value any, ctx *shared.Context) (string, error) {
		fieldName, err := subject.Rego(ctx)
		if err != nil {
			return "", err
		}
		assertion := fmt.Sprintf("%s == %s", fieldName, fieldName)
		var expected, isBool bool
		expectedStr, isString := value.(string)
		if isString {
			expected = expectedStr == "true"
		}
		_, isBool = value.(bool)
		if isBool {
			expected = value.(bool)
		}
		if !expected {
			assertion = "not " + assertion
		}
		return assertion, nil
	}, ctx)
}
