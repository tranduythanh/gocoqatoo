package rules

import (
	"fmt"
	"strings"

	"github.com/tranduythanh/gocoqatoo/coq"
)

type Apply struct {
	bundle map[string]string
}

func NewApply(bundle map[string]string) *Apply {
	return &Apply{bundle: bundle}
}

func (a *Apply) Apply(
	input *coq.Input,
	output *coq.Output,
	before, after map[string]*coq.Assumption,
	previousOutput *coq.Output,
) string {
	lemmaName := strings.Replace(strings.Split(input.Value, " ")[1], ".", "", -1)
	var lemmaDefinition string

	for b := range before {
		if before[b].Name == lemmaName {
			lemmaDefinition = before[b].Typ
			break
		}
	}

	indexOfLastImplication := strings.LastIndex(lemmaDefinition, "->")
	if indexOfLastImplication != -1 {
		propositionsLeftToProve := strings.Split(lemmaDefinition[:indexOfLastImplication], "->")
		for i, proposition := range propositionsLeftToProve {
			propositionsLeftToProve[i] = strings.TrimSpace(proposition)
		}
		commaSeparatedPropositionsLeftToProve := strings.Join(propositionsLeftToProve, ", ")
		return fmt.Sprintf("apply: %s, %s, %s", lemmaDefinition, previousOutput.Goal.Value, commaSeparatedPropositionsLeftToProve)
	}

	return ""
}
