package operation

import (
	"github.com/tfmodtest/azopaform/pkg/condition"
	"strings"

	"github.com/tfmodtest/azopaform/pkg/shared"
)

var _ shared.Rego = Count{}

type Count struct {
	Where          Operation
	CountExp       string
	Subject        shared.Rego
	countFieldName string
}

func NewCount(input any, ctx *shared.Context) (Count, error) {
	items := input.(map[string]any)
	subject, err := tryParseSubject(items, ctx)
	if err != nil {
		return Count{}, err
	}
	var countFieldName string
	if field, ok := subject.(condition.FieldValue); ok {
		ctx.EnterCountRego(field.Name)
		defer ctx.ExitCountRego()
		countFieldName = field.Name
		field.Name = strings.ReplaceAll(field.Name, "[*]", "[_]")
		subject = field
	}
	var whereBody Operation
	if where, ok := items[shared.Where].(map[string]any); ok {
		whereBody, err = NewWhere(where, ctx)
		if err != nil {
			return Count{}, err
		}
	}
	countFieldConverted, err := subject.Rego(ctx)
	if err != nil {
		return Count{}, err
	}
	var countBody string
	if whereBody != nil {
		countBody = shared.Count + "({x|x:=" + countFieldConverted + ";" + whereBody.HelperFunctionName() + "(r, x)})"
	} else {
		countBody = shared.Count + "({x|x:=" + countFieldConverted + "})"
	}
	countBody = strings.ReplaceAll(countBody, "[*]", "[_]")
	return Count{
		Where:          whereBody,
		CountExp:       countBody,
		Subject:        subject,
		countFieldName: countFieldName,
	}, nil
}

func (c Count) Rego(ctx *shared.Context) (string, error) {
	ctx.EnterCountRego(c.countFieldName)
	defer ctx.ExitCountRego()
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
