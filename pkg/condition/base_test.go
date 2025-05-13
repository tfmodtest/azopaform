package condition

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tfmodtest/azopaform/pkg/shared"
)

func Test_ResolveParameterValue_NonStringInput(t *testing.T) {
	// Arrange
	ctx := shared.NewContext()
	input := 42

	// Act
	result := ResolveParameterValue[int](input, ctx)

	// Assert
	assert.Equal(t, 42, result)
}

func Test_ResolveParameterValue_StringButNotParameter(t *testing.T) {
	// Arrange
	ctx := shared.NewContext()
	input := "regular string"

	// Act
	result := ResolveParameterValue[string](input, ctx)

	// Assert
	assert.Equal(t, "regular string", result)
}

func Test_ResolveParameterValue_StringWithParameterFormat(t *testing.T) {
	// Arrange
	ctx := shared.NewContext()
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
	ctx := shared.NewContext()
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
	ctx := shared.NewContext()
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
	ctx := shared.NewContext()
	input := "some text with [parameters('paramName')] in the middle"

	// Act
	result := ResolveParameterValue[string](input, ctx)

	// Assert
	assert.Equal(t, input, result)
}

func Test_NewCondition_WithParameterValue(t *testing.T) {
	// Arrange
	conditionType := "equals"
	subject := shared.StringRego("field.name")
	paramValue := "[parameters('testParam')]"
	expectedResolvedValue := "resolved-value"

	ctx := shared.NewContext()
	ctx.GetParameterFunc = func(name string) (any, bool) {
		if name == "testParam" {
			return expectedResolvedValue, true
		}
		return nil, false
	}

	// Act
	condition := NewCondition(conditionType, subject, paramValue, ctx)

	// Assert
	// Verify condition is not nil
	require.NotNil(t, condition)

	// Verify it's the right type of condition
	equalsCondition, ok := condition.(Equals)
	require.True(t, ok, "Expected Equals condition")

	// Verify the condition's value is the resolved parameter value
	assert.Equal(t, expectedResolvedValue, equalsCondition.Value)
}

func Test_NewCondition_WithArrayParameterValue(t *testing.T) {
	// Arrange
	conditionType := "in"
	subject := shared.StringRego("field.name")
	paramValue := "[parameters('allowedValues')]"
	expectedResolvedValue := []any{"value1", "value2", "value3"}

	ctx := shared.NewContext()
	ctx.GetParameterFunc = func(name string) (any, bool) {
		if name == "allowedValues" {
			return expectedResolvedValue, true
		}
		return nil, false
	}

	// Act
	condition := NewCondition(conditionType, subject, paramValue, ctx)

	// Assert
	// Verify condition is not nil
	require.NotNil(t, condition)

	// Verify it's the right type of condition
	inCondition, ok := condition.(In)
	require.True(t, ok, "Expected In condition")

	// Verify the condition's values match the resolved parameter values
	expectedStrings := []string{"value1", "value2", "value3"}
	assert.Equal(t, expectedStrings, inCondition.Values)
}

func Test_NewCondition_InvalidConditionType(t *testing.T) {
	// Arrange
	conditionType := "nonExistentCondition"
	subject := shared.StringRego("field.name")
	value := "test"
	ctx := shared.NewContext()

	// Act
	condition := NewCondition(conditionType, subject, value, ctx)

	// Assert
	assert.Nil(t, condition, "Should return nil for invalid condition type")
}
