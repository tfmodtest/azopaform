package pkg

import "context"

var _ Rego = ConditionSet{}

type ConditionSet struct {
	Conditions []Rego
}

func (c ConditionSet) Rego(ctx context.Context) (string, error) {
	var res string
	for _, item := range c.Conditions {
		condition, err := item.Rego(ctx)
		if err != nil {
			return "", err
		}
		if res != "" {
			res = res + "\n"
		}
		res += condition
	}
	return res, nil
}
