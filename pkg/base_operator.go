package pkg

type baseOperator struct {
	conditionSetName string
}

func (o baseOperator) GetConditionSetName() string {
	return o.conditionSetName
}
