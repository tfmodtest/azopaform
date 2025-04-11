package shared

import (
	"github.com/open-policy-agent/opa/rego"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func AssertRegoAllow(t *testing.T, cfg string, input map[string]any, allowed bool, ctx *Context) {
	eval, err := rego.New(rego.Query("data.main.allow"), rego.Module("test.rego", cfg)).PrepareForEval(ctx)
	if input != nil {
		eval, err = rego.New(rego.Query("data.main.allow"), rego.Module("test.rego", cfg), rego.Input(input)).PrepareForEval(ctx)
	}
	require.NoError(t, err)
	result, err := eval.Eval(ctx)
	require.NoError(t, err)
	assert.Equal(t, allowed, result.Allowed())
}

const RegoTestTemplate = `package main

import rego.v1

default allow := false
allow if { 
  %s
}`
