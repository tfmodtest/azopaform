package pkg

import (
	"context"
	"fmt"
	"strings"
)

func (i *If) Rego(ctx context.Context) (string, error) {
	if i.rego == nil {
		i.rego = func() Rego {
			conditionMap := i.body
			var subject Rego
			var creator func(subject Rego, input any) Rego
			var cv any
			for key, conditionValue := range conditionMap {
				key = strings.ToLower(key)
				if key == count {
					operationFactory, ok := operatorFactories[key]
					if !ok {
						panic(fmt.Sprintf("unknown operation: %s", key))
					}
					conditionSet := operationFactory(conditionValue, ctx)
					subject = conditionSet
					continue
				}
				if key == allOf {
					operationFactory, ok := operatorFactories[key]
					if !ok {
						panic(fmt.Sprintf("unknown operation: %s", key))
					}
					conditionSet := operationFactory(conditionValue, ctx)
					return conditionSet
				}
				if key == anyOf {
					operationFactory, ok := operatorFactories[key]
					if !ok {
						panic(fmt.Sprintf("unknown operation: %s", key))
					}
					conditionSet := operationFactory(conditionValue, ctx)
					return conditionSet
				}
				if key == not {
					operationFactory, ok := operatorFactories[key]
					if !ok {
						panic(fmt.Sprintf("unknown operation: %s", key))
					}
					conditionSet := operationFactory(conditionValue, ctx)
					return conditionSet
				}
				if key == field {
					if conditionValue == typeOfResource {
						pushResourceType(ctx, conditionValue.(string))
					}
					subject = OperationField(conditionValue.(string))
					continue
				}
				if key == value {
					subject = OperationValue(conditionValue.(string))
					continue
				}
				factory, ok := conditionFactory[key]
				if !ok {
					panic(fmt.Sprintf("unknown condition: %s", key))
				}
				creator = factory
				cv = conditionValue
			}
			return creator(subject, cv)
		}()

	}
	return i.rego.Rego(ctx)
}

func (i *If) ConditionName(defaultConditionName string) string {
	if operator, ok := i.rego.(Operator); ok {
		return operator.GetConditionSetName()
	}
	return defaultConditionName
}

var _ Rego = &If{}

type If struct {
	body IfBody
	rego Rego
}
