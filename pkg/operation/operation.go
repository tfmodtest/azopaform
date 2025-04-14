package operation

import (
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"json-rule-finder/pkg/condition"
	"json-rule-finder/pkg/shared"
	"json-rule-finder/pkg/value"
	"strings"

	"github.com/xyproto/randomstring"
)

var usedRandomString = hashset.New()

type baseOperation struct {
	helperFunctionName string
}

func (o baseOperation) HelperFunctionName() string {
	return o.helperFunctionName
}

type Operation interface {
	shared.Rego
	HelperFunctionName() string
}

func NewOperationOrCondition(input map[string]any, ctx *shared.Context) (shared.Rego, error) {
	var operation shared.Rego
	var err error
	operation, err = tryParseOperation(input, ctx)
	if err != nil {
		return nil, err
	}
	if operation != nil {
		return operation, nil
	}
	subject, err := tryParseSubject(input, ctx)
	if err != nil {
		return nil, err
	}
	if cond := tryParseCondition(subject, input); cond != nil {
		return cond, nil
	}
	return nil, fmt.Errorf("unknown operation or condition: %v", input)
}

func NewOperation(operationType string, body any, ctx *shared.Context) (shared.Rego, error) {
	switch operationType {
	case shared.AllOf:
		return ParseAllOf(body, ctx), nil
	case shared.AnyOf:
		return ParseAnyOf(body, ctx), nil
	case shared.Not:
		return ParseNot(body, ctx)
	}
	return nil, nil
}

var NeoConditionNameGenerator = func(ctx *shared.Context) string {
	var randomSuffix string
	for {
		randomSuffix = randomstring.HumanFriendlyEnglishString(10)
		if !usedRandomString.Contains(randomSuffix) {
			usedRandomString.Add(randomSuffix)
			break
		}
	}
	return fmt.Sprintf("condition%s", randomSuffix)
}

func parseOperationBody(input any, ctx *shared.Context) ([]shared.Rego, string, error) {
	items := input.([]any)
	var bodies []shared.Rego
	for _, item := range items {
		itemMap := item.(map[string]any)
		body, err := NewOperationOrCondition(itemMap, ctx)
		if err != nil {
			return nil, "", err
		}
		bodies = append(bodies, body)
	}
	conditionSetName := NeoConditionNameGenerator(ctx)
	return bodies, conditionSetName, nil
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

func tryParseOperation(conditionMap map[string]any, ctx *shared.Context) (shared.Rego, error) {
	for key, conditionValue := range conditionMap {
		key = strings.ToLower(key)
		operation, err := NewOperation(key, conditionValue, ctx)
		if err != nil {
			return nil, err
		}
		if operation != nil {
			return operation, nil
		}
	}
	return nil, nil
}

func tryParseSubject(conditionMap map[string]any, ctx *shared.Context) (shared.Rego, error) {
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
			return value.NewFieldValue(conditionValue, ctx), nil
		case shared.Value:
			return value.NewLiteralValue(conditionValue, ctx), nil
		}
	}
	return nil, nil
}
