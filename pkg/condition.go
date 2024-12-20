package pkg

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/emirpasic/gods/stacks"
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
	Rego
	GetReverseRego(ctx context.Context) (string, error)
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

var _ Rego = NotEqualsCondition{}
var _ Condition = NotEqualsCondition{}

type NotEqualsCondition struct {
	condition
	Value any
}

func (n NotEqualsCondition) Rego(ctx context.Context) (string, error) {
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

func (n NotEqualsCondition) GetReverseRego(ctx context.Context) (string, error) {
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

var _ Rego = LikeCondition{}
var _ Condition = LikeCondition{}

type LikeCondition struct {
	condition
	Value string
}

func (l LikeCondition) Rego(ctx context.Context) (string, error) {
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

func (l LikeCondition) GetReverseRego(ctx context.Context) (string, error) {
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

var _ Rego = NotLikeCondition{}
var _ Condition = NotLikeCondition{}

type NotLikeCondition struct {
	condition
	Value string
}

func (n NotLikeCondition) Rego(ctx context.Context) (string, error) {
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

func (n NotLikeCondition) GetReverseRego(ctx context.Context) (string, error) {
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

var _ Rego = ContainsCondition{}

type ContainsCondition struct {
	condition
	Value string
}

func (c ContainsCondition) Rego(ctx context.Context) (string, error) {
	fieldName, err := c.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{regexExp, "(", "\"", ".*", fmt.Sprint(c.Value), ".*", "\"", ",", fieldName, ")"}, ""), nil

}

var _ Rego = NotContainsCondition{}

type NotContainsCondition struct {
	condition
	Value string
}

func (n NotContainsCondition) Rego(context.Context) (string, error) {
	return "", fmt.Errorf("`notContains` condition is not supported, yet")
}

var _ Rego = InCondition{}
var _ Condition = InCondition{}

type InCondition struct {
	condition
	Values []string
}

func (i InCondition) Rego(ctx context.Context) (string, error) {
	fieldName, err := i.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{"some", fieldName, "in", SliceConstructor(i.Values)}, " "), nil
}

func (i InCondition) GetReverseRego(ctx context.Context) (string, error) {
	fieldName, err := i.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{not, fieldName, "in", SliceConstructor(i.Values)}, " "), nil
}

var _ Rego = NotInCondition{}
var _ Condition = NotInCondition{}

type NotInCondition struct {
	condition
	Values []string
}

func (n NotInCondition) Rego(ctx context.Context) (string, error) {
	fieldName, err := n.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{not, fieldName, "in", SliceConstructor(n.Values)}, " "), nil
}

func (n NotInCondition) GetReverseRego(ctx context.Context) (string, error) {
	fieldName, err := n.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{"some", fieldName, "in", SliceConstructor(n.Values)}, " "), nil
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

var _ Rego = LessCondition{}
var _ Condition = LessCondition{}

type LessCondition struct {
	condition
	Value any
}

func (l LessCondition) Rego(ctx context.Context) (string, error) {
	fieldName, err := l.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
		fieldName = replaceIndex(fieldName)
	}
	return strings.Join([]string{fieldName, "<", fmt.Sprint(l.Value)}, " "), nil
}

func (l LessCondition) GetReverseRego(ctx context.Context) (string, error) {
	fieldName, err := l.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
		fieldName = replaceIndex(fieldName)
	}
	return strings.Join([]string{fieldName, ">=", fmt.Sprint(l.Value)}, " "), nil
}

var _ Rego = LessOrEqualsCondition{}
var _ Condition = LessOrEqualsCondition{}

type LessOrEqualsCondition struct {
	condition
	Value any
}

func (l LessOrEqualsCondition) Rego(ctx context.Context) (string, error) {
	fieldName, err := l.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
		fieldName = replaceIndex(fieldName)
	}
	return strings.Join([]string{fieldName, "<=", fmt.Sprint(l.Value)}, " "), nil
}

func (l LessOrEqualsCondition) GetReverseRego(ctx context.Context) (string, error) {
	fieldName, err := l.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
		fieldName = replaceIndex(fieldName)
	}
	return strings.Join([]string{fieldName, ">", fmt.Sprint(l.Value)}, " "), nil
}

var _ Rego = GreaterCondition{}
var _ Condition = GreaterCondition{}

type GreaterCondition struct {
	condition
	Value any
}

func (g GreaterCondition) Rego(ctx context.Context) (string, error) {
	fieldName, err := g.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
		fieldName = replaceIndex(fieldName)
	}
	return strings.Join([]string{fieldName, ">", fmt.Sprint(g.Value)}, " "), nil
}

func (g GreaterCondition) GetReverseRego(ctx context.Context) (string, error) {
	fieldName, err := g.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
		fieldName = replaceIndex(fieldName)
	}
	return strings.Join([]string{fieldName, "<=", fmt.Sprint(g.Value)}, " "), nil
}

var _ Rego = GreaterOrEqualsCondition{}
var _ Condition = GreaterOrEqualsCondition{}

type GreaterOrEqualsCondition struct {
	condition
	Value any
}

func (g GreaterOrEqualsCondition) Rego(ctx context.Context) (string, error) {
	fieldName, err := g.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
		fieldName = replaceIndex(fieldName)
	}
	return strings.Join([]string{fieldName, ">=", fmt.Sprint(g.Value)}, " "), nil
}

func (g GreaterOrEqualsCondition) GetReverseRego(ctx context.Context) (string, error) {
	fieldName, err := g.Subject.Rego(ctx)
	if err != nil {
		return "", err
	}
	if ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"] != nil && ctx.Value("context").(map[string]stacks.Stack)["fieldNameReplacer"].(stacks.Stack).Size() > 0 {
		fieldName = replaceIndex(fieldName)
	}
	return strings.Join([]string{fieldName, "<", fmt.Sprint(g.Value)}, " "), nil
}

var _ Rego = ExistsCondition{}
var _ Condition = ExistsCondition{}

type ExistsCondition struct {
	condition
	Value any
}

func (e ExistsCondition) Rego(ctx context.Context) (string, error) {
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

func (e ExistsCondition) GetReverseRego(ctx context.Context) (string, error) {
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
