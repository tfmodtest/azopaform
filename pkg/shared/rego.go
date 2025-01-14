package shared

type Rego interface {
	Rego(ctx *Context) (string, error)
}
