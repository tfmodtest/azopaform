package pkg

import (
	"fmt"
	"json-rule-finder/pkg/condition"
	"json-rule-finder/pkg/shared"
	"strings"
)

func (i *If) Rego(ctx *shared.Context) (string, error) {
	if i.rego == nil {
		i.rego = func() shared.Rego {
			conditionMap := i.body
			var subject shared.Rego
			var creator func(subject shared.Rego, input any) shared.Rego
			var cv any
			for key, conditionValue := range conditionMap {
				key = strings.ToLower(key)
				if key == shared.Count {
					operationFactory, ok := operatorFactories[key]
					if !ok {
						panic(fmt.Sprintf("unknown operation: %s", key))
					}
					conditionSet := operationFactory(conditionValue, ctx)
					subject = conditionSet
					continue
				}
				if key == shared.AllOf {
					operationFactory, ok := operatorFactories[key]
					if !ok {
						panic(fmt.Sprintf("unknown operation: %s", key))
					}
					conditionSet := operationFactory(conditionValue, ctx)
					return conditionSet
				}
				if key == shared.AnyOf {
					operationFactory, ok := operatorFactories[key]
					if !ok {
						panic(fmt.Sprintf("unknown operation: %s", key))
					}
					conditionSet := operationFactory(conditionValue, ctx)
					return conditionSet
				}
				if key == shared.Not {
					operationFactory, ok := operatorFactories[key]
					if !ok {
						panic(fmt.Sprintf("unknown operation: %s", key))
					}
					conditionSet := operationFactory(conditionValue, ctx)
					return conditionSet
				}
				if key == shared.Field {
					if conditionValue == shared.TypeOfResource {
						ctx.PushResourceType(conditionValue.(string))
					}
					subject = OperationField(conditionValue.(string))
					continue
				}
				if key == shared.Value {
					subject = OperationValue(conditionValue.(string))
					continue
				}
				factory, ok := condition.ConditionFactory[key]
				if !ok {
					panic(fmt.Sprintf("unknown BaseCondition: %s", key))
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

var _ shared.Rego = &If{}

type If struct {
	body IfBody
	rego shared.Rego
}
