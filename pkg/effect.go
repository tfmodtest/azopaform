package pkg

type EffectBody struct {
	DefaultValue string `json:"defaultValue"`
}

func (e *EffectBody) GetDefaultValue() string {
	if e == nil {
		return ""
	}
	return e.DefaultValue
}
