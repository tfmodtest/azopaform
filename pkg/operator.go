package pkg

import "context"

var operatorFactories = make(map[string]func(input any) Rego)

func init() {
	operatorFactories["allOf"] = func(input any) Rego {
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
	//TODO implement me
	panic("implement me")
}

var _ Rego = &AnyOf{}

type AnyOf []Rego

func (a AnyOf) Rego(ctx context.Context) (string, error) {
	//TODO implement me
	panic("implement me")
}
