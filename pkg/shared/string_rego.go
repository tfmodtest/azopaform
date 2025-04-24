package shared

import "strings"

var _ Rego = StringRego("")

type StringRego string

func (s StringRego) Rego(ctx *Context) (string, error) {
	return strings.Replace(string(s), "[*]", "[_]", -1), nil
}
