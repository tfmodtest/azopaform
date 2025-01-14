package shared

import (
	"context"
)

var _ Rego = StringRego("")

type StringRego string

func (s StringRego) Rego(ctx context.Context) (string, error) {
	return string(s), nil
}
