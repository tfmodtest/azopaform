package condition

import (
	"reflect"

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

func NewCondition(conditionType string, subject shared.Rego, value any) shared.Rego {
	if cf, ok := ConditionFactory[conditionType]; ok {
		return cf(subject, value)
	}
	return nil
}

var ConditionFactory = map[string]func(shared.Rego, any) shared.Rego{
	"equals": func(s shared.Rego, input any) shared.Rego {
		return Equals{
			BaseCondition: BaseCondition{Subject: s},
			Value:         input,
		}
	},
	"notequals": func(s shared.Rego, input any) shared.Rego {
		return NotEquals{
			BaseCondition: BaseCondition{Subject: s},
			Value:         input,
		}
	},
	"like": func(s shared.Rego, input any) shared.Rego {
		return Like{
			BaseCondition: BaseCondition{Subject: s},
			Value:         input.(string),
		}
	},
	"notlike": func(s shared.Rego, input any) shared.Rego {
		return NotLike{
			BaseCondition: BaseCondition{Subject: s},
			Value:         input.(string),
		}
	},
	"match": func(s shared.Rego, input any) shared.Rego {
		return Match{
			BaseCondition: BaseCondition{Subject: s},
			Value:         input.(string),
		}
	},
	"matchinsensitively": func(s shared.Rego, input any) shared.Rego {
		return MatchInsensitivelyCondition{
			BaseCondition: BaseCondition{Subject: s},
			Value:         input.(string),
		}
	},
	"notmatch": func(s shared.Rego, input any) shared.Rego {
		return NotMatch{
			BaseCondition: BaseCondition{Subject: s},
			Value:         input.(string),
		}
	},
	"notmatchinsensitively": func(s shared.Rego, input any) shared.Rego {
		return NotMatchInsensitively{
			BaseCondition: BaseCondition{Subject: s},
			Value:         input.(string),
		}
	},
	"contains": func(s shared.Rego, input any) shared.Rego {
		return Contains{
			BaseCondition: BaseCondition{Subject: s},
			Value:         input.(string),
		}
	},
	"notcontains": func(s shared.Rego, input any) shared.Rego {
		return NotContains{
			BaseCondition: BaseCondition{Subject: s},
			Value:         input.(string),
		}
	},
	"in": func(s shared.Rego, input any) shared.Rego {
		return In{
			BaseCondition: BaseCondition{Subject: s},
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
	"notin": func(s shared.Rego, input any) shared.Rego {
		return NotIn{
			BaseCondition: BaseCondition{Subject: s},
			Values: func() []string {
				var v []string
				for _, i := range input.([]any) {
					v = append(v, i.(string))
				}
				return v
			}(),
		}
	},
	"containskey": func(s shared.Rego, input any) shared.Rego {
		return ContainsKey{
			BaseCondition: BaseCondition{Subject: s},
			KeyName:       input.(string),
		}
	},
	"notcontainskey": func(s shared.Rego, input any) shared.Rego {
		return NotContainsKey{
			BaseCondition: BaseCondition{Subject: s},
			KeyName:       input.(string),
		}
	},
	"less": func(s shared.Rego, input any) shared.Rego {
		return Less{
			BaseCondition: BaseCondition{Subject: s},
			Value:         input,
		}
	},
	"lessorequals": func(s shared.Rego, input any) shared.Rego {
		return LessOrEquals{
			BaseCondition: BaseCondition{Subject: s},
			Value:         input,
		}
	},
	"greater": func(s shared.Rego, input any) shared.Rego {
		return Greater{
			BaseCondition: BaseCondition{Subject: s},
			Value:         input,
		}
	},
	"greaterorequals": func(s shared.Rego, input any) shared.Rego {
		return GreaterOrEquals{
			BaseCondition: BaseCondition{Subject: s},
			Value:         input,
		}
	},
	"exists": func(s shared.Rego, input any) shared.Rego {
		return Exists{
			BaseCondition: BaseCondition{Subject: s},
			Value:         input,
		}
	},
}
