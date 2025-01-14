package shared

import (
	"context"
	"github.com/open-policy-agent/opa/rego"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func AssertRegoAllow(t *testing.T, cfg string, input *rego.EvalOption, allowed bool, ctx context.Context) {
	eval, err := rego.New(rego.Query("data.main.allow"), rego.Module("test.rego", cfg)).PrepareForEval(ctx)
	require.NoError(t, err)
	var result rego.ResultSet
	if input == nil {
		result, err = eval.Eval(ctx)
	} else {
		result, err = eval.Eval(ctx, *input)
	}
	require.NoError(t, err)
	assert.Equal(t, allowed, result.Allowed())
}
