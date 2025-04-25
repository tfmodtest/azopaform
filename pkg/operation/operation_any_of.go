package operation

import (
	"strings"

	"github.com/tfmodtest/azopaform/pkg/shared"
)

var _ Operation = &AnyOf{}

type AnyOf struct {
	baseOperation
	Conditions []shared.Rego
}

func NewAnyOf(conditionSetName string, conditions []shared.Rego) AnyOf {
	return AnyOf{
		baseOperation: baseOperation{
			helperFunctionName: conditionSetName,
		},
		Conditions: conditions,
	}
}

func ParseAnyOf(input any, ctx *shared.Context) shared.Rego {
	body, err := parseOperationBody(input, ctx)
	if err != nil {
		panic(err)
	}
	return AnyOf{
		Conditions:    body,
		baseOperation: newBaseOperation(),
	}
}

func (a AnyOf) Rego(ctx *shared.Context) (string, error) {
	sb := strings.Builder{}
	for _, item := range a.Conditions {
		var funcDef string
		var err error
		if operation, ok := item.(Operation); ok {
			funcDef, err = a.forkFunctionForOperation(operation, ctx)
		} else {
			funcDef, err = a.forkFunctionForCondition(item, ctx)
		}
		if err != nil {
			return "", err
		}
		sb.WriteString(funcDef)
	}
	return sb.String(), nil
}
