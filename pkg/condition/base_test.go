package condition

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tfmodtest/azopaform/pkg/shared"
)

func Test_NewCondition_WithParameterValue(t *testing.T) {
	// Arrange
	conditionType := "equals"
	subject := shared.StringRego("field.name")
	paramValue := "[parameters('testParam')]"
	expectedResolvedValue := "resolved-value"

	ctx := shared.NewContext()
	ctx.GetParameterFunc = func(name string) (any, bool, error) {
		if name == "testParam" {
			return expectedResolvedValue, true, nil
		}
		return nil, false, nil
	}

	// Act
	condition, err := NewCondition(conditionType, subject, paramValue, ctx)
	require.NoError(t, err)

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
	ctx.GetParameterFunc = func(name string) (any, bool, error) {
		if name == "allowedValues" {
			return expectedResolvedValue, true, nil
		}
		return nil, false, nil
	}

	// Act
	condition, err := NewCondition(conditionType, subject, paramValue, ctx)
	require.NoError(t, err)

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
	condition, err := NewCondition(conditionType, subject, value, ctx)
	require.NoError(t, err)

	// Assert
	assert.Nil(t, condition, "Should return nil for invalid condition type")
}
