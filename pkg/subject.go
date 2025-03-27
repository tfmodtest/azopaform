package pkg

import (
	"fmt"
	"json-rule-finder/pkg/shared"
	"json-rule-finder/pkg/value"
)

func NewSubject(subjectKey string, body any, ctx *shared.Context) shared.Rego {
	switch subjectKey {
	case shared.Field:
		return value.NewFieldValue(body, ctx)
	case shared.Value:
		return value.NewValue(body, ctx)
	case shared.Count:
		return NewCount(body, ctx)
	}
	panic(fmt.Errorf("subjectKey %s not found", subjectKey))
}
