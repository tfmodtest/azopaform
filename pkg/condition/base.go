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

func NewCondition(conditionType string, subject shared.Rego, value any, ctx *shared.Context) shared.Rego {
	if cf, ok := ConditionFactory[conditionType]; ok {
		return cf(subject, value, ctx)
	}
	return nil
}

var ConditionFactory = map[string]func(shared.Rego, any, *shared.Context) shared.Rego{
	"equals": func(s shared.Rego, input any, ctx *shared.Context) shared.Rego {
		return Equals{
			BaseCondition: BaseCondition{Subject: s},
			Value:         shared.ResolveParameterValue[any](input, ctx),
		}
	},
	"notequals": func(s shared.Rego, input any, ctx *shared.Context) shared.Rego {
		return NotEquals{
			BaseCondition: BaseCondition{Subject: s},
			Value:         shared.ResolveParameterValue[any](input, ctx),
		}
	},
	"like": func(s shared.Rego, input any, ctx *shared.Context) shared.Rego {
		return Like{
			BaseCondition: BaseCondition{Subject: s},
			Value:         shared.ResolveParameterValue[string](input, ctx),
		}
	},
	"notlike": func(s shared.Rego, input any, ctx *shared.Context) shared.Rego {
		return NotLike{
			BaseCondition: BaseCondition{Subject: s},
			Value:         shared.ResolveParameterValue[string](input, ctx),
		}
	},
	"match": func(s shared.Rego, input any, ctx *shared.Context) shared.Rego {
		return Match{
			BaseCondition: BaseCondition{Subject: s},
			Value:         shared.ResolveParameterValue[string](input, ctx),
		}
	},
	"matchinsensitively": func(s shared.Rego, input any, ctx *shared.Context) shared.Rego {
		return MatchInsensitivelyCondition{
			BaseCondition: BaseCondition{Subject: s},
			Value:         shared.ResolveParameterValue[string](input, ctx),
		}
	},
	"notmatch": func(s shared.Rego, input any, ctx *shared.Context) shared.Rego {
		return NotMatch{
			BaseCondition: BaseCondition{Subject: s},
			Value:         shared.ResolveParameterValue[string](input, ctx),
		}
	},
	"notmatchinsensitively": func(s shared.Rego, input any, ctx *shared.Context) shared.Rego {
		return NotMatchInsensitively{
			BaseCondition: BaseCondition{Subject: s},
			Value:         shared.ResolveParameterValue[string](input, ctx),
		}
	},
	"contains": func(s shared.Rego, input any, ctx *shared.Context) shared.Rego {
		return Contains{
			BaseCondition: BaseCondition{Subject: s},
			Value:         shared.ResolveParameterValue[string](input, ctx),
		}
	},
	"notcontains": func(s shared.Rego, input any, ctx *shared.Context) shared.Rego {
		return NotContains{
			BaseCondition: BaseCondition{Subject: s},
			Value:         shared.ResolveParameterValue[string](input, ctx),
		}
	},
	"in": func(s shared.Rego, input any, ctx *shared.Context) shared.Rego {
		return In{
			BaseCondition: BaseCondition{Subject: s},
			Values: func() []string {
				var v []string
				for _, i := range shared.ResolveParameterValue[[]any](input, ctx) {
					v = append(v, i.(string))
				}
				return v
			}(),
		}
	},
	"notin": func(s shared.Rego, input any, ctx *shared.Context) shared.Rego {
		return NotIn{
			BaseCondition: BaseCondition{Subject: s},
			Values: func() []string {
				var v []string
				for _, i := range shared.ResolveParameterValue[[]any](input, ctx) {
					v = append(v, i.(string))
				}
				return v
			}(),
		}
	},
	"containskey": func(s shared.Rego, input any, ctx *shared.Context) shared.Rego {
		return ContainsKey{
			BaseCondition: BaseCondition{Subject: s},
			KeyName:       shared.ResolveParameterValue[string](input, ctx),
		}
	},
	"notcontainskey": func(s shared.Rego, input any, ctx *shared.Context) shared.Rego {
		return NotContainsKey{
			BaseCondition: BaseCondition{Subject: s},
			KeyName:       shared.ResolveParameterValue[string](input, ctx),
		}
	},
	"less": func(s shared.Rego, input any, ctx *shared.Context) shared.Rego {
		return Less{
			BaseCondition: BaseCondition{Subject: s},
			Value:         shared.ResolveParameterValue[any](input, ctx),
		}
	},
	"lessorequals": func(s shared.Rego, input any, ctx *shared.Context) shared.Rego {
		return LessOrEquals{
			BaseCondition: BaseCondition{Subject: s},
			Value:         shared.ResolveParameterValue[any](input, ctx),
		}
	},
	"greater": func(s shared.Rego, input any, ctx *shared.Context) shared.Rego {
		return Greater{
			BaseCondition: BaseCondition{Subject: s},
			Value:         shared.ResolveParameterValue[any](input, ctx),
		}
	},
	"greaterorequals": func(s shared.Rego, input any, ctx *shared.Context) shared.Rego {
		return GreaterOrEquals{
			BaseCondition: BaseCondition{Subject: s},
			Value:         shared.ResolveParameterValue[any](input, ctx),
		}
	},
	"exists": func(s shared.Rego, input any, ctx *shared.Context) shared.Rego {
		return Exists{
			BaseCondition: BaseCondition{Subject: s},
			Value:         shared.ResolveParameterValue[any](input, ctx),
		}
	},
}
