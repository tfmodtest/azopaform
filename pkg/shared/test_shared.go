package shared

import (
	"encoding/json"
	"github.com/open-policy-agent/opa/rego"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func AssertRegoAllow(t *testing.T, cfg string, input *rego.EvalOption, allowed bool, ctx *Context) {
	AssertRego(t, "data.main.allow", cfg, input, allowed, ctx)
}

func AssertRego(t *testing.T, query, cfg string, input *rego.EvalOption, expected bool, ctx *Context) {
	eval, err := rego.New(rego.Query(query), rego.Module("test.rego", cfg)).PrepareForEval(ctx)
	require.NoError(t, err)
	var resultSet rego.ResultSet
	if input == nil {
		resultSet, err = eval.Eval(ctx)
	} else {
		j, _ := json.Marshal(*input)
		println(string(j))
		resultSet, err = eval.Eval(ctx, *input)
	}
	require.NoError(t, err)
	valueAsBool := len(resultSet) > 0 && resultSet[0].Expressions[0].Value != false
	assert.Equal(t, expected, valueAsBool)
}

const RegoTestTemplate = `package main

import rego.v1

default allow := false
allow if %s`
