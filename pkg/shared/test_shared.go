package shared

import (
	"github.com/open-policy-agent/opa/rego"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func AssertRego(t *testing.T, query, cfg string, input map[string]any, expected bool, ctx *Context) {
	eval, err := rego.New(rego.Query(query), rego.Module("test.rego", cfg)).PrepareForEval(ctx)
	if input != nil {
		eval, err = rego.New(rego.Query(query), rego.Module("test.rego", cfg), rego.Input(input)).PrepareForEval(ctx)
	}
	require.NoError(t, err)
	result, err := eval.Eval(ctx)
	require.NoError(t, err)
	assert.Equal(t, expected, result.Allowed())
}

func AssertRegoAllow(t *testing.T, cfg string, input map[string]any, allowed bool, ctx *Context) {
	AssertRego(t, "data.main.allow", cfg, input, allowed, ctx)
}

const RegoTestTemplate = `package main

import rego.v1

default allow := false
allow if { 
  %s
}`
