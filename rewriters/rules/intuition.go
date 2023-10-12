package rules

import (
	"github.com/tranduythanh/gocoqatoo/coq"
)

type Intuition struct {
	bundle map[string]string
}

func NewIntuition(bundle map[string]string) *Intuition {
	return &Intuition{bundle: bundle}
}

func (d *Intuition) Apply(
	input *coq.Input,
	output *coq.Output,
	before, after map[coq.Assumption]struct{},
	previousOutput *coq.Output,
) string {
	return ""
}
