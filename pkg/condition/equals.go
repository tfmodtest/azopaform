package condition

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/tfmodtest/azopaform/pkg/shared"
	"github.com/tfmodtest/azopaform/pkg/value"
)

var _ Condition = Equals{}

type Equals struct {
	BaseCondition
	Value any
}

// Rego For conditions under 'where' operator, "[[0-9]+]" should be replaced with "[x]"
func (e Equals) Rego(ctx *shared.Context) (string, error) {
	fieldName, err := e.GetSubject(ctx).Rego(ctx)
	if err != nil {
		return "", err
	}
	var v string
	if reflect.TypeOf(e.Value).Kind() == reflect.String {
		v = strings.Join([]string{"\"", fmt.Sprint(e.Value), "\""}, "")
	} else if reflect.TypeOf(e.Value).Kind() == reflect.Bool {
		v = fmt.Sprint(e.Value)
	} else {
		v = fmt.Sprint(e.Value)
	}
	equals := strings.Join([]string{fieldName, "==", v}, " ")
	prefix := ""
	if utilLibraryName := ctx.UtilLibraryPackageName(); utilLibraryName != "" {
		prefix = fmt.Sprintf("data.%s.", utilLibraryName)
	}
	if field, ok := e.Subject.(value.FieldValue); ok && field.Name == shared.TypeOfResource {
		equals = fmt.Sprintf(`%sis_azure_type(%s, %s)`, prefix, shared.ResourcePathPrefix, v)
	}
	return equals, nil
}
