package shared

import (
	"fmt"
	"testing"

	"github.com/open-policy-agent/opa/v1/rego"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func AssertRego(t *testing.T, query, cfg string, input map[string]any, expected bool, ctx *Context) {
	var allowed bool
	r := EvaluateRego(t, query, cfg, input, ctx)
	if r != nil {
		allowed = r.(bool)
	}
	assert.Equal(t, expected, allowed)
}

func EvaluateRego(t *testing.T, query, cfg string, input map[string]any, ctx *Context) any {
	eval, err := rego.New(rego.Query(query), rego.Module("test.rego", cfg)).PrepareForEval(ctx)
	if input != nil {
		eval, err = rego.New(rego.Query(query), rego.Module("test.rego", cfg), rego.Input(input)).PrepareForEval(ctx)
	}
	require.NoError(t, err)
	result, err := eval.Eval(ctx)
	require.NoError(t, err)
	if result == nil {
		return nil
	}
	return result[0].Expressions[0].Value
}

func AssertRegoAllow(t *testing.T, cfg string, input map[string]any, allowed bool, ctx *Context) {
	AssertRego(t, "data.main.allow", cfg, input, allowed, ctx)
}

func WithUtilFunctions(exp string) string {
	return fmt.Sprintf(RegoTestTemplate, exp) + "\n" + UtilsRego
}

const RegoTestTemplate = `package main

import rego.v1

default allow := false
allow if { 
  %s
}`
