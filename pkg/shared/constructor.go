package shared

import (
	"fmt"
	"strings"
)

func FieldNameProcessor(fieldName string, ctx *Context) (string, error) {
	if strings.Contains(fieldName, VarInCountWhere) {
		return strings.ReplaceAll(fieldName, VarInCountWhere, "x"), nil
	}
	if fieldName == TypeOfResource || fieldName == KindOfResource {
		return fmt.Sprintf("%s.%s", ResourcePathPrefix, fieldName), nil
	}
	resourceType, err := currentResourceType(ctx)
	// No resource type defined, return the field name as is
	if err != nil {
		return processedFieldName(fieldName)
	}
	if strings.HasPrefix(fieldName, resourceType) {
		fieldName = strings.TrimPrefix(fieldName, resourceType)
	}
	currentVarName, ok := ctx.VarNameForField()
	if !ok {
		currentVarName = "r.values.body.properties"
	}
	if fieldName == "" {
		return currentVarName, nil
	}
	return currentVarName + "." + FieldNameParser(fieldName, ctx), nil
}

func processedFieldName(name string) (string, error) {
	name = fixUnescapedSingleQuotes(name)
	if !strings.Contains(name, "/") {
		return name, nil
	}
	split := strings.Split(name, "/")
	propertyPath := split[len(split)-1]
	return fmt.Sprintf("%s.properties.%s", ResourcePathPrefix, propertyPath), nil
}

// with double quotes to ensure they're valid for Rego language
func fixUnescapedSingleQuotes(content string) string {
	// Step 1: Temporarily mark all escaped single quotes with a unique marker
	tempMarker := "##ESCAPED_SINGLE_QUOTE##"
	intermediate := strings.ReplaceAll(content, `\'`, tempMarker)

	// Step 2: Replace all remaining single quotes with double quotes
	intermediate = strings.ReplaceAll(intermediate, `'`, `"`)

	// Step 3: Restore the escaped single quotes
	return strings.ReplaceAll(intermediate, tempMarker, `\'`)
}

func SliceConstructor(input any) string {
	var array []string
	var res string
	switch typedInput := input.(type) {
	case []any:
		for _, v := range typedInput {
			array = append(array, "\""+fmt.Sprint(v)+"\"")
		}
	case []string:
		for _, v := range typedInput {
			array = append(array, "\""+fmt.Sprint(v)+"\"")
		}
	case string:
		array = append(array, fmt.Sprint(typedInput))
	}

	res = strings.Join(array, ",")
	res = strings.Join([]string{"[", res, "]"}, "")
	return res
}

func currentResourceType(ctx *Context) (string, error) {
	resourceType, ok := ctx.currentResourceType()
	if !ok {
		return "", fmt.Errorf("cannot find the resource type in the context")
	}
	return resourceType, nil
}

func FieldNameParser(fieldNameRaw string, ctx *Context) string {
	if fieldNameRaw == TypeOfResource {
		return fieldNameRaw
	}

	return ConvertAzurePathToObjectPath(fieldNameRaw, ctx)
}

func ConvertAzurePathToObjectPath(prop string, ctx *Context) string {
	resourceType, _ := ctx.currentResourceType()
	prop = strings.TrimPrefix(prop, strings.ReplaceAll(resourceType, "/", "."))
	prop = strings.Replace(prop, ".", "/", -1)
	prop = strings.Replace(prop, "[x]", "/*", -1)
	prop = strings.Replace(prop, "[*]", "/*", -1)
	prop = strings.TrimPrefix(prop, "/")
	prop = strings.ReplaceAll(prop, "/", ".")
	return prop
}
