package operation

import (
	"strings"

	"github.com/tfmodtest/azopaform/pkg/shared"
)

var _ Operation = &AllOf{}

type AllOf struct {
	baseOperation
	Conditions []shared.Rego
}

func NewAllOf(conditionSetName string, conditions []shared.Rego) AllOf {
	return AllOf{
		baseOperation: baseOperation{
			helperFunctionName: conditionSetName,
		},
		Conditions: conditions,
	}
}

func ParseAllOf(input any, ctx *shared.Context) (shared.Rego, error) {
	body, err := parseOperationBody(input, ctx)
	if err != nil {
		return nil, err
	}
	return AllOf{
		baseOperation: newBaseOperation(),
		Conditions:    body,
	}, nil
}

func (a AllOf) Rego(ctx *shared.Context) (string, error) {
	return a.wrapToFunction(func() (string, error) {
		sb := strings.Builder{}
		for _, item := range a.Conditions {
			if operation, ok := item.(Operation); ok {
				funcDecl, err := a.asFunctionForOperation(operation, ctx)
				if err != nil {
					return "", err
				}
				sb.WriteString("\n")
				sb.WriteString(funcDecl)
				continue
			}

			condition, err := item.Rego(ctx)
			if err != nil {
				return "", err
			}
			sb.WriteString("\n" + condition)
		}
		return sb.String(), nil
	}, ctx)
}
