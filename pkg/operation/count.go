package operation

import (
	"github.com/tfmodtest/azopaform/pkg/shared"
	"strings"
)

var _ shared.Rego = Count{}

type Count struct {
	Where    Operation
	CountExp string
}

func NewCount(input any, ctx *shared.Context) (Count, error) {
	items := input.(map[string]any)

	var whereBody Operation
	var err error
	if where, ok := items[shared.Where].(map[string]any); ok {
		whereBody, err = NewWhere(where, ctx)
		if err != nil {
			return Count{}, err
		}
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
		countBody = shared.Count + "({x|x:=" + countFieldConverted + ";" + whereBody.HelperFunctionName() + "(x)})"
	} else {
		countBody = shared.Count + "({x|x:=" + countFieldConverted + "})"
	}
	countBody = strings.Replace(countBody, "*", "x", -1)
	return Count{
		Where:    whereBody,
		CountExp: countBody,
	}, nil
}

func (c Count) Rego(ctx *shared.Context) (string, error) {
	res := c.CountExp
	if c.Where != nil {
		whereHelperDef, err := c.Where.Rego(ctx)
		if err != nil {
			return "", err
		}
		ctx.EnqueueHelperFunction(whereHelperDef)
	}
	return res, nil
}
