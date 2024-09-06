package pkg

import (
	"context"
	"fmt"
	"strings"
)

type OperationValue string
type OperationField string

func (o OperationValue) Rego(ctx context.Context) (string, error) {
	processed, _, err := FieldNameProcessor(string(o), ctx)
	return processed, err
}

func (o OperationField) Rego(ctx context.Context) (string, error) {
	processed, _, err := FieldNameProcessor(string(o), ctx)
	return processed, err
}

type Rego interface {
	Rego(ctx context.Context) (string, error)
}

type operation struct {
	Subject Rego
}

var conditionFactory = map[string]func(Rego, any) Rego{
	//"anyOf": func(s Rego, input any) Rego {
	//	return AnyOf(input.([]Rego))
	//},
	//"not": func(s Rego, input any) Rego {
	//	return NotOperator{
	//		Body: input.(Rego),
	//	}
	//},
	"equals": func(subject Rego, input any) Rego {
		return EqualsOperation{
			operation: operation{Subject: subject},
			Value:     input.(string),
		}
	},
	"notEquals": func(s Rego, input any) Rego {
		return NotEqualsOperation{
			operation: operation{Subject: s},
			Value:     input.(string),
		}
	},
	"like": func(s Rego, input any) Rego {
		return LikeOperation{
			operation: operation{Subject: s},
			Value:     input.(string),
		}
	},
	"notLike": func(s Rego, input any) Rego {
		return NotLikeOperation{
			operation: operation{Subject: s},
			Value:     input.(string),
		}
	},
	"match": func(s Rego, input any) Rego {
		return MatchOperation{
			operation: operation{Subject: s},
			Value:     input.(string),
		}
	},
	"matchInsensitively": func(s Rego, input any) Rego {
		return MatchInsensitivelyOperation{
			operation: operation{Subject: s},
			Value:     input.(string),
		}
	},
	"notMatch": func(s Rego, input any) Rego {
		return NotMatchOperation{
			operation: operation{Subject: s},
			Value:     input.(string),
		}
	},
	"notMatchInsensitively": func(s Rego, input any) Rego {
		return NotMatchInsensitivelyOperation{
			operation: operation{Subject: s},
			Value:     input.(string),
		}
	},
	"contains": func(s Rego, input any) Rego {
		return ContainsOperation{
			operation: operation{Subject: s},
			Value:     input.(string),
		}
	},
	"notContains": func(s Rego, input any) Rego {
		return NotContainsOperation{
			operation: operation{Subject: s},
			Value:     input.(string),
		}
	},
	"in": func(s Rego, input any) Rego {
		return InOperation{
			operation: operation{Subject: s},
			Values: func() []string {
				var v []string
				for _, i := range input.([]any) {
					v = append(v, i.(string))
				}
				return v
			}(),
		}
	},
	"notIn": func(s Rego, input any) Rego {
		return NotInOperation{
			operation: operation{Subject: s},
			Values: func() []string {
				var v []string
				for _, i := range input.([]any) {
					v = append(v, i.(string))
				}
				return v
			}(),
		}
	},
	"containsKey": func(s Rego, input any) Rego {
		return ContainsKeyOperation{
			operation: operation{Subject: s},
			KeyName:   input.(string),
		}
	},
	"notContainsKey": func(s Rego, input any) Rego {
		return NotContainsKeyOperation{
			operation: operation{Subject: s},
			KeyName:   input.(string),
		}
	},
	"less": func(s Rego, input any) Rego {
		return LessOperation{
			operation: operation{Subject: s},
			Value:     input,
		}
	},
	"lessOrEquals": func(s Rego, input any) Rego {
		return LessOrEqualsOperation{
			operation: operation{Subject: s},
			Value:     input,
		}
	},
	"greater": func(s Rego, input any) Rego {
		return GreaterOperation{
			operation: operation{Subject: s},
			Value:     input,
		}
	},
	"greaterOrEquals": func(s Rego, input any) Rego {
		return GreaterOrEqualsOperation{
			operation: operation{Subject: s},
			Value:     input,
		}
	},
	"exists": func(s Rego, input any) Rego {
		return ExistsOperation{
			operation: operation{Subject: s},
			Value:     input.(bool),
		}
	},
}

var _ Rego = EqualsOperation{}

type EqualsOperation struct {
	operation
	Value string
}

func (e EqualsOperation) Rego(ctx context.Context) (string, error) {
	fieldName, err := e.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{fieldName, "==", fmt.Sprint(e.Value)}, " "), nil
}

var _ Rego = NotEqualsOperation{}

type NotEqualsOperation struct {
	operation
	Value string
}

func (n NotEqualsOperation) Rego(ctx context.Context) (string, error) {
	fieldName, err := n.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{fieldName, "!=", fmt.Sprint(n.Value)}, " "), nil
}

var _ Rego = LikeOperation{}

type LikeOperation struct {
	operation
	Value string
}

func (l LikeOperation) Rego(ctx context.Context) (string, error) {
	fieldName, err := l.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{regexExp, "(", fmt.Sprint(l.Value), ",", fieldName, ")"}, ""), nil
}

var _ Rego = NotLikeOperation{}

type NotLikeOperation struct {
	operation
	Value string
}

func (n NotLikeOperation) Rego(ctx context.Context) (string, error) {
	fieldName, err := n.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{not, " ", regexExp, "(", fmt.Sprint(n.Value), ",", fieldName, ")"}, ""), nil
}

var _ Rego = MatchOperation{}

type MatchOperation struct {
	operation
	Value string
}

func (m MatchOperation) Rego(ctx context.Context) (string, error) {
	return "", fmt.Errorf("`match` condition is not supported, yet")
}

var _ Rego = MatchInsensitivelyOperation{}

type MatchInsensitivelyOperation struct {
	operation
	Value string
}

func (m MatchInsensitivelyOperation) Rego(context.Context) (string, error) {
	return "", fmt.Errorf("`matchInsensitively` condition is not supported, yet")
}

var _ Rego = NotMatchOperation{}

type NotMatchOperation struct {
	operation
	Value string
}

func (n NotMatchOperation) Rego(context.Context) (string, error) {
	return "", fmt.Errorf("`notMatch` condition is not supported, yet")
}

var _ Rego = NotMatchInsensitivelyOperation{}

type NotMatchInsensitivelyOperation struct {
	operation
	Value string
}

func (n NotMatchInsensitivelyOperation) Rego(context.Context) (string, error) {
	return "", fmt.Errorf("`notMatchInsensitively` condition is not supported, yet")
}

var _ Rego = ContainsOperation{}

type ContainsOperation struct {
	operation
	Value string
}

func (c ContainsOperation) Rego(ctx context.Context) (string, error) {
	fieldName, err := c.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{regexExp, "(", "\"", ".*", fmt.Sprint(c.Value), ".*", "\"", ",", fieldName, ")"}, ""), nil

}

var _ Rego = NotContainsOperation{}

type NotContainsOperation struct {
	operation
	Value string
}

func (n NotContainsOperation) Rego(context.Context) (string, error) {
	return "", fmt.Errorf("`notContains` condition is not supported, yet")
}

var _ Rego = InOperation{}

type InOperation struct {
	operation
	Values []string
}

func (i InOperation) Rego(ctx context.Context) (string, error) {
	fieldName, err := i.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{"some", fieldName, "in", SliceConstructor(i.Values)}, " "), nil
}

var _ Rego = NotInOperation{}

type NotInOperation struct {
	operation
	Values []string
}

func (n NotInOperation) Rego(ctx context.Context) (string, error) {
	fieldName, err := n.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{"not", fieldName, "in", SliceConstructor(n.Values)}, " "), nil
}

var _ Rego = ContainsKeyOperation{}

type ContainsKeyOperation struct {
	operation
	KeyName string
}

func (c ContainsKeyOperation) Rego(context.Context) (string, error) {
	return "", fmt.Errorf("`containsKey` condition is not supported, yet")
}

var _ Rego = NotContainsKeyOperation{}

type NotContainsKeyOperation struct {
	operation
	KeyName string
}

func (n NotContainsKeyOperation) Rego(context.Context) (string, error) {
	return "", fmt.Errorf("`notContainsKey` condition is not supported, yet")
}

var _ Rego = LessOperation{}

type LessOperation struct {
	operation
	Value any
}

func (l LessOperation) Rego(ctx context.Context) (string, error) {
	fieldName, err := l.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{fieldName, "<", fmt.Sprint(l.Value)}, " "), nil
}

var _ Rego = LessOrEqualsOperation{}

type LessOrEqualsOperation struct {
	operation
	Value any
}

func (l LessOrEqualsOperation) Rego(ctx context.Context) (string, error) {
	fieldName, err := l.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{fieldName, "<=", fmt.Sprint(l.Value)}, " "), nil
}

var _ Rego = GreaterOperation{}

type GreaterOperation struct {
	operation
	Value any
}

func (g GreaterOperation) Rego(ctx context.Context) (string, error) {
	fieldName, err := g.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{fieldName, ">", fmt.Sprint(g.Value)}, " "), nil
}

var _ Rego = GreaterOrEqualsOperation{}

type GreaterOrEqualsOperation struct {
	operation
	Value any
}

func (g GreaterOrEqualsOperation) Rego(ctx context.Context) (string, error) {
	fieldName, err := g.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{fieldName, ">=", fmt.Sprint(g.Value)}, " "), nil
}

var _ Rego = ExistsOperation{}

type ExistsOperation struct {
	operation
	Value bool
}

func (e ExistsOperation) Rego(ctx context.Context) (string, error) {
	fieldName, err := e.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	if e.Value {
		return fieldName, nil
	} else {
		return strings.Join([]string{not, fieldName}, " "), nil
	}
}

var _ Rego = CountOperation{}

type CountOperation struct {
	field          string
	whereCondition Rego
	operation
}

func (c CountOperation) Rego(ctx context.Context) (string, error) {
	fieldName, err := c.operation.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}

	conditions, err := c.whereCondition.Rego(ctx)
	if err != nil {
		return "", err
	}

	whereConditionName := c.whereCondition.(WhereOperator).ConditionSetName

	res := fmt.Sprintf(count+"{"+fieldName+" | %s}", whereConditionName)
	res += "\n" + conditions

	return res, nil
}
