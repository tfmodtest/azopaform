package pkg

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

func (p *PolicyRuleParameters) GetParameter(name string) (any, bool) {
	if p == nil {
		return nil, false
	}
	if p.Parameters == nil {
		return nil, false
	}
	parameter, ok := p.Parameters[name]
	if !ok || parameter == nil {
		return nil, false
	}
	return parameter.DefaultValue, true
}
