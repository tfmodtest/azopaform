package pkg

import (
	"context"
	"reflect"
)

var operatorFactories = make(map[string]func(input any) Rego)

func init() {
	operatorFactories[allOf] = func(input any) Rego {
		items := input.([]any)
		var body []Rego
		for _, item := range items {
			itemMap := item.(map[string]any)
			var cf func(Rego, any) Rego
			var conditionKey string
			var subjectKey string
			var of func(any) Rego
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
				subject := subjectFactories[subjectKey](itemMap[subjectKey])
				body = append(body, cf(subject, itemMap[conditionKey]))
			} else if of != nil {
				body = append(body, of(operatorValue))
			}
		}
		return AllOf(body)
	}
	operatorFactories[anyOf] = func(input any) Rego {
		items := input.([]any)
		var body []Rego
		for _, item := range items {
			itemMap := item.(map[string]any)
			var cf func(Rego, any) Rego
			var conditionKey string
			var subjectKey string
			var of func(any) Rego
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
				subject := subjectFactories[subjectKey](itemMap[subjectKey])
				body = append(body, cf(subject, itemMap[conditionKey]))
			} else if of != nil {
				body = append(body, of(operatorValue))
			}
		}
		return AnyOf(body)
	}
	operatorFactories[not] = func(input any) Rego {
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
		subject := subjectFactories[subjectKey](itemMap[subjectKey])
		return NotOperator{
			Body: cf(subject, itemMap[conditionKey]),
		}
	}
}

var _ Rego = &NotOperator{}

type NotOperator struct {
	Body Rego
}

func (n NotOperator) Rego(ctx context.Context) (string, error) {
	panic("implement me")
}

var _ Rego = &AllOf{}

type AllOf []Rego

func (a AllOf) Rego(ctx context.Context) (string, error) {
	var res string
	for _, item := range a {
		condition, err := item.Rego(ctx)
		if err != nil {
			return "", err
		}
		if res != "" {
			res = res + "\n"
		}
		res += condition
	}
	return res, nil
}

var _ Rego = &AnyOf{}

type AnyOf []Rego

func (a AnyOf) Rego(ctx context.Context) (string, error) {
	var res string
	for _, item := range a {
		if reflect.TypeOf(item) == reflect.TypeOf(EqualsOperation{}) {

		}
		condition, err := item.Rego(ctx)
		if err != nil {
			return "", err
		}
		if res != "" {
			res = res + "\n"
		}
		res += not + " " + condition
	}
	return res, nil
}
