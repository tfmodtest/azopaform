package shared

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ResolveParameterValue_NonStringInput(t *testing.T) {
	// Arrange
	ctx := NewContext()
	input := 42

	// Act
	result := ResolveParameterValue[int](input, ctx)

	// Assert
	assert.Equal(t, 42, result)
}

func Test_ResolveParameterValue_StringButNotParameter(t *testing.T) {
	// Arrange
	ctx := NewContext()
	input := "regular string"

	// Act
	result := ResolveParameterValue[string](input, ctx)

	// Assert
	assert.Equal(t, "regular string", result)
}

func Test_ResolveParameterValue_StringWithParameterFormat(t *testing.T) {
	// Arrange
	ctx := NewContext()
	ctx.GetParameterFunc = func(name string) (any, bool) {
		if name == "testParam" {
			return "param-value", true
		}
		return nil, false
	}
	input := "[parameters('testParam')]"

	// Act
	result := ResolveParameterValue[string](input, ctx)

	// Assert - this will fail with the current implementation
	assert.Equal(t, "param-value", result)
}

func Test_ResolveParameterValue_NonExistentParameter(t *testing.T) {
	// Arrange
	ctx := NewContext()
	ctx.GetParameterFunc = func(name string) (any, bool) {
		return nil, false
	}
	input := "[parameters('nonExistentParam')]"

	// Act
	result := ResolveParameterValue[string](input, ctx)

	// Assert
	assert.Equal(t, input, result)
}

func Test_ResolveParameterValue_ArrayParameter(t *testing.T) {
	// Arrange
	ctx := NewContext()
	expectedValue := []string{"rapid", "stable", "patch"}
	ctx.GetParameterFunc = func(name string) (any, bool) {
		if name == "allowedClusterAutoUpgradeChannels" {
			return expectedValue, true
		}
		return nil, false
	}
	input := "[parameters('allowedClusterAutoUpgradeChannels')]"

	// Act
	result := ResolveParameterValue[[]string](input, ctx)

	// Assert
	assert.Equal(t, expectedValue, result)
}

func Test_ResolveParameterValue_ComplexExpression(t *testing.T) {
	// Arrange
	ctx := NewContext()
	input := "some text with [parameters('paramName')] in the middle"

	// Act
	result := ResolveParameterValue[string](input, ctx)

	// Assert
	assert.Equal(t, input, result)
}
