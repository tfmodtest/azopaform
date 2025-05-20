package shared

import (
	"fmt"
	"strconv"
	"strings"
)

func FieldNameProcessor(fieldName string, ctx *Context) (string, error) {
	var result string
	if fieldName == TypeOfResource || fieldName == KindOfResource {
		return fmt.Sprintf("%s.%s", ResourcePathPrefix, fieldName), nil
	}
	rt, err := currentResourceType(ctx)
	if err != nil {
		return processedFieldName(fieldName)
	}
	res, err := FieldNameParser(fieldName, rt, "")
	if err != nil {
		return "", err
	}
	result = TFNameMapping(res)

	return result, nil
}

func processedFieldName(name string) (string, error) {
	name = fixUnescapedSingleQuotes(name)
	if !strings.Contains(name, "/") {
		return name, nil
	}
	split := strings.Split(name, "/")
	propertyPath := split[len(split)-1]
	propertyPath = strings.ReplaceAll(propertyPath, "[*]", "[_]")
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

func FieldNameParser(fieldNameRaw, resourceType, version string) (string, error) {
	if fieldNameRaw == TypeOfResource {
		return fieldNameRaw, nil
	}
	prop := strings.TrimPrefix(fieldNameRaw, resourceType)

	prop = strings.Replace(prop, ".", "/", -1)
	prop = strings.Replace(prop, "[x]", "/*", -1)
	prop = strings.Replace(prop, "[*]", "/*", -1)
	prop = strings.TrimPrefix(prop, "/")
	prop = strings.ReplaceAll(prop, "/*", "[_]")
	prop = strings.ReplaceAll(prop, "/", ".")
	prop = "properties." + prop
	return prop, nil
}

func TFNameMapping(fieldName string) string {
	var result string
	attributes := strings.Split(fieldName, "/")
	for _, v := range attributes {
		if v == "" {
			continue
		}
		next := result + "." + v
		if _, err := strconv.Atoi(v); err == nil || v == "*" {
			next = result + "[" + v + "]"
		}
		result = next
	}
	result = ResourcePathPrefix + result

	return result
}
