package coq

type Goal struct {
	Value string
}

func NewGoal(value string) *Goal {
	return &Goal{Value: value}
}
