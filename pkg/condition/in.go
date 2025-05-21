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
	prefix := ""
	if utilLibraryName := ctx.UtilLibraryPackageName(); utilLibraryName != "" {
		prefix = fmt.Sprintf("data.%s.", utilLibraryName)
	}
	if field, ok := i.GetSubject(ctx).(FieldValue); ok && field.Name == "type" {
		return fmt.Sprintf("%sis_azure_type(%s, r.values)", prefix, shared.SliceConstructor(i.Values)), nil
	}
	fieldName, err := i.GetSubject(ctx).Rego(ctx)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%sarraycontains(%s, %s)", prefix, shared.SliceConstructor(i.Values), fieldName), nil
}
