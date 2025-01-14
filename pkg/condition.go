package pkg

import (
	"context"
	"fmt"
	"json-rule-finder/pkg/shared"
	"reflect"
)

type OperationValue string
type OperationField string

func (o OperationValue) Rego(ctx context.Context) (string, error) {
	processed, _, err := shared.FieldNameProcessor(string(o), ctx)
	return processed, err
}

func (o OperationField) Rego(ctx context.Context) (string, error) {
	processed, _, err := shared.FieldNameProcessor(string(o), ctx)
	return processed, err
}

type Rego interface {
	Rego(ctx context.Context) (string, error)
}

type Condition interface {
	Rego
}

type Operation interface {
	Rego
}

type condition struct {
	Subject Rego
}

var conditionFactory = map[string]func(Rego, any) Rego{
	"equals": func(s Rego, input any) Rego {
		return EqualsCondition{
			condition: condition{Subject: s},
			Value:     input,
		}
	},
	"notequals": func(s Rego, input any) Rego {
		return NotEqualsCondition{
			condition: condition{Subject: s},
			Value:     input,
		}
	},
	"like": func(s Rego, input any) Rego {
		return LikeCondition{
			condition: condition{Subject: s},
			Value:     input.(string),
		}
	},
	"notlike": func(s Rego, input any) Rego {
		return NotLikeCondition{
			condition: condition{Subject: s},
			Value:     input.(string),
		}
	},
	"match": func(s Rego, input any) Rego {
		return MatchCondition{
			condition: condition{Subject: s},
			Value:     input.(string),
		}
	},
	"matchinsensitively": func(s Rego, input any) Rego {
		return MatchInsensitivelyCondition{
			condition: condition{Subject: s},
			Value:     input.(string),
		}
	},
	"notmatch": func(s Rego, input any) Rego {
		return NotMatchCondition{
			condition: condition{Subject: s},
			Value:     input.(string),
		}
	},
	"notmatchinsensitively": func(s Rego, input any) Rego {
		return NotMatchInsensitivelyCondition{
			condition: condition{Subject: s},
			Value:     input.(string),
		}
	},
	"contains": func(s Rego, input any) Rego {
		return ContainsCondition{
			condition: condition{Subject: s},
			Value:     input.(string),
		}
	},
	"notcontains": func(s Rego, input any) Rego {
		return NotContainsCondition{
			condition: condition{Subject: s},
			Value:     input.(string),
		}
	},
	"in": func(s Rego, input any) Rego {
		return InCondition{
			condition: condition{Subject: s},
			Values: func() []string {
				var v []string
				if reflect.TypeOf(input).Kind() != reflect.Slice {
					return nil
				}
				for _, i := range input.([]any) {
					v = append(v, i.(string))
				}
				return v
			}(),
		}
	},
	"notin": func(s Rego, input any) Rego {
		return NotInCondition{
			condition: condition{Subject: s},
			Values: func() []string {
				var v []string
				for _, i := range input.([]any) {
					v = append(v, i.(string))
				}
				return v
			}(),
		}
	},
	"containskey": func(s Rego, input any) Rego {
		return ContainsKeyCondition{
			condition: condition{Subject: s},
			KeyName:   input.(string),
		}
	},
	"notcontainskey": func(s Rego, input any) Rego {
		return NotContainsKeyCondition{
			condition: condition{Subject: s},
			KeyName:   input.(string),
		}
	},
	"less": func(s Rego, input any) Rego {
		return LessCondition{
			condition: condition{Subject: s},
			Value:     input,
		}
	},
	"lessorequals": func(s Rego, input any) Rego {
		return LessOrEqualsCondition{
			condition: condition{Subject: s},
			Value:     input,
		}
	},
	"greater": func(s Rego, input any) Rego {
		return GreaterCondition{
			condition: condition{Subject: s},
			Value:     input,
		}
	},
	"greaterorequals": func(s Rego, input any) Rego {
		return GreaterOrEqualsCondition{
			condition: condition{Subject: s},
			Value:     input,
		}
	},
	"exists": func(s Rego, input any) Rego {
		return ExistsCondition{
			condition: condition{Subject: s},
			Value:     input,
		}
	},
}

var _ Rego = MatchCondition{}

type MatchCondition struct {
	condition
	Value string
}

func (m MatchCondition) Rego(ctx context.Context) (string, error) {
	return "", fmt.Errorf("`match` condition is not supported, yet")
}

var _ Rego = MatchInsensitivelyCondition{}

type MatchInsensitivelyCondition struct {
	condition
	Value string
}

func (m MatchInsensitivelyCondition) Rego(context.Context) (string, error) {
	return "", fmt.Errorf("`matchInsensitively` condition is not supported, yet")
}

var _ Rego = NotMatchCondition{}

type NotMatchCondition struct {
	condition
	Value string
}

func (n NotMatchCondition) Rego(context.Context) (string, error) {
	return "", fmt.Errorf("`notMatch` condition is not supported, yet")
}

var _ Rego = NotMatchInsensitivelyCondition{}

type NotMatchInsensitivelyCondition struct {
	condition
	Value string
}

func (n NotMatchInsensitivelyCondition) Rego(context.Context) (string, error) {
	return "", fmt.Errorf("`notMatchInsensitively` condition is not supported, yet")
}

var _ Rego = ContainsKeyCondition{}

type ContainsKeyCondition struct {
	condition
	KeyName string
}

func (c ContainsKeyCondition) Rego(context.Context) (string, error) {
	return "", fmt.Errorf("`containsKey` condition is not supported, yet")
}

var _ Rego = NotContainsKeyCondition{}

type NotContainsKeyCondition struct {
	condition
	KeyName string
}

func (n NotContainsKeyCondition) Rego(context.Context) (string, error) {
	return "", fmt.Errorf("`notContainsKey` condition is not supported, yet")
}
