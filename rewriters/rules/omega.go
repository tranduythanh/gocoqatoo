package rules

import (
	"github.com/tranduythanh/gocoqatoo/coq"
)

type Omega struct {
	bundle map[string]string
}

func NewOmega(bundle map[string]string) *Omega {
	return &Omega{bundle: bundle}
}

func (d *Omega) Apply(
	input *coq.Input,
	output *coq.Output,
	before, after map[coq.Assumption]struct{},
	previousOutput *coq.Output,
) string {
	return ""
}
