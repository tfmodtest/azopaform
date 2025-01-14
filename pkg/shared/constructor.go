package shared

import (
	"context"
	"fmt"
	"github.com/emirpasic/gods/stacks"
	"github.com/magodo/aztfq/aztfq"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type LookupTable aztfq.LookupTable

func (t LookupTable) QueryProperty(resourceType, apiVersion, propertyAddress string) ([]aztfq.TFResult, bool) {
	m, ok := t.QueryResource(resourceType, apiVersion)
	if !ok {
		return nil, false
	}
	r, ok := m[propertyAddress]
	return r, ok
}

func (t LookupTable) QueryParentProperty(resourceType, apiVersion, propertyAddress string) string {
	var result string
	m, ok := t.QueryResource(resourceType, apiVersion)
	if !ok {
		return ""
	}
	_, ok = m[propertyAddress]
	if !ok {
		for k, v := range m {
			if strings.HasPrefix(k, propertyAddress) {
				childAddr := v[0].PropertyAddr
				addrArray := strings.Split(childAddr, "/")
				for i := len(addrArray) - 1; i >= 0; i-- {
					if _, err := strconv.Atoi(addrArray[i]); err == nil {
						continue
					}
					result = strings.Join(addrArray[:i], "/")
					break
				}
			}
		}
	}
	return result
}

func (t LookupTable) QueryResource(resourceType, apiVersion string) (map[string][]aztfq.TFResult, bool) {
	l2, ok := t[resourceType]
	if !ok {
		return nil, false
	}
	l3, ok := l2[apiVersion]
	if !ok {
		return nil, false
	}
	return l3, true
}

var ResourceTypeLookupTable = func() LookupTable {
	b, err := os.ReadFile("../output.json")
	if err != nil {
		panic(err.Error())
	}
	t, err := aztfq.BuildLookupTable(b, nil)
	if err != nil {
		panic(err.Error())
	}
	return LookupTable(t)
}()

func FieldNameProcessor(fieldName interface{}, ctx context.Context) (string, string, error) {
	var result string
	var rules string
	switch fn := fieldName.(type) {
	case string:
		if fn == TypeOfResource || fn == KindOfResource {
			return strings.Join([]string{"r.", fn}, ""), "", nil
		}
		if strings.Contains(fn, "count") {
			return fn, "", nil
		}
		rt, err := currentResourceType(ctx)
		if err != nil {
			return "", "", err
		}
		res, err := FieldNameParser(fn, rt, "")
		if err != nil {
			return "", "", err
		}
		result = TFNameMapping(res)
	}

	return result, rules, nil
}

func SliceConstructor(input any) string {
	var array []string
	var res string
	//fmt.Printf("the input type is %+v\n", reflect.TypeOf(input))
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

func currentResourceType(ctx context.Context) (string, error) {
	resourceTypeStack := ctx.Value("context").(map[string]stacks.Stack)["resourceType"]
	if resourceTypeStack == nil {
		return "", fmt.Errorf("cannot find the resource type in the context")
	}
	resourceType, ok := resourceTypeStack.Peek()
	if !ok {
		return "", fmt.Errorf("cannot find the resource type in the context")
	}
	rt, ok := resourceType.(string)
	if !ok {
		return "", fmt.Errorf("cannot convert the resource type to string")
	}
	return rt, nil
}

func FieldNameParser(fieldNameRaw, resourceType, version string) (string, error) {
	if fieldNameRaw == TypeOfResource {
		return fieldNameRaw, nil
	}
	//if strings.Contains(fieldNameRaw, "count") {
	//	return fieldNameRaw, nil
	//}
	if strings.HasPrefix(strings.ToLower(fieldNameRaw), strings.ToLower(resourceType)) {
		rtLen := len(resourceType)
		fieldNameRaw = fieldNameRaw[rtLen:]
	}
	//some attributes has "properties/" in the middle of the path after the list name, need to address this case
	prop := fieldNameRaw
	prop = strings.Replace(prop, ".", "/", -1)
	prop = strings.Replace(prop, "[x]", "/*", -1)
	prop = strings.Replace(prop, "[*]", "/*", -1)
	prop = strings.TrimPrefix(prop, "/")
	//fmt.Printf("the prop is %s\n", prop)
	upperRt := strings.ToUpper(resourceType)
	if results, ok := ResourceTypeLookupTable.QueryProperty(upperRt, version, prop); ok {
		return results[0].PropertyAddr, nil
	}
	prop = "properties/" + prop
	if results, ok := ResourceTypeLookupTable.QueryProperty(upperRt, version, prop); ok {
		return results[0].PropertyAddr, nil
	}
	prop = strings.Replace(prop, "*/", "*/properties/", -1)
	if results, ok := ResourceTypeLookupTable.QueryProperty(upperRt, version, prop); ok {
		return results[0].PropertyAddr, nil
	}

	parentPropAddr := ResourceTypeLookupTable.QueryParentProperty(upperRt, version, prop)
	if parentPropAddr != "" {
		return parentPropAddr, nil
	}

	prop = strings.Replace(prop, "properties/", "", -1)
	prop = ToSnakeCase(prop)
	return prop, nil
}

func ToSnakeCase(str string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
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
	result = "r.change.after" + result

	return result
}
