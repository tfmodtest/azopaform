package pkg

import (
	"context"
	"fmt"
	"github.com/emirpasic/gods/stacks"
	"reflect"
	"strings"
)

var operatorFactories = make(map[string]func(input any, ctx context.Context) Rego)

func init() {
	operatorFactories[count] = func(input any, ctx context.Context) Rego {
		items := input.(map[string]any)
		var whereBody Rego
		whereMap := items[where].(map[string]any)
		of := operatorFactories[where]
		whereBody = of(whereMap, ctx)
		countField := items[field].(string)
		countBody := count + "(" + "{" + "x" + "|" + countField + ";" + whereBody.(WhereOperator).ConditionSetName + "(x)" + "}" + ")"
		countBody = strings.Replace(countBody, "*", "x", -1)
		return CountOperator{
			Where:    whereBody,
			CountExp: countBody,
		}
	}
	operatorFactories[allOf] = func(input any, ctx context.Context) Rego {
		items := input.([]any)
		var body []Rego
		for _, item := range items {
			itemMap := item.(map[string]any)
			var cf func(Rego, any) Rego
			var conditionKey string
			var subjectKey string
			var of func(any, context.Context) Rego
			var operatorValue any
			var containsTypeOfResource bool
			for k, v := range itemMap {
				if k == field && v == typeOfResource {
					containsTypeOfResource = true
				}
				if f, ok := conditionFactory[k]; ok {
					cf = f
					conditionKey = k
					continue
				}
			}
			for k, v := range itemMap {
				if f, ok := operatorFactories[k]; ok {
					of = f
					operatorValue = v
					continue
				}
			}
			if containsTypeOfResource {
				for k, v := range itemMap {
					if k == field && v == typeOfResource {
						continue
					}
					pushResourceType(ctx, v.(string))
				}
			}
			if cf != nil {
				for k, _ := range itemMap {
					if k == conditionKey {
						continue
					}
					subjectKey = k
					fmt.Printf("subject key is %v\n", subjectKey)
				}
				subject := subjectFactories[subjectKey](itemMap[subjectKey], ctx)
				if reflect.TypeOf(subject) == reflect.TypeOf(Count{}) {
					body = append(body, cf(subject, itemMap[conditionKey]))
					body = append(body, subject.(Count).ConditionSet)
				} else {
					body = append(body, cf(subject, itemMap[conditionKey]))
				}
			} else if of != nil {
				body = append(body, of(operatorValue, ctx))
			}
		}
		conditionSetName := conditionNameGenerator(andConditionLen, charNum)
		return AllOf{
			Conditions:       body,
			ConditionSetName: conditionSetName,
		}
	}
	operatorFactories[anyOf] = func(input any, ctx context.Context) Rego {
		items := input.([]any)
		var body []Rego
		for _, item := range items {
			itemMap := item.(map[string]any)
			var cf func(Rego, any) Rego
			var conditionKey string
			var subjectKey string
			var of func(any, context.Context) Rego
			var operatorValue any
			for k, _ := range itemMap {
				if f, ok := conditionFactory[k]; ok {
					cf = f
					conditionKey = k
					continue
				}
			}
			for k, v := range itemMap {
				if f, ok := operatorFactories[k]; ok {
					of = f
					operatorValue = v
					continue
				}
			}
			if cf != nil {
				for k, _ := range itemMap {
					if k == conditionKey {
						continue
					}
					subjectKey = k
				}
				subject := subjectFactories[subjectKey](itemMap[subjectKey], ctx)
				body = append(body, cf(subject, itemMap[conditionKey]))
			} else if of != nil {
				body = append(body, of(operatorValue, ctx))
			}
		}
		return AnyOf{
			Conditions:       body,
			ConditionSetName: conditionNameGenerator(orConditionLen, charNum),
		}
	}
	operatorFactories[not] = func(input any, ctx context.Context) Rego {
		itemMap := input.(map[string]any)
		var cf func(Rego, any) Rego
		var conditionKey string
		var subjectKey string
		for k, _ := range itemMap {
			if f, ok := conditionFactory[k]; ok {
				cf = f
				conditionKey = k
				continue
			}
		}
		for k, _ := range itemMap {
			if k == conditionKey {
				continue
			}
			subjectKey = k
		}
		subject := subjectFactories[subjectKey](itemMap[subjectKey], ctx)
		return NotOperator{
			Body:             cf(subject, itemMap[conditionKey]),
			ConditionSetName: conditionNameGenerator(singleConditionLen, charNum),
		}
	}
	operatorFactories[where] = func(input any, ctx context.Context) Rego {
		itemMap := input.(map[string]any)
		var body Rego
		var cf func(Rego, any) Rego
		var conditionKey string
		var subjectKey string
		var of func(any, context.Context) Rego
		var operatorValue any
		for k, _ := range itemMap {
			if f, ok := conditionFactory[k]; ok {
				cf = f
				conditionKey = k
				continue
			}
		}
		for k, v := range itemMap {
			k = strings.ToLower(k)
			if f, ok := operatorFactories[k]; ok {
				of = f
				operatorValue = v
				continue
			}
		}
		if cf != nil {
			for k, _ := range itemMap {
				if k == conditionKey {
					continue
				}
				subjectKey = k
			}
			subject := subjectFactories[subjectKey](itemMap[subjectKey], ctx)
			body = cf(subject, itemMap[conditionKey])
		} else if of != nil {
			body = of(operatorValue, ctx)
		}
		conditionSetName := conditionNameGenerator(whereConditionLen, charNum)
		return WhereOperator{
			Condition:        body,
			ConditionSetName: conditionSetName,
		}
	}
}

var _ Rego = &NotOperator{}

type NotOperator struct {
	Body             Rego
	ConditionSetName string
}

func (n NotOperator) Rego(ctx context.Context) (string, error) {
	var res string
	condition, err := n.Body.Rego(ctx)
	if err != nil {
		return "", err
	}
	res = n.ConditionSetName + " {\n" + condition + "\n}"
	return res, nil
}

var _ Rego = &CountOperator{}

type CountOperator struct {
	Where    Rego
	CountExp string
}

func (c CountOperator) Rego(ctx context.Context) (string, error) {
	var res string
	whereSubset, err := c.Where.Rego(ctx)
	if err != nil {
		return "", err
	}
	if err != nil {
		return "", err
	}
	res = c.CountExp + "\n" + whereSubset
	return res, nil
}

var _ Rego = &WhereOperator{}

type WhereOperator struct {
	Condition        Rego
	ConditionSetName string
}

func (w WhereOperator) Rego(ctx context.Context) (string, error) {
	var res string
	var subSets []string
	item := w.Condition

	if reflect.TypeOf(item) == reflect.TypeOf(AnyOf{}) {
		// (x) should be added to subset names, potentially use ctx to pass it?
		res += not + " " + item.(AnyOf).ConditionSetName + "(x)"
		fieldNameReplacerStack := ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"]
		fieldNameReplacerStack.Push("x")
		subSet, err := item.Rego(ctx)
		if err != nil {
			return "", err
		}
		subSets = append(subSets, subSet)
	} else if reflect.TypeOf(item) == reflect.TypeOf(NotOperator{}) {
		res += not + " " + item.(NotOperator).ConditionSetName + "(x)"
		fieldNameReplacerStack := ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"]
		fieldNameReplacerStack.Push("x")
		subSet, err := item.Rego(ctx)
		if err != nil {
			return "", err
		}
		subSets = append(subSets, subSet)
	} else if reflect.TypeOf(item) == reflect.TypeOf(AllOf{}) {
		res += item.(AllOf).ConditionSetName + "(x)"
		fieldNameReplacerStack := ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"]
		fieldNameReplacerStack.Push("x")
		subSet, err := item.Rego(ctx)
		if err != nil {
			return "", err
		}
		subSets = append(subSets, subSet)
	} else {
		condition, err := item.Rego(ctx)
		if err != nil {
			return "", err
		}
		res += condition
	}

	res = w.ConditionSetName + "(x)" + " " + ifCondition + " {\n" + res
	res = res + "\n" + "}"

	for _, subSet := range subSets {
		res += "\n" + subSet
	}

	// add condition set body at the end
	return res, nil
}

var _ Rego = &AllOf{}

type AllOf struct {
	Conditions       []Rego
	ConditionSetName string
}

func (a AllOf) Rego(ctx context.Context) (string, error) {
	var res string
	var subSets []string

	if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
		res = a.ConditionSetName + "(x)" + " " + ifCondition + " {"
	} else {
		res = a.ConditionSetName + " " + ifCondition + " {"
	}

	for _, item := range a.Conditions {
		if reflect.TypeOf(item) == reflect.TypeOf(AnyOf{}) {
			if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
				res += not + " " + item.(AnyOf).ConditionSetName + "(x)"
			} else {
				res += not + " " + item.(AnyOf).ConditionSetName
			}
			subSet, err := item.Rego(ctx)
			if err != nil {
				return "", err
			}
			subSets = append(subSets, subSet)
			continue
		}
		if reflect.TypeOf(item) == reflect.TypeOf(NotOperator{}) {
			if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
				res += not + " " + item.(NotOperator).ConditionSetName + "(x)"
			} else {
				res += not + " " + item.(NotOperator).ConditionSetName
			}
			subSet, err := item.Rego(ctx)
			if err != nil {
				return "", err
			}
			subSets = append(subSets, subSet)
			continue
		}
		if reflect.TypeOf(item) == reflect.TypeOf(AllOf{}) {
			if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
				res += item.(AllOf).ConditionSetName + "(x)"
			} else {
				res += item.(AllOf).ConditionSetName
			}
			subSet, err := item.Rego(ctx)
			if err != nil {
				return "", err
			}
			subSets = append(subSets, subSet)
			continue
		}
		if reflect.TypeOf(item) == reflect.TypeOf(WhereOperator{}) {
			subSet, err := item.Rego(ctx)
			if err != nil {
				return "", err
			}
			subSets = append(subSets, subSet)
			continue
		}
		condition, err := item.Rego(ctx)
		if err != nil {
			return "", err
		}
		res += "\n" + condition
	}

	res = res + "\n}"

	for _, subSet := range subSets {
		res += "\n" + subSet
	}
	fmt.Printf("current res: %v\n", res)

	// add condition set body at the end
	return res, nil
}

var _ Rego = &AnyOf{}

type AnyOf struct {
	Conditions       []Rego
	ConditionSetName string
}

func (a AnyOf) Rego(ctx context.Context) (string, error) {
	var res string
	var subSets []string
	for _, item := range a.Conditions {
		if res != "" {
			res = res + "\n"
		}
		switch reflect.TypeOf(item) {
		case reflect.TypeOf(AllOf{}):
			if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
				res += not + " " + item.(AllOf).ConditionSetName + "(x)"
			} else {
				res += not + " " + item.(AllOf).ConditionSetName
			}
			subSet, err := item.Rego(ctx)
			if err != nil {
				return "", err
			}
			subSets = append(subSets, subSet)
		case reflect.TypeOf(AnyOf{}):
			if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
				res += item.(AnyOf).ConditionSetName + "(x)"
			} else {
				res += item.(AnyOf).ConditionSetName
			}
			subSet, err := item.Rego(ctx)
			if err != nil {
				return "", err
			}
			subSets = append(subSets, subSet)
		case reflect.TypeOf(NotOperator{}):
			if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
				res += item.(NotOperator).ConditionSetName + "(x)"
			} else {
				res += item.(NotOperator).ConditionSetName
			}
			subSet, err := item.Rego(ctx)
			if err != nil {
				return "", err
			}
			subSets = append(subSets, subSet)
		case reflect.TypeOf(EqualsOperation{}):
			oppoItem := NotEqualsOperation{
				operation: item.(EqualsOperation).operation,
				Value:     item.(EqualsOperation).Value,
			}
			oppoCondition, err := oppoItem.Rego(ctx)
			if err != nil {
				return "", err
			}
			res += oppoCondition
		case reflect.TypeOf(NotEqualsOperation{}):
			oppoItem := EqualsOperation{
				operation: item.(NotEqualsOperation).operation,
				Value:     item.(NotEqualsOperation).Value,
			}
			oppoCondition, err := oppoItem.Rego(ctx)
			if err != nil {
				return "", err
			}
			res += oppoCondition
		case reflect.TypeOf(LikeOperation{}):
			oppoItem := NotLikeOperation{
				operation: item.(LikeOperation).operation,
				Value:     item.(LikeOperation).Value,
			}
			oppoCondition, err := oppoItem.Rego(ctx)
			if err != nil {
				return "", err
			}
			res += oppoCondition
		case reflect.TypeOf(NotLikeOperation{}):
			oppoItem := LikeOperation{
				operation: item.(NotLikeOperation).operation,
				Value:     item.(NotLikeOperation).Value,
			}
			oppoCondition, err := oppoItem.Rego(ctx)
			if err != nil {
				return "", err
			}
			res += oppoCondition
		case reflect.TypeOf(ContainsOperation{}):
			oppoItem := NotContainsOperation{
				operation: item.(ContainsOperation).operation,
				Value:     item.(ContainsOperation).Value,
			}
			oppoCondition, err := oppoItem.Rego(ctx)
			if err != nil {
				return "", err
			}
			res += oppoCondition
		case reflect.TypeOf(NotContainsOperation{}):
			oppoItem := ContainsOperation{
				operation: item.(NotContainsOperation).operation,
				Value:     item.(NotContainsOperation).Value,
			}
			oppoCondition, err := oppoItem.Rego(ctx)
			if err != nil {
				return "", err
			}
			res += oppoCondition
		case reflect.TypeOf(InOperation{}):
			oppoItem := NotInOperation{
				operation: item.(InOperation).operation,
				Values:    item.(InOperation).Values,
			}
			oppoCondition, err := oppoItem.Rego(ctx)
			if err != nil {
				return "", err
			}
			res += oppoCondition
		case reflect.TypeOf(NotInOperation{}):
			oppoItem := InOperation{
				operation: item.(NotInOperation).operation,
				Values:    item.(NotInOperation).Values,
			}
			oppoCondition, err := oppoItem.Rego(ctx)
			if err != nil {
				return "", err
			}
			res += oppoCondition
		case reflect.TypeOf(LessOrEqualsOperation{}):
			oppoItem := GreaterOperation{
				operation: item.(LessOrEqualsOperation).operation,
				Value:     item.(LessOrEqualsOperation).Value,
			}
			oppoCondition, err := oppoItem.Rego(ctx)
			if err != nil {
				return "", err
			}
			res += oppoCondition
		case reflect.TypeOf(GreaterOperation{}):
			oppoItem := LessOrEqualsOperation{
				operation: item.(GreaterOperation).operation,
				Value:     item.(GreaterOperation).Value,
			}
			oppoCondition, err := oppoItem.Rego(ctx)
			if err != nil {
				return "", err
			}
			res += oppoCondition
		case reflect.TypeOf(LessOperation{}):
			oppoItem := GreaterOrEqualsOperation{
				operation: item.(LessOperation).operation,
				Value:     item.(LessOperation).Value,
			}
			oppoCondition, err := oppoItem.Rego(ctx)
			if err != nil {
				return "", err
			}
			res += oppoCondition
		case reflect.TypeOf(GreaterOrEqualsOperation{}):
			oppoItem := LessOperation{
				operation: item.(GreaterOrEqualsOperation).operation,
				Value:     item.(GreaterOrEqualsOperation).Value,
			}
			oppoCondition, err := oppoItem.Rego(ctx)
			if err != nil {
				return "", err
			}
			res += oppoCondition
		}
	}
	if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
		res = a.ConditionSetName + "(x)" + " " + ifCondition + " {\n" + res
	} else {
		res = a.ConditionSetName + " " + ifCondition + " {\n" + res
	}
	res = res + "\n" + "}"

	for _, subSet := range subSets {
		res += "\n" + subSet
	}
	return res, nil
}

func conditionNameGenerator(strLen int, charSet string) string {
	result := make([]byte, strLen)
	for i := 0; i < strLen; i++ {
		result[i] = charSet[RandIntRange(0, len(charSet))]
	}
	return string(result)
}
