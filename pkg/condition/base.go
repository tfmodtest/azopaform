package condition

import (
	"github.com/tfmodtest/azopaform/pkg/shared"
)

type BaseCondition struct {
	Subject shared.Rego
}

func (b BaseCondition) GetSubject(ctx *shared.Context) shared.Rego {
	localName, ok := ctx.FieldNameReplacer()
	if ok {
		return shared.StringRego(localName)
	}
	return b.Subject
}

func NewCondition(conditionType string, subject shared.Rego, value any, ctx *shared.Context) (shared.Rego, error) {
	if cf, ok := ConditionFactory[conditionType]; ok {
		return cf(subject, value, ctx)
	}
	return nil, nil
}

var ConditionFactory = map[string]func(shared.Rego, any, *shared.Context) (shared.Rego, error){
	"equals": func(s shared.Rego, input any, ctx *shared.Context) (shared.Rego, error) {
		value, err := shared.ResolveParameterValue[any](input, ctx)
		return Equals{
			BaseCondition: BaseCondition{Subject: s},
			Value:         value,
		}, err
	},
	"notequals": func(s shared.Rego, input any, ctx *shared.Context) (shared.Rego, error) {
		value, err := shared.ResolveParameterValue[any](input, ctx)
		return NotEquals{
			BaseCondition: BaseCondition{Subject: s},
			Value:         value,
		}, err
	},
	"like": func(s shared.Rego, input any, ctx *shared.Context) (shared.Rego, error) {
		value, err := shared.ResolveParameterValue[string](input, ctx)
		return Like{
			BaseCondition: BaseCondition{Subject: s},
			Value:         value,
		}, err
	},
	"notlike": func(s shared.Rego, input any, ctx *shared.Context) (shared.Rego, error) {
		value, err := shared.ResolveParameterValue[string](input, ctx)
		return NotLike{
			BaseCondition: BaseCondition{Subject: s},
			Value:         value,
		}, err
	},
	"match": func(s shared.Rego, input any, ctx *shared.Context) (shared.Rego, error) {
		value, err := shared.ResolveParameterValue[string](input, ctx)
		return Match{
			BaseCondition: BaseCondition{Subject: s},
			Value:         value,
		}, err
	},
	"matchinsensitively": func(s shared.Rego, input any, ctx *shared.Context) (shared.Rego, error) {
		value, err := shared.ResolveParameterValue[string](input, ctx)
		return MatchInsensitivelyCondition{
			BaseCondition: BaseCondition{Subject: s},
			Value:         value,
		}, err
	},
	"notmatch": func(s shared.Rego, input any, ctx *shared.Context) (shared.Rego, error) {
		value, err := shared.ResolveParameterValue[string](input, ctx)
		return NotMatch{
			BaseCondition: BaseCondition{Subject: s},
			Value:         value,
		}, err
	},
	"notmatchinsensitively": func(s shared.Rego, input any, ctx *shared.Context) (shared.Rego, error) {
		value, err := shared.ResolveParameterValue[string](input, ctx)
		return NotMatchInsensitively{
			BaseCondition: BaseCondition{Subject: s},
			Value:         value,
		}, err
	},
	"contains": func(s shared.Rego, input any, ctx *shared.Context) (shared.Rego, error) {
		value, err := shared.ResolveParameterValue[string](input, ctx)
		return Contains{
			BaseCondition: BaseCondition{Subject: s},
			Value:         value,
		}, err
	},
	"notcontains": func(s shared.Rego, input any, ctx *shared.Context) (shared.Rego, error) {
		value, err := shared.ResolveParameterValue[string](input, ctx)
		return NotContains{
			BaseCondition: BaseCondition{Subject: s},
			Value:         value,
		}, err
	},
	"in": func(s shared.Rego, input any, ctx *shared.Context) (shared.Rego, error) {
		value, err := func() ([]string, error) {
			var v []string
			items, err := shared.ResolveParameterValue[[]any](input, ctx)
			if err != nil {
				return nil, err
			}
			for _, i := range items {
				v = append(v, i.(string))
			}
			return v, nil
		}()
		return In{
			BaseCondition: BaseCondition{Subject: s},
			Values:        value,
		}, err
	},
	"notin": func(s shared.Rego, input any, ctx *shared.Context) (shared.Rego, error) {
		value, err := func() ([]string, error) {
			var v []string
			items, err := shared.ResolveParameterValue[[]any](input, ctx)
			if err != nil {
				return nil, err
			}
			for _, i := range items {
				v = append(v, i.(string))
			}
			return v, nil
		}()
		return NotIn{
			BaseCondition: BaseCondition{Subject: s},
			Values:        value,
		}, err
	},
	"containskey": func(s shared.Rego, input any, ctx *shared.Context) (shared.Rego, error) {
		value, err := shared.ResolveParameterValue[string](input, ctx)
		return ContainsKey{
			BaseCondition: BaseCondition{Subject: s},
			KeyName:       value,
		}, err
	},
	"notcontainskey": func(s shared.Rego, input any, ctx *shared.Context) (shared.Rego, error) {
		value, err := shared.ResolveParameterValue[string](input, ctx)
		return NotContainsKey{
			BaseCondition: BaseCondition{Subject: s},
			KeyName:       value,
		}, err
	},
	"less": func(s shared.Rego, input any, ctx *shared.Context) (shared.Rego, error) {
		value, err := shared.ResolveParameterValue[any](input, ctx)
		return Less{
			BaseCondition: BaseCondition{Subject: s},
			Value:         value,
		}, err
	},
	"lessorequals": func(s shared.Rego, input any, ctx *shared.Context) (shared.Rego, error) {
		value, err := shared.ResolveParameterValue[any](input, ctx)
		return LessOrEquals{
			BaseCondition: BaseCondition{Subject: s},
			Value:         value,
		}, err
	},
	"greater": func(s shared.Rego, input any, ctx *shared.Context) (shared.Rego, error) {
		value, err := shared.ResolveParameterValue[any](input, ctx)
		return Greater{
			BaseCondition: BaseCondition{Subject: s},
			Value:         value,
		}, err
	},
	"greaterorequals": func(s shared.Rego, input any, ctx *shared.Context) (shared.Rego, error) {
		value, err := shared.ResolveParameterValue[any](input, ctx)
		return GreaterOrEquals{
			BaseCondition: BaseCondition{Subject: s},
			Value:         value,
		}, err
	},
	"exists": func(s shared.Rego, input any, ctx *shared.Context) (shared.Rego, error) {
		value, err := shared.ResolveParameterValue[any](input, ctx)
		return Exists{
			BaseCondition: BaseCondition{Subject: s},
			Value:         value,
		}, err
	},
}
