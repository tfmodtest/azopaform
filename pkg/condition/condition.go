package condition

import "json-rule-finder/pkg/shared"

type Condition interface {
	shared.Rego
	GetSubject(*shared.Context) shared.Rego
}
