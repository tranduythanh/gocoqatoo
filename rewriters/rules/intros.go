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
	before map[coq.Assumption]struct{},
	after map[coq.Assumption]struct{},
	previousOutput *coq.Output,
) string {
	var result, variables, hypotheses string
	typesToVariables := make(map[string][]string)

	// Helper function to determine if a map contains a specific key
	containsKey := func(m map[coq.Assumption]struct{}, key coq.Assumption) bool {
		_, exists := m[key]
		return exists
	}

	for a := range after {
		assumptionTypeIsAlsoVariableName := false
		for b := range before {
			if containsKey(before, b) || containsKey(after, b) {
				if a.Typ == b.Name {
					assumptionTypeIsAlsoVariableName = true
				}
			}
		}

		if !strings.Contains(a.Typ, " ") && !assumptionTypeIsAlsoVariableName {
			// Need to use intros.given
			typesToVariables[a.Typ] = append(typesToVariables[a.Typ], a.Name)
		} else {
			// Need to use intros.suppose
			hypotheses += fmt.Sprintf("%s, ", a.Typ)
		}
	}

	for typ, vars := range typesToVariables {
		variables += fmt.Sprintf("%s : %s, ", strings.Join(vars, ", "), typ)
	}

	if variables != "" {
		variables = strings.TrimSuffix(variables, ", ")
		result += fmt.Sprintf("intros.given: %s", variables)
	}
	if hypotheses != "" {
		hypotheses = strings.TrimSuffix(hypotheses, ", ")
		result += fmt.Sprintf("intros.suppose: %s", hypotheses)
	}
	result += fmt.Sprintf("intros.goal: %s", output.Goal.Value)

	return result
}
