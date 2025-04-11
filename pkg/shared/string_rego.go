package shared

var _ Rego = StringRego("")

type StringRego string

func (s StringRego) Rego(ctx *Context) (string, error) {
	return ReplaceIndex(string(s)), nil
}
