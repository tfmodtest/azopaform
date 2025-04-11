package operation

import (
	"json-rule-finder/pkg/shared"
	"strings"
)

var _ shared.Rego = Count{}

type Count struct {
	Where    Operation
	CountExp string
}

func NewCount(input any, ctx *shared.Context) Count {
	items := input.(map[string]any)

	var whereBody Operation
	if where, ok := items[shared.Where].(map[string]any); ok {
		whereBody = NewWhere(where, ctx)
	}
	fieldName := items[shared.Field]
	if items[shared.Field] == nil {
		fieldName = items[shared.Value]
	}
	countField, err := shared.FieldNameProcessor(fieldName.(string), ctx)
	if err != nil {
		countField = items[shared.Field].(string)
	}
	countFieldConverted := replaceIndex(countField)
	var countBody string
	if whereBody != nil {
		countBody = shared.Count + "({x|x:=" + countFieldConverted + "[_];" + whereBody.GetConditionSetName() + "(x)})"
	} else {
		countBody = shared.Count + "({x|x:=" + countFieldConverted + "[_]})"
	}
	countBody = strings.Replace(countBody, "*", "x", -1)
	return Count{
		Where:    whereBody,
		CountExp: countBody,
	}
}

func (c Count) Rego(ctx *shared.Context) (string, error) {
	res := c.CountExp
	if c.Where != nil {
		whereSubset, err := c.Where.Rego(ctx)
		if err != nil {
			return "", err
		}
		res = c.CountExp + "\n" + whereSubset
	}
	return res, nil
}
