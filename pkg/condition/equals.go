package condition

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/tfmodtest/azopaform/pkg/shared"
)

var _ Condition = Equals{}

type Equals struct {
	BaseCondition
	Value any
}

// Rego For conditions under 'where' operator, "[[0-9]+]" should be replaced with "[x]"
func (e Equals) Rego(ctx *shared.Context) (string, error) {
	return subjectRego(e.GetSubject(ctx), e.Value, func(subject shared.Rego, value any, ctx *shared.Context) (string, error) {
		fieldName, err := subject.Rego(ctx)
		if err != nil {
			return "", err
		}
		var v string
		//TODO:refactor this
		if reflect.TypeOf(value).Kind() == reflect.String {
			v = strings.Join([]string{"\"", fmt.Sprint(value), "\""}, "")
		} else if reflect.TypeOf(value).Kind() == reflect.Bool {
			v = fmt.Sprint(value)
		} else {
			v = fmt.Sprint(value)
		}
		equals := strings.Join([]string{fieldName, "==", v}, " ")
		prefix := ""
		if utilLibraryName := ctx.UtilLibraryPackageName(); utilLibraryName != "" {
			prefix = fmt.Sprintf("data.%s.", utilLibraryName)
		}
		if field, ok := subject.(FieldValue); ok && field.Name == shared.TypeOfResource {
			equals = fmt.Sprintf(`%sis_azure_type(%s, %s)`, prefix, shared.ResourcePathPrefix, v)
		}
		return equals, nil
	}, ctx)
}
