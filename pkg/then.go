package pkg

import (
	"fmt"
	"strings"

	"github.com/tfmodtest/azopaform/pkg/shared"
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
	if defaultEffect == shared.DeployIfNotExists {
		return shared.DeployIfNotExists, nil
	}

	return "", fmt.Errorf("unexpected input, effect is %s, defaultEffect is %s", effect, defaultEffect)
}

func (t *ThenBody) Action(ruleName, result, helperFunctionName string, rule *Rule) (string, error) {
	action, err := t.MapEffectToAction(rule.Properties.Parameters.GetEffect().GetDefaultValue())
	if err != nil {
		fmt.Printf("cannot map effect to action %+v\n", err)
		return "", err
	}
	var collection string
	var prefix string
	switch action {
	case shared.Deny:
		fallthrough
	case shared.Disabled:
		collection = shared.Deny
	case shared.Warn:
		collection = shared.Warn
	case shared.DeployIfNotExists:
		{
			collection = shared.Deny
			prefix = "not "
		}
	}
	if ruleName != "" {
		collection = collection + "_" + ruleName
	}
	if helperFunctionName != "" {
		return fmt.Sprintf(`%s if {
  res := resource(input, "azapi_resource")[_]
 %s%s(res)
}
%s`, collection, prefix, helperFunctionName, result), nil // collection + " if {\n " + helperFunctionName + "\n}\n" + result, nil
	}
	return fmt.Sprintf(`%s if {
 %s%s
}
`, collection, prefix, result), nil //collection + " if {\n " + result + "\n}\n", nil
}
