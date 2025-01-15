package pkg

import (
	"fmt"
	"json-rule-finder/pkg/condition"
	"json-rule-finder/pkg/shared"
	"reflect"
	"strings"
)

type Operator interface {
	shared.Rego
	GetConditionSetName() string
}

var operatorFactories = make(map[string]func(input any, ctx *shared.Context) shared.Rego)

func init() {
	operatorFactories[shared.Count] = func(input any, ctx *shared.Context) shared.Rego {
		items := input.(map[string]any)
		var whereBody shared.Rego
		if items[shared.Where] != nil {
			whereMap := items[shared.Where].(map[string]any)
			of := operatorFactories[shared.Where]
			whereBody = of(whereMap, ctx)
		}
		fieldName := items[shared.Field]
		if items[shared.Field] == nil {
			fieldName = items[shared.Value]
		}
		countField, _, err := shared.FieldNameProcessor(fieldName.(string), ctx)
		if err != nil {
			countField = items[shared.Field].(string)
			fmt.Printf("error in field name processor: %v\n", err)
		}
		fmt.Printf("count field is %v\n", countField)
		countFieldConverted := replaceIndex(countField)
		var countBody string
		if whereBody != nil {
			countBody = shared.Count + "(" + "{" + "x" + "|" + countFieldConverted + ";" + whereBody.(WhereOperator).ConditionSetName + "(x)" + "}" + ")"
		} else {
			countBody = shared.Count + "(" + "{" + "x" + "|" + countFieldConverted + "}" + ")"
		}
		countBody = strings.Replace(countBody, "*", "x", -1)
		return CountOperator{
			Where:    whereBody,
			CountExp: countBody,
		}
	}
	operatorFactories[shared.AllOf] = func(input any, ctx *shared.Context) shared.Rego {
		items := input.([]any)
		var body []shared.Rego
		for _, item := range items {
			itemMap := item.(map[string]any)
			var cf func(shared.Rego, any) shared.Rego
			var conditionKey string
			var subjectKey string
			var of func(any, *shared.Context) shared.Rego
			var operatorValue any
			var containsTypeOfResource bool
			for k, v := range itemMap {
				if k == shared.Field && v == shared.TypeOfResource {
					containsTypeOfResource = true
				}
				if f, ok := condition.ConditionFactory[strings.ToLower(k)]; ok {
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
					if k == shared.Field && v == shared.TypeOfResource {
						continue
					}
					if reflect.TypeOf(v).Kind() == reflect.String {
						ctx.PushResourceType(v.(string))
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
				if subjectKey == shared.Field && subjectItem == shared.TypeOfResource {
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
			fmt.Printf("error in BaseCondition name generator: %v\n", err)
			return nil
		}
		return AllOf{
			baseOperator: baseOperator{
				conditionSetName: conditionSetName,
			},
			Conditions: body,
		}
	}
	operatorFactories[shared.AnyOf] = func(input any, ctx *shared.Context) shared.Rego {
		items := input.([]any)
		var body []shared.Rego
		for _, item := range items {
			itemMap := item.(map[string]any)
			var cf func(shared.Rego, any) shared.Rego
			var conditionKey string
			var subjectKey string
			var of func(any, *shared.Context) shared.Rego
			var operatorValue any
			var containsTypeOfResource bool
			for k, v := range itemMap {
				if k == shared.Field && v == shared.TypeOfResource {
					containsTypeOfResource = true
				}
				if f, ok := condition.ConditionFactory[strings.ToLower(k)]; ok {
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
					if k == shared.Field && v == shared.TypeOfResource {
						continue
					}
					ctx.PushResourceType(v.(string))
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
				if subjectKey == shared.Field && subjectItem == shared.TypeOfResource {
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
	operatorFactories[shared.Not] = func(input any, ctx *shared.Context) shared.Rego {
		itemMap := input.(map[string]any)
		fmt.Printf("item map is %v\n", itemMap)
		var cf func(shared.Rego, any) shared.Rego
		var conditionKey string
		var subjectKey string
		var of func(any, *shared.Context) shared.Rego
		var subject shared.Rego
		var body shared.Rego
		var operatorValue any
		for k, _ := range itemMap {
			if f, ok := condition.ConditionFactory[strings.ToLower(k)]; ok {
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
		conditionName, err := NeoConditionNameGenerator(ctx)
		if err != nil {
			return nil
		}
		return NotOperator{
			Body:             body,
			ConditionSetName: conditionName,
		}
	}
	operatorFactories[shared.Where] = func(input any, ctx *shared.Context) shared.Rego {
		itemMap := input.(map[string]any)
		var body []shared.Rego
		var cf func(shared.Rego, any) shared.Rego
		var conditionKey string
		var subjectKey string
		var of func(any, *shared.Context) shared.Rego
		var operatorValue any
		for k, _ := range itemMap {
			if f, ok := condition.ConditionFactory[k]; ok {
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

var _ shared.Rego = &CountOperator{}

type CountOperator struct {
	Where    shared.Rego
	CountExp string
}

func (c CountOperator) Rego(ctx *shared.Context) (string, error) {
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

var _ shared.Rego = &WhereOperator{}

var _ Operator = &WhereOperator{}

type WhereOperator struct {
	Conditions       []shared.Rego
	ConditionSetName string
}

func (w WhereOperator) GetConditionSetName() string {
	return w.ConditionSetName
}

func (w WhereOperator) GetConditionSetNameRev() string {
	return strings.Join([]string{"not", w.ConditionSetName}, " ")
}

func (w WhereOperator) Rego(ctx *shared.Context) (string, error) {
	var res string
	var subSets []string

	for _, item := range w.Conditions {
		if _, ok := item.(Operator); ok {
			res += item.(Operator).GetConditionSetName() + "(x)"
			ctx.PushFieldName("x")
			subSet, err := item.Rego(ctx)
			if err != nil {
				return "", err
			}
			subSets = append(subSets, subSet)
		} else {
			ctx.PushFieldName("x")
			condition, err := item.Rego(ctx)
			if err != nil {
				return "", err
			}
			res += condition
		}
	}

	res = w.ConditionSetName + "(x)" + " " + shared.IfCondition + " {\n" + res
	res = res + "\n" + "}"

	for _, subSet := range subSets {
		res += "\n" + subSet
	}

	// add BaseCondition set body at the end
	return res, nil
}

var NeoConditionNameGenerator = func(ctx *shared.Context) (string, error) {
	index, ok := ctx.PopConditionNameCounter()
	if !ok {
		return "", fmt.Errorf("conditionNameStack is empty")
	}
	conditionName := "condition" + fmt.Sprintf("%d", index)
	return conditionName, nil
}
