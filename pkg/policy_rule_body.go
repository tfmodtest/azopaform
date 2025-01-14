package pkg

import (
	"json-rule-finder/pkg/shared"
)

type PolicyRuleBody struct {
	Then   *ThenBody
	If     IfBody `json:"if,omitempty"`
	IfBody shared.Rego
}

func (p *PolicyRuleBody) GetIf() *If {
	return &If{
		body: p.If,
	}
}

func (p *PolicyRuleBody) GetThen() *ThenBody {
	if p == nil {
		return nil
	}
	return p.Then
}

func (p *PolicyRuleBody) BuildIfBody(ctx *shared.Context) *PolicyRuleBody {
	if p == nil {
		return nil
	}
	return NewPolicyRuleBody(p.If, ctx)
}

func (p *PolicyRuleBody) GetIfBody() *If {
	return &If{
		body: p.If,
	}
}
