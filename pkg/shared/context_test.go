package shared

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ResolveParameterValue_NonStringInput(t *testing.T) {
	// Arrange
	ctx := NewContext()
	input := 42

	// Act
	result, err := ResolveParameterValue[int](input, ctx)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, 42, result)
}

func Test_ResolveParameterValue_StringButNotParameter(t *testing.T) {
	// Arrange
	ctx := NewContext()
	input := "regular string"

	// Act
	result, err := ResolveParameterValue[string](input, ctx)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, "regular string", result)
}

func Test_ResolveParameterValue_StringWithParameterFormat(t *testing.T) {
	// Arrange
	ctx := NewContext()
	ctx.GetParameterFunc = func(name string) (any, bool, error) {
		if name == "testParam" {
			return "param-value", true, nil
		}
		return nil, false, nil
	}
	input := "[parameters('testParam')]"

	// Act
	result, err := ResolveParameterValue[string](input, ctx)
	require.NoError(t, err)

	// Assert - this will fail with the current implementation
	assert.Equal(t, "param-value", result)
}

func Test_ResolveParameterValue_NonExistentParameter(t *testing.T) {
	// Arrange
	ctx := NewContext()
	ctx.GetParameterFunc = func(name string) (any, bool, error) {
		return nil, false, nil
	}
	input := "[parameters('nonExistentParam')]"

	// Act
	_, err := ResolveParameterValue[string](input, ctx)
	assert.Equal(t, "parameter nonExistentParam not found", err.Error())
}

func Test_ResolveParameterValue_ArrayParameter(t *testing.T) {
	// Arrange
	ctx := NewContext()
	expectedValue := []string{"rapid", "stable", "patch"}
	ctx.GetParameterFunc = func(name string) (any, bool, error) {
		if name == "allowedClusterAutoUpgradeChannels" {
			return expectedValue, true, nil
		}
		return nil, false, nil
	}
	input := "[parameters('allowedClusterAutoUpgradeChannels')]"

	// Act
	result, err := ResolveParameterValue[[]string](input, ctx)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, expectedValue, result)
}

func Test_ResolveParameterValue_ComplexExpression(t *testing.T) {
	// Arrange
	ctx := NewContext()
	input := "some text with [parameters('paramName')] in the middle"

	// Act
	result, err := ResolveParameterValue[string](input, ctx)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, input, result)
}
