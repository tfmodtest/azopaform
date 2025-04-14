package operation

import (
	"json-rule-finder/pkg/shared"
	"strings"
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

func ParseAllOf(input any, ctx *shared.Context) shared.Rego {
	body, err := parseOperationBody(input, ctx)
	if err != nil {
		panic(err)
	}
	return AllOf{
		baseOperation: newBaseOperation(),
		Conditions:    body,
	}
}

func (a AllOf) Rego(ctx *shared.Context) (string, error) {
	return a.WithFunction(func() (string, error) {
		sb := strings.Builder{}
		for _, item := range a.Conditions {
			if _, ok := item.(Operation); ok {
				if _, ok := ctx.FieldNameReplacer(); ok {
					sb.WriteString("\n" + item.(Operation).HelperFunctionName() + "(x)")
				} else {
					sb.WriteString("\n" + item.(Operation).HelperFunctionName())
				}
				subFunction, err := item.Rego(ctx)
				if err != nil {
					return "", err
				}
				ctx.EnqueueHelperFunction(subFunction)
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
