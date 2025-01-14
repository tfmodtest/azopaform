package pkg

import (
	"context"
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

type Condition interface {
	shared.Rego
}

type BaseCondition struct {
	Subject shared.Rego
}

var ConditionFactory = map[string]func(shared.Rego, any) shared.Rego{
	"equals": func(s shared.Rego, input any) shared.Rego {
		return EqualsCondition{
			BaseCondition: BaseCondition{Subject: s},
			Value:         input,
		}
	},
	"notequals": func(s shared.Rego, input any) shared.Rego {
		return NotEqualsCondition{
			BaseCondition: BaseCondition{Subject: s},
			Value:         input,
		}
	},
	"like": func(s shared.Rego, input any) shared.Rego {
		return LikeCondition{
			BaseCondition: BaseCondition{Subject: s},
			Value:         input.(string),
		}
	},
	"notlike": func(s shared.Rego, input any) shared.Rego {
		return NotLikeCondition{
			BaseCondition: BaseCondition{Subject: s},
			Value:         input.(string),
		}
	},
	"match": func(s shared.Rego, input any) shared.Rego {
		return MatchCondition{
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
		return NotMatchCondition{
			BaseCondition: BaseCondition{Subject: s},
			Value:         input.(string),
		}
	},
	"notmatchinsensitively": func(s shared.Rego, input any) shared.Rego {
		return NotMatchInsensitivelyCondition{
			BaseCondition: BaseCondition{Subject: s},
			Value:         input.(string),
		}
	},
	"contains": func(s shared.Rego, input any) shared.Rego {
		return ContainsCondition{
			BaseCondition: BaseCondition{Subject: s},
			Value:         input.(string),
		}
	},
	"notcontains": func(s shared.Rego, input any) shared.Rego {
		return NotContainsCondition{
			BaseCondition: BaseCondition{Subject: s},
			Value:         input.(string),
		}
	},
	"in": func(s shared.Rego, input any) shared.Rego {
		return InCondition{
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
		return NotInCondition{
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
		return ContainsKeyCondition{
			BaseCondition: BaseCondition{Subject: s},
			KeyName:       input.(string),
		}
	},
	"notcontainskey": func(s shared.Rego, input any) shared.Rego {
		return NotContainsKeyCondition{
			BaseCondition: BaseCondition{Subject: s},
			KeyName:       input.(string),
		}
	},
	"less": func(s shared.Rego, input any) shared.Rego {
		return LessCondition{
			BaseCondition: BaseCondition{Subject: s},
			Value:         input,
		}
	},
	"lessorequals": func(s shared.Rego, input any) shared.Rego {
		return LessOrEqualsCondition{
			BaseCondition: BaseCondition{Subject: s},
			Value:         input,
		}
	},
	"greater": func(s shared.Rego, input any) shared.Rego {
		return GreaterCondition{
			BaseCondition: BaseCondition{Subject: s},
			Value:         input,
		}
	},
	"greaterorequals": func(s shared.Rego, input any) shared.Rego {
		return GreaterOrEqualsCondition{
			BaseCondition: BaseCondition{Subject: s},
			Value:         input,
		}
	},
	"exists": func(s shared.Rego, input any) shared.Rego {
		return ExistsCondition{
			BaseCondition: BaseCondition{Subject: s},
			Value:         input,
		}
	},
}
