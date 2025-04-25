package condition

import (
	"github.com/tfmodtest/azopaform/pkg/shared"
)

type Condition interface {
	shared.Rego
	GetSubject(*shared.Context) shared.Rego
}
