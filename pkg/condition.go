package pkg

import (
	"context"
	"fmt"
	"github.com/emirpasic/gods/stacks"
	"reflect"
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

type Condition interface {
	GetReverseRego(ctx context.Context) (string, error)
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
			Value:     input,
		}
	},
	"notequals": func(s Rego, input any) Rego {
		return NotEqualsOperation{
			operation: operation{Subject: s},
			Value:     input,
		}
	},
	"like": func(s Rego, input any) Rego {
		return LikeOperation{
			operation: operation{Subject: s},
			Value:     input.(string),
		}
	},
	"notlike": func(s Rego, input any) Rego {
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
	"matchinsensitively": func(s Rego, input any) Rego {
		return MatchInsensitivelyOperation{
			operation: operation{Subject: s},
			Value:     input.(string),
		}
	},
	"notmatch": func(s Rego, input any) Rego {
		return NotMatchOperation{
			operation: operation{Subject: s},
			Value:     input.(string),
		}
	},
	"notmatchinsensitively": func(s Rego, input any) Rego {
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
	"notcontains": func(s Rego, input any) Rego {
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
	"containskey": func(s Rego, input any) Rego {
		return ContainsKeyOperation{
			operation: operation{Subject: s},
			KeyName:   input.(string),
		}
	},
	"notcontainskey": func(s Rego, input any) Rego {
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
	"lessorequals": func(s Rego, input any) Rego {
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
	"greaterorequals": func(s Rego, input any) Rego {
		return GreaterOrEqualsOperation{
			operation: operation{Subject: s},
			Value:     input,
		}
	},
	"exists": func(s Rego, input any) Rego {
		return ExistsOperation{
			operation: operation{Subject: s},
			Value:     input,
		}
	},
}

var _ Rego = EqualsOperation{}
var _ Condition = EqualsOperation{}

type EqualsOperation struct {
	operation
	Value any
}

// Rego For conditions under 'where' operator, "[[0-9]+]" should be replaced with "[x]"
func (e EqualsOperation) Rego(ctx context.Context) (string, error) {
	fieldName, err := e.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
		fieldName = replaceIndex(fieldName)
	}
	var v string
	if reflect.TypeOf(e.Value).Kind() == reflect.String {
		v = strings.Join([]string{"\"", fmt.Sprint(e.Value), "\""}, "")
	} else if reflect.TypeOf(e.Value).Kind() == reflect.Bool {
		v = fmt.Sprint(e.Value)
	} else {
		v = fmt.Sprint(e.Value)
	}
	return strings.Join([]string{fieldName, "==", v}, " "), nil
}

func (e EqualsOperation) GetReverseRego(ctx context.Context) (string, error) {
	fieldName, err := e.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
		fieldName = replaceIndex(fieldName)
	}
	var v string
	if reflect.TypeOf(e.Value).Kind() == reflect.String {
		v = strings.Join([]string{"\"", fmt.Sprint(e.Value), "\""}, "")
	} else if reflect.TypeOf(e.Value).Kind() == reflect.Bool {
		v = fmt.Sprint(e.Value)
	} else {
		v = fmt.Sprint(e.Value)
	}
	return strings.Join([]string{fieldName, "!=", v}, " "), nil
}

var _ Rego = NotEqualsOperation{}
var _ Condition = NotEqualsOperation{}

type NotEqualsOperation struct {
	operation
	Value any
}

func (n NotEqualsOperation) Rego(ctx context.Context) (string, error) {
	fieldName, err := n.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
		fieldName = replaceIndex(fieldName)
	}
	var v string
	if reflect.TypeOf(n.Value).Kind() == reflect.String {
		v = strings.Join([]string{"\"", fmt.Sprint(n.Value), "\""}, "")
	} else if reflect.TypeOf(n.Value).Kind() == reflect.Bool {
		v = fmt.Sprint(n.Value)
	} else {
		v = fmt.Sprint(n.Value)
	}
	return strings.Join([]string{fieldName, "!=", v}, " "), nil
}

func (n NotEqualsOperation) GetReverseRego(ctx context.Context) (string, error) {
	fieldName, err := n.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
		fieldName = replaceIndex(fieldName)
	}
	var v string
	if reflect.TypeOf(n.Value).Kind() == reflect.String {
		v = strings.Join([]string{"\"", fmt.Sprint(n.Value), "\""}, "")
	} else if reflect.TypeOf(n.Value).Kind() == reflect.Bool {
		v = fmt.Sprint(n.Value)
	} else {
		v = fmt.Sprint(n.Value)
	}
	return strings.Join([]string{fieldName, "==", v}, " "), nil
}

var _ Rego = LikeOperation{}
var _ Condition = LikeOperation{}

type LikeOperation struct {
	operation
	Value string
}

func (l LikeOperation) Rego(ctx context.Context) (string, error) {
	fieldName, err := l.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
		fieldName = replaceIndex(fieldName)
	}
	v := strings.Join([]string{"\"", fmt.Sprint(l.Value), "\""}, "")
	return strings.Join([]string{regexExp, "(", v, ",", fieldName, ")"}, ""), nil
}

func (l LikeOperation) GetReverseRego(ctx context.Context) (string, error) {
	fieldName, err := l.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
		fieldName = replaceIndex(fieldName)
	}
	v := strings.Join([]string{"\"", fmt.Sprint(l.Value), "\""}, "")
	return strings.Join([]string{not, " ", regexExp, "(", v, ",", fieldName, ")"}, ""), nil
}

var _ Rego = NotLikeOperation{}
var _ Condition = NotLikeOperation{}

type NotLikeOperation struct {
	operation
	Value string
}

func (n NotLikeOperation) Rego(ctx context.Context) (string, error) {
	fieldName, err := n.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
		fieldName = replaceIndex(fieldName)
	}
	v := strings.Join([]string{"`", fmt.Sprint(n.Value), "`"}, "")
	return strings.Join([]string{not, " ", regexExp, "(", v, ",", fieldName, ")"}, ""), nil
}

func (n NotLikeOperation) GetReverseRego(ctx context.Context) (string, error) {
	fieldName, err := n.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
		fieldName = replaceIndex(fieldName)
	}
	v := strings.Join([]string{"`", fmt.Sprint(n.Value), "`"}, "")
	return strings.Join([]string{regexExp, "(", v, ",", fieldName, ")"}, ""), nil
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
var _ Condition = InOperation{}

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

func (i InOperation) GetReverseRego(ctx context.Context) (string, error) {
	fieldName, err := i.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{not, fieldName, "in", SliceConstructor(i.Values)}, " "), nil
}

var _ Rego = NotInOperation{}
var _ Condition = NotInOperation{}

type NotInOperation struct {
	operation
	Values []string
}

func (n NotInOperation) Rego(ctx context.Context) (string, error) {
	fieldName, err := n.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{not, fieldName, "in", SliceConstructor(n.Values)}, " "), nil
}

func (n NotInOperation) GetReverseRego(ctx context.Context) (string, error) {
	fieldName, err := n.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{"some", fieldName, "in", SliceConstructor(n.Values)}, " "), nil
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
var _ Condition = LessOperation{}

type LessOperation struct {
	operation
	Value any
}

func (l LessOperation) Rego(ctx context.Context) (string, error) {
	fieldName, err := l.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
		fieldName = replaceIndex(fieldName)
	}
	return strings.Join([]string{fieldName, "<", fmt.Sprint(l.Value)}, " "), nil
}

func (l LessOperation) GetReverseRego(ctx context.Context) (string, error) {
	fieldName, err := l.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
		fieldName = replaceIndex(fieldName)
	}
	return strings.Join([]string{fieldName, ">=", fmt.Sprint(l.Value)}, " "), nil
}

var _ Rego = LessOrEqualsOperation{}
var _ Condition = LessOrEqualsOperation{}

type LessOrEqualsOperation struct {
	operation
	Value any
}

func (l LessOrEqualsOperation) Rego(ctx context.Context) (string, error) {
	fieldName, err := l.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
		fieldName = replaceIndex(fieldName)
	}
	return strings.Join([]string{fieldName, "<=", fmt.Sprint(l.Value)}, " "), nil
}

func (l LessOrEqualsOperation) GetReverseRego(ctx context.Context) (string, error) {
	fieldName, err := l.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
		fieldName = replaceIndex(fieldName)
	}
	return strings.Join([]string{fieldName, ">", fmt.Sprint(l.Value)}, " "), nil
}

var _ Rego = GreaterOperation{}
var _ Condition = GreaterOperation{}

type GreaterOperation struct {
	operation
	Value any
}

func (g GreaterOperation) Rego(ctx context.Context) (string, error) {
	fieldName, err := g.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
		fieldName = replaceIndex(fieldName)
	}
	return strings.Join([]string{fieldName, ">", fmt.Sprint(g.Value)}, " "), nil
}

func (g GreaterOperation) GetReverseRego(ctx context.Context) (string, error) {
	fieldName, err := g.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
		fieldName = replaceIndex(fieldName)
	}
	return strings.Join([]string{fieldName, "<=", fmt.Sprint(g.Value)}, " "), nil
}

var _ Rego = GreaterOrEqualsOperation{}
var _ Condition = GreaterOrEqualsOperation{}

type GreaterOrEqualsOperation struct {
	operation
	Value any
}

func (g GreaterOrEqualsOperation) Rego(ctx context.Context) (string, error) {
	fieldName, err := g.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
		fieldName = replaceIndex(fieldName)
	}
	return strings.Join([]string{fieldName, ">=", fmt.Sprint(g.Value)}, " "), nil
}

func (g GreaterOrEqualsOperation) GetReverseRego(ctx context.Context) (string, error) {
	fieldName, err := g.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
		fieldName = replaceIndex(fieldName)
	}
	return strings.Join([]string{fieldName, "<", fmt.Sprint(g.Value)}, " "), nil
}

var _ Rego = ExistsOperation{}
var _ Condition = ExistsOperation{}

type ExistsOperation struct {
	operation
	Value any
}

func (e ExistsOperation) Rego(ctx context.Context) (string, error) {
	fieldName, err := e.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
		fieldName = replaceIndex(fieldName)
	}
	if (reflect.TypeOf(e.Value).Kind() == reflect.Bool && e.Value.(bool)) || (reflect.TypeOf(e.Value).Kind() == reflect.String && e.Value.(string) == "true") {
		return fieldName, nil
	} else {
		return strings.Join([]string{not, fieldName}, " "), nil
	}
}

func (e ExistsOperation) GetReverseRego(ctx context.Context) (string, error) {
	fieldName, err := e.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
		fieldName = replaceIndex(fieldName)
	}
	if (reflect.TypeOf(e.Value).Kind() == reflect.Bool && e.Value.(bool)) || (reflect.TypeOf(e.Value).Kind() == reflect.String && e.Value.(string) == "true") {
		return strings.Join([]string{not, fieldName}, " "), nil
	} else {
		return fieldName, nil
	}
}
