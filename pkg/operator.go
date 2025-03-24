package pkg

import (
	"fmt"
	"json-rule-finder/pkg/shared"
)

type Operator interface {
	shared.Rego
	GetConditionSetName() string
}

var otherFactories map[string]func(input any, ctx *shared.Context) shared.Rego
var operators map[string]func(input any, ctx *shared.Context) shared.Rego

func init() {
	operators = map[string]func(input any, ctx *shared.Context) shared.Rego{
		shared.AllOf: NewAllOf,
		shared.AnyOf: NewAnyOf,
		shared.Not:   NewNot,
	}
	otherFactories = map[string]func(input any, ctx *shared.Context) shared.Rego{
		shared.Count: NewCountOperator,
		shared.Where: NewWhere,
	}
}

var NeoConditionNameGenerator = func(ctx *shared.Context) (string, error) {
	index, ok := ctx.PopConditionNameCounter()
	if !ok {
		return "", fmt.Errorf("conditionNameStack is empty")
	}
	conditionName := "condition" + fmt.Sprintf("%d", index)
	return conditionName, nil
}
