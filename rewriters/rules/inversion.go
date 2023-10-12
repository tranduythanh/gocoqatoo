package rules

import (
	"fmt"
	"strings"

	"github.com/tranduythanh/gocoqatoo/coq"
)

type Inversion struct {
	bundle map[string]string
}

func NewInversion(bundle map[string]string) *Inversion {
	return &Inversion{bundle: bundle}
}

func (i *Inversion) Apply(
	input *coq.Input,
	output *coq.Output,
	before, after map[coq.Assumption]struct{},
	previousOutput *coq.Output,
) string {
	inversionLemmaName := strings.Replace(strings.Split(input.Value, " ")[1], ".", "", -1)
	var inversionLemmaDefinition string

	for a := range before {
		if a.Name == inversionLemmaName {
			inversionLemmaDefinition = a.Typ
			break
		}
	}

	var enumerationOfAddedAssumptions []string
	for a := range after {
		if !strings.Contains(a.Typ, " ") {
			enumerationOfAddedAssumptions = append(enumerationOfAddedAssumptions, a.Typ)
		}
	}
	enumerationStr := strings.Join(enumerationOfAddedAssumptions, ", ")

	return fmt.Sprintf("inversion: %s, %s", inversionLemmaDefinition, enumerationStr)
}
