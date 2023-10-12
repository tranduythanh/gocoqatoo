package rules

import (
	"github.com/tranduythanh/gocoqatoo/coq"
)

type Unfold struct {
	bundle map[string]string
}

func NewUnfold(bundle map[string]string) *Unfold {
	return &Unfold{bundle: bundle}
}

func (d *Unfold) Apply(
	input *coq.Input,
	output *coq.Output,
	before, after map[coq.Assumption]struct{},
	previousOutput *coq.Output,
) string {
	return ""
}
