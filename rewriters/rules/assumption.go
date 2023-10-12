package rules

import "github.com/tranduythanh/gocoqatoo/coq"

type Assumption struct {
	bundle map[string]string
}

func NewAssumption(bundle map[string]string) *Assumption {
	return &Assumption{bundle: bundle}
}

func (a *Assumption) Apply(
	input *coq.Input,
	output *coq.Output,
	before, after map[coq.Assumption]struct{},
	previousOutput *coq.Output,
) string {
	return a.bundle["assumption"]
}
