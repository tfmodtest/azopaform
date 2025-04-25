package pkg

import (
	"fmt"
	"github.com/tfmodtest/azopaform/pkg/shared"
	"strings"
)

type ThenBody struct {
	Effect string `json:"effect,omitempty"`
}

func (t *ThenBody) GetEffect() string {
	if t == nil {
		return ""
	}
	return t.Effect
}

func (t *ThenBody) MapEffectToAction(defaultEffect string) (string, error) {
	effect := t.GetEffect()
	if effect == "" {
		return "", fmt.Errorf("unexpected input, effect is %s, defaultEffect is %s", effect, defaultEffect)
	}
	effect = strings.ToLower(effect)
	if effect == shared.Deny {
		return shared.Deny, nil
	}
	if effect != "[parameters('effect')]" {
		return "", fmt.Errorf("unexpected input, effect is %s, defaultEffect is %s", effect, defaultEffect)
	}
	defaultEffect = strings.ToLower(defaultEffect)
	if defaultEffect == shared.Audit {
		return shared.Warn, nil
	}
	if defaultEffect == shared.Modify || defaultEffect == shared.Deny || defaultEffect == shared.Disabled {
		return shared.Deny, nil
	}
	return "", fmt.Errorf("unexpected input, effect is %s, defaultEffect is %s", effect, defaultEffect)
}

func (t *ThenBody) Action(result, conditionName string, rule *Rule) (string, error) {
	action, err := t.MapEffectToAction(rule.Properties.Parameters.GetEffect().GetDefaultValue())
	if err != nil {
		fmt.Printf("cannot map effect to action %+v\n", err)
		return "", err
	}
	var top string
	switch action {
	case shared.Deny:
		fallthrough
	case shared.Disabled:
		top = "deny if {\n" + " " + conditionName + "\n}\n"
	case shared.Warn:
		top = "warn if {\n" + " " + conditionName + "\n}\n"
	}
	return top + result, nil
}
