package condition

import (
	"fmt"
	"github.com/tfmodtest/azopaform/pkg/shared"
)

var _ Condition = In{}

type In struct {
	BaseCondition
	Values []string
}

func (i In) Rego(ctx *shared.Context) (string, error) {
	return subjectRego(i.GetSubject(ctx), i.Values, func(subject shared.Rego, value any, ctx *shared.Context) (string, error) {
		prefix := ""
		if utilLibraryName := ctx.UtilLibraryPackageName(); utilLibraryName != "" {
			prefix = fmt.Sprintf("data.%s.", utilLibraryName)
		}
		values := value.([]string)
		if field, ok := subject.(FieldValue); ok && field.Name == "type" {
			return fmt.Sprintf("%sis_azure_type(%s, r.values)", prefix, shared.SliceConstructor(values)), nil
		}
		fieldName, err := subject.Rego(ctx)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%sarraycontains(%s, %s)", prefix, shared.SliceConstructor(values), fieldName), nil
	}, ctx)
}
