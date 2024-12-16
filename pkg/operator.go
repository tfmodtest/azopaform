package pkg

import (
	"context"
	"fmt"
	"github.com/emirpasic/gods/stacks"
	"reflect"
	"strings"
)

type Operator interface {
	GetConditionSetName() string
	GetConditionSetNameRev() string
}

var operatorFactories = make(map[string]func(input any, ctx context.Context) Rego)

func init() {
	operatorFactories[count] = func(input any, ctx context.Context) Rego {
		items := input.(map[string]any)
		var whereBody Rego
		if items[where] != nil {
			whereMap := items[where].(map[string]any)
			of := operatorFactories[where]
			whereBody = of(whereMap, ctx)
		}
		fieldName := items[field]
		if items[field] == nil {
			fieldName = items[value]
		}
		countField, _, err := FieldNameProcessor(fieldName.(string), ctx)
		if err != nil {
			countField = items[field].(string)
			fmt.Printf("error in field name processor: %v\n", err)
		}
		fmt.Printf("count field is %v\n", countField)
		countFieldConverted := replaceIndex(countField)
		var countBody string
		if whereBody != nil {
			countBody = count + "(" + "{" + "x" + "|" + countFieldConverted + ";" + whereBody.(WhereOperator).ConditionSetName + "(x)" + "}" + ")"
		} else {
			countBody = count + "(" + "{" + "x" + "|" + countFieldConverted + "}" + ")"
		}
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
				if f, ok := conditionFactory[strings.ToLower(k)]; ok {
					cf = f
					conditionKey = k
					continue
				}
			}
			for k, v := range itemMap {
				if f, ok := operatorFactories[strings.ToLower(k)]; ok {
					fmt.Printf("operator is %v\n", k)
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
					if reflect.TypeOf(v).Kind() == reflect.String {
						pushResourceType(ctx, v.(string))
					}
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
				subjectItem := itemMap[subjectKey]
				if subjectKey == field && subjectItem == typeOfResource {
					if reflect.TypeOf(itemMap[conditionKey]).Kind() == reflect.String {
						rawType := itemMap[conditionKey]
						translatedType, err := ResourceTypeParser(rawType.(string))
						if err != nil {
							fmt.Printf("error in resource type parser: %v\n", err)
							return nil
						}
						body = append(body, cf(subjectFactories[subjectKey](subjectItem, ctx), translatedType))
						continue
					} else {
						rawTypes := itemMap[conditionKey]
						var translatedTypes []string
						for _, rawType := range rawTypes.([]interface{}) {
							translatedType, err := ResourceTypeParser(rawType.(string))
							if err != nil {
								fmt.Printf("error in resource type parser: %v\n", err)
								return nil
							}
							translatedTypes = append(translatedTypes, translatedType)
						}
						body = append(body, cf(subjectFactories[subjectKey](subjectItem, ctx), translatedTypes))
					}
				}
				subject := subjectFactories[subjectKey](subjectItem, ctx)
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
		//conditionSetName := conditionNameGenerator(andConditionLen, charNum)
		conditionSetName, err := NeoConditionNameGenerator(ctx)
		if err != nil {
			fmt.Printf("error in condition name generator: %v\n", err)
			return nil
		}
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
			var containsTypeOfResource bool
			for k, v := range itemMap {
				if k == field && v == typeOfResource {
					containsTypeOfResource = true
				}
				if f, ok := conditionFactory[strings.ToLower(k)]; ok {
					cf = f
					conditionKey = k
					continue
				}
			}
			for k, v := range itemMap {
				if f, ok := operatorFactories[strings.ToLower(k)]; ok {
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
				}
				subjectItem := itemMap[subjectKey]
				if subjectKey == field && subjectItem == typeOfResource {
					rawType := itemMap[conditionKey]
					translatedType, err := ResourceTypeParser(rawType.(string))
					if err != nil {
						return nil
					}
					body = append(body, cf(subjectFactories[subjectKey](subjectItem, ctx), translatedType))
					continue
				}
				subject := subjectFactories[subjectKey](subjectItem, ctx)
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
		//conditionName := conditionNameGenerator(orConditionLen, charNum)
		conditionName, err := NeoConditionNameGenerator(ctx)
		if err != nil {
			return nil
		}
		return AnyOf{
			Conditions:       body,
			ConditionSetName: conditionName,
		}
	}
	operatorFactories[not] = func(input any, ctx context.Context) Rego {
		itemMap := input.(map[string]any)
		fmt.Printf("item map is %v\n", itemMap)
		var cf func(Rego, any) Rego
		var conditionKey string
		var subjectKey string
		var of func(any, context.Context) Rego
		var subject Rego
		var body Rego
		var operatorValue any
		for k, _ := range itemMap {
			if f, ok := conditionFactory[strings.ToLower(k)]; ok {
				cf = f
				conditionKey = k
				continue
			}
		}
		for k, v := range itemMap {
			if f, ok := operatorFactories[strings.ToLower(k)]; ok {
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
			fmt.Printf("subject key is %v\n", subjectKey)
			subject = subjectFactories[subjectKey](itemMap[subjectKey], ctx)
			body = cf(subject, itemMap[conditionKey])
		} else if of != nil {
			body = of(operatorValue, ctx)
		}
		fmt.Printf("subject is %v\n", subject)
		//conditionName := conditionNameGenerator(singleConditionLen, charNum)
		conditionName, err := NeoConditionNameGenerator(ctx)
		if err != nil {
			return nil
		}
		return NotOperator{
			Body:             body,
			ConditionSetName: conditionName,
		}
	}
	operatorFactories[where] = func(input any, ctx context.Context) Rego {
		itemMap := input.(map[string]any)
		var body []Rego
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
			body = append(body, cf(subject, itemMap[conditionKey]))
		} else if of != nil {
			body = append(body, of(operatorValue, ctx))
			//fmt.Printf("where body is %v\n", body)
		}
		//conditionSetName := conditionNameGenerator(whereConditionLen, charNum)
		conditionSetName, err := NeoConditionNameGenerator(ctx)
		if err != nil {
			return nil
		}
		return WhereOperator{
			Conditions:       body,
			ConditionSetName: conditionSetName,
		}
	}
}

var _ Rego = &NotOperator{}

var _ Operator = &NotOperator{}

type NotOperator struct {
	Body             Rego
	ConditionSetName string
}

func (n NotOperator) GetConditionSetName() string {
	return strings.Join([]string{"not", n.ConditionSetName}, " ")
}

func (n NotOperator) GetConditionSetNameRev() string {
	return n.ConditionSetName
}

func (n NotOperator) Rego(ctx context.Context) (string, error) {
	var res string
	var subSets []string

	res = n.ConditionSetName + " " + ifCondition + " {\n"
	if _, ok := n.Body.(Operator); ok {
		if reflect.TypeOf(n.Body) != reflect.TypeOf(WhereOperator{}) {
			if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
				res += n.Body.(Operator).GetConditionSetName() + "(x)"
			} else {
				res += n.Body.(Operator).GetConditionSetName()
			}
		}
		subSet, err := n.Body.Rego(ctx)
		if err != nil {
			return "", err
		}
		subSets = append(subSets, subSet)
	} else {
		condition, err := n.Body.Rego(ctx)
		if err != nil {
			return "", err
		}
		res += condition
	}

	res += "\n}"
	for _, subSet := range subSets {
		res += "\n" + subSet
	}
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

var _ Operator = &WhereOperator{}

type WhereOperator struct {
	Conditions       []Rego
	ConditionSetName string
}

func (w WhereOperator) GetConditionSetName() string {
	return w.ConditionSetName
}

func (w WhereOperator) GetConditionSetNameRev() string {
	return strings.Join([]string{"not", w.ConditionSetName}, " ")
}

func (w WhereOperator) Rego(ctx context.Context) (string, error) {
	var res string
	var subSets []string

	for _, item := range w.Conditions {
		if _, ok := item.(Operator); ok {
			res += item.(Operator).GetConditionSetName() + "(x)"
			fieldNameReplacerStack := ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"]
			fieldNameReplacerStack.Push("x")
			subSet, err := item.Rego(ctx)
			if err != nil {
				return "", err
			}
			subSets = append(subSets, subSet)
		} else {
			fieldNameReplacerStack := ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"]
			fieldNameReplacerStack.Push("x")
			condition, err := item.Rego(ctx)
			if err != nil {
				return "", err
			}
			res += condition
		}
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

var _ Operator = &AllOf{}

type AllOf struct {
	Conditions       []Rego
	ConditionSetName string
}

func (a AllOf) GetConditionSetName() string {
	return a.ConditionSetName
}

func (a AllOf) GetConditionSetNameRev() string {
	return strings.Join([]string{"not", a.ConditionSetName}, " ")
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
		fmt.Printf("the item has type %v\n", reflect.TypeOf(item))
		if reflect.TypeOf(item) == nil {
			continue
		}
		if _, ok := item.(Operator); ok {
			if reflect.TypeOf(item) != reflect.TypeOf(WhereOperator{}) {
				if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
					res += "\n" + item.(Operator).GetConditionSetName() + "(x)"
				} else {
					res += "\n" + item.(Operator).GetConditionSetName()
				}
			}
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

	// add condition set body at the end
	return res, nil
}

var _ Rego = &AnyOf{}

type AnyOf struct {
	Conditions       []Rego
	ConditionSetName string
}

func (a AnyOf) GetConditionSetName() string {
	return strings.Join([]string{"not", a.ConditionSetName}, " ")
}

func (a AnyOf) GetConditionSetNameRev() string {
	return a.ConditionSetName
}

func (a AnyOf) Rego(ctx context.Context) (string, error) {
	var res string
	var subSets []string
	for _, item := range a.Conditions {
		if res != "" {
			res = res + "\n"
		}
		if _, ok := item.(Operator); ok {
			if reflect.TypeOf(item) != reflect.TypeOf(WhereOperator{}) {
				if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
					res += item.(Operator).GetConditionSetNameRev() + "(x)"
				} else {
					res += item.(Operator).GetConditionSetNameRev()
				}
			}
			subSet, err := item.Rego(ctx)
			if err != nil {
				return "", err
			}
			subSets = append(subSets, subSet)
			continue
		}

		if _, ok := item.(Condition); ok {
			condition, err := item.(Condition).GetReverseRego(ctx)
			if err != nil {
				return "", err
			}
			res += condition
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

var NeoConditionNameGenerator = func(ctx context.Context) (string, error) {
	conditionNameStack := ctx.Value("context").(map[string]stacks.Stack)["conditionNameCounter"]
	if conditionNameStack == nil {
		return "", fmt.Errorf("conditionNameStack is nil")
	}
	index, ok := conditionNameStack.Pop()
	if !ok {
		return "", fmt.Errorf("conditionNameStack is empty")
	}
	conditionName := "condition" + fmt.Sprintf("%v", index)
	return conditionName, nil
}
