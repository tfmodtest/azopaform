package ctrl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/magodo/azure-rest-api-bridge/log"
	"github.com/magodo/azure-rest-api-bridge/mockserver"
	"github.com/magodo/azure-rest-api-bridge/mockserver/swagger"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/function/stdlib"
)

type Option struct {
	ConfigFile    string
	ContinueOnErr bool
	ServerOption  mockserver.Option
	ExecFrom      string
	ExecTo        string
}

type Ctrl struct {
	ExecSpec      Config
	ContinueOnErr bool
	MockServer    mockserver.Server

	ExecFrom  string
	ExecTo    string
	execState ExecutionState
}

type ExecutionState int

const (
	ExecutionStateBeforeRun ExecutionState = iota
	ExecutionStateRunning
	ExecutionStateAfterRun
)

func NewCtrl(opt Option) (*Ctrl, error) {
	parser := hclparse.NewParser()
	f, diags := parser.ParseHCLFile(opt.ConfigFile)
	if diags.HasErrors() {
		return nil, fmt.Errorf("parsing %s: %v", opt.ConfigFile, diags.Error())
	}

	homedir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("getting user home dir: %v", err)
	}
	var execSpec Config
	ctx := &hcl.EvalContext{
		Functions: map[string]function.Function{
			"jsonencode": stdlib.JSONEncodeFunc,
		},
		Variables: map[string]cty.Value{
			"home":        cty.StringVal(homedir),
			"server_addr": cty.StringVal(fmt.Sprintf("%s:%d", opt.ServerOption.Addr, opt.ServerOption.Port)),
		},
	}
	if diags := gohcl.DecodeBody(f.Body, ctx, &execSpec); diags.HasErrors() {
		return nil, fmt.Errorf("decoding %s: %v", opt.ConfigFile, diags.Error())
	}

	if err := validateExecSpec(execSpec); err != nil {
		return nil, fmt.Errorf("invalid exec spec: %v", err)
	}

	srv, err := mockserver.New(opt.ServerOption)
	if err != nil {
		return nil, fmt.Errorf("creating mock server: %v", err)
	}

	return &Ctrl{
		ExecSpec:      execSpec,
		ContinueOnErr: opt.ContinueOnErr,
		MockServer:    *srv,
		ExecFrom:      opt.ExecFrom,
		ExecTo:        opt.ExecTo,
		execState:     ExecutionStateBeforeRun,
	}, nil
}

func validateExecSpec(spec Config) error {
	validateOverride := func(ovs []Override) error {
		for _, ov := range ovs {
			if ov.ResponseBody+ov.ResponseSelectorMerge+ov.ResponseSelectorJSON+ov.ResponsePatchJSON+ov.ResponsePatchMerge == "" && len(ov.ResponseHeader) == 0 && ov.ExpanderOption == nil && ov.SynthOption == nil {
				return fmt.Errorf("empty override block is not allowed")
			}
			if ov.ResponseBody != "" {
				if ov.ResponseSelectorMerge+ov.ResponseSelectorJSON+ov.ResponsePatchJSON+ov.ResponsePatchMerge != "" || ov.ExpanderOption != nil || ov.SynthOption != nil {
					return fmt.Errorf("`response_body` can only be exclusive specified")
				}
				continue
			}
			if ov.ResponsePatchJSON != "" && ov.ResponsePatchMerge != "" {
				return fmt.Errorf("`response_patch_merge` conflicts with `response_patch_json`")
			}
			if ov.ResponseSelectorMerge != "" && ov.ResponseSelectorJSON != "" {
				return fmt.Errorf("`response_selector_merge` conflicts with `response_selector_json`")
			}
		}
		return nil
	}

	if err := validateOverride(spec.Overrides); err != nil {
		return err
	}

	execNames := map[string]map[string]bool{}
	for _, exec := range spec.Executions {
		if exec.Skip && exec.SkipReason == "" {
			return fmt.Errorf("skipped execution %s must have a skip_reason", exec)
		}
		if err := validateOverride(exec.Overrides); err != nil {
			return err
		}
		m, ok := execNames[exec.Name]
		if !ok {
			m = map[string]bool{}
			execNames[exec.Name] = m
		}
		if m[exec.Type] {
			return fmt.Errorf("duplicated execution %s", exec)
		}
		m[exec.Type] = true
	}

	return nil
}

func (ctrl *Ctrl) Run(ctx context.Context) error {
	// Start mock server
	log.Info("Starting the mock server")
	if err := ctrl.MockServer.Start(); err != nil {
		return err
	}

	results := map[string][]SingleModelMap{}

	execTotal := len(ctrl.ExecSpec.Executions)
	execSkip := 0
	execSucceed := 0
	execFail := 0

	expCache := swagger.NewExpanderCache()

	// Launch each execution
	for i, execution := range ctrl.ExecSpec.Executions {
		switch ctrl.execState {
		case ExecutionStateBeforeRun:
			if ctrl.ExecFrom == "" || ctrl.ExecFrom == execution.String() {
				ctrl.execState = ExecutionStateRunning
			} else {
				log.Info(fmt.Sprintf("Skipping %s (%d/%d): skipped by -from", execution, i+1, execTotal))
				execSkip++
				continue
			}
		case ExecutionStateRunning:
			if ctrl.ExecTo != "" && ctrl.ExecTo == execution.String() {
				ctrl.execState = ExecutionStateAfterRun
				log.Info(fmt.Sprintf("Skipping %s (%d/%d): skipped by -to", execution, i+1, execTotal))
				execSkip++
				continue
			}
		case ExecutionStateAfterRun:
			log.Info(fmt.Sprintf("Skipping %s (%d/%d): skipped by -to", execution, i+1, execTotal))
			execSkip++
			continue
		}

		if execution.Skip {
			log.Info(fmt.Sprintf("Skipping %s (%d/%d): %s", execution, i+1, execTotal, execution.SkipReason))
			execSkip++
			continue
		}

		run := func(execution Execution) error {

			overrides := append([]Override{}, execution.Overrides...)
			overrides = append(overrides, ctrl.ExecSpec.Overrides...)

			var ovs []mockserver.Override
			for _, override := range overrides {
				ov := mockserver.Override{
					PathPattern:           *regexp.MustCompile(override.PathPattern),
					ResponseSelectorMerge: override.ResponseSelectorMerge,
					ResponseSelectorJSON:  override.ResponseSelectorJSON,
					ResponseBody:          override.ResponseBody,
					ResponsePatchMerge:    override.ResponsePatchMerge,
					ResponsePatchJSON:     override.ResponsePatchJSON,
					ResponseHeader:        override.ResponseHeader,
					SynthOption:           &swagger.SynthesizerOption{},
					ExpanderOption: &swagger.ExpanderOption{
						Cache: expCache,
					},
				}
				if opt := override.SynthOption; opt != nil {
					if opt.UseEnumValue {
						ov.SynthOption.UseEnumValues = true
					}
					var del []swagger.SynthDuplicateElement
					for _, eopt := range opt.DuplicateElement {
						cnt := 1
						if eopt.Count != nil {
							cnt = *eopt.Count
						}
						del = append(del, swagger.SynthDuplicateElement{
							Cnt:  cnt,
							Addr: swagger.ParseAddr(eopt.Addr),
						})
					}
					ov.SynthOption.DuplicateElements = del
				}
				if opt := override.ExpanderOption; opt != nil {
					if opt.EmptyObjAsStr {
						ov.ExpanderOption.EmptyObjAsStr = true
					}
					if opt.DisableCache {
						ov.ExpanderOption.Cache = nil
					}
				}

				ovs = append(ovs, ov)
			}

			ctrl.MockServer.InitExecution(ovs)

			env := os.Environ()
			for k, v := range execution.Env {
				env = append(env, k+"="+v)
			}

			var stdout bytes.Buffer
			var stderr bytes.Buffer

			cmd := exec.Cmd{
				Path:   execution.Path,
				Args:   append([]string{filepath.Base(execution.Path)}, execution.Args...),
				Env:    env,
				Dir:    execution.Dir,
				Stdout: &stdout,
				Stderr: &stderr,
			}

			log.Info(fmt.Sprintf("Executing %s (%d/%d)", execution, i+1, execTotal))

			log.Debug("execution detail", "path", execution.Path, "args", execution.Args, "env", env, "dir", execution.Dir)

			if err := cmd.Run(); err != nil {
				log.Error("run failure", "stdout", stdout.String(), "stderr", stderr.String())
				return fmt.Errorf("running execution %q: %v", execution, err)
			}

			log.Debug("execution result", "stdout", stdout.String())

			var appModel interface{}
			if err := json.Unmarshal(stdout.Bytes(), &appModel); err != nil {
				log.Error("post-execution unmarshal failure", "error", err, "stdout", stdout.String())
				return fmt.Errorf("post-execution %q unmarshal: %v", execution, err)
			}

			m, err := MapSingleAppModel(appModel, ctrl.MockServer.Records()...)
			if err != nil {
				log.Error("post-execution map models", "error", err)
				return fmt.Errorf("post-execution %q map models: %v", execution, err)
			}

			if err := m.AddLink(ctrl.MockServer.Idx.Commit, ctrl.MockServer.Specdir); err != nil {
				log.Error("post-execution model map adding link", "error", err)
				return fmt.Errorf("post-execution model map adding link: %v", err)
			}
			if err := m.RelativeLocalLink(ctrl.MockServer.Specdir); err != nil {
				log.Error("post-execution model map relative local link", "error", err)
				return fmt.Errorf("post-execution model map relative local link: %v", err)
			}

			results[execution.Name] = append(results[execution.Name], m)

			return nil
		}

		if err := run(execution); err != nil {
			execFail++
			if ctrl.ContinueOnErr {
				continue
			}
			return err
		} else {
			execSucceed++
		}
	}

	if ctrl.ContinueOnErr {
		log.Info("Summary", "total", execTotal, "succeed", execSucceed, "fail", execFail, "skip", execSkip)
	}

	if err := ctrl.WriteResult(ctx, results); err != nil {
		log.Error("Write Result", "err", err.Error())
		return err
	}

	// Stop mock server
	log.Info("Stopping the mock server")
	if err := ctrl.MockServer.Stop(ctx); err != nil {
		return err
	}

	return nil
}

func (ctrl *Ctrl) WriteResult(ctx context.Context, results map[string][]SingleModelMap) error {
	outputs := map[string]ModelMap{}
	for execName, models := range results {
		outputs[execName] = NewModelMap(models)
	}

	b, err := json.MarshalIndent(outputs, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling output: %v", err)
	}

	fmt.Println(string(b))
	return nil
}
