package shared

import (
	"fmt"
	"strconv"
	"strings"
)

func FieldNameProcessor(fieldName string, ctx *Context) (string, error) {
	var result string
	if fieldName == TypeOfResource || fieldName == KindOfResource {
		return fmt.Sprintf("r.values.%s", fieldName), nil
	}
	if strings.Contains(fieldName, "count") {
		return fieldName, nil
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
	if !strings.Contains(name, "/") {
		return name, nil
	}
	split := strings.Split(name, "/")
	propertyPath := split[len(split)-1]
	propertyPath = strings.ReplaceAll(propertyPath, "[*]", "[_]")
	return fmt.Sprintf("r.values.properties.%s", propertyPath), nil
}

func SliceConstructor(input any) string {
	var array []string
	var res string
	switch input.(type) {
	case []interface{}:
		for _, v := range input.([]interface{}) {
			array = append(array, "\""+fmt.Sprint(v)+"\"")
		}
	case []string:
		for _, v := range input.([]string) {
			array = append(array, "\""+fmt.Sprint(v)+"\"")
		}
	case string:
		array = append(array, fmt.Sprint(input))
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
	result = "r.values" + result

	return result
}
