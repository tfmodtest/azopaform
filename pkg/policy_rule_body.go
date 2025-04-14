package pkg

import (
	"json-rule-finder/pkg/shared"
)

type PolicyRuleBody struct {
	Then *ThenBody
	If   IfBody `json:"if,omitempty"`
}

func (p *PolicyRuleBody) GetIf(ctx *shared.Context) (*If, error) {
	return NewIf(p.If, ctx)
}

func (p *PolicyRuleBody) GetThen() *ThenBody {
	if p == nil {
		return nil
	}
	return p.Then
}
