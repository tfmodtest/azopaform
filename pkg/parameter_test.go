package pkg

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PolicyRuleParameters_GetParameter_NonExistent(t *testing.T) {
	// Test with non-existent parameter
	p := &PolicyRuleParameters{
		Parameters: map[string]*PolicyRuleParameter{
			"existingParam": {
				Name:         "existingParam",
				Type:         "String",
				DefaultValue: "value",
			},
		},
	}
	value, ok, err := p.GetParameter("nonExistentParam")
	require.NoError(t, err)
	assert.False(t, ok)
	assert.Nil(t, value)
}

func Test_PolicyRuleParameters_GetParameter_NilParameter(t *testing.T) {
	// Test with nil parameter
	p := &PolicyRuleParameters{
		Parameters: map[string]*PolicyRuleParameter{
			"nilParam": nil,
		},
	}
	value, ok, err := p.GetParameter("nilParam")
	require.NoError(t, err)
	assert.False(t, ok)
	assert.Nil(t, value)
}

func Test_PolicyRuleParameters_GetParameter_Success(t *testing.T) {
	// Test with valid parameter
	p := &PolicyRuleParameters{
		Parameters: map[string]*PolicyRuleParameter{
			"stringParam": {
				Name:         "stringParam",
				Type:         "String",
				DefaultValue: "testValue",
			},
			"arrayParam": {
				Name:         "arrayParam",
				Type:         "Array",
				DefaultValue: []any{"rapid", "stable", "patch"},
			},
		},
	}

	// Test string parameter
	value, ok, err := p.GetParameter("stringParam")
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, `"testValue"`, value)

	// Test array parameter
	value, ok, err = p.GetParameter("arrayParam")
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, []any{"rapid", "stable", "patch"}, value)
}
