package rules

import (
	"github.com/tranduythanh/gocoqatoo/coq"
)

type Simpl struct {
	bundle map[string]string
}

func NewSimpl(bundle map[string]string) *Simpl {
	return &Simpl{bundle: bundle}
}

func (d *Simpl) Apply(
	input *coq.Input,
	output *coq.Output,
	before, after map[coq.Assumption]struct{},
	previousOutput *coq.Output,
) string {
	return ""
}
