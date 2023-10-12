package rules

import (
	"github.com/tranduythanh/gocoqatoo/coq"
)

type Reflexitivity struct {
	bundle map[string]string
}

func NewReflexitivity(bundle map[string]string) *Reflexitivity {
	return &Reflexitivity{bundle: bundle}
}

func (d *Reflexitivity) Apply(
	input *coq.Input,
	output *coq.Output,
	before, after map[coq.Assumption]struct{},
	previousOutput *coq.Output,
) string {
	return ""
}
