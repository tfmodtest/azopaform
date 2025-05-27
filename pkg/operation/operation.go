package operation

import (
	"fmt"
	"strings"

	"github.com/emirpasic/gods/sets/hashset"
	"github.com/tfmodtest/azopaform/pkg/condition"
	"github.com/tfmodtest/azopaform/pkg/shared"
)

var usedRandomString = hashset.New()

type baseOperation struct {
	helperFunctionName string
}

func newBaseOperation() baseOperation {
	return baseOperation{
		helperFunctionName: RandomHelperFunctionNameGenerator(),
	}
}

func (o baseOperation) wrapToFunction(body func() (string, error), ctx *shared.Context) (string, error) {
	bodyContent, err := body()
	if err != nil {
		return "", err
	}

	res := o.HelperFunctionName() + "(r) " + shared.IfCondition + " {"
	//if _, ok := ctx.VarNameForField(); ok {
	//	res = o.HelperFunctionName() + "(x)" + " " + shared.IfCondition + " {"
	//}
	sb := strings.Builder{}
	sb.WriteString(res)
	sb.WriteString("\n")
	sb.WriteString(bodyContent)
	sb.WriteString("\n}")
	return sb.String(), nil
}

func (o baseOperation) asFunctionForOperation(operation Operation, ctx *shared.Context) (string, error) {
	sb := strings.Builder{}
	call := operation.HelperFunctionName()
	//if _, ok := ctx.VarNameForField(); ok {
	//	call += "(x)"
	//}
	operationDecl, err := operation.Rego(ctx)
	if err != nil {
		return "", err
	}
	sb.WriteString(call)
	sb.WriteString("(r)")
	sb.WriteString("\n")
	ctx.EnqueueHelperFunction(operationDecl)
	return sb.String(), nil
}

func (o baseOperation) forkFunctionForOperation(operation Operation, ctx *shared.Context) (string, error) {
	sb := strings.Builder{}
	subSet, err := operation.Rego(ctx)
	if err != nil {
		return "", err
	}
	funcDef, _ := o.wrapToFunction(func() (string, error) {
		if _, ok := ctx.VarNameForField(); ok {
			return operation.HelperFunctionName() + "(r, x)", nil
		}
		return operation.HelperFunctionName() + "(r)", nil
	}, ctx)
	sb.WriteString(funcDef)
	sb.WriteString("\n")
	ctx.EnqueueHelperFunction(subSet)
	return sb.String(), nil
}

func (o baseOperation) forkFunctionForCondition(cond shared.Rego, ctx *shared.Context) (string, error) {
	sb := strings.Builder{}
	condStr, err := cond.Rego(ctx)
	if err != nil {
		return "", err
	}
	funcDef, _ := o.wrapToFunction(func() (string, error) {
		return condStr, nil
	}, ctx)
	sb.WriteString(funcDef)
	sb.WriteString("\n")
	return sb.String(), nil
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
	cond, err := tryParseCondition(subject, input, ctx)
	if err != nil {
		return nil, err
	}
	if cond != nil {
		return cond, nil
	}
	return nil, fmt.Errorf("unknown operation or condition: %v", input)
}

func NewOperation(operationType string, body any, ctx *shared.Context) (shared.Rego, error) {
	switch operationType {
	case shared.AllOf:
		return ParseAllOf(body, ctx)
	case shared.AnyOf:
		return ParseAnyOf(body, ctx)
	case shared.Not:
		return parseNot(body, ctx)
	}
	return nil, nil
}

var RandomHelperFunctionNameGenerator = func() string {
	var randomSuffix string
	for {
		randomSuffix = shared.HumanFriendlyEnglishString(10)
		if !usedRandomString.Contains(randomSuffix) {
			usedRandomString.Add(randomSuffix)
			break
		}
	}
	return fmt.Sprintf("condition_%s", randomSuffix)
}

func parseOperationBody(input any, ctx *shared.Context) ([]shared.Rego, error) {
	items := input.([]any)
	var bodies []shared.Rego
	for _, item := range items {
		itemMap := item.(map[string]any)
		body, err := NewOperationOrCondition(itemMap, ctx)
		if err != nil {
			return nil, err
		}
		bodies = append(bodies, body)
	}
	return bodies, nil
}

func tryParseCondition(subject shared.Rego, input map[string]any, ctx *shared.Context) (shared.Rego, error) {
	for key, conditionValue := range input {
		key = strings.ToLower(key)
		cond, err := condition.NewCondition(key, subject, conditionValue, ctx)
		if err != nil {
			return nil, err
		}
		if cond != nil {
			return cond, nil
		}
	}
	return nil, nil
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
			return condition.NewFieldValue(conditionValue, ctx)
		case shared.Value:
			return condition.NewLiteralValue(conditionValue, ctx)
		}
	}
	return nil, nil
}
