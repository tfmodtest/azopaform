package pkg

import (
	"fmt"
	"json-rule-finder/pkg/condition"
	"json-rule-finder/pkg/shared"
	"strings"
)

func (i *If) Rego(ctx *shared.Context) (string, error) {
	if i.rego != nil {
		return i.rego.Rego(ctx)
	}
	conditionMap := i.body
	var subject shared.Rego
	var creator func(subject shared.Rego, input any) shared.Rego
	var cv any
	for key, conditionValue := range conditionMap {
		key = strings.ToLower(key)
		operation := NewOperation(key, conditionValue, ctx)
		if operation != nil {
			i.rego = operation
			return i.rego.Rego(ctx)
		}
		if key == shared.Count {
			continue
		}
		if key == shared.Field {
			if conditionValue == shared.TypeOfResource {
				ctx.PushResourceType(conditionValue.(string))
			}
			subject = NewSubject(shared.Field, conditionValue, ctx)
			continue
		}
		if key == shared.Value {
			subject = NewSubject(shared.Value, conditionValue, ctx)
			continue
		}
		factory, ok := condition.ConditionFactory[key]
		if !ok {
			panic(fmt.Sprintf("unknown condition: %s", key))
		}
		creator = factory
		cv = conditionValue
	}
	i.rego = creator(subject, cv)
	return i.rego.Rego(ctx)
}

func (i *If) ConditionName(defaultConditionName string) string {
	if operator, ok := i.rego.(Operation); ok {
		return operator.GetConditionSetName()
	}
	return defaultConditionName
}

var _ shared.Rego = &If{}

type If struct {
	body IfBody
	rego shared.Rego
}
