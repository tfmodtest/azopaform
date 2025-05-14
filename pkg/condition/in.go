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
	fieldName, err := i.GetSubject(ctx).Rego(ctx)
	if err != nil {
		return "", err
	}
	prefix := ""
	if utilLibraryName := ctx.UtilLibraryPackageName(); utilLibraryName != "" {
		prefix = fmt.Sprintf("data.%s.", utilLibraryName)
	}
	return fmt.Sprintf("%sarraycontains(%s, %s)", prefix, shared.SliceConstructor(i.Values), fieldName), nil
}
