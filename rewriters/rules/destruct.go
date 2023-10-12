package rules

import (
	"fmt"
	"strings"

	"github.com/tranduythanh/gocoqatoo/coq"
)

type Destruct struct {
	bundle map[string]string
}

func NewDestruct(bundle map[string]string) *Destruct {
	return &Destruct{bundle: bundle}
}

func (d *Destruct) Apply(
	input *coq.Input,
	output *coq.Output,
	before, after map[coq.Assumption]struct{},
	previousOutput *coq.Output,
) string {
	destructedObject := strings.TrimSpace(input.Value[strings.Index(input.Value, " ") : len(input.Value)-1]) // Obtains the "(A B)" in "destruct (A B)."
	return fmt.Sprintf(d.bundle["destruct"], destructedObject)
}
