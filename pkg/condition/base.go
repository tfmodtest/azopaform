package condition

import (
	"regexp"
	
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
			Value:         ResolveParameterValue[any](input, ctx),
		}
	},
	"notequals": func(s shared.Rego, input any, ctx *shared.Context) shared.Rego {
		return NotEquals{
			BaseCondition: BaseCondition{Subject: s},
			Value:         ResolveParameterValue[any](input, ctx),
		}
	},
	"like": func(s shared.Rego, input any, ctx *shared.Context) shared.Rego {
		return Like{
			BaseCondition: BaseCondition{Subject: s},
			Value:         ResolveParameterValue[string](input, ctx),
		}
	},
	"notlike": func(s shared.Rego, input any, ctx *shared.Context) shared.Rego {
		return NotLike{
			BaseCondition: BaseCondition{Subject: s},
			Value:         ResolveParameterValue[string](input, ctx),
		}
	},
	"match": func(s shared.Rego, input any, ctx *shared.Context) shared.Rego {
		return Match{
			BaseCondition: BaseCondition{Subject: s},
			Value:         ResolveParameterValue[string](input, ctx),
		}
	},
	"matchinsensitively": func(s shared.Rego, input any, ctx *shared.Context) shared.Rego {
		return MatchInsensitivelyCondition{
			BaseCondition: BaseCondition{Subject: s},
			Value:         ResolveParameterValue[string](input, ctx),
		}
	},
	"notmatch": func(s shared.Rego, input any, ctx *shared.Context) shared.Rego {
		return NotMatch{
			BaseCondition: BaseCondition{Subject: s},
			Value:         ResolveParameterValue[string](input, ctx),
		}
	},
	"notmatchinsensitively": func(s shared.Rego, input any, ctx *shared.Context) shared.Rego {
		return NotMatchInsensitively{
			BaseCondition: BaseCondition{Subject: s},
			Value:         ResolveParameterValue[string](input, ctx),
		}
	},
	"contains": func(s shared.Rego, input any, ctx *shared.Context) shared.Rego {
		return Contains{
			BaseCondition: BaseCondition{Subject: s},
			Value:         ResolveParameterValue[string](input, ctx),
		}
	},
	"notcontains": func(s shared.Rego, input any, ctx *shared.Context) shared.Rego {
		return NotContains{
			BaseCondition: BaseCondition{Subject: s},
			Value:         ResolveParameterValue[string](input, ctx),
		}
	},
	"in": func(s shared.Rego, input any, ctx *shared.Context) shared.Rego {
		return In{
			BaseCondition: BaseCondition{Subject: s},
			Values: func() []string {
				var v []string
				for _, i := range ResolveParameterValue[[]any](input, ctx) {
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
				for _, i := range ResolveParameterValue[[]any](input, ctx) {
					v = append(v, i.(string))
				}
				return v
			}(),
		}
	},
	"containskey": func(s shared.Rego, input any, ctx *shared.Context) shared.Rego {
		return ContainsKey{
			BaseCondition: BaseCondition{Subject: s},
			KeyName:       ResolveParameterValue[string](input, ctx),
		}
	},
	"notcontainskey": func(s shared.Rego, input any, ctx *shared.Context) shared.Rego {
		return NotContainsKey{
			BaseCondition: BaseCondition{Subject: s},
			KeyName:       ResolveParameterValue[string](input, ctx),
		}
	},
	"less": func(s shared.Rego, input any, ctx *shared.Context) shared.Rego {
		return Less{
			BaseCondition: BaseCondition{Subject: s},
			Value:         ResolveParameterValue[any](input, ctx),
		}
	},
	"lessorequals": func(s shared.Rego, input any, ctx *shared.Context) shared.Rego {
		return LessOrEquals{
			BaseCondition: BaseCondition{Subject: s},
			Value:         ResolveParameterValue[any](input, ctx),
		}
	},
	"greater": func(s shared.Rego, input any, ctx *shared.Context) shared.Rego {
		return Greater{
			BaseCondition: BaseCondition{Subject: s},
			Value:         ResolveParameterValue[any](input, ctx),
		}
	},
	"greaterorequals": func(s shared.Rego, input any, ctx *shared.Context) shared.Rego {
		return GreaterOrEquals{
			BaseCondition: BaseCondition{Subject: s},
			Value:         ResolveParameterValue[any](input, ctx),
		}
	},
	"exists": func(s shared.Rego, input any, ctx *shared.Context) shared.Rego {
		return Exists{
			BaseCondition: BaseCondition{Subject: s},
			Value:         ResolveParameterValue[any](input, ctx),
		}
	},
}

var paramRegex = regexp.MustCompile(`\[parameters\('([^']+)'\)\]`)

// ResolveParameterValue checks if the input is a parameter reference string.
// If it is, it extracts the parameter name and retrieves the value using ctx.GetParameterFunc.
// Otherwise, it returns the original input.
func ResolveParameterValue[T any](input any, ctx *shared.Context) T {
	// If input is not a string, return it as-is
	str, ok := input.(string)
	if !ok {
		return input.(T)
	}

	// Check if string matches the parameters reference format
	if matches := paramRegex.FindStringSubmatch(str); len(matches) > 1 {
		// Extract parameter name from the first capture group
		paramName := matches[1]

		// Call GetParameterFunc to get the value if available
		if ctx.GetParameterFunc != nil {
			if value, ok := ctx.GetParameterFunc(paramName); ok {
				return value.(T)
			}
		}
	}

	// Return original input if not a parameter reference or parameter not found
	return input.(T)
}
