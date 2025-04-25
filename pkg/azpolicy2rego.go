package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/afero"
	"github.com/tfmodtest/azopaform/pkg/shared"
	"path/filepath"
)

var Fs = afero.NewOsFs()

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

func AzurePolicyToRego(policyPath string, dir string, ctx *shared.Context) error {
	var paths []string
	var err error

	//For batch translation
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
		rule, err := loadRule(path, ctx)
		if err != nil {
			return fmt.Errorf("error when loading rule from path %s, error is %+v", path, err)
		}
		err = rule.SaveToDisk()
		if err != nil {
			return fmt.Errorf("error when saving parsed rule to disk, error is %+v", err)
		}
	}
	if utilLibraryName := ctx.UtilLibraryPackageName(); utilLibraryName == "" {
		return saveUtilRegoFile(ctx)
	}
	return nil
}

func loadRule(path string, ctx *shared.Context) (*Rule, error) {
	rule, err := readRuleFromFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot find rules %+v", err)
	}
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

func readRuleFromFile(path string) (*Rule, error) {
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

func saveUtilRegoFile(ctx *shared.Context) error {
	err := afero.WriteFile(Fs, ctx.UtilRegoFileName(), []byte(fmt.Sprintf(`package %s

import rego.v1

%s`, ctx.PackageName(), shared.UtilsRego)), 0644)
	if err != nil {
		return fmt.Errorf("cannot save file utils.rego, error is %+v", err)
	}
	return nil
}
