package operation

import (
	"fmt"

	"github.com/tfmodtest/azopaform/pkg/shared"
)

var _ Operation = &Not{}

type Not struct {
	baseOperation
	Body shared.Rego
}

func NewNot(conditionSetName string, body shared.Rego) Not {
	return Not{
		Body: body,
		baseOperation: baseOperation{
			helperFunctionName: conditionSetName,
		},
	}
}

func parseNot(input any, ctx *shared.Context) (shared.Rego, error) {
	itemMap := input.(map[string]any)
	body, err := NewOperationOrCondition(itemMap, ctx)
	if err != nil {
		return nil, err
	}
	return Not{
		Body:          body,
		baseOperation: newBaseOperation(),
	}, nil
}

func (n Not) Rego(ctx *shared.Context) (string, error) {
	body, ok := n.Body.(Operation)
	if !ok {
		body = &AllOf{
			Conditions: []shared.Rego{
				n.Body,
			},
			baseOperation: baseOperation{
				helperFunctionName: fmt.Sprintf("%s_%s", n.HelperFunctionName(), "negation"),
			},
		}
	}
	
	bodyRes := body.HelperFunctionName()
	subSet, err := body.Rego(ctx)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`%s(r) if {
  not %s(r)
}

%s`, n.HelperFunctionName(), bodyRes, subSet), nil
}
