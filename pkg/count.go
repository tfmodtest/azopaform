package pkg

import (
	"json-rule-finder/pkg/shared"
	"strings"
)

var _ shared.Rego = &CountOperator{}

type CountOperator struct {
	Where    shared.Rego
	CountExp string
}

func NewCountOperator(input any, ctx *shared.Context) shared.Rego {
	items := input.(map[string]any)
	var whereBody shared.Rego
	if items[shared.Where] != nil {
		whereMap := items[shared.Where].(map[string]any)
		of := otherFactories[shared.Where]
		whereBody = of(whereMap, ctx)
	}
	fieldName := items[shared.Field]
	if items[shared.Field] == nil {
		fieldName = items[shared.Value]
	}
	countField, _, err := shared.FieldNameProcessor(fieldName.(string), ctx)
	if err != nil {
		countField = items[shared.Field].(string)
	}
	countFieldConverted := replaceIndex(countField)
	var countBody string
	if whereBody != nil {
		countBody = shared.Count + "(" + "{" + "x" + "|" + countFieldConverted + ";" + whereBody.(WhereOperator).ConditionSetName + "(x)" + "}" + ")"
	} else {
		countBody = shared.Count + "(" + "{" + "x" + "|" + countFieldConverted + "}" + ")"
	}
	countBody = strings.Replace(countBody, "*", "x", -1)
	return CountOperator{
		Where:    whereBody,
		CountExp: countBody,
	}
}

func (c CountOperator) Rego(ctx *shared.Context) (string, error) {
	var res string
	whereSubset, err := c.Where.Rego(ctx)
	if err != nil {
		return "", err
	}
	res = c.CountExp + "\n" + whereSubset
	return res, nil
}

var _ shared.Rego = &Count{}

type Count struct {
	Count        string
	ConditionSet shared.Rego
}

func (c Count) Rego(ctx *shared.Context) (string, error) {
	return c.Count, nil
}
