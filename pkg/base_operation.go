package pkg

type baseOperation struct {
	conditionSetName string
}

func (o baseOperation) GetConditionSetName() string {
	return o.conditionSetName
}
