package pkg

import (
	"errors"
)

type PolicyRuleParameterType string

type PolicyRuleParameterMetaData struct {
	Description string
	DisplayName string
	Deprecated  bool
}

type PolicyRuleParameter struct {
	Name         string
	Type         PolicyRuleParameterType
	DefaultValue any
	MetaData     *PolicyRuleParameterMetaData
}

type PolicyRuleParameters struct {
	Effect     *EffectBody
	Parameters map[string]*PolicyRuleParameter
}

func (p *PolicyRuleParameters) GetEffect() *EffectBody {
	if p == nil {
		return nil
	}
	return p.Effect
}

func (p *PolicyRuleParameters) GetParameter(name string) (any, bool, error) {
	if p == nil {
		return nil, false, nil
	}
	if p.Parameters == nil {
		return nil, false, nil
	}
	parameter, ok := p.Parameters[name]
	if !ok || parameter == nil {
		return nil, false, nil
	}
	if parameter.DefaultValue == nil {
		return nil, false, errors.New("only support parameter with default value now")
	}
	value := parameter.DefaultValue
	return value, true, nil
}
