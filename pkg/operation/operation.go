package operation

import (
	"fmt"
	"json-rule-finder/pkg/condition"
	"json-rule-finder/pkg/shared"
	"json-rule-finder/pkg/value"
	"strings"
)

type baseOperation struct {
	conditionSetName string
}

func (o baseOperation) GetConditionSetName() string {
	return o.conditionSetName
}

type Operation interface {
	shared.Rego
	GetConditionSetName() string
}

func NewOperationOrCondition(input map[string]any, ctx *shared.Context) shared.Rego {
	if operation := tryParseOperation(input, ctx); operation != nil {
		return operation
	}
	subject := tryParseSubject(input, ctx)
	if cond := tryParseCondition(subject, input); cond != nil {
		return cond
	}
	return nil
}

func NewOperation(operationType string, body any, ctx *shared.Context) shared.Rego {
	switch operationType {
	case shared.AllOf:
		return ParseAllOf(body, ctx)
	case shared.AnyOf:
		return ParseAnyOf(body, ctx)
	case shared.Not:
		return ParseNot(body, ctx)
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

func parseOperationBody(input any, ctx *shared.Context) ([]shared.Rego, baseOperation, error) {
	items := input.([]any)
	var body []shared.Rego
	for _, item := range items {
		itemMap := item.(map[string]any)
		rego := NewOperationOrCondition(itemMap, ctx)
		body = append(body, rego)
	}
	conditionSetName, err := NeoConditionNameGenerator(ctx)
	if err != nil {
		return nil, baseOperation{}, err
	}
	return body, baseOperation{
		conditionSetName: conditionSetName,
	}, nil
}

func tryParseCondition(subject shared.Rego, input map[string]any) shared.Rego {
	for key, conditionValue := range input {
		key = strings.ToLower(key)
		if ifBody := condition.NewCondition(key, subject, conditionValue); ifBody != nil {
			return ifBody
		}
	}
	return nil
}

func tryParseOperation(conditionMap map[string]any, ctx *shared.Context) shared.Rego {
	for key, conditionValue := range conditionMap {
		key = strings.ToLower(key)
		operation := NewOperation(key, conditionValue, ctx)
		if operation != nil {
			return operation
		}
	}
	return nil
}

func tryParseSubject(conditionMap map[string]any, ctx *shared.Context) shared.Rego {
	for key, conditionValue := range conditionMap {
		key = strings.ToLower(key)
		switch key {
		case shared.Count:
			return NewCount(conditionValue, ctx)
		case shared.Field:
			if conditionValue == shared.TypeOfResource {
				if resourceType, ok := conditionMap["equals"].(string); ok {
					ctx.PushResourceType(resourceType)
				}
			}
			return value.NewFieldValue(conditionValue, ctx)
		case shared.Value:
			return value.NewLiteralValue(conditionValue, ctx)
		}
	}
	return nil
}
