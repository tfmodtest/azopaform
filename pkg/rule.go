package pkg

import (
	"fmt"
	"github.com/open-policy-agent/opa/format"
	"github.com/spf13/afero"
	"json-rule-finder/pkg/shared"
	"path/filepath"
	"strings"
)

var _ shared.Rego = &Rule{}

type Rule struct {
	Properties *PolicyRuleModel
	Id         string
	Name       string
	path       string
	result     string
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

%s`, ctx.PackageName(), rego, ctx.HelperFunctionsRego()), nil
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
	return nil
}

func (r *Rule) SaveToDisk() error {
	fileName := strings.TrimSuffix(filepath.Base(r.path), filepath.Ext(r.path)) + ".rego"
	err := afero.WriteFile(Fs, fileName, []byte(r.result), 0644)
	if err != nil {
		return fmt.Errorf("cannot save file %s, error is %+v", fileName, err)
	}
	return nil
}
