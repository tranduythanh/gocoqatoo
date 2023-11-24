package rules

import (
	"fmt"
	"strings"

	"github.com/tranduythanh/gocoqatoo/coq"
)

type Intros struct {
	bundle map[string]string
}

func NewIntros(bundle map[string]string) *Intros {
	return &Intros{bundle: bundle}
}

func (i *Intros) Apply(
	input *coq.Input,
	output *coq.Output,
	before map[string]*coq.Assumption,
	after map[string]*coq.Assumption,
	previousOutput *coq.Output,
) string {
	var result, variables, hypotheses string
	typesToVariables := make(map[string][]string)

	// Helper function to determine if a map contains a specific key
	containsKey := func(m map[string]*coq.Assumption, key string) bool {
		_, exists := m[key]
		return exists
	}

	for a := range after {
		assumptionTypeIsAlsoVariableName := false
		for b := range before {
			if containsKey(before, b) || containsKey(after, b) {
				if after[a].Typ == before[b].Name {
					assumptionTypeIsAlsoVariableName = true
				}
			}
		}

		if !strings.Contains(after[a].Typ, " ") && !assumptionTypeIsAlsoVariableName {
			// Need to use intros.given
			typesToVariables[after[a].Typ] = append(typesToVariables[after[a].Typ], after[a].Name)
		} else {
			// Need to use intros.suppose
			hypotheses += fmt.Sprintf("%s, ", after[a].Typ)
		}
	}

	for typ, vars := range typesToVariables {
		variables += fmt.Sprintf("%s : %s, ", strings.Join(vars, ", "), typ)
	}

	if variables != "" {
		variables = strings.TrimSuffix(variables, ", ")

		result += fmt.Sprintf("%s:\t%s", i.bundle["intros.given"], variables)
	}
	if hypotheses != "" {
		hypotheses = strings.TrimSuffix(hypotheses, ", ")
		result += fmt.Sprintf("%s:\t%s", i.bundle["intros.suppose"], hypotheses)
	}
	result += fmt.Sprintf("%s:\t%s", i.bundle["intros.goal"], output.Goal.Value)

	return result
}
