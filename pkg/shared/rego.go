package shared

import "context"

type Rego interface {
	Rego(ctx context.Context) (string, error)
}
