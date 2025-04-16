package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/open-policy-agent/opa/format"
	"json-rule-finder/pkg/shared"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

var _ shared.Rego = &Rule{}

type Rule struct {
	Properties  *PolicyRuleModel
	Id          string
	Name        string
	path        string
	result      string
	packageName string
}

func (r *Rule) Rego(ctx *shared.Context) (string, error) {
	ifBody, err := r.Properties.PolicyRule.GetIf(ctx)
	if err != nil {
		return "", err
	}
	ifRego, err := ifBody.Rego(ctx)
	if err != nil {
		return "", err
	}
	then := r.Properties.PolicyRule.GetThen()
	conditionName := ifBody.ConditionName(ifRego)
	rego, err := then.Action(ifRego, conditionName, r)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`package %s

import rego.v1

%s

%s`, r.PackageName(), rego, ctx.HelperFunctionsRego()), nil
}

func (r *Rule) PackageName() string {
	return getOrDefault(r.packageName, "main")
}

func (r *Rule) Parse(ctx *shared.Context) error {
	ruleRego, err := r.Rego(ctx)
	if err != nil {
		return err
	}
	formattedSrc, err := format.Source("output.rego", []byte(ruleRego))
	if err != nil {
		return fmt.Errorf("invalid rego code: %w", err)
	}
	r.result = string(formattedSrc)
	//r.result = ruleRego
	return nil
}

func (r *Rule) SaveToDisk() error {
	fileName := strings.TrimSuffix(filepath.Base(r.path), filepath.Ext(r.path)) + ".rego"
	err := afero.WriteFile(Fs, fileName, []byte(r.result), 0644)
	if err != nil {
		return fmt.Errorf("cannot save file %s, error is %+v", fileName, err)
	}
	err = afero.WriteFile(Fs, "utils.rego", []byte(fmt.Sprintf(`package %s

import rego.v1

%s`, r.PackageName(), shared.UTILS_REGO)), 0644)
	if err != nil {
		return fmt.Errorf("cannot save file utils.rego, error is %+v", err)
	}
	return nil
}

func NewPolicyRuleBody(input map[string]any) *PolicyRuleBody {
	return &PolicyRuleBody{
		If: input,
	}
}

type PolicyRuleMetaData struct {
	Version  string
	Category string
}

type PolicyRuleModel struct {
	PolicyRule  *PolicyRuleBody
	Parameters  *PolicyRuleParameters
	DisplayName string
	PolicyType  string
	Mode        string
	Description string
	Version     string
	Metadata    *PolicyRuleMetaData
}

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
	if defaultEffect == shared.Deny || defaultEffect == shared.Disabled {
		return defaultEffect, nil
	}
	return "", fmt.Errorf("unexpected input, effect is %s, defaultEffect is %s", effect, defaultEffect)
}

func (t *ThenBody) Action(result, conditionName string, rule *Rule) (string, error) {
	action, err := t.MapEffectToAction(rule.Properties.Parameters.GetEffect().GetDefaultValue())
	if err != nil {
		fmt.Printf("cannot map effect to action %+v\n", err)
		return "", err
	}
	if action == shared.Disabled {
		result = "default allow := true\n\n" + result
	} else if action == shared.Deny {
		top := "deny if {\n" + " " + conditionName + "\n}\n"
		result = top + result
	} else if action == shared.Warn {
		top := "warn if {\n" + " " + conditionName + "\n}\n"
		result = top + result
	}
	return result, nil
}

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

type EffectBody struct {
	DefaultValue string `json:"defaultValue"`
}

func (e *EffectBody) GetDefaultValue() string {
	if e == nil {
		return ""
	}
	return e.DefaultValue
}

type OperatorModel struct {
	Name  string
	Value any
}

var Fs = afero.NewOsFs()

func AzurePolicyToRego(policyPath string, dir string, options Options, ctx *shared.Context) error {
	var paths []string
	var err error

	////For batch translation
	if dir != "" {
		paths, err = jsonFiles(dir)
		if err != nil {
			return err
		}
	}

	//Override for hard cases
	if policyPath != "" {
		paths = []string{policyPath}
	}
	for _, path := range paths {
		rule, err := LoadRule(path, options, ctx)
		if err != nil {
			return fmt.Errorf("error when loading rule from path %s, error is %+v", path, err)
		}
		err = rule.SaveToDisk()
		if err != nil {
			return fmt.Errorf("error when saving parsed rule to disk, error is %+v", err)
		}
	}
	return nil
}

func LoadRule(path string, option Options, ctx *shared.Context) (*Rule, error) {
	rule, err := ReadRuleFromFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot find rules %+v", err)
	}
	rule.packageName = option.PackageName
	err = rule.Parse(ctx)
	return rule, err
}

func jsonFiles(dir string) ([]string, error) {
	res, err := readJsonFilePaths(dir)
	if err != nil {
		fmt.Printf("cannot find files in directory %+v\n", err)
		return nil, err
	}
	return res, nil
}

func readJsonFilePaths(path string) ([]string, error) {
	var filePaths []string

	entries, err := afero.ReadDir(Fs, path)
	if err != nil {
		return nil, fmt.Errorf("cannot find files in directory %+v\n", err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			subFilePaths, err := readJsonFilePaths(path + "/" + entry.Name())
			if err != nil {
				return nil, err
			}
			filePaths = append(filePaths, subFilePaths...)
			continue
		}
		if filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		fileName := entry.Name()
		filePath := path + "/" + fileName
		filePaths = append(filePaths, filePath)
	}
	return filePaths, nil
}

func ReadRuleFromFile(path string) (*Rule, error) {
	data, err := afero.ReadFile(Fs, path)
	if err != nil {
		return nil, err
	}

	var rule Rule
	err = json.Unmarshal(data, &rule)
	if err != nil {
		fmt.Printf("unable to unmarshal json %+v\n\n", err)
		return nil, err
	}
	m := make(map[string]any)
	err = json.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}
	rule.path = path

	return &rule, nil
}

func getOrDefault[T comparable](value, defaultValue T) T {
	var defaultTValue T
	if value == defaultTValue {
		return defaultValue
	}
	return value
}
