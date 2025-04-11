package condition

import (
	"fmt"
	"json-rule-finder/pkg/shared"
	"strings"
)

var _ Condition = Like{}

type Like struct {
	BaseCondition
	Value string
}

func (l Like) Rego(ctx *shared.Context) (string, error) {
	fieldName, err := l.GetSubject(ctx).Rego(ctx)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{shared.RegexExp, "(", "\"", fmt.Sprintf(l.Value), "\"", ",", "\"", fieldName, "\"", ")"}, ""), nil
}
