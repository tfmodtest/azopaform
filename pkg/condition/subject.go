package condition

import (
	"fmt"
	"strings"

	"github.com/tfmodtest/azopaform/pkg/shared"
)

func subjectRego(subject shared.Rego, value any, callback func(shared.Rego, any, *shared.Context) (string, error), ctx *shared.Context) (string, error) {
	if field, ok := subject.(FieldValue); ok && strings.Contains(field.Name, "[*]") && !ctx.IsInCountRego() {
		return conditionInEvery(strings.ReplaceAll(field.Name, "/", "."), value, callback, ctx)
	}
	if field, ok := subject.(FieldValue); ok && ctx.IsInCountRego() {
		field.Name = strings.TrimPrefix(field.Name, ctx.CurrentCountFieldName())
		if strings.HasPrefix(field.Name, ".") {
			field.Name = field.Name[1:]
		}
		subject = field
	}
	return callback(subject, value, ctx)
}

func conditionInEvery(name string, value any, callback func(shared.Rego, any, *shared.Context) (string, error), ctx *shared.Context) (string, error) {
	// Find the first occurrence of [*]
	idx := strings.Index(name, "[*]")
	if idx == -1 {
		return callback(FieldValue{Name: name}, value, ctx)
	}

	// Split the path into array part and remaining part
	arrayPath := name[:idx]
	remaining := ""
	if idx+3 < len(name) {
		remaining = name[idx+3:] // Skip "[*]"
	}

	// Generate a variable name based on the array path
	parts := strings.Split(arrayPath, ".")
	currentLeft := "item"
	if len(parts) > 0 {
		lastPart := parts[len(parts)-1]
		if lastPart != "" {
			currentLeft = lastPart
		}
	}
	currentLeft += shared.HumanFriendlyString(4)
	currentRightPrefix, ok := ctx.VarNameForField()
	if !ok {
		currentRightPrefix = "r.values.properties"
	}
	collection := currentRightPrefix + "." + shared.ConvertAzurePathToObjectPath(arrayPath, ctx)

	var itemField shared.Rego
	ctx.PushVarNameForField(currentLeft)
	defer ctx.PopVarNameForField()
	// Clean up the remaining path if it starts with a dot
	if strings.HasPrefix(remaining, ".") {
		remaining = remaining[1:]
	}

	// If there are more [*] in the remaining path, handle them recursively
	if strings.Contains(remaining, "[*]") {

		// Recursively process the nested array
		nestedCondition, err := conditionInEvery(remaining, value, callback, ctx)
		if err != nil {
			return "", err
		}

		// Format the outer "every" block
		return fmt.Sprintf("every %s in %s {\n %s \n}", currentLeft, collection, nestedCondition), nil
	}

	if remaining == "" {
		itemField = FieldValue{Name: ""}
	} else {
		itemField = FieldValue{Name: remaining}
	}

	condition, err := callback(itemField, value, ctx)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("every %s in %s {\n %s \n}", currentLeft, collection, condition), nil
}
