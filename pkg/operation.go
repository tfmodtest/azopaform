package pkg

import (
	"fmt"
	"json-rule-finder/pkg/shared"
)

type Operation interface {
	shared.Rego
	GetConditionSetName() string
}

func NewOperationOrCondition(input map[string]any, ctx *shared.Context) shared.Rego {
	if operation := extractOperation(input, ctx); operation != nil {
		return operation
	}
	subject := extractSubject(input, ctx)
	if cond := extractCondition(subject, input); cond != nil {
		return cond
	}
	return nil
}

func NewOperation(operationType string, body any, ctx *shared.Context) shared.Rego {
	switch operationType {
	case shared.AllOf:
		return NewAllOf(body, ctx)
	case shared.AnyOf:
		return NewAnyOf(body, ctx)
	case shared.Not:
		return NewNot(body, ctx)
	}
	return nil
}

var NeoConditionNameGenerator = func(ctx *shared.Context) (string, error) {
	index, ok := ctx.PopConditionNameCounter()
	if !ok {
		return "", fmt.Errorf("conditionNameStack is empty")
	}
	conditionName := "condition" + fmt.Sprintf("%d", index)
	return conditionName, nil
}
